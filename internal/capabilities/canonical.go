package capabilities

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// CanonicalMarkdown applies the bounded pilot implementation of the §4.1.2
// canonical profile: LF endings, exactly one terminal LF, insignificant trailing
// whitespace removed, blank-line runs collapsed, and table rows normalized to
// declared pipe/separator syntax. Whitespace inside fenced code blocks is
// preserved exactly. Indented literal blocks and alignment-colon preservation are
// outside this bounded v1.
func CanonicalMarkdown(payload []byte) ([]byte, error) {
	text := strings.ReplaceAll(string(payload), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	var out []string
	inFence := false
	fenceMarker := ""
	blankRun := 0
	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if !inFence && (strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")) {
			inFence = true
			fenceMarker = trimmed[:3]
			blankRun = 0
			out = append(out, strings.TrimRight(line, " \t"))
			i++
			continue
		}
		if inFence {
			out = append(out, line)
			if strings.HasPrefix(trimmed, fenceMarker) {
				inFence = false
			}
			i++
			continue
		}
		if trimmed == "" {
			blankRun++
			i++
			continue
		}
		if blankRun > 0 && len(out) > 0 {
			out = append(out, "")
		}
		blankRun = 0
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			for i < len(lines) {
				row := strings.TrimSpace(lines[i])
				if !strings.HasPrefix(row, "|") || !strings.HasSuffix(row, "|") {
					break
				}
				out = append(out, canonicalTableRow(row))
				i++
			}
			continue
		}
		out = append(out, strings.TrimRight(line, " \t"))
		i++
	}
	if len(out) == 0 {
		return []byte("\n"), nil
	}
	return []byte(strings.Join(out, "\n") + "\n"), nil
}

func canonicalTableRow(row string) string {
	rawCells := strings.Split(row, "|")
	if len(rawCells) < 2 {
		return row
	}
	cells := rawCells[1 : len(rawCells)-1]
	normalized := make([]string, 0, len(cells))
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		cell = strings.Trim(cell, ":")
		cell = strings.TrimSpace(cell)
		if cell != "" && strings.Trim(cell, "-") == "" {
			cell = "---"
		}
		normalized = append(normalized, cell)
	}
	return "| " + strings.Join(normalized, " | ") + " |"
}

// SemanticContentFingerprint canonicalizes a semantic payload and excludes the
// named carrier field line before hashing.
func SemanticContentFingerprint(payload []byte, carrier string) (string, error) {
	if carrier == "" {
		return "", fmt.Errorf("canonical: empty carrier name")
	}
	canonical, err := CanonicalMarkdown(payload)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSuffix(string(canonical), "\n"), "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **"+carrier+":**") || strings.HasPrefix(trimmed, "- **"+carrier+"**") {
			continue
		}
		kept = append(kept, line)
	}
	return Hash([]byte(strings.Join(kept, "\n") + "\n")), nil
}

// ValidEnvelopeKind reports whether the kind is a declared envelope kind.
func ValidEnvelopeKind(kind string) bool {
	switch kind {
	case "CERTIFICATION", "DECISION-APPROVAL", "CNF-DIGEST-SIGNATURE":
		return true
	}
	return false
}

// EnvelopeFingerprint computes the canonical envelope fingerprint. Preimage:
// "ENVELOPE|" + path + "|" + targetID + "|" + kind + "|" + version + LF +
// canonical region.
func EnvelopeFingerprint(path, targetID, kind string, version int, region []byte) (string, error) {
	if path == "" || targetID == "" {
		return "", fmt.Errorf("canonical: envelope path and target ID are required")
	}
	if !ValidEnvelopeKind(kind) {
		return "", fmt.Errorf("canonical: unknown envelope kind %q", kind)
	}
	if version < 1 {
		return "", fmt.Errorf("canonical: envelope version must be positive")
	}
	canonicalRegion, err := CanonicalMarkdown(region)
	if err != nil {
		return "", err
	}
	preimage := "ENVELOPE|" + path + "|" + targetID + "|" + kind + "|" + strconv.Itoa(version) + "\n" + string(canonicalRegion)
	return Hash([]byte(preimage)), nil
}

// PackageMemberFingerprint computes the complete path-bound file fingerprint.
func PackageMemberFingerprint(path string, file []byte) (string, error) {
	if path == "" {
		return "", fmt.Errorf("canonical: package-member path is required")
	}
	canonical, err := CanonicalMarkdown(file)
	if err != nil {
		return "", err
	}
	return Hash([]byte("PACKAGE-MEMBER|" + path + "\n" + string(canonical))), nil
}

// TransportBinding is one containing-file transport tuple.
type TransportBinding struct {
	RecordType  string
	RecordID    string
	Version     string
	Fingerprint string
}

