package work

import (
	"strings"
	"testing"
)

func ticketPayload(status, trace, reason string) []byte {
	return []byte("- **source-pov-file-invocation:** SRC-1 | 01-FRONTEND-POV.md | ALFA-AD-POV-1-001\n" +
		"- **target-pov-ledger:** POV-2 | personas/_shared/02-BACKEND-POV.md\n" +
		"- **status:** " + status + "\n" +
		"- **trace-history:** " + trace + "\n" +
		"- **affected-ids:** SRC-FILE-AAAAAAAAAAAA\n" +
		"- **memory-anchor:** bounded static read of the entry surface\n" +
		"- **description:** one technical ambiguity about validation order\n" +
		"- **static-investigation:** bounded static paths attempted; result inconclusive\n" +
		"- **escalation-reason:** " + reason + "\n" +
		"- **probe-feasibility:** not applicable to static evidence\n" +
		"- **coverage-promotion-effect:** blocks the affected matrix cell only\n")
}

func TestTicketLifecycleValidation(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-TICKET-TEST", Mode: "Ticket Resolution", RecordID: "TKT-AD-POV-1-001", ResultType: "TICKET", POV: "POV-2"}
	commit := func(payload []byte, update bool) error {
		_, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "TICKET", Records: []Record{{Payload: payload}}}, update)
		return err
	}
	if err := commit(ticketPayload("T-OPEN", "origin:ALFA -> T-OPEN", "None"), false); err != nil {
		t.Fatalf("initial ticket rejected: %v", err)
	}
	if err := commit(ticketPayload("T-PENDING", "origin:ALFA -> T-OPEN; target:ALFA -> T-PENDING", "None"), true); err != nil {
		t.Fatalf("declared transition rejected: %v", err)
	}
	if err := commit(ticketPayload("T-OPEN", "origin:ALFA -> T-OPEN without target investigation", "None"), true); err == nil || !strings.Contains(err.Error(), "not declared") {
		t.Fatalf("undeclared transition accepted: %v", err)
	}
}

// TestTicketActorRights enforces §9.2: only the target POV moves a ticket to
// pending/probe, only the origin accepts or rejects, unbound units cannot claim
// either right, and a recorded human review is exempt.
func TestTicketActorRights(t *testing.T) {
	db, doc := testDoc(t)
	commitAs := func(pov string, payload []byte, update bool) error {
		unit := Unit{ID: "U-TICKET-ACTOR-" + pov, Mode: "Ticket Resolution", RecordID: "TKT-AD-POV-1-003", ResultType: "TICKET", POV: pov}
		_, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "TICKET", Records: []Record{{Payload: payload}}}, update)
		return err
	}
	if err := commitAs("POV-2", ticketPayload("T-OPEN", "origin:ALFA -> T-OPEN", "None"), false); err != nil {
		t.Fatal(err)
	}
	// Origin POV cannot move the ticket to pending.
	if err := commitAs("POV-1", ticketPayload("T-PENDING", "target:ALFA -> T-PENDING", "None"), true); err == nil || !strings.Contains(err.Error(), "requires active POV POV-2") {
		t.Fatalf("origin POV accepted as target: %v", err)
	}
	// The target POV can.
	if err := commitAs("POV-2", ticketPayload("T-PENDING", "target:ALFA -> T-PENDING", "None"), true); err != nil {
		t.Fatalf("target POV rejected: %v", err)
	}
	// Target POV cannot accept its own answer; the origin must.
	accepted := ticketPayload("T-COMPLETE", "target:ALFA -> T-PENDING; origin:ALFA -> T-COMPLETE", "None")
	if err := commitAs("POV-2", accepted, true); err == nil || !strings.Contains(err.Error(), "requires active POV POV-1") {
		t.Fatalf("target POV accepted as origin: %v", err)
	}
	if err := commitAs("POV-1", accepted, true); err != nil {
		t.Fatalf("origin POV rejected: %v", err)
	}
}

func TestTicketUnboundActorAndHumanReview(t *testing.T) {
	db, doc := testDoc(t)
	commitAs := func(pov string, payload []byte, update bool) error {
		unit := Unit{ID: "U-TICKET-UNBOUND", Mode: "Ticket Resolution", RecordID: "TKT-AD-POV-1-004", ResultType: "TICKET", POV: pov}
		_, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "TICKET", Records: []Record{{Payload: payload}}}, update)
		return err
	}
	if err := commitAs("", ticketPayload("T-OPEN", "origin:ALFA -> T-OPEN", "None"), false); err != nil {
		t.Fatal(err)
	}
	if err := commitAs("", ticketPayload("T-PENDING", "target:ALFA -> T-PENDING", "None"), true); err == nil || !strings.Contains(err.Error(), "cannot claim actor rights") {
		t.Fatalf("unbound unit claimed target rights: %v", err)
	}
	if err := commitAs("", ticketPayload("T-PENDING", "owner-review:ALFA -> T-PENDING", "None"), true); err != nil {
		t.Fatalf("human review actor rejected: %v", err)
	}
}

func TestTicketCreationRules(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{ID: "U-TICKET-CREATE", Mode: "Ticket Resolution", RecordID: "TKT-AD-POV-1-002", ResultType: "TICKET"}
	commit := func(payload []byte) error {
		_, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "TICKET", Records: []Record{{Payload: payload}}}, false)
		return err
	}
	if err := commit(ticketPayload("T-PENDING", "origin -> T-PENDING", "None")); err == nil || !strings.Contains(err.Error(), "must start in T-OPEN") {
		t.Fatalf("non-initial creation accepted: %v", err)
	}
	if err := commit(ticketPayload("T-HUMAN-REQUIRED", "origin -> T-HUMAN-REQUIRED", "NO-AUTHORITY-TO-DECIDE")); err == nil || !strings.Contains(err.Error(), "must start in T-OPEN") {
		t.Fatalf("direct T-HUMAN-REQUIRED creation accepted: %v", err)
	}
	if err := commit(ticketPayload("T-OPEN", "origin:ALFA -> created without the state token", "None")); err == nil || !strings.Contains(err.Error(), "trace-history") {
		t.Fatalf("missing transition history accepted: %v", err)
	}
	if err := commit(ticketPayload("T-PROBE-REQUIRED", "origin -> T-PROBE-REQUIRED", "None")); err == nil || !strings.Contains(err.Error(), "RUNTIME-REQUIRED") {
		t.Fatalf("probe state without RUNTIME-REQUIRED accepted: %v", err)
	}
}

func TestTicketBlocksExitA(t *testing.T) {
	for state, want := range map[string]bool{
		"T-OPEN": true, "T-PENDING": true, "T-PROBE-REQUIRED": true,
		"T-FOLLOW-UP": true, "T-HUMAN-REQUIRED": true,
		"T-COMPLETE": false, "T-INVALID": false, "T-SUPERSEDED": false, "T-OUT-OF-SCOPE": false,
	} {
		if got := TicketBlocksExitA(state); got != want {
			t.Fatalf("TicketBlocksExitA(%s) = %v, want %v", state, got, want)
		}
	}
}
