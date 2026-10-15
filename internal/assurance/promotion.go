package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// promotionLockStates are the component staging states whose promotion locks
// must hold: verified staging and promoted components.
var promotionLockStates = map[string]bool{
	"[S-STAGING-VERIFIED]": true,
	"[S-PROMOTED]":         true,
}

// ValidatePromotionLocks validates the §6.3 promotion locks for one component:
// usage uniqueness (exactly one persona across all tracks), ambiguity (no
// active ticket and no open material contradiction referencing it), dependency
// (no other discovered component depends on it), coverage (its source units are
// C-COVERED), and evidence (no material [E-UNKNOWN] claim). It returns the
// first violated lock.
func ValidatePromotionLocks(db string, doc *cdl.EIRDoc, cmpID string) error {
	names, _, _, err := work.FieldContract(doc, "ATOMIC-COMPONENT")
	if err != nil {
		return err
	}
	payload, err := store.Get(db, doc, "ATOMIC-COMPONENT", cmpID)
	if err != nil {
		return fmt.Errorf("promotion: component %s: %w", cmpID, err)
	}
	values, err := work.PayloadText(names, payload)
	if err != nil {
		return err
	}

	// Usage uniqueness: exactly one persona prefix.
	personas := []string{}
	for _, token := range splitTokens(values["personas"]) {
		prefix := strings.ToUpper(token)
		if prefix == "" || prefix == "NONE" || prefix == "SHARED" || prefix == "UNKNOWN" {
			continue
		}
		personas = append(personas, prefix)
	}
	if len(personas) != 1 {
		return fmt.Errorf("promotion: %s usage uniqueness: expected exactly one persona, got %v", cmpID, personas)
	}

	// Ambiguity: no active ticket or open contradiction referencing it.
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
		if strings.Contains(ticket["affected-ids"], cmpID) && work.TicketBlocksExitA(ticket["status"]) {
			return fmt.Errorf("promotion: %s ambiguity: active ticket %s (%s)", cmpID, id, strings.TrimSpace(ticket["status"]))
		}
	}
	contradictionNames, _, _, err := work.FieldContract(doc, "CONTRADICTION")
	if err != nil {
		return err
	}
	contradictionIDs, err := store.List(db, doc, "CONTRADICTION")
	if err != nil {
		return err
	}
	sort.Strings(contradictionIDs)
	for _, id := range contradictionIDs {
		contradictionPayload, err := store.Get(db, doc, "CONTRADICTION", id)
		if err != nil {
			return err
		}
		contradiction, err := work.PayloadText(contradictionNames, contradictionPayload)
		if err != nil {
			return err
		}
		if !strings.Contains(contradiction["affected-records-and-stale-impact"], cmpID) {
			continue
		}
		state := strings.TrimSpace(contradiction["resolution-state"])
		switch state {
		case "Open", "Decision-Required":
			return fmt.Errorf("promotion: %s ambiguity: open contradiction %s (%s)", cmpID, id, state)
		}
	}

	// Dependency: no other component depends on it.
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
		if strings.Contains(other["depends-on-depended-on-by"], cmpID) {
			return fmt.Errorf("promotion: %s dependency: component %s depends on it", cmpID, id)
		}
	}

	// Coverage: every source unit is C-COVERED.
	dispositions, err := coverageDispositions(db, doc)
	if err != nil {
		return err
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return err
	}
	statusBySrc := map[string]string{}
	for _, id := range coverageIDs {
		coveragePayload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return err
		}
		coverage, err := work.PayloadText(coverageNames, coveragePayload)
		if err != nil {
			return err
		}
		statusBySrc[firstTokenField(coverage["src-version"])] = strings.TrimSpace(coverage["c-status"])
	}
	for _, source := range splitTokens(values["source-ids"]) {
		if dispositions[source] == "Approved-Excluded" {
			continue
		}
		status, ok := statusBySrc[source]
		if !ok {
			return fmt.Errorf("promotion: %s coverage: source unit %s has no coverage row", cmpID, source)
		}
		if status != "[C-COVERED]" {
			return fmt.Errorf("promotion: %s coverage: source unit %s is %s", cmpID, source, status)
		}
	}

	// Evidence: no material [E-UNKNOWN] claim.
	claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return err
	}
	for _, token := range splitTokens(values["material-claims"]) {
		if !strings.HasPrefix(token, "CLM-") {
			continue
		}
		claimPayload, err := store.Get(db, doc, "CLAIM", token)
		if err != nil {
			return err
		}
		claim, err := work.PayloadText(claimNames, claimPayload)
		if err != nil {
			return err
		}
		material := strings.EqualFold(strings.TrimSpace(claim["materiality"]), "Material")
		unknown := strings.TrimSpace(claim["evidence-level"]) == "[E-UNKNOWN]"
		if material && unknown {
			return fmt.Errorf("promotion: %s evidence: material claim %s is [E-UNKNOWN]", cmpID, token)
		}
	}
	return nil
}

// PromotionLockDiagnostics validates every component declaring a verified or
// promoted staging state and returns sorted diagnostics.
func PromotionLockDiagnostics(db string, doc *cdl.EIRDoc) ([]string, error) {
	names, _, _, err := work.FieldContract(doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	var diagnostics []string
	for _, id := range ids {
		payload, err := store.Get(db, doc, "ATOMIC-COMPONENT", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		if !promotionLockStates[strings.TrimSpace(values["staging-status"])] {
			continue
		}
		if err := ValidatePromotionLocks(db, doc, id); err != nil {
			diagnostics = append(diagnostics, err.Error())
		}
	}
	sort.Strings(diagnostics)
	return diagnostics, nil
}
