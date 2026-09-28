package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestHandbookAndDecisionsDrift checks the §13.1 audience roles and the §14.5
// equivalence-test kinds against protocol.md.
func TestHandbookAndDecisionsDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/handbook-and-decisions.cdl")
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text131, err := model.SectionText("13.1. Audience paths and progressive disclosure")
	if err != nil {
		t.Fatal(err)
	}
	var roles []string
	for _, line := range strings.Split(text131, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		if idx := strings.Index(trimmed, ":"); idx > 2 {
			roles = append(roles, strings.TrimSpace(trimmed[2:idx]))
		}
	}
	if len(roles) == 0 {
		t.Fatal("no audience roles parsed")
	}
	compareSets(t, "AUDIENCE-ROLE", enums["AUDIENCE-ROLE"], roles)

	text145, err := model.SectionText("14.5. Equivalence and acceptance suite")
	if err != nil {
		t.Fatal(err)
	}
	var kinds []string
	for _, line := range strings.Split(text145, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		body := strings.TrimPrefix(trimmed, "- ")
		cut := len(body)
		for _, marker := range []string{" from ", " that ", " for "} {
			if idx := strings.Index(body, marker); idx > 0 && idx < cut {
				cut = idx
			}
		}
		kinds = append(kinds, strings.TrimSpace(body[:cut]))
	}
	if len(kinds) == 0 {
		t.Fatal("no equivalence kinds parsed")
	}
	compareSets(t, "EQUIVALENCE-TEST-KIND", enums["EQUIVALENCE-TEST-KIND"], kinds)
}
