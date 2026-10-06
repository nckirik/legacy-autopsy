package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specSynthesis runs the Synthesis pipeline for one catalog.
func specSynthesis(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec synthesis <run> [flags]")
	}
	if args[0] != "run" {
		return fmt.Errorf("spec synthesis operation %q is unsupported; supported operations: run", args[0])
	}
	set := flags("spec synthesis run", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	catalogID := set.String("catalog", "MOD", "synthesis catalog ID (MOD, ER, REL, SM, DR, BR, UC, IF, DEP, CFG, SCHED, NFR, SEC, FLT)")
	maxClaims := set.Int("max-claims", 200, "bounded claims in the unit context")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	showModel := set.String("show-model", "off", "stream raw model output: off|stderr|stdout")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	llmReasoning := set.String("llm-reasoning", "", "override "+llm.EnvReasoning)
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec synthesis run requires --db")
	}
	catalog, ok := synthesis.Catalogs[strings.ToUpper(*catalogID)]
	if !ok {
		return fmt.Errorf("unsupported synthesis catalog %q", *catalogID)
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	claims, err := synthesis.LoadClaims(*db, res.EIR)
	if err != nil {
		return err
	}
	if len(claims) == 0 {
		return fmt.Errorf("synthesis: no committed claims (run spec assurance claims first)")
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
	client, err := llm.NewClient(config)
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
	if iteration := header["current-iteration"]; iteration != "" {
		runConfig.Iteration = iteration
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
	unit := synthesis.Plan(catalog, claims, synthesis.Directives{
		Namespace:      header["system-namespace"],
		Snapshot:       header["included-environments-snapshots"],
		StorageBinding: header["canonical-storage-binding"],
		MaxClaims:      *maxClaims,
	})
	session, err := work.NewSession(*db, res.EIR, runConfig, executor, reviewer)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "claims: %d, catalog: %s\n", len(claims), catalog.FieldID)
	report, runErr := session.Run([]work.Unit{unit})
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-20s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "synthesis stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "synthesis committed %d record(s)\n", len(report.Units[0].MutationIDs))
	return nil
}
