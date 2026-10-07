package store

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// RewriteRecord updates declared column values in place, regenerates the
// canonical payload and fingerprint, and commits one row. It is a runtime-owned
// mechanical mutation for derived state (for example coverage dispositions);
// executors never call it.
func RewriteRecord(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string, updates map[string]string) error {
	field, err := tableField(doc, fieldID)
	if err != nil {
		return err
	}
	payload, err := Get(dbPath, doc, fieldID, recordID)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSuffix(string(payload), "\n"), "\n")
	applied := map[string]bool{}
	for name := range updates {
		applied[name] = false
	}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- **") {
			continue
		}
		end := strings.Index(trimmed[2:], ":**")
		if end < 0 {
			continue
		}
		label := trimmed[4 : end+2]
		column := capabilities.ColumnName(label)
		for name, value := range updates {
			if capabilities.ColumnName(name) != column {
				continue
			}
			if applied[name] {
				return fmt.Errorf("store: rewrite %s %s: field %q appears more than once", fieldID, recordID, name)
			}
			applied[name] = true
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = indent + "- **" + label + ":** " + value
		}
	}
	for name, seen := range applied {
		if !seen {
			return fmt.Errorf("store: rewrite %s %s: field %q is not present in the payload", fieldID, recordID, name)
		}
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(strings.Join(lines, "\n") + "\n"))
	if err != nil {
		return err
	}
	setClauses := []string{
		"la_semantic_fingerprint = " + literal(capabilities.Hash(canonical)),
		"la_canonical_payload = " + literal(string(canonical)),
	}
	for name, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%q = %s", ident(name), literal(value)))
	}
	sql := fmt.Sprintf("UPDATE %q SET %s WHERE la_record_id = %s;", ident(field.ID), strings.Join(setClauses, ", "), literal(recordID))
	return execSQL(dbPath, sql)
}

// ReplaceRecord atomically replaces a record's payload and version, regenerating
// the canonical payload, fingerprint, and derived columns. It is a
// runtime-owned mechanical mutation for explicit re-synthesis with changed
// content at the same evidence-anchored identity.
func ReplaceRecord(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string, version int, payload []byte) error {
	field, err := tableField(doc, fieldID)
	if err != nil {
		return err
	}
	if version < 1 {
		return fmt.Errorf("store: replace %s %s: version must be positive", fieldID, recordID)
	}
	record := RecordInput{RecordID: recordID, Version: version, Payload: payload}
	statement, err := ingestStatement(doc, field, record)
	if err != nil {
		return err
	}
	script := "BEGIN;\nDELETE FROM " + fmt.Sprintf("%q", ident(field.ID)) + " WHERE la_record_id = " + literal(recordID) + ";\n" + statement + "\nCOMMIT;\n"
	return execSQL(dbPath, script)
}
