package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var classificationValues = regexp.MustCompile("`Classification` is `([^`]+)`")

// TestTraceabilityDrift checks the §12.1 row shapes, the classification enum, and
// the §12.5 evidence-binding table shape against protocol.md.
func TestTraceabilityDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/traceability.cdl")
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
	text, err := model.SectionText("12.1. Traceability — `20-TRACEABILITY.md`")
	if err != nil {
		t.Fatal(err)
	}
	columns := map[string][]string{}
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cols := tableColumns(trimmed)
		if len(cols) == 0 {
			continue
		}
		switch cols[0] {
		case "src":
			columns["TRACEABILITY-ROW"] = cols
		case "scope-id":
			columns["LIFECYCLE-ELIGIBILITY"] = cols
		}
	}
	if len(columns) != 2 {
		t.Fatalf("found %d §12.1 tables, want 2", len(columns))
	}
	for id, cols := range columns {
		compareSets(t, id, fields[id], cols)
	}
	if match := classificationValues.FindStringSubmatch(text); match != nil {
		var values []string
		for _, v := range strings.Split(match[1], "|") {
			values = append(values, strings.TrimSpace(v))
		}
		compareSets(t, "CLASSIFICATION", enums["CLASSIFICATION"], values)
	} else {
		t.Fatal("classification values not found")
	}

	text125, err := model.SectionText("12.5. Deterministic validation summary")
	if err != nil {
		t.Fatal(err)
	}
	var binding []string
	for _, line := range strings.Split(text125, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "| Row Kind |") {
			binding = tableColumns(trimmed)
		}
	}
	if len(binding) == 0 {
		t.Fatal("§12.5 evidence-binding table not found")
	}
	compareSets(t, "EVIDENCE-BINDING", fields["EVIDENCE-BINDING"], binding)
}

func tableColumns(line string) []string {
	var out []string
	for _, cell := range strings.Split(strings.Trim(line, "|"), "|") {
		out = append(out, normalizeLabel(cell))
	}
	return out
}
