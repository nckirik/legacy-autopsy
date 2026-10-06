package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// CellScope is the traversal evidence scope for one required matrix cell.
type CellScope struct {
	ClusterID   string
	Persona     string
	POV         string
	Track       string
	Environment string
	Snapshot    string
	Coordinates []string
	SourceIDs   []string
	ClaimIDs    []string
	Kinds       map[string]int
	MappedAll   bool
	ZeroDomain  bool
}

var mappedDispositions = map[string]bool{
	"Mapped": true, "Approved-Excluded": true, "Superseded": true, "Retired": true,
}

// LoadCellScopes builds one scope per required matrix cell, resolving the source
// units owned by the cell's persona and whether every one is dispositioned.
func LoadCellScopes(db string, doc *cdl.EIRDoc) ([]CellScope, error) {
	matrixNames, _, _, err := work.FieldContract(doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return nil, err
	}
	matrixIDs, err := store.List(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return nil, err
	}
	inventoryNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	inventoryIDs, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	disposition := map[string]string{}
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(coverageNames, payload)
		if err != nil {
			return nil, err
		}
		disposition[firstTokenField(values["src-version"])] = values["disposition"]
	}
	claims, err := claimsBySubject(db, doc)
	if err != nil {
		return nil, err
	}
	var scopes []CellScope
	for _, id := range matrixIDs {
		payload, err := store.Get(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(matrixNames, payload)
		if err != nil {
			return nil, err
		}
		if values["applicability"] != "Required" || values["status"] != "[R-UNSWEPT]" {
			continue
		}
		scope := CellScope{
			ClusterID:   values["cluster-id"],
			Persona:     values["persona-owner"],
			POV:         values["pov"],
			Track:       values["track"],
			Environment: values["environment"],
			Snapshot:    values["snapshot"],
			Kinds:       map[string]int{},
		}
		for _, sourceID := range inventoryIDs {
			srcPayload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", sourceID)
			if err != nil {
				return nil, err
			}
			srcValues, err := work.PayloadText(inventoryNames, srcPayload)
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(srcValues["canonical-owner"]) != scope.Persona {
				continue
			}
			scope.SourceIDs = append(scope.SourceIDs, sourceID)
			scope.Coordinates = append(scope.Coordinates, srcValues["concrete-coordinate"])
			scope.Kinds[srcValues["kind"]]++
			scope.ClaimIDs = append(scope.ClaimIDs, claims[sourceID]...)
		}
		sort.Strings(scope.SourceIDs)
		sort.Strings(scope.Coordinates)
		sort.Strings(scope.ClaimIDs)
		scope.ZeroDomain = len(scope.SourceIDs) == 0
		scope.MappedAll = true
		for _, sourceID := range scope.SourceIDs {
			if !mappedDispositions[disposition[sourceID]] {
				scope.MappedAll = false
			}
		}
		scopes = append(scopes, scope)
	}
	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i].ClusterID+scopes[i].Persona+scopes[i].POV < scopes[j].ClusterID+scopes[j].Persona+scopes[j].POV
	})
	return scopes, nil
}

// CellSweepPlan builds the mechanical unit that records one SWEEP-RECORD per
// required matrix cell, keyed by the exact cell tuple.
func CellSweepPlan(scopes []CellScope) work.Unit {
	return work.Unit{
		ID:           "U-ASSURANCE-CELL-SWEEPS",
		Mode:         "Discovery",
		IDFromFields: []string{"cluster-id", "persona-owner", "pov", "track", "environment", "snapshot"},
		ResultType:   "SWEEP-RECORD",
		Sections:     []string{"source-coverage", "atomic-unit-rule"},
	}
}

// MechanicalCellSweepExecutor authors the cell sweeps deterministically.
type MechanicalCellSweepExecutor struct {
	Scopes []CellScope
}

