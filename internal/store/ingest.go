package store

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/certify"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// RecordInput is one record to ingest.
type RecordInput struct {
	RecordID string
	Version  int
	Payload  []byte
}

// Ingest validates one Markdown record payload against its declared FIELD
// schema, derives the queryable columns, and inserts one row.
func Ingest(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string, version int, payload []byte) error {
	return IngestMany(dbPath, doc, fieldID, []RecordInput{{RecordID: recordID, Version: version, Payload: payload}})
}

// IngestMany validates and inserts records in one transaction, all or nothing.
// The canonical payload bytes are stored verbatim in la_canonical_payload, so
// record payload digests stay binding-independent; queryable columns are derived
// projections rebuildable from the payload. Unknown, duplicate, missing-required,
// wrongly typed, or noncanonical values fail closed, as does a non-table storage
// class, an empty/duplicate record ID, or a non-positive version.
func IngestMany(dbPath string, doc *cdl.EIRDoc, fieldID string, records []RecordInput) error {
	if len(records) == 0 {
		return fmt.Errorf("store: ingest %s: no records", fieldID)
	}
	field, err := tableField(doc, fieldID)
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	var script strings.Builder
	script.WriteString("BEGIN;\n")
	for _, record := range records {
		if seen[record.RecordID] {
			return fmt.Errorf("store: ingest %s: duplicate record ID %q", fieldID, record.RecordID)
		}
		seen[record.RecordID] = true
		statement, err := ingestStatement(doc, field, record)
		if err != nil {
			return err
		}
		script.WriteString(statement)
		script.WriteString("\n")
	}
	script.WriteString("COMMIT;\n")
	if err := execSQL(dbPath, script.String()); err != nil {
		return fmt.Errorf("store: ingest %s: %w", fieldID, err)
	}
	return nil
}

func ingestStatement(doc *cdl.EIRDoc, field *cdl.EIRField, record RecordInput) (string, error) {
	if record.RecordID == "" {
		return "", fmt.Errorf("store: ingest %s: empty record ID", field.ID)
	}
	if record.Version < 1 {
		return "", fmt.Errorf("store: ingest %s %s: version must be positive", field.ID, record.RecordID)
	}
	canonical, err := capabilities.CanonicalMarkdown(record.Payload)
	if err != nil {
		return "", err
	}
	values, err := extractValues(field.ID, field.Fields, canonical)
	if err != nil {
		return "", err
	}
	if err := certify.ValidateRecord(doc, field.ID, values); err != nil {
		return "", fmt.Errorf("%w (parsed declared fields: %s)", err, strings.Join(sortedKeys(values), ", "))
	}
	if err := certify.ValidateSemantics(field.ID, values); err != nil {
		return "", err
	}
	columns := []string{"la_record_id", "la_record_version", "la_semantic_fingerprint", "la_canonical_payload"}
	literals := []string{literal(record.RecordID), literal(int64(record.Version)), literal(capabilities.Hash(canonical)), literal(string(canonical))}
	for _, spec := range field.Fields {
		column := ident(spec.Name)
		value, present := values[spec.Name]
		if !present || value == nil {
			columns = append(columns, column)
			literals = append(literals, "NULL")
			continue
		}
		if items, ok := value.([]any); ok {
			canonicalSet, err := certify.CanonicalStringSet(items)
			if err != nil {
				return "", fmt.Errorf("store: ingest %s.%s: %w", field.ID, spec.Name, err)
			}
			value = strings.Join(canonicalSet, ", ")
		}
		columns = append(columns, column)
		literals = append(literals, literal(value))
	}
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = fmt.Sprintf("%q", column)
	}
	return fmt.Sprintf("INSERT INTO %q (%s) VALUES (%s);", ident(field.ID), strings.Join(quoted, ", "), strings.Join(literals, ", ")), nil
}

func findField(doc *cdl.EIRDoc, fieldID string) (*cdl.EIRField, error) {
	for i := range doc.Declarations.Fields {
		if doc.Declarations.Fields[i].ID == fieldID {
			return &doc.Declarations.Fields[i], nil
		}
	}
	return nil, fmt.Errorf("store: unknown FIELD %q", fieldID)
}

// extractValues maps payload field labels to declared specs. Every payload label
// must match exactly one declared column, and every declared column may appear at
// most once.
func extractValues(fieldID string, specs []cdl.EIRFieldSpec, payload []byte) (map[string]any, error) {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return nil, fmt.Errorf("store: parse %s payload: %w", fieldID, err)
	}
	byColumn := map[string]string{}
	for _, node := range parsed.Nodes {
		if node.Kind != md.Field {
			continue
		}
		column := capabilities.ColumnName(node.Label)
		if _, ok := byColumn[column]; ok {
			return nil, fmt.Errorf("store: %s payload repeats field %q", fieldID, node.Label)
		}
		byColumn[column] = node.Value
	}
	declared := map[string]bool{}
	values := map[string]any{}
	for _, spec := range specs {
		column := capabilities.ColumnName(spec.Name)
		declared[column] = true
		text, present := byColumn[column]
		if !present || strings.TrimSpace(text) == "" {
			continue
		}
		switch spec.Type {
		case "int":
			number, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("store: %s.%s expects an integer, got %q", fieldID, spec.Name, text)
			}
			values[spec.Name] = number
		case "set<string>":
			var items []any
			for _, part := range strings.Split(text, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					items = append(items, part)
				}
			}
			values[spec.Name] = items
		default:
			values[spec.Name] = strings.TrimSpace(text)
		}
	}
	for column := range byColumn {
		if !declared[column] {
			return nil, fmt.Errorf("store: %s payload has unknown field %q", fieldID, column)
		}
	}
	return values, nil
}

func execSQL(dbPath, sql string) error {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return fmt.Errorf("sqlite3 CLI not found on PATH: %w", err)
	}
	cmd := exec.Command("sqlite3", "--bail", dbPath)
	cmd.Stdin = strings.NewReader(sql + "\n")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func sortedKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return []string{"none"}
	}
	return keys
}
