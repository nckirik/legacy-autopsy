package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/llm"
)

// specLLM probes the configured endpoint with a real prompt and reports
// time-to-first-token, total time, token usage, and throughput so reasoning
// behavior can be witnessed directly.
func specLLM(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "probe" {
		return fmt.Errorf("usage: legacy-autopsy spec llm probe [flags]")
	}
	set := flags("spec llm probe", stderr)
	systemFile := set.String("system-file", "", "system prompt file")
	userFile := set.String("user-file", "", "user prompt file (required)")
	effort := set.String("effort", "", "reasoning effort override (low|medium|high|xhigh|off)")
	maxTokens := set.Int("max-tokens", 0, "max completion tokens override")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if *userFile == "" {
		return fmt.Errorf("spec llm probe requires --user-file")
	}
	config, err := llm.ConfigFromEnv()
	if err != nil {
		return err
	}
	if *maxTokens > 0 {
		config.MaxTokens = *maxTokens
	}
	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}
	system := "You are a terse assistant."
	if *systemFile != "" {
		data, err := os.ReadFile(*systemFile)
		if err != nil {
			return err
		}
		system = string(data)
	}
	user, err := os.ReadFile(*userFile)
	if err != nil {
		return err
	}
	resolvedEffort := config.Reasoning
	if *effort != "" {
		resolvedEffort = *effort
	}
	if strings.EqualFold(resolvedEffort, "off") || strings.EqualFold(resolvedEffort, "none") {
		resolvedEffort = ""
	}
	fmt.Fprintf(stdout, "model=%s effort=%q prompt_chars=%d max_tokens=%d\n", config.Model, resolvedEffort, len(user), config.MaxTokens)
	var usage *llm.Usage
	var firstDelta time.Duration
	start := time.Now()
	text, err := client.CompleteEffortStream(context.Background(), system, string(user), resolvedEffort, func(delta llm.Delta) {
		if firstDelta == 0 && (delta.Content != "" || delta.Reasoning != "") {
			firstDelta = time.Since(start)
		}
		if delta.Content != "" {
			fmt.Fprint(stdout, delta.Content)
		}
		if delta.Usage != nil {
			usage = delta.Usage
		}
	})
	total := time.Since(start)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "\n---\nfirst_delta=%s total=%s content_chars=%d\n", firstDelta.Round(time.Millisecond), total.Round(time.Millisecond), len(text))
	if usage != nil {
		seconds := total.Seconds()
		fmt.Fprintf(stdout, "usage: prompt=%d completion=%d reasoning=%d\n", usage.PromptTokens, usage.CompletionTokens, usage.ReasoningTokens)
		if seconds > 0 {
			fmt.Fprintf(stdout, "throughput: %.1f tok/s (completion/seconds)\n", float64(usage.CompletionTokens)/seconds)
		}
	} else {
		fmt.Fprintln(stdout, "usage: not reported by provider")
	}
	return nil
}
