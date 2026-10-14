package handbook

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
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

func seedModules(t *testing.T, db string, doc *cdl.EIRDoc, count int) []string {
	t.Helper()
	var claims []string
	for i := 1; i <= count; i++ {
		claim := "CLM-0000000000" + string(rune('0'+(i/10))) + string(rune('0'+(i%10)))
		claims = append(claims, claim)
		id := "MOD-00000000000" + string(rune('0'+i))
		payload := "- **source-bindings:** " + claim + "\n" +
			"- **block-id:** " + id + "\n" +
			"- **record-title:** Module " + string(rune('A'+i-1)) + "\n" +
			"- **domain-path:** app/m" + string(rune('0'+i)) + "\n" +
			"- **responsibility-domain-boundary:** The module owns one boundary.\n" +
			"- **public-interfaces:** None\n" +
			"- **encapsulated-state-entities:** None\n" +
			"- **dependencies-dependents:** None\n" +
			"- **data-and-events-passed:** None\n" +
			"- **failure-propagation:** None\n" +
			"- **invariants-rules-owned:** None\n" +
			"- **observed-legacy-placement:** app\n" +
			"- **target-decision:** Pending\n" +
			"- **personas-use-cases-capabilities:** None\n"
		if err := store.Ingest(db, doc, "ARCHITECTURE-MODULE", id, 1, []byte(payload)); err != nil {
			t.Fatal(err)
		}
		claimPayload := "- **subject-id-record-version:** SRC-FILE-AAAAAAAAAAAA\n- **field-step:** file\n- **claim:** module " + string(rune('A'+i-1)) + " owns one boundary\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** app/m" + string(rune('0'+i)) + ":1\n- **source-fingerprint-snapshot:** sha256:x\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n"
		if err := store.Ingest(db, doc, "CLAIM", claim, 1, []byte(claimPayload)); err != nil {
			t.Fatal(err)
		}
	}
	return claims
}

func goodNarrative(claims []string) string {
	var b strings.Builder
	b.WriteString("## Summary\nThe operator reads this chapter to understand the system.\n\n## Observed Legacy Behavior\n")
	for i, claim := range claims {
		b.WriteString("The system gives module " + string(rune('A'+i)) + " one clear responsibility (" + claim + ").\n")
	}
	b.WriteString("Each module owns its boundary and its interfaces.\n\n## Target Decisions\nNot-Applicable: no records exist for this chapter in the current scope.\n")
	return b.String()
}

func TestNarrateComposesProseChapter(t *testing.T) {
	doc, db := testDocDB(t)
	claims := seedModules(t, db, doc, 8)
	completer := &fakeCompleter{replies: []string{goodNarrative(claims)}}
	narrative, err := Narrate(context.Background(), completer, db, doc, "01")
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{
		"## 01. System Overview",
		"### Summary",
		"### Observed Legacy Behavior",
		"The system gives module A one clear responsibility",
		"### Target Decisions",
		"### Navigation",
		"### Record Details",
		"<details><summary>MOD-",
	} {
		if !strings.Contains(narrative.Markdown, marker) {
			t.Fatalf("narrative missing %q:\n%s", marker, narrative.Markdown)
		}
	}
	if len(narrative.Findings) != 0 {
		t.Fatalf("findings = %v", narrative.Findings)
	}
}

func TestNarrateRetriesOnCitationsAndStyle(t *testing.T) {
	doc, db := testDocDB(t)
	claims := seedModules(t, db, doc, 8)
	bad := "## Summary\nThe operator can't read module A (CLM-FFFFFFFFFFFF).\n\n## Observed Legacy Behavior\nA short paragraph.\n\n## Target Decisions\nNone.\n"
	completer := &fakeCompleter{replies: []string{bad, goodNarrative(claims)}}
	narrative, err := Narrate(context.Background(), completer, db, doc, "01")
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 2 {
		t.Fatalf("calls = %d, want 2", completer.calls)
	}
	if !strings.Contains(narrative.Markdown, "The system gives module A") {
		t.Fatalf("retry output not used:\n%s", narrative.Markdown)
	}
}

func TestNarrateFailsClosedOnPersistentProblems(t *testing.T) {
	doc, db := testDocDB(t)
	seedModules(t, db, doc, 8)
	bad := "## Summary\nThe operator can't read this (CLM-FFFFFFFFFFFF).\n\n## Observed Legacy Behavior\nMissing citations.\n\n## Target Decisions\nNone.\n"
	completer := &fakeCompleter{replies: []string{bad, bad, bad, bad, bad}}
	if _, err := Narrate(context.Background(), completer, db, doc, "01"); err == nil {
		t.Fatal("persistent validation failure accepted")
	}
}

func TestNarrateNumbersCitationsAndIndexesEvidence(t *testing.T) {
	doc, db := testDocDB(t)
	claims := seedModules(t, db, doc, 8)
	completer := &fakeCompleter{replies: []string{goodNarrative(claims)}}
	narrative, err := Narrate(context.Background(), completer, db, doc, "01")
	if err != nil {
		t.Fatal(err)
	}
	body := narrative.Markdown
	if !strings.Contains(body, "one clear responsibility [^1].") {
		t.Fatalf("citations not numbered:\n%s", body)
	}
	if strings.Contains(strings.Split(body, "### Evidence Notes")[0], "CLM-") {
		t.Fatalf("machine IDs leaked into the reading surface:\n%s", body)
	}
	for i, claim := range claims {
		marker := fmt.Sprintf("[^%d]:", i+1)
		link := fmt.Sprintf("[%s](evidence.md#%s)", claim, ClaimAnchor(claim))
		if !strings.Contains(body, marker) || !strings.Contains(body, link) {
			t.Fatalf("evidence notes missing %q or %q", marker, link)
		}
	}
}
