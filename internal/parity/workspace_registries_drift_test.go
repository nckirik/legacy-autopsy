package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestWorkspaceRegistriesDrift checks every §3.2 record block field list against
// protocol.md.
func TestWorkspaceRegistriesDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/workspace-registries.cdl")
	fields := map[string][]string{}
	for _, f := range res.EIR.Declarations.Fields {
		names := make([]string, 0, len(f.Fields))
		for _, spec := range f.Fields {
			names = append(names, spec.Name)
		}
		fields[f.ID] = names
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("3.2. Foundational registries")
	if err != nil {
		t.Fatal(err)
	}
	blockField := map[string]string{
		"SHR":   "SHARED-ENTRY",
		"AUTH":  "AUTH-MECHANISM",
		"PRV":   "PRIVACY-FINDING",
		"CAP":   "CAPABILITY-RECORD",
		"CMP":   "ROUTINE-CROSS-REFERENCE",
		"INV":   "INVARIANT",
		"STATE": "SHARED-STATE",
	}
	labels := map[string][]string{}
	current := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### [") {
			id := strings.TrimPrefix(trimmed, "### [")
			if idx := strings.Index(id, "-ID]"); idx > 0 {
				current = id[:idx]
			}
			continue
		}
		if current == "" {
			continue
		}
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			fieldID := blockField[current]
			labels[fieldID] = append(labels[fieldID], normalizeLabel(match[1]))
		}
	}
	if len(labels) != 7 {
		t.Fatalf("found %d record blocks, want 7", len(labels))
	}
	for fieldID, want := range labels {
		compareSets(t, fieldID, fields[fieldID], want)
	}
}
