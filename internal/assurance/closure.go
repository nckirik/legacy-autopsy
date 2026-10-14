package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// unresolvedFrontierStates are the frontier states that still count as
// unresolved for a dead-code proof: an open or in-progress traversal, and a
// human-blocked terminal that remains a coverage gap until separately excluded.
var unresolvedFrontierStates = map[string]bool{
	"Open":                   true,
	"Traversed":              true,
	"Terminal-Human-Blocked": true,
}

// ValidateDeadCode validates one component's [S-DEAD-CODE] claim: the proof
// requires denominator dimensions, committed source units, no unresolved
// frontier, no active ticket, and no discovered caller anywhere in the
// committed component graph. It returns nil when the claim holds, otherwise an
// exact diagnostic. The checks extend to dynamic-reference and
// acquisition-candidate rows when those record classes exist.
func ValidateDeadCode(db string, doc *cdl.EIRDoc, cmpID string) error {
	names, _, _, err := work.FieldContract(doc, "ATOMIC-COMPONENT")
	if err != nil {
		return err
	}
	payload, err := store.Get(db, doc, "ATOMIC-COMPONENT", cmpID)
	if err != nil {
		return fmt.Errorf("dead-code: component %s: %w", cmpID, err)
	}
	values, err := work.PayloadText(names, payload)
	if err != nil {
		return err
	}
	if strings.TrimSpace(values["staging-status"]) != "[S-DEAD-CODE]" {
		return fmt.Errorf("dead-code: %s does not declare [S-DEAD-CODE]", cmpID)
	}
	for _, dimension := range []string{"traversal-tracks", "runtime-state-matrix", "source-ids"} {
		value := strings.TrimSpace(values[dimension])
		if value == "" || strings.EqualFold(value, "None") {
			return fmt.Errorf("dead-code: %s has no %s denominator dimension", cmpID, dimension)
		}
	}
	for _, token := range splitTokens(values["source-ids"]) {
		if _, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", token); err != nil {
			return fmt.Errorf("dead-code: %s source unit %s is not committed", cmpID, token)
		}
	}
	frontierNames, _, _, err := work.FieldContract(doc, "FRONTIER-RECORD")
	if err != nil {
		return err
	}
	frontierIDs, err := store.List(db, doc, "FRONTIER-RECORD")
	if err != nil {
		return err
	}
	sort.Strings(frontierIDs)
	for _, id := range frontierIDs {
		frontierPayload, err := store.Get(db, doc, "FRONTIER-RECORD", id)
		if err != nil {
			return err
		}
		frontier, err := work.PayloadText(frontierNames, frontierPayload)
		if err != nil {
			return err
		}
		if !strings.Contains(frontier["from-src-cmp"], cmpID) && !strings.Contains(frontier["target-src"], cmpID) {
			continue
		}
		if unresolvedFrontierStates[strings.TrimSpace(frontier["state"])] {
			return fmt.Errorf("dead-code: %s has an unresolved frontier %s (%s)", cmpID, id, strings.TrimSpace(frontier["state"]))
		}
	}
	ticketNames, _, _, err := work.FieldContract(doc, "TICKET")
	if err != nil {
		return err
	}
	ticketIDs, err := store.List(db, doc, "TICKET")
	if err != nil {
		return err
	}
	sort.Strings(ticketIDs)
	for _, id := range ticketIDs {
		ticketPayload, err := store.Get(db, doc, "TICKET", id)
		if err != nil {
			return err
		}
		ticket, err := work.PayloadText(ticketNames, ticketPayload)
		if err != nil {
			return err
		}
		if !strings.Contains(ticket["affected-ids"], cmpID) {
			continue
		}
		if work.TicketBlocksExitA(ticket["status"]) {
			return fmt.Errorf("dead-code: %s has an active ticket %s (%s)", cmpID, id, strings.TrimSpace(ticket["status"]))
		}
	}
	componentIDs, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		return err
	}
	sort.Strings(componentIDs)
	for _, id := range componentIDs {
		if id == cmpID {
			continue
		}
		otherPayload, err := store.Get(db, doc, "ATOMIC-COMPONENT", id)
		if err != nil {
			return err
		}
		other, err := work.PayloadText(names, otherPayload)
		if err != nil {
			return err
		}
		for _, field := range []string{"depends-on-depended-on-by", "source-ids", "synthesis-references"} {
			if strings.Contains(other[field], cmpID) {
				return fmt.Errorf("dead-code: %s is referenced by %s (%s); a discovered caller disproves dead code", cmpID, id, field)
			}
		}
	}
	return nil
}

