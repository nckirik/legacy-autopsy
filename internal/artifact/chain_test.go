package artifact

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// stage builds a schema-valid artifact payload: every declared field is present
// in order, enum fields take their first declared value, and callers override
// the fields that matter for the chain.
func stage(t *testing.T, doc *cdl.EIRDoc, fieldID string, overrides map[string]string, checks []string) []byte {
	t.Helper()
	var field *cdl.EIRField
	for i := range doc.Declarations.Fields {
		if doc.Declarations.Fields[i].ID == fieldID {
			field = &doc.Declarations.Fields[i]
		}
	}
	if field == nil {
		t.Fatalf("unknown FIELD %s", fieldID)
	}
	enums := map[string][]string{}
	for _, enum := range doc.Declarations.Enums {
		enums[enum.ID] = enum.Values
	}
	var b strings.Builder
	for _, spec := range field.Fields {
		value := "x"
		switch {
		case spec.Type == "int":
			value = "1"
		case len(enums[spec.Type]) > 0:
			value = enums[spec.Type][0]
		}
		if override, ok := overrides[spec.Name]; ok {
			value = override
		}
		if spec.Name == "check-results" && len(checks) > 0 {
			b.WriteString("- **check-results:**\n")
			b.WriteString("  | Check ID | Requirement | Result | Evidence Fingerprint |\n")
			b.WriteString("  | :-- | :-- | :-- | :-- |\n")
			for _, check := range checks {
				b.WriteString("  | " + check + " | requirement | Pass | sha256:x |\n")
			}
			continue
		}
		if spec.Name == "package-members" {
			b.WriteString("- **package-members:**\n")
			b.WriteString("  | Artifact Type | Normalized Path | Package-Member File Fingerprint |\n")
			b.WriteString("  | :-- | :-- | :-- |\n")
			continue
		}
		b.WriteString("- **" + spec.Name + ":** " + value + "\n")
	}
	return []byte(b.String())
}

func TestVerifyChain(t *testing.T) {
	d := doc(t)
	snapshot := "demo | 4.2.1 | dev | snap-1"
	finalRegistry := RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", RegistryCheckIDs(d, "FINAL-CHECK"))

	candidate := stage(t, d, "EXIT-E-CANDIDATE-REPORT",
		map[string]string{"system-protocol-snapshot": snapshot, "candidate-result": "Ready-For-Human-Review"},
		RegistryCheckIDs(d, "EXIT-E-CANDIDATE-CHECKS"))
	candidateFingerprint, err := PayloadFingerprint("EXIT-E-CANDIDATE-REPORT", "23-candidate.md", candidate)
	if err != nil {
		t.Fatal(err)
	}
	manifest := stage(t, d, "CANDIDATE-PAYLOAD-MANIFEST",
		map[string]string{"system-protocol-snapshot": snapshot, "exit-e-candidate-report-fingerprint": candidateFingerprint}, nil)
	manifestFingerprint, _ := PayloadFingerprint("CANDIDATE-PAYLOAD-MANIFEST", "24-manifest.md", manifest)
	readiness := stage(t, d, "EXIT-E-CONTENT-READINESS-REPORT",
		map[string]string{"system-protocol-snapshot": snapshot, "exit-e-candidate-report-fingerprint": candidateFingerprint, "candidate-payload-manifest-fingerprint": manifestFingerprint, "content-readiness-result": "Ready-For-Certificate", "exit-e-prepackage-state": "Pending"},
		RegistryCheckIDs(d, "EXIT-E-CONTENT-READINESS-CHECKS"))
	readinessFingerprint, _ := PayloadFingerprint("EXIT-E-CONTENT-READINESS-REPORT", "25-readiness.md", readiness)
	certificate := stage(t, d, "SCOPE-CERTIFICATE",
		map[string]string{"system-protocol-snapshot": snapshot, "exit-e-content-readiness-report-fingerprint": readinessFingerprint}, nil)
	certificateFingerprint, _ := PayloadFingerprint("SCOPE-CERTIFICATE", "26-certificate.md", certificate)
	outer := stage(t, d, "OUTER-BUNDLE-MANIFEST",
		map[string]string{
			"system-protocol-snapshot":                      snapshot,
			"final-verification-check-registry-fingerprint": finalRegistry,
			"exit-e-content-readiness-report-fingerprint":   readinessFingerprint,
			"scope-certificate-fingerprint":                 certificateFingerprint,
			"exit-e-status":                                 "Passed",
		}, nil)
	stages := []ChainStage{
		{Type: "EXIT-E-CANDIDATE-REPORT", Path: "23-candidate.md", Payload: candidate},
		{Type: "CANDIDATE-PAYLOAD-MANIFEST", Path: "24-manifest.md", Payload: manifest},
		{Type: "EXIT-E-CONTENT-READINESS-REPORT", Path: "25-readiness.md", Payload: readiness},
		{Type: "SCOPE-CERTIFICATE", Path: "26-certificate.md", Payload: certificate},
		{Type: "OUTER-BUNDLE-MANIFEST", Path: "27-outer.md", Payload: outer},
	}
	if err := VerifyChain(d, stages); err != nil {
		t.Fatalf("valid chain rejected: %v", err)
	}

	broken := append([]ChainStage(nil), stages...)
	broken[2] = ChainStage{Type: "EXIT-E-CONTENT-READINESS-REPORT", Path: "25-readiness.md", Payload: []byte(strings.Replace(string(readiness), candidateFingerprint, "sha256:stale", 1))}
	if err := VerifyChain(d, broken); err == nil || !strings.Contains(err.Error(), "does not reference") {
		t.Fatalf("broken reference accepted: %v", err)
	}
	outOfOrder := append([]ChainStage(nil), stages...)
	outOfOrder[2], outOfOrder[3] = outOfOrder[3], outOfOrder[2]
	if err := VerifyChain(d, outOfOrder); err == nil {
		t.Fatal("out-of-order chain accepted")
	}
	passedEarly := append([]ChainStage(nil), stages...)
	passedEarly[2] = ChainStage{Type: "EXIT-E-CONTENT-READINESS-REPORT", Path: "25-readiness.md", Payload: []byte(strings.Replace(string(readiness), "- **exit-e-prepackage-state:** Pending", "- **exit-e-prepackage-state:** Passed", 1))}
	if err := VerifyChain(d, passedEarly); err == nil || !strings.Contains(err.Error(), "before the outer manifest") {
		t.Fatalf("early Passed accepted: %v", err)
	}
	missingChecks := append([]ChainStage(nil), stages...)
	missingChecks[0] = ChainStage{Type: "EXIT-E-CANDIDATE-REPORT", Path: "23-candidate.md", Payload: []byte(strings.ReplaceAll(string(candidate), "CANDIDATE-01-EXIT-A", "CANDIDATE-99-UNKNOWN"))}
	if err := VerifyChain(d, missingChecks); err == nil {
		t.Fatal("unknown check accepted")
	}
}
