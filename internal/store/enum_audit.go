package store

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// AuditEnumRegistries checks the store-facing enum registry for completeness:
// every declared enum carries at least one unique value, and every enum-typed
// column in the generated DDL has a CHECK constraint that lists exactly the
// declared values — no missing and no extra value. Diagnostics are
// deterministic and sorted; an empty result means the registry is complete.
func AuditEnumRegistries(doc *cdl.EIRDoc) ([]string, error) {
	var diagnostics []string
	for _, enum := range doc.Declarations.Enums {
		if len(enum.Values) == 0 {
			diagnostics = append(diagnostics, fmt.Sprintf("enum %s declares no values", enum.ID))
		}
		seen := map[string]bool{}
		for _, value := range enum.Values {
			if seen[value] {
				diagnostics = append(diagnostics, fmt.Sprintf("enum %s repeats value %q", enum.ID, value))
			}
			seen[value] = true
		}
	}
	ddl, err := Schema(doc, Provenance{Generator: "store/enum-audit", SourceFingerprint: doc.Envelope.SourceFingerprint})
	if err != nil {
		return nil, err
	}
	text := string(ddl)
	for _, field := range doc.Declarations.Fields {
		switch field.Storage {
		case "record", "registry", "singleton", "log":
		default:
			continue
		}
		table := ident(field.ID)
		block := tableBlock(text, table)
		if block == "" {
			diagnostics = append(diagnostics, fmt.Sprintf("field %s has no generated table block", field.ID))
			continue
		}
		for _, spec := range field.Fields {
			values, ok := declaredEnum(doc, spec.Type)
			if !ok {
				continue
			}
			column := ident(spec.Name)
			got, found := enumCheckValues(block, column)
			if !found {
				diagnostics = append(diagnostics, fmt.Sprintf("field %s column %s is typed %s but the generated DDL has no enum CHECK", field.ID, spec.Name, spec.Type))
				continue
			}
			want := map[string]bool{}
			for _, value := range values {
				want[value] = true
			}
			if len(values) == 0 {
				diagnostics = append(diagnostics, fmt.Sprintf("field %s column %s is typed %s which declares no values", field.ID, spec.Name, spec.Type))
			}
			have := map[string]bool{}
			for _, value := range got {
				have[value] = true
				if !want[value] {
					diagnostics = append(diagnostics, fmt.Sprintf("field %s column %s CHECK lists undeclared value %q", field.ID, spec.Name, value))
				}
			}
			for _, value := range values {
				if !have[value] {
					diagnostics = append(diagnostics, fmt.Sprintf("field %s column %s CHECK is missing declared value %q", field.ID, spec.Name, value))
				}
			}
		}
	}
	sort.Strings(diagnostics)
	return diagnostics, nil
}

func declaredEnum(doc *cdl.EIRDoc, id string) ([]string, bool) {
	for _, enum := range doc.Declarations.Enums {
		if enum.ID == id {
			return enum.Values, true
		}
	}
	return nil, false
}

var enumCheckPattern = regexp.MustCompile(`(?m)CHECK \("([^"]+)" IN \(([^)]*)\)\)`)

// tableBlock returns the CREATE TABLE block for one generated table, so enum
// CHECK validation is context-qualified: the same column name in two contexts
// (for example target_decision on different tables) never cross-checks.
func tableBlock(ddl, table string) string {
	marker := "CREATE TABLE \"" + table + "\" ("
	start := strings.Index(ddl, marker)
	if start < 0 {
		return ""
	}
	end := strings.Index(ddl[start:], ");")
	if end < 0 {
		return ""
	}
	return ddl[start : start+end]
}

// enumCheckValues extracts the CHECK value list for one column from generated
// DDL. Values are unquoted and SQL-unescaped.
func enumCheckValues(ddl, column string) ([]string, bool) {
	for _, match := range enumCheckPattern.FindAllStringSubmatch(ddl, -1) {
		if match[1] != column {
			continue
		}
		var values []string
		for _, cell := range strings.Split(match[2], ",") {
			cell = strings.TrimSpace(cell)
			if !strings.HasPrefix(cell, "'") || !strings.HasSuffix(cell, "'") {
				continue
			}
			values = append(values, strings.ReplaceAll(cell[1:len(cell)-1], "''", "'"))
		}
		return values, true
	}
	return nil, false
}
