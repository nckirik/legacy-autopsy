package parity

import (
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/spec"
)

// TestProseAndStructureAreJointlyPresent enforces the joint-normativity
// contract: every section carries normative prose, and every section-scoped
// structured declaration belongs to a section that carries prose. A structure-
// only or prose-only section is a source defect.
func TestProseAndStructureAreJointlyPresent(t *testing.T) {
	res, err := spec.Compile(root(t))
	if err != nil {
		t.Fatal(err)
	}
	hasGoal := map[string]bool{}
	for _, section := range res.EIR.Sections {
		if section.Goal == "" {
			t.Fatalf("section %s (%s) has structured content but no normative prose", section.Number, section.ID)
		}
		hasGoal[section.ID] = true
	}
	for _, rule := range res.EIR.Declarations.Rules {
		if rule.Section == "" {
			continue
		}
		if !hasGoal[rule.Section] {
			t.Fatalf("rule %s belongs to section %q without prose", rule.ID, rule.Section)
		}
	}
	for _, enum := range res.EIR.Declarations.Enums {
		if enum.Section != "" && !hasGoal[enum.Section] {
			t.Fatalf("enum %s belongs to section %q without prose", enum.ID, enum.Section)
		}
	}
	for _, field := range res.EIR.Declarations.Fields {
		if field.Section != "" && !hasGoal[field.Section] {
			t.Fatalf("field %s belongs to section %q without prose", field.ID, field.Section)
		}
	}
	for _, state := range res.EIR.Declarations.States {
		if state.Section != "" && !hasGoal[state.Section] {
			t.Fatalf("state %s belongs to section %q without prose", state.ID, state.Section)
		}
	}
	for _, table := range res.EIR.Declarations.Tables {
		if table.Section != "" && !hasGoal[table.Section] {
			t.Fatalf("table %s belongs to section %q without prose", table.ID, table.Section)
		}
	}
}
