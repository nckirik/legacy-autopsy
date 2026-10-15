package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// tableField resolves a FIELD ID that has a table storage class.
func tableField(doc *cdl.EIRDoc, fieldID string) (*cdl.EIRField, error) {
	field, err := findField(doc, fieldID)
	if err != nil {
		return nil, err
	}
	switch field.Storage {
	case "record", "registry", "singleton", "log", "child":
		return field, nil
	default:
		return nil, fmt.Errorf("store: %s: STORAGE %s is not a table class", fieldID, field.Storage)
	}
}

// Get returns the exact canonical payload bytes stored for one record. These are
// the §4.1.2 canonical bytes; queryable columns are derived and never returned as
// the record's identity.
func Get(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string) ([]byte, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return nil, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select la_canonical_payload from %q where la_record_id = %s;", ident(fieldID), literal(recordID)))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("store: %s %s not found", fieldID, recordID)
	}
	payload, ok := rows[0]["la_canonical_payload"].(string)
	if !ok {
		return nil, fmt.Errorf("store: %s %s has a non-text canonical payload", fieldID, recordID)
	}
	return []byte(payload), nil
}

// List returns the sorted record IDs of one table-class FIELD.
func List(dbPath string, doc *cdl.EIRDoc, fieldID string) ([]string, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return nil, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select la_record_id from %q order by la_record_id;", ident(fieldID)))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		id, ok := row["la_record_id"].(string)
		if !ok {
			return nil, fmt.Errorf("store: %s has a non-text record ID", fieldID)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// FingerprintOf returns the stored semantic fingerprint of one record. It is the
// value read-set provenance compares against at commit time.
func FingerprintOf(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string) (string, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return "", err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select la_semantic_fingerprint from %q where la_record_id = %s;", ident(fieldID), literal(recordID)))
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", fmt.Errorf("store: %s %s not found", fieldID, recordID)
	}
	fingerprint, ok := rows[0]["la_semantic_fingerprint"].(string)
	if !ok {
		return "", fmt.Errorf("store: %s %s has no fingerprint", fieldID, recordID)
	}
	return fingerprint, nil
}

// Exists reports whether a record ID is present in a table-class FIELD.
func Exists(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string) (bool, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return false, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select 1 as present from %q where la_record_id = %s limit 1;", ident(fieldID), literal(recordID)))
	if err != nil {
		return false, err
	}
	return len(rows) == 1, nil
}

// Search returns the sorted record IDs of a table whose canonical payload
// contains the exact needle.
func Search(dbPath string, doc *cdl.EIRDoc, fieldID, needle string) ([]string, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return nil, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select la_record_id from %q where instr(la_canonical_payload, %s) > 0 order by la_record_id;", ident(fieldID), literal(needle)))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		id, ok := row["la_record_id"].(string)
		if !ok {
			return nil, fmt.Errorf("store: %s has a non-text record ID", fieldID)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// Render returns a clean human-readable view of every record in a table: the
// canonical payloads, sorted by record ID, with no mechanical metadata. It is a
// display projection, not an operational record.
func Render(dbPath string, doc *cdl.EIRDoc, fieldID string) ([]byte, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return nil, err
	}
	ids, err := List(dbPath, doc, fieldID)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", fieldID)
	for _, id := range ids {
		payload, err := Get(dbPath, doc, fieldID, id)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&b, "## %s\n\n%s\n", id, ensureTrailingNewline(string(payload)))
	}
	return []byte(b.String()), nil
}

// Inspect renders a debug-only, non-authoritative Markdown view of every row in a
// table, including the mechanical la_* metadata. It is a diagnostic render, never
// an operational record or a package member.
func Inspect(dbPath string, doc *cdl.EIRDoc, fieldID string) ([]byte, error) {
	field, err := tableField(doc, fieldID)
	if err != nil {
		return nil, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select * from %q order by la_record_id;", ident(fieldID)))
	if err != nil {
		return nil, err
	}
	columns := []string{"la_record_id", "la_record_version", "la_semantic_fingerprint"}
	for _, spec := range field.Fields {
		columns = append(columns, ident(spec.Name))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "<!-- DEBUG RENDER - NON-AUTHORITATIVE: mechanical store state -->\n")
	fmt.Fprintf(&b, "# %s (debug)\n\n- **rows:** %d\n\n", fieldID, len(rows))
	for _, row := range rows {
		id, _ := row["la_record_id"].(string)
		fmt.Fprintf(&b, "## %s\n\n", id)
		for _, column := range columns {
			fmt.Fprintf(&b, "- **%s:** %s\n", column, displayValue(row[column]))
		}
		payload, _ := row["la_canonical_payload"].(string)
		fmt.Fprintf(&b, "\n- **la_canonical_payload:**\n\n```markdown\n%s```\n\n", ensureTrailingNewline(payload))
	}
	return []byte(b.String()), nil
}

func displayValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "None"
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ensureTrailingNewline(value string) string {
	if strings.HasSuffix(value, "\n") {
		return value
	}
	return value + "\n"
}

// VersionOf returns the stored record version.
func VersionOf(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string) (int, error) {
	if _, err := tableField(doc, fieldID); err != nil {
		return 0, err
	}
	rows, err := queryRows(dbPath, fmt.Sprintf("select la_record_version from %q where la_record_id = %s;", ident(fieldID), literal(recordID)))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, fmt.Errorf("store: %s %s not found", fieldID, recordID)
	}
	number, ok := rows[0]["la_record_version"].(json.Number)
	if !ok {
		return 0, fmt.Errorf("store: %s %s has a non-numeric version", fieldID, recordID)
	}
	version, err := number.Int64()
	if err != nil {
		return 0, err
	}
	return int(version), nil
}
