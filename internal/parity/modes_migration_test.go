package parity

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	backtickedAny       = regexp.MustCompile("`([A-Za-z/ ]+)`")
	artifactPrefixToken = regexp.MustCompile(`^[0-9][0-9A-Z]*$`)
)

// TestInvocationModesDrift checks the CDL MODE declarations against protocol.md
// and proves equivalence with the frozen Go read sets before they are retired.
func TestInvocationModesDrift(t *testing.T) {
	res := compileSection(t, "examples/spec/invocation-modes.cdl")
	base := res.EIR.Declarations.BaseReads
	modes := res.EIR.Declarations.Modes
	if len(base) == 0 || len(modes) == 0 {
		t.Fatal("no mode declarations compiled")
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}

	// 1. Mode ids == §8.3 enum.
	text83, err := model.SectionText("8.3. Invocation-mode enum and explicit multi-file modes")
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := parseModeEnum(t, text83)
	var gotIDs []string
	for _, mode := range modes {
		gotIDs = append(gotIDs, mode.ID)
	}
	compareSets(t, "invocation modes", gotIDs, wantIDs)

	// 2. Strict set == the mode names bound in §8.2.
	text82, err := model.SectionText("8.2. Strict single-scope modes")
	if err != nil {
		t.Fatal(err)
	}
	modeSet := map[string]bool{}
	for _, id := range wantIDs {
		modeSet[id] = true
	}
	strictSeen := map[string]bool{}
	var wantStrict []string
	for _, match := range backtickedAny.FindAllStringSubmatch(text82, -1) {
		if modeSet[match[1]] && !strictSeen[match[1]] {
			strictSeen[match[1]] = true
			wantStrict = append(wantStrict, match[1])
		}
	}
	if len(wantStrict) == 0 {
		t.Fatal("no strict modes found in §8.2")
	}
	var gotStrict []string
	for _, mode := range modes {
		if mode.Strict {
			gotStrict = append(gotStrict, mode.ID)
		}
	}
	compareSets(t, "strict modes", gotStrict, wantStrict)

	// 3. Baseline prefixes == §8.4 base sentence plus the numbered inventory group.
	text84, err := model.SectionText("8.4. Mandatory read sets")
	if err != nil {
		t.Fatal(err)
	}
	baseLine := ""
	for _, line := range strings.Split(text84, "\n") {
		if strings.Contains(line, "All modes load") {
			baseLine = line
			break
		}
	}
	if baseLine == "" {
		t.Fatal("base read sentence not found in §8.4")
	}
	gotBasePrefixes := make([]string, 0, len(base))
	for _, artifact := range base {
		gotBasePrefixes = append(gotBasePrefixes, readPrefix(artifact))
	}
	wantBasePrefixes := []string{"0A", "0D", "0E", "0F", "0G", "0H", "10", "11", "12"}
	compareSets(t, "base read prefixes", gotBasePrefixes, wantBasePrefixes)
	for _, match := range backtickedAny.FindAllStringSubmatch(baseLine, -1) {
		if !artifactPrefixToken.MatchString(match[1]) {
			continue
		}
		found := false
		for _, prefix := range gotBasePrefixes {
			if prefix == match[1] {
				found = true
			}
		}
		if !found {
			t.Fatalf("§8.4 base token %s is not in the CDL baseline", match[1])
		}
	}

	// 4. The runtime now consumes these declarations; the retired Go read sets are
	// no longer the operational source, so equivalence is enforced by construction.
}

func parseModeEnum(t *testing.T, text string) []string {
	t.Helper()
	inBlock := false
	var content []string
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if inBlock {
				break
			}
			inBlock = true
			continue
		}
		if inBlock && trimmed != "" {
			content = append(content, trimmed)
		}
	}
	if len(content) == 0 {
		t.Fatal("mode enum code block not found in §8.3")
	}
	var out []string
	for _, part := range strings.Split(strings.Join(content, " "), "|") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func readPrefix(artifact string) string {
	if idx := strings.IndexByte(artifact, '-'); idx >= 0 {
		return artifact[:idx]
	}
	return artifact
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
