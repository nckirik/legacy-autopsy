package equivalence

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Prompter asks the human reviewer one question and returns the answer.
type Prompter func(question string) (string, error)

// ReviewRow is one recorded per-record equivalence review decision.
type ReviewRow struct {
	RecordType string
	RecordID   string
	Verdict    string // Accepted | Excluded
	Reason     string
}

// SuiteReview is one recorded human review of the acceptance suite.
type SuiteReview struct {
	Reviewer    string
	Timestamp   string
	Scope       string
	Limitations string
	Result      string // Pass | Incomplete
	Rows        []ReviewRow
}

// CollectSuiteReview runs the interactive per-record review.
func CollectSuiteReview(prompt Prompter, out io.Writer, mappings []Mapping, review SuiteReview, timestamp string) (SuiteReview, error) {
	if len(mappings) == 0 {
		return review, fmt.Errorf("equivalence review: the suite declares no mappings")
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
		answer, err := prompt("Review scope [every mapped contract]: ")
		if err != nil {
			return review, err
		}
		review.Scope = strings.TrimSpace(answer)
		if review.Scope == "" {
			review.Scope = "every mapped contract"
		}
	}
	if strings.TrimSpace(review.Limitations) == "" {
		answer, err := prompt("Review limitations [document-level review of the skeletons]: ")
		if err != nil {
			return review, err
		}
		review.Limitations = strings.TrimSpace(answer)
		if review.Limitations == "" {
			review.Limitations = "document-level review of the skeletons"
		}
	}
	review.Timestamp = timestamp
	review.Result = "Pass"
	fmt.Fprintf(out, "\nFor each item: accept the skeleton as an adequate check, exclude it with a reason, or skip it.\n")
	for _, mapping := range mappings {
		fmt.Fprintf(out, "\n---- %s %s (%s) ----\n", mapping.RecordType, mapping.RecordID, mapping.Skeleton)
		if strings.TrimSpace(mapping.Body) != "" {
			fmt.Fprintf(out, "%s\n", plainBody(mapping.Body))
		}
		for {
			answer, err := prompt("[a]ccept / [e]xclude / [s]kip: ")
			if err != nil {
				return review, err
			}
			switch strings.ToLower(strings.TrimSpace(answer)) {
			case "a", "accept":
				review.Rows = append(review.Rows, ReviewRow{RecordType: mapping.RecordType, RecordID: mapping.RecordID, Verdict: "Accepted"})
			case "e", "exclude":
				reason, err := prompt("Exclusion reason (required): ")
				if err != nil {
					return review, err
				}
				reason = strings.TrimSpace(reason)
				if reason == "" {
					fmt.Fprintln(out, "An exclusion needs a reason. Record the reason or accept the skeleton.")
					continue
				}
				review.Rows = append(review.Rows, ReviewRow{RecordType: mapping.RecordType, RecordID: mapping.RecordID, Verdict: "Excluded", Reason: reason})
			case "s", "skip":
				review.Result = "Incomplete"
			default:
				fmt.Fprintln(out, "Answer a, e, or s.")
				continue
			}
			break
		}
	}
	return review, nil
}

// RenderSuiteReview serializes the review record.
func RenderSuiteReview(review SuiteReview) []byte {
	var b strings.Builder
	b.WriteString("# Acceptance Suite Review\n\n")
	fmt.Fprintf(&b, "- **Reviewer / Role / Authority Basis:** %s\n", review.Reviewer)
	fmt.Fprintf(&b, "- **Review Timestamp:** %s\n", review.Timestamp)
	fmt.Fprintf(&b, "- **Review Scope:** %s\n", review.Scope)
	fmt.Fprintf(&b, "- **Review Limitations:** %s\n", review.Limitations)
	fmt.Fprintf(&b, "- **Result:** %s\n\n", review.Result)
	b.WriteString("| Record Type | Typed Record ID | Verdict | Reason |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, row := range review.Rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", row.RecordType, row.RecordID, row.Verdict, cell(row.Reason))
	}
	return []byte(b.String())
}

// ParseSuiteReview reads a recorded review back.
func ParseSuiteReview(content []byte) (SuiteReview, error) {
	var review SuiteReview
	inTable := false
	header := true
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if value, ok := fieldValue(trimmed, "Reviewer / Role / Authority Basis"); ok {
			review.Reviewer = value
			continue
		}
		if value, ok := fieldValue(trimmed, "Review Timestamp"); ok {
			review.Timestamp = value
			continue
		}
		if value, ok := fieldValue(trimmed, "Review Scope"); ok {
			review.Scope = value
			continue
		}
		if value, ok := fieldValue(trimmed, "Review Limitations"); ok {
			review.Limitations = value
			continue
		}
		if value, ok := fieldValue(trimmed, "Result"); ok {
			review.Result = value
			continue
		}
		if strings.HasPrefix(trimmed, "| Record Type |") {
			inTable = true
			header = true
			continue
		}
		if inTable && !strings.HasPrefix(trimmed, "|") {
			inTable = false
			continue
		}
		if !inTable {
			continue
		}
		cells := splitCells(trimmed)
		if header {
			header = false
			continue
		}
		if len(cells) < 3 || strings.Trim(cells[0], "-: ") == "" {
			continue
		}
		review.Rows = append(review.Rows, ReviewRow{
			RecordType: cells[0],
			RecordID:   cells[1],
			Verdict:    cells[2],
			Reason:     optionalCell(cells, 3),
		})
	}
	if review.Reviewer == "" || review.Timestamp == "" || review.Result == "" {
		return review, fmt.Errorf("equivalence: review file is incomplete")
	}
	switch review.Result {
	case "Pass", "Incomplete":
	default:
		return review, fmt.Errorf("equivalence: review result %q is invalid", review.Result)
	}
	return review, nil
}

// LoadSuiteReview reads one review record; missing files report ok=false.
func LoadSuiteReview(path string) (SuiteReview, bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SuiteReview{}, false, nil
		}
		return SuiteReview{}, false, err
	}
	review, err := ParseSuiteReview(content)
	if err != nil {
		return SuiteReview{}, false, err
	}
	return review, true, nil
}

func fieldValue(line, label string) (string, bool) {
	prefix := "- **" + label + ":**"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
}

func splitCells(line string) []string {
	line = strings.Trim(strings.TrimSpace(line), "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func optionalCell(cells []string, index int) string {
	if index < len(cells) {
		return strings.TrimSpace(cells[index])
	}
	return ""
}

func cell(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// plainBody renders a Markdown skeleton body for terminal reading: emphasis
// markers and backticks are removed, indentation is preserved.
func plainBody(body string) string {
	body = strings.ReplaceAll(body, "**", "")
	body = strings.ReplaceAll(body, "`", "")
	return body
}
