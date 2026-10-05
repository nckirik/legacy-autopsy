package render

import (
	"strings"
	"testing"
)

func TestDecisionBatchMarkdown(t *testing.T) {
	out := string(Markdown(Batch{
		Title:   "Scope decisions",
		Changed: []string{"kept ledger retry policy"},
		Needed:  []string{"approve exclusion of legacy batch jobs"},
		Items: []Item{
			{ID: "DEC-0001", Summary: "retry policy retained", Payload: []byte("- **Type:** Scope\n- **Rationale:** bounded retries suffice.")},
			{ID: "DEC-0002", Payload: []byte("- **Type:** Exclusion")},
		},
	}))
	for _, marker := range []string{
		"# Scope decisions",
		"## First screen",
		"- **Changed:**\n  - kept ledger retry policy",
		"- **Blocked:** None",
		"- **Needed:**\n  - approve exclusion of legacy batch jobs",
		"## Index",
		"- DEC-0001 - retry policy retained",
		"- DEC-0002 - no summary",
		"### DEC-0001",
		"### DEC-0002",
	} {
		if !strings.Contains(out, marker) {
			t.Fatalf("batch missing %q:\n%s", marker, out)
		}
	}
	if strings.Count(out, "- DEC-0001 - ") != 1 {
		t.Fatal("index must be one line per item")
	}
	if strings.Index(out, "## First screen") > strings.Index(out, "## Index") {
		t.Fatal("first screen must precede the index")
	}
}

func TestQuestionsMarkdown(t *testing.T) {
	out := string(QuestionsMarkdown([]Question{{
		ID:     "Q-1",
		Prompt: "Should legacy batch jobs be excluded from scope?",
		Options: []Option{
			{Label: "Exclude", Consequence: "denominator shrinks; residual risk recorded"},
			{Label: "Keep", Consequence: "static coverage continues"},
		},
	}}))
	for _, marker := range []string{"# Operator questions", "## Q-1", "Should legacy batch jobs", "- **Exclude** - denominator shrinks", "- **Keep** - static coverage continues"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("questions missing %q:\n%s", marker, out)
		}
	}
	if out := string(QuestionsMarkdown(nil)); !strings.Contains(out, "None") {
		t.Fatalf("empty questions = %q", out)
	}
}
