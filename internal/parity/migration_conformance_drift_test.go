package parity

import (
	"os"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestMigrationAndConformanceDrift checks the §17.2 conformance-case enum and the
// Part 18 artifact-ownership matrix shape against protocol.md.
func TestMigrationAndConformanceDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/conformance.cdl")
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
	text192, err := model.SectionText("17.2. Normative conformance cases")
	if err != nil {
		t.Fatal(err)
	}
	var cases []string
	for _, line := range strings.Split(text192, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- **") {
			continue
		}
		rest := strings.TrimPrefix(trimmed, "- **")
		if idx := strings.Index(rest, ":**"); idx > 0 {
			cases = append(cases, rest[:idx])
		}
	}
	if len(cases) == 0 {
		t.Fatal("no conformance cases parsed")
	}
	// The frozen oracle's case set is a lower bound: a revision may add cases but
	// must never silently remove or rename one.
	current := map[string]bool{}
	for _, value := range enums["CONFORMANCE-CASE"] {
		current[value] = true
	}
	for _, c := range cases {
		if !current[c] {
			t.Errorf("legacy CONFORMANCE-CASE %q was removed or renamed", c)
		}
	}

	raw, err := os.ReadFile(protocol.OraclePath(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	inPart := false
	var matrix []string
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# Part 18.") {
			inPart = true
			continue
		}
		if !inPart || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cols := tableColumns(trimmed)
		if len(cols) == 5 && cols[0] == "artifact-record" {
			matrix = cols
		}
	}
	if len(matrix) == 0 {
		t.Fatal("Part 18 ownership matrix not found")
	}
	for i, col := range matrix {
		matrix[i] = strings.ReplaceAll(col, "authorized-writer-s", "authorized-writers")
	}
	compareSets(t, "ARTIFACT-OWNERSHIP-ROW", fields["ARTIFACT-OWNERSHIP-ROW"], matrix)
}
