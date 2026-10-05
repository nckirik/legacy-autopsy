package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// DumpVersion identifies the canonical dump serialization.
const DumpVersion = 2

// CanonicalDump renders a deterministic, restorable SQL dump of the store: the
// generated relational schema followed by all rows. Before emitting anything it
// verifies the live database schema against the EIR-derived expected schema, so a
// modified or partial database cannot produce a certified dump.
//
// Tables are ordered by name; rows by la_record_id; columns by the declared EIR
// field order after the fixed metadata columns. The returned fingerprint covers
// only the canonical content block (schema plus rows, BEGIN through COMMIT), not
// the provenance header, so tooling metadata changes do not alter content
// identity. The provenance header remains part of the returned bytes.
func CanonicalDump(dbPath string, doc *cdl.EIRDoc, prov Provenance) ([]byte, string, error) {
	stmts, err := Statements(doc)
	if err != nil {
		return nil, "", err
	}
	if err := verifyLiveSchema(dbPath, stmts); err != nil {
		return nil, "", err
	}
	specs, err := tableSpecs(doc)
	if err != nil {
		return nil, "", err
	}

	var content strings.Builder
	content.WriteString("BEGIN;\n")
	for _, stmt := range stmts {
		content.WriteString(stmt.SQL)
		content.WriteString("\n")
	}
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
			fmt.Fprintf(&content, "INSERT INTO %q (%s) VALUES (%s);\n", spec.Name, strings.Join(quoted, ", "), strings.Join(values, ", "))
		}
	}
	content.WriteString("COMMIT;\n")

	var b strings.Builder
	fmt.Fprintf(&b, "-- GENERATED canonical store dump by %s - do not edit\n", prov.Generator)
	fmt.Fprintf(&b, "-- schema-version: %d\n", SchemaVersion)
	fmt.Fprintf(&b, "-- dump-format: %d\n", DumpVersion)
	fmt.Fprintf(&b, "-- source-fingerprint: %s\n", prov.SourceFingerprint)
	b.WriteString(content.String())
	dump := []byte(b.String())
	return dump, cdl.Fingerprint([]byte(content.String())), nil
}

type liveObject struct {
	Type string `json:"type"`
	Name string `json:"name"`
	SQL  string `json:"sql"`
}

// verifyLiveSchema fails closed unless the database contains exactly the tables,
// triggers, and indexes the generated DDL declares, with matching SQL normalized
// for whitespace and the trailing semicolon. This binds CHECK, foreign-key, index,
// and trigger definitions into the certified dump. SQLite-internal objects
// (sqlite_*) are ignored.
func verifyLiveSchema(dbPath string, expected []Statement) error {
	query := "select type, name, sql from sqlite_master where type in ('table','trigger','index') and name not like 'sqlite_%' order by name;"
	cmd := exec.Command("sqlite3", "--json", dbPath, query)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("store: inspect live schema: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	var live []liveObject
	if trimmed := strings.TrimSpace(string(out)); trimmed != "" && trimmed != "[]" {
		if err := json.Unmarshal([]byte(trimmed), &live); err != nil {
			return fmt.Errorf("store: parse live schema: %w", err)
		}
	}
	want := map[string]string{}
	for _, stmt := range expected {
		want[stmt.Name] = normalizeSQL(stmt.SQL)
	}
	got := map[string]string{}
	for _, obj := range live {
		got[obj.Name] = normalizeSQL(obj.SQL)
	}
	for name, sql := range want {
		liveSQL, ok := got[name]
		if !ok {
			return fmt.Errorf("store: live schema is missing %s %q", objectType(expected, name), name)
		}
		if liveSQL != sql {
			return fmt.Errorf("store: live schema for %q does not match the generated EIR schema", name)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			return fmt.Errorf("store: live schema has object %q not present in the generated EIR schema", name)
		}
	}
	return nil
}

func objectType(stmts []Statement, name string) string {
	for _, stmt := range stmts {
		if stmt.Name == name {
			return stmt.Type
		}
	}
	return "object"
}

func normalizeSQL(sql string) string {
	return strings.TrimSuffix(strings.Join(strings.Fields(sql), " "), ";")
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
