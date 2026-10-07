package synthesis

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
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
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "synthesis/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

type fakeExecutor func(work.Invocation) (work.Proposal, error)

func (f fakeExecutor) Execute(inv work.Invocation) (work.Proposal, error) { return f(inv) }

func TestPlanInjectsBlockIDAndDerivesSemanticID(t *testing.T) {
	doc, db := testDocDB(t)
	claims := []ClaimRef{
		{ID: "CLM-000000000001", Subject: "SRC-FILE-AAAAAAAAAAAA", Text: "app/routes.js registers the public and API routes"},
		{ID: "CLM-000000000002", Subject: "SRC-FILE-BBBBBBBBBBBB", Text: "the passport strategy authenticates users via SAML"},
	}
	for _, claim := range claims {
		inventory := "- **kind:** FILE\n- **concrete-coordinate:** " + claim.Subject + ".js\n- **parent-container-src:** None\n- **source-fingerprint-snapshot-src-version:** sha256:x version 1\n- **discovery-method:** filesystem\n- **scope-summary:** WEB\n- **environment-runtime-state:** dev\n- **traversal-tracks:** Normal\n- **canonical-owner:** WEB\n- **disposition-summary:** Unknown\n- **mappings:** None\n- **frontier-references:** None\n- **coverage-rows:** None\n- **coverage-summary:** C-GAP\n- **claims:** None\n"
		if err := store.Ingest(db, doc, "SOURCE-INVENTORY-RECORD", claim.Subject, 1, []byte(inventory)); err != nil {
			t.Fatal(err)
		}
		claimPayload := "- **subject-id-record-version:** " + claim.Subject + "\n- **field-step:** file\n- **claim:** " + claim.Text + "\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** x:1\n- **source-fingerprint-snapshot:** sha256:x\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n"
		if err := store.Ingest(db, doc, "CLAIM", claim.ID, 1, []byte(claimPayload)); err != nil {
			t.Fatal(err)
		}
	}
	unit := Plan(Catalogs["MOD"], claims, Directives{Namespace: "demo", Snapshot: "dev"})
	if unit.RuntimeFields[0] != "block-id" || unit.IDFromBindings != "source-bindings" {
		t.Fatalf("unit = %+v", unit)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		context := string(inv.Context.Bytes)
		if !strings.Contains(context, "app/routes.js registers") {
			t.Error("context lacks claims")
		}
		payload := "- **source-bindings:** CLM-000000000001, CLM-000000000002\n" +
			"- **block-id:** PENDING\n" +
			"- **record-title:** Request Routing and Authentication\n" +
			"- **domain-path:** app/routes\n" +
			"- **responsibility-domain-boundary:** HTTP request routing and SAML authentication\n" +
			"- **public-interfaces:** /login, /api/*\n" +
			"- **encapsulated-state-entities:** session state\n" +
			"- **dependencies-dependents:** passport, controllers\n" +
			"- **data-and-events-passed:** HTTP requests\n" +
			"- **failure-propagation:** authentication failures redirect to /login\n" +
			"- **invariants-rules-owned:** authenticated routes require a session\n" +
			"- **observed-legacy-placement:** index.js and app/middlewares\n" +
			"- **target-decision:** Pending\n" +
			"- **personas-use-cases-capabilities:** CAN-1\n"
		return work.Proposal{FieldID: "ARCHITECTURE-MODULE", Records: []work.Record{{Payload: []byte(payload)}}}, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]work.Unit{unit})
	if err != nil {
		t.Fatal(err)
	}
	if !report.Units[0].Committed {
		t.Fatalf("stats = %+v", report.Units[0])
	}
	want, err := capabilities.TypedID("MOD", "demo", "SRC-FILE-AAAAAAAAAAAA", "MOD")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(report.Units[0].MutationIDs, ",") != want {
		t.Fatalf("mutation IDs = %v, want %s", report.Units[0].MutationIDs, want)
	}
	payload, err := store.Get(db, doc, "ARCHITECTURE-MODULE", want)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), "- **block-id:** "+want) {
		t.Fatalf("block-id not injected:\n%s", payload)
	}
}

func TestInjectionFailsWhenCarrierMissing(t *testing.T) {
	doc, db := testDocDB(t)
	unit := Plan(Catalogs["MOD"], nil, Directives{Namespace: "demo", Snapshot: "dev"})
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		payload := "- **source-bindings:** None\n- **record-title:** X\n- **domain-path:** app/x\n- **responsibility-domain-boundary:** x\n- **public-interfaces:** None\n- **encapsulated-state-entities:** None\n- **dependencies-dependents:** None\n- **data-and-events-passed:** None\n- **failure-propagation:** None\n- **invariants-rules-owned:** None\n- **observed-legacy-placement:** None\n- **target-decision:** Pending\n- **personas-use-cases-capabilities:** None\n"
		return work.Proposal{FieldID: "ARCHITECTURE-MODULE", Records: []work.Record{{Payload: []byte(payload)}}}, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run([]work.Unit{unit}); err == nil {
		t.Fatal("missing block-id carrier accepted")
	}
}

