package handbook

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/ste"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
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

const narrativeSystemPrompt = `You write one chapter of a software architecture reconstruction handbook. Write it as a person explaining the system to a colleague: direct, natural, and clear. Avoid template formulas, repeated sentence openings, and filler. Do not open with "a fresh developer" or describe who reads the chapter.

The canonical protocol state is authoritative. The structured records below are your evidence. Do not invent behavior, names, or evidence.

` + ste.Directive + `

Output exactly three sections and nothing else:

## Summary
<two or three plain paragraphs about what this part of the system does for its users and how work moves through it. Generalized terms are welcome: web application, database, database view, container, build pipeline, version control, template, session, token, sign-in, and the like. Do not name specific products, vendors, libraries, languages, frameworks, hosts, ports, file or directory paths, routes, code identifiers, or record identifiers, and add no citations. Describe the surface and the overall flow, not the implementation details. This section replaces the older audience-purpose paragraph, so it may say briefly what the reader gets from the chapter, but it must mostly describe the system.>

## Observed Legacy Behavior
<several paragraphs of narrative prose that explain the observed behavior as a coherent design description, with the concrete technical detail. Write flowing prose: connect related facts in one sentence instead of one fact per sentence, and vary how sentences begin. Use short paragraphs and lists only when a list is the clear form. Put each cited statement on its own line (a single line break between statements; a blank line starts a new paragraph). State the observed facts fully: cover the evidence thoroughly and aim for at least eight separate cited statements whenever the evidence supports them. Cite claims after each statement: (CLM-0123ABCDEF45). Every record listed under Inputs must be represented by at least one cited statement. Never cite a claim that is not in the Inputs evidence list.>

## Target Decisions
<one paragraph, or "Not-Applicable: <reason>" when no decision applies>

Do not add other headings, navigation, record IDs, or field tables. The runtime adds navigation and a collapsed record appendix.

Diagram and pseudocode blocks:
- Designated flow-heavy chapters (02 Architecture, 03 Domain and Data, 05 Use Cases, 08 Background Processing) must include at least one fenced Mermaid diagram (code fence with info string mermaid) that projects the cited records: use flowchart/graph for flows, stateDiagram-v2 for lifecycles, erDiagram for entity relationships (chapter 03 must show the entity/table inventory with keys and the relationship topology as an erDiagram), and sequenceDiagram for message flows. When the chapter has no backing records at all, write a reason-bearing "Not-Applicable: <reason>" instead.
- Chapters 03 and 08 carry procedures by default (database routines, faults, schedules): when a procedure or lifecycle cannot be stated precisely in prose (validation order, locking, retry/rollback, state transitions), include a language-agnostic pseudocode block (code fence with info string text) using neutral constructs and no product, language, or framework syntax. Otherwise write a reason-bearing "Not-Applicable: <reason>".
- Diagram and pseudocode content is a projection of the cited records: every element traces to a claim you cite, and citations stay in the surrounding prose. Never assert a fact in a diagram or pseudocode block that the cited records do not support.`

// Narrate authors the full chapter narrative from its committed records and
// validates citations, the STE rules, and the substance floor. It retries with
// the validation findings as feedback.
func Narrate(ctx context.Context, completer Completer, db string, doc *cdl.EIRDoc, chapterNumber, feedback string) (ChapterNarrative, error) {
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
	claimNotes, err := loadClaimNotes(db, doc, evidence)
	if err != nil {
		return ChapterNarrative{}, err
	}
	user := narrativeUserPrompt(*spec, observed, evidence)
	if strings.TrimSpace(feedback) != "" {
		user += "\n\nReviewer feedback to address in this revision:\n" + strings.TrimSpace(feedback) + "\n"
	}
	var rejection string
	lastRejection := ""
	for attempt := 0; attempt < 5; attempt++ {
		if rejection != "" {
			user = narrativeUserPrompt(*spec, observed, evidence) + "\n\nPrevious attempt rejected:\n" + rejection
		}
		text, err := completer.Complete(ctx, narrativeSystemPrompt, user)
		if err != nil {
			return ChapterNarrative{}, err
		}
		purpose, behavior, decisions := splitNarrative(text)
		var problems []string
		if purpose == "" || behavior == "" || decisions == "" {
			problems = append(problems, "output must contain ## Summary, ## Observed Legacy Behavior, and ## Target Decisions sections")
		}
		if typedIDPattern.MatchString(purpose) {
			problems = append(problems, "the Summary must not contain record identifiers")
		}
		if strings.Contains(purpose, "CLM-") || strings.Contains(purpose, "[^") {
			problems = append(problems, "the Summary must not contain citations")
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
			numbered, index := numberCitations(narrative, claimNotes)
			full := numbered + index + recordAppendix(observed)
			return ChapterNarrative{Chapter: spec.Number, Markdown: full}, nil
		}
		rejection = strings.Join(problems, "\n")
		lastRejection = rejection
		for _, problem := range problems {
			if strings.Contains(problem, "HBK-SUBSTANCE") {
				rejection += "\nHint: split combined statements and add one cited sentence per distinct observed fact. Each finding usually supports several facts: its category, its affected scope, its risk, its safeguard gap, and its mitigation."
				break
			}
		}
	}
	return ChapterNarrative{}, fmt.Errorf("handbook: chapter %s did not pass validation: %s", chapterNumber, lastRejection)
}

var claimPattern = regexp.MustCompile(`CLM-[A-F0-9]{12}`)

var typedIDPattern = regexp.MustCompile(`[A-Z][A-Z0-9]{1,7}-[A-F0-9]{12}`)

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
	fmt.Fprintf(&b, "### Summary\n\n%s\n\n", purpose)
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
// claimNotes maps claim IDs to a readable note (claim text plus subject).
func loadClaimNotes(db string, doc *cdl.EIRDoc, evidence map[string]bool) (map[string]string, error) {
	notes := map[string]string{}
	if len(evidence) == 0 {
		return notes, nil
	}
	names, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	for claim := range evidence {
		payload, err := store.Get(db, doc, "CLAIM", claim)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		text := strings.Join(strings.Fields(values["claim"]), " ")
		subject := strings.Fields(values["subject-id-record-version"])
		note := text
		if len(subject) > 0 {
			note = fmt.Sprintf("%s - source %s", text, subject[0])
		}
		notes[claim] = note
	}
	return notes, nil
}

func numberCitations(narrative string, notes map[string]string) (string, string) {
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
		note := notes[claim]
		if note == "" {
			note = "evidence record"
		}
		fmt.Fprintf(&b, "[^%d]: %s (see [%s](%s#%s))\n", i+1, note, claim, EvidenceRegisterFile, ClaimAnchor(claim))
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
		case strings.HasPrefix(trimmed, "## Summary"):
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
