package assurance

import (
	"fmt"
	"sort"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// GateCheck is one deterministic gate condition result.
type GateCheck struct {
	ID     string
	Result string // Pass | Fail | Unsupported
	Detail string
}

// GateResult aggregates the Exit A evaluation.
type GateResult struct {
	Passed bool
	Checks []GateCheck
}

// EvaluateExitA computes the gate conditions that the current canonical state
// can prove deterministically. Conditions whose models do not exist yet are
// reported Unsupported, never silently passed.
func EvaluateExitA(db string, doc *cdl.EIRDoc) (GateResult, error) {
	result := GateResult{Passed: true}
	add := func(check GateCheck) {
		if check.Result != "Pass" {
			result.Passed = false
		}
		result.Checks = append(result.Checks, check)
	}

	inventoryIDs, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return result, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return result, err
	}
	frontierIDs, err := store.List(db, doc, "FRONTIER-RECORD")
	if err != nil {
		return result, err
	}
	tickets, err := store.List(db, doc, "TICKET")
	if err != nil {
		return result, err
	}

	if len(inventoryIDs) == 0 {
		add(GateCheck{ID: "COVERAGE-SOURCE-INVENTORY", Result: "Fail", Detail: "no source inventory records exist"})
	} else if len(coverageIDs) == len(inventoryIDs) {
		add(GateCheck{ID: "COVERAGE-SOURCE-INVENTORY", Result: "Pass", Detail: fmt.Sprintf("%d source units, %d coverage rows", len(inventoryIDs), len(coverageIDs))})
	} else {
		add(GateCheck{ID: "COVERAGE-SOURCE-INVENTORY", Result: "Fail", Detail: fmt.Sprintf("%d source units but %d coverage rows", len(inventoryIDs), len(coverageIDs))})
	}

	summary, err := CoverageSummary(db, doc)
	if err != nil {
		return result, err
	}
	mapped := summary.CoverageByDisposition["Mapped"] + summary.CoverageByDisposition["Approved-Excluded"]
	total := len(coverageIDs)
	if total == 0 {
		add(GateCheck{ID: "COVERAGE-DENOMINATOR", Result: "Fail", Detail: "coverage denominator is empty"})
	} else if mapped == total {
		add(GateCheck{ID: "COVERAGE-DENOMINATOR", Result: "Pass", Detail: fmt.Sprintf("%d/%d mapped or excluded", mapped, total)})
	} else {
		add(GateCheck{ID: "COVERAGE-DENOMINATOR", Result: "Fail", Detail: fmt.Sprintf("%d/%d mapped or excluded", mapped, total)})
	}

	openFrontiers := 0
	frontierNames, _, _, err := work.FieldContract(doc, "FRONTIER-RECORD")
	if err != nil {
		return result, err
	}
	for _, id := range frontierIDs {
		payload, err := store.Get(db, doc, "FRONTIER-RECORD", id)
		if err != nil {
			return result, err
		}
		values, err := work.PayloadText(frontierNames, payload)
		if err != nil {
			return result, err
		}
		if values["state"] == "Open" {
			openFrontiers++
		}
	}
	if openFrontiers == 0 {
		add(GateCheck{ID: "FRONTIER-CLOSURE", Result: "Pass", Detail: fmt.Sprintf("%d frontiers, none open", len(frontierIDs))})
	} else {
		add(GateCheck{ID: "FRONTIER-CLOSURE", Result: "Fail", Detail: fmt.Sprintf("%d open frontiers", openFrontiers)})
	}

	if len(tickets) == 0 {
		add(GateCheck{ID: "ACTIVE-TICKETS", Result: "Pass", Detail: "no active terminal tickets"})
	} else {
		active := 0
		for range tickets {
			active++
		}
		add(GateCheck{ID: "ACTIVE-TICKETS", Result: "Fail", Detail: fmt.Sprintf("%d ticket records outstanding", active)})
	}

	matrixNames, _, _, err := work.FieldContract(doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return result, err
	}
	matrixIDs, err := store.List(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return result, err
	}
	closed := 0
	for _, id := range matrixIDs {
		payload, err := store.Get(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", id)
		if err != nil {
			return result, err
		}
		values, err := work.PayloadText(matrixNames, payload)
		if err != nil {
			return result, err
		}
		if values["status"] == "[R-SWEPT]" || (values["applicability"] == "Not-Applicable" && values["status"] == "N/A") {
			closed++
		}
	}
	if len(matrixIDs) == 0 {
		add(GateCheck{ID: "TRAVERSAL-MATRIX-CLOSURE", Result: "Fail", Detail: "no required traversal-matrix rows exist"})
	} else if closed == len(matrixIDs) {
		add(GateCheck{ID: "TRAVERSAL-MATRIX-CLOSURE", Result: "Pass", Detail: fmt.Sprintf("%d/%d matrix cells closed", closed, len(matrixIDs))})
	} else {
		add(GateCheck{ID: "TRAVERSAL-MATRIX-CLOSURE", Result: "Fail", Detail: fmt.Sprintf("%d/%d matrix cells closed", closed, len(matrixIDs))})
	}
	add(GateCheck{ID: "SYNTHESIS-COVERAGE", Result: "Unsupported", Detail: "synthesis records and coverage are not generated yet"})

	sort.Slice(result.Checks, func(i, j int) bool { return result.Checks[i].ID < result.Checks[j].ID })
	return result, nil
}
