package assurance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func seedOpenTicket(t *testing.T, db string, doc interface{ ID() string }, scope, claim, ticketID string) {
	t.Helper()
	_ = doc
}

func TestLoadDecisions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decisions.md")
	content := "# ticket decisions\n\nTKT-NONE-POV-4-001 | Ineligible | Owner review: stateless helper.\n\nTKT-NONE-POV-4-002 | Eligible | should fail\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadDecisions(path); err == nil || !strings.Contains(err.Error(), "author and commit a complete SM") {
		t.Fatalf("Eligible decision accepted: %v", err)
	}
	content = "TKT-NONE-POV-4-001 | Ineligible | Owner review: stateless helper.\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	decisions, err := LoadDecisions(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(decisions) != 1 || decisions[0].TicketID != "TKT-NONE-POV-4-001" || decisions[0].Classification != "Ineligible" {
		t.Fatalf("decisions = %+v", decisions)
	}
	content = "TKT-NONE-POV-4-001 | Unknown | nope\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadDecisions(path); err == nil || !strings.Contains(err.Error(), "must be Ineligible or Eligible") {
		t.Fatalf("unknown classification accepted: %v", err)
	}
}

func TestResolutionPlanAppliesOwnerDecision(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	seedClaim(t, db, doc, "CLM-AAAAAAAAAAAA", "SRC-FILE-AAAAAAAAAAAA")
	seedUnknownLifecycle(t, db, doc, "SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA")
	ticket := ticketTestPayload("SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA")
	if err := store.Ingest(db, doc, "TICKET", "TKT-NONE-POV-4-001", 1, ticket); err != nil {
		t.Fatal(err)
	}
	decisions := []Decision{{
		TicketID:       "TKT-NONE-POV-4-001",
		Classification: "Ineligible",
		Reason:         "Owner review 2026-10-15: utility module without a managed entity lifecycle.",
	}}
	units, payloads, err := ResolutionPlan(db, doc, decisions)
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 3 {
		t.Fatalf("units = %+v", units)
	}
	if units[0].RecordID != "TKT-NONE-POV-4-001" || units[2].RecordID != "SRC-FILE-AAAAAAAAAAAA" {
		t.Fatalf("unit identities = %+v", units)
	}
	runConfig := work.DefaultConfig()
	runConfig.UpdateExisting = true
	session, err := work.NewSession(db, doc, runConfig, MechanicalResolutionExecutor{Payloads: payloads}, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run(units)
	if err != nil {
		t.Fatal(err)
	}
	ticketStored, err := store.Get(db, doc, "TICKET", "TKT-NONE-POV-4-001")
	if err != nil {
		t.Fatal(err)
	}
	text := string(ticketStored)
	for _, marker := range []string{"- **status:** T-COMPLETE", "owner-review:ALFA -> T-PENDING", "owner-review:ALFA -> T-COMPLETE"} {
		if !strings.Contains(text, marker) {
			t.Fatalf("ticket payload missing %q:\n%s", marker, text)
		}
	}
	lifecycleStored, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", "SRC-FILE-AAAAAAAAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	lifecycleText := string(lifecycleStored)
	if !strings.Contains(lifecycleText, "- **classification:** Ineligible") || !strings.Contains(lifecycleText, "utility module without a managed entity lifecycle") {
		t.Fatalf("lifecycle payload:\n%s", lifecycleText)
	}
	if !strings.Contains(lifecycleText, "- **classified-by-invocation:** "+report.Units[2].InvocationID) {
		t.Fatalf("classifying invocation not injected:\n%s", lifecycleText)
	}
	if diagnostics, err := ValidateLifecycleEligibility(db, doc); err != nil || len(diagnostics) != 0 {
		t.Fatalf("lifecycle diagnostics = %v (err %v)", diagnostics, err)
	}
	scopes, err := TicketlessUnknowns(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 0 {
		t.Fatalf("ticketless unknowns remain: %+v", scopes)
	}

	// Re-applying the same decision fails closed on the terminal ticket.
	if _, _, err := ResolutionPlan(db, doc, decisions); err == nil || !strings.Contains(err.Error(), "already terminal") {
		t.Fatalf("terminal ticket accepted: %v", err)
	}
}
