package parity

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
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
	root := root(t)
	src, err := os.ReadFile(filepath.Join(root, "examples/spec/typed-id.cdl"))
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := cdl.LoadLedger(filepath.Join(root, "examples/spec/identity-ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := cdl.Compile(cdl.CompileInput{
		Sources: []cdl.Source{{Path: "examples/spec/typed-id.cdl", Bytes: src}},
		Ledger:  ledger,
	})
	if err != nil {
		t.Fatal(err)
	}
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(filepath.Join(root, "protocol.md"))
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
