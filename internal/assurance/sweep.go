package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SweepScope is one owner/persona sweep group over committed source units.
type SweepScope struct {
	Owner       string
	Coordinates []string
	SourceIDs   []string
	ClaimIDs    []string
	Kinds       map[string]int
}

// SweepPlan returns one mechanical unit per owner group; the runtime authors the
// sweep records and the 0G entries.
func SweepPlan(scopes []SweepScope, directives CoverageDirectives) work.Unit {
	return work.Unit{
		ID:                   "U-ASSURANCE-SWEEPS",
		Mode:                 "Discovery",
		Scope:                directives.Namespace,
		IDOwnerField:         "persona-cluster-pov-track",
		IDDiscriminatorField: "snapshot-environment",
		IDType:               "SWP",
		Namespace:            directives.Namespace,
		ExpectedRecords:      len(scopes),
		ResultType:           "SWEEP-RECORD",
		Sections:             []string{"source-coverage", "atomic-unit-rule"},
	}
}

// LoadSweepScopes groups committed source units by canonical owner with their
// binding claims and kind counts.
func LoadSweepScopes(db string, doc *cdl.EIRDoc) ([]SweepScope, error) {
	inventoryNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	claims, err := claimsBySubject(db, doc)
	if err != nil {
		return nil, err
	}
	byOwner := map[string]*SweepScope{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(inventoryNames, payload)
		if err != nil {
			return nil, err
		}
		owner := values["canonical-owner"]
		if strings.TrimSpace(owner) == "" {
			owner = "None"
		}
		scope, ok := byOwner[owner]
		if !ok {
			scope = &SweepScope{Owner: owner, Kinds: map[string]int{}}
			byOwner[owner] = scope
		}
		scope.SourceIDs = append(scope.SourceIDs, id)
		scope.Coordinates = append(scope.Coordinates, values["concrete-coordinate"])
		scope.Kinds[values["kind"]]++
		scope.ClaimIDs = append(scope.ClaimIDs, claims[id]...)
	}
	scopes := make([]SweepScope, 0, len(byOwner))
	for _, scope := range byOwner {
		sort.Strings(scope.SourceIDs)
		sort.Strings(scope.Coordinates)
		sort.Strings(scope.ClaimIDs)
		scopes = append(scopes, *scope)
	}
	sort.Slice(scopes, func(i, j int) bool { return scopes[i].Owner < scopes[j].Owner })
	return scopes, nil
}

// MechanicalSweepExecutor authors one sweep record per scope deterministically.
type MechanicalSweepExecutor struct {
	Scopes     []SweepScope
	Directives CoverageDirectives
}

// Execute implements work.Executor.
func (e MechanicalSweepExecutor) Execute(invocation work.Invocation) (work.Proposal, error) {
	proposal := work.Proposal{FieldID: "SWEEP-RECORD"}
	for _, scope := range e.Scopes {
		counts := make([]string, 0, len(scope.Kinds))
		kinds := make([]string, 0, len(scope.Kinds))
		for kind := range scope.Kinds {
			kinds = append(kinds, kind)
		}
		sort.Strings(kinds)
		for _, kind := range kinds {
			counts = append(counts, fmt.Sprintf("%s: %d expected / 0 mapped / 0 excluded", kind, scope.Kinds[kind]))
		}
		claims := "None"
		if len(scope.ClaimIDs) > 0 {
			claims = strings.Join(scope.ClaimIDs, ", ")
		}
		payload := "- **snapshot-environment:** " + e.Directives.Snapshot + " | " + e.Directives.Environment + "\n" +
			"- **persona-cluster-pov-track:** " + scope.Owner + " | None | None | Normal\n" +
			"- **expanded-concrete-entry-coordinates:** " + strings.Join(scope.Coordinates, ", ") + "\n" +
			fmt.Sprintf("- **source-denominator:** %d\n", len(scope.SourceIDs)) +
			"- **expected-mapped-excluded-counts-by-kind:** " + strings.Join(counts, " ; ") + "\n" +
			"- **open-terminal-frontier-counts:** 0\n" +
			"- **active-terminal-ticket-counts:** 0\n" +
			"- **export-reconciliation:** None\n" +
			"- **coverage-result:** C-GAP\n" +
			"- **closing-invocation:** " + invocation.InvocationID + "\n" +
			"- **evidence-claims:** " + claims + "\n"
		proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
	}
	return proposal, nil
}
