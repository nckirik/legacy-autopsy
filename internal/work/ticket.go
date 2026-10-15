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
// executor assign a terminal or probe state merely by writing its token, and
// persona/POV-bound units additionally pass the §9.2 actor-right checks.
func validateTicketLifecycle(doc *cdl.EIRDoc, unit Unit, caps *Caps, recordID string, values map[string]string) error {
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
	declared := false
	for _, allowed := range state.Allows {
		if allowed[0] == from && allowed[1] == status {
			declared = true
			break
		}
	}
	if !declared {
		return fmt.Errorf("TICKET transition %s -> %s is not declared by TICKET-STATE", from, status)
	}
	return validateTicketActor(unit, values, trace, status)
}

// validateTicketActor enforces the §9.2 actor rights: only the target POV
// updates a ticket to T-PENDING or T-PROBE-REQUIRED, and only the origin POV
// accepts or rejects to T-COMPLETE or T-FOLLOW-UP. Transitions recorded as a
// human-authorized review (the `owner-review` actor) are exempt; a unit without
// a bound POV cannot claim target/origin rights and fails closed.
func validateTicketActor(unit Unit, values map[string]string, trace, to string) error {
	actor := transitionActor(trace, to)
	switch actor {
	case "", "owner-review", "human-hatch", "human":
		return nil
	}
	require := ""
	switch to {
	case "T-PENDING", "T-PROBE-REQUIRED":
		require = ticketTargetPOV(values["target-pov-ledger"])
	case "T-COMPLETE", "T-FOLLOW-UP":
		require = ticketSourcePOV(values["source-pov-file-invocation"])
	}
	if require == "" {
		return nil
	}
	pov := strings.ToUpper(strings.TrimSpace(unit.POV))
	if pov == "" {
		return fmt.Errorf("TICKET transition to %s requires active POV %s (%s); unbound units cannot claim actor rights", to, require, actor)
	}
	if pov != require {
		return fmt.Errorf("TICKET transition to %s requires active POV %s for the %s actor, got %s", to, require, actor, pov)
	}
	return nil
}

// transitionActor extracts the actor token of the trace step that records the
// given status.
func transitionActor(trace, status string) string {
	steps := strings.Split(trace, ";")
	for i := len(steps) - 1; i >= 0; i-- {
		step := strings.TrimSpace(steps[i])
		if !strings.Contains(step, status) {
			continue
		}
		if idx := strings.Index(step, ":"); idx > 0 {
			return strings.ToLower(strings.TrimSpace(step[:idx]))
		}
		return ""
	}
	return ""
}

// ticketSourcePOV extracts the active POV of the ticket's origin from
// `source-pov-file-invocation` (the persona POV file segment).
func ticketSourcePOV(raw string) string {
	for _, segment := range strings.Split(raw, "|") {
		segment = strings.TrimSpace(segment)
		if !strings.HasSuffix(segment, "-POV.md") {
			continue
		}
		name := strings.TrimSuffix(segment, ".md")
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		if pov := povFromFileName(name); pov != "" {
			return pov
		}
	}
	return ""
}

// ticketTargetPOV extracts the target POV from `target-pov-ledger`.
func ticketTargetPOV(raw string) string {
	token := strings.ToUpper(strings.TrimSpace(strings.SplitN(raw, "|", 2)[0]))
	fields := strings.Fields(token)
	if len(fields) == 0 {
		return ""
	}
	if pov := povFromToken(fields[0]); pov != "" {
		return pov
	}
	return ""
}

// povFromFileName maps a numbered POV file stem (for example
// "01-FRONTEND-POV") to its POV token.
func povFromFileName(name string) string {
	fields := strings.SplitN(strings.ToUpper(name), "-", 2)
	if len(fields) == 0 {
		return ""
	}
	return povFromToken("POV-" + strings.TrimLeft(fields[0], "0"))
}

// povFromToken validates a POV-n token.
func povFromToken(token string) string {
	token = strings.ToUpper(strings.TrimSpace(token))
	if !strings.HasPrefix(token, "POV-") {
		return ""
	}
	number := strings.TrimPrefix(token, "POV-")
	if number < "1" || number > "5" || len(number) != 1 {
		return ""
	}
	return "POV-" + number
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
