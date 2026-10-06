package preflight

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// MatrixDirectives are runtime-owned scope values for matrix expansion.
type MatrixDirectives struct {
	Environment string
	Snapshot    string
}

// MatrixPOVs are the five traversal POV dimensions.
var MatrixPOVs = []string{"POV-1", "POV-2", "POV-3", "POV-4", "POV-5"}

// ExpandMatrix mechanically creates the Required Traversal Matrix rows for every
// committed cluster, persona, and POV: each cell is a keyed
// REQUIRED-TRAVERSAL-MATRIX-ROW record (fail-inclusive `Required`/`[R-UNSWEPT]`).
// Existing rows are left untouched, so the operation is idempotent.
func ExpandMatrix(db string, doc *cdl.EIRDoc, directives MatrixDirectives) (int, error) {
	clusterNames, _, _, err := work.FieldContract(doc, "ENTRY-CLUSTER")
	if err != nil {
		return 0, err
	}
	clusterIDs, err := store.List(db, doc, "ENTRY-CLUSTER")
	if err != nil {
		return 0, err
	}
	personas, err := store.List(db, doc, "PERSONA-REGISTRY")
	if err != nil {
		return 0, err
	}
	existing := map[string]bool{}
	existingIDs, err := store.List(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return 0, err
	}
	for _, id := range existingIDs {
		existing[id] = true
	}
	var records []store.RecordInput
	for _, clusterID := range clusterIDs {
		payload, err := store.Get(db, doc, "ENTRY-CLUSTER", clusterID)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(clusterNames, payload)
		if err != nil {
			return 0, err
		}
		track := values["traversal-track"]
		if track == "" {
			track = "Normal"
		}
		concreteEntry := values["concrete-entry-coordinates"]
		if strings.TrimSpace(concreteEntry) == "" {
			concreteEntry = "None"
		}
		for _, persona := range personas {
			for _, pov := range MatrixPOVs {
				recordID := strings.Join([]string{clusterID, persona, pov, track, directives.Environment, directives.Snapshot}, "|")
				if existing[recordID] {
					continue
				}
				var b strings.Builder
				fmt.Fprintf(&b, "- **cluster-id:** %s\n", clusterID)
				fmt.Fprintf(&b, "- **persona-owner:** %s\n", persona)
				fmt.Fprintf(&b, "- **concrete-entry:** %s\n", concreteEntry)
				fmt.Fprintf(&b, "- **track:** %s\n", track)
				fmt.Fprintf(&b, "- **environment:** %s\n", directives.Environment)
				fmt.Fprintf(&b, "- **snapshot:** %s\n", directives.Snapshot)
				fmt.Fprintf(&b, "- **pov:** %s\n", pov)
				b.WriteString("- **applicability:** Required\n")
				b.WriteString("- **status:** [R-UNSWEPT]\n")
				b.WriteString("- **swp-id:** None\n")
				b.WriteString("- **claims:** None\n")
				records = append(records, store.RecordInput{RecordID: recordID, Version: 1, Payload: []byte(b.String())})
			}
		}
	}
	if len(records) == 0 {
		return 0, nil
	}
	if err := store.IngestMany(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", records); err != nil {
		return 0, err
	}
	return len(records), nil
}
