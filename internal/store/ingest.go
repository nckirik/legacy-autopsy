package store

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/certify"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// Ingest validates a Markdown record payload against its declared FIELD schema,
// derives the queryable columns, and inserts one row. The canonical payload bytes
// are stored verbatim in la_canonical_payload, so record payload digests stay
// binding-independent; the queryable columns are derived projections and can be
// rebuilt from the payload. Unknown, duplicate, missing-required, wrongly typed,
// or noncanonical set values fail closed, as does inserting into a non-table
// storage class or a duplicate record ID.
func Ingest(dbPath string, doc *cdl.EIRDoc, fieldID, recordID string, version int, payload []byte) error {
	if recordID == "" {
		return fmt.Errorf("store: ingest %s: empty record ID", fieldID)
	}
	if version < 1 {
		return fmt.Errorf("store: ingest %s %s: version must be positive", fieldID, recordID)
	}
	field, err := tableField(doc, fieldID)
	if err != nil {
		return err
	}
	canonical, err := capabilities.CanonicalMarkdown(payload)
	if err != nil {
		return err
	}
	values, err := extractValues(fieldID, field.Fields, canonical)
	if err != nil {
		return err
	}
	if err := certify.ValidateRecord(doc, fieldID, values); err != nil {
		return err
	}
	columns := []string{"la_record_id", "la_record_version", "la_semantic_fingerprint", "la_canonical_payload"}
	literals := []string{literal(recordID), literal(int64(version)), literal(capabilities.Hash(canonical)), literal(string(canonical))}
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
				return fmt.Errorf("store: ingest %s.%s: %w", fieldID, spec.Name, err)
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
	sql := fmt.Sprintf("INSERT INTO %q (%s) VALUES (%s);", ident(fieldID), strings.Join(quoted, ", "), strings.Join(literals, ", "))
	if err := execSQL(dbPath, sql); err != nil {
		return fmt.Errorf("store: ingest %s %s: %w", fieldID, recordID, err)
	}
	return nil
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
