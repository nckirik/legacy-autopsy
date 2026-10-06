// Package work implements bounded semantic work units: declared dependencies,
// read-set provenance, stale-checked proposals, independent review, blocker
// classification, and per-unit instrumentation. It owns deterministic mechanics
// only; semantic content comes from an Executor.
package work

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// Ref identifies one stored record dependency.
type Ref struct {
	FieldID  string `json:"field_id"`
	RecordID string `json:"record_id"`
}

func (r Ref) String() string { return r.FieldID + "/" + r.RecordID }

// Unit is one bounded semantic work item with an explicit read closure and a
// single result contract. It never carries conversational history.
// Input is a bounded runtime-supplied context block (for example a deterministic
// repository scan). Inputs are immutable for the unit and fingerprinted into the
// context; they are not store records and carry no stale check.
type Input struct {
	Name  string `json:"name"`
	Bytes []byte `json:"bytes,omitempty"`
}

type Unit struct {
	ID                   string `json:"id"`
	Mode                 string `json:"mode"`
	Scope                string `json:"scope,omitempty"`
	Persona              string `json:"persona,omitempty"`
	POV                  string `json:"pov,omitempty"`
	RecordID             string `json:"record_id,omitempty"`
	IDFromField          string `json:"id_from_field,omitempty"`
	IDFromCoordinate     string `json:"id_from_coordinate,omitempty"`
	IDOwnerField         string `json:"id_owner_field,omitempty"`
	IDDiscriminatorField string `json:"id_discriminator_field,omitempty"`
	// IDFromFields composes a natural composite record ID from the named
	// payload fields joined by "|" (for keyed supporting rows).
	IDFromFields    []string `json:"id_from_fields,omitempty"`
	IDType          string   `json:"id_type,omitempty"`
	Namespace       string   `json:"namespace,omitempty"`
	ExpectedRecords int      `json:"expected_records,omitempty"`
	// RuntimeFields are payload fields the runtime fills with the resolved
	// record ID after identity assignment (for example a declared `block-id`).
	RuntimeFields   []string `json:"runtime_fields,omitempty"`
	Reads           []Ref    `json:"reads,omitempty"`
	ResultType      string   `json:"result_type"`
	AllowedEffects  []string `json:"allowed_effects,omitempty"`
	RequiresReview  bool     `json:"requires_review,omitempty"`
	Inputs          []Input  `json:"inputs,omitempty"`
	Sections        []string `json:"sections,omitempty"`
	ReasoningEffort string   `json:"reasoning_effort,omitempty"`
}

// Validate checks the unit against the compiled schema: every declared read and
// the result type must be a table-class FIELD.
func (u Unit) Validate(doc *cdl.EIRDoc) error {
	if u.ID == "" || u.Mode == "" || u.ResultType == "" {
		return fmt.Errorf("work: unit requires ID, mode, and result type")
	}
	modes := 0
	for _, set := range []bool{u.RecordID != "", u.IDFromField != "", u.IDFromCoordinate != "", u.IDOwnerField != "", len(u.IDFromFields) > 0} {
		if set {
			modes++
		}
	}
	if modes != 1 {
		return fmt.Errorf("work: unit %s requires exactly one record-ID mode", u.ID)
	}
	if (u.IDFromCoordinate != "" || u.IDOwnerField != "") && (u.IDType == "" || u.Namespace == "") {
		return fmt.Errorf("work: unit %s runtime IDs require an ID type and namespace", u.ID)
	}
	if u.IDOwnerField != "" && u.IDDiscriminatorField == "" {
		return fmt.Errorf("work: unit %s typed-field IDs require an owner and discriminator field", u.ID)
	}
	if err := requireTableField(doc, u.ResultType); err != nil {
		return fmt.Errorf("work: unit %s result: %w", u.ID, err)
	}
	for _, ref := range u.Sections {
		if _, err := SectionText(doc, ref); err != nil {
			return fmt.Errorf("work: unit %s section %s: %w", u.ID, ref, err)
		}
	}
	for _, ref := range u.Reads {
		if err := requireTableField(doc, ref.FieldID); err != nil {
			return fmt.Errorf("work: unit %s dependency %s: %w", u.ID, ref, err)
		}
	}
	return nil
}

// SectionText resolves a normative section by ID or derived number.
func SectionText(doc *cdl.EIRDoc, reference string) (string, error) {
	for _, section := range doc.Sections {
		if section.ID == reference || section.Number == reference {
			return strings.TrimSpace(section.Goal), nil
		}
	}
	return "", fmt.Errorf("unknown section %q", reference)
}

// FieldContract returns the declared field names in order, the required subset,
// and the closed option domain for every enum-typed field, so executors can be
// given the exact payload contract without guessing.
func FieldContract(doc *cdl.EIRDoc, fieldID string) (names []string, required []string, options map[string][]string, err error) {
	enums := map[string][]string{}
	for _, enum := range doc.Declarations.Enums {
		enums[enum.ID] = enum.Values
	}
	for i := range doc.Declarations.Fields {
		field := &doc.Declarations.Fields[i]
		if field.ID != fieldID {
			continue
		}
		for _, spec := range field.Fields {
			names = append(names, spec.Name)
			if spec.Required {
				required = append(required, spec.Name)
			}
			if values, ok := enums[spec.Type]; ok {
				if options == nil {
					options = map[string][]string{}
				}
				options[spec.Name] = values
			}
		}
		return names, required, options, nil
	}
	return nil, nil, nil, fmt.Errorf("unknown FIELD %s", fieldID)
}

func requireTableField(doc *cdl.EIRDoc, fieldID string) error {
	for i := range doc.Declarations.Fields {
		field := &doc.Declarations.Fields[i]
		if field.ID != fieldID {
			continue
		}
		switch field.Storage {
		case "record", "registry", "singleton", "log":
			return nil
		default:
			return fmt.Errorf("FIELD %s has STORAGE %s, not a table class", fieldID, field.Storage)
		}
	}
	return fmt.Errorf("unknown FIELD %s", fieldID)
}
