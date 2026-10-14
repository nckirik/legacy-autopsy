package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

func TestSignOuterPreservesPayloadFingerprint(t *testing.T) {
	payload := "<!-- ARTIFACT-PAYLOAD: OUTER-BUNDLE-MANIFEST -->\n\n" +
		"- **EXIT-E-STATUS:** Passed\n" +
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n" +
		"<!-- END-ARTIFACT-PAYLOAD: OUTER-BUNDLE-MANIFEST -->\n\n" +
		"<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->\n\n" +
		"- **Required Human Signatories / Signature Timestamps:** None\n" +
		"- **Envelope Audit:** pending human signature\n" +
		"<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->\n"
	path := filepath.Join(t.TempDir(), "outer.md")
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	fingerprint, err := artifact.PayloadFingerprint(OuterManifestArtifactType, OuterManifestDefaultPath, []byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	withCarrier := strings.Replace(payload, "- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1)
	if err := os.WriteFile(path, []byte(withCarrier), 0o644); err != nil {
		t.Fatal(err)
	}
	signedFingerprint, envelopeFingerprint, err := SignOuter(path, OuterManifestDefaultPath, "Ali (architect)", "", "2026-10-14T12:00:00Z", "approved by Ali; no cryptographic signature")
	if err != nil {
		t.Fatal(err)
	}
	if signedFingerprint != fingerprint {
		t.Fatalf("signing changed the payload fingerprint: %s -> %s", fingerprint, signedFingerprint)
	}
	if envelopeFingerprint == "" {
		t.Fatalf("empty envelope fingerprint")
	}
	signed, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(signed), "Ali (architect) at 2026-10-14T12:00:00Z") {
		t.Fatalf("signature envelope missing approver")
	}
	if outerSignatureEnvelope(signed) == "" {
		t.Fatalf("outer signature envelope not recognized")
	}
}
