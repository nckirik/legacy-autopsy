package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/preflight"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specPreflight scans a target repository and runs the Preflight pipeline.
func specPreflight(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec preflight <scan|run> [flags]")
	}
	op := args[0]
	switch op {
	case "scan":
		return specPreflightScan(args[1:], stdout, stderr)
	case "run":
		return specPreflightRun(args[1:], stdout, stderr)
	case "matrix":
		return specPreflightMatrix(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec preflight operation %q is unsupported; supported operations: scan, run, matrix", op)
	}
}

// specPreflightScan prints the bounded deterministic repository scan.
func specPreflightScan(args []string, stdout, stderr io.Writer) error {
	set := flags("spec preflight scan", stderr)
	target := set.String("target", "", "target legacy repository (required)")
	maxFiles := set.Int("max-files", 5000, "bounded scan file count")
	maxDepth := set.Int("max-depth", 6, "bounded scan depth")
	maxEntries := set.Int("max-entry-points", 50, "bounded entry-point list")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *target == "" {
		return fmt.Errorf("spec preflight scan requires --target")
	}
	scan, err := preflight.ScanRepo(*target, *maxFiles, *maxEntries, *maxDepth)
	if err != nil {
		return err
	}
	stdout.Write(scan.Markdown())
	fmt.Fprintf(stdout, "\nfingerprint: %s\n", scan.Fingerprint)
	return nil
}

// specPreflightMatrix mechanically expands required traversal-matrix rows.
func specPreflightMatrix(args []string, stdout, stderr io.Writer) error {
	set := flags("spec preflight matrix", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	environment := set.String("environment", "", "environment label (default: preflight header)")
	snapshot := set.String("snapshot", "", "snapshot label (default: preflight header)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec preflight matrix requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	env := *environment
	if env == "" {
		env = header["included-environments-snapshots"]
	}
	snap := *snapshot
	if snap == "" {
		snap = header["included-environments-snapshots"]
	}
	created, err := preflight.ExpandMatrix(*db, res.EIR, preflight.MatrixDirectives{Environment: env, Snapshot: snap})
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "matrix rows created: %d\n", created)
	return nil
}

func specPreflightRun(args []string, stdout, stderr io.Writer) error {
	set := flags("spec preflight run", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	target := set.String("target", "", "target legacy repository (required)")
	db := set.String("db", "", "SQLite database path (run)")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxFiles := set.Int("max-files", 5000, "bounded scan file count")
	maxDepth := set.Int("max-depth", 6, "bounded scan depth")
	maxEntries := set.Int("max-entry-points", 50, "bounded entry-point list")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	llmReasoning := set.String("llm-reasoning", "", "override "+llm.EnvReasoning)
	timeout := set.Duration("llm-timeout", 0, "per-request timeout (default 120s)")
	showModel := set.String("show-model", "off", "stream raw model output: off|stderr|stdout (diagnostic, non-authoritative)")
	maxRetries := set.Int("max-retries", 0, "override the retry budget for rejected proposals")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *target == "" {
		return fmt.Errorf("spec preflight run requires --target")
	}
	scan, err := preflight.ScanRepo(*target, *maxFiles, *maxEntries, *maxDepth)
	if err != nil {
		return err
	}
	if *db == "" {
		return fmt.Errorf("spec preflight run requires --db")
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
	if *llmReasoning != "" {
		config.Reasoning = *llmReasoning
	}
	if *timeout > 0 {
		config.Timeout = *timeout
	}
	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	runConfig := work.DefaultConfig()
	if *reviewMode != "" {
		runConfig.ReviewMode = work.ReviewMode(*reviewMode)
	}
	if *maxRetries > 0 {
		runConfig.MaxRetries = *maxRetries
	}
	if err := runConfig.Validate(); err != nil {
		return err
	}
	executor := llm.NewExecutor(client)
	reviewer := llm.NewReviewer(client)
	if *llmReasoning != "" {
		executor.ForceReasoning = llmReasoning
	}
	switch strings.ToLower(*showModel) {
	case "", "off", "none":
		wireActivity(&runConfig, executor, stdout)
	case "stderr":
		executor.SetDebug(stderr)
		reviewer.Debug = stderr
	case "stdout":
		executor.SetDebug(stdout)
		reviewer.Debug = stdout
	default:
		return fmt.Errorf("--show-model must be off, stderr, or stdout")
	}
	session, err := work.NewSession(*db, res.EIR, runConfig, executor, reviewer)
	if err != nil {
		return err
	}
	directives := preflight.DefaultDirectives()
	runConfig.Iteration = directives.Iteration
	units := preflight.Plan(scan, directives)
	report, runErr := session.Run(units)
	fmt.Fprintf(stdout, "%-24s %8s %5s %6s %7s %-12s %s\n", "unit", "context", "deps", "fetch", "retries", "review", "escalated")
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-24s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "preflight stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "preflight committed %d unit(s)\n", len(report.Units))
	return nil
}
