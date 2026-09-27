package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var coldResumeFieldLine = regexp.MustCompile(`^- \*\*(.+?):\*\*`)

// TestColdResumeFieldDrift checks the §8.5 field list against the block in
// protocol.md and requires the carrier to be the declared capability constant.
func TestColdResumeFieldDrift(t *testing.T) {
	res := compileSection(t, "examples/spec/cold-resume.cdl")
	var fields []string
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "COLD-RESUME-FIELD" {
			fields = e.Values
		}
	}
	if len(fields) == 0 {
		t.Fatal("COLD-RESUME-FIELD not compiled")
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("8.5. Cold resume check")
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, line := range strings.Split(text, "\n") {
		if match := coldResumeFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			want = append(want, match[1])
		}
	}
	if len(want) == 0 {
		t.Fatal("no field lines found in §8.5")
	}
	compareSets(t, "COLD-RESUME-FIELD", fields, want)
	if !containsString(fields, capabilities.ColdResumeCarrier) {
		t.Fatalf("carrier %q missing from COLD-RESUME-FIELD", capabilities.ColdResumeCarrier)
	}
}

func containsString(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}
