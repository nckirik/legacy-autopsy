// Package envelope emits, parses, and verifies protocol approval/signature
// envelope regions. The runtime owns envelope mechanics; executors and humans
// author only the fields inside a region.
package envelope

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// Kind is one closed envelope kind.
type Kind string

const (
	Certification      Kind = "CERTIFICATION"
	DecisionApproval   Kind = "DECISION-APPROVAL"
	CNFDigestSignature Kind = "CNF-DIGEST-SIGNATURE"
)

// Markers returns the exact begin and end markers for a kind. The begin marker
// preserves the protocol's explanatory suffix form.
func Markers(kind Kind) (string, string, error) {
	switch kind {
	case Certification:
		return "<!-- CERTIFICATION-ENVELOPE: excluded from semantic version/fingerprint -->", "<!-- END-CERTIFICATION-ENVELOPE -->", nil
	case DecisionApproval:
		return "<!-- DECISION-APPROVAL-ENVELOPE: excluded from decision-content version/fingerprint -->", "<!-- END-DECISION-APPROVAL-ENVELOPE -->", nil
	case CNFDigestSignature:
		return "<!-- CNF-DIGEST-SIGNATURE-ENVELOPE: excluded from confirmation payload version/fingerprint -->", "<!-- END-CNF-DIGEST-SIGNATURE-ENVELOPE -->", nil
	default:
		return "", "", fmt.Errorf("envelope: unknown kind %q", kind)
	}
}

// Emit renders a canonical envelope region for a target record at a positive
// version. Field lines are emitted in the given order.
func Emit(kind Kind, version int, fields []string) ([]byte, error) {
	if version < 1 {
		return nil, fmt.Errorf("envelope: version must be positive")
	}
	begin, end, err := Markers(kind)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString(begin)
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "- **Envelope Kind:** %s\n", kind)
	fmt.Fprintf(&b, "- **Envelope Version:** %d\n", version)
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		b.WriteString(field)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(end)
	b.WriteString("\n")
	return []byte(b.String()), nil
}

// Region extracts the exact bytes between the begin and end markers, including
// the markers themselves. Exactly one region is allowed.
func Region(payload []byte, kind Kind) ([]byte, error) {
	begin, end, err := Markers(kind)
	if err != nil {
		return nil, err
	}
	text := string(payload)
	start := strings.Index(text, begin)
	if start < 0 {
		return nil, fmt.Errorf("envelope: %s begin marker not found", kind)
	}
	if strings.Index(text[start+len(begin):], begin) >= 0 {
		return nil, fmt.Errorf("envelope: %s has duplicate begin markers", kind)
	}
	stop := strings.Index(text[start:], end)
	if stop < 0 {
		return nil, fmt.Errorf("envelope: %s end marker not found", kind)
	}
	stop += start + len(end)
	return []byte(text[start:stop]), nil
}

// Verify recomputes the canonical envelope fingerprint binding the containing
// path, target record, kind, and version over the exact region bytes.
func Verify(payload []byte, path, targetID string, kind Kind, version int) (string, error) {
	if targetID == "" {
		return "", fmt.Errorf("envelope: empty target record ID")
	}
	region, err := Region(payload, kind)
	if err != nil {
		return "", err
	}
	return capabilities.EnvelopeFingerprint(path, targetID, string(kind), version, region)
}
