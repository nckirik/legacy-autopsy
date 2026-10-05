package store

import (
	"fmt"
	"sort"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// LogFields returns every log-class FIELD ID, sorted.
func LogFields(doc *cdl.EIRDoc) []string {
	var fields []string
	for _, field := range doc.Declarations.Fields {
		if field.Storage == "log" {
			fields = append(fields, field.ID)
		}
	}
	sort.Strings(fields)
	return fields
}

// LogState snapshots the current record IDs of every log-class table. The
// runtime records it before a transaction and re-checks it afterwards so that
// append-only is a runtime post-condition, not only a database trigger.
func LogState(dbPath string, doc *cdl.EIRDoc) (map[string][]string, error) {
	state := map[string][]string{}
	for _, fieldID := range LogFields(doc) {
		ids, err := List(dbPath, doc, fieldID)
		if err != nil {
			return nil, err
		}
		state[fieldID] = ids
	}
	return state, nil
}

// VerifyMonotonicLogs fails closed if any previously present log record ID
// disappeared. Additions are allowed; deletion and replacement are not.
func VerifyMonotonicLogs(before, after map[string][]string) error {
	for fieldID, previous := range before {
		present := map[string]bool{}
		for _, id := range after[fieldID] {
			present[id] = true
		}
		for _, id := range previous {
			if !present[id] {
				return fmt.Errorf("store: append-only violation: %s %s disappeared", fieldID, id)
			}
		}
	}
	return nil
}
