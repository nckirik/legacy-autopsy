package work

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func testDoc(t *testing.T) (string, *cdl.EIRDoc) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "work/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return db, res.EIR
}

func mutateInvariant(t *testing.T, db, id, constraint string) {
	t.Helper()
	sql := "update invariant set la_semantic_fingerprint = 'sha256:' || '" + constraint + "', la_canonical_payload = '- **Constraint:** " + constraint + "' where la_record_id = '" + id + "';"
	cmd := exec.Command("sqlite3", "--bail", db)
	cmd.Stdin = strings.NewReader(sql)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mutate %s: %v: %s", id, err, out)
	}
}

func seedInvariant(t *testing.T, db string, doc *cdl.EIRDoc, id, constraint string) {
	t.Helper()
	payload := "- **Constraint:** " + constraint + "\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n"
	if err := store.Ingest(db, doc, "INVARIANT", id, 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func invariantProposal(id, constraint string) Proposal {
	return Proposal{
		FieldID:  "INVARIANT",
		RecordID: id,
		Version:  1,
		Payload:  []byte("- **Constraint:** " + constraint + "\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n"),
	}
}

type funcExecutor func(Invocation) (Proposal, error)

func (f funcExecutor) Execute(inv Invocation) (Proposal, error) { return f(inv) }

type funcReviewer func(ReviewRequest) (ReviewResult, error)

func (f funcReviewer) Review(request ReviewRequest) (ReviewResult, error) { return f(request) }

func TestSessionCommitsBoundedUnit(t *testing.T) {
	db, doc := testDoc(t)
	seedInvariant(t, db, doc, "INV-DEP", "dependency")
	seedInvariant(t, db, doc, "INV-OTHER", "unrelated record")
	unit := Unit{
		ID: "U-1", Mode: "Discovery", Scope: "demo", ResultType: "INVARIANT",
		Reads: []Ref{{FieldID: "INVARIANT", RecordID: "INV-DEP"}},
	}
	var sawContext bool
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		text := string(inv.Context.Bytes)
		sawContext = strings.Contains(text, "INV-DEP") && !strings.Contains(text, "INV-OTHER")
		if len(inv.Context.Closure) != 1 {
			t.Errorf("closure = %v", inv.Context.Closure)
		}
		return invariantProposal("INV-NEW", "proposed"), nil
	})
	session, err := NewSession(db, doc, Config{ReviewMode: ReviewCritical, OperatorMinimums: OperatorMinimums{DecisionBatchMinutes: 10}}, executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err != nil {
		t.Fatal(err)
	}
	if !sawContext {
		t.Fatal("context materialized records outside the declared closure")
	}
	if len(report.Units) != 1 {
		t.Fatalf("report = %+v", report)
	}
	stats := report.Units[0]
	if stats.ContextBytes == 0 || stats.Dependencies != 1 || stats.Fetches != 1 || stats.Retries != 0 {
		t.Fatalf("stats = %+v", stats)
	}
	if _, err := store.Get(db, doc, "INVARIANT", "INV-NEW"); err != nil {
		t.Fatal(err)
	}
}

func TestCommitRejectsStaleDependency(t *testing.T) {
	db, doc := testDoc(t)
	seedInvariant(t, db, doc, "INV-DEP", "original")
	unit := Unit{
		ID: "U-STALE", Mode: "Discovery", ResultType: "INVARIANT",
		Reads: []Ref{{FieldID: "INVARIANT", RecordID: "INV-DEP"}},
	}
	attempt := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		attempt++
		mutateInvariant(t, db, "INV-DEP", fmt.Sprintf("changed-%d", attempt))
		return invariantProposal("INV-STALE", "proposed"), nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err == nil {
		t.Fatal("stale commit accepted")
	}
	if report.Units[0].Retries != 3 {
		t.Fatalf("retries = %d, want 3", report.Units[0].Retries)
	}
	if _, err := store.Get(db, doc, "INVARIANT", "INV-STALE"); err == nil {
		t.Fatal("stale proposal was committed")
	}
}

func TestStaleErrorType(t *testing.T) {
	db, doc := testDoc(t)
	seedInvariant(t, db, doc, "INV-DEP", "original")
	caps := NewCaps(db, doc)
	if _, err := caps.Get(Ref{FieldID: "INVARIANT", RecordID: "INV-DEP"}); err != nil {
		t.Fatal(err)
	}
	mutateInvariant(t, db, "INV-DEP", "changed after read")
	err := Commit(caps, Unit{ID: "U", Mode: "M", ResultType: "INVARIANT"}, invariantProposal("INV-NEW2", "p"))
	var stale *StaleError
	if !errors.As(err, &stale) {
		t.Fatalf("err = %v, want StaleError", err)
	}
}

func TestReviewPolicyAndRetry(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-REVIEW", Mode: "Synthesis", ResultType: "INVARIANT", RequiresReview: true}
	executions := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		executions++
		return invariantProposal("INV-REVIEW", "proposed"), nil
	})
	reviews := 0
	reviewer := funcReviewer(func(request ReviewRequest) (ReviewResult, error) {
		reviews++
		if !strings.Contains(string(request.Context.Bytes), "Proposed record") {
			t.Error("review context lacks the proposed record")
		}
		if request.Target != (Ref{FieldID: "INVARIANT", RecordID: "INV-REVIEW"}) {
			t.Errorf("review target = %v", request.Target)
		}
		if reviews == 1 {
			return ReviewResult{Outcome: ReviewRejected, Findings: []string{"unsupported"}}, nil
		}
		return ReviewResult{Outcome: ReviewApproved}, nil
	})
	session, err := NewSession(db, doc, Config{ReviewMode: ReviewCritical, OperatorMinimums: OperatorMinimums{DecisionBatchMinutes: 10}}, executor, reviewer)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err != nil {
		t.Fatal(err)
	}
	if executions != 2 || reviews != 2 {
		t.Fatalf("executions=%d reviews=%d", executions, reviews)
	}
	if report.Units[0].Review != ReviewApproved || report.Units[0].Retries != 1 {
		t.Fatalf("stats = %+v", report.Units[0])
	}
	if _, err := store.Get(db, doc, "INVARIANT", "INV-REVIEW"); err != nil {
		t.Fatal(err)
	}
}

