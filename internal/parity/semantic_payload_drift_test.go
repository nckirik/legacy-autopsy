package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var hyphenatedToken = regexp.MustCompile(`\b[A-Z][A-Z0-9]+(?:-[A-Z0-9]+)+\b`)

// TestSemanticPayloadTokenDrift requires every hyphenated normative token in
// protocol.md §4.1.1 to appear in the CDL section text. Adding or renaming a field
// in protocol.md therefore fails until the CDL section is updated.
func TestSemanticPayloadTokenDrift(t *testing.T) {
	res := compileSection(t, "protocol/semantic-payload-identity.cdl")
	cdlText := res.EIR.Sections[0].Goal
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += "\n" + r.Goal
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("4.1.1. Semantic payload identity and certification envelopes")
	if err != nil {
		t.Fatal(err)
	}
	tokens := map[string]bool{}
	for _, token := range hyphenatedToken.FindAllString(text, -1) {
		tokens[token] = true
	}
	if len(tokens) < 8 {
		t.Fatalf("expected at least 8 normative tokens in §4.1.1, found %d", len(tokens))
	}
	for token := range tokens {
		if !strings.Contains(cdlText, token) {
			t.Fatalf("§4.1.1 token %s is missing from the CDL section text", token)
		}
	}
}