// TransportFingerprint computes the ordered containing-file transport fingerprint.
func TransportFingerprint(path string, bindings []TransportBinding) (string, error) {
	if path == "" {
		return "", fmt.Errorf("canonical: transport path is required")
	}
	var b strings.Builder
	b.WriteString("FILE-TRANSPORT|" + path + "\n")
	for _, binding := range bindings {
		b.WriteString(binding.RecordType + "|" + binding.RecordID + "|" + binding.Version + "|" + binding.Fingerprint + "\n")
	}
	return Hash([]byte(b.String())), nil
}

// EvidenceBinding is one evidence-set tuple.
type EvidenceBinding struct {
	InputPath        string
	Kind             string
	RecordOrArtifact string
	Version          string
	Fingerprint      string
}

var evidenceBindingKinds = map[string]bool{
	"RAW-SOURCE-BYTES": true, "RECORD-PAYLOAD": true, "SEMANTIC-CONTENT": true,
	"DECISION-CONTENT": true, "CNF-PAYLOAD": true, "ARTIFACT-PAYLOAD": true,
	"CANONICAL-ENVELOPE": true, "FILE-TRANSPORT": true, "PACKAGE-MEMBER-FILE": true,
	"VALIDATION-SUMMARY": true, "DENOMINATOR-QUERY": true, "PROTOCOL-REGISTRY": true,
}

var anchorBindingKinds = map[string]bool{
	"RAW-SOURCE-BYTES": true, "RECORD-PAYLOAD": true, "SEMANTIC-CONTENT": true,
	"DECISION-CONTENT": true, "CNF-PAYLOAD": true, "ARTIFACT-PAYLOAD": true,
	"DENOMINATOR-QUERY": true, "PROTOCOL-REGISTRY": true,
}

// EvidenceSetFingerprint computes the §4.1.2 evidence-set fingerprint. Preimage:
// "EVIDENCE-SET|" + artifact type + LF + row kind + LF + row ID + LF + System
// Protocol Snapshot + LF + canonical table rows (sorted bytewise by complete
// tuple, columns in tuple order).
func EvidenceSetFingerprint(artifactType, rowKind, rowID, systemProtocolSnapshot string, bindings []EvidenceBinding) (string, error) {
	if artifactType == "" || rowID == "" {
		return "", fmt.Errorf("canonical: evidence-set artifact type and row ID are required")
	}
	if rowKind != "CHECK" && rowKind != "BLOCKER" {
		return "", fmt.Errorf("canonical: evidence-set row kind %q is not CHECK or BLOCKER", rowKind)
	}
	type tuple struct {
		key    string
		fields [5]string
	}
	tuples := make([]tuple, 0, len(bindings))
	seen := map[string]bool{}
	anchored := false
	for _, binding := range bindings {
		if !evidenceBindingKinds[binding.Kind] {
			return "", fmt.Errorf("canonical: unknown evidence binding kind %q", binding.Kind)
		}
		t := tuple{fields: [5]string{binding.InputPath, binding.Kind, binding.RecordOrArtifact, binding.Version, binding.Fingerprint}}
		t.key = strings.Join(t.fields[:], "|")
		if seen[t.key] {
			return "", fmt.Errorf("canonical: duplicate evidence binding tuple")
		}
		seen[t.key] = true
		tuples = append(tuples, t)
		if anchorBindingKinds[binding.Kind] {
			anchored = true
		}
	}
	if !anchored {
		return "", fmt.Errorf("canonical: evidence set lacks a denominator, registry, or record/artifact binding")
	}
	sort.Slice(tuples, func(i, j int) bool { return tuples[i].key < tuples[j].key })
	rows := make([]string, 0, len(tuples))
	for _, t := range tuples {
		rows = append(rows, "| "+strings.Join(t.fields[:], " | ")+" |")
	}
	preimage := "EVIDENCE-SET|" + artifactType + "\n" + rowKind + "\n" + rowID + "\n" + systemProtocolSnapshot + "\n" + strings.Join(rows, "\n")
	return Hash([]byte(preimage)), nil
}

// HashDomainIdentity returns the artifact-type plus normalized-path identity.
func HashDomainIdentity(artifactType, path string) (string, error) {
	if artifactType == "" || path == "" {
		return "", fmt.Errorf("canonical: hash-domain artifact type and path are required")
	}
	return artifactType + "|" + path, nil
}

// PostHashArtifactInstanceBinding returns the post-hash artifact instance binding.
func PostHashArtifactInstanceBinding(artifactType, path, payloadFingerprint string) (string, error) {
	domain, err := HashDomainIdentity(artifactType, path)
	if err != nil {
		return "", err
	}
	if payloadFingerprint == "" {
		return "", fmt.Errorf("canonical: payload fingerprint is required")
	}
	return domain + "|" + payloadFingerprint, nil
}
