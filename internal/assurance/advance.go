package assurance

import (
	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// AdvanceCoverage recomputes source coverage dispositions from committed
// DERIVED-FROM-EDGE records: a source unit referenced by at least one synthesis
// binding becomes `Mapped`. Other dispositions (exclusions, human blocks) are
// preserved. It is a runtime-owned mechanical derivation.
func AdvanceCoverage(db string, doc *cdl.EIRDoc) (int, error) {
	edgeNames, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return 0, err
	}
	edgeIDs, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return 0, err
	}
	mapped := map[string]bool{}
	for _, id := range edgeIDs {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(edgeNames, payload)
		if err != nil {
			return 0, err
		}
		source := firstTokenField(values["source-record-id"])
		if source != "" && source != "None" {
			mapped[source] = true
		}
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return changed, err
		}
		values, err := work.PayloadText(coverageNames, payload)
		if err != nil {
			return changed, err
		}
		source := firstTokenField(values["src-version"])
		if !mapped[source] {
			continue
		}
		switch values["disposition"] {
		case "Mapped", "Superseded", "Retired", "Approved-Excluded", "Human-Blocked":
			continue
		}
		if err := store.RewriteRecord(db, doc, "SOURCE-COVERAGE-ROW", id, map[string]string{"disposition": "Mapped"}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}
