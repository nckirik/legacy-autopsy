package equivalence

import (
	"strings"
	"testing"
)

func TestSuiteReviewRoundTrip(t *testing.T) {
	mappings := []Mapping{
		{RecordType: "UC", RecordID: "UC-ABCDEF012345", Skeleton: "#uc-abcdef012345"},
		{RecordType: "BR", RecordID: "BR-ABCDEF012345", Skeleton: "#br-abcdef012345"},
	}
	answers := []string{"Ayse (operator)", "", "", "a", "e", "covered by approved exclusion DEC-X"}
	index := 0
	prompt := func(question string) (string, error) {
		if index >= len(answers) {
			t.Fatalf("prompt %q has no scripted answer", question)
		}
		answer := answers[index]
		index++
		return answer, nil
	}
	review, err := CollectSuiteReview(prompt, &strings.Builder{}, mappings, SuiteReview{}, "2026-10-14T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if review.Result != "Pass" || len(review.Rows) != 2 || review.Rows[1].Verdict != "Excluded" {
		t.Fatalf("review = %+v", review)
	}
	rendered := RenderSuiteReview(review)
	parsed, err := ParseSuiteReview(rendered)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Result != "Pass" || parsed.Rows[1].Reason == "" {
		t.Fatalf("parsed = %+v", parsed)
	}
}

func TestSuiteReviewSkipIsIncomplete(t *testing.T) {
	mappings := []Mapping{{RecordType: "UC", RecordID: "UC-ABCDEF012345", Skeleton: "#uc"}}
	answers := []string{"Ayse", "", "", "s"}
	index := 0
	prompt := func(question string) (string, error) {
		answer := answers[index]
		index++
		return answer, nil
	}
	review, err := CollectSuiteReview(prompt, &strings.Builder{}, mappings, SuiteReview{}, "2026-10-14T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if review.Result != "Incomplete" {
		t.Fatalf("result = %q, want Incomplete", review.Result)
	}
}