func TestBindingsAndCoverageAdvance(t *testing.T) {
	doc, db := testDocDB(t)
	// seed one source unit and one claim
	srcID := "SRC-FILE-AAAAAAAAAAAA"
	inventoryPayload := "- **kind:** FILE\n- **concrete-coordinate:** app/a.js\n- **parent-container-src:** None\n- **source-fingerprint-snapshot-src-version:** sha256:aaa version 1\n- **discovery-method:** filesystem\n- **scope-summary:** WEB\n- **environment-runtime-state:** dev\n- **traversal-tracks:** Normal\n- **canonical-owner:** WEB\n- **disposition-summary:** Unknown\n- **mappings:** None\n- **frontier-references:** None\n- **coverage-rows:** None\n- **coverage-summary:** C-GAP\n- **claims:** None\n"
	if err := store.Ingest(db, doc, "SOURCE-INVENTORY-RECORD", srcID, 1, []byte(inventoryPayload)); err != nil {
		t.Fatal(err)
	}
	claimID := "CLM-000000000001"
	claimPayload := "- **subject-id-record-version:** " + srcID + "\n- **field-step:** file\n- **claim:** api exposes a route\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** app/a.js:1\n- **source-fingerprint-snapshot:** sha256:aaa\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n"
	if err := store.Ingest(db, doc, "CLAIM", claimID, 1, []byte(claimPayload)); err != nil {
		t.Fatal(err)
	}
	// seed a synthesis record binding that source
	modPayload := "- **source-bindings:** " + srcID + "\n- **block-id:** MOD-000000000000\n- **record-title:** API Layer\n- **domain-path:** app/routes\n- **responsibility-domain-boundary:** routing\n- **public-interfaces:** /api\n- **encapsulated-state-entities:** None\n- **dependencies-dependents:** None\n- **data-and-events-passed:** None\n- **failure-propagation:** None\n- **invariants-rules-owned:** None\n- **observed-legacy-placement:** app\n- **target-decision:** Pending\n- **personas-use-cases-capabilities:** None\n"
	if err := store.Ingest(db, doc, "ARCHITECTURE-MODULE", "MOD-000000000000", 1, []byte(modPayload)); err != nil {
		t.Fatal(err)
	}
	if _, err := assuranceExpandForTest(t, db, doc); err != nil {
		t.Fatal(err)
	}
	records, err := LoadSynthesisRecords(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || len(records[0].Bindings) != 1 || records[0].Bindings[0] != srcID {
		t.Fatalf("records = %+v", records)
	}
	runConfig := work.DefaultConfig()
	session, err := work.NewSession(db, doc, runConfig, MechanicalBindingExecutor{
		Records:   records,
		Inventory: map[string]string{srcID: "sha256:aaa version 1"},
		Claims:    map[string]string{claimID: srcID},
		Snapshot:  "dev",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run([]work.Unit{BindingPlan("demo")})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Units[0].MutationIDs) != 1 {
		t.Fatalf("edges = %v", report.Units[0].MutationIDs)
	}
	changed, err := advanceCoverageForTest(t, db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if changed != 1 {
		t.Fatalf("advanced = %d", changed)
	}
	summary, err := summaryForTest(t, db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if summary.CoverageByDisposition["Mapped"] != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}

func assuranceExpandForTest(t *testing.T, db string, doc *cdl.EIRDoc) (int, error) {
	t.Helper()
	return assurance.ExpandCoverage(db, doc, assurance.CoverageDirectives{Namespace: "demo", Environment: "dev", Snapshot: "snap-1"})
}

func advanceCoverageForTest(t *testing.T, db string, doc *cdl.EIRDoc) (int, error) {
	t.Helper()
	return assurance.AdvanceCoverage(db, doc)
}

func summaryForTest(t *testing.T, db string, doc *cdl.EIRDoc) (assurance.Summary, error) {
	t.Helper()
	return assurance.CoverageSummary(db, doc)
}

func TestAllCatalogsUseEvidenceAnchoredIdentity(t *testing.T) {
	for id, catalog := range Catalogs {
		unit := Plan(catalog, nil, Directives{Namespace: "demo", Snapshot: "dev"})
		if unit.IDFromBindings != "source-bindings" {
			t.Errorf("catalog %s is not evidence-anchored: IDFromBindings=%q", id, unit.IDFromBindings)
		}
		if len(unit.RuntimeFields) != 1 || unit.RuntimeFields[0] != "block-id" {
			t.Errorf("catalog %s does not inject block-id: %v", id, unit.RuntimeFields)
		}
		if unit.MinRecords != 1 || !unit.RequiresReview {
			t.Errorf("catalog %s lacks floor/review: MinRecords=%d RequiresReview=%v", id, unit.MinRecords, unit.RequiresReview)
		}
		if unit.IDDiscriminatorField != catalog.DiscriminatorField {
			t.Errorf("catalog %s discriminator = %q, want %q", id, unit.IDDiscriminatorField, catalog.DiscriminatorField)
		}
		if !unit.AllowZeroDomain {
			t.Errorf("catalog %s does not allow evidence-backed zero domains", id)
		}
	}
}
