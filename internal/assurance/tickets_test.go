package assurance

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func seedClaim(t *testing.T, db string, doc *cdl.EIRDoc, id, subject string) {
	t.Helper()
	payload := []byte("- **subject-id-record-version:** " + subject + " version 1\n" +
		"- **field-step:** lifecycle\n" +
		"- **claim:** The unit performs static validation only.\n" +
		"- **materiality:** Supporting\n" +
		"- **evidence-level:** [E-DIRECT]\n" +
		"- **anchor:** app/a.py:1\n" +
		"- **source-fingerprint-snapshot:** sha256:aaa\n" +
		"- **evidence-route-lineage-id:** static-filesystem\n" +
		"- **projection-classification:** raw-visible\n" +
		"- **independence-basis:** single static route\n" +
		"- **created-last-revalidated-by:** ALFA-NONE-NONE-001\n")
	if err := store.Ingest(db, doc, "CLAIM", id, 1, payload); err != nil {
		t.Fatal(err)
	}
}

func seedUnknownLifecycle(t *testing.T, db string, doc *cdl.EIRDoc, scope, claim string) {
	t.Helper()
	row := []byte("- **scope-id:** " + scope + "\n- **classification:** Unknown\n- **reason:** no committed claim establishes a meaningful lifecycle\n- **clm-evidence-bindings:** " + claim + "\n- **classified-by-invocation:** ALFA-NONE-NONE-099\n- **snapshot:** dev\n")
	if err := store.Ingest(db, doc, "LIFECYCLE-ELIGIBILITY", scope, 1, row); err != nil {
		t.Fatal(err)
	}
}

func ticketTestPayload(scope, claim string) []byte {
	return []byte("- **source-pov-file-invocation:** " + scope + " | 04-DATA-POV.md | ALFA-NONE-NONE-099\n" +
		"- **target-pov-ledger:** POV-4 | personas/_shared/04-DATA-POV.md\n" +
		"- **status:** T-OPEN\n" +
		"- **trace-history:** origin:ALFA -> T-OPEN\n" +
		"- **affected-ids:** " + scope + ", " + claim + "\n" +
		"- **memory-anchor:** lifecycle row " + scope + " classified Unknown at ALFA-NONE-NONE-099\n" +
		"- **description:** Lifecycle eligibility for " + scope + " remains Unknown after bounded static analysis.\n" +
		"- **static-investigation:** Committed claims reviewed; no fact establishes a meaningful lifecycle.\n" +
		"- **escalation-reason:** DOMAIN-MEANING-REQUIRED\n" +
		"- **probe-feasibility:** A safe probe adds no lifecycle evidence.\n" +
		"- **coverage-promotion-effect:** Affects only " + scope + " and its dependency closure.\n")
}

func TestTicketPlanOriginatesTicketsForUnknownScopes(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	seedClaim(t, db, doc, "CLM-AAAAAAAAAAAA", "SRC-FILE-AAAAAAAAAAAA")
	seedUnknownLifecycle(t, db, doc, "SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA")

	scopes, err := TicketlessUnknowns(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 || scopes[0].ScopeID != "SRC-FILE-AAAAAAAAAAAA" || len(scopes[0].Claims) != 1 {
		t.Fatalf("scopes = %+v", scopes)
	}
	units, err := TicketPlan(db, doc, "demo", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 1 || units[0].RecordID != "TKT-NONE-POV-4-001" || units[0].POV != "POV-4" {
		t.Fatalf("units = %+v", units)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		if !strings.Contains(string(inv.Context.Bytes), "SRC-FILE-AAAAAAAAAAAA") || !strings.Contains(string(inv.Context.Bytes), "CLM-AAAAAAAAAAAA") {
			t.Error("ticket context lacks the unresolved scope or its claims")
		}
		return work.Proposal{FieldID: "TICKET", Records: []work.Record{{Payload: ticketTestPayload("SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA")}}}, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run(units); err != nil {
		t.Fatal(err)
	}
	stored, err := store.Get(db, doc, "TICKET", "TKT-NONE-POV-4-001")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stored), "- **status:** T-OPEN") {
		t.Fatalf("ticket payload:\n%s", stored)
	}
	scopes, err = TicketlessUnknowns(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 0 {
		t.Fatalf("ticketless unknowns remain: %+v", scopes)
	}
	if diagnostics, err := ValidateLifecycleEligibility(db, doc); err != nil || len(diagnostics) != 0 {
		t.Fatalf("lifecycle diagnostics after ticketing = %v (err %v)", diagnostics, err)
	}

	// A second unknown scope allocates the next free sequence.
	seedSource(t, db, doc, "SRC-FILE-BBBBBBBBBBBB")
	seedClaim(t, db, doc, "CLM-BBBBBBBBBBBB", "SRC-FILE-BBBBBBBBBBBB")
	seedUnknownLifecycle(t, db, doc, "SRC-FILE-BBBBBBBBBBBB", "CLM-BBBBBBBBBBBB")
	units, err = TicketPlan(db, doc, "demo", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 1 || units[0].RecordID != "TKT-NONE-POV-4-002" {
		t.Fatalf("second allocation = %+v", units)
	}
}
