package parity

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	fsmLine     = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9 /-]*) --[^>]*--> (.+)$`)
	ticketToken = regexp.MustCompile(`T-[A-Z-]+`)
)

// TestTicketFSMDrift checks the §9.2 transition table, the ticket state domain,
// and the §9.1 escalation-reason enum against protocol.md.
func TestTicketFSMDrift(t *testing.T) {
	res := compileSection(t, "protocol/ticket-fsm.cdl")
	var allows [][2]string
	var states []string
	for _, st := range res.EIR.Declarations.States {
		if st.ID == "TICKET-STATE" {
			allows = st.Allows
			states = st.Values
		}
	}
	if len(allows) == 0 || len(states) == 0 {
		t.Fatal("TICKET-STATE not compiled")
	}
	cdlText := ""
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += r.Goal + "\n"
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.SectionText("9.2. FSM and write rights")
	if err != nil {
		t.Fatal(err)
	}

	var wantPairs [][2]string
	wildcard := false
	for _, line := range strings.Split(text, "\n") {
		match := fsmLine.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		left := match[1]
		right := strings.TrimSpace(match[2])
		if left == "Active state" {
			wildcard = true
			continue
		}
		sources := strings.Split(left, "/")
		if idx := strings.Index(right, "->"); idx >= 0 {
			first := ticketToken.FindString(right[:idx])
			rest := ticketToken.FindAllString(right[idx+2:], -1)
			for _, source := range sources {
				wantPairs = append(wantPairs, [2]string{source, first})
			}
			for _, target := range rest {
				wantPairs = append(wantPairs, [2]string{first, target})
			}
			continue
		}
		for _, source := range sources {
			for _, target := range ticketToken.FindAllString(right, -1) {
				wantPairs = append(wantPairs, [2]string{source, target})
			}
		}
	}
	if len(wantPairs) == 0 {
		t.Fatal("no transitions parsed from §9.2")
	}
	comparePairs(t, "TICKET-STATE ALLOWS", allows, wantPairs)
	if !wildcard {
		t.Fatal("active-state wildcard transitions not found in §9.2")
	}
	for _, target := range []string{"T-INVALID", "T-SUPERSEDED", "T-OUT-OF-SCOPE"} {
		if !strings.Contains(cdlText, target) {
			t.Fatalf("wildcard target %s missing from the CDL rules", target)
		}
	}

	seenStates := map[string]bool{}
	var wantStates []string
	for _, token := range ticketToken.FindAllString(text, -1) {
		if !seenStates[token] {
			seenStates[token] = true
			wantStates = append(wantStates, token)
		}
	}
	compareSets(t, "TICKET-STATE", states, wantStates)

	text91, err := model.SectionText("9.1. Ticket schema and canonical IDs")
	if err != nil {
		t.Fatal(err)
	}
	var wantReasons []string
	for _, line := range strings.Split(text91, "\n") {
		idx := strings.Index(line, "Escalation Reason:**")
		if idx < 0 {
			continue
		}
		for _, value := range strings.Split(line[idx+len("Escalation Reason:**"):], "|") {
			value = strings.TrimSpace(value)
			if value != "" {
				wantReasons = append(wantReasons, value)
			}
		}
		break
	}
	var reasons []string
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "TICKET-ESCALATION-REASON" {
			reasons = e.Values
		}
	}
	if len(reasons) == 0 || len(wantReasons) == 0 {
		t.Fatal("escalation reasons not found")
	}
	compareSets(t, "TICKET-ESCALATION-REASON", reasons, wantReasons)
}

func comparePairs(t *testing.T, name string, got, want [][2]string) {
	t.Helper()
	key := func(p [2]string) string { return p[0] + "->" + p[1] }
	g := make([]string, 0, len(got))
	for _, p := range got {
		g = append(g, key(p))
	}
	w := make([]string, 0, len(want))
	for _, p := range want {
		w = append(w, key(p))
	}
	sort.Strings(g)
	sort.Strings(w)
	if strings.Join(g, ",") != strings.Join(w, ",") {
		t.Fatalf("%s drift:\n cdl:      %v\n protocol: %v", name, g, w)
	}
}
