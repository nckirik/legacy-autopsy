package handbook

import (
	"context"
	"strings"
	"testing"
)

type fakeCompleter struct {
	replies []string
	calls   int
}

func (f *fakeCompleter) Complete(ctx context.Context, system, user string) (string, error) {
	reply := f.replies[f.calls]
	f.calls++
	return reply, nil
}

func TestNarrateComposesAndLints(t *testing.T) {
	doc, db := testDocDB(t)
	completer := &fakeCompleter{replies: []string{
		"## Audience Purpose\nThe operator reads this chapter to understand the system.\n\n## Target Decisions\nNot-Applicable: no records exist for this chapter in the current scope.\n",
	}}
	narrative, err := Narrate(context.Background(), completer, db, doc, "01")
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{"## 01. System Overview", "### Audience Purpose", "The operator reads this chapter", "### Observed Legacy Behavior", "### Navigation"} {
		if !strings.Contains(narrative.Markdown, marker) {
			t.Fatalf("narrative missing %q:\n%s", marker, narrative.Markdown)
		}
	}
	if len(narrative.Findings) != 0 {
		t.Fatalf("findings = %v", narrative.Findings)
	}
}

func TestNarrateRetriesWithLintFeedback(t *testing.T) {
	doc, db := testDocDB(t)
	completer := &fakeCompleter{replies: []string{
		"## Audience Purpose\nThe operator can't read this; it must be simpler.\n\n## Target Decisions\nNot-Applicable: no records exist for this chapter in the current scope.\n",
		"## Audience Purpose\nThe operator can read this chapter.\n\n## Target Decisions\nNot-Applicable: no records exist for this chapter in the current scope.\n",
	}}
	narrative, err := Narrate(context.Background(), completer, db, doc, "01")
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 2 {
		t.Fatalf("calls = %d, want 2", completer.calls)
	}
	if !strings.Contains(narrative.Markdown, "The operator can read this chapter.") {
		t.Fatalf("retry output not used:\n%s", narrative.Markdown)
	}
}

func TestNarrateFailsClosedOnPersistentLint(t *testing.T) {
	doc, db := testDocDB(t)
	bad := "## Audience Purpose\nThe operator can't read this.\n\n## Target Decisions\nNot-Applicable: no records exist for this chapter in the current scope.\n"
	completer := &fakeCompleter{replies: []string{bad, bad, bad}}
	if _, err := Narrate(context.Background(), completer, db, doc, "01"); err == nil {
		t.Fatal("persistent lint failure accepted")
	}
}
