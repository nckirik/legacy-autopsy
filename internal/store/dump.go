package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// CanonicalDump renders a deterministic SQL dump of the store and returns it with
// its sha256 fingerprint. Tables are ordered by name; rows by la_record_id; columns
// by the declared EIR field order after the fixed metadata columns. The fingerprint
// covers the complete dump including its provenance header.
func CanonicalDump(dbPath string, doc *cdl.EIRDoc, prov Provenance) ([]byte, string, error) {
	specs, err := tableSpecs(doc)
	if err != nil {
		return nil, "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "-- GENERATED canonical store dump by %s - do not edit\n", prov.Generator)
	fmt.Fprintf(&b, "-- schema-version: %d\n", SchemaVersion)
	fmt.Fprintf(&b, "-- source-fingerprint: %s\n", prov.SourceFingerprint)
	b.WriteString("BEGIN;\n")
	for _, spec := range specs {
		columns := []string{"la_record_id", "la_record_version", "la_semantic_fingerprint", "la_canonical_payload"}
		for _, f := range spec.Field.Fields {
			columns = append(columns, ident(f.Name))
		}
		quoted := make([]string, len(columns))
		for i, c := range columns {
			quoted[i] = fmt.Sprintf("%q", c)
		}
		query := fmt.Sprintf("select %s from %q order by %q;", strings.Join(quoted, ", "), spec.Name, "la_record_id")
		cmd := exec.Command("sqlite3", "--json", dbPath, query)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			return nil, "", fmt.Errorf("store: dump %s: %v: %s", spec.Name, err, strings.TrimSpace(stderr.String()))
		}
		trimmed := strings.TrimSpace(string(out))
		if trimmed == "" || trimmed == "[]" {
			continue
		}
		var rows []map[string]any
		dec := json.NewDecoder(strings.NewReader(trimmed))
		dec.UseNumber()
		if err := dec.Decode(&rows); err != nil {
			return nil, "", fmt.Errorf("store: parse dump rows for %s: %w", spec.Name, err)
		}
		for _, row := range rows {
			values := make([]string, len(columns))
			for i, c := range columns {
				values[i] = literal(row[c])
			}
			fmt.Fprintf(&b, "INSERT INTO %q (%s) VALUES (%s);\n", spec.Name, strings.Join(quoted, ", "), strings.Join(values, ", "))
		}
	}
	b.WriteString("COMMIT;\n")
	dump := []byte(b.String())
	return dump, cdl.Fingerprint(dump), nil
}

func literal(v any) string {
	switch value := v.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + strings.ReplaceAll(value, "'", "''") + "'"
	case json.Number:
		return value.String()
	case bool:
		if value {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("'%v'", value)
	}
}
