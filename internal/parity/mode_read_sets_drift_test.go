package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestModeReadSetsDrift checks the §8.3 exhaustive invocation-mode enum against the
// compiled registry and the mandatory §8.4 read-set rules.
func TestModeReadSetsDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/invocation-context.cdl")
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("8.3. Invocation-mode enum and explicit multi-file modes")
	if err != nil {
		t.Fatal(err)
	}
	var values []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if !inFence {
			continue
		}
		for _, value := range strings.Split(trimmed, "|") {
			if value = strings.TrimSpace(value); value != "" {
				values = append(values, value)
			}
		}
	}
	if len(values) != 13 {
		t.Fatalf("parsed %d invocation modes, want 13", len(values))
	}
	compareSets(t, "REGISTRY-INVOCATIONHEADER-INVOCATION-MODE", enums["REGISTRY-INVOCATIONHEADER-INVOCATION-MODE"], values)
}
