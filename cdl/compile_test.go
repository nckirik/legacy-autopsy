package cdl

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden assets and identity ledger")

const (
	examplePath      = "examples/spec/export-reconciliation.cdl"
	ledgerPath       = "examples/spec/identity-ledger.json"
	goldenEIRPath    = "examples/spec/golden/export-reconciliation.eir.json"
	goldenPromptPath = "examples/spec/golden/export-reconciliation.light.prompt.md"
	generatedAt      = "2026-09-25T00:00:00Z"
)

func repoPath(rel string) string { return filepath.Join("..", rel) }

func exampleSource(t *testing.T) Source {
	t.Helper()
	b, err := os.ReadFile(repoPath(examplePath))
	if err != nil {
		t.Fatal(err)
	}
	return Source{Path: examplePath, Bytes: b}
}

func compileExample(t *testing.T) *Result {
	t.Helper()
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatalf("load ledger (run: go test ./cdl -run TestUpdateAssets -update): %v", err)
	}
	res, err := Compile(CompileInput{Sources: []Source{exampleSource(t)}, Ledger: ledger})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return res
}

func TestUpdateAssets(t *testing.T) {
	if !*update {
		t.Skip("run with -update to rewrite golden assets and identity ledger")
	}
	src := exampleSource(t)
	prog, diags := Parse(src.Path, src.Bytes)
	if len(diags) > 0 {
		t.Fatalf("parse: %v", diags)
	}
	semDiags, _ := resolve(prog)
	if len(semDiags) > 0 {
		t.Fatalf("resolve: %v", semDiags)
	}
	if err := os.MkdirAll(filepath.Dir(repoPath(goldenEIRPath)), 0o755); err != nil {
		t.Fatal(err)
	}
	fingerprint := sourceFingerprint([]Source{src})
	if err := WriteLedger(repoPath(ledgerPath), collectIDs(prog), LedgerProvenance{
		Generator:         GeneratorVersion,
		Source:            examplePath,
		SourceFingerprint: fingerprint,
		GeneratedAt:       generatedAt,
	}); err != nil {
		t.Fatal(err)
	}
	doc, eirBytes, _, err := BuildEIR(prog, Versions{}, fingerprint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoPath(goldenEIRPath), append(eirBytes, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	prompt, err := RenderPrompt(doc, "LIGHT", "prompt", Provenance{
		SourcePath:        examplePath,
		SourceFingerprint: fingerprint,
		EIRHash:           doc.Envelope.EIRHash,
		GeneratedAt:       generatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoPath(goldenPromptPath), prompt, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestGoldenEIR(t *testing.T) {
	res := compileExample(t)
	want, err := os.ReadFile(repoPath(goldenEIRPath))
	if err != nil {
		t.Fatal(err)
	}
	got := append(append([]byte(nil), res.EIRBytes...), '\n')
	if string(got) != string(want) {
		t.Fatalf("EIR golden mismatch\ngot:\n%s", got)
	}
}

func TestGoldenPrompt(t *testing.T) {
	res := compileExample(t)
	prompt, err := RenderPrompt(res.EIR, "LIGHT", "prompt", Provenance{
		SourcePath:        examplePath,
		SourceFingerprint: res.SourceFingerprint,
		EIRHash:           res.EIRHash,
		GeneratedAt:       generatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(repoPath(goldenPromptPath))
	if err != nil {
		t.Fatal(err)
	}
	if string(prompt) != string(want) {
		t.Fatalf("prompt golden mismatch\ngot:\n%s", prompt)
	}
}

func TestCompileExample(t *testing.T) {
	res := compileExample(t)
	if !strings.HasPrefix(res.EIRHash, "sha256:") {
		t.Fatalf("bad eir hash %q", res.EIRHash)
	}
	if len(res.IDs) == 0 {
		t.Fatal("no identities collected")
	}
	if res.EIR.Envelope.Protocol != ProtocolVersion {
		t.Fatalf("protocol %q", res.EIR.Envelope.Protocol)
	}
}

func TestLedgerRejectsUnknownIdentity(t *testing.T) {
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = Compile(CompileInput{Sources: []Source{exampleSource(t)}, Ledger: ledger})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	trimmed := map[string]IDRecord{}
	for k, v := range ledger {
		if k != "reconcile.guard" {
			trimmed[k] = v
		}
	}
	_, err = Compile(CompileInput{Sources: []Source{exampleSource(t)}, Ledger: trimmed})
	diags, ok := err.(Diagnostics)
	if !ok || !diags.Has("CDL_ID_LEDGER_UNKNOWN_ID") {
		t.Fatalf("expected CDL_ID_LEDGER_UNKNOWN_ID, got %v", err)
	}
}

func TestNegativeFixtures(t *testing.T) {
	base := string(exampleSource(t).Bytes)
	cases := []struct {
		name   string
		mutate func(string) string
		code   string
	}{
		{
			name: "omission missing supplied-by",
			mutate: func(s string) string {
				return strings.Replace(s, "    SUPPLIED-BY native-harness\n", "", 1)
			},
			code: "CDL_OMISSION_MISSING_SUPPLIED_BY",
		},
		{
			name: "set targets a field",
			mutate: func(s string) string {
				return strings.Replace(s,
					"VALUE fingerprint\n  TYPE hash\n  LIFETIME section\nEND",
					"FIELD fingerprint\n  value -> hash required\nEND", 1)
			},
			code: "CDL_SET_TARGET_INVALID",
		},
		{
			name: "guard lacks explicit comparison",
			mutate: func(s string) string {
				return strings.Replace(s, "  GUARD reconciliation-complete == true", "  GUARD reconciliation-complete", 1)
			},
			code: "CDL_GUARD_NOT_BOOLEAN",
		},
		{
			name: "COMPUTE is not in the grammar",
			mutate: func(s string) string {
				return strings.Replace(s, "  SET reconciliation-complete := reconciliation-completeness.complete",
					"  SET reconciliation-complete := COMPUTE reconciliation-completeness", 1)
			},
			code: "CDL_UNKNOWN_INSTRUCTION",
		},
		{
			name: "block target is undeclared",
			mutate: func(s string) string {
				return strings.Replace(s, "  OTHERWISE BLOCK A-STRUCTURALLY-COMPLETE, R-SWEPT, Exit-A",
					"  OTHERWISE BLOCK undeclared-target", 1)
			},
			code: "CDL_BLOCK_TARGET_UNRESOLVED",
		},
		{
			name: "allows violation",
			mutate: func(s string) string {
				return strings.Replace(s, "    unknown -> true\n", "", 1)
			},
			code: "CDL_ALLOWS_VIOLATION",
		},
		{
			name: "requires summary mismatch",
			mutate: func(s string) string {
				return strings.Replace(s, "REQUIRES CAPABILITY hash.sha256, table.serialize",
					"REQUIRES CAPABILITY hash.sha256", 1)
			},
			code: "CDL_REQUIRES_SUMMARY_MISMATCH",
		},
		{
			name: "unresolved reference",
			mutate: func(s string) string {
				return strings.Replace(s, "USES REGISTRY normalized-units", "USES REGISTRY missing-registry", 1)
			},
			code: "CDL_UNRESOLVED_REFERENCE",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Compile(CompileInput{
				Sources: []Source{{Path: examplePath, Bytes: []byte(tc.mutate(base))}},
				Ledger:  map[string]IDRecord{},
			})
			diags, ok := err.(Diagnostics)
			if !ok {
				t.Fatalf("expected Diagnostics, got %T (%v)", err, err)
			}
			if !diags.Has(tc.code) {
				t.Fatalf("expected %s, got %v", tc.code, diags)
			}
		})
	}
}

func TestDependencyBoundary(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "legacy-autopsy/internal") {
			t.Fatalf("%s imports an internal runtime package", f)
		}
	}
}
