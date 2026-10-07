package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/discovery"
	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specDiscovery enumerates the target repository and runs the Discovery
// inventory pipeline against the canonical store.
func specDiscovery(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec discovery <run|frontier> [flags]")
	}
	op := args[0]
	switch op {
	case "run":
		return specDiscoveryRun(args[1:], stdout, stderr)
	case "frontier":
		return specDiscoveryFrontier(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec discovery operation %q is unsupported; supported operations: run, frontier", op)
	}
}

// specDiscoveryFrontier records one mechanical terminal frontier per committed
// source unit. It needs no model; the runtime authors the rows and the 0G entry.
func specDiscoveryFrontier(args []string, stdout, stderr io.Writer) error {
	set := flags("spec discovery frontier", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	iteration := set.String("iteration", "", "override the §10.6 iteration token")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec discovery frontier requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	inventory, err := assurance.LoadInventory(*db, res.EIR)
	if err != nil {
		return err
	}
	if len(inventory) == 0 {
		return fmt.Errorf("discovery frontier: no committed source inventory (run spec discovery run first)")
	}
	units := make([]discovery.SourceUnit, 0, len(inventory))
	for _, item := range inventory {
		units = append(units, discovery.SourceUnit{ID: item.ID, Coordinate: item.Coordinate, Fingerprint: item.Fingerprint})
	}
	runConfig := work.DefaultConfig()
	if iteration := header["current-iteration"]; iteration != "" {
		runConfig.Iteration = iteration
	}
	if *iteration != "" {
		runConfig.Iteration = *iteration
	}
	session, err := work.NewSession(*db, res.EIR, runConfig, discovery.MechanicalFrontierExecutor{Inventory: units}, nil)
	if err != nil {
		return err
	}
	report, err := session.Run([]work.Unit{discovery.FrontierPlan(units, header["system-namespace"])})
	if err != nil {
		return err
	}
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%s: invocation %s committed=%v\n", stats.UnitID, stats.InvocationID, stats.Committed)
	}
	fmt.Fprintf(stdout, "terminal frontiers recorded: %d\n", len(inventory))
	return nil
}

func specDiscoveryRun(args []string, stdout, stderr io.Writer) error {
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
	windowOffset := set.Int("window-offset", 0, "skip the first N sorted candidates (windowed runs)")
	showModel := set.String("show-model", "off", "stream raw model output: off|stderr|stdout")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	llmReasoning := set.String("llm-reasoning", "", "override "+llm.EnvReasoning)
	if err := set.Parse(args); err != nil {
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
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	scopePaths := []string{}
	for _, cluster := range clusters {
		scopePaths = append(scopePaths, cluster.ScopePaths...)
	}
	candidates, err := discovery.EnumerateWindow(*target, scopePaths, *windowOffset, *maxFiles)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		return fmt.Errorf("discovery: no candidates under scope %v at offset %d", scopePaths, *windowOffset)
	}
	if manifest, err := writeWindowManifest(*db, *target, directives.Namespace, *windowOffset, candidates); err != nil {
		return err
	} else {
		fmt.Fprintf(stdout, "window manifest: %s\n", manifest)
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

// writeWindowManifest records which files were in a bounded discovery window,
// so successive windowed runs are auditable. The manifest is a debug/audit
// derivative, never operational state.
func writeWindowManifest(db, target, namespace string, offset int, candidates []discovery.Candidate) (string, error) {
	type manifestCandidate struct {
		Path        string `json:"path"`
		Size        int64  `json:"size"`
		Fingerprint string `json:"fingerprint"`
	}
	manifest := struct {
		Target     string              `json:"target"`
		Namespace  string              `json:"namespace"`
		Offset     int                 `json:"offset"`
		Count      int                 `json:"count"`
		Candidates []manifestCandidate `json:"candidates"`
	}{
		Target: target, Namespace: namespace, Offset: offset, Count: len(candidates),
	}
	for _, candidate := range candidates {
		manifest.Candidates = append(manifest.Candidates, manifestCandidate{Path: candidate.Path, Size: candidate.Size, Fingerprint: candidate.Fingerprint})
	}
	dir := filepath.Join(filepath.Dir(db), "windows")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("discovery-%05d-%05d.json", offset, len(candidates)))
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
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
