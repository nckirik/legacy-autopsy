package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// CoverageDirectives are runtime-owned scope values for coverage expansion.
type CoverageDirectives struct {
	Namespace   string
	Environment string
	Snapshot    string
}

// ExpandCoverage mechanically creates one SOURCE-COVERAGE-ROW per committed
// source unit from the immutable identity tuple
// `src-version | persona-owner | track | environment | snapshot`. The COV ID is
// runtime-computed from that tuple; the row starts `Unmapped` (fail-inclusive)
// in `In-Scope`, and any committed claims for the subject are bound into the
// row. Existing rows are left untouched.
func ExpandCoverage(db string, doc *cdl.EIRDoc, directives CoverageDirectives) (int, error) {
	inventoryNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return 0, err
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return 0, err
	}
	claimsBySubject, err := claimsBySubject(db, doc)
	if err != nil {
		return 0, err
	}
	existing := map[string]bool{}
	existingIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	for _, id := range existingIDs {
		existing[id] = true
	}
	var records []store.RecordInput
	for _, srcID := range ids {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", srcID)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(inventoryNames, payload)
		if err != nil {
			return 0, err
		}
		owner := values["canonical-owner"]
		if strings.TrimSpace(owner) == "" {
			owner = "None"
		}
		tuple := strings.Join([]string{srcID, "version 1", owner, "Normal", directives.Environment, directives.Snapshot}, "|")
		covID, err := capabilities.TypedID("COV", directives.Namespace, tuple, "coverage-row")
		if err != nil {
			return 0, err
		}
		if existing[covID] {
			continue
		}
		claimIDs := claimsBySubject[srcID]
		if len(claimIDs) == 0 {
			claimIDs = []string{"None"}
		}
		var b strings.Builder
		fmt.Fprintf(&b, "- **cov-id:** %s\n", covID)
		fmt.Fprintf(&b, "- **src-version:** %s version 1\n", srcID)
		fmt.Fprintf(&b, "- **persona-owner:** %s\n", owner)
		b.WriteString("- **track:** Normal\n")
		fmt.Fprintf(&b, "- **environment:** %s\n", directives.Environment)
		fmt.Fprintf(&b, "- **snapshot:** %s\n", directives.Snapshot)
		b.WriteString("- **scope:** In-Scope\n")
		b.WriteString("- **disposition:** Unmapped\n")
		b.WriteString("- **c-status:** C-GAP\n")
		b.WriteString("- **exclusion-dec:** None\n")
		b.WriteString("- **mappings:** None\n")
		fmt.Fprintf(&b, "- **frontier-ticket-claim-evidence:** %s\n", strings.Join(claimIDs, ", "))
		records = append(records, store.RecordInput{RecordID: covID, Version: 1, Payload: []byte(b.String())})
	}
	if len(records) == 0 {
		return 0, nil
	}
	if err := store.IngestMany(db, doc, "SOURCE-COVERAGE-ROW", records); err != nil {
		return 0, err
	}
	return len(records), nil
}

// Summary is the deterministic coverage arithmetic over committed inventory.
type Summary struct {
	SourcesByKind         map[string]int
	CoverageByDisposition map[string]int
}

// CoverageSummary computes per-kind source counts and per-disposition coverage
// counts. It is the deterministic substrate for gate arithmetic.
func CoverageSummary(db string, doc *cdl.EIRDoc) (Summary, error) {
	summary := Summary{SourcesByKind: map[string]int{}, CoverageByDisposition: map[string]int{}}
	inventoryNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return summary, err
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return summary, err
	}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return summary, err
		}
		values, err := work.PayloadText(inventoryNames, payload)
		if err != nil {
			return summary, err
		}
		summary.SourcesByKind[values["kind"]]++
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return summary, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return summary, err
	}
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return summary, err
		}
		values, err := work.PayloadText(coverageNames, payload)
		if err != nil {
			return summary, err
		}
		summary.CoverageByDisposition[values["disposition"]]++
	}
	return summary, nil
}

// SortedKeys returns the sorted keys of a count map for deterministic reports.
func SortedKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func claimsBySubject(db string, doc *cdl.EIRDoc) (map[string][]string, error) {
	names, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	bySubject := map[string][]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "CLAIM", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		subject := firstTokenField(values["subject-id-record-version"])
		if subject == "" {
			continue
		}
		bySubject[subject] = append(bySubject[subject], id)
	}
	for subject := range bySubject {
		sort.Strings(bySubject[subject])
	}
	return bySubject, nil
}
