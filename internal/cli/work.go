package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specWork runs bounded semantic work units through an OpenAI-compatible model
// adapter. Secrets come from environment variables only and are never logged.
func specWork(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec work <ping|run> [flags]")
	}
	op := args[0]
	switch op {
	case "ping", "run":
	default:
		return fmt.Errorf("spec work operation %q is unsupported; supported operations: ping, run", op)
	}
	set := flags("spec work "+op, stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (run)")
	unitsPath := set.String("units", "", "JSON work-unit list (run)")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	timeout := set.Duration("llm-timeout", 0, "per-request timeout (default 120s)")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
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
	if *timeout > 0 {
		config.Timeout = *timeout
	}
	if err := config.Validate(); err != nil {
		return err
	}
	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}

	if op == "ping" {
		start := time.Now()
		reply, err := client.Complete(context.Background(), "You are a connectivity probe. Reply with exactly OK.", "Reply with exactly OK.")
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "endpoint: %s/chat/completions\n", strings.TrimRight(config.BaseURL, "/"))
		fmt.Fprintf(stdout, "model:    %s\n", config.Model)
		fmt.Fprintf(stdout, "latency:  %s\n", time.Since(start).Round(time.Millisecond))
		fmt.Fprintf(stdout, "reply:    %s\n", strings.TrimSpace(reply))
		return nil
	}

	if *db == "" || *unitsPath == "" {
		return fmt.Errorf("spec work run requires --db and --units")
	}
	unitsBytes, err := os.ReadFile(*unitsPath)
	if err != nil {
		return err
	}
	var units []work.Unit
	if err := json.Unmarshal(unitsBytes, &units); err != nil {
		return fmt.Errorf("parse units %s: %w", *unitsPath, err)
	}
	if len(units) == 0 {
		return fmt.Errorf("units file %s declares no work items", *unitsPath)
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	runConfig := work.DefaultConfig()
	if *reviewMode != "" {
		runConfig.ReviewMode = work.ReviewMode(*reviewMode)
	}
	if err := runConfig.Validate(); err != nil {
		return err
	}
	session, err := work.NewSession(*db, res.EIR, runConfig, llm.NewExecutor(client), llm.NewReviewer(client))
	if err != nil {
		return err
	}
	report, runErr := session.Run(units)
	fmt.Fprintf(stdout, "%-20s %8s %5s %6s %7s %-12s %s\n", "unit", "context", "deps", "fetch", "retries", "review", "escalated")
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-20s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "run stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "committed %d unit(s)\n", len(report.Units))
	return nil
}
