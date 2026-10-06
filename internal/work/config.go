package work

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ReviewMode selects when independent review is required.
type ReviewMode string

const (
	ReviewAll         ReviewMode = "all"
	ReviewCritical    ReviewMode = "critical"
	ReviewUncertainty ReviewMode = "uncertainty"
)

// OperatorMinimums are the operator-presentation minimums asked before a run.
type OperatorMinimums struct {
	DecisionBatchMinutes   int
	PlainLanguageQuestions bool
}

// Config is the per-run runtime configuration asked before the run starts.
type Config struct {
	ReviewMode       ReviewMode
	OperatorMinimums OperatorMinimums
	MaxRetries       int
	// Observer receives lifecycle events for progress surfaces; it never
	// affects execution outcome.
	Observer *Observer
}

// DefaultConfig is the protocol default: independent review for units that
// declare themselves critical, plain-language questions, ten minutes per
// fifty-decision batch.
func DefaultConfig() Config {
	return Config{
		ReviewMode:       ReviewCritical,
		OperatorMinimums: OperatorMinimums{DecisionBatchMinutes: 10, PlainLanguageQuestions: true},
		MaxRetries:       3,
	}
}

// Validate checks the closed review-mode domain and the operator minimums.
func (c Config) Validate() error {
	switch c.ReviewMode {
	case ReviewAll, ReviewCritical, ReviewUncertainty:
	default:
		return fmt.Errorf("work: unknown review mode %q", c.ReviewMode)
	}
	if c.OperatorMinimums.DecisionBatchMinutes <= 0 {
		return fmt.Errorf("work: decision batch minutes must be positive")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("work: max retries must not be negative")
	}
	return nil
}

// NeedsReview applies the configured policy. uncertainty is set when the worker
// or a prior review was inconclusive.
func (c Config) NeedsReview(unit Unit, uncertainty bool) bool {
	switch c.ReviewMode {
	case ReviewAll:
		return true
	case ReviewUncertainty:
		return uncertainty
	default:
		return unit.RequiresReview
	}
}

// ConfigQuestion is one plain-language pre-run question with a default.
type ConfigQuestion struct {
	Prompt  string
	Options []string
	Default string
}

// RunQuestions are asked before a run starts; answers override defaults.
func RunQuestions() []ConfigQuestion {
	return []ConfigQuestion{
		{
			Prompt:  "Independent review policy",
			Options: []string{string(ReviewAll), string(ReviewCritical), string(ReviewUncertainty)},
			Default: string(ReviewCritical),
		},
		{
			Prompt:  "Target minutes for a 50-decision review batch",
			Default: "10",
		},
		{
			Prompt:  "Plain-language questions and decision indexes",
			Options: []string{"yes", "no"},
			Default: "yes",
		},
	}
}

// AskConfig asks the pre-run questions and returns the resulting configuration.
// Empty answers take the default.
func AskConfig(ask func(ConfigQuestion, string) (string, error)) (Config, error) {
	config := DefaultConfig()
	answers := make([]string, 0, 3)
	for _, question := range RunQuestions() {
		answer, err := ask(question, question.Default)
		if err != nil {
			return Config{}, err
		}
		answers = append(answers, strings.TrimSpace(answer))
	}
	if answers[0] != "" {
		config.ReviewMode = ReviewMode(answers[0])
	}
	if answers[1] != "" {
		minutes, err := strconv.Atoi(answers[1])
		if err != nil {
			return Config{}, fmt.Errorf("work: decision batch minutes %q is not a number", answers[1])
		}
		config.OperatorMinimums.DecisionBatchMinutes = minutes
	}
	if answers[2] != "" {
		config.OperatorMinimums.PlainLanguageQuestions = strings.EqualFold(answers[2], "yes")
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

// AskConfigInteractive reads answers line by line from in, writing prompts to out.
func AskConfigInteractive(in io.Reader, out io.Writer) (Config, error) {
	reader := bufio.NewReader(in)
	return AskConfig(func(question ConfigQuestion, fallback string) (string, error) {
		fmt.Fprintf(out, "%s", question.Prompt)
		if len(question.Options) > 0 {
			fmt.Fprintf(out, " [%s]", strings.Join(question.Options, "/"))
		}
		fmt.Fprintf(out, " (default %s): ", fallback)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		answer := strings.TrimSpace(line)
		if answer == "" {
			return fallback, nil
		}
		return answer, nil
	})
}
