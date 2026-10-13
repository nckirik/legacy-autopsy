package handbook

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/ste"
)

// Completer is one bounded model call (satisfied by the llm adapter).
type Completer interface {
	Complete(ctx context.Context, system, user string) (string, error)
}

// ChapterNarrative is the rendered narrative for one chapter with lint findings.
type ChapterNarrative struct {
	Chapter  string
	Markdown string
	Findings []Finding
}

const narrativeSystemPrompt = `You write one chapter of a software architecture reconstruction handbook. Fresh developers read this chapter to understand the legacy system. Write it as a design document in prose, not as bullet lists or field dumps.

The canonical protocol state is authoritative. The structured records below are your evidence. Do not invent behavior, names, or evidence.

Citation rules:
- Every observed statement must cite the claim IDs that support it, in parentheses: (CLM-0123ABCDEF45).
- Every record listed under Inputs must be represented by at least one cited statement.
- Never cite a claim that is not in the Inputs evidence list.

` + ste.Directive + `

Output exactly three sections and nothing else:

## Audience Purpose
<one paragraph: who reads this chapter and what they need from it>

## Observed Legacy Behavior
<several paragraphs of narrative prose that explain the observed behavior as a coherent design description. Use short paragraphs and lists only when a list is the clear form. Cite claims after each statement.>

## Target Decisions
<one paragraph, or "Not-Applicable: <reason>" when no decision applies>

Do not add other headings, navigation, record IDs, or field tables. The runtime adds navigation and a collapsed record appendix.`

// Narrate authors the full chapter narrative from its committed records and
// validates citations, the STE rules, and the substance floor. It retries with
// the validation findings as feedback.
func Narrate(ctx context.Context, completer Completer, db string, doc *cdl.EIRDoc, chapterNumber string) (ChapterNarrative, error) {
	specs := DefaultChapters()
	var spec *ChapterSpec
	for i := range specs {
		if specs[i].Number == chapterNumber {
			spec = &specs[i]
		}
	}
	if spec == nil {
		return ChapterNarrative{}, fmt.Errorf("handbook: unknown chapter %q", chapterNumber)
	}
	var observed []Record
	evidence := map[string]bool{}
	recordClaims := map[string][]string{}
	for _, catalogID := range spec.Catalogs {
		records, err := loadRecords(db, doc, catalogID)
		if err != nil {
			return ChapterNarrative{}, err
		}
		for _, record := range records {
			claims := claimsIn(string(record.Payload))
			recordClaims[record.ID] = claims
			for _, claim := range claims {
				evidence[claim] = true
			}
		}
		observed = append(observed, records...)
	}
	user := narrativeUserPrompt(*spec, observed, evidence)
	var feedback string
	for attempt := 0; attempt < 3; attempt++ {
		if feedback != "" {
			user = narrativeUserPrompt(*spec, observed, evidence) + "\n\nPrevious attempt rejected:\n" + feedback
		}
		text, err := completer.Complete(ctx, narrativeSystemPrompt, user)
		if err != nil {
			return ChapterNarrative{}, err
		}
		purpose, behavior, decisions := splitNarrative(text)
		var problems []string
		if purpose == "" || behavior == "" || decisions == "" {
			problems = append(problems, "output must contain ## Audience Purpose, ## Observed Legacy Behavior, and ## Target Decisions sections")
		}
		for _, claim := range claimsIn(text) {
			if !evidence[claim] {
				problems = append(problems, fmt.Sprintf("citation %s is not in the chapter evidence list", claim))
			}
		}
		for recordID, claims := range recordClaims {
			if len(claims) == 0 {
				continue
			}
			covered := false
			for _, claim := range claims {
				if strings.Contains(text, claim) {
					covered = true
					break
				}
			}
			if !covered {
				problems = append(problems, fmt.Sprintf("record %s is not represented by any cited statement", recordID))
			}
		}
		narrative := composeNarrative(spec, purpose, behavior, decisions)
		for _, finding := range LintChapter(spec.Number, narrative) {
			if !strings.HasPrefix(finding.Code, "STE-") {
				problems = append(problems, finding.String())
			}
		}
		for _, finding := range ste.Lint(spec.Number, purpose+"\n"+behavior+"\n"+decisions) {
			problems = append(problems, finding.String())
		}
		if len(problems) == 0 {
			numbered, index := numberCitations(narrative)
			full := numbered + index + recordAppendix(observed)
			return ChapterNarrative{Chapter: spec.Number, Markdown: full}, nil
		}
		feedback = strings.Join(problems, "\n")
	}
	return ChapterNarrative{}, fmt.Errorf("handbook: chapter %s did not pass validation: %s", chapterNumber, feedback)
}

