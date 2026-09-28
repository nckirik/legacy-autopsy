package parity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestMigrationAndConformanceDrift checks the §19.2 conformance-case enum and the
// Part 20 artifact-ownership matrix shape against protocol.md.
func TestMigrationAndConformanceDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/migration-and-conformance.cdl")
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
	text192, err := model.SectionText("19.2. Normative conformance cases")
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
	compareSets(t, "CONFORMANCE-CASE", enums["CONFORMANCE-CASE"], cases)

	raw, err := os.ReadFile(filepath.Join(root(t), "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	inPart := false
	var matrix []string
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# Part 20.") {
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
		t.Fatal("Part 20 ownership matrix not found")
	}
	for i, col := range matrix {
		matrix[i] = strings.ReplaceAll(col, "authorized-writer-s", "authorized-writers")
	}
	compareSets(t, "ARTIFACT-OWNERSHIP-ROW", fields["ARTIFACT-OWNERSHIP-ROW"], matrix)
}
