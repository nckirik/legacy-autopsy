package parity

import (
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestExtractionEvolutionDrift checks the §6.1 atomic-component schema, the §6.2
// shared-reference schema, and the block-confidence-rank enum against protocol.md.
func TestExtractionEvolutionDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/extraction-evolution.cdl")
	fields := map[string][]string{}
	for _, f := range res.EIR.Declarations.Fields {
		names := make([]string, 0, len(f.Fields))
		for _, spec := range f.Fields {
			names = append(names, spec.Name)
		}
		fields[f.ID] = names
	}
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text61, err := model.SectionText("6.1. Atomic component schema")
	if err != nil {
		t.Fatal(err)
	}
	compareSets(t, "ATOMIC-COMPONENT", fields["ATOMIC-COMPONENT"], sectionLabels(text61, false))
	compareSets(t, "BLOCK-CONFIDENCE-RANK", enums["BLOCK-CONFIDENCE-RANK"], enumValuesFromFieldLine(text61, "block-confidence-rank"))

	text62, err := model.SectionText("6.2. Dual-entry discovery")
	if err != nil {
		t.Fatal(err)
	}
	compareSets(t, "SHARED-REFERENCE", fields["SHARED-REFERENCE"], sectionLabels(text62, false))
}
