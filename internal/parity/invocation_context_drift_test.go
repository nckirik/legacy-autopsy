package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	yamlHeaderKey = regexp.MustCompile(`^([A-Z][A-Za-z0-9 /-]+):`)
	logFieldLine  = regexp.MustCompile(`^- \*\*(.+?):\*\*`)
)

// TestInvocationContextDrift checks the §8.1 identity-header field list and the
// §8.8 invocation-log field list against protocol.md, and requires every
// hyphenated normative token in those sections to appear in the CDL text.
func TestInvocationContextDrift(t *testing.T) {
	res := compileSection(t, "protocol/invocation-context.cdl")
	fields := map[string][]string{}
	for _, f := range res.EIR.Declarations.Fields {
		names := make([]string, 0, len(f.Fields))
		for _, spec := range f.Fields {
			names = append(names, spec.Name)
		}
		fields[f.ID] = names
	}
	cdlText := ""
	for _, sec := range res.EIR.Sections {
		cdlText += sec.Goal + "\n"
	}
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += r.Goal + "\n"
	}
	for _, names := range fields {
		cdlText += strings.Join(names, "\n") + "\n"
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}

	// 1. §8.1 header keys == FIELD INVOCATION-HEADER names.
	text81, err := model.SectionText("8.1. Resume identity header")
	if err != nil {
		t.Fatal(err)
	}
	var wantHeader []string
	inYAML := false
	for _, line := range strings.Split(text81, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```yaml") {
			inYAML = true
			continue
		}
		if !inYAML {
			continue
		}
		if strings.HasPrefix(trimmed, "```") {
			break
		}
		if match := yamlHeaderKey.FindStringSubmatch(trimmed); match != nil {
			wantHeader = append(wantHeader, normalizeLabel(match[1]))
		}
	}
	if len(wantHeader) == 0 {
		t.Fatal("no header keys found in §8.1")
	}
	compareSets(t, "INVOCATION-HEADER", fields["INVOCATION-HEADER"], wantHeader)

	// 2. §8.8 log labels == FIELD INVOCATION-LOG names.
	text88, err := model.SectionText("8.8. `0G` invocation log")
	if err != nil {
		t.Fatal(err)
	}
	var wantLog []string
	for _, line := range strings.Split(text88, "\n") {
		if match := logFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			wantLog = append(wantLog, normalizeLabel(match[1]))
		}
	}
	if len(wantLog) == 0 {
		t.Fatal("no log fields found in §8.8")
	}
	compareSets(t, "INVOCATION-LOG", fields["INVOCATION-LOG"], wantLog)

	// 3. Normative token coverage for §8.1-8.2 and §8.6-8.8.
	for _, ref := range []string{
		"8.1. Resume identity header",
		"8.2. Strict single-scope modes",
		"8.6. Concurrent persona traversal",
		"8.7. Stale checkpoint guard",
		"8.8. `0G` invocation log",
	} {
		text, err := model.SectionText(ref)
		if err != nil {
			t.Fatal(err)
		}
		for _, token := range hyphenatedToken.FindAllString(text, -1) {
			if !strings.Contains(cdlText, token) {
				t.Fatalf("%s token %s is missing from the CDL section text", ref, token)
			}
		}
	}
}

// normalizeLabel lowercases a protocol label and joins alphanumeric runs with
// hyphens: "Cold-Resume Result / Check Fingerprint" -> "cold-resume-result-check-fingerprint".
func normalizeLabel(label string) string {
	var b strings.Builder
	lastHyphen := false
	for _, r := range strings.ToLower(label) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastHyphen = false
		default:
			if !lastHyphen && b.Len() > 0 {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}
	return strings.TrimRight(b.String(), "-")
}
