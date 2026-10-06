package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/runtime"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

const (
	defaultSpecAssembly = "protocol/assembly.json"
	defaultSpecLedger   = "protocol/identity-ledger.json"
)

func specCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec <compile|render|run|schema|store|record|work|artifact|preflight|discovery|assurance|synthesis> [flags]")
	}
	switch args[0] {
	case "compile":
		return specCompile(args[1:], stdout, stderr)
	case "render":
		return specRender(args[1:], stdout, stderr)
	case "run":
		return specRun(args[1:], stdout, stderr)
	case "schema":
		return specSchema(args[1:], stdout, stderr)
	case "store":
		return specStore(args[1:], stdout, stderr)
	case "record":
		return specRecord(args[1:], stdout, stderr)
	case "work":
		return specWork(args[1:], stdout, stderr)
	case "artifact":
		return specArtifact(args[1:], stdout, stderr)
	case "preflight":
		return specPreflight(args[1:], stdout, stderr)
	case "discovery":
		return specDiscovery(args[1:], stdout, stderr)
	case "assurance":
		return specAssurance(args[1:], stdout, stderr)
	case "synthesis":
		return specSynthesis(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec operation %q is unsupported; supported operations: compile, render, run, schema, store, record, work, artifact, preflight, discovery, assurance, synthesis", args[0])
	}
}

func specRecord(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec record <get|list|put|render|inspect> [flags]")
	}
	op := args[0]
	switch op {
	case "get", "list", "put", "render", "inspect":
	default:
		return fmt.Errorf("spec record operation %q is unsupported; supported operations: get, list, put, render, inspect", op)
	}
	set := flags("spec record "+op, stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path")
	field := set.String("field", "", "FIELD ID (table class)")
	id := set.String("id", "", "record ID")
	version := set.Int("version", 1, "record version for put")
	payloadPath := set.String("payload", "", "Markdown payload file for put")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" || *field == "" {
		return fmt.Errorf("spec record %s requires --db and --field", op)
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	switch op {
	case "get":
		if *id == "" {
			return fmt.Errorf("spec record get requires --id")
		}
		payload, err := store.Get(*db, res.EIR, *field, *id)
		if err != nil {
			return err
		}
		_, err = stdout.Write(payload)
		return err
	case "list":
		ids, err := store.List(*db, res.EIR, *field)
		if err != nil {
			return err
		}
		for _, recordID := range ids {
			fmt.Fprintln(stdout, recordID)
		}
		return nil
	case "render":
		view, err := store.Render(*db, res.EIR, *field)
		if err != nil {
			return err
		}
		_, err = stdout.Write(view)
		return err
	case "inspect":
		view, err := store.Inspect(*db, res.EIR, *field)
		if err != nil {
			return err
		}
		_, err = stdout.Write(view)
		return err
	default: // put
		if *id == "" || *payloadPath == "" {
			return fmt.Errorf("spec record put requires --id and --payload")
		}
		payload, err := os.ReadFile(*payloadPath)
		if err != nil {
			return err
		}
		if err := store.Ingest(*db, res.EIR, *field, *id, *version, payload); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "put %s %s (version %d)\n", *field, *id, *version)
		return nil
	}
}

// compileSpec compiles one CDL source with its committed identity ledger. The
// source path is logical for fingerprinting, so results do not depend on cwd.
// compileSpec compiles the ordered document assembly, or the shared globals plus
// one named section when sourcePath is set. The returned label is the logical
// compilation input for provenance.
func compileSpec(repo, assemblyPath, sourcePath, ledgerPath string) (*cdl.Result, string, error) {
	var sources []cdl.Source
	logical := assemblyPath
	if sourcePath != "" {
		globalsPath := "protocol/globals.cdl"
		globalsBytes, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(globalsPath)))
		if err != nil {
			return nil, "", err
		}
		sourceAbs := sourcePath
		if !filepath.IsAbs(sourceAbs) {
			sourceAbs = filepath.Join(repo, filepath.FromSlash(sourcePath))
		}
		sourceBytes, err := os.ReadFile(sourceAbs)
		if err != nil {
			return nil, "", err
		}
		logical = filepath.ToSlash(sourcePath)
		sources = []cdl.Source{
			{Path: globalsPath, Bytes: globalsBytes},
			{Path: logical, Bytes: sourceBytes},
		}
	} else {
		loaded, err := cdl.LoadAssembly(repo, assemblyPath)
		if err != nil {
			return nil, "", err
		}
		sources = loaded
	}
	ledgerAbs := ledgerPath
	if !filepath.IsAbs(ledgerAbs) {
		ledgerAbs = filepath.Join(repo, filepath.FromSlash(ledgerPath))
	}
	ledger, err := cdl.LoadLedger(ledgerAbs)
	if err != nil {
		return nil, "", fmt.Errorf("identity ledger: %w", err)
	}
	res, err := cdl.Compile(cdl.CompileInput{Sources: sources, Ledger: ledger})
	if err != nil {
		return nil, "", err
	}
	return res, logical, nil
}

func specCompile(args []string, stdout, stderr io.Writer) error {
	set := flags("spec compile", stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	golden := set.String("golden", "", "expected EIR golden path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	res, logical, err := compileSpec(*repo, *assembly, *source, *ledger)
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

func specSchema(args []string, stdout, stderr io.Writer) error {
	set := flags("spec schema", stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	out := set.String("out", "", "write DDL to this path instead of stdout")
	apply := set.String("apply", "", "apply DDL to this SQLite database path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	ddl, err := store.Schema(res.EIR, store.Provenance{
		Generator:         cdl.GeneratorVersion,
		SourceFingerprint: res.SourceFingerprint,
	})
	if err != nil {
		return err
	}
	if *apply != "" {
		if err := store.Apply(*apply, ddl); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "applied schema to %s\n", *apply)
	}
	if *out != "" {
		path := *out
		if !filepath.IsAbs(path) {
			path = filepath.Join(*repo, filepath.FromSlash(path))
		}
		if err := os.WriteFile(path, []byte(ddl), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "wrote schema to %s\n", *out)
		return nil
	}
	if *apply == "" {
		_, err := io.WriteString(stdout, ddl)
		return err
	}
	return nil
}

func specStore(args []string, stdout, stderr io.Writer) error {
	set := flags("spec store", stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path")
	out := set.String("out", "", "write the canonical dump to this path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec store requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	dump, fingerprint, err := store.CanonicalDump(*db, res.EIR, store.Provenance{
		Generator:         cdl.GeneratorVersion,
		SourceFingerprint: res.SourceFingerprint,
	})
	if err != nil {
		return err
	}
	if *out != "" {
		path := *out
		if !filepath.IsAbs(path) {
			path = filepath.Join(*repo, filepath.FromSlash(path))
		}
		if err := os.WriteFile(path, dump, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "wrote canonical dump to %s\nfingerprint: %s\n", *out, fingerprint)
		return nil
	}
	if _, err := stdout.Write(dump); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "fingerprint: %s\n", fingerprint)
	return nil
}

func specRender(args []string, stdout, stderr io.Writer) error {
	set := flags("spec render", stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
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
	res, logical, err := compileSpec(*repo, *assembly, *source, *ledger)
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
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
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
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
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
