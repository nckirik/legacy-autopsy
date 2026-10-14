package work

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// TicketTerminalStates are the declared terminal TICKET states; a ticket in one
// of them no longer blocks Exit A.
var TicketTerminalStates = map[string]bool{
	"T-COMPLETE":     true,
	"T-INVALID":      true,
	"T-SUPERSEDED":   true,
	"T-OUT-OF-SCOPE": true,
}

// TicketBlocksExitA reports whether a committed ticket state still blocks
// closure. T-HUMAN-REQUIRED remains a coverage gap or partial unless the unit
// is separately approved as excluded.
func TicketBlocksExitA(state string) bool {
	state = strings.TrimSpace(state)
	if state == "" {
		return true
	}
	return !TicketTerminalStates[state]
}

// validateTicketLifecycle enforces the declared TICKET-STATE machine at commit:
// a new ticket starts in the declared initial state, updates follow declared
// transitions, the transition token is recorded in the trace history, and
// probe/human states carry a matching escalation reason. It never lets an
// executor assign a terminal or probe state merely by writing its token.
func validateTicketLifecycle(doc *cdl.EIRDoc, caps *Caps, recordID string, values map[string]string) error {
	state, ok := declaredState(doc, "TICKET-STATE")
	if !ok {
		return fmt.Errorf("work: TICKET-STATE is not declared in the EIR")
	}
	status := strings.TrimSpace(values["status"])
	if !stringIn(state.Values, status) {
		return fmt.Errorf("TICKET.status %q is not a declared TICKET-STATE value", status)
	}
	trace := strings.TrimSpace(values["trace-history"])
	if trace == "" || !strings.Contains(trace, status) {
		return fmt.Errorf("TICKET.trace-history must record the %s transition, not only the resulting status token", status)
	}
	reason := strings.TrimSpace(values["escalation-reason"])
	if status == "T-PROBE-REQUIRED" && reason != "RUNTIME-REQUIRED" {
		return fmt.Errorf("TICKET %s requires escalation-reason RUNTIME-REQUIRED, got %q", status, reason)
	}
	if status == "T-HUMAN-REQUIRED" && (reason == "" || strings.EqualFold(reason, "None")) {
		return fmt.Errorf("TICKET %s requires an exact escalation reason", status)
	}
	existing, err := store.Get(caps.db, caps.doc, "TICKET", recordID)
	if err != nil {
		if status != state.Initial {
			return fmt.Errorf("a new ticket must start in %s, not %s; T-HUMAN-REQUIRED is never created directly", state.Initial, status)
		}
		return nil
	}
	names, _, _, err := FieldContract(doc, "TICKET")
	if err != nil {
		return err
	}
	previous, err := PayloadText(names, existing)
	if err != nil {
		return err
	}
	from := strings.TrimSpace(previous["status"])
	if from == status {
		return nil
	}
	for _, allowed := range state.Allows {
		if allowed[0] == from && allowed[1] == status {
			return nil
		}
	}
	return fmt.Errorf("TICKET transition %s -> %s is not declared by TICKET-STATE", from, status)
}

func declaredState(doc *cdl.EIRDoc, id string) (cdl.EIRState, bool) {
	for _, state := range doc.Declarations.States {
		if state.ID == id {
			return state, true
		}
	}
	return cdl.EIRState{}, false
}

func stringIn(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
