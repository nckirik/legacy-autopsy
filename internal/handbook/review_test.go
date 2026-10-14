package handbook

import (
	"strings"
	"testing"
)

func scriptedPrompter(answers []string, t *testing.T) (Prompter, *int) {
	index := 0
	return func(question string) (string, error) {
		if index >= len(answers) {
			t.Fatalf("prompt %q has no scripted answer", question)
		}
		answer := answers[index]
		index++
		return answer, nil
	}, &index
}

func TestReadabilityReviewRoundTrip(t *testing.T) {
	chapters := []ChapterPrompt{
		{Path: "chapter-01.md", HBKID: "HBK-AAA000000001", Fingerprint: "sha256:a", Title: "01. Overview", Summary: "The system serves pages."},
		{Path: "chapter-02.md", HBKID: "HBK-AAA000000002", Fingerprint: "sha256:b", Title: "02. Data", Summary: "Data lives in views."},
		{Path: "chapter-03.md", HBKID: "HBK-AAA000000003", Fingerprint: "sha256:c", Title: "03. Jobs", Summary: "Jobs run at night."},
	}
	prompt, _ := scriptedPrompter([]string{"Ali (architect)", "", "", "p", "f", "missing job schedule detail", "s"}, t)
	review, err := CollectReadabilityReview(prompt, &strings.Builder{}, chapters, ReadabilityReviewFile{}, "2026-10-14T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if review.Reviewer != "Ali (architect)" || review.Result != "Fail" || len(review.Verdicts) != 2 {
		t.Fatalf("review = %+v", review)
	}
	rendered := RenderReadabilityReview(review)
	parsed, err := ParseReadabilityReview(rendered)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Result != "Fail" || len(parsed.Verdicts) != 2 || parsed.Verdicts[1].Notes != "missing job schedule detail" {
		t.Fatalf("parsed = %+v", parsed)
	}
	reviews := ReviewsFromFile(parsed)
	findings := VerifyReadabilityCoverage("handbook", []Binding{
		{HBKID: "HBK-AAA000000001", Fingerprint: "sha256:a"},
		{HBKID: "HBK-AAA000000002", Fingerprint: "sha256:b"},
		{HBKID: "HBK-AAA000000003", Fingerprint: "sha256:c"},
	}, reviews)
	if len(findings) == 0 || !strings.Contains(findings[0].Detail, "not covered") {
		t.Fatalf("expected uncovered chapter finding, got %+v", findings)
	}
}

func TestReadabilityReviewRejectsCoverageDrift(t *testing.T) {
	review := ReadabilityReviewFile{
		Reviewer:    "Ali",
		Timestamp:   "2026-10-14T10:00:00Z",
		Scope:       "all",
		Limitations: "static read",
		Result:      "Pass",
		Verdicts: []Verdict{
			{Chapter: "chapter-01.md", HBKID: "HBK-AAA000000001", Fingerprint: "sha256:DRIFT", Verdict: "Pass"},
		},
	}
	findings := VerifyReadabilityCoverage("handbook", []Binding{{HBKID: "HBK-AAA000000001", Fingerprint: "sha256:a"}}, ReviewsFromFile(review))
	if len(findings) == 0 || !strings.Contains(findings[0].Detail, "fingerprint drift") {
		t.Fatalf("expected drift finding, got %+v", findings)
	}
}
