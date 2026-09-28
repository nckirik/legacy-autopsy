package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var acquisitionStateToken = regexp.MustCompile(`\[(A-[A-Z-]+)\]`)

// TestAcquisitionStateDrift checks the §7.5 acquisition state domain against the
// bracket state tokens in protocol.md, and requires normative token coverage.
func TestAcquisitionStateDrift(t *testing.T) {
	res := compileSection(t, "protocol/export-acquisition-loop.cdl")
	var states []string
	cdlText := res.EIR.Sections[0].Goal
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += "\n" + r.Goal
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "ACQUISITION-STATE" {
			states = e.Values
		}
		cdlText += "\n" + strings.Join(e.Values, "\n")
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("7.5. n8n human acquisition loop")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	var wantStates []string
	for _, match := range acquisitionStateToken.FindAllStringSubmatch(text, -1) {
		if !seen[match[1]] {
			seen[match[1]] = true
			wantStates = append(wantStates, match[1])
		}
	}
	if len(wantStates) == 0 {
		t.Fatal("no acquisition state tokens found in §7.5")
	}
	compareSets(t, "ACQUISITION-STATE", states, wantStates)

	for _, token := range hyphenatedToken.FindAllString(text, -1) {
		if !strings.Contains(cdlText, token) {
			t.Fatalf("§7.5 token %s is missing from the CDL section text", token)
		}
	}
}
