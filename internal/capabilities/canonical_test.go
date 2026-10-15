package capabilities

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func shaHex(preimage string) string {
	sum := sha256.Sum256([]byte(preimage))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func TestCanonicalMarkdown(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{
			name: "lf, trailing space, blank collapse",
			in:   "## T\r\n\r\n\r\n- **A:** x   \r\n",
			want: "## T\n\n- **A:** x\n",
		},
		{
			name: "table normalization",
			in:   "| a |  b   |\n| :--- | ---: |\n| 1 | 2 |\n",
			want: "| a | b |\n| --- | --- |\n| 1 | 2 |\n",
		},
		{
			name: "code fence preserved verbatim",
			in:   "```\na   \nb\n```\n",
			want: "```\na   \nb\n```\n",
		},
		{
			name: "trailing blanks removed",
			in:   "x\n\n\n",
			want: "x\n",
		},
		{
			name: "empty payload",
			in:   "",
			want: "\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CanonicalMarkdown([]byte(tc.in))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestSemanticContentFingerprintExcludesCarrier(t *testing.T) {
	payload := "## Block\n\n- **PRF-ID:** PRF-Q38FN\n- **SEMANTIC-CONTENT-FINGERPRINT:** sha256:placeholder\n\nBody text.\n"
	got, err := SemanticContentFingerprint([]byte(payload), "SEMANTIC-CONTENT-FINGERPRINT")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := CanonicalMarkdown([]byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	var kept []string
	for _, line := range strings.Split(strings.TrimSuffix(string(canonical), "\n"), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- **SEMANTIC-CONTENT-FINGERPRINT:**") {
			continue
		}
		kept = append(kept, line)
	}
	want := shaHex(strings.Join(kept, "\n") + "\n")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestDecisionContentFingerprintEnvelopeInvariance(t *testing.T) {
	content := "### DEC-ABCDEF012345 Choose a migration path\n\n" +
		"- **DECISION-CONTENT-VERSION:** 1\n" +
		"- **DECISION-CONTENT-FINGERPRINT:** sha256:placeholder\n" +
		"- **Type:** target\n" +
		"- **Observed Inputs:** MOD-111111111111 version 1 sha256:aaa\n" +
		"- **Options / Decision / Rationale / Risks:** Keep the module.\n" +
		"- **Exact Scope and Denominator Effect:** None.\n" +
		"- **Affected Records / Rollout / Rollback:** None.\n" +
		"- **Content Supersession:** None\n"
	base, err := DecisionContentFingerprint([]byte(content), "DECISION-CONTENT-FINGERPRINT")
	if err != nil {
		t.Fatal(err)
	}
	approved := content + "\n<!-- DECISION-APPROVAL-ENVELOPE: excluded from decision-content version/fingerprint -->\n\n" +
		"- **Envelope Kind:** DECISION-APPROVAL\n" +
		"- **Envelope Version:** 1\n" +
		"- **Approval Status:** Approved\n" +
		"- **Owner / Approver / Authority Basis / Approval Timestamp:** operator 2026-10-14\n" +
		"- **Approved Content Binding:** DEC-ABCDEF012345 version 1 " + base + "\n" +
		"- **Approval Signature / Audit:** review decisions\n" +
		"<!-- END-DECISION-APPROVAL-ENVELOPE -->\n"
	withEnvelope, err := DecisionContentFingerprint([]byte(approved), "DECISION-CONTENT-FINGERPRINT")
	if err != nil {
		t.Fatal(err)
	}
	if base != withEnvelope {
		t.Fatalf("envelope changed the decision-content fingerprint: %s vs %s", base, withEnvelope)
	}
	changed := strings.Replace(approved, "Keep the module.", "Retire the module.", 1)
	mutated, err := DecisionContentFingerprint([]byte(changed), "DECISION-CONTENT-FINGERPRINT")
	if err != nil {
		t.Fatal(err)
	}
	if mutated == base {
		t.Fatal("content change did not change the decision-content fingerprint")
	}
	if _, err := DecisionContentFingerprint([]byte(content+"<!-- DECISION-APPROVAL-ENVELOPE: x -->\n"), "DECISION-CONTENT-FINGERPRINT"); err == nil {
		t.Fatal("partial envelope region was accepted")
	}
}

func TestDecisionContentFingerprintCarrierExclusion(t *testing.T) {
	withoutCarrier := "### DEC-ABCDEF012345 Choose\n\n- **Type:** target\n"
	withCarrier := withoutCarrier + "- **DECISION-CONTENT-FINGERPRINT:** sha256:placeholder\n"
	first, err := DecisionContentFingerprint([]byte(withCarrier), "DECISION-CONTENT-FINGERPRINT")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := CanonicalMarkdown([]byte(withoutCarrier))
	if err != nil {
		t.Fatal(err)
	}
	if first != shaHex(string(canonical)) {
		t.Fatalf("carrier exclusion mismatch: %s", first)
	}
}

func TestEnvelopeFingerprint(t *testing.T) {
	region := []byte("<!-- CERTIFICATION: X -->\n- **Status:** B-CONFIRMED\n<!-- END -->\n")
	canonical, err := CanonicalMarkdown(region)
	if err != nil {
		t.Fatal(err)
	}
	want := shaHex("ENVELOPE|a/b.md|PRF-Q38FN|CERTIFICATION|2\n" + string(canonical))
	got, err := EnvelopeFingerprint("a/b.md", "PRF-Q38FN", "CERTIFICATION", 2, region)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if _, err := EnvelopeFingerprint("a/b.md", "PRF-Q38FN", "UNKNOWN", 1, region); err == nil {
		t.Fatal("unknown kind accepted")
	}
	if _, err := EnvelopeFingerprint("a/b.md", "PRF-Q38FN", "CERTIFICATION", 0, region); err == nil {
		t.Fatal("versionless envelope accepted")
	}
	if _, err := EnvelopeFingerprint("", "PRF-Q38FN", "CERTIFICATION", 1, region); err == nil {
		t.Fatal("empty path accepted")
	}
}

func TestPackageMemberFingerprint(t *testing.T) {
	file := []byte("# Doc\r\n\r\n| a |  b |\r\n")
	canonical, err := CanonicalMarkdown(file)
	if err != nil {
		t.Fatal(err)
	}
	want := shaHex("PACKAGE-MEMBER|pkg/x.md\n" + string(canonical))
	got, err := PackageMemberFingerprint("pkg/x.md", file)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestTransportFingerprint(t *testing.T) {
	bindings := []TransportBinding{
		{RecordType: "CMP", RecordID: "CMP-AAAA", Version: "2", Fingerprint: "sha256:aa"},
		{RecordType: "CLM", RecordID: "CLM-BBBB", Version: "None", Fingerprint: "sha256:bb"},
	}
	want := shaHex("FILE-TRANSPORT|pkg/x.md\nCMP|CMP-AAAA|2|sha256:aa\nCLM|CLM-BBBB|None|sha256:bb\n")
	got, err := TransportFingerprint("pkg/x.md", bindings)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestEvidenceSetFingerprint(t *testing.T) {
	bindings := []EvidenceBinding{
		{InputPath: "p/z.md", Kind: "PROTOCOL-REGISTRY", RecordOrArtifact: "protocol", Version: "4.1.2", Fingerprint: "sha256:zz"},
		{InputPath: "p/a.md", Kind: "RECORD-PAYLOAD", RecordOrArtifact: "CMP-AAAA", Version: "3", Fingerprint: "sha256:aa"},
	}
	want := shaHex("EVIDENCE-SET|EXIT-E-CANDIDATE-REPORT\nCHECK\nFINAL-01\nSystem / Protocol / Snapshot\n" +
		"| p/a.md | RECORD-PAYLOAD | CMP-AAAA | 3 | sha256:aa |\n" +
		"| p/z.md | PROTOCOL-REGISTRY | protocol | 4.1.2 | sha256:zz |")
	got, err := EvidenceSetFingerprint("EXIT-E-CANDIDATE-REPORT", "CHECK", "FINAL-01", "System / Protocol / Snapshot", bindings)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}

	if _, err := EvidenceSetFingerprint("T", "CHECK", "R", "S", []EvidenceBinding{{Kind: "NOT-A-KIND"}}); err == nil {
		t.Fatal("unknown kind accepted")
	}
	if _, err := EvidenceSetFingerprint("T", "CHECK", "R", "S", []EvidenceBinding{
		{Kind: "PROTOCOL-REGISTRY"}, {Kind: "PROTOCOL-REGISTRY"},
	}); err == nil {
		t.Fatal("duplicate tuple accepted")
	}
	if _, err := EvidenceSetFingerprint("T", "CHECK", "R", "S", []EvidenceBinding{{Kind: "VALIDATION-SUMMARY"}}); err == nil {
		t.Fatal("unanchored evidence set accepted")
	}
	if _, err := EvidenceSetFingerprint("T", "PENDING", "R", "S", []EvidenceBinding{{Kind: "PROTOCOL-REGISTRY"}}); err == nil {
		t.Fatal("invalid row kind accepted")
	}
}

func TestHashDomainAndPostHashBinding(t *testing.T) {
	domain, err := HashDomainIdentity("SCOPE-CERTIFICATE", "pkg/cert.md")
	if err != nil {
		t.Fatal(err)
	}
	if domain != "SCOPE-CERTIFICATE|pkg/cert.md" {
		t.Fatalf("domain %s", domain)
	}
	binding, err := PostHashArtifactInstanceBinding("SCOPE-CERTIFICATE", "pkg/cert.md", "sha256:cc")
	if err != nil {
		t.Fatal(err)
	}
	if binding != "SCOPE-CERTIFICATE|pkg/cert.md|sha256:cc" {
		t.Fatalf("binding %s", binding)
	}
}

func TestSemanticContentFingerprintCarrierCaseInsensitive(t *testing.T) {
	lower := "## Block\n\n- **field:** value\n- **semantic-content-fingerprint:** sha256:placeholder\n"
	upper := "## Block\n\n- **field:** value\n- **SEMANTIC-CONTENT-FINGERPRINT:** sha256:placeholder\n"
	first, err := SemanticContentFingerprint([]byte(lower), "semantic-content-fingerprint")
	if err != nil {
		t.Fatal(err)
	}
	second, err := SemanticContentFingerprint([]byte(upper), "semantic-content-fingerprint")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("carrier exclusion is case-sensitive: %s vs %s", first, second)
	}
}
