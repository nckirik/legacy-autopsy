package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/equivalence"
	"github.com/nckirik/legacy-autopsy/internal/handbook"
)

// specHandoffReview dispatches the interactive human review acts.
func specHandoffReview(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec handoff review <readability|equivalence> [flags]")
	}
	switch args[0] {
	case "readability":
		return specHandoffReviewReadability(args[1:], stdout, stderr)
	case "equivalence":
		return specHandoffReviewEquivalence(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec handoff review operation %q is unsupported; supported operations: readability, equivalence", args[0])
	}
}

// specHandoffReviewReadability records the human per-chapter readability review.
func specHandoffReviewReadability(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff review readability", stderr)
	dir := set.String("dir", "preview", "clean human handbook view directory")
	handbookDir := set.String("handbook-dir", "preview/package/handbook/sections", "minted HBK section payload directory")
	out := set.String("out", "preview/readability-review.md", "review record output file")
	reviewer := set.String("reviewer", "", "reviewer name / role (prompted when empty)")
	scope := set.String("scope", "", "review scope (prompted when empty)")
	limitations := set.String("limitations", "", "review limitations (prompted when empty)")
	answers := set.String("answers", "", "scripted answers file (non-interactive)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	sections, err := handbook.LoadMintedSections(*handbookDir)
	if err != nil {
		return err
	}
	if len(sections) == 0 {
		return fmt.Errorf("handoff review readability: no minted HBK sections in %s", *handbookDir)
	}
	var chapters []handbook.ChapterPrompt
	for _, section := range sections {
		view, err := os.ReadFile(filepath.Join(*dir, filepath.FromSlash(section.Path)))
		if err != nil {
			return err
		}
		if string(view) != handbook.SectionContent(section.Bytes) {
			return fmt.Errorf("handoff review readability: %s does not match its bound HBK payload", section.Path)
		}
		chapters = append(chapters, handbook.ChapterPrompt{
			Path:        section.Path,
			HBKID:       section.ID,
			Fingerprint: section.Fingerprint,
			Title:       section.Title,
			Summary:     firstSummary(view),
		})
	}
	prompt, closeFn, err := reviewPrompter(*answers, stdout)
	if err != nil {
		return err
	}
	if closeFn != nil {
		defer closeFn()
	}
	review := handbook.ReadabilityReviewFile{Reviewer: *reviewer, Scope: *scope, Limitations: *limitations}
	review, err = handbook.CollectReadabilityReview(handbook.Prompter(prompt), stdout, chapters, review, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, handbook.RenderReadabilityReview(review), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "\nreadability review: %s (%d verdicts, result %s)\n", *out, len(review.Verdicts), review.Result)
	return nil
}

// specHandoffReviewEquivalence records the human per-record suite review.
func specHandoffReviewEquivalence(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff review equivalence", stderr)
	suite := set.String("suite", "preview/package/equivalence/acceptance-suite.md", "generated acceptance-suite file")
	out := set.String("out", "", "review record output file (default beside the suite)")
	reviewer := set.String("reviewer", "", "reviewer name / role (prompted when empty)")
	scope := set.String("scope", "", "review scope (prompted when empty)")
	limitations := set.String("limitations", "", "review limitations (prompted when empty)")
	answers := set.String("answers", "", "scripted answers file (non-interactive)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	content, err := os.ReadFile(*suite)
	if err != nil {
		return err
	}
	mappings := equivalence.ParseMappings(content)
	if len(mappings) == 0 {
		return fmt.Errorf("handoff review equivalence: %s declares no mappings", *suite)
	}
	outPath := *out
	if outPath == "" {
		outPath = filepath.Join(filepath.Dir(*suite), "equivalence-review.md")
	}
	prompt, closeFn, err := reviewPrompter(*answers, stdout)
	if err != nil {
		return err
	}
	if closeFn != nil {
		defer closeFn()
	}
	review := equivalence.SuiteReview{Reviewer: *reviewer, Scope: *scope, Limitations: *limitations}
	review, err = equivalence.CollectSuiteReview(equivalence.Prompter(prompt), stdout, mappings, review, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, equivalence.RenderSuiteReview(review), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "\nequivalence review: %s (%d decisions, result %s)\n", outPath, len(review.Rows), review.Result)
	return nil
}

// reviewPrompter reads scripted answers when provided, otherwise prompts on the
// terminal.
func reviewPrompter(answersPath string, stdout io.Writer) (func(string) (string, error), func(), error) {
	if answersPath != "" {
		content, err := os.ReadFile(answersPath)
		if err != nil {
			return nil, nil, err
		}
		var answers []string
		for _, line := range strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "#") {
				continue
			}
			answers = append(answers, strings.TrimSpace(line))
		}
		index := 0
		prompt := func(question string) (string, error) {
			if index >= len(answers) {
				return "", fmt.Errorf("review answers are exhausted at: %s", question)
			}
			answer := answers[index]
			index++
			fmt.Fprintf(stdout, "%s%s\n", question, answer)
			return answer, nil
		}
		return prompt, nil, nil
	}
	reader := bufio.NewReader(os.Stdin)
	prompt := func(question string) (string, error) {
		fmt.Fprint(stdout, question)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF && strings.TrimSpace(line) == "" {
			return "", fmt.Errorf("standard input ended before the review completed")
		}
		return strings.TrimSpace(line), nil
	}
	return prompt, nil, nil
}

// firstSummary extracts the opening paragraph of the chapter Summary section.
func firstSummary(view []byte) string {
	lines := strings.Split(string(view), "\n")
	inSummary := false
	var paragraph []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### Summary") {
			inSummary = true
			continue
		}
		if !inSummary {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			break
		}
		if trimmed == "" {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		paragraph = append(paragraph, trimmed)
	}
	return strings.Join(paragraph, " ")
}
