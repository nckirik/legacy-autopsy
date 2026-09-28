package parity

import (
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestHumanHatchProbeDrift checks the §9.4 probe specification, the §9.4 safety
// classification enum, and the §9.6 no-mock fallback payload against protocol.md.
func TestHumanHatchProbeDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/ticket-fsm.cdl")
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
	text94, err := model.SectionText("9.4. Probe specification and asynchronous handoff")
	if err != nil {
		t.Fatal(err)
	}
	compareSets(t, "PROBE-SPECIFICATION", fields["PROBE-SPECIFICATION"], sectionLabels(text94, false))
	compareSets(t, "SAFETY-CLASSIFICATION", enums["SAFETY-CLASSIFICATION"], enumValuesFromFieldLine(text94, "safety-classification"))

	text96, err := model.SectionText("9.6. No-mock fallback payload")
	if err != nil {
		t.Fatal(err)
	}
	compareSets(t, "NO-MOCK-FALLBACK", fields["NO-MOCK-FALLBACK"], sectionLabels(text96, false))
}
