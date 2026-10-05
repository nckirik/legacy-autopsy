package certify

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func TestCanonicalStringSet(t *testing.T) {
	got, err := CanonicalStringSet([]any{"b", "a"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("set not canonicalized: %v", got)
	}
	if _, err := CanonicalStringSet([]any{"a", "a"}); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate element accepted: %v", err)
	}
	if _, err := CanonicalStringSet([]any{"a", 1}); err == nil || !strings.Contains(err.Error(), "non-string") {
		t.Fatalf("non-string element accepted: %v", err)
	}
	if _, err := CanonicalStringSet("nope"); err == nil {
		t.Fatal("non-set accepted")
	}
}

func TestCheckValueRejectsNoncanonicalSet(t *testing.T) {
	spec := cdl.EIRFieldSpec{Name: "x", Type: "set<string>", Required: true}
	if err := checkValue(spec, []string{"b", "a"}, nil); err != nil {
		t.Fatalf("unsorted set rejected: %v", err)
	}
	if err := checkValue(spec, []any{"a", "a"}, nil); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate set accepted: %v", err)
	}
}

func TestValidateRecord(t *testing.T) {
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	doc := res.EIR
	if err := ValidateRecord(doc, "SEMANTIC-RECORD-COUNT", map[string]any{"record-type": "CMP", "count": 3}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRecord(doc, "SEMANTIC-RECORD-COUNT", map[string]any{"record-type": "CMP"}); err == nil || !strings.Contains(err.Error(), "missing required") {
		t.Fatalf("missing required accepted: %v", err)
	}
	if err := ValidateRecord(doc, "SEMANTIC-RECORD-COUNT", map[string]any{"record-type": "CMP", "count": 1, "extra": "x"}); err == nil || !strings.Contains(err.Error(), "unknown fields") {
		t.Fatalf("unknown field accepted: %v", err)
	}
	if err := ValidateRecord(doc, "SEMANTIC-RECORD-COUNT", map[string]any{"record-type": "CMP", "count": "three"}); err == nil || !strings.Contains(err.Error(), "integer required") {
		t.Fatalf("wrong type accepted: %v", err)
	}
	if err := ValidateRecord(doc, "CAPABILITY-STATE-BLOCK", map[string]any{
		"capability-state": "Nonsense", "runtime-state-matrix": "m",
		"historical-intended-personas": "h", "dependencies": "d", "last-known-state": "l",
	}); err == nil || !strings.Contains(err.Error(), "not in enum") {
		t.Fatalf("invalid enum accepted: %v", err)
	}
	if err := ValidateRecord(doc, "UNKNOWN-SCHEMA", nil); err == nil || !strings.Contains(err.Error(), "unknown record schema") {
		t.Fatalf("unknown schema accepted: %v", err)
	}
}
