package handoff

import (
	"strings"
	"testing"
)

func TestBuildAndValidateConfirmation(t *testing.T) {
	bindings := []ManifestBinding{
		{RecordType: "PRF", RecordID: "PRF-ABCDEF012345", Version: "1", Fingerprint: "sha256:prf"},
		{RecordType: "MOD", RecordID: "MOD-ABCDEF012345", Version: "1", Fingerprint: "sha256:mod"},
	}
	byToken := map[string]ManifestBinding{}
	for _, binding := range bindings {
		byToken[binding.BindingToken()] = binding
	}
	id, payload, err := BuildConfirmation("demo", "PRF-ABCDEF012345", "Ali (architect)", "scope", "limitations", "2026-10-14T10:00:00Z", "Not-Applicable", bindings[:1], "sha256:report", "sha256:manifest")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(id, "CNF-") {
		t.Fatalf("cnf id = %q", id)
	}
	if err := ValidateConfirmation(payload, byToken, "sha256:report", "sha256:manifest"); err != nil {
		t.Fatalf("valid confirmation rejected: %v", err)
	}
	if len(ConfirmationBindings(payload)) != 1 {
		t.Fatalf("bindings = %v", ConfirmationBindings(payload))
	}
	tampered := []byte(strings.Replace(string(payload), "PRF-ABCDEF012345@1@sha256:prf", "PRF-ABCDEF012345@1@sha256:other", 1))
	if err := ValidateConfirmation(tampered, byToken, "sha256:report", "sha256:manifest"); err == nil {
		t.Fatalf("tampered confirmation accepted")
	}
	groups := ConfirmationGroups(bindings)
	if len(groups) != 2 || groups[0].Type != "synthesis" || groups[1].Type != "PRF" {
		t.Fatalf("groups = %+v", groups)
	}
}
