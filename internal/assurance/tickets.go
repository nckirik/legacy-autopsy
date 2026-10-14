package assurance

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// TicketMoniker is the harness ticket moniker for tickets originated by
// runtime-driven discovery of unresolved semantic predicates.
const TicketMoniker = "TKT"

// UnknownScope is one lifecycle scope classified Unknown without a live ticket.
type UnknownScope struct {
	ScopeID      string
	Reason       string
	Claims       []string
	Coordinate   string
	ClassifiedBy string
}

// TicketlessUnknowns enumerates lifecycle-eligibility rows classified Unknown
// that no live ticket references, in scope order. These are the unresolved
// semantic predicates that honest coverage requires to be visible as tickets
// (protocol/ticket-fsm.cdl §9.1 and protocol/traceability.cdl §12.4).
func TicketlessUnknowns(db string, doc *cdl.EIRDoc) ([]UnknownScope, error) {
	names, _, _, err := work.FieldContract(doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return nil, err
	}
	inventory, err := LoadInventory(db, doc)
	if err != nil {
		return nil, err
	}
	coordinates := map[string]string{}
	for _, item := range inventory {
		coordinates[item.ID] = item.Coordinate
	}
	sort.Strings(ids)
	var scopes []UnknownScope
	for _, id := range ids {
		payload, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(values["classification"]) != "Unknown" {
			continue
		}
		scope := firstTokenField(values["scope-id"])
		if scope == "" || ticketReferences(db, doc, scope) {
			continue
		}
		claimTokens := []string{}
		for _, token := range splitTokens(values["clm-evidence-bindings"]) {
			if strings.HasPrefix(token, "CLM-") {
				claimTokens = append(claimTokens, token)
			}
		}
		scopes = append(scopes, UnknownScope{
			ScopeID:      scope,
			Reason:       strings.Join(strings.Fields(values["reason"]), " "),
			Claims:       claimTokens,
			Coordinate:   coordinates[scope],
			ClassifiedBy: strings.Join(strings.Fields(values["classified-by-invocation"]), " "),
		})
	}
	return scopes, nil
}

// nextTicketID allocates the next free ticket sequence for the moniker,
// persona, and target POV. Allocation is deterministic over the committed log
// and never reuses a sequence.
func nextTicketID(db string, doc *cdl.EIRDoc, persona, pov string) (int, string, error) {
	if persona == "" {
		persona = "NONE"
	}
	if pov == "" {
		pov = "NONE"
	}
	prefix := TicketMoniker + "-" + persona + "-" + pov + "-"
	ids, err := store.List(db, doc, "TICKET")
	if err != nil {
		return 0, "", err
	}
	highest := 0
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			continue
		}
		sequence, err := strconv.Atoi(strings.TrimPrefix(id, prefix))
		if err == nil && sequence > highest {
			highest = sequence
		}
	}
	return highest + 1, prefix, nil
}

