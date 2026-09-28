package parity

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	boldToken     = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	authorityItem = regexp.MustCompile(`(?m)^\d+\. \*\*([^*]+):\*\*`)
	planeHeading  = regexp.MustCompile(`^## 1\.[1-4]\. Plane \d+ — (.+)$`)
)

// TestFoundationsDrift checks §0.1 normative keywords/authority classes/predicate
// kinds, §0.4 flow stages, and the §1.1–§1.4 plane enum against protocol.md.
func TestFoundationsDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/foundations.cdl")
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text01, err := model.SectionText("0.1. Normative language")
	if err != nil {
		t.Fatal(err)
	}
	keywords := boldToken.FindAllStringSubmatch(sentenceBefore(text01, "are normative"), -1)
	var keywordValues []string
	for _, match := range keywords {
		keywordValues = append(keywordValues, match[1])
	}
	compareSets(t, "NORMATIVE-KEYWORD", enums["NORMATIVE-KEYWORD"], keywordValues)

	var classes []string
	for _, match := range authorityItem.FindAllStringSubmatch(text01, -1) {
		classes = append(classes, match[1])
	}
	if len(classes) == 0 {
		t.Fatal("no authority classes parsed")
	}
	compareSets(t, "AUTHORITY-CLASS", enums["AUTHORITY-CLASS"], classes)

	var predicates []string
	for _, match := range boldToken.FindAllStringSubmatch(sentenceAfter(text01, "MUST be one of:"), -1) {
		predicates = append(predicates, match[1])
	}
	if len(predicates) == 0 {
		t.Fatal("no predicate kinds parsed")
	}
	compareSets(t, "PREDICATE-KIND", enums["PREDICATE-KIND"], predicates)

	text04, err := model.SectionText("0.4. Conceptual flow")
	if err != nil {
		t.Fatal(err)
	}
	var stages []string
	inFence := false
	for _, line := range strings.Split(text04, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if !inFence {
			continue
		}
		stages = append(stages, strings.TrimSpace(strings.TrimPrefix(trimmed, "->")))
	}
	if len(stages) == 0 {
		t.Fatal("no flow stages parsed")
	}
	compareSets(t, "FLOW-STAGE", enums["FLOW-STAGE"], stages)

	raw, err := os.ReadFile(filepath.Join(root(t), "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	var planes []string
	for _, line := range strings.Split(string(raw), "\n") {
		if match := planeHeading.FindStringSubmatch(line); match != nil {
			planes = append(planes, match[1])
		}
	}
	if len(planes) != 4 {
		t.Fatalf("parsed %d planes, want 4", len(planes))
	}
	compareSets(t, "PLANE", enums["PLANE"], planes)
}

func sentenceBefore(text, marker string) string {
	idx := strings.Index(text, marker)
	if idx < 0 {
		return text
	}
	return text[:idx]
}

func sentenceAfter(text, marker string) string {
	idx := strings.Index(text, marker)
	if idx < 0 {
		return ""
	}
	tail := text[idx:]
	if cut := strings.Index(tail, "\n\n"); cut >= 0 {
		tail = tail[:cut]
	}
	return tail
}
