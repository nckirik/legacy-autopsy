package assurance

import (
	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// ReconcilePlan builds the mechanical unit that declares export reconciliation
// state. When no export artifacts exist in scope, the unit records an
// evidence-backed zero domain instead of fabricating rows.
func ReconcilePlan(namespace string) work.Unit {
	return work.Unit{
		ID:              "U-ASSURANCE-RECONCILE",
		Mode:            "Export Acquisition",
		Scope:           namespace,
		RecordID:        "reconciliation",
		ResultType:      "RECONCILIATION-ROW",
		AllowZeroDomain: true,
		Sections:        []string{"source-coverage"},
	}
}

// MechanicalReconcileExecutor declares the zero domain with the scanned
// inventory as evidence.
type MechanicalReconcileExecutor struct {
	SourceIDs []string
}

// Execute implements work.Executor.
func (e MechanicalReconcileExecutor) Execute(invocation work.Invocation) (work.Proposal, error) {
	return work.Proposal{
		FieldID: "RECONCILIATION-ROW",
		ZeroDomain: &work.ZeroDomain{
			Reason:   "no export artifacts are present in the scanned scope. The export reconciliation family is empty.",
			Evidence: e.SourceIDs,
		},
	}, nil
}

// ReconciliationSources returns the committed source-unit IDs used as zero-domain
// evidence.
func ReconciliationSources(db string, doc *cdl.EIRDoc) ([]string, error) {
	return store.List(db, doc, "SOURCE-INVENTORY-RECORD")
}
