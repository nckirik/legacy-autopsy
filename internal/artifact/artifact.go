// Package artifact validates, fingerprints, packages, and verifies protocol
// packaging artifacts. Executors author payload Markdown; this package owns the
// deterministic side: declared field sequence and order, canonical payload
// fingerprints, artifact-instance bindings, package-member fingerprints, and
// fail-closed bundle verification.
package artifact

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// Carrier is the artifact payload fingerprint carrier field label.
const Carrier = "ARTIFACT-PAYLOAD-FINGERPRINT"

// FieldSequence returns the declared field names of an artifact FIELD in order.
func FieldSequence(doc *cdl.EIRDoc, fieldID string) ([]string, error) {
	field, err := findArtifactField(doc, fieldID)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(field.Fields))
	for i, spec := range field.Fields {
		names[i] = spec.Name
	}
	return names, nil
}

// ValidatePayload checks an executor-authored artifact payload against its
// declared FIELD schema: every payload label must match a declared field through
// the label-to-column mapping, labels must appear in declared relative order,
// duplicates are rejected, and required fields must be present.
func ValidatePayload(doc *cdl.EIRDoc, fieldID string, payload []byte) error {
	field, err := findArtifactField(doc, fieldID)
	if err != nil {
		return err
	}
	position := map[string]int{}
	for i, spec := range field.Fields {
		position[spec.Name] = i
	}
	parsed, err := md.Parse(string(ExtractPayloadBody(payload)))
	if err != nil {
		return fmt.Errorf("artifact: parse %s payload: %w", fieldID, err)
	}
	seen := map[string]bool{}
	last := -1
	for _, node := range parsed.Nodes {
		if node.Kind != md.Field {
			continue
		}
		name := capabilities.ColumnName(node.Label)
		at, ok := position[name]
		if !ok {
			return fmt.Errorf("artifact: %s payload has unknown field %q", fieldID, node.Label)
		}
		if seen[name] {
			return fmt.Errorf("artifact: %s payload repeats field %q", fieldID, node.Label)
		}
		if at <= last {
			return fmt.Errorf("artifact: %s payload field %q is out of declared order", fieldID, node.Label)
		}
		seen[name] = true
		last = at
	}
	for _, spec := range field.Fields {
		if spec.Required && !seen[spec.Name] {
			return fmt.Errorf("artifact: %s payload is missing required field %q", fieldID, spec.Name)
		}
	}
	return nil
}

// PayloadFingerprint computes the canonical artifact-payload fingerprint. Its
// preimage is the hash-domain identity (artifact type + "|" + normalized path)
// followed by one LF and the canonical payload with the payload's own carrier
// line excluded. The returned fingerprint never includes the carrier.
func PayloadFingerprint(artifactType, path string, payload []byte) (string, error) {
	identity, err := capabilities.HashDomainIdentity(artifactType, path)
	if err != nil {
		return "", err
	}
	body := ExtractPayloadBody(payload)
	canonical, err := capabilities.CanonicalMarkdown(body)
	if err != nil {
		return "", err
	}
	withoutCarrier, err := excludeCarrier(canonical)
	if err != nil {
		return "", err
	}
	return capabilities.Hash([]byte(identity + "\n" + string(withoutCarrier))), nil
}

// ExtractPayloadBody returns the exact AST content between the generic
// artifact-payload boundary markers. Payloads without markers are returned
// unchanged (envelope-region payloads).
func ExtractPayloadBody(payload []byte) []byte {
	text := string(payload)
	start := strings.Index(text, "<!-- ARTIFACT-PAYLOAD:")
	if start < 0 {
		return payload
	}
	startEnd := strings.Index(text[start:], "-->")
	if startEnd < 0 {
		return payload
	}
	bodyStart := start + startEnd + len("-->")
	end := strings.Index(text[bodyStart:], "<!-- END-ARTIFACT-PAYLOAD")
	if end < 0 {
		return payload
	}
	return []byte(text[bodyStart : bodyStart+end])
}

// InstanceBinding returns the post-hash artifact-instance binding.
func InstanceBinding(artifactType, path, fingerprint string) (string, error) {
	return capabilities.PostHashArtifactInstanceBinding(artifactType, path, fingerprint)
}

func excludeCarrier(canonical []byte) ([]byte, error) {
	lines := strings.Split(strings.TrimSuffix(string(canonical), "\n"), "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **"+Carrier+":**") || strings.HasPrefix(trimmed, "- **"+Carrier+"**") {
			continue
		}
		kept = append(kept, line)
	}
	return []byte(strings.Join(kept, "\n") + "\n"), nil
}

func findArtifactField(doc *cdl.EIRDoc, fieldID string) (*cdl.EIRField, error) {
	for i := range doc.Declarations.Fields {
		field := &doc.Declarations.Fields[i]
		if field.ID != fieldID {
			continue
		}
		if field.Storage != "artifact" {
			return nil, fmt.Errorf("artifact: FIELD %s has STORAGE %s, not artifact", fieldID, field.Storage)
		}
		return field, nil
	}
	return nil, fmt.Errorf("artifact: unknown artifact FIELD %s", fieldID)
}

// Member is one package member with its exact finalized bytes.
type Member struct {
	Path  string
	Bytes []byte
}

// ManifestEntry binds one package-relative path to its package-member file
// fingerprint.
type ManifestEntry struct {
	Path        string
	Fingerprint string
}

// PackageMembers computes the sorted manifest for the complete member set using
// path-bound package-member fingerprints.
func PackageMembers(members []Member) ([]ManifestEntry, error) {
	entries := make([]ManifestEntry, 0, len(members))
	seen := map[string]bool{}
	for _, member := range members {
		if seen[member.Path] {
			return nil, fmt.Errorf("artifact: duplicate package member %q", member.Path)
		}
		seen[member.Path] = true
		fingerprint, err := capabilities.PackageMemberFingerprint(member.Path, member.Bytes)
		if err != nil {
			return nil, err
		}
		entries = append(entries, ManifestEntry{Path: member.Path, Fingerprint: fingerprint})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries, nil
}

// DumpMember wraps a canonical store dump as the binding's canonical
// serialization package member.
func DumpMember(dump []byte) Member {
	return Member{Path: "store/legacy-autopsy.sql", Bytes: dump}
}

// VerifyMembers recomputes every member fingerprint and fails closed on a
// missing member, an extra member, a duplicate binding, or any byte change.
func VerifyMembers(members []Member, manifest []ManifestEntry) error {
	actual, err := PackageMembers(members)
	if err != nil {
		return err
	}
	if len(actual) != len(manifest) {
		return fmt.Errorf("artifact: manifest declares %d members, bundle has %d", len(manifest), len(actual))
	}
	want := map[string]string{}
	for _, entry := range manifest {
		if _, ok := want[entry.Path]; ok {
			return fmt.Errorf("artifact: manifest repeats member %q", entry.Path)
		}
		want[entry.Path] = entry.Fingerprint
	}
	for _, entry := range actual {
		declared, ok := want[entry.Path]
		if !ok {
			return fmt.Errorf("artifact: bundle member %q is not in the manifest", entry.Path)
		}
		if declared != entry.Fingerprint {
			return fmt.Errorf("artifact: member %q fingerprint mismatch", entry.Path)
		}
		delete(want, entry.Path)
	}
	if len(want) != 0 {
		for path := range want {
			return fmt.Errorf("artifact: manifest member %q is missing from the bundle", path)
		}
	}
	return nil
}
