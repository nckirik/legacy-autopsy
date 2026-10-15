package store

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

func TestRecordGetListInspect(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	earlier := "- **Constraint:** a\n- **Verification Vectors:** v2\n- **Discovered In:** SRC-1\n"
	later := "- **Constraint:** b\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-2\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-B", 1, []byte(later)); err != nil {
		t.Fatal(err)
	}
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-A", 1, []byte(earlier)); err != nil {
		t.Fatal(err)
	}
	payload, err := Get(db, res.EIR, "INVARIANT", "INV-A")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(earlier))
	if err != nil {
		t.Fatal(err)
	}
	if string(payload) != string(canonical) {
		t.Fatalf("Get returned %q, want canonical payload", payload)
	}
	ids, err := List(db, res.EIR, "INVARIANT")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "INV-A,INV-B" {
		t.Fatalf("List = %v, want sorted IDs", ids)
	}
	if _, err := Get(db, res.EIR, "INVARIANT", "NOPE"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("missing record: %v", err)
	}
	if _, err := Get(db, res.EIR, "BLOCK-EVIDENCE-SUMMARY", "BES-1"); err == nil || !strings.Contains(err.Error(), "not a table class") {
		t.Fatalf("non-table class: %v", err)
	}
	view, err := Inspect(db, res.EIR, "INVARIANT")
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{
		"DEBUG RENDER - NON-AUTHORITATIVE",
		"# INVARIANT (debug)",
		"## INV-A",
		"## INV-B",
		"la_semantic_fingerprint",
		"```markdown",
	} {
		if !strings.Contains(string(view), marker) {
			t.Fatalf("inspect view missing %q:\n%s", marker, view)
		}
	}
}
