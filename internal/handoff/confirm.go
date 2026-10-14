package handoff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// ManifestBinding is one parsed candidate manifest semantic binding.
type ManifestBinding struct {
	RecordType  string
	RecordID    string
	Version     string
	Fingerprint string
	Path        string
}

// ParseManifestBindings reads the candidate manifest Semantic Bindings table.
func ParseManifestBindings(payload []byte) []ManifestBinding {
	rows := manifestBindings(payload)
	bindings := make([]ManifestBinding, 0, len(rows))
	for _, row := range rows {
		bindings = append(bindings, ManifestBinding{
			RecordType:  row.recordType,
			RecordID:    row.recordID,
			Version:     row.version,
			Fingerprint: row.fingerprint,
			Path:        row.path,
		})
	}
	sort.Slice(bindings, func(i, j int) bool {
		if bindings[i].RecordType != bindings[j].RecordType {
			return bindings[i].RecordType < bindings[j].RecordType
		}
		return bindings[i].RecordID < bindings[j].RecordID
	})
	return bindings
}

// BindingToken is the exact `ID@version@fingerprint` confirmation token.
func (b ManifestBinding) BindingToken() string {
	return b.RecordID + "@" + b.Version + "@" + b.Fingerprint
}

// BuildConfirmation constructs one CNF record payload for a group of exact
// candidate bindings. The machine only records a human act supplied by the
// caller; it never decides the approval.
func BuildConfirmation(namespace, group, approver, scope, limitations, timestamp, readability string, bindings []ManifestBinding, reportFingerprint, manifestFingerprint string) (string, []byte, error) {
	if strings.TrimSpace(namespace) == "" || strings.TrimSpace(group) == "" {
		return "", nil, fmt.Errorf("handoff: confirmation requires a namespace and group")
	}
	if strings.TrimSpace(approver) == "" {
		return "", nil, fmt.Errorf("handoff: confirmation requires the approving human")
	}
	if len(bindings) == 0 {
		return "", nil, fmt.Errorf("handoff: confirmation group %s binds no records", group)
	}
	switch readability {
	case "Pass", "Fail", "Not-Applicable":
	default:
		return "", nil, fmt.Errorf("handoff: confirmation readability value %q is invalid", readability)
	}
	tokens := make([]string, 0, len(bindings))
	for _, binding := range bindings {
		tokens = append(tokens, binding.BindingToken())
	}
	sort.Strings(tokens)
	var b strings.Builder
	b.WriteString("- **human-hatch-action:** Confirmation\n")
	fmt.Fprintf(&b, "- **approver-role-authority-basis:** %s\n", approver)
	fmt.Fprintf(&b, "- **review-timestamp-review-scope-limitations:** %s; scope: %s; limitations: %s\n", timestamp, scope, limitations)
	fmt.Fprintf(&b, "- **handbook-readability-review:** %s\n", readability)
	fmt.Fprintf(&b, "- **confirmed-record-bindings:** %s\n", strings.Join(tokens, ", "))
	fmt.Fprintf(&b, "- **pinned-exit-a-exit-e-candidate-report-candidate-payload-manifest-fingerprints:** %s %s\n", reportFingerprint, manifestFingerprint)
	b.WriteString("- **required-dec-content-versions-fingerprints-and-approval-envelopes-reviewed:** None\n")
	b.WriteString("- **invalidated-by:** None\n")
	b.WriteString("- **cnf-signature-envelope-version:** 1\n")
	b.WriteString("- **cnf-payload-fingerprint:** PENDING\n")
	fmt.Fprintf(&b, "- **signature-or-approval-reference-signature-timestamp:** %s at %s\n", approver, timestamp)
	payload := []byte(b.String())
	fingerprint, err := capabilities.SemanticContentFingerprint(payload, "cnf-payload-fingerprint")
	if err != nil {
		return "", nil, err
	}
	payload = []byte(strings.Replace(string(payload),
		"- **cnf-payload-fingerprint:** PENDING", "- **cnf-payload-fingerprint:** "+fingerprint, 1))
	owner := "manifest-" + strings.TrimPrefix(manifestFingerprint, "sha256:") + "/" + group
	recordID, err := capabilities.TypedID("CNF", namespace, owner, "confirmation")
	if err != nil {
		return "", nil, err
	}
	return recordID, payload, nil
}

