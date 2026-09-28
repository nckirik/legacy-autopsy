package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestPreflightRegistryDrift checks every §3.1 record block field list against
// protocol.md.
func TestPreflightRegistryDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/preflight-registry.cdl")
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
	text, err := model.SectionText("3.1. `0A-PREFLIGHT.md`")
	if err != nil {
		t.Fatal(err)
	}
	blockField := map[string]string{
		"CLUSTER":   "ENTRY-CLUSTER",
		"DISCOVERY": "UNMAPPED-DISCOVERY",
		"DSC":       "DISCOVERY-BUFFER",
	}
	var introLabels []string
	labels := map[string][]string{}
	current := ""
	seenHeading := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### ") {
			seenHeading = true
			current = ""
			if strings.HasPrefix(trimmed, "### [") {
				id := strings.TrimPrefix(trimmed, "### [")
				if idx := strings.Index(id, "-ID]"); idx > 0 {
					current = id[:idx]
				}
			}
			continue
		}
		match := logFieldLine.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		if !seenHeading {
			introLabels = append(introLabels, normalizeLabel(match[1]))
			continue
		}
		if current == "" {
			continue
		}
		fieldID := blockField[current]
		labels[fieldID] = append(labels[fieldID], normalizeLabel(match[1]))
	}
	if len(introLabels) == 0 {
		t.Fatal("no preflight header fields found")
	}
	compareSets(t, "PREFLIGHT-HEADER", fields["PREFLIGHT-HEADER"], introLabels)
	if len(labels) != 3 {
		t.Fatalf("found %d blocks, want 3", len(labels))
	}
	for fieldID, want := range labels {
		compareSets(t, fieldID, fields[fieldID], want)
	}

	tableFields := map[string]string{"persona": "PERSONA-REGISTRY", "persona-owner": "REQUIRED-TRAVERSAL-MATRIX"}
	seen := 0
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cols := tableColumns(trimmed)
		if len(cols) == 0 {
			continue
		}
		if fieldID, ok := tableFields[cols[0]]; ok {
			compareSets(t, fieldID, fields[fieldID], cols)
			seen++
		}
	}
	if seen != 2 {
		t.Fatalf("found %d §3.1 tables, want 2", seen)
	}
}
