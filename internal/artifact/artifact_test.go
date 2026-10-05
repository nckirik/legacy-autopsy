package artifact

import (
	"fmt"
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

func outerManifest(members []ManifestEntry, snapshot string) []byte {
	var b strings.Builder
	b.WriteString("- **ARTIFACT-TYPE:** OUTER-BUNDLE-MANIFEST\n")
	b.WriteString("- **ARTIFACT-PATH:** handoff/27-OUTER-BUNDLE-MANIFEST.md\n")
	b.WriteString("- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** 4.2.1\n")
	b.WriteString("- **System / Protocol / Snapshot:** " + snapshot + "\n")
	b.WriteString("- **Final Verification Check Registry / Fingerprint:** EXIT-E-FINAL-CHECKS-v1 sha256:x\n")
	b.WriteString("- **Pre-Signature Validation Input-Set Fingerprint:** sha256:y\n")
	b.WriteString("- **Exit E Content-Readiness Report Fingerprint:** sha256:z\n")
	b.WriteString("- **Scope Certificate Fingerprint:** sha256:w\n")
	b.WriteString("- **Package Members:**\n")
	b.WriteString("  | Artifact Type | Normalized Path | Package-Member File Fingerprint |\n")
	b.WriteString("  | :------------ | :-------------- | :------------------------------ |\n")
	for _, entry := range members {
		fmt.Fprintf(&b, "  | OUTER | %s | %s |\n", entry.Path, entry.Fingerprint)
	}
	b.WriteString("- **EXIT-E-STATUS:** Passed\n")
	return []byte(b.String())
}

func TestVerifyBundle(t *testing.T) {
	others := []Member{
		{Path: "handoff/handbook/01.md", Bytes: []byte("- **System / Protocol / Snapshot:** demo | 4.2.1 | dev | snap-1\n")},
		{Path: "handoff/26-SCOPE-CERTIFICATE.md", Bytes: []byte("- **System / Protocol / Snapshot:** demo | 4.2.1 | dev | snap-1\n")},
	}
	manifest, err := PackageMembers(others)
	if err != nil {
		t.Fatal(err)
	}
	outer := outerManifest(manifest, "demo | 4.2.1 | dev | snap-1")
	members := append(append([]Member(nil), others...), Member{Path: "handoff/27-OUTER-BUNDLE-MANIFEST.md", Bytes: outer})
	if err := VerifyBundle("OUTER-BUNDLE-MANIFEST", "handoff/27-OUTER-BUNDLE-MANIFEST.md", outer, members); err != nil {
		t.Fatalf("valid bundle rejected: %v", err)
	}
	tampered := append([]Member(nil), members...)
	for i := range tampered {
		if tampered[i].Path == "handoff/handbook/01.md" {
			tampered[i] = Member{Path: tampered[i].Path, Bytes: append(append([]byte(nil), tampered[i].Bytes...), 'x')}
		}
	}
	if err := VerifyBundle("OUTER-BUNDLE-MANIFEST", "handoff/27-OUTER-BUNDLE-MANIFEST.md", outer, tampered); err == nil {
		t.Fatal("tampered member accepted")
	}
	mixed := append([]Member(nil), members...)
	for i := range mixed {
		if mixed[i].Path == "handoff/26-SCOPE-CERTIFICATE.md" {
			mixed[i] = Member{Path: mixed[i].Path, Bytes: []byte("- **System / Protocol / Snapshot:** other | 4.2.1 | dev | snap-1\n")}
		}
	}
	if err := VerifyBundle("OUTER-BUNDLE-MANIFEST", "handoff/27-OUTER-BUNDLE-MANIFEST.md", outer, mixed); err == nil || !strings.Contains(err.Error(), "snapshot mismatch") {
		t.Fatalf("mixed snapshot accepted: %v", err)
	}
	failedStatus := strings.Replace(string(outer), "**EXIT-E-STATUS:** Passed", "**EXIT-E-STATUS:** Pending", 1)
	if err := VerifyBundle("OUTER-BUNDLE-MANIFEST", "handoff/27-OUTER-BUNDLE-MANIFEST.md", []byte(failedStatus), members); err == nil || !strings.Contains(err.Error(), "EXIT-E-STATUS") {
		t.Fatalf("Pending status accepted: %v", err)
	}
	extra := append(append([]Member(nil), members...), Member{Path: "handoff/extra.md", Bytes: []byte("x")})
	if err := VerifyBundle("OUTER-BUNDLE-MANIFEST", "handoff/27-OUTER-BUNDLE-MANIFEST.md", outer, extra); err == nil {
		t.Fatal("extra member accepted")
	}
}
