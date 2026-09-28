package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	numberedCondition = regexp.MustCompile(`^\d+\. \*\*(.+?):\*\*`)
	iterationBacktick = regexp.MustCompile("`([A-Z-]+)`")
)

// TestCoverageAndExitsDrift checks the §10.3 condition list, the §10.1 sweep-record
// fields, and the §10.6 iteration tokens against protocol.md.
func TestCoverageAndExitsDrift(t *testing.T) {
	res := compileSection(t, "protocol/coverage-and-exits.cdl")
	enums := map[string][]string{}
	fields := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}
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

	// 1. §10.3 numbered conditions == EXIT-A-CONDITION.
	text103, err := model.SectionText("10.3. Composite Exit A - Deconstruction Closure")
	if err != nil {
		t.Fatal(err)
	}
	var wantConditions []string
	for _, line := range strings.Split(text103, "\n") {
		if match := numberedCondition.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			wantConditions = append(wantConditions, match[1])
		}
	}
	if len(wantConditions) == 0 {
		t.Fatal("no numbered conditions found in §10.3")
	}
	compareSets(t, "EXIT-A-CONDITION", enums["EXIT-A-CONDITION"], wantConditions)

	// 2. §10.1 sweep-record labels == FIELD SWEEP-RECORD.
	text101, err := model.SectionText("10.1. Evidence-backed `[R-SWEPT]`")
	if err != nil {
		t.Fatal(err)
	}
	var wantFields []string
	for _, line := range strings.Split(text101, "\n") {
		if match := logFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			wantFields = append(wantFields, normalizeLabel(match[1]))
		}
	}
	if len(wantFields) == 0 {
		t.Fatal("no sweep-record fields found in §10.1")
	}
	compareSets(t, "SWEEP-RECORD", fields["SWEEP-RECORD"], wantFields)

	// 3. §10.6 iteration tokens == the ID-ITERATION enum owned by typed-id.cdl.
	typed := compileSection(t, "protocol/typed-id.cdl")
	var iteration []string
	for _, e := range typed.EIR.Declarations.Enums {
		if e.ID == "ID-ITERATION" {
			iteration = e.Values
		}
	}
	text106, err := model.SectionText("10.6. Exit D - Limit exhaustion and deterministic iteration accounting")
	if err != nil {
		t.Fatal(err)
	}
	tokenLine := ""
	for _, line := range strings.Split(text106, "\n") {
		if strings.Contains(line, "Named iterations use exactly these canonical tokens") {
			tokenLine = line
			break
		}
	}
	if tokenLine == "" {
		t.Fatal("iteration token sentence not found in §10.6")
	}
	seen := map[string]bool{}
	var wantIteration []string
	for _, match := range iterationBacktick.FindAllStringSubmatch(tokenLine, -1) {
		if !seen[match[1]] {
			seen[match[1]] = true
			wantIteration = append(wantIteration, match[1])
		}
	}
	if len(wantIteration) != 26 {
		t.Fatalf("parsed %d iteration tokens from §10.6, want 26", len(wantIteration))
	}
	compareSets(t, "ID-ITERATION", iteration, wantIteration)
}
