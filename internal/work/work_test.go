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

func invariantProposal(constraint string) Proposal {
	return Proposal{
		FieldID: "INVARIANT",
		Records: []Record{{Payload: []byte("- **Constraint:** " + constraint + "\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")}},
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
		ID: "U-1", Mode: "Discovery", Scope: "demo", RecordID: "INV-NEW", ResultType: "INVARIANT",
		Reads: []Ref{{FieldID: "INVARIANT", RecordID: "INV-DEP"}},
	}
	var sawContext bool
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		text := string(inv.Context.Bytes)
		sawContext = strings.Contains(text, "INV-DEP") && !strings.Contains(text, "INV-OTHER")
		if len(inv.Context.Closure) != 1 {
			t.Errorf("closure = %v", inv.Context.Closure)
		}
		return invariantProposal("proposed"), nil
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
		ID: "U-STALE", Mode: "Discovery", RecordID: "INV-STALE", ResultType: "INVARIANT",
		Reads: []Ref{{FieldID: "INVARIANT", RecordID: "INV-DEP"}},
	}
	attempt := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		attempt++
		mutateInvariant(t, db, "INV-DEP", fmt.Sprintf("changed-%d", attempt))
		return invariantProposal("proposed"), nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err == nil {
		t.Fatal("stale commit accepted")
	}
	if report.Units[0].Retries != 4 {
		t.Fatalf("retries = %d, want 4", report.Units[0].Retries)
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
	err := Commit(caps, Unit{ID: "U", Mode: "M", RecordID: "INV-NEW2", ResultType: "INVARIANT"}, invariantProposal("p"))
	var stale *StaleError
	if !errors.As(err, &stale) {
		t.Fatalf("err = %v, want StaleError", err)
	}
}

func TestReviewPolicyAndRetry(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-REVIEW", Mode: "Synthesis", RecordID: "INV-REVIEW", ResultType: "INVARIANT", RequiresReview: true}
	executions := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		executions++
		return invariantProposal("proposed"), nil
	})
	reviews := 0
	reviewer := funcReviewer(func(request ReviewRequest) (ReviewResult, error) {
		reviews++
		if !strings.Contains(string(request.Context.Bytes), "Proposed record") {
			t.Error("review context lacks the proposed record")
		}
		if len(request.Targets) != 1 || request.Targets[0] != (Ref{FieldID: "INVARIANT", RecordID: "INV-REVIEW"}) {
			t.Errorf("review targets = %v", request.Targets)
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
		return invariantProposal("resolved"), nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{{ID: "U-BLOCK", Mode: "Discovery", RecordID: "INV-BLOCK", ResultType: "INVARIANT"}})
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
	report, err = session.Run([]Unit{{ID: "U-HUMAN", Mode: "Discovery", RecordID: "INV-HUMAN", ResultType: "INVARIANT"}})
	if err == nil {
		t.Fatal("human-authority blocker did not stop the run")
	}
	if !report.Units[0].Escalated || !ReachesOperator(report.Units[0].Blocker) {
		t.Fatalf("stats = %+v", report.Units[0])
	}
}

func TestUnitValidation(t *testing.T) {
	_, doc := testDoc(t)
	if err := (Unit{ID: "U", Mode: "M", RecordID: "R", ResultType: "NOPE"}).Validate(doc); err == nil {
		t.Fatal("unknown result type accepted")
	}
	if err := (Unit{ID: "U", Mode: "M", RecordID: "R", ResultType: "ENTITY-COLUMN"}).Validate(doc); err == nil {
		t.Fatal("non-table result type accepted")
	}
	if err := (Unit{ID: "U", Mode: "M", RecordID: "R", ResultType: "INVARIANT"}).Validate(doc); err != nil {
		t.Fatal(err)
	}
	if err := (Unit{ID: "U", Mode: "M", ResultType: "INVARIANT"}).Validate(doc); err == nil {
		t.Fatal("unit without a target record ID accepted")
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

func TestBatchProposalDerivesIDs(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-BATCH", Mode: "Preflight", IDFromField: "constraint", ResultType: "INVARIANT"}
	proposal := Proposal{FieldID: "INVARIANT", Records: []Record{
		{Payload: []byte("- **Constraint:** ROW-1\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")},
		{Payload: []byte("- **Constraint:** ROW-2\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")},
	}}
	if err := Commit(NewCaps(db, doc), unit, proposal); err != nil {
		t.Fatal(err)
	}
	ids, err := store.List(db, doc, "INVARIANT")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "ROW-1,ROW-2" {
		t.Fatalf("ids = %v", ids)
	}
	duplicate := Proposal{FieldID: "INVARIANT", Records: []Record{
		{Payload: []byte("- **Constraint:** ROW-3\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")},
		{Payload: []byte("- **Constraint:** ROW-3\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")},
	}}
	if err := Commit(NewCaps(db, doc), unit, duplicate); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate batch accepted: %v", err)
	}
	if _, err := store.Get(db, doc, "INVARIANT", "ROW-3"); err == nil {
		t.Fatal("partial batch committed")
	}
}

func TestValidationFailureRetriesWithFeedback(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-FEEDBACK", Mode: "Discovery", RecordID: "INV-FEEDBACK", ResultType: "INVARIANT"}
	attempts := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		attempts++
		if attempts == 1 {
			if err := Commit(NewCaps(db, doc), unit, Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte("- **Constraint:** c\n- **Unknown:** x\n- **Discovered In:** SRC-1\n")}}}); err == nil {
				t.Fatal("invalid payload committed")
			}
			return Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte("- **Constraint:** c\n- **Unknown:** x\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n")}}}, nil
		}
		if !strings.Contains(string(inv.Context.Bytes), "Previous attempt rejected") {
			t.Fatal("retry context lacks rejection feedback")
		}
		return invariantProposal("corrected"), nil
	})
	session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]Unit{unit})
	if err != nil {
		t.Fatal(err)
	}
	if report.Units[0].Retries != 1 || attempts != 2 {
		t.Fatalf("stats = %+v attempts = %d", report.Units[0], attempts)
	}
}

