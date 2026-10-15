package parity

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	srcKindForm = regexp.MustCompile(`^SRC-\(([A-Z|]+)\)-\[A-F0-9\]`)
	prefixForm  = regexp.MustCompile(`^\(([A-Z|]+)\)-\[A-F0-9\]`)
	iterForm    = regexp.MustCompile(`^\(([A-Z|-]+)\)-\[A-Z\]\[A-Z0-9\]`)
)

// TestTypedIDRegistryDrift checks that the CDL §4.1 registries exactly match the
// prefix and iteration forms declared in protocol.md.
func TestTypedIDRegistryDrift(t *testing.T) {
	res := compileSection(t, "protocol/typed-id.cdl")
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("4.1. ID generation")
	if err != nil {
		t.Fatal(err)
	}
	var wantPrefixes, wantIterations []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case iterForm.MatchString(line):
			wantIterations = append(wantIterations, strings.Split(iterForm.FindStringSubmatch(line)[1], "|")...)
		case srcKindForm.MatchString(line):
			for _, v := range strings.Split(srcKindForm.FindStringSubmatch(line)[1], "|") {
				wantPrefixes = append(wantPrefixes, "SRC-"+v)
			}
		case prefixForm.MatchString(line):
			wantPrefixes = append(wantPrefixes, strings.Split(prefixForm.FindStringSubmatch(line)[1], "|")...)
		}
	}
	if len(wantPrefixes) == 0 || len(wantIterations) == 0 {
		t.Fatal("failed to parse §4.1 prefix/iteration forms from protocol.md")
	}
	compareSets(t, "ID-PREFIX", enums["ID-PREFIX"], wantPrefixes)
	compareSets(t, "ID-ITERATION", enums["ID-ITERATION"], wantIterations)
}

// compareSetsWithAdditions compares a CDL set against the frozen bootstrap
// oracle set while allowing declared post-4.1.3 additions: every oracle value
// must remain present (nothing removed or renamed) and every extra value must
// be listed as an accepted revision addition.
func compareSetsWithAdditions(t *testing.T, name string, got, want, additions []string) {
	t.Helper()
	allowed := map[string]bool{}
	for _, addition := range additions {
		allowed[addition] = true
	}
	filtered := make([]string, 0, len(got))
	for _, value := range got {
		if allowed[value] {
			continue
		}
		filtered = append(filtered, value)
	}
	for _, addition := range additions {
		found := false
		for _, value := range got {
			if value == addition {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s declared addition %q is missing from the CDL", name, addition)
		}
	}
	compareSets(t, name, filtered, want)
}

func compareSets(t *testing.T, name string, got, want []string) {
	t.Helper()
	g := append([]string(nil), got...)
	w := append([]string(nil), want...)
	sort.Strings(g)
	sort.Strings(w)
	if strings.Join(g, ",") != strings.Join(w, ",") {
		t.Fatalf("%s drift:\n cdl:      %v\n protocol: %v", name, g, w)
	}
}
