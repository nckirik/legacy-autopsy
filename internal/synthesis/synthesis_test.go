package synthesis

import (
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
		{ID: "CLM-000000000001", Subject: "SRC-FILE-A", Text: "app/routes.js registers the public and API routes"},
		{ID: "CLM-000000000002", Subject: "SRC-FILE-B", Text: "the passport strategy authenticates users via SAML"},
	}
	unit := Plan(Catalogs["MOD"], claims, Directives{Namespace: "demo", Snapshot: "dev"})
	if unit.RuntimeFields[0] != "block-id" || unit.IDOwnerField != "domain-path" {
		t.Fatalf("unit = %+v", unit)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		context := string(inv.Context.Bytes)
		if !strings.Contains(context, "app/routes.js registers") {
			t.Error("context lacks claims")
		}
		payload := "- **block-id:** PENDING\n" +
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
	want, err := capabilities.TypedID("MOD", "demo", "app/routes", "Request Routing and Authentication")
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
		payload := "- **record-title:** X\n- **domain-path:** app/x\n- **responsibility-domain-boundary:** x\n- **public-interfaces:** None\n- **encapsulated-state-entities:** None\n- **dependencies-dependents:** None\n- **data-and-events-passed:** None\n- **failure-propagation:** None\n- **invariants-rules-owned:** None\n- **observed-legacy-placement:** None\n- **target-decision:** Pending\n- **personas-use-cases-capabilities:** None\n"
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
