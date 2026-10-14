package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// specAssurance runs the Assurance claim pipeline over committed inventory.
func specAssurance(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec assurance <claims> [flags]")
	}
	switch args[0] {
	case "claims":
		return specAssuranceClaims(args[1:], stdout, stderr)
	case "coverage":
		return specAssuranceCoverage(args[1:], stdout, stderr)
	case "sweep":
		return specAssuranceSweep(args[1:], stdout, stderr)
	case "gate":
		return specAssuranceGate(args[1:], stdout, stderr)
	case "advance":
		return specAssuranceAdvance(args[1:], stdout, stderr)
	case "reconcile":
		return specAssuranceReconcile(args[1:], stdout, stderr)
	case "lifecycle":
		return specAssuranceLifecycle(args[1:], stdout, stderr)
	case "privacy":
		return specAssurancePrivacy(args[1:], stdout, stderr)
	case "closure":
		return specAssuranceClosure(args[1:], stdout, stderr)
	case "tickets":
		return specAssuranceTickets(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec assurance operation %q is unsupported; supported operations: claims, coverage, sweep, gate, advance, reconcile, lifecycle, privacy, closure, tickets", args[0])
	}
}

// specAssuranceTickets originates T-OPEN tickets for ticketless unresolved
// semantic predicates (unknown lifecycle classifications) so honest coverage
// is visible as ticket state instead of silent gaps.
func specAssuranceTickets(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance tickets", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace tickets whose identity exists with changed content")
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
	if *db == "" {
		return fmt.Errorf("spec assurance tickets requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	units, err := assurance.TicketPlan(*db, res.EIR, header["system-namespace"], header["included-environments-snapshots"])
	if err != nil {
		return err
	}
	if len(units) == 0 {
		fmt.Fprintf(stdout, "tickets: 0 (no ticketless unresolved predicates)\n")
		return nil
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
	fmt.Fprintf(stdout, "ticket units: %d\n", len(units))
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-40s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "ticket origination stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "tickets originated: %d\n", len(report.Units))
	return nil
}

// specAssuranceClosure validates the deterministic dead-code and
// lifecycle-eligibility closure rules (T11.3).
func specAssuranceClosure(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance closure", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec assurance closure requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	names, _, _, err := work.FieldContract(res.EIR, "ATOMIC-COMPONENT")
	if err != nil {
		return err
	}
	componentIDs, err := store.List(*db, res.EIR, "ATOMIC-COMPONENT")
	if err != nil {
		return err
	}
	checked := 0
	invalid := 0
	for _, id := range componentIDs {
		payload, err := store.Get(*db, res.EIR, "ATOMIC-COMPONENT", id)
		if err != nil {
			return err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return err
		}
		if strings.TrimSpace(values["staging-status"]) != "[S-DEAD-CODE]" {
			continue
		}
		checked++
		if err := assurance.ValidateDeadCode(*db, res.EIR, id); err != nil {
			invalid++
			fmt.Fprintf(stdout, "%-24s FAIL %s\n", id, err)
			continue
		}
		fmt.Fprintf(stdout, "%-24s PASS dead-code proof holds\n", id)
	}
	diagnostics, err := assurance.ValidateLifecycleEligibility(*db, res.EIR)
	if err != nil {
		return err
	}
	for _, diagnostic := range diagnostics {
		fmt.Fprintf(stdout, "lifecycle: %s\n", diagnostic)
	}
	fmt.Fprintf(stdout, "dead-code components: %d checked, %d invalid; lifecycle diagnostics: %d\n", checked, invalid, len(diagnostics))
	if invalid > 0 || len(diagnostics) > 0 {
		return fmt.Errorf("closure validation failed")
	}
	return nil
}

// specAssuranceSweep mechanically records one sweep per owner scope.
func specAssuranceSweep(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance sweep", stderr)
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
		return fmt.Errorf("spec assurance sweep requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	scopes, err := assurance.LoadCellScopes(*db, res.EIR)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return fmt.Errorf("assurance sweep: no open required matrix cells (run spec preflight matrix first)")
	}
	env := *environment
	if env == "" {
		env = header["included-environments-snapshots"]
	}
	snap := *snapshot
	if snap == "" {
		snap = header["included-environments-snapshots"]
	}
	_ = env
	_ = snap
	runConfig := work.DefaultConfig()
	if iteration := header["current-iteration"]; iteration != "" {
		runConfig.Iteration = iteration
	}
	// Cell sweeps are a rebuildable projection of committed coverage state.
	runConfig.UpdateExisting = true
	session, err := work.NewSession(*db, res.EIR, runConfig, assurance.MechanicalCellSweepExecutor{Scopes: scopes}, nil)
	if err != nil {
		return err
	}
	report, err := session.Run([]work.Unit{assurance.CellSweepPlan(scopes)})
	if err != nil {
		return err
	}
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%s: invocation %s committed=%v\n", stats.UnitID, stats.InvocationID, stats.Committed)
	}
	closed, err := assurance.AdvanceMatrix(*db, res.EIR)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "cell sweeps recorded: %d, matrix cells closed: %d\n", len(scopes), closed)
	return nil
}

