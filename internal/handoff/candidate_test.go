package handoff

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
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
	if results["CANDIDATE-04-HANDBOOK-QUALITY"] != "Unsupported" {
		t.Fatalf("handbook = %q, want Unsupported without a rendered handbook", results["CANDIDATE-04-HANDBOOK-QUALITY"])
	}
	if results["CANDIDATE-12-PROFILES"] != "Fail" {
		t.Fatalf("profiles = %q, want Fail with no persona registry rows", results["CANDIDATE-12-PROFILES"])
	}
	if results["CANDIDATE-09-PRIVACY-SECURITY"] != "Fail" {
		t.Fatalf("privacy = %q, want Fail when no findings are recorded", results["CANDIDATE-09-PRIVACY-SECURITY"])
	}
	if results["CANDIDATE-11-FAULTS"] != "Fail" {
		t.Fatalf("faults = %q, want Fail when no fault records exist", results["CANDIDATE-11-FAULTS"])
	}
	if results["CANDIDATE-10-DISABLED-CAPABILITIES"] != "Pass" {
		t.Fatalf("disabled = %q, want Pass on a resolved zero denominator", results["CANDIDATE-10-DISABLED-CAPABILITIES"])
	}
}

func TestBuildCandidateReportEmptyState(t *testing.T) {
	doc, db := testDocDB(t)
	report, err := BuildCandidateReport(db, doc, ReportOptions{
		Namespace:             "demo",
		ProtocolVersion:       "4.6",
		EnvironmentsSnapshots: "dev",
		ArtifactPath:          "exit-e/EXIT-E-CANDIDATE-REPORT.md",
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != "Blocked" {
		t.Fatalf("result = %q, want Blocked on empty state", report.Result)
	}
	if report.Fingerprint == "" || report.InstanceBinding == "" || len(report.Payload) == 0 {
		t.Fatalf("report artifact is incomplete: %+v", report)
	}
	recomputed, err := artifact.PayloadFingerprint(CandidateReportArtifactType, "exit-e/EXIT-E-CANDIDATE-REPORT.md", report.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if recomputed != report.Fingerprint {
		t.Fatalf("fingerprint %s does not recompute to %s", report.Fingerprint, recomputed)
	}
	if err := artifact.ValidatePayload(doc, CandidateReportArtifactType, report.Payload); err != nil {
		t.Fatalf("payload schema: %v", err)
	}
	text := string(report.Payload)
	if !strings.Contains(text, "| CANDIDATE-14-SNAPSHOT-INPUTS |") {
		t.Fatalf("report is missing the candidate check rows")
	}
	if strings.Contains(text, "BLOCKER-CANDIDATE-01-EXIT-A | CANDIDATE-01-EXIT-A | all deterministic") {
		t.Fatalf("passing check produced a blocker")
	}
}