func TestBlockerRouting(t *testing.T) {
	cases := map[string]struct {
		class    BlockerClass
		operator bool
	}{
		"human-authority":   {BlockerHumanAuthority, true},
		"decision-approval": {BlockerHumanAuthority, true},
		"uncertainty":       {BlockerUncertainty, false},
		"harness":           {BlockerHarness, false},
		"weird-unknown":     {BlockerProtocol, false},
	}
	for kind, want := range cases {
		class := ClassifyBlocker(Blocker{Kind: kind})
		if class != want.class || ReachesOperator(class) != want.operator {
			t.Errorf("kind %q: class=%s operator=%v", kind, class, ReachesOperator(class))
		}
	}
	db, doc := testDoc(t)
	attempts := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		attempts++
		if attempts <= 2 {
			return Proposal{}, &BlockerError{Blocker: Blocker{Kind: "uncertainty", Message: "not sure"}}
		}
		return invariantProposal("INV-BLOCK", "resolved"), nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{{ID: "U-BLOCK", Mode: "Discovery", ResultType: "INVARIANT"}})
	if err != nil {
		t.Fatal(err)
	}
	if report.Units[0].Retries != 2 {
		t.Fatalf("stats = %+v", report.Units[0])
	}
	human := funcExecutor(func(inv Invocation) (Proposal, error) {
		return Proposal{}, &BlockerError{Blocker: Blocker{Kind: "scope-approval", Message: "needs scope decision"}}
	})
	session, err = NewSession(db, doc, DefaultConfig(), human, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err = session.Run([]Unit{{ID: "U-HUMAN", Mode: "Discovery", ResultType: "INVARIANT"}})
	if err == nil {
		t.Fatal("human-authority blocker did not stop the run")
	}
	if !report.Units[0].Escalated || !ReachesOperator(report.Units[0].Blocker) {
		t.Fatalf("stats = %+v", report.Units[0])
	}
}

func TestUnitValidation(t *testing.T) {
	_, doc := testDoc(t)
	if err := (Unit{ID: "U", Mode: "M", ResultType: "NOPE"}).Validate(doc); err == nil {
		t.Fatal("unknown result type accepted")
	}
	if err := (Unit{ID: "U", Mode: "M", ResultType: "ENTITY-COLUMN"}).Validate(doc); err == nil {
		t.Fatal("non-table result type accepted")
	}
	if err := (Unit{ID: "U", Mode: "M", ResultType: "INVARIANT"}).Validate(doc); err != nil {
		t.Fatal(err)
	}
}

func TestAskConfig(t *testing.T) {
	answers := map[string]string{}
	config, err := AskConfig(func(q ConfigQuestion, fallback string) (string, error) {
		return answers[q.Prompt], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if config != DefaultConfig() {
		t.Fatalf("empty answers = %+v", config)
	}
	answers["Independent review policy"] = "all"
	answers["Target minutes for a 50-decision review batch"] = "20"
	answers["Plain-language questions and decision indexes"] = "no"
	config, err = AskConfig(func(q ConfigQuestion, fallback string) (string, error) {
		return answers[q.Prompt], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.ReviewMode != ReviewAll || config.OperatorMinimums.DecisionBatchMinutes != 20 || config.OperatorMinimums.PlainLanguageQuestions {
		t.Fatalf("config = %+v", config)
	}
	answers["Independent review policy"] = "sometimes"
	if _, err := AskConfig(func(q ConfigQuestion, fallback string) (string, error) {
		return answers[q.Prompt], nil
	}); err == nil {
		t.Fatal("invalid review mode accepted")
	}
}
