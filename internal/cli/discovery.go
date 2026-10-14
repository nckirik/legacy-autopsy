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
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/discovery"
	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specDiscovery enumerates the target repository and runs the Discovery
// inventory pipeline against the canonical store.
func specDiscovery(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec discovery <run|hits|components|pov-manifest|frontier> [flags]")
	}
	op := args[0]
	switch op {
	case "run":
		return specDiscoveryRun(args[1:], stdout, stderr)
	case "hits":
		return specDiscoveryHits(args[1:], stdout, stderr)
	case "components":
		return specDiscoveryComponents(args[1:], stdout, stderr)
	case "pov-manifest":
		return specDiscoveryPOVManifest(args[1:], stdout, stderr)
	case "frontier":
		return specDiscoveryFrontier(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec discovery operation %q is unsupported; supported operations: run, hits, components, pov-manifest, frontier", op)
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
	// Frontiers are a rebuildable projection of committed inventory.
	runConfig.UpdateExisting = true
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
	updateExisting := set.Bool("update-existing", false, "replace records whose evidence-anchored identity exists with changed content")
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
	runConfig.UpdateExisting = *updateExisting
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

// specDiscoveryHits enumerates route, query, and job source units inside the
// enumerated file window and runs the Discovery inventory pipeline over them.
// The hit scan is deterministic and bounded; the model authors one inventory
// record per hit and the runtime computes every SRC identity.
func specDiscoveryHits(args []string, stdout, stderr io.Writer) error {
	set := flags("spec discovery hits", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	target := set.String("target", "", "target legacy repository (required)")
	db := set.String("db", "", "SQLite database path (required)")
	maxFiles := set.Int("max-files", 2000, "bounded file scan count")
	maxHits := set.Int("max-hits", 500, "bounded hit count")
	batchSize := set.Int("batch-size", 25, "hits per inventory unit")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace records whose evidence-anchored identity exists with changed content")
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
		return fmt.Errorf("spec discovery hits requires --target and --db")
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
	files, err := discovery.EnumerateWindow(*target, scopePaths, 0, *maxFiles)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("discovery hits: no files under scope %v", scopePaths)
	}
	hits, err := discovery.EnumerateHits(*target, files, *maxHits)
	if err != nil {
		return err
	}
	if len(hits) == 0 {
		return fmt.Errorf("discovery hits: no route/query/job hits found under scope %v (the absence proof must be authored separately)", scopePaths)
	}
	if manifest, err := writeHitsWindowManifest(*db, *target, directives.Namespace, hits); err != nil {
		return err
	} else {
		fmt.Fprintf(stdout, "hits window manifest: %s\n", manifest)
	}
	prefixes, err := personaPrefixSet(*db, res.EIR)
	if err != nil {
		return err
	}
	directives.PersonaPrefixes = prefixes
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
	runConfig.UpdateExisting = *updateExisting
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
	fmt.Fprintf(stdout, "hits: %d (files scanned: %d across %d cluster scope(s))\n", len(hits), len(files), len(clusters))
	units := discovery.HitPlan(hits, clusters, directives)
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-26s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "discovery hits stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "discovery committed %d inventory unit(s) over %d hit(s)\n", len(report.Units), len(hits))
	return nil
}

// specDiscoveryComponents authors one ATOMIC-COMPONENT per enumerated hit with
// runtime-anchored identity in the hit's committed source unit, then
// materializes the dual-entry shared-reference pointer rows mechanically (§6.2).
func specDiscoveryComponents(args []string, stdout, stderr io.Writer) error {
	set := flags("spec discovery components", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	target := set.String("target", "", "target legacy repository (required)")
	db := set.String("db", "", "SQLite database path (required)")
	maxFiles := set.Int("max-files", 2000, "bounded file scan count")
	maxHits := set.Int("max-hits", 200, "bounded hit count")
	batchSize := set.Int("batch-size", 5, "hits per component unit")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace records whose evidence-anchored identity exists with changed content")
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
		return fmt.Errorf("spec discovery components requires --target and --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	directives, clusters, err := discoveryContext(*db, res.EIR)
	if err != nil {
		return err
	}
	prefixes, err := personaPrefixSet(*db, res.EIR)
	if err != nil {
		return err
	}
	directives.PersonaPrefixes = prefixes
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	scopePaths := []string{}
	for _, cluster := range clusters {
		scopePaths = append(scopePaths, cluster.ScopePaths...)
	}
	files, err := discovery.EnumerateWindow(*target, scopePaths, 0, *maxFiles)
	if err != nil {
		return err
	}
	hits, err := discovery.EnumerateHits(*target, files, *maxHits)
	if err != nil {
		return err
	}
	if len(hits) == 0 {
		return fmt.Errorf("discovery components: no route/query/job hits found under scope %v", scopePaths)
	}
	for _, hit := range hits {
		srcID, err := capabilities.SourceID(directives.Namespace, string(hit.Kind), hit.Coordinate)
		if err != nil {
			return err
		}
		if _, err := store.Get(*db, res.EIR, "SOURCE-INVENTORY-RECORD", srcID); err != nil {
			return fmt.Errorf("discovery components: hit %s has no committed source unit %s; run spec discovery hits first", hit.Coordinate, srcID)
		}
	}
	directives.BatchSize = *batchSize
	units, err := discovery.ComponentPlan(*target, hits, clusters, directives)
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
	runConfig.UpdateExisting = *updateExisting
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
	fmt.Fprintf(stdout, "component units: %d over %d hit(s)\n", len(units), len(hits))
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-34s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "discovery components stopped: %s\n", runErr)
		return runErr
	}
	components, err := discovery.LoadComponents(*db, res.EIR)
	if err != nil {
		return err
	}
	creating, err := discovery.LoadCreatingInvocations(*db, res.EIR)
	if err != nil {
		return err
	}
	for i := range components {
		if invocation, ok := creating[components[i].ID]; ok {
			components[i].CreatingInvocation = invocation
		}
	}
	refUnits := discovery.SharedReferencePlan(components)
	if len(refUnits) == 0 {
		fmt.Fprintf(stdout, "shared-reference pointers: 0 (no persona-bound components)\n")
		return nil
	}
	mechConfig := work.DefaultConfig()
	mechConfig.UpdateExisting = true
	if iteration := header["current-iteration"]; iteration != "" {
		mechConfig.Iteration = iteration
	}
	refSession, err := work.NewSession(*db, res.EIR, mechConfig, discovery.MechanicalSharedReferenceExecutor{Components: components}, nil)
	if err != nil {
		return err
	}
	refReport, err := refSession.Run(refUnits)
	if err != nil {
		fmt.Fprintf(stderr, "shared-reference materialization stopped: %s\n", err)
		return err
	}
	pointers := 0
	for _, stats := range refReport.Units {
		pointers += len(stats.MutationIDs)
	}
	fmt.Fprintf(stdout, "components committed: %d, shared-reference pointers: %d\n", len(components), pointers)
	return nil
}

// specDiscoveryPOVManifest projects committed shared references into
// persona-local POV manifest files under the workspace personas/ root (or the
// declared --out directory). The projection is deterministic and
// non-authoritative.
func specDiscoveryPOVManifest(args []string, stdout, stderr io.Writer) error {
	set := flags("spec discovery pov-manifest", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	out := set.String("out", "", "personas root (default: <workspace>/personas)")
	dryRun := set.Bool("dry-run", false, "print the manifest paths without writing")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec discovery pov-manifest requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	manifests, err := discovery.POVManifests(*db, res.EIR)
	if err != nil {
		return err
	}
	if len(manifests) == 0 {
		fmt.Fprintf(stdout, "pov manifests: 0 (no committed shared references with registered personas)\n")
		return nil
	}
	root := *out
	if root == "" {
		root = filepath.Join(filepath.Dir(*db), "personas")
	}
	for _, manifest := range manifests {
		fmt.Fprintf(stdout, "%-45s %5d bytes\n", manifest.Path, len(manifest.Markdown))
		if *dryRun {
			continue
		}
		target := filepath.Join(root, manifest.Directory, manifest.FileName)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, manifest.Markdown, 0o644); err != nil {
			return err
		}
	}
	if *dryRun {
		fmt.Fprintf(stdout, "pov manifests: %d (dry run, nothing written)\n", len(manifests))
		return nil
	}
	fmt.Fprintf(stdout, "pov manifests written: %d (root=%s)\n", len(manifests), root)
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
			PrimaryPOV: values["primary-pov"],
			Track:      values["traversal-track"],
			Entries:    splitList(values["concrete-entry-coordinates"]),
			ScopePaths: splitList(values["scope-paths"]),
		})
	}
	return directives, clusters, nil
}

// personaPrefixSet returns the registered persona prefixes from the committed
// 0A registry, so discovery units bind Persona only to registered owners.
func personaPrefixSet(db string, doc *cdl.EIRDoc) (map[string]bool, error) {
	names, _, _, err := work.FieldContract(doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	prefixes := map[string]bool{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "PERSONA-REGISTRY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		prefix := strings.ToUpper(strings.TrimSpace(values["prefix"]))
		if prefix != "" {
			prefixes[prefix] = true
		}
	}
	return prefixes, nil
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

// writeHitsWindowManifest records which route/query/job hits were in a bounded
// hit window, so a run is auditable. The manifest is a debug/audit derivative,
// never operational state.
func writeHitsWindowManifest(db, target, namespace string, hits []discovery.Hit) (string, error) {
	manifest := struct {
		Target    string          `json:"target"`
		Namespace string          `json:"namespace"`
		Count     int             `json:"count"`
		Hits      []discovery.Hit `json:"hits"`
	}{
		Target: target, Namespace: namespace, Count: len(hits), Hits: hits,
	}
	dir := filepath.Join(filepath.Dir(db), "windows")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("discovery-hits-%05d.json", len(hits)))
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
