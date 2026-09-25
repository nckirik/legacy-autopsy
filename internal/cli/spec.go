package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/runtime"
)

const (
	defaultSpecSource = "examples/spec/export-reconciliation.cdl"
	defaultSpecLedger = "examples/spec/identity-ledger.json"
)

func specCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec <compile|render|run> [flags]")
	}
	switch args[0] {
	case "compile":
		return specCompile(args[1:], stdout, stderr)
	case "render":
		return specRender(args[1:], stdout, stderr)
	case "run":
		return specRun(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec operation %q is unsupported; supported operations: compile, render, run", args[0])
	}
}

// compileSpec compiles one CDL source with its committed identity ledger. The
// source path is logical for fingerprinting, so results do not depend on cwd.
func compileSpec(repo, sourcePath, ledgerPath string) (*cdl.Result, string, error) {
	sourceAbs := sourcePath
	if !filepath.IsAbs(sourceAbs) {
		sourceAbs = filepath.Join(repo, filepath.FromSlash(sourcePath))
	}
	sourceBytes, err := os.ReadFile(sourceAbs)
	if err != nil {
		return nil, "", err
	}
	logical := filepath.ToSlash(sourcePath)
	ledgerAbs := ledgerPath
	if !filepath.IsAbs(ledgerAbs) {
		ledgerAbs = filepath.Join(repo, filepath.FromSlash(ledgerPath))
	}
	ledger, err := cdl.LoadLedger(ledgerAbs)
	if err != nil {
		return nil, "", fmt.Errorf("identity ledger: %w", err)
	}
	res, err := cdl.Compile(cdl.CompileInput{
		Sources: []cdl.Source{{Path: logical, Bytes: sourceBytes}},
		Ledger:  ledger,
	})
	if err != nil {
		return nil, "", err
	}
	return res, logical, nil
}

func specCompile(args []string, stdout, stderr io.Writer) error {
	set := flags("spec compile", stderr)
	repo := set.String("repo", ".", "repository root")
	source := set.String("source", defaultSpecSource, "CDL source path")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	golden := set.String("golden", "", "expected EIR golden path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	res, logical, err := compileSpec(*repo, *source, *ledger)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "compiled %s\n", logical)
	fmt.Fprintf(stdout, "source:     %s\n", res.SourceFingerprint)
	fmt.Fprintf(stdout, "eir:        %s\n", res.EIRHash)
	fmt.Fprintf(stdout, "identities: %d\n", len(res.IDs))
	if *golden != "" {
		path := *golden
		if !filepath.IsAbs(path) {
			path = filepath.Join(*repo, filepath.FromSlash(path))
		}
		want, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.TrimRight(string(want), "\n") != string(res.EIRBytes) {
			return fmt.Errorf("EIR golden mismatch: %s", path)
		}
		fmt.Fprintf(stdout, "golden:     %s ok\n", *golden)
	}
	return nil
}

func specRender(args []string, stdout, stderr io.Writer) error {
	set := flags("spec render", stderr)
	repo := set.String("repo", ".", "repository root")
	source := set.String("source", defaultSpecSource, "CDL source path")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	projection := set.String("projection", "LIGHT", "projection id")
	channel := set.String("channel", "prompt", "channel: prompt or native")
	out := set.String("out", "", "output path (default stdout)")
	generatedAt := set.String("generated-at", "unspecified", "declared generation time for the provenance header")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	res, logical, err := compileSpec(*repo, *source, *ledger)
	if err != nil {
		return err
	}
	prompt, err := cdl.RenderPrompt(res.EIR, strings.ToUpper(*projection), strings.ToLower(*channel), cdl.Provenance{
		SourcePath:        logical,
		SourceFingerprint: res.SourceFingerprint,
		EIRHash:           res.EIRHash,
		GeneratedAt:       *generatedAt,
	})
	if err != nil {
		return err
	}
	if *out != "" {
		return os.WriteFile(*out, prompt, 0o644)
	}
	_, err = stdout.Write(prompt)
	return err
}

func specRun(args []string, stdout, stderr io.Writer) error {
	set := flags("spec run", stderr)
	repo := set.String("repo", ".", "repository root")
	source := set.String("source", defaultSpecSource, "CDL source path")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	fixture := set.String("fixture", "", "execution fixture path (required)")
	jsonOut := set.Bool("json", false, "emit the canonical trace as JSON")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *fixture == "" {
		return fmt.Errorf("spec run requires --fixture")
	}
	res, _, err := compileSpec(*repo, *source, *ledger)
	if err != nil {
		return err
	}
	fixturePath := *fixture
	if !filepath.IsAbs(fixturePath) {
		fixturePath = filepath.Join(*repo, filepath.FromSlash(fixturePath))
	}
	fx, err := runtime.LoadFixture(fixturePath)
	if err != nil {
		return err
	}
	caps := map[string]bool{}
	for _, c := range fx.Capabilities {
		caps[c] = true
	}
	trace, err := runtime.Execute(res.EIR, fx, caps)
	if err != nil {
		return err
	}
	if *jsonOut {
		blob, err := cdl.CanonicalBytes(trace)
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, string(blob))
		return nil
	}
	for _, rec := range trace.Steps {
		fmt.Fprintf(stdout, "%-24s %-8s %-9s %s\n", rec.Step, rec.Owner, rec.Outcome, rec.Authority)
		for _, effect := range rec.Effects {
			value := ""
			if effect.Value != "" {
				value = " = " + effect.Value
			}
			fmt.Fprintf(stdout, "    %-10s %s%s\n", effect.Operation, effect.Target, value)
		}
		for _, diag := range rec.Diagnostics {
			fmt.Fprintf(stdout, "    diag       %s: %s\n", diag.Code, diag.Subject)
		}
	}
	hash, err := trace.Hash()
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "trace: %s\n", hash)
	return nil
}
