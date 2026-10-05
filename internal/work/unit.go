// Package work implements bounded semantic work units: declared dependencies,
// read-set provenance, stale-checked proposals, independent review, blocker
// classification, and per-unit instrumentation. It owns deterministic mechanics
// only; semantic content comes from an Executor.
package work

import (
	"fmt"

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
type Unit struct {
	ID             string   `json:"id"`
	Mode           string   `json:"mode"`
	Scope          string   `json:"scope,omitempty"`
	Reads          []Ref    `json:"reads,omitempty"`
	ResultType     string   `json:"result_type"`
	AllowedEffects []string `json:"allowed_effects,omitempty"`
	RequiresReview bool     `json:"requires_review,omitempty"`
}

// Validate checks the unit against the compiled schema: every declared read and
// the result type must be a table-class FIELD.
func (u Unit) Validate(doc *cdl.EIRDoc) error {
	if u.ID == "" || u.Mode == "" || u.ResultType == "" {
		return fmt.Errorf("work: unit requires ID, mode, and result type")
	}
	if err := requireTableField(doc, u.ResultType); err != nil {
		return fmt.Errorf("work: unit %s result: %w", u.ID, err)
	}
	for _, ref := range u.Reads {
		if err := requireTableField(doc, ref.FieldID); err != nil {
			return fmt.Errorf("work: unit %s dependency %s: %w", u.ID, ref, err)
		}
	}
	return nil
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