// specAssuranceAdvance recomputes coverage dispositions from derived-from edges.
func specAssuranceAdvance(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance advance", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec assurance advance requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	changed, err := assurance.AdvanceCoverage(*db, res.EIR)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "coverage rows advanced to Mapped: %d\n", changed)
	return nil
}

// specAssuranceReconcile declares export reconciliation state, including the
// evidence-backed zero domain when no exports exist.
func specAssuranceReconcile(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance reconcile", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec assurance reconcile requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	sources, err := assurance.ReconciliationSources(*db, res.EIR)
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return fmt.Errorf("assurance reconcile: no committed source inventory to cite as evidence")
	}
	runConfig := work.DefaultConfig()
	if iteration := header["current-iteration"]; iteration != "" {
		runConfig.Iteration = iteration
	}
	session, err := work.NewSession(*db, res.EIR, runConfig, assurance.MechanicalReconcileExecutor{SourceIDs: sources}, nil)
	if err != nil {
		return err
	}
	report, err := session.Run([]work.Unit{assurance.ReconcilePlan(header["system-namespace"])})
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "reconciliation: invocation %s, zero-domain=%v\n", report.Units[0].InvocationID, report.Units[0].ZeroDomain != "")
	return nil
}

// specAssuranceLifecycle authors lifecycle eligibility per source unit.
func specAssuranceLifecycle(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance lifecycle", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	batchSize := set.Int("batch-size", 10, "source units per lifecycle unit")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace eligibility whose identity exists with changed content")
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
	if *db == "" {
		return fmt.Errorf("spec assurance lifecycle requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	units, err := assurance.LifecyclePlan(*db, res.EIR, header["system-namespace"], header["included-environments-snapshots"], *batchSize)
	if err != nil {
		return err
	}
	if len(units) == 0 {
		return fmt.Errorf("assurance lifecycle: no committed source inventory")
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
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-26s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "lifecycle stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "lifecycle units committed: %d\n", len(report.Units))
	return nil
}

// specAssurancePrivacy authors 0C privacy findings mirrored to SEC records.
func specAssurancePrivacy(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance privacy", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace findings whose identity exists with changed content")
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
	if *db == "" {
		return fmt.Errorf("spec assurance privacy requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	units, err := assurance.PrivacyPlan(*db, res.EIR, header["system-namespace"], header["included-environments-snapshots"])
	if err != nil {
		return err
	}
	if len(units) == 0 {
		fmt.Fprintf(stdout, "privacy findings: no committed SEC findings to mirror\n")
		return nil
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
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-26s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "privacy stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "privacy units committed: %d\n", len(report.Units))
	return nil
}

// specAssuranceGate evaluates the deterministic Exit A conditions.
func specAssuranceGate(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance gate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec assurance gate requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	gate, err := assurance.EvaluateExitA(*db, res.EIR)
	if err != nil {
		return err
	}
	for _, check := range gate.Checks {
		fmt.Fprintf(stdout, "%-32s %-11s %s\n", check.ID, check.Result, check.Detail)
	}
	if gate.Passed {
		fmt.Fprintln(stdout, "Exit A: Pass")
		return nil
	}
	fmt.Fprintln(stdout, "Exit A: Fail")
	return fmt.Errorf("Exit A conditions are not met")
}

// specAssuranceCoverage mechanically expands coverage rows and prints the
// deterministic coverage arithmetic.
func specAssuranceCoverage(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance coverage", stderr)
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
		return fmt.Errorf("spec assurance coverage requires --db")
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
	created, err := assurance.ExpandCoverage(*db, res.EIR, assurance.CoverageDirectives{
		Namespace: header["system-namespace"], Environment: env, Snapshot: snap,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "coverage rows created: %d\n", created)
	summary, err := assurance.CoverageSummary(*db, res.EIR)
	if err != nil {
		return err
	}
	for _, kind := range assurance.SortedKeys(summary.SourcesByKind) {
		fmt.Fprintf(stdout, "sources[%s]: %d\n", kind, summary.SourcesByKind[kind])
	}
	for _, disposition := range assurance.SortedKeys(summary.CoverageByDisposition) {
		fmt.Fprintf(stdout, "coverage[%s]: %d\n", disposition, summary.CoverageByDisposition[disposition])
	}
	return nil
}

func specAssuranceClaims(args []string, stdout, stderr io.Writer) error {
	set := flags("spec assurance claims", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	target := set.String("target", "", "target legacy repository (required)")
	db := set.String("db", "", "SQLite database path (required)")
	batchSize := set.Int("batch-size", 10, "source units per claim unit")
	maxExcerpt := set.Int("max-excerpt", 1200, "bytes of file excerpt per source unit")
	reviewMode := set.String("review-mode", "", "override review policy: all|critical|uncertainty")
	maxRetries := set.Int("max-retries", 0, "override the retry budget")
	updateExisting := set.Bool("update-existing", false, "replace claims whose evidence-anchored identity exists with changed prose")
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
		return fmt.Errorf("spec assurance claims requires --target and --db")
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
		return fmt.Errorf("assurance: no committed source inventory (run spec discovery run first)")
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
	units, err := assurance.Plan(*target, inventory, assurance.Directives{
		Namespace:       header["system-namespace"],
		Snapshot:        header["included-environments-snapshots"],
		StorageBinding:  header["canonical-storage-binding"],
		BatchSize:       *batchSize,
		MaxExcerptBytes: *maxExcerpt,
	})
	if err != nil {
		return err
	}
	session, err := work.NewSession(*db, res.EIR, runConfig, executor, reviewer)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "sources: %d in %d claim unit(s)\n", len(inventory), len(units))
	report, runErr := session.Run(units)
	for _, stats := range report.Units {
		fmt.Fprintf(stdout, "%-26s %8d %5d %6d %7d %-12s %v\n",
			stats.UnitID, stats.ContextBytes, stats.Dependencies, stats.Fetches, stats.Retries, stats.Review, stats.Escalated)
		if stats.LastError != "" {
			fmt.Fprintf(stdout, "  last rejection: %s\n", stats.LastError)
		}
	}
	if runErr != nil {
		fmt.Fprintf(stderr, "assurance stopped: %s\n", runErr)
		return runErr
	}
	fmt.Fprintf(stdout, "assurance committed %d claim unit(s)\n", len(report.Units))
	return nil
}

// preflightHeaderValues reads the committed preflight header record.
func preflightHeaderValues(db string, doc *cdl.EIRDoc) (map[string]string, error) {
	payload, err := store.Get(db, doc, "PREFLIGHT-HEADER", "0A")
	if err != nil {
		return nil, fmt.Errorf("requires a completed preflight (PREFLIGHT-HEADER/0A): %w", err)
	}
	names, _, _, err := work.FieldContract(doc, "PREFLIGHT-HEADER")
	if err != nil {
		return nil, err
	}
	return work.PayloadText(names, payload)
}