func TestObserverLifecycleEvents(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-OBS", Mode: "Discovery", RecordID: "INV-OBS", ResultType: "INVARIANT", RequiresReview: true}
	attempts := 0
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		attempts++
		if attempts == 1 {
			return Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte("- **Constraint:** c\n- **Unknown:** x\n- **Discovered In:** SRC-1\n")}}}, nil
		}
		return invariantProposal("observed"), nil
	})
	reviewer := funcReviewer(func(request ReviewRequest) (ReviewResult, error) {
		return ReviewResult{Outcome: ReviewApproved}, nil
	})
	config := DefaultConfig()
	var events []string
	config.Observer = &Observer{
		UnitStarted: func(work Unit) { events = append(events, "start:"+work.ID) },
		Reviewing:   func(Unit) { events = append(events, "review") },
		Retrying:    func(Unit, string) { events = append(events, "retry") },
		UnitDone: func(stats UnitStats) {
			events = append(events, fmt.Sprintf("done:%s:%v", stats.UnitID, stats.Committed))
		},
	}
	session, err := NewSession(db, doc, config, executor, reviewer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run([]Unit{unit}); err != nil {
		t.Fatal(err)
	}
	want := []string{"start:U-OBS", "review", "retry", "review", "done:U-OBS:true"}
	if strings.Join(events, ",") != strings.Join(want, ",") {
		t.Fatalf("events = %v, want %v", events, want)
	}
}

