package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var backtickedValue = regexp.MustCompile("`([^`]+)`")

// TestPersonasPovsDrift checks the §2.5 capability-state block, runtime-state matrix,
// and the §2.6 target-decision enum against protocol.md.
func TestPersonasPovsDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/personas-povs.cdl")
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
	text25, err := model.SectionText("2.5. Capability and runtime state")
	if err != nil {
		t.Fatal(err)
	}
	var labels []string
	var matrix []string
	for _, line := range strings.Split(text25, "\n") {
		trimmed := strings.TrimSpace(line)
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			labels = append(labels, normalizeLabel(match[1]))
			continue
		}
		if strings.HasPrefix(trimmed, "|") {
			cols := tableColumns(trimmed)
			if len(cols) > 0 && cols[0] == "environment" {
				matrix = cols
			}
		}
	}
	if len(labels) == 0 || len(matrix) == 0 {
		t.Fatalf("§2.5 parsed %d labels and %d matrix columns", len(labels), len(matrix))
	}
	compareSets(t, "CAPABILITY-STATE-BLOCK", fields["CAPABILITY-STATE-BLOCK"], labels)
	compareSets(t, "RUNTIME-STATE-MATRIX", fields["RUNTIME-STATE-MATRIX"], matrix)

	text26, err := model.SectionText("2.6. Capability grouping and target decision")
	if err != nil {
		t.Fatal(err)
	}
	idx := strings.Index(text26, "target decision: ")
	if idx < 0 {
		t.Fatal("target decision values not found")
	}
	tail := text26[idx:]
	if cut := strings.IndexAny(tail, ".;"); cut >= 0 {
		tail = tail[:cut]
	}
	var values []string
	for _, match := range backtickedValue.FindAllStringSubmatch(tail, -1) {
		values = append(values, match[1])
	}
	if len(values) == 0 {
		t.Fatal("no target decision values parsed")
	}
	compareSets(t, "TARGET-DECISION", enums["TARGET-DECISION"], values)
}
