package capabilities

import (
	"strings"
	"testing"
)

func coldResumeFixture() []ColdResumeField {
	fields := []ColdResumeField{}
	for _, label := range coldResumeRequiredFields {
		fields = append(fields, ColdResumeField{Label: label, Value: "value for " + label})
	}
	return fields
}

func TestColdResumeFingerprint(t *testing.T) {
	fields := append(coldResumeFixture(), ColdResumeField{Label: ColdResumeCarrier, Value: "ignored"})
	got, err := ColdResumeFingerprint("example", "INV-1", fields)
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for _, field := range coldResumeFixture() {
		b.WriteString("- **" + field.Label + ":** " + field.Value + "\n")
	}
	canonical, err := CanonicalMarkdown([]byte(b.String()))
	if err != nil {
		t.Fatal(err)
	}
	want := shaHex("COLD-RESUME|example|INV-1\n" + string(canonical))
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestColdResumeFingerprintRejectsInvalid(t *testing.T) {
	if _, err := ColdResumeFingerprint("", "INV-1", coldResumeFixture()); err == nil {
		t.Fatal("empty namespace accepted")
	}
	if _, err := ColdResumeFingerprint("example", "", coldResumeFixture()); err == nil {
		t.Fatal("empty invocation accepted")
	}
	missing := coldResumeFixture()
	missing = missing[:len(missing)-1]
	if _, err := ColdResumeFingerprint("example", "INV-1", missing); err == nil {
		t.Fatal("missing field accepted")
	}
	unknown := append(coldResumeFixture(), ColdResumeField{Label: "Not A Field", Value: "x"})
	if _, err := ColdResumeFingerprint("example", "INV-1", unknown); err == nil {
		t.Fatal("unknown field accepted")
	}
	duplicate := append(coldResumeFixture(), coldResumeFixture()[0])
	if _, err := ColdResumeFingerprint("example", "INV-1", duplicate); err == nil {
		t.Fatal("duplicate field accepted")
	}
}
