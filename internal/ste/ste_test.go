package ste

import (
	"strings"
	"testing"
)

func TestLintPassesCleanText(t *testing.T) {
	text := "The scheduler retries failed jobs. The operator sets the retry limit. The system writes one log entry for each failure."
	if findings := Lint("chapter", text); len(findings) != 0 {
		t.Fatalf("clean text rejected: %v", findings)
	}
}

func TestLintFindsViolations(t *testing.T) {
	long := "The system must write one audit entry for every accepted proposal and the runtime must also write one entry for every rejected proposal and every abandoned work item so that traceability stays complete forever"
	cases := map[string]struct {
		text string
		code string
	}{
		"long sentence": {long, "STE-SENTENCE-LENGTH"},
		"contraction":   {"The operator can't delete the record.", "STE-CONTRACTION"},
		"semicolon":     {"The system stores claims; the runtime recomputes hashes.", "STE-SEMICOLON"},
		"n/a":           {"The value is N/A for this record.", "STE-EMPTY-VALUE"},
		"should":        {"The operator should approve the decision.", "STE-SHOULD"},
	}
	for name, tc := range cases {
		findings := Lint("text", tc.text)
		found := false
		for _, finding := range findings {
			if finding.Code == tc.code {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s: expected %s in %v", name, tc.code, findings)
		}
	}
}

func TestDirectiveCoversKeyRules(t *testing.T) {
	for _, marker := range []string{"ASD-STE100", "20 words", "25 words", "active voice", `"must"`, `"None"`} {
		if !strings.Contains(Directive, marker) {
			t.Fatalf("directive missing %q", marker)
		}
	}
}
