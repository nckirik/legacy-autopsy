package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/handbook"
	"github.com/nckirik/legacy-autopsy/internal/llm"
)

// specHandbook renders the human handbook draft from committed synthesis records.
func specHandbook(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec handbook <draft|narrate> [flags]")
	}
	if args[0] == "narrate" {
		return specHandbookNarrate(args[1:], stdout, stderr)
	}
	if args[0] != "draft" {
		return fmt.Errorf("spec handbook operation %q is unsupported; supported operations: draft, narrate", args[0])
	}
	set := flags("spec handbook draft", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	out := set.String("out", "", "write the draft to this path (default stdout)")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handbook draft requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	draft, findings, err := handbook.Draft(*db, res.EIR)
	if err != nil {
		return err
	}
	if *out != "" {
		if err := os.WriteFile(*out, draft, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "handbook draft: %s\n", *out)
	} else {
		stdout.Write(draft)
	}
	fmt.Fprintf(stdout, "lint findings: %d\n", len(findings))
	for _, finding := range findings {
		fmt.Fprintf(stdout, "  %s\n", finding.String())
	}
	return nil
}

// specHandbookNarrate authors one chapter's narrative with STE lint feedback.
func specHandbookNarrate(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handbook narrate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	chapter := set.String("chapter", "", "chapter number (for example 01)")
	out := set.String("out", "", "write the chapter to this path (default stdout)")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	llmReasoning := set.String("llm-reasoning", "", "override "+llm.EnvReasoning)
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" || *chapter == "" {
		return fmt.Errorf("spec handbook narrate requires --db and --chapter")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	config, err := llm.ConfigFromEnv()
	if err != nil {
		return err
	}
	if *llmURL != "" {
		config.BaseURL = *llmURL
	}
	if *llmModel != "" {
		config.Model = *llmModel
	}
	effort := config.Reasoning
	if *llmReasoning != "" {
		effort = *llmReasoning
	}
	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}
	narrative, err := handbook.Narrate(context.Background(), &handbookCompleter{client: client, effort: effort}, *db, res.EIR, *chapter)
	if err != nil {
		return err
	}
	if *out != "" {
		if err := os.WriteFile(*out, []byte(narrative.Markdown), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "chapter %s: %s\n", *chapter, *out)
		return nil
	}
	_, err = io.WriteString(stdout, narrative.Markdown)
	return err
}

type handbookCompleter struct {
	client *llm.Client
	effort string
}

func (h *handbookCompleter) Complete(ctx context.Context, system, user string) (string, error) {
	return h.client.CompleteEffort(ctx, system, user, h.effort)
}
