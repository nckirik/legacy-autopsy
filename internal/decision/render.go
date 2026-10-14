package decision

import (
	"fmt"
	"strings"
)

// DraftPath is the default candidate decision-draft file.
const DraftPath = "preview/decisions/DECISIONS-DRAFT.md"

// ReviewPath is the default recorded human decision-review file.
const ReviewPath = "preview/decisions/decisions-review.md"

const draftBegin = "<!-- DECISION-DRAFT:"
const draftEnd = "<!-- END-DECISION-DRAFT -->"

// RenderDraft renders the candidate decision drafts as a reviewable Markdown
// file with per-entry markers so the interactive review can parse it back.
func RenderDraft(drafts []Draft, generatedAt string) ([]byte, error) {
	var b strings.Builder
	b.WriteString("# Decision drafts (candidate content)\n\n")
	b.WriteString("> Generated candidate content for human decision review. This file is not canonical state; canonical decision records are created only after review.\n\n")
	fmt.Fprintf(&b, "- **Generated At:** %s\n", generatedAt)
	fmt.Fprintf(&b, "- **Drafts:** %d\n\n", len(drafts))
	for _, draft := range drafts {
		sections, err := DraftContent(draft.Content)
		if err != nil {
			return nil, fmt.Errorf("decision: draft %s: %w", draft.Target.SubjectID, err)
		}
		fmt.Fprintf(&b, "%s %s -->\n", draftBegin, draft.Target.SubjectID)
		fmt.Fprintf(&b, "### %s %s decision\n\n", draft.Target.SubjectID, draft.Target.DecisionType)
		fmt.Fprintf(&b, "- **Subject:** %s\n", draft.Target.SubjectLine)
		fmt.Fprintf(&b, "- **Decision Type:** %s\n", draft.Target.DecisionType)
		fmt.Fprintf(&b, "- **Observed Inputs:** %s\n", strings.Join(draft.Target.ObservedLines, "; "))
		b.WriteString("- **Draft Content:**\n\n")
		for _, want := range []string{"Options", "Decision", "Rationale", "Risks", "Exact Scope and Denominator Effect", "Affected Records / Rollout / Rollback"} {
			fmt.Fprintf(&b, "## %s\n\n%s\n\n", want, sections[want])
		}
		b.WriteString(draftEnd)
		b.WriteString("\n\n")
	}
	return []byte(b.String()), nil
}

// DraftEntry is one parsed decision draft from a rendered draft file.
type DraftEntry struct {
	SubjectID     string
	DecisionType  string
	SubjectLine   string
	ObservedLines []string
	Sections      map[string]string
}

// ParseDraft reads the rendered draft file back into bounded entries.
func ParseDraft(content []byte) ([]DraftEntry, error) {
	text := string(content)
	var entries []DraftEntry
	position := 0
	for {
		begin := strings.Index(text[position:], draftBegin)
		if begin < 0 {
			break
		}
		begin += position
		markerEnd := strings.Index(text[begin:], "-->")
		if markerEnd < 0 {
			return nil, fmt.Errorf("decision: unterminated draft marker")
		}
		markerEnd += begin
		subjectID := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(text[begin+len(draftBegin):markerEnd]), "-->"))
		end := strings.Index(text[markerEnd:], draftEnd)
		if end < 0 {
			return nil, fmt.Errorf("decision: unterminated draft entry %s", subjectID)
		}
		end += markerEnd
		block := text[markerEnd:end]
		entry := DraftEntry{SubjectID: subjectID, Sections: map[string]string{}}
		current := ""
		var body []string
		flush := func() {
			if current != "" {
				entry.Sections[current] = strings.TrimSpace(strings.Join(body, "\n"))
			}
			body = nil
		}
		for _, line := range strings.Split(block, "\n") {
			trimmed := strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(trimmed, "- **Subject:**"):
				entry.SubjectLine = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Subject:**"))
				continue
			case strings.HasPrefix(trimmed, "- **Decision Type:**"):
				entry.DecisionType = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Decision Type:**"))
				continue
			case strings.HasPrefix(trimmed, "- **Observed Inputs:**"):
				for _, item := range strings.Split(strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Observed Inputs:**")), ";") {
					item = strings.TrimSpace(item)
					if item != "" {
						entry.ObservedLines = append(entry.ObservedLines, item)
					}
				}
				continue
			case strings.HasPrefix(trimmed, "- **Draft Content:**"):
				continue
			case strings.HasPrefix(trimmed, "## "):
				flush()
				current = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
				continue
			}
			if current != "" {
				body = append(body, line)
			}
		}
		flush()
		for _, want := range []string{"Options", "Decision", "Rationale", "Risks", "Exact Scope and Denominator Effect", "Affected Records / Rollout / Rollback"} {
			if strings.TrimSpace(entry.Sections[want]) == "" {
				return nil, fmt.Errorf("decision: draft %s is missing the %q section", subjectID, want)
			}
		}
		if entry.SubjectLine == "" || entry.DecisionType == "" {
			return nil, fmt.Errorf("decision: draft %s is missing its subject or decision type", subjectID)
		}
		entries = append(entries, entry)
		position = end + len(draftEnd)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("decision: no decision drafts found")
	}
	return entries, nil
}
