package decision

import (
	"fmt"
	"io"
	"strings"
)

// Prompter asks the human reviewer one question and returns the answer.
type Prompter func(question string) (string, error)

// ReviewRow is one recorded human verdict for a decision draft.
type ReviewRow struct {
	SubjectID string
	Type      string
	Verdict   string // Approve | Reject | Defer
	Reason    string
}

// Review is one recorded human review of the decision drafts.
type Review struct {
	Reviewer    string
	Timestamp   string
	Scope       string
	Limitations string
	Result      string // Pass | Incomplete
	Rows        []ReviewRow
}

// CollectReview runs the interactive per-decision review. Approving a draft
// authorizes the reviewed commit path to create the canonical decision record;
// rejecting or deferring leaves it uncreated.
func CollectReview(prompt Prompter, out io.Writer, entries []DraftEntry, review Review, timestamp string) (Review, error) {
	if len(entries) == 0 {
		return review, fmt.Errorf("decision review: no decision drafts to review")
	}
	if strings.TrimSpace(review.Reviewer) == "" {
		for {
			answer, err := prompt("Reviewer name / role: ")
			if err != nil {
				return review, err
			}
			if strings.TrimSpace(answer) != "" {
				review.Reviewer = strings.TrimSpace(answer)
				break
			}
			fmt.Fprintln(out, "A reviewer name or role is required.")
		}
	}
	if strings.TrimSpace(review.Scope) == "" {
		answer, err := prompt("Review scope [every drafted decision]: ")
		if err != nil {
			return review, err
		}
		review.Scope = strings.TrimSpace(answer)
		if review.Scope == "" {
			review.Scope = "every drafted decision"
		}
	}
	if strings.TrimSpace(review.Limitations) == "" {
		answer, err := prompt("Review limitations [content review of the draft decisions]: ")
		if err != nil {
			return review, err
		}
		review.Limitations = strings.TrimSpace(answer)
		if review.Limitations == "" {
			review.Limitations = "content review of the draft decisions"
		}
	}
	review.Timestamp = timestamp
	review.Result = "Pass"
	fmt.Fprintf(out, "\nFor each decision: [a]pprove, [r]eject, or [d]efer.\n")
	for _, entry := range entries {
		fmt.Fprintf(out, "\n---- %s %s decision ----\n", entry.SubjectID, entry.DecisionType)
		fmt.Fprintf(out, "Subject: %s\n", entry.SubjectLine)
		for _, want := range []string{"Options", "Decision", "Rationale", "Risks", "Exact Scope and Denominator Effect", "Affected Records / Rollout / Rollback"} {
			fmt.Fprintf(out, "\n%s:\n%s\n", want, entry.Sections[want])
		}
		row := ReviewRow{SubjectID: entry.SubjectID, Type: entry.DecisionType}
		for {
			answer, err := prompt("[a]pprove / [r]eject / [d]efer: ")
			if err != nil {
				return review, err
			}
			switch strings.ToLower(strings.TrimSpace(answer)) {
			case "a", "approve":
				row.Verdict = "Approve"
			case "r", "reject":
				row.Verdict = "Reject"
			case "d", "defer":
				row.Verdict = "Defer"
			default:
				fmt.Fprintln(out, "Answer a, r, or d.")
				continue
			}
			break
		}
		answer, err := prompt("Reason / conditions: ")
		if err != nil {
			return review, err
		}
		row.Reason = strings.TrimSpace(answer)
		if row.Reason == "" {
			row.Reason = "None"
		}
		if row.Verdict != "Approve" {
			review.Result = "Incomplete"
		}
		review.Rows = append(review.Rows, row)
	}
	return review, nil
}

// RenderReview renders the recorded review as Markdown.
func RenderReview(review Review) []byte {
	var b strings.Builder
	b.WriteString("# Decision review\n\n")
	fmt.Fprintf(&b, "- **Reviewer:** %s\n", review.Reviewer)
	fmt.Fprintf(&b, "- **Timestamp:** %s\n", review.Timestamp)
	fmt.Fprintf(&b, "- **Scope:** %s\n", review.Scope)
	fmt.Fprintf(&b, "- **Limitations:** %s\n", review.Limitations)
	fmt.Fprintf(&b, "- **Result:** %s\n\n", review.Result)
	b.WriteString("| Subject | Type | Verdict | Reason |\n| :-- | :-- | :-- | :-- |\n")
	for _, row := range review.Rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", row.SubjectID, row.Type, row.Verdict, sanitizeCell(row.Reason))
	}
	return []byte(b.String())
}

// ParseReview reads a recorded decision review.
func ParseReview(content []byte) (Review, error) {
	var review Review
	var inTable bool
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "- **Reviewer:**"):
			review.Reviewer = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Reviewer:**"))
		case strings.HasPrefix(trimmed, "- **Timestamp:**"):
			review.Timestamp = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Timestamp:**"))
		case strings.HasPrefix(trimmed, "- **Scope:**"):
			review.Scope = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Scope:**"))
		case strings.HasPrefix(trimmed, "- **Limitations:**"):
			review.Limitations = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Limitations:**"))
		case strings.HasPrefix(trimmed, "- **Result:**"):
			review.Result = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Result:**"))
		case strings.HasPrefix(trimmed, "| Subject |"):
			inTable = true
		case inTable && strings.HasPrefix(trimmed, "| :--"):
			continue
		case inTable && strings.HasPrefix(trimmed, "|"):
			cells := splitRow(trimmed)
			if len(cells) != 4 {
				return review, fmt.Errorf("decision review: malformed verdict row %q", trimmed)
			}
			review.Rows = append(review.Rows, ReviewRow{SubjectID: cells[0], Type: cells[1], Verdict: cells[2], Reason: cells[3]})
		}
	}
	if review.Reviewer == "" || review.Timestamp == "" || len(review.Rows) == 0 {
		return review, fmt.Errorf("decision review: review record is incomplete")
	}
	return review, nil
}

func splitRow(line string) []string {
	line = strings.TrimSpace(strings.TrimPrefix(line, "|"))
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func sanitizeCell(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "|", "/"), "\n", " ")
}