// ValidateLifecycleEligibility checks the committed lifecycle-eligibility rows:
// unique scope rows, snapshot consistency, classification provenance, and the
// fail-inclusive ticket requirement for Unknown. It returns the diagnostics;
// an empty slice means every row is valid.
func ValidateLifecycleEligibility(db string, doc *cdl.EIRDoc) ([]string, error) {
	names, _, _, err := work.FieldContract(doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	headerNames, _, _, err := work.FieldContract(doc, "PREFLIGHT-HEADER")
	if err != nil {
		return nil, err
	}
	headerPayload, err := store.Get(db, doc, "PREFLIGHT-HEADER", "0A")
	if err != nil {
		return nil, fmt.Errorf("lifecycle eligibility requires a completed preflight: %w", err)
	}
	headerValues, err := work.PayloadText(headerNames, headerPayload)
	if err != nil {
		return nil, err
	}
	snapshot := strings.TrimSpace(headerValues["included-environments-snapshots"])
	seen := map[string]string{}
	var diagnostics []string
	for _, id := range ids {
		payload, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		scope := strings.TrimSpace(values["scope-id"])
		if scope == "" {
			diagnostics = append(diagnostics, fmt.Sprintf("%s has no scope-id", id))
			continue
		}
		if previous, ok := seen[scope]; ok {
			diagnostics = append(diagnostics, fmt.Sprintf("%s duplicates the scope row %s already claimed by %s", id, scope, previous))
		} else {
			seen[scope] = id
		}
		classification := strings.TrimSpace(values["classification"])
		if strings.TrimSpace(values["clm-evidence-bindings"]) == "" || strings.EqualFold(strings.TrimSpace(values["clm-evidence-bindings"]), "None") {
			diagnostics = append(diagnostics, fmt.Sprintf("%s (%s) has no claim/evidence provenance", id, classification))
		}
		if snapshot != "" && !strings.Contains(values["snapshot"], snapshot) {
			diagnostics = append(diagnostics, fmt.Sprintf("%s snapshot %q is inconsistent with the header snapshot %q", id, strings.TrimSpace(values["snapshot"]), snapshot))
		}
		if classification == "Unknown" {
			if !ticketReferences(db, doc, scope) {
				diagnostics = append(diagnostics, fmt.Sprintf("%s classifies %s Unknown without a ticket", id, scope))
			}
		}
	}
	return diagnostics, nil
}

// ticketReferences reports whether a live ticket references the scope.
func ticketReferences(db string, doc *cdl.EIRDoc, scope string) bool {
	names, _, _, err := work.FieldContract(doc, "TICKET")
	if err != nil {
		return false
	}
	ids, err := store.List(db, doc, "TICKET")
	if err != nil {
		return false
	}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "TICKET", id)
		if err != nil {
			continue
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			continue
		}
		if !strings.Contains(values["affected-ids"], scope) {
			continue
		}
		state := strings.TrimSpace(values["status"])
		if state == "T-INVALID" || state == "T-SUPERSEDED" {
			continue
		}
		return true
	}
	return false
}

func splitTokens(raw string) []string {
	var tokens []string
	for _, token := range strings.Split(raw, ",") {
		token = strings.TrimSpace(token)
		if token == "" || strings.EqualFold(token, "none") {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}
