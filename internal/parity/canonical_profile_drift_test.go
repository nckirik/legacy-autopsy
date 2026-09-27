package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	backtickedUpper = regexp.MustCompile("`([A-Z][A-Z0-9-]+)`")
	domainPrefix    = regexp.MustCompile("`([A-Z-]+)\\|`")
)

// TestCanonicalProfileDrift checks the §4.1.2 declarations against protocol.md:
// every hyphenated normative token must appear in the CDL text, and the binding,
// envelope, artifact-type, row-kind, and domain-prefix registries must match the
// protocol tables and literals exactly.
func TestCanonicalProfileDrift(t *testing.T) {
	res := compileSection(t, "examples/spec/canonical-hash-profile.cdl")
	cdlText := res.EIR.Sections[0].Goal
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += "\n" + r.Goal
	}
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("4.1.2. Normative canonical hash profile and packaging artifacts")
	if err != nil {
		t.Fatal(err)
	}

	// 1. Normative token coverage.
	tokens := map[string]bool{}
	for _, token := range hyphenatedToken.FindAllString(text, -1) {
		tokens[token] = true
	}
	if len(tokens) < 20 {
		t.Fatalf("expected at least 20 normative tokens in §4.1.2, found %d", len(tokens))
	}
	for token := range tokens {
		if !strings.Contains(cdlText, token) {
			t.Fatalf("§4.1.2 token %s is missing from the CDL section text", token)
		}
	}

	// 2. Binding-kind table == ENUM EVIDENCE-BINDING-KIND.
	var bindingKinds []string
	inTable := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "| Binding Kind") {
			inTable = true
			continue
		}
		if !inTable {
			continue
		}
		if !strings.HasPrefix(trimmed, "|") {
			break
		}
		cell := strings.TrimSpace(strings.Split(trimmed, "|")[1])
		if cell == "" || strings.HasPrefix(cell, ":") || strings.Contains(cell, "---") {
			continue
		}
		bindingKinds = append(bindingKinds, cell)
	}
	compareSets(t, "EVIDENCE-BINDING-KIND", enums["EVIDENCE-BINDING-KIND"], bindingKinds)

	// 3. Envelope kinds from the parenthesized literal list.
	envelopeLine := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "envelope kind") {
			envelopeLine = line
			break
		}
	}
	if envelopeLine == "" {
		t.Fatal("envelope kind literal not found in §4.1.2")
	}
	var envelopeKinds []string
	for _, match := range backtickedUpper.FindAllStringSubmatch(envelopeLine, -1) {
		envelopeKinds = append(envelopeKinds, match[1])
	}
	compareSets(t, "ENVELOPE-KIND", enums["ENVELOPE-KIND"], envelopeKinds)

	// 4. Artifact types from the ARTIFACT-TYPE field literal.
	artifactLine := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "ARTIFACT-TYPE:**") {
			artifactLine = line
			break
		}
	}
	if artifactLine == "" {
		t.Fatal("ARTIFACT-TYPE literal not found in §4.1.2")
	}
	var artifactTypes []string
	for _, token := range hyphenatedToken.FindAllString(artifactLine, -1) {
		if token != "ARTIFACT-TYPE" {
			artifactTypes = append(artifactTypes, token)
		}
	}
	compareSets(t, "ARTIFACT-TYPE", enums["ARTIFACT-TYPE"], artifactTypes)

	// 5. Row kinds.
	compareSets(t, "BINDING-ROW-KIND", enums["BINDING-ROW-KIND"], []string{"CHECK", "BLOCKER"})
	for _, kind := range []string{"CHECK", "BLOCKER"} {
		if !strings.Contains(text, "`"+kind+"`") {
			t.Fatalf("row kind %s not found in §4.1.2", kind)
		}
	}

	// 6. Domain prefixes.
	var prefixes []string
	for _, match := range domainPrefix.FindAllStringSubmatch(text, -1) {
		prefixes = append(prefixes, match[1])
	}
	compareSets(t, "HASH-DOMAIN-PREFIX", enums["HASH-DOMAIN-PREFIX"], prefixes)
}
