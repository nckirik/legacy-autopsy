package discovery

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SourceUnit is the committed inventory form needed for frontier creation.
type SourceUnit struct {
	ID          string
	Coordinate  string
	Fingerprint string
}

// FrontierPlan builds the single mechanical unit that records one terminal
// frontier per committed source unit. The unit is executed by
// MechanicalFrontierExecutor, not by a model.
func FrontierPlan(inventory []SourceUnit, namespace string) work.Unit {
	return work.Unit{
		ID:               "U-DISCOVERY-FRONTIERS",
		Mode:             "Discovery",
		Scope:            namespace,
		IDFromCoordinate: "target-concrete-coordinate-or-resolution-rule",
		IDType:           "FRT",
		Namespace:        namespace,
		ExpectedRecords:  len(inventory),
		ResultType:       "FRONTIER-RECORD",
		Sections:         []string{"traversal-frontier"},
	}
}

// MechanicalFrontierExecutor authors terminal frontier payloads deterministically.
// Every enumerated file is a terminal static boundary: traversal reached the unit
// and stops at file granularity, so the frontier is closed (state `Traversed`)
// by this invocation.
type MechanicalFrontierExecutor struct {
	Inventory []SourceUnit
}

// Execute implements work.Executor.
func (e MechanicalFrontierExecutor) Execute(invocation work.Invocation) (work.Proposal, error) {
	proposal := work.Proposal{FieldID: "FRONTIER-RECORD"}
	for _, item := range e.Inventory {
		payload := "- **origin-invocation-cluster-pov-track:** " + invocation.InvocationID + " | None | None | Normal\n" +
			"- **from-src-cmp:** " + item.ID + "\n" +
			"- **target-concrete-coordinate-or-resolution-rule:** " + item.Coordinate + "\n" +
			"- **target-src:** " + item.ID + "\n" +
			"- **discovery-method:** other\n" +
			fmt.Sprintf("- **depth:** %d\n", strings.Count(item.Coordinate, "/")+1) +
			"- **snapshot-fingerprint:** " + item.Fingerprint + "\n" +
			"- **state:** Traversed\n" +
			"- **disposition-evidence:** terminal static file unit inventoried; no further file-granularity traversal\n" +
			"- **child-frontier-ids:** None\n" +
			"- **closing-invocation:** " + invocation.InvocationID + "\n"
		proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
	}
	return proposal, nil
}