var claimPattern = regexp.MustCompile(`CLM-[A-F0-9]{12}`)

func claimsIn(text string) []string {
	seen := map[string]bool{}
	var claims []string
	for _, match := range claimPattern.FindAllString(text, -1) {
		if !seen[match] {
			seen[match] = true
			claims = append(claims, match)
		}
	}
	sort.Strings(claims)
	return claims
}

func composeNarrative(spec *ChapterSpec, purpose, behavior, decisions string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s. %s\n\n", spec.Number, spec.Title)
	fmt.Fprintf(&b, "### Audience Purpose\n\n%s\n\n", purpose)
	fmt.Fprintf(&b, "### Observed Legacy Behavior\n\n%s\n\n", behavior)
	fmt.Fprintf(&b, "### Target Decisions\n\n%s\n\n", decisions)
	fmt.Fprintf(&b, "### Navigation\n\n")
	specs := DefaultChapters()
	for _, item := range navigationFor(specs, spec.Number) {
		fmt.Fprintf(&b, "- %s\n", item)
	}
	return b.String()
}

// numberCitations replaces claim citations with bracketed numbers in order of
// first appearance and returns the chapter text plus its evidence index, so the
// reading surface carries no machine IDs.
func numberCitations(narrative string) (string, string) {
	order := []string{}
	numbers := map[string]int{}
	replace := func(match string) string {
		claim := claimPattern.FindString(match)
		if claim == "" {
			return match
		}
		if _, ok := numbers[claim]; !ok {
			order = append(order, claim)
			numbers[claim] = len(order)
		}
		return fmt.Sprintf("[^%d]", numbers[claim])
	}
	text := regexp.MustCompile(`\(CLM-[A-F0-9]{12}\)`).ReplaceAllStringFunc(narrative, replace)
	text = claimPattern.ReplaceAllStringFunc(text, replace)
	if len(order) == 0 {
		return text, ""
	}
	var b strings.Builder
	b.WriteString("\n### Evidence Notes\n\n")
	for i, claim := range order {
		fmt.Fprintf(&b, "[^%d]: `%s`\n", i+1, claim)
	}
	return text, b.String()
}

func recordAppendix(observed []Record) string {
	if len(observed) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n### Record Details\n\n")
	for _, record := range observed {
		fmt.Fprintf(&b, "<details><summary>%s</summary>\n\n%s\n</details>\n\n", record.ID, record.Payload)
	}
	return b.String()
}

func narrativeUserPrompt(spec ChapterSpec, observed []Record, evidence map[string]bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Chapter %s: %s\n\n", spec.Number, spec.Title)
	if len(observed) == 0 {
		b.WriteString("Inputs: None. Write a short chapter. Use a reason-bearing Not-Applicable for target decisions.\n")
		return b.String()
	}
	b.WriteString("Inputs (structured evidence records):\n\n")
	for _, record := range observed {
		fmt.Fprintf(&b, "%s\n\n", record.Payload)
	}
	claims := make([]string, 0, len(evidence))
	for claim := range evidence {
		claims = append(claims, claim)
	}
	sort.Strings(claims)
	fmt.Fprintf(&b, "Evidence claim IDs (cite only these): %s\n", strings.Join(claims, ", "))
	return b.String()
}

func splitNarrative(text string) (string, string, string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var purpose, behavior, decisions strings.Builder
	current := &purpose
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "## Audience Purpose"):
			current = &purpose
			continue
		case strings.HasPrefix(trimmed, "## Observed Legacy Behavior"):
			current = &behavior
			continue
		case strings.HasPrefix(trimmed, "## Target Decisions"):
			current = &decisions
			continue
		case strings.HasPrefix(trimmed, "#"):
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")
	}
	return strings.TrimSpace(purpose.String()), strings.TrimSpace(behavior.String()), strings.TrimSpace(decisions.String())
}

func navigationFor(specs []ChapterSpec, number string) []string {
	var navigation []string
	for i := range specs {
		if specs[i].Number != number {
			continue
		}
		if i > 0 {
			navigation = append(navigation, fmt.Sprintf("previous: %s-%s", specs[i-1].Number, specs[i-1].Title))
		}
		if i < len(specs)-1 {
			navigation = append(navigation, fmt.Sprintf("next: %s-%s", specs[i+1].Number, specs[i+1].Title))
		}
	}
	return navigation
}
