package artifact

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// ArtifactSchemaVersion is the packaging artifact schema version carried in the
// `PROTOCOL / ARTIFACT-SCHEMA VERSION` payload field.
const ArtifactSchemaVersion = "artifact-schema/1"

// EvidenceBinding is one §12.5 evidence-set tuple.
type EvidenceBinding struct {
	Path        string
	Kind        string
	RecordID    string
	Version     string
	Fingerprint string
}

var evidenceKinds = map[string]bool{
	"RAW-SOURCE-BYTES":    true,
	"RECORD-PAYLOAD":      true,
	"SEMANTIC-CONTENT":    true,
	"DECISION-CONTENT":    true,
	"CNF-PAYLOAD":         true,
	"ARTIFACT-PAYLOAD":    true,
	"CANONICAL-ENVELOPE":  true,
	"FILE-TRANSPORT":      true,
	"PACKAGE-MEMBER-FILE": true,
	"VALIDATION-SUMMARY":  true,
	"DENOMINATOR-QUERY":   true,
	"PROTOCOL-REGISTRY":   true,
}

// EvidenceSetFingerprint computes the §12.5 evidence-set fingerprint for one
// check or blocker row. The preimage is `EVIDENCE-SET|`, the artifact type, the
// row kind, the row ID, the canonical System / Protocol / Snapshot value, and
// the canonically serialized sorted binding table.
func EvidenceSetFingerprint(artifactType, rowKind, rowID, snapshot string, bindings []EvidenceBinding) (string, error) {
	if artifactType == "" || rowKind == "" || rowID == "" || snapshot == "" {
		return "", fmt.Errorf("evidence: artifact type, row kind, row ID, and snapshot are required")
	}
	if len(bindings) == 0 {
		return "", fmt.Errorf("evidence: row %s has no evidence bindings", rowID)
	}
	sorted := append([]EvidenceBinding(nil), bindings...)
	for _, binding := range sorted {
		if binding.Path == "" || binding.Fingerprint == "" {
			return "", fmt.Errorf("evidence: row %s has a binding without path or fingerprint", rowID)
		}
		if !evidenceKinds[binding.Kind] {
			return "", fmt.Errorf("evidence: row %s has unknown binding kind %q", rowID, binding.Kind)
		}
		if binding.RecordID == "" {
			return "", fmt.Errorf("evidence: row %s has a binding without a record/artifact ID", rowID)
		}
		if binding.Version == "" {
			return "", fmt.Errorf("evidence: row %s has a binding without a version", rowID)
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return evidenceTuple(sorted[i]) < evidenceTuple(sorted[j])
	})
	for i := 1; i < len(sorted); i++ {
		if evidenceTuple(sorted[i]) == evidenceTuple(sorted[i-1]) {
			return "", fmt.Errorf("evidence: row %s repeats a binding tuple", rowID)
		}
	}
	var b strings.Builder
	b.WriteString("| Normalized Authoritative Input Path | Binding Kind | Typed Record or Artifact ID | Version | Authoritative Fingerprint |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, binding := range sorted {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", binding.Path, binding.Kind, binding.RecordID, binding.Version, binding.Fingerprint)
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(b.String()))
	if err != nil {
		return "", err
	}
	preimage := "EVIDENCE-SET|" + artifactType + "\n" + rowKind + "\n" + rowID + "\n" + snapshot + "\n" + string(canonical)
	return capabilities.Hash([]byte(preimage)), nil
}

// BindingTableFingerprint hashes a labelled, canonically sorted binding table
// for input-set projections that are not check or blocker rows.
func BindingTableFingerprint(domain, snapshot string, bindings []EvidenceBinding) (string, error) {
	if domain == "" || len(bindings) == 0 {
		return "", fmt.Errorf("evidence: binding table requires a domain and bindings")
	}
	sorted := append([]EvidenceBinding(nil), bindings...)
	for _, binding := range sorted {
		if binding.Path == "" || binding.Fingerprint == "" {
			return "", fmt.Errorf("evidence: binding table has a binding without path or fingerprint")
		}
		if !evidenceKinds[binding.Kind] {
			return "", fmt.Errorf("evidence: binding table has unknown binding kind %q", binding.Kind)
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return evidenceTuple(sorted[i]) < evidenceTuple(sorted[j])
	})
	var b strings.Builder
	b.WriteString("| Normalized Authoritative Input Path | Binding Kind | Typed Record or Artifact ID | Version | Authoritative Fingerprint |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, binding := range sorted {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", binding.Path, binding.Kind, binding.RecordID, binding.Version, binding.Fingerprint)
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(b.String()))
	if err != nil {
		return "", err
	}
	return capabilities.Hash([]byte(domain + "|" + snapshot + "\n" + string(canonical))), nil
}

func evidenceTuple(binding EvidenceBinding) string {
	return strings.Join([]string{binding.Path, binding.Kind, binding.RecordID, binding.Version, binding.Fingerprint}, "|")
}

// ValidationSummaryFingerprint computes a deterministic validator-output
// fingerprint over a snapshot and sorted key/value rows.
func ValidationSummaryFingerprint(validatorID, snapshot string, rows [][2]string) string {
	sorted := append([][2]string(nil), rows...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i][0] != sorted[j][0] {
			return sorted[i][0] < sorted[j][0]
		}
		return sorted[i][1] < sorted[j][1]
	})
	var b strings.Builder
	b.WriteString("| Key | Value |\n| --- | --- |\n")
	for _, row := range sorted {
		fmt.Fprintf(&b, "| %s | %s |\n", row[0], row[1])
	}
	canonical, _ := capabilities.CanonicalMarkdown([]byte(b.String()))
	return capabilities.Hash([]byte("VALIDATION-SUMMARY|" + validatorID + "\n" + snapshot + "\n" + string(canonical)))
}

// CanonicalSetFingerprint computes a deterministic fingerprint over a labelled
// sorted value set (input-set and similar projections).
func CanonicalSetFingerprint(domain, snapshot string, values []string) string {
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	var b strings.Builder
	b.WriteString("| Value |\n| --- |\n")
	for _, value := range sorted {
		fmt.Fprintf(&b, "| %s |\n", value)
	}
	canonical, _ := capabilities.CanonicalMarkdown([]byte(b.String()))
	return capabilities.Hash([]byte(domain + "|" + snapshot + "\n" + string(canonical)))
}

// FileTransportFingerprint computes the durable fingerprint of one transport
// file: domain `FILE-TRANSPORT|`, normalized path, and the canonicalized bytes.
func FileTransportFingerprint(path string, content []byte) (string, error) {
	canonical, err := capabilities.CanonicalMarkdown(content)
	if err != nil {
		return "", err
	}
	return capabilities.Hash([]byte("FILE-TRANSPORT|" + path + "\n" + string(canonical))), nil
}

// CheckRequirements extracts the normative requirement text for every report
// check ID from the declared registry table in the compiled protocol sources.
func CheckRequirements(doc *cdl.EIRDoc) map[string]string {
	requirements := map[string]string{}
	for _, section := range doc.Sections {
		for _, line := range strings.Split(section.Goal, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "|") {
				continue
			}
			cells := splitRow(line)
			if len(cells) != 3 {
				continue
			}
			scope := strings.TrimSpace(cells[0])
			id := strings.TrimSpace(cells[1])
			if scope != "CANDIDATE" && scope != "CONTENT-READINESS" && scope != "FINAL" {
				continue
			}
			if id == "" || id == "Check ID" {
				continue
			}
			requirements[id] = strings.TrimSpace(cells[2])
		}
	}
	return requirements
}