// ConfirmationGroup is one human-confirmation target: one synthesis batch, one
// PRF profile, or one HBK section.
type ConfirmationGroup struct {
	Key      string
	Label    string
	Type     string // synthesis | PRF | HBK
	Bindings []ManifestBinding
}

// ConfirmationGroups splits the manifest bindings into ordered confirmation
// groups.
func ConfirmationGroups(bindings []ManifestBinding) []ConfirmationGroup {
	var synthesis []ManifestBinding
	var singles []ManifestBinding
	for _, binding := range bindings {
		switch binding.RecordType {
		case "PRF", "HBK":
			singles = append(singles, binding)
		default:
			synthesis = append(synthesis, binding)
		}
	}
	var groups []ConfirmationGroup
	if len(synthesis) > 0 {
		groups = append(groups, ConfirmationGroup{
			Key:      "synthesis-records",
			Label:    fmt.Sprintf("%d synthesis records", len(synthesis)),
			Type:     "synthesis",
			Bindings: synthesis,
		})
	}
	for _, binding := range singles {
		groupType := binding.RecordType
		label := binding.RecordID
		if binding.RecordType == "HBK" {
			label = binding.Path + " " + binding.RecordID
		}
		groups = append(groups, ConfirmationGroup{
			Key:      binding.RecordID,
			Label:    label,
			Type:     groupType,
			Bindings: []ManifestBinding{binding},
		})
	}
	return groups
}

// ValidateConfirmation checks one stored CNF payload against the manifest
// bindings and the pinned candidate fingerprints.
func ValidateConfirmation(payload []byte, bindings map[string]ManifestBinding, reportFingerprint, manifestFingerprint string) error {
	values := parseConfirmation(payload)
	recomputed, err := capabilities.SemanticContentFingerprint(payload, "cnf-payload-fingerprint")
	if err != nil {
		return err
	}
	if values["cnf-payload-fingerprint"] != recomputed {
		return fmt.Errorf("handoff: CNF payload fingerprint does not recompute")
	}
	if !strings.Contains(values["pinned-exit-a-exit-e-candidate-report-candidate-payload-manifest-fingerprints"], reportFingerprint) ||
		!strings.Contains(values["pinned-exit-a-exit-e-candidate-report-candidate-payload-manifest-fingerprints"], manifestFingerprint) {
		return fmt.Errorf("handoff: CNF does not pin the candidate report and manifest fingerprints")
	}
	switch values["handbook-readability-review"] {
	case "Pass", "Fail", "Not-Applicable":
	default:
		return fmt.Errorf("handoff: CNF has invalid readability value %q", values["handbook-readability-review"])
	}
	tokens := strings.Split(values["confirmed-record-bindings"], ",")
	resolved := 0
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		matched := false
		for _, binding := range bindings {
			if token == binding.BindingToken() {
				matched = true
				resolved++
				break
			}
		}
		if !matched {
			return fmt.Errorf("handoff: CNF binding %q is not an exact candidate manifest entry", token)
		}
	}
	if resolved == 0 {
		return fmt.Errorf("handoff: CNF binds no candidate manifest entries")
	}
	return nil
}

// ConfirmationBindings returns the manifest tokens confirmed by a CNF payload.
func ConfirmationBindings(payload []byte) []string {
	values := parseConfirmation(payload)
	var tokens []string
	for _, token := range strings.Split(values["confirmed-record-bindings"], ",") {
		token = strings.TrimSpace(token)
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func parseConfirmation(payload []byte) map[string]string {
	values := map[string]string{}
	for _, line := range strings.Split(string(payload), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- **") {
			continue
		}
		end := strings.Index(trimmed, ":**")
		if end < 0 {
			continue
		}
		values[capabilities.ColumnName(trimmed[4:end+2])] = strings.TrimSpace(trimmed[end+4:])
	}
	return values
}

// ReportFingerprintFromManifest returns the candidate report fingerprint
// declared by the manifest.
func ReportFingerprintFromManifest(payload []byte) string {
	return reportField(payload, "exit-e-candidate-report-fingerprint")
}

// ArtifactFingerprint wraps artifact.PayloadFingerprint for callers.
func ArtifactFingerprint(artifactType, path string, payload []byte) (string, error) {
	return artifact.PayloadFingerprint(artifactType, path, payload)
}
