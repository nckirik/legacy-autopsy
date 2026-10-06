package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/discovery"
	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specDiscovery enumerates the target repository and runs the Discovery
// inventory pipeline against the canonical store.
func specDiscovery(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec discovery <run> [flags]")
	}
	op := args[0]
	if op != "run" {
		return fmt.Errorf("spec discovery operation %q is unsupported; supported operations: run", op)
	}
	set := flags("spec discovery run", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	target := set.String("target", "", "target legacy repository (required)")
	db := set.String("db", "", "SQLite database path (required)")
	maxFiles := set.Int("max-files", 2000, "bounded candidate count")
	batchSize := set.Int("batch-size", 25, "candidates per inventory unit")
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
	if *target == "" || *db == "" {
		return fmt.Errorf("spec discovery run requires --target and --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	directives, clusters, err := discoveryContext(*db, res.EIR)
	if err != nil {
		return err
	}
	scopePaths := []string{}
	for _, cluster := range clusters {
		scopePaths = append(scopePaths, cluster.ScopePaths...)
	}
	candidates, err := discovery.Enumerate(*target, scopePaths, *maxFiles)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		return fmt.Errorf("discovery: no candidates under scope %v", scopePaths)
	}
	directives.BatchSize = *batchSize
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
	fmt.Fprintf(stdout, "candidates: %d across %d cluster scope(s)\n", len(candidates), len(clusters))
	units := discovery.Plan(candidates, clusters, directives)
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-26s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "discovery stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "discovery committed %d inventory unit(s)\n", len(report.Units))
	return nil
}

// discoveryContext reads the committed preflight header and clusters from the
// store; Discovery requires a completed Preflight.
func discoveryContext(db string, doc *cdl.EIRDoc) (discovery.Directives, []discovery.Cluster, error) {
	headerPayload, err := store.Get(db, doc, "PREFLIGHT-HEADER", "0A")
	if err != nil {
		return discovery.Directives{}, nil, fmt.Errorf("discovery requires a completed preflight (PREFLIGHT-HEADER/0A): %w", err)
	}
	headerNames, _, _, err := work.FieldContract(doc, "PREFLIGHT-HEADER")
	if err != nil {
		return discovery.Directives{}, nil, err
	}
	headerValues, err := work.PayloadText(headerNames, headerPayload)
	if err != nil {
		return discovery.Directives{}, nil, err
	}
	directives := discovery.Directives{
		Namespace:      headerValues["system-namespace"],
		Snapshot:       headerValues["included-environments-snapshots"],
		StorageBinding: headerValues["canonical-storage-binding"],
	}
	names, _, _, err := work.FieldContract(doc, "ENTRY-CLUSTER")
	if err != nil {
		return discovery.Directives{}, nil, err
	}
	ids, err := store.List(db, doc, "ENTRY-CLUSTER")
	if err != nil {
		return discovery.Directives{}, nil, err
	}
	var clusters []discovery.Cluster
	for _, id := range ids {
		payload, err := store.Get(db, doc, "ENTRY-CLUSTER", id)
		if err != nil {
			return discovery.Directives{}, nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return discovery.Directives{}, nil, err
		}
		clusters = append(clusters, discovery.Cluster{
			ID:         values["cluster-id"],
			Title:      values["cluster-title"],
			Owner:      values["canonical-owner-persona"],
			ScopePaths: splitList(values["scope-paths"]),
		})
	}
	return directives, clusters, nil
}

func splitList(raw string) []string {
	var values []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" || strings.EqualFold(part, "none") {
			continue
		}
		values = append(values, part)
	}
	return values
}
