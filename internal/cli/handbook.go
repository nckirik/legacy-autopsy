package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/handbook"
)

// specHandbook renders the human handbook draft from committed synthesis records.
func specHandbook(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "draft" {
		return fmt.Errorf("usage: legacy-autopsy spec handbook draft [flags]")
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
