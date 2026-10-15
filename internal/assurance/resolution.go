package assurance

import (
	"fmt"
	"os"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Decision is one recorded human disposition for an open ticket.
type Decision struct {
	TicketID       string
	Classification string // Ineligible (Eligible requires a committed SM first)
	Reason         string
}

// LoadDecisions parses a plain decisions file: one `TICKET-ID | Classification |
// Reason` line per ticket; blank lines and `#` comments are ignored.
func LoadDecisions(path string) ([]Decision, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var decisions []Decision
	for number, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("decisions line %d: expected `TICKET-ID | Classification | Reason`", number+1)
		}
		decision := Decision{
			TicketID:       strings.TrimSpace(parts[0]),
			Classification: strings.TrimSpace(parts[1]),
			Reason:         strings.TrimSpace(parts[2]),
		}
		if decision.TicketID == "" || decision.Reason == "" {
			return nil, fmt.Errorf("decisions line %d: ticket ID and reason are required", number+1)
		}
		switch decision.Classification {
		case "Ineligible":
		case "Eligible":
			return nil, fmt.Errorf("decisions line %d: %s is Eligible; author and commit a complete SM record before applying the decision", number+1, decision.TicketID)
		default:
			return nil, fmt.Errorf("decisions line %d: classification must be Ineligible or Eligible", number+1)
		}
		decisions = append(decisions, decision)
	}
	if len(decisions) == 0 {
		return nil, fmt.Errorf("decisions file %s declares no decisions", path)
	}
	return decisions, nil
}

// ResolutionPlan builds the mechanical units that apply reviewed decisions —
// two FSM-legal ticket transitions per ticket (T-OPEN -> T-PENDING ->
// T-COMPLETE) and one lifecycle reclassification per scope
// (Unknown -> Ineligible with the recorded reason, provenance, and classifying
// invocation) — together with each unit's exact payload. The plan fails closed
// on unknown tickets, terminal or non-open tickets, a scope that is not
// currently Unknown, or an unsupported classification.
func ResolutionPlan(db string, doc *cdl.EIRDoc, decisions []Decision) ([]work.Unit, map[string][]byte, error) {
	ticketNames, _, _, err := work.FieldContract(doc, "TICKET")
	if err != nil {
		return nil, nil, err
	}
	lifecycleNames, _, _, err := work.FieldContract(doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return nil, nil, err
	}
	var units []work.Unit
	payloads := map[string][]byte{}
	for _, decision := range decisions {
		ticketPayload, err := store.Get(db, doc, "TICKET", decision.TicketID)
		if err != nil {
			return nil, nil, fmt.Errorf("resolution: ticket %s: %w", decision.TicketID, err)
		}
		ticket, err := work.PayloadText(ticketNames, ticketPayload)
		if err != nil {
			return nil, nil, err
		}
		status := strings.TrimSpace(ticket["status"])
		if !work.TicketBlocksExitA(status) {
			return nil, nil, fmt.Errorf("resolution: ticket %s is already terminal (%s)", decision.TicketID, status)
		}
		if status != "T-OPEN" {
			return nil, nil, fmt.Errorf("resolution: ticket %s is %s; only T-OPEN tickets can be resolved by this path", decision.TicketID, status)
		}
		scope := firstTokenField(ticket["affected-ids"])
		if scope == "" {
			return nil, nil, fmt.Errorf("resolution: ticket %s declares no affected scope", decision.TicketID)
		}
		lifecyclePayload, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", scope)
		if err != nil {
			return nil, nil, fmt.Errorf("resolution: lifecycle row for %s: %w", scope, err)
		}
		lifecycle, err := work.PayloadText(lifecycleNames, lifecyclePayload)
		if err != nil {
			return nil, nil, err
		}
		if strings.TrimSpace(lifecycle["classification"]) != "Unknown" {
			return nil, nil, fmt.Errorf("resolution: %s is not currently Unknown; reclassification is not applicable", scope)
		}
		trace := strings.TrimSpace(ticket["trace-history"])
		pendingTrace := appendTrace(trace, "owner-review:ALFA -> T-PENDING")
		completeTrace := appendTrace(pendingTrace, "owner-review:ALFA -> T-COMPLETE")
		pending := replaceField(replaceField(string(ticketPayload), "trace-history", pendingTrace), "status", "T-PENDING")
		complete := replaceField(replaceField(pending, "trace-history", completeTrace), "status", "T-COMPLETE")
		reason := strings.Join(strings.Fields(decision.Reason), " ")
		reclassified := replaceField(replaceField(string(lifecyclePayload), "classification", "Ineligible"), "reason", reason)

		stepOne := "U-ASSURANCE-TICKET-RESOLVE-" + decision.TicketID + "-1"
		stepTwo := "U-ASSURANCE-TICKET-RESOLVE-" + decision.TicketID + "-2"
		lifecycleUnit := "U-ASSURANCE-LIFECYCLE-RESOLVE-" + scope
		payloads[stepOne] = []byte(pending)
		payloads[stepTwo] = []byte(complete)
		payloads[lifecycleUnit] = []byte(reclassified)
		units = append(units,
			work.Unit{ID: stepOne, Mode: "Discovery", RecordID: decision.TicketID, ResultType: "TICKET", Sections: []string{"ticket-fsm"}},
			work.Unit{ID: stepTwo, Mode: "Discovery", RecordID: decision.TicketID, ResultType: "TICKET", Sections: []string{"ticket-fsm"}},
			work.Unit{ID: lifecycleUnit, Mode: "Discovery", RecordID: scope, RuntimeInvocationFields: []string{"classified-by-invocation"}, ResultType: "LIFECYCLE-ELIGIBILITY", Sections: []string{"dependency-provenance"}},
		)
	}
	return units, payloads, nil
}

func appendTrace(trace, step string) string {
	if strings.TrimSpace(trace) == "" || strings.EqualFold(strings.TrimSpace(trace), "None") {
		return step
	}
	return trace + "; " + step
}

// replaceField rewrites one `- **field:** value` line, preserving order.
func replaceField(payload, field, value string) string {
	lines := strings.Split(payload, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **"+field+":**") {
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = indent + "- **" + field + ":** " + value
		}
	}
	return strings.Join(lines, "\n")
}

// MechanicalResolutionExecutor writes one precomputed payload per unit. The
// session supplies the invocation identity for invocation-owned fields.
type MechanicalResolutionExecutor struct {
	Payloads map[string][]byte
}

// Execute implements work.Executor.
func (e MechanicalResolutionExecutor) Execute(inv work.Invocation) (work.Proposal, error) {
	payload, ok := e.Payloads[inv.Unit.ID]
	if !ok {
		return work.Proposal{}, fmt.Errorf("resolution: no payload planned for unit %s", inv.Unit.ID)
	}
	return work.Proposal{FieldID: inv.Unit.ResultType, Records: []work.Record{{Payload: payload}}}, nil
}