// TicketPlan builds one bounded Ticket Resolution unit per ticketless Unknown
// lifecycle scope. The model authors the technical ambiguity, the static
// investigation summary, the probe feasibility, and the escalation reason from
// the closed enum; the runtime allocates the T-OPEN ticket identity. Lifecycle
// ambiguity binds to POV-4 Data Movement, whose scope explicitly includes the
// persistence lifecycle.
func TicketPlan(db string, doc *cdl.EIRDoc, namespace, snapshot string) ([]work.Unit, error) {
	scopes, err := TicketlessUnknowns(db, doc)
	if err != nil {
		return nil, err
	}
	claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	claimText := map[string]string{}
	for _, scope := range scopes {
		for _, claim := range scope.Claims {
			payload, err := store.Get(db, doc, "CLAIM", claim)
			if err != nil {
				return nil, err
			}
			values, err := work.PayloadText(claimNames, payload)
			if err != nil {
				return nil, err
			}
			claimText[claim] = strings.Join(strings.Fields(values["claim"]), " ")
		}
	}
	var units []work.Unit
	sequence, prefix, err := nextTicketID(db, doc, "NONE", "POV-4")
	if err != nil {
		return nil, err
	}
	for _, scope := range scopes {
		id := fmt.Sprintf("%s%03d", prefix, sequence)
		sequence++
		var b strings.Builder
		b.WriteString("# Ticket authoring context (runtime-generated, bounded)\n\n")
		b.WriteString("Author exactly one TICKET record for the unresolved lifecycle classification below.\n\n")
		b.WriteString("## Runtime directives (must be used exactly)\n\n")
		fmt.Fprintf(&b, "- **System Namespace:** %s\n", namespace)
		fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", snapshot)
		fmt.Fprintf(&b, "- `status` is `T-OPEN` (the ticket is newly originated; never open in another state).\n")
		fmt.Fprintf(&b, "- `trace-history` must record this origin transition, for example `origin:ALFA -> T-OPEN`.\n")
		fmt.Fprintf(&b, "- `source-pov-file-invocation` is `%s | 04-DATA-POV.md | %s` with the exact values shown.\n", scope.ScopeID, scope.ClassifiedBy)
		fmt.Fprintf(&b, "- `target-pov-ledger` is `POV-4 | personas/_shared/04-DATA-POV.md`.\n")
		fmt.Fprintf(&b, "- `affected-ids` is `%s` and the claim IDs shown below, comma-separated.\n", scope.ScopeID)
		fmt.Fprintf(&b, "- `memory-anchor` restates the exact committed inputs (lifecycle row, source unit, claims) so the ticket is stateless.\n")
		b.WriteString("- `description` is exactly one technical ambiguity: why the lifecycle classification remains Unknown and what fact would resolve it.\n")
		b.WriteString("- `static-investigation` summarizes the bounded static paths already attempted and their result (do not invent new reads).\n")
		b.WriteString("- `escalation-reason` is one of: RUNTIME-REQUIRED, THIRD-PARTY-OPAQUE, UNAUTHORIZED-ACCESS, UNSAFE, DESTRUCTIVE, NO-SANDBOX, NO-AUTHORITY-TO-DECIDE, DOMAIN-MEANING-REQUIRED, EXTERNAL-FACT-UNAVAILABLE. Use `None` only if static analysis can already resolve it; otherwise choose the exact blocker.\n")
		b.WriteString("- `probe-feasibility` states whether a safe probe could resolve the ambiguity, or why not.\n")
		b.WriteString("- `coverage-promotion-effect` states the exact coverage and promotion effect, limited to the affected unit and its dependency closure.\n\n")
		b.WriteString("## Unresolved lifecycle classification\n\n")
		fmt.Fprintf(&b, "- **Scope:** %s\n", scope.ScopeID)
		fmt.Fprintf(&b, "- **Concrete coordinate:** %s\n", orNoneText(scope.Coordinate))
		fmt.Fprintf(&b, "- **Classification:** Unknown\n")
		fmt.Fprintf(&b, "- **Recorded reason:** %s\n", scope.Reason)
		fmt.Fprintf(&b, "- **Classified by:** %s\n", scope.ClassifiedBy)
		if len(scope.Claims) == 0 {
			b.WriteString("- **Claims:** None\n")
		} else {
			b.WriteString("- **Claims:**\n")
			for _, claim := range scope.Claims {
				fmt.Fprintf(&b, "  - %s: %s\n", claim, claimText[claim])
			}
		}
		units = append(units, work.Unit{
			ID:              "U-ASSURANCE-TICKET-" + scope.ScopeID,
			Mode:            "Discovery",
			Scope:           namespace,
			Persona:         "",
			POV:             "POV-4",
			RecordID:        id,
			ExpectedRecords: 1,
			ResultType:      "TICKET",
			Reads:           []work.Ref{{FieldID: "SOURCE-INVENTORY-RECORD", RecordID: scope.ScopeID}},
			Inputs:          []work.Input{{Name: "ticket-context.md", Bytes: []byte(b.String())}},
			Sections:        []string{"ticket-schema", "ticket-fsm", "dependency-provenance"},
			ReasoningEffort: "off",
		})
	}
	return units, nil
}

func orNoneText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "None"
	}
	return value
}
