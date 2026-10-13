package handbook

import (
	"context"
	"fmt"
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

const narrativeSystemPrompt = `You write one handbook chapter for a legacy-system reconstruction. The canonical protocol state is authoritative; you never invent records or evidence.

` + ste.Directive + `

Output exactly two sections, nothing else:

## Audience Purpose
<one short paragraph>

## Target Decisions
<one short paragraph, or "Not-Applicable: <reason>" when no decision applies>

Do not add headings, record lists, navigation, or commentary. The runtime adds the observed records and navigation.`

// Narrate authors the narrative sections for one chapter from its committed
// records and returns the composed chapter with lint findings. It retries once
// with the lint findings as feedback.
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
	for _, catalogID := range spec.Catalogs {
		records, err := loadRecords(db, doc, catalogID)
		if err != nil {
			return ChapterNarrative{}, err
		}
		observed = append(observed, records...)
	}
	user := narrativeUserPrompt(*spec, observed)
	var feedback string
	for attempt := 0; attempt < 3; attempt++ {
		if feedback != "" {
			user = narrativeUserPrompt(*spec, observed) + "\n\nPrevious attempt rejected by lint:\n" + feedback
		}
		text, err := completer.Complete(ctx, narrativeSystemPrompt, user)
		if err != nil {
			return ChapterNarrative{}, err
		}
		purpose, decisions := splitNarrative(text)
		if purpose == "" || decisions == "" {
			feedback = "output must contain both ## Audience Purpose and ## Target Decisions sections"
			continue
		}
		chapter := Chapter{
			Number:          spec.Number,
			Title:           spec.Title,
			AudiencePurpose: purpose,
			Observed:        observed,
			TargetDecisions: decisions,
			Navigation:      navigationFor(specs, spec.Number),
		}
		body := chapterBody(Assemble([]Chapter{chapter}))
		// The STE gate judges the generated narrative. Observed record prose is
		// governed at synthesis time; historical records must not fail the
		// narrative retry loop.
		narrativeOnly := chapter
		narrativeOnly.Observed = nil
		findings := LintChapter(spec.Number, string(chapterBody(Assemble([]Chapter{narrativeOnly}))))
		if len(findings) == 0 {
			return ChapterNarrative{Chapter: spec.Number, Markdown: string(body)}, nil
		}
		var messages []string
		for _, finding := range findings {
			messages = append(messages, finding.String())
		}
		feedback = strings.Join(messages, "\n")
	}
	return ChapterNarrative{}, fmt.Errorf("handbook: chapter %s did not pass lint: %s", chapterNumber, feedback)
}

func narrativeUserPrompt(spec ChapterSpec, observed []Record) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Chapter %s: %s\n\nObserved synthesis records:\n", spec.Number, spec.Title)
	if len(observed) == 0 {
		b.WriteString("- None. Use a reason-bearing Not-Applicable for target decisions.\n")
	}
	for _, record := range observed {
		fmt.Fprintf(&b, "\n### %s\n\n%s\n", record.ID, record.Payload)
	}
	return b.String()
}

func splitNarrative(text string) (string, string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var purpose, decisions strings.Builder
	current := &purpose
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "## Audience Purpose"):
			current = &purpose
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
	return strings.TrimSpace(purpose.String()), strings.TrimSpace(decisions.String())
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
