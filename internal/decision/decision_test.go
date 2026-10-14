package decision

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func testDocAndDB(t *testing.T) (*cdl.EIRDoc, string) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "decision/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

type fakeCompleter struct{ content string }

func (f fakeCompleter) Complete(ctx context.Context, system, user string) (string, error) {
	return f.content, nil
}

func seedSubjects(t *testing.T, db string, doc *cdl.EIRDoc) {
	t.Helper()
	header := []byte("- **system-namespace:** demo\n- **target-repository-evidence-roots:** app\n- **primary-technology-stack:** Python\n- **protocol-version:** 4.2.1\n- **canonical-storage-binding:** sqlite/1\n- **current-iteration:** ALFA\n- **default-max-traversal-depth:** 3\n- **included-environments-snapshots:** dev\n- **explicit-scope-policy:** repo\n")
	if err := store.Ingest(db, doc, "PREFLIGHT-HEADER", "0A", 1, header); err != nil {
		t.Fatal(err)
	}
	claim := []byte("- **subject-id-record-version:** MOD-AAAAAAAAAAAA version 1\n- **field-step:** responsibility-domain-boundary\n- **claim:** The module validates every inbound request.\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** app/a.py:1\n- **source-fingerprint-snapshot:** sha256:aaa\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independence-basis:** single static route\n- **created-last-revalidated-by:** ALFA-NONE-NONE-001\n")
	if err := store.Ingest(db, doc, "CLAIM", "CLM-AAAAAAAAAAAA", 1, claim); err != nil {
		t.Fatal(err)
	}
	module := []byte("- **source-bindings:** CLM-AAAAAAAAAAAA\n- **block-id:** MOD-AAAAAAAAAAAA\n- **record-title:** Request validation\n- **domain-path:** demo/core\n- **responsibility-domain-boundary:** Validates requests.\n- **public-interfaces:** None\n- **encapsulated-state-entities:** None\n- **dependencies-dependents:** None\n- **data-and-events-passed:** None\n- **invariants-rules-owned:** Unknown\n- **personas-use-cases-capabilities:** None\n- **failure-propagation:** Unknown\n- **observed-legacy-placement:** inline\n- **target-decision:** None\n")
	if err := store.Ingest(db, doc, "ARCHITECTURE-MODULE", "MOD-AAAAAAAAAAAA", 1, module); err != nil {
		t.Fatal(err)
	}
	contradiction := []byte("- **claims-in-conflict:** CLM-AAAAAAAAAAAA\n- **scope-environments-snapshots:** dev\n- **conflict-description:** Validation order is unclear.\n- **resolution-state:** Decision-Required\n- **affected-records-and-stale-impact:** MOD-AAAAAAAAAAAA\n")
	if err := store.Ingest(db, doc, "CONTRADICTION", "CON-BBBBBBBBBBBB", 1, contradiction); err != nil {
		t.Fatal(err)
	}
}

const draftContent = `## Options

Keep the inline validation or extract a shared validator.

## Decision

Change-By-Approved-Migration

## Rationale

The observed claim shows one shared validation path.

## Risks

Unknown callers may bypass the extracted validator.

## Exact Scope and Denominator Effect

Affects MOD-AAAAAAAAAAAA coverage only.

## Affected Records / Rollout / Rollback

MOD-AAAAAAAAAAAA; staged rollout; rollback restores inline validation.`

func TestTargetsAndDraftRoundTrip(t *testing.T) {
	doc, db := testDocAndDB(t)
	seedSubjects(t, db, doc)
	targets, err := Targets(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 {
		t.Fatalf("targets = %+v", targets)
	}
	if targets[0].DecisionType != "contradiction" || targets[0].SubjectID != "CON-BBBBBBBBBBBB" {
		t.Fatalf("target[0] = %+v", targets[0])
	}
	if targets[1].DecisionType != "target" || targets[1].SubjectID != "MOD-AAAAAAAAAAAA" {
		t.Fatalf("target[1] = %+v", targets[1])
	}
	if len(targets[1].ObservedLines) != 2 || !strings.Contains(targets[1].ObservedLines[1], "CLAIM/CLM-AAAAAAAAAAAA") {
		t.Fatalf("observed inputs = %v", targets[1].ObservedLines)
	}
	drafts, err := DraftAll(context.Background(), fakeCompleter{content: draftContent}, db, doc)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := RenderDraft(drafts, "2026-10-14T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := ParseDraft(rendered)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %+v", entries)
	}
	if entries[0].SubjectID != "CON-BBBBBBBBBBBB" || entries[1].SubjectID != "MOD-AAAAAAAAAAAA" {
		t.Fatalf("entry order = %+v", entries)
	}
	if entries[1].Sections["Decision"] != "Change-By-Approved-Migration" {
		t.Fatalf("decision section = %q", entries[1].Sections["Decision"])
	}
}

func TestReviewAndCommitApproved(t *testing.T) {
	doc, db := testDocAndDB(t)
	seedSubjects(t, db, doc)
	drafts, err := DraftAll(context.Background(), fakeCompleter{content: draftContent}, db, doc)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := RenderDraft(drafts, "2026-10-14T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := ParseDraft(rendered)
	if err != nil {
		t.Fatal(err)
	}
	answers := []string{"operator", "every drafted decision", "content review", "r", "conflict needs field data", "a", "approve the module decision"}
	next := 0
	prompt := func(question string) (string, error) {
		answer := answers[next]
		next++
		return answer, nil
	}
	review, err := CollectReview(prompt, &strings.Builder{}, entries, Review{}, "2026-10-14T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if review.Result != "Incomplete" || len(review.Rows) != 2 {
		t.Fatalf("review = %+v", review)
	}
	renderedReview := RenderReview(review)
	parsedReview, err := ParseReview(renderedReview)
	if err != nil {
		t.Fatal(err)
	}
	if parsedReview.Reviewer != "operator" || parsedReview.Rows[0].Verdict != "Reject" || parsedReview.Rows[1].Verdict != "Approve" {
		t.Fatalf("parsed review = %+v", parsedReview)
	}
	created, err := CommitApproved(db, doc, "demo", entries, parsedReview)
	if err != nil {
		t.Fatal(err)
	}
	wantID, err := capabilities.TypedID("DEC", "demo", "MOD-AAAAAAAAAAAA", "target")
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0] != wantID {
		t.Fatalf("created = %v, want %s", created, wantID)
	}
	payload, err := store.Get(db, doc, "DECISION", wantID)
	if err != nil {
		t.Fatal(err)
	}
	values, err := work.PayloadText([]string{"decision-content-fingerprint", "approval-status", "approved-content-binding", "owner-approver-authority-basis-approval-timestamp"}, payload)
	if err != nil {
		t.Fatal(err)
	}
	if values["approval-status"] != "Approved" {
		t.Fatalf("approval status = %q", values["approval-status"])
	}
	fingerprint, err := capabilities.DecisionContentFingerprint(payload, "decision-content-fingerprint")
	if err != nil {
		t.Fatal(err)
	}
	if values["decision-content-fingerprint"] != fingerprint {
		t.Fatalf("stored fingerprint %s does not match recomputation %s", values["decision-content-fingerprint"], fingerprint)
	}
	if !strings.Contains(values["approved-content-binding"], fingerprint) || !strings.Contains(values["owner-approver-authority-basis-approval-timestamp"], "operator") {
		t.Fatalf("binding/approval fields = %+v", values)
	}
}
