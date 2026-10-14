package handbook

import (
	"os"
	"path/filepath"
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

func TestAttachConfirmationPreservesFingerprint(t *testing.T) {
	dir := t.TempDir()
	chapter := "## 01. Overview\n\n### Summary\n\nPlain prose.\n\n### Observed Legacy Behavior\n\nFact [^1].\n"
	if err := os.WriteFile(filepath.Join(dir, "chapter-01.md"), []byte(chapter), 0o644); err != nil {
		t.Fatal(err)
	}
	resolve := func(id string) (UpstreamRef, bool) { return UpstreamRef{}, false }
	sections, err := MintSections(dir, "demo", resolve)
	if err != nil {
		t.Fatal(err)
	}
	attached, err := AttachConfirmation(sections[0].Bytes, "CNF-ABCDEF012345", "Ali (architect)", "2026-10-14T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	status, confirmation := EnvelopeStatus(attached)
	if status != "B-CONFIRMED" || confirmation != "CNF-ABCDEF012345" {
		t.Fatalf("envelope = %q %q", status, confirmation)
	}
	if err := os.WriteFile(filepath.Join(dir, "chapter-01.md"), attached, 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadMintedSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded[0].Fingerprint != sections[0].Fingerprint {
		t.Fatalf("confirmation changed the semantic fingerprint")
	}
	if string(SectionContent(loaded[0].Bytes)) != SectionContent(sections[0].Bytes) {
		t.Fatalf("confirmation changed the payload content")
	}
}
