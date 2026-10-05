// Package certify validates protocol records against their compiled EIR schema
// before any fingerprint, coverage arithmetic, or packaging binding is computed.
package certify

import (
	"fmt"
	"sort"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// ValidateRecord checks that values match the declared FIELD block: no unknown
// fields, all required fields present, declared types respected, and enum-typed
// values inside their closed domain. The first violation is reported.
func ValidateRecord(doc *cdl.EIRDoc, fieldID string, values map[string]any) error {
	var field *cdl.EIRField
	for i := range doc.Declarations.Fields {
		if doc.Declarations.Fields[i].ID == fieldID {
			field = &doc.Declarations.Fields[i]
			break
		}
	}
	if field == nil {
		return fmt.Errorf("certify: unknown record schema %q", fieldID)
	}
	enums := map[string][]string{}
	for _, e := range doc.Declarations.Enums {
		enums[e.ID] = e.Values
	}
	declared := map[string]cdl.EIRFieldSpec{}
	for _, spec := range field.Fields {
		declared[spec.Name] = spec
	}
	var unknown []string
	for name := range values {
		if _, ok := declared[name]; !ok {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return fmt.Errorf("certify: %s has unknown fields %v", fieldID, unknown)
	}
	for _, spec := range field.Fields {
		value, present := values[spec.Name]
		if !present || value == nil {
			if spec.Required {
				return fmt.Errorf("certify: %s is missing required field %q", fieldID, spec.Name)
			}
			continue
		}
		if err := checkValue(spec, value, enums); err != nil {
			return fmt.Errorf("certify: %s.%s: %w", fieldID, spec.Name, err)
		}
	}
	return nil
}

func checkValue(spec cdl.EIRFieldSpec, value any, enums map[string][]string) error {
	if values, ok := enums[spec.Type]; ok {
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("enum value must be a string, got %T", value)
		}
		for _, allowed := range values {
			if text == allowed {
				return nil
			}
		}
		return fmt.Errorf("value %q is not in enum %s", text, spec.Type)
	}
	switch spec.Type {
	case "int":
		switch v := value.(type) {
		case int:
			return nil
		case int64:
			return nil
		case float64:
			if v != float64(int64(v)) {
				return fmt.Errorf("integer required, got %v", v)
			}
			return nil
		default:
			return fmt.Errorf("integer required, got %T", value)
		}
	case "set<string>":
		items, ok := value.([]any)
		if !ok {
			if _, ok := value.([]string); ok {
				return nil
			}
			return fmt.Errorf("string set required, got %T", value)
		}
		for _, item := range items {
			if _, ok := item.(string); !ok {
				return fmt.Errorf("string set contains non-string %T", item)
			}
		}
		return nil
	default:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("string required, got %T", value)
		}
		return nil
	}
}