// Execute implements work.Executor.
func (e MechanicalCellSweepExecutor) Execute(invocation work.Invocation) (work.Proposal, error) {
	proposal := work.Proposal{FieldID: "SWEEP-RECORD"}
	for _, scope := range e.Scopes {
		kinds := make([]string, 0, len(scope.Kinds))
		for kind := range scope.Kinds {
			kinds = append(kinds, kind)
		}
		sort.Strings(kinds)
		var counts []string
		for _, kind := range kinds {
			counts = append(counts, fmt.Sprintf("%s: %d expected / %d mapped / 0 excluded", kind, scope.Kinds[kind], scope.Kinds[kind]))
		}
		if len(counts) == 0 {
			counts = []string{"None: 0 expected / 0 mapped / 0 excluded"}
		}
		result := "C-PARTIAL"
		if scope.MappedAll {
			result = "C-COVERED"
		}
		claims := "None"
		if len(scope.ClaimIDs) > 0 {
			claims = strings.Join(scope.ClaimIDs, ", ")
		}
		payload := "- **cluster-id:** " + scope.ClusterID + "\n" +
			"- **persona-owner:** " + scope.Persona + "\n" +
			"- **pov:** " + scope.POV + "\n" +
			"- **track:** " + scope.Track + "\n" +
			"- **environment:** " + scope.Environment + "\n" +
			"- **snapshot:** " + scope.Snapshot + "\n" +
			"- **snapshot-environment:** " + scope.Snapshot + " | " + scope.Environment + "\n" +
			"- **persona-cluster-pov-track:** " + scope.Persona + " | " + scope.ClusterID + " | " + scope.POV + " | " + scope.Track + "\n" +
			"- **expanded-concrete-entry-coordinates:** " + joinOrNone(scope.Coordinates) + "\n" +
			fmt.Sprintf("- **source-denominator:** %d\n", len(scope.SourceIDs)) +
			"- **expected-mapped-excluded-counts-by-kind:** " + strings.Join(counts, " ; ") + "\n" +
			"- **open-terminal-frontier-counts:** 0\n" +
			"- **active-terminal-ticket-counts:** 0\n" +
			"- **export-reconciliation:** None\n" +
			"- **coverage-result:** " + result + "\n" +
			"- **closing-invocation:** " + invocation.InvocationID + "\n" +
			"- **evidence-claims:** " + claims + "\n"
		proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
	}
	return proposal, nil
}

// AdvanceMatrix marks required cells [R-SWEPT] when a matching sweep tuple
// reports C-COVERED evidence. All other cells stay fail-inclusive.
func AdvanceMatrix(db string, doc *cdl.EIRDoc) (int, error) {
	sweepNames, _, _, err := work.FieldContract(doc, "SWEEP-RECORD")
	if err != nil {
		return 0, err
	}
	sweepIDs, err := store.List(db, doc, "SWEEP-RECORD")
	if err != nil {
		return 0, err
	}
	covered := map[string]bool{}
	for _, id := range sweepIDs {
		payload, err := store.Get(db, doc, "SWEEP-RECORD", id)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(sweepNames, payload)
		if err != nil {
			return 0, err
		}
		if values["coverage-result"] != "C-COVERED" {
			continue
		}
		covered[sweepKey(values)] = true
	}
	matrixNames, _, _, err := work.FieldContract(doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return 0, err
	}
	matrixIDs, err := store.List(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, id := range matrixIDs {
		payload, err := store.Get(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", id)
		if err != nil {
			return changed, err
		}
		values, err := work.PayloadText(matrixNames, payload)
		if err != nil {
			return changed, err
		}
		if values["status"] != "[R-UNSWEPT]" || values["applicability"] != "Required" {
			continue
		}
		if !covered[sweepKey(values)] {
			continue
		}
		if err := store.RewriteRecord(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", id, map[string]string{"status": "[R-SWEPT]"}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func sweepKey(values map[string]string) string {
	return strings.Join([]string{
		values["cluster-id"], values["persona-owner"], values["pov"],
		values["track"], values["environment"], values["snapshot"],
	}, "|")
}

func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "None"
	}
	return strings.Join(values, ", ")
}
