package handbook

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ChapterPrompt is the bounded review material for one chapter.
type ChapterPrompt struct {
	Path        string
	HBKID       string
	Fingerprint string
	Title       string
	Summary     string
}

// Verdict is one recorded chapter readability verdict.
type Verdict struct {
	Chapter     string
	HBKID       string
	Fingerprint string
	Verdict     string // Pass | Fail
	Notes       string
}

// ReadabilityReviewFile is one combined human readability review record.
type ReadabilityReviewFile struct {
	Reviewer    string
	Timestamp   string
	Scope       string
	Limitations string
	Result      string
	Verdicts    []Verdict
}

// Prompter asks the human reviewer one question and returns the answer.
type Prompter func(question string) (string, error)

// CollectReadabilityReview runs the interactive per-chapter review. Fail and
// Pass verdicts are recorded; skipped chapters stay uncovered and fail closed.
func CollectReadabilityReview(prompt Prompter, out io.Writer, chapters []ChapterPrompt, review ReadabilityReviewFile, timestamp string) (ReadabilityReviewFile, error) {
	if len(chapters) == 0 {
		return review, fmt.Errorf("readability review: no chapters to review")
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
		answer, err := prompt("Review scope [all candidates chapters]: ")
		if err != nil {
			return review, err
		}
		review.Scope = strings.TrimSpace(answer)
		if review.Scope == "" {
			review.Scope = "all candidate handbook chapters"
		}
	}
	if strings.TrimSpace(review.Limitations) == "" {
		answer, err := prompt("Review limitations [static read of the rendered chapters only]: ")
		if err != nil {
			return review, err
		}
		review.Limitations = strings.TrimSpace(answer)
		if review.Limitations == "" {
			review.Limitations = "static read of the rendered chapters only"
		}
	}
	review.Timestamp = timestamp
	review.Result = "Pass"
	fmt.Fprintln(out, "\nRead each chapter file listed below, then answer [p]ass, [f]ail (with a reason), or [s]kip.")
	for _, chapter := range chapters {
		fmt.Fprintf(out, "\n%s  %s  %s\n%s\n", chapter.Path, chapter.HBKID, chapter.Title, chapter.Summary)
		for {
			answer, err := prompt("[p]ass / [f]ail / [s]kip: ")
			if err != nil {
				return review, err
			}
			switch strings.ToLower(strings.TrimSpace(answer)) {
			case "p", "pass":
				review.Verdicts = append(review.Verdicts, Verdict{Chapter: chapter.Path, HBKID: chapter.HBKID, Fingerprint: chapter.Fingerprint, Verdict: "Pass"})
			case "f", "fail":
				notes, err := prompt("Failure notes (required): ")
				if err != nil {
					return review, err
				}
				notes = strings.TrimSpace(notes)
				if notes == "" {
					fmt.Fprintln(out, "A failure needs a reason. Record the reason or pass the chapter.")
					continue
				}
				review.Verdicts = append(review.Verdicts, Verdict{Chapter: chapter.Path, HBKID: chapter.HBKID, Fingerprint: chapter.Fingerprint, Verdict: "Fail", Notes: notes})
				review.Result = "Fail"
			case "s", "skip":
			default:
				fmt.Fprintln(out, "Answer p, f, or s.")
				continue
			}
			break
		}
	}
	if len(review.Verdicts) == 0 {
		return review, fmt.Errorf("readability review: no verdicts were recorded")
	}
	return review, nil
}

// RenderReadabilityReview serializes the review record.
func RenderReadabilityReview(review ReadabilityReviewFile) []byte {
	var b strings.Builder
	b.WriteString("# Handbook Readability Review\n\n")
	fmt.Fprintf(&b, "- **Reviewer / Role / Authority Basis:** %s\n", review.Reviewer)
	fmt.Fprintf(&b, "- **Review Timestamp:** %s\n", review.Timestamp)
	fmt.Fprintf(&b, "- **Review Scope:** %s\n", review.Scope)
	fmt.Fprintf(&b, "- **Review Limitations:** %s\n", review.Limitations)
	fmt.Fprintf(&b, "- **Result:** %s\n\n", review.Result)
	b.WriteString("| Chapter | HBK ID | Semantic Fingerprint | Verdict | Notes |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, verdict := range review.Verdicts {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", verdict.Chapter, verdict.HBKID, verdict.Fingerprint, verdict.Verdict, cell(verdict.Notes))
	}
	return []byte(b.String())
}

// ParseReadabilityReview reads a recorded review back.
func ParseReadabilityReview(content []byte) (ReadabilityReviewFile, error) {
	var review ReadabilityReviewFile
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
		if strings.HasPrefix(trimmed, "| Chapter |") {
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
		if len(cells) < 4 || strings.Trim(cells[0], "-: ") == "" {
			continue
		}
		review.Verdicts = append(review.Verdicts, Verdict{
			Chapter:     cells[0],
			HBKID:       cells[1],
			Fingerprint: cells[2],
			Verdict:     cells[3],
			Notes:       optionalCell(cells, 4),
		})
	}
	if review.Reviewer == "" || review.Timestamp == "" || review.Result == "" {
		return review, fmt.Errorf("handbook: readability review file is incomplete")
	}
	if review.Result != "Pass" && review.Result != "Fail" {
		return review, fmt.Errorf("handbook: readability review result %q is invalid", review.Result)
	}
	return review, nil
}

// EnvelopeStatus reports the certification envelope status of one minted
// section.
func EnvelopeStatus(section []byte) (status, confirmation string) {
	content := string(section)
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **HANDBOOK-STATUS:**") {
			status = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **HANDBOOK-STATUS:**"))
		}
		if strings.HasPrefix(trimmed, "- **Confirmation:**") {
			confirmation = strings.TrimSpace(strings.TrimPrefix(trimmed, "- **Confirmation:**"))
		}
	}
	return status, confirmation
}

// ReviewsFromFile converts one recorded file into the PI-12 review model.
// Fail verdicts still cover their chapter; the failing result is the gate.
func ReviewsFromFile(review ReadabilityReviewFile) []Review {
	var bindings []Binding
	for _, verdict := range review.Verdicts {
		bindings = append(bindings, Binding{HBKID: verdict.HBKID, Fingerprint: verdict.Fingerprint})
	}
	value := review.Result
	if value == "" {
		value = "Fail"
	}
	return []Review{{
		ReadabilityReview: ReadabilityReview{
			Value:       value,
			Scope:       review.Scope,
			Limitations: review.Limitations,
			Timestamp:   review.Timestamp,
		},
		Bindings: bindings,
	}}
}

// LoadReadabilityReviews reads every readability-review.md file in a directory.
func LoadReadabilityReviews(dir string) ([]Review, []ReadabilityReviewFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "readability-") || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	var reviews []Review
	var files []ReadabilityReviewFile
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, nil, err
		}
		file, err := ParseReadabilityReview(content)
		if err != nil {
			return nil, nil, err
		}
		files = append(files, file)
		reviews = append(reviews, ReviewsFromFile(file)...)
	}
	return reviews, files, nil
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
