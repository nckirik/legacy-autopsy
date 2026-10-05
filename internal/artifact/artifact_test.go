package artifact

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func doc(t *testing.T) *cdl.EIRDoc {
	t.Helper()
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	return res.EIR
}

func validEnvelope() []byte {
	return []byte("- **Required Human Signatories / Signature Timestamps:** Alice 2026-10-05\n- **Envelope Audit:** reviewed\n")
}

func TestValidatePayload(t *testing.T) {
	d := doc(t)
	if err := ValidatePayload(d, "OUTER-SIGNATURE-ENVELOPE", validEnvelope()); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}
	unknown := []byte("- **Required Human Signatories / Signature Timestamps:** Alice\n- **Surprise Field:** x\n- **Envelope Audit:** y\n")
	if err := ValidatePayload(d, "OUTER-SIGNATURE-ENVELOPE", unknown); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("unknown field accepted: %v", err)
	}
	reordered := []byte("- **Envelope Audit:** y\n- **Required Human Signatories / Signature Timestamps:** Alice\n")
	if err := ValidatePayload(d, "OUTER-SIGNATURE-ENVELOPE", reordered); err == nil || !strings.Contains(err.Error(), "order") {
		t.Fatalf("reordered payload accepted: %v", err)
	}
	duplicate := []byte("- **Envelope Audit:** y\n- **Envelope Audit:** z\n- **Required Human Signatories / Signature Timestamps:** Alice\n")
	if err := ValidatePayload(d, "OUTER-SIGNATURE-ENVELOPE", duplicate); err == nil || !strings.Contains(err.Error(), "repeats") {
		t.Fatalf("duplicate field accepted: %v", err)
	}
	missing := []byte("- **Envelope Audit:** y\n")
	if err := ValidatePayload(d, "OUTER-SIGNATURE-ENVELOPE", missing); err == nil || !strings.Contains(err.Error(), "missing required") {
		t.Fatalf("missing required field accepted: %v", err)
	}
	if err := ValidatePayload(d, "INVARIANT", []byte("- **Constraint:** x\n")); err == nil || !strings.Contains(err.Error(), "not artifact") {
		t.Fatalf("non-artifact FIELD accepted: %v", err)
	}
}

func TestPayloadFingerprintAndInstanceBinding(t *testing.T) {
	first, err := PayloadFingerprint("OUTER-BUNDLE-MANIFEST", "handoff/27.md", validEnvelope())
	if err != nil {
		t.Fatal(err)
	}
	second, err := PayloadFingerprint("OUTER-BUNDLE-MANIFEST", "handoff/27.md", []byte(strings.ReplaceAll(string(validEnvelope()), "\r\n", "\n")))
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("transport variance changed the artifact fingerprint")
	}
	carrier := append(validEnvelope(), []byte("- **ARTIFACT-PAYLOAD-FINGERPRINT:** stale-value\n")...)
	carrierFingerprint, err := PayloadFingerprint("OUTER-BUNDLE-MANIFEST", "handoff/27.md", carrier)
	if err != nil {
		t.Fatal(err)
	}
	if carrierFingerprint != first {
		t.Fatal("carrier value participates in its own preimage")
	}
	otherPath, err := PayloadFingerprint("OUTER-BUNDLE-MANIFEST", "handoff/28.md", validEnvelope())
	if err != nil {
		t.Fatal(err)
	}
	if otherPath == first {
		t.Fatal("path is not bound into the fingerprint")
	}
	binding, err := InstanceBinding("OUTER-BUNDLE-MANIFEST", "handoff/27.md", first)
	if err != nil || !strings.Contains(binding, first) {
		t.Fatalf("binding = %q err = %v", binding, err)
	}
}

func TestPackageMembersAndVerification(t *testing.T) {
	members := []Member{
		{Path: "handoff/27.md", Bytes: validEnvelope()},
		{Path: "handoff/26.md", Bytes: []byte("# certificate\n")},
	}
	manifest, err := PackageMembers(members)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest) != 2 || manifest[0].Path != "handoff/26.md" {
		t.Fatalf("manifest = %+v", manifest)
	}
	if err := VerifyMembers(members, manifest); err != nil {
		t.Fatal(err)
	}
	tampered := []Member{{Path: "handoff/27.md", Bytes: append(validEnvelope(), 'x')}, members[1]}
	if err := VerifyMembers(tampered, manifest); err == nil || !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("tamper accepted: %v", err)
	}
	extra := append(append([]Member(nil), members...), Member{Path: "handoff/25.md", Bytes: []byte("x")})
	if err := VerifyMembers(extra, manifest); err == nil || !strings.Contains(err.Error(), "declares") {
		t.Fatalf("extra member accepted: %v", err)
	}
	missing := members[:1]
	if err := VerifyMembers(missing, manifest); err == nil {
		t.Fatal("missing member accepted")
	}
	stableA, _ := PackageMembers(members)
	stableB, _ := PackageMembers([]Member{members[1], members[0]})
	if len(stableA) != len(stableB) || stableA[0] != stableB[0] || stableA[1] != stableB[1] {
		t.Fatal("member manifest depends on input order")
	}
}
