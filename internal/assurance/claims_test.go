package assurance

import (
	"fmt"
	"os"
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

func testTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "api.js"), []byte("app.get('/api/menuler', handler)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

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
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "assurance/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	header := []byte("- **system-namespace:** demo\n- **target-repository-evidence-roots:** .\n- **primary-technology-stack:** JavaScript\n- **protocol-version:** 4.2.1\n- **canonical-storage-binding:** sqlite/1\n- **current-iteration:** ALFA\n- **default-max-traversal-depth:** 3\n- **included-environments-snapshots:** dev\n- **explicit-scope-policy:** repo\n")
	if err := store.Ingest(db, res.EIR, "PREFLIGHT-HEADER", "0A", 1, header); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

func seedInventory(t *testing.T, db string, doc *cdl.EIRDoc, kind, coordinate, fingerprint string) string {
	t.Helper()
	id, err := capabilities.SourceID("demo", kind, coordinate)
	if err != nil {
		t.Fatal(err)
	}
	payload := "- **kind:** " + kind + "\n" +
		"- **concrete-coordinate:** " + coordinate + "\n" +
		"- **parent-container-src:** None\n" +
		"- **source-fingerprint-snapshot-src-version:** " + fingerprint + " version 1\n" +
		"- **discovery-method:** filesystem\n" +
		"- **scope-summary:** WEB\n" +
		"- **environment-runtime-state:** dev\n" +
		"- **traversal-tracks:** Normal\n" +
		"- **canonical-owner:** WEB\n" +
		"- **disposition-summary:** Unknown\n" +
		"- **mappings:** None\n" +
		"- **frontier-references:** None\n" +
		"- **coverage-rows:** None\n" +
		"- **coverage-summary:** C-GAP\n" +
		"- **claims:** None\n"
	if err := store.Ingest(db, doc, "SOURCE-INVENTORY-RECORD", id, 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	return id
}

type fakeExecutor func(work.Invocation) (work.Proposal, error)

func (f fakeExecutor) Execute(inv work.Invocation) (work.Proposal, error) { return f(inv) }

func TestLoadInventoryAndPlanClaims(t *testing.T) {
	doc, db := testDocDB(t)
	root := testTree(t)
	first := seedInventory(t, db, doc, "FILE", "api.js", "sha256:aaa")
	second := seedInventory(t, db, doc, "FILE", "web.js", "sha256:bbb")
	inventory, err := LoadInventory(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 2 {
		t.Fatalf("inventory = %+v", inventory)
	}
	fingerprints := map[string]bool{}
	for _, item := range inventory {
		fingerprints[item.Fingerprint] = true
	}
	if !fingerprints["sha256:aaa"] || !fingerprints["sha256:bbb"] {
		t.Fatalf("inventory fingerprints = %+v", inventory)
	}
	units, err := Plan(root, inventory, Directives{Namespace: "demo", Snapshot: "dev", BatchSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 1 || units[0].ExpectedRecords != 2 || len(units[0].Reads) != 2 {
		t.Fatalf("units = %+v", units[0])
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		context := string(inv.Context.Bytes)
		if !strings.Contains(context, "app.get('/api/menuler'") {
			t.Error("context lacks bounded excerpts")
		}
		proposal := work.Proposal{FieldID: "CLAIM"}
		for _, item := range inventory {
			payload := fmt.Sprintf("- **subject-id-record-version:** %s\n- **field-step:** file\n- **claim:** %s exposes an HTTP entry point\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** %s:1\n- **source-fingerprint-snapshot:** %s\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n",
				item.ID, item.Coordinate, item.Coordinate, item.Fingerprint)
			proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
		}
		return proposal, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run(units); err != nil {
		t.Fatal(err)
	}
	ids, err := store.List(db, doc, "CLAIM")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 {
		t.Fatalf("claims = %v", ids)
	}
	want, _ := capabilities.TypedID("CLM", "demo", first, "api.js:1")
	found := false
	for _, id := range ids {
		if id == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected %s in %v", want, ids)
	}
	_ = second
}

func TestClaimCrossReferenceValidation(t *testing.T) {
	doc, db := testDocDB(t)
	seedInventory(t, db, doc, "FILE", "api.js", "sha256:aaa")
	id, _ := capabilities.SourceID("demo", "FILE", "api.js")
	unit := work.Unit{ID: "U-CLAIM", Mode: "Discovery", IDType: "CLM", Namespace: "demo", IDOwnerField: "subject-id-record-version", IDDiscriminatorField: "claim", ResultType: "CLAIM"}
	claim := func(fingerprint, anchor string) []byte {
		return []byte("- **subject-id-record-version:** " + id + "\n- **field-step:** file\n- **claim:** fact\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** " + anchor + "\n- **source-fingerprint-snapshot:** " + fingerprint + "\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n")
	}
	err := work.Commit(work.NewCaps(db, doc), unit, work.Proposal{FieldID: "CLAIM", Records: []work.Record{{Payload: claim("sha256:wrong", "api.js:1")}}})
	if err == nil || !strings.Contains(err.Error(), "does not match source unit") {
		t.Fatalf("wrong fingerprint accepted: %v", err)
	}
	err = work.Commit(work.NewCaps(db, doc), unit, work.Proposal{FieldID: "CLAIM", Records: []work.Record{{Payload: claim("sha256:aaa", "other.js:1")}}})
	if err == nil || !strings.Contains(err.Error(), "must start at the source unit coordinate") {
		t.Fatalf("wrong anchor accepted: %v", err)
	}
	if err := work.Commit(work.NewCaps(db, doc), unit, work.Proposal{FieldID: "CLAIM", Records: []work.Record{{Payload: claim("sha256:aaa", "api.js:1")}}}); err != nil {
		t.Fatalf("valid claim rejected: %v", err)
	}
}

func TestExpandCoverageAndSummary(t *testing.T) {
	doc, db := testDocDB(t)
	first := seedInventory(t, db, doc, "FILE", "api.js", "sha256:aaa")
	second := seedInventory(t, db, doc, "FILE", "web.js", "sha256:bbb")
	claimPayload := "- **subject-id-record-version:** " + first + "\n- **field-step:** file\n- **claim:** api exposes a route\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** api.js:1\n- **source-fingerprint-snapshot:** sha256:aaa\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n"
	claimID, err := capabilities.TypedID("CLM", "demo", first, "api.js:1")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Ingest(db, doc, "CLAIM", claimID, 1, []byte(claimPayload)); err != nil {
		t.Fatal(err)
	}
	created, err := ExpandCoverage(db, doc, CoverageDirectives{Namespace: "demo", Environment: "dev", Snapshot: "snap-1"})
	if err != nil {
		t.Fatal(err)
	}
	if created != 2 {
		t.Fatalf("created = %d, want 2", created)
	}
	again, err := ExpandCoverage(db, doc, CoverageDirectives{Namespace: "demo", Environment: "dev", Snapshot: "snap-1"})
	if err != nil || again != 0 {
		t.Fatalf("second expand = %d err = %v", again, err)
	}
	covIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		t.Fatal(err)
	}
	if len(covIDs) != 2 {
		t.Fatalf("coverage rows = %v", covIDs)
	}
	found := false
	for _, id := range covIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(payload), claimID) {
			found = true
			if !strings.Contains(string(payload), "- **disposition:** Unmapped") {
				t.Fatalf("row not fail-inclusive: %s", payload)
			}
		}
	}
	if !found {
		t.Fatal("claim not bound into any coverage row")
	}
	summary, err := CoverageSummary(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if summary.SourcesByKind["FILE"] != 2 || summary.CoverageByDisposition["Unmapped"] != 2 {
		t.Fatalf("summary = %+v", summary)
	}
	_ = second
}

func TestCellSweepsAndExitAGate(t *testing.T) {
	doc, db := testDocDB(t)
	seedInventory(t, db, doc, "FILE", "api.js", "sha256:aaa")
	seedInventory(t, db, doc, "FILE", "web.js", "sha256:bbb")
	if _, err := ExpandCoverage(db, doc, CoverageDirectives{Namespace: "demo", Environment: "dev", Snapshot: "snap-1"}); err != nil {
		t.Fatal(err)
	}
	cluster := "- **cluster-id:** CL-1\n- **cluster-title:** Web entry\n- **concrete-entry-coordinates:** api.js\n- **canonical-owner-persona:** WEB\n- **primary-pov:** POV-1\n- **traversal-track:** Normal\n- **capability-state:** Active\n- **environment-snapshot:** dev\n- **source-ids:** None\n- **max-depth:** 3\n- **traversal-status:** [R-UNSWEPT]\n- **required-traversal-matrix:** None\n"
	if err := store.Ingest(db, doc, "ENTRY-CLUSTER", "CL-1", 1, []byte(cluster)); err != nil {
		t.Fatal(err)
	}
	matrixRow := "- **cluster-id:** CL-1\n- **persona-owner:** WEB\n- **concrete-entry:** api.js\n- **track:** Normal\n- **environment:** dev\n- **snapshot:** snap-1\n- **pov:** POV-1\n- **applicability:** Required\n- **status:** [R-UNSWEPT]\n- **swp-id:** None\n- **claims:** None\n"
	if err := store.Ingest(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", "CL-1|WEB|POV-1|Normal|dev|snap-1", 1, []byte(matrixRow)); err != nil {
		t.Fatal(err)
	}
	scopes, err := LoadCellScopes(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 || len(scopes[0].SourceIDs) != 2 || scopes[0].MappedAll {
		t.Fatalf("scopes = %+v", scopes)
	}
	session, err := work.NewSession(db, doc, work.DefaultConfig(), MechanicalCellSweepExecutor{Scopes: scopes}, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]work.Unit{CellSweepPlan(scopes)})
	if err != nil {
		t.Fatal(err)
	}
	if !report.Units[0].Committed {
		t.Fatalf("cell sweep not committed: %+v", report.Units[0])
	}
	payload, err := store.Get(db, doc, "SWEEP-RECORD", report.Units[0].MutationIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{"- **cluster-id:** CL-1", "- **pov:** POV-1", "- **source-denominator:** 2", "- **coverage-result:** C-PARTIAL"} {
		if !strings.Contains(string(payload), marker) {
			t.Fatalf("sweep payload missing %q:\n%s", marker, payload)
		}
	}
	closed, err := AdvanceMatrix(db, doc)
	if err != nil || closed != 0 {
		t.Fatalf("closed = %d err = %v", closed, err)
	}
	gate, err := EvaluateExitA(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	results := map[string]string{}
	for _, check := range gate.Checks {
		results[check.ID] = check.Result
	}
	if results["COVERAGE-DENOMINATOR"] != "Fail" || results["TRAVERSAL-MATRIX-CLOSURE"] != "Fail" || results["SYNTHESIS-COVERAGE"] != "Fail" {
		t.Fatalf("gate checks = %+v", gate.Checks)
	}
}
