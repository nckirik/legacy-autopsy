package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestAcquisitionTrustDrift checks the §7.2 export register fields and the §7.4
// per-platform inventory kinds against protocol.md.
func TestAcquisitionTrustDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/acquisition-trust.cdl")
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
	text72, err := model.SectionText("7.2. Acquisition register")
	if err != nil {
		t.Fatal(err)
	}
	var labels []string
	for _, line := range strings.Split(text72, "\n") {
		if match := logFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			labels = append(labels, normalizeLabel(match[1]))
		}
	}
	if len(labels) == 0 {
		t.Fatal("no export register fields found in §7.2")
	}
	compareSets(t, "EXPORT", fields["EXPORT"], labels)

	text74, err := model.SectionText("7.4. Required explosion coverage")
	if err != nil {
		t.Fatal(err)
	}
	compareSets(t, "N8N-INVENTORY-KIND", enums["N8N-INVENTORY-KIND"], inventoryKinds(t, text74, "For n8n, inventory "))
	compareSets(t, "APPSMITH-INVENTORY-KIND", enums["APPSMITH-INVENTORY-KIND"], inventoryKinds(t, text74, "For Appsmith, inventory "))
}

func inventoryKinds(t *testing.T, text, prefix string) []string {
	t.Helper()
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		body := strings.TrimPrefix(line, prefix)
		if idx := strings.Index(body, ". "); idx >= 0 {
			body = body[:idx]
		}
		body = strings.TrimSuffix(body, ".")
		body = strings.ReplaceAll(body, ", and ", ", ")
		return strings.Split(body, ", ")
	}
	t.Fatalf("inventory sentence %q not found", prefix)
	return nil
}