func TestClusterReferenceValidation(t *testing.T) {
	db, doc := testDoc(t)
	clusterPayload := func(owner string) []byte {
		return []byte("- **cluster-id:** WEB\n" +
			"- **cluster-title:** Web UI\n" +
			"- **concrete-entry-coordinates:** web/index.js\n" +
			"- **canonical-owner-persona:** " + owner + "\n" +
			"- **primary-pov:** POV-1\n" +
			"- **traversal-track:** Normal\n" +
			"- **capability-state:** Active\n" +
			"- **environment-snapshot:** dev\n" +
			"- **source-ids:** None\n" +
			"- **max-depth:** 3\n" +
			"- **traversal-status:** [R-UNSWEPT]\n" +
			"- **required-traversal-matrix:** view\n")
	}
	unit := Unit{ID: "U-CL", Mode: "Preflight", RecordID: "WEB", ResultType: "ENTRY-CLUSTER"}
	err := Commit(NewCaps(db, doc), unit, Proposal{FieldID: "ENTRY-CLUSTER", Records: []Record{{Payload: clusterPayload("NOPE")}}})
	if err == nil || !strings.Contains(err.Error(), "does not reference a committed persona") {
		t.Fatalf("unknown owner accepted: %v", err)
	}
	persona := []byte("- **persona:** Web UI\n- **prefix:** WEB\n- **persona-slug:** web-ui\n- **persona-directory:** WEB-web-ui\n- **actor-class:** human\n- **canonical-entry-ids:** web/index.js\n- **scope-paths:** web/\n- **registry-state:** Active\n")
	if err := store.Ingest(db, doc, "PERSONA-REGISTRY", "WEB", 1, persona); err != nil {
		t.Fatal(err)
	}
	if err := Commit(NewCaps(db, doc), unit, Proposal{FieldID: "ENTRY-CLUSTER", Records: []Record{{Payload: clusterPayload("WEB")}}}); err != nil {
		t.Fatalf("committed owner rejected: %v", err)
	}
	if err := Commit(NewCaps(db, doc), unit, Proposal{FieldID: "ENTRY-CLUSTER", Records: []Record{{Payload: []byte(strings.Replace(string(clusterPayload("WEB")), "- **source-ids:** None", "- **source-ids:** passport/index.js", 1))}}}); err == nil {
		t.Fatal("bad source token accepted")
	}
}

func TestInvocationIdentityAndLog(t *testing.T) {
	db, doc := testDoc(t)
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		return invariantProposal("logged"), nil
	})
	run := func(id, recordID string) UnitStats {
		session, err := NewSession(db, doc, DefaultConfig(), executor, nil)
		if err != nil {
			t.Fatal(err)
		}
		report, err := session.Run([]Unit{{ID: id, Mode: "Discovery", RecordID: recordID, ResultType: "INVARIANT"}})
		if err != nil {
			t.Fatal(err)
		}
		return report.Units[0]
	}
	first := run("U-I1", "INV-I1")
	second := run("U-I2", "INV-I2")
	if first.InvocationID != "ALFA-NONE-NONE-001" || second.InvocationID != "ALFA-NONE-NONE-002" {
		t.Fatalf("invocation IDs = %q %q", first.InvocationID, second.InvocationID)
	}
	logIDs, err := store.List(db, doc, "INVOCATION-LOG")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(logIDs, ",") != "ALFA-NONE-NONE-001,ALFA-NONE-NONE-002" {
		t.Fatalf("log IDs = %v", logIDs)
	}
	payload, err := store.Get(db, doc, "INVOCATION-LOG", "ALFA-NONE-NONE-001")
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{"- **src-cmp-clm-ids-updated:** INV-I1", "- **mutation-summary:** committed 1 record(s)", "mode=Discovery"} {
		if !strings.Contains(string(payload), marker) {
			t.Fatalf("log payload missing %q:\n%s", marker, payload)
		}
	}
	third := run("U-I3", "INV-I3")
	if third.InvocationID != "ALFA-NONE-NONE-003" {
		t.Fatalf("sequence not continued: %q", third.InvocationID)
	}
}

func TestReRunIsIdempotent(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-IDEM", Mode: "Discovery", RecordID: "INV-IDEM", ResultType: "INVARIANT"}
	proposal := invariantProposal("stable")
	for i := 0; i < 2; i++ {
		session, err := NewSession(db, doc, DefaultConfig(), funcExecutor(func(inv Invocation) (Proposal, error) {
			return proposal, nil
		}), nil)
		if err != nil {
			t.Fatal(err)
		}
		report, err := session.Run([]Unit{unit})
		if err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
		if !report.Units[0].Committed {
			t.Fatalf("run %d not committed: %+v", i+1, report.Units[0])
		}
		if i == 1 && len(report.Units[0].MutationIDs) != 0 {
			t.Fatalf("identical re-run mutated records: %v", report.Units[0].MutationIDs)
		}
	}
	ids, err := store.List(db, doc, "INVARIANT")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "INV-IDEM" {
		t.Fatalf("ids = %v", ids)
	}
	conflict := Unit{ID: "U-IDEM", Mode: "Discovery", RecordID: "INV-IDEM", ResultType: "INVARIANT"}
	session, err := NewSession(db, doc, DefaultConfig(), funcExecutor(func(inv Invocation) (Proposal, error) {
		return invariantProposal("changed"), nil
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run([]Unit{conflict}); err == nil {
		t.Fatal("conflicting re-run accepted")
	}
}
