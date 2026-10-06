package envelope

import (
	"strings"
	"testing"
)

func TestEmitParseVerify(t *testing.T) {
	region, err := Emit(DecisionApproval, 2, []string{
		"- **Approver:** Alice",
		"- **Authority Basis:** operator",
		"- **Timestamp:** 2026-10-05T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte("# DEC-0001\n\n- **Decision:** keep retries\n\n" + string(region))
	first, err := Verify(payload, "15-DECISIONS.md", "DEC-0001", DecisionApproval, 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Verify(payload, "15-DECISIONS.md", "DEC-0001", DecisionApproval, 2)
	if first != second {
		t.Fatal("envelope fingerprint is not deterministic")
	}
	otherPath, _ := Verify(payload, "archive/15-DECISIONS.md", "DEC-0001", DecisionApproval, 2)
	if otherPath == first {
		t.Fatal("path not bound into the envelope fingerprint")
	}
	otherTarget, _ := Verify(payload, "15-DECISIONS.md", "DEC-0002", DecisionApproval, 2)
	if otherTarget == first {
		t.Fatal("target not bound into the envelope fingerprint")
	}
	tampered := []byte(strings.Replace(string(payload), "Alice", "Mallory", 1))
	tamperedFingerprint, _ := Verify(tampered, "15-DECISIONS.md", "DEC-0001", DecisionApproval, 2)
	if tamperedFingerprint == first {
		t.Fatal("tampering did not change the envelope fingerprint")
	}
}

func TestEmitRejectsAndRegionFailsClosed(t *testing.T) {
	if _, err := Emit(Certification, 0, nil); err == nil {
		t.Fatal("zero version accepted")
	}
	if _, err := Emit(Kind("NOPE"), 1, nil); err == nil {
		t.Fatal("unknown kind accepted")
	}
	if _, err := Region([]byte("no envelopes here"), Certification); err == nil {
		t.Fatal("missing marker accepted")
	}
	region, _ := Emit(Certification, 1, []string{"- **HANDBOOK-STATUS:** [B-CONFIRMED]"})
	duplicated := append(append([]byte(nil), region...), region...)
	if _, err := Region(duplicated, Certification); err == nil {
		t.Fatal("duplicate begin markers accepted")
	}
	unterminated := region[:len(region)-len("<!-- END-CERTIFICATION-ENVELOPE -->\n")]
	if _, err := Region(unterminated, Certification); err == nil {
		t.Fatal("missing end marker accepted")
	}
	if _, err := Verify(region, "x.md", "", Certification, 1); err == nil {
		t.Fatal("empty target accepted")
	}
}
