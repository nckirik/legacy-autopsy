package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SchemaCoverageDiagnostics reports discovered schema-object and output units
// that lack committed coverage, so silent under-capture becomes visible:
//
//   - an SRC-TABLE or SRC-VIEW unit must carry at least one committed claim and
//     that claim must be bound by at least one ENTITY-RECORD;
//   - an SRC-DR unit must carry at least one committed claim bound by at least
//     one DATABASE-ROUTINE record;
//   - an SRC-ACTION unit (user-visible output/export behavior) must carry at
//     least one committed claim so the behavior is not invisible.
//
// Diagnostics are deterministic and sorted; an empty result means every
// discovered schema/output unit is represented.
func SchemaCoverageDiagnostics(db string, doc *cdl.EIRDoc) ([]string, error) {
	inventory, err := LoadInventory(db, doc)
	if err != nil {
		return nil, err
	}
	claimsBySubject, err := claimsBySubject(db, doc)
	if err != nil {
		return nil, err
	}
	erBindings, err := recordBindings(db, doc, "ENTITY-RECORD")
	if err != nil {
		return nil, err
	}
	drBindings, err := recordBindings(db, doc, "DATABASE-ROUTINE")
	if err != nil {
		return nil, err
	}
	dispositions, err := coverageDispositions(db, doc)
	if err != nil {
		return nil, err
	}
	var diagnostics []string
	for _, unit := range inventory {
		switch dispositions[unit.ID] {
		case "Approved-Excluded", "Superseded", "Retired":
			continue
		}
		kind := sourceKind(unit.ID)
		claims := claimsBySubject[unit.ID]
		switch kind {
		case "TABLE", "VIEW":
			if len(claims) == 0 {
				diagnostics = append(diagnostics, fmt.Sprintf("declared %s %s has no committed claim", strings.ToLower(kind), unit.Coordinate))
				continue
			}
			if !anyBinding(erBindings, claims) {
				diagnostics = append(diagnostics, fmt.Sprintf("declared %s %s is not represented by an entity record", strings.ToLower(kind), unit.Coordinate))
			}
		case "DR":
			if len(claims) == 0 {
				diagnostics = append(diagnostics, fmt.Sprintf("declared database routine %s has no committed claim", unit.Coordinate))
				continue
			}
			if !anyBinding(drBindings, claims) {
				diagnostics = append(diagnostics, fmt.Sprintf("declared database routine %s is not represented by a database-routine record", unit.Coordinate))
			}
		case "ACTION":
			if len(claims) == 0 {
				diagnostics = append(diagnostics, fmt.Sprintf("discovered output action %s has no committed claim", unit.Coordinate))
			}
		}
	}
	sort.Strings(diagnostics)
	return diagnostics, nil
}

// coverageDispositions maps each source unit to its scope-specific coverage
// disposition, so an exactly approved exclusion (or a superseded/retired unit)
// resolves a schema/output unit without further coverage records.
func coverageDispositions(db string, doc *cdl.EIRDoc) (map[string]string, error) {
	names, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	dispositions := map[string]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		src := firstTokenField(values["src-version"])
		if src == "" {
			continue
		}
		dispositions[src] = strings.TrimSpace(values["disposition"])
	}
	return dispositions, nil
}

// sourceKind extracts the SRC kind token from a source-unit ID.
func sourceKind(id string) string {
	parts := strings.SplitN(id, "-", 3)
	if len(parts) != 3 || parts[0] != "SRC" {
		return ""
	}
	return parts[1]
}

// recordBindings loads the declared source bindings of every committed record
// of one field, in bounded form.
func recordBindings(db string, doc *cdl.EIRDoc, fieldID string) ([]string, error) {
	names, _, _, err := work.FieldContract(doc, fieldID)
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, fieldID)
	if err != nil {
		return nil, err
	}
	var bindings []string
	for _, id := range ids {
		payload, err := store.Get(db, doc, fieldID, id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		for _, field := range []string{"source-bindings", "material-claims"} {
			for _, token := range splitTokens(values[field]) {
				if strings.HasPrefix(token, "CLM-") {
					bindings = append(bindings, token)
				}
			}
		}
	}
	return bindings, nil
}

func anyBinding(bindings, claims []string) bool {
	set := map[string]bool{}
	for _, binding := range bindings {
		set[binding] = true
	}
	for _, claim := range claims {
		if set[claim] {
			return true
		}
	}
	return false
}
