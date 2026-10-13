package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/handoff"
)

// specHandoff prints the §15.1 candidate condition report.
func specHandoff(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "candidate" {
		return fmt.Errorf("usage: legacy-autopsy spec handoff candidate [flags]")
	}
	set := flags("spec handoff candidate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff candidate requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	checks, err := handoff.EvaluateCandidateChecks(*db, res.EIR)
	if err != nil {
		return err
	}
	failures := 0
	unsupported := 0
	for _, check := range checks {
		fmt.Fprintf(stdout, "%-40s %-12s %s\n", check.ID, check.Result, check.Detail)
		switch check.Result {
		case "Fail":
			failures++
		case "Unsupported":
			unsupported++
		}
	}
	fmt.Fprintf(stdout, "candidate: %d pass, %d fail, %d unsupported\n", len(checks)-failures-unsupported, failures, unsupported)
	if failures > 0 || unsupported > 0 {
		return fmt.Errorf("candidate conditions are not ready")
	}
	return nil
}
