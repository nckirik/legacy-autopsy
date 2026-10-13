package handoff

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func testDocDB(t *testing.T) (*cdl.EIRDoc, string) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "handoff/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

func TestEvaluateCandidateChecksEmptyState(t *testing.T) {
	doc, db := testDocDB(t)
	checks, err := EvaluateCandidateChecks(db, doc, Options{})
	if err != nil {
		t.Fatal(err)
	}
	results := map[string]string{}
	for _, check := range checks {
		results[check.ID] = check.Result
	}
	if len(checks) != 14 {
		t.Fatalf("checks = %d, want 14", len(checks))
	}
	for _, id := range []string{"CANDIDATE-01-EXIT-A", "CANDIDATE-02-SOURCE-CLOSURE", "CANDIDATE-03-SYNTHESIS-DOMAINS", "CANDIDATE-07-USE-CASES", "CANDIDATE-08-CONFIGURATION", "CANDIDATE-14-SNAPSHOT-INPUTS"} {
		if results[id] != "Fail" {
			t.Fatalf("%s = %q, want Fail on empty state", id, results[id])
		}
	}
	if results["CANDIDATE-06-LIFECYCLE"] != "Fail" {
		t.Fatalf("lifecycle = %q, want Fail on empty state", results["CANDIDATE-06-LIFECYCLE"])
	}
	if results["CANDIDATE-13-PENDING-AND-CONTRADICTIONS"] != "Pass" {
		t.Fatalf("pending/contradictions = %q, want Pass on empty state", results["CANDIDATE-13-PENDING-AND-CONTRADICTIONS"])
	}
	if results["CANDIDATE-05-REFERENTIAL-INTEGRITY"] != "Pass" {
		t.Fatalf("referential integrity = %q, want vacuous Pass on empty state", results["CANDIDATE-05-REFERENTIAL-INTEGRITY"])
	}
	for _, id := range []string{"CANDIDATE-04-HANDBOOK-QUALITY", "CANDIDATE-09-PRIVACY-SECURITY", "CANDIDATE-10-DISABLED-CAPABILITIES", "CANDIDATE-11-FAULTS", "CANDIDATE-12-PROFILES"} {
		if results[id] != "Unsupported" {
			t.Fatalf("%s = %q, want Unsupported", id, results[id])
		}
	}
}
