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
	examplePath          = "examples/spec/export-reconciliation.cdl"
	typedIDPath          = "examples/spec/typed-id.cdl"
	semanticPayloadPath  = "examples/spec/semantic-payload-identity.cdl"
	canonicalProfilePath = "examples/spec/canonical-hash-profile.cdl"
	acquisitionPath      = "examples/spec/export-acquisition-loop.cdl"
	normalizedMapsPath   = "examples/spec/normalized-maps.cdl"
	ledgerPath           = "examples/spec/identity-ledger.json"
	goldenEIRPath        = "examples/spec/golden/export-reconciliation.eir.json"
	goldenPromptPath     = "examples/spec/golden/export-reconciliation.light.prompt.md"
	generatedAt          = "2026-09-25T00:00:00Z"
)

func repoPath(rel string) string { return filepath.Join("..", rel) }

func readRepoFile(t *testing.T, rel string) []byte {
	t.Helper()
	b, err := os.ReadFile(repoPath(rel))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func exampleSource(t *testing.T) Source {
	t.Helper()
	return Source{Path: examplePath, Bytes: readRepoFile(t, examplePath)}
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
	var ledgerSources []LedgerSource
	for _, rel := range []string{examplePath, typedIDPath, semanticPayloadPath, canonicalProfilePath, acquisitionPath, normalizedMapsPath} {
		src := Source{Path: rel, Bytes: readRepoFile(t, rel)}
		prog, diags := Parse(src.Path, src.Bytes)
		if len(diags) > 0 {
			t.Fatalf("%s parse: %v", rel, diags)
		}
		if semDiags, _ := resolve(prog); len(semDiags) > 0 {
			t.Fatalf("%s resolve: %v", rel, semDiags)
		}
		ledgerSources = append(ledgerSources, LedgerSource{
			Path:              rel,
			SourceFingerprint: sourceFingerprint([]Source{src}),
			IDs:               collectIDs(prog),
		})
	}
	if err := os.MkdirAll(filepath.Dir(repoPath(goldenEIRPath)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := WriteLedger(repoPath(ledgerPath), ledgerSources, LedgerProvenance{
		Generator:   GeneratorVersion,
		GeneratedAt: generatedAt,
	}); err != nil {
		t.Fatal(err)
	}
	src := exampleSource(t)
	prog, diags := Parse(src.Path, src.Bytes)
	if len(diags) > 0 {
		t.Fatalf("parse: %v", diags)
	}
	if semDiags, _ := resolve(prog); len(semDiags) > 0 {
		t.Fatalf("resolve: %v", semDiags)
	}
	fingerprint := sourceFingerprint([]Source{src})
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

func TestTypedIDSectionCompiles(t *testing.T) {
	src := Source{Path: typedIDPath, Bytes: readRepoFile(t, typedIDPath)}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compile(CompileInput{Sources: []Source{src}, Ledger: ledger})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "4.1" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	counts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		counts[e.ID] = len(e.Values)
	}
	if counts["ID-PREFIX"] != 59 {
		t.Fatalf("ID-PREFIX has %d values, want 59", counts["ID-PREFIX"])
	}
	if counts["ID-ITERATION"] != 26 {
		t.Fatalf("ID-ITERATION has %d values, want 26", counts["ID-ITERATION"])
	}
	rules := map[string]bool{}
	for _, r := range res.EIR.Declarations.Rules {
		rules[r.ID] = true
	}
	for _, want := range []string{
		"canonical-key", "src-kind-and-coordinates", "normalized-relative-path",
		"prf-hbk-identity", "collision-extension", "id-registry-checks",
	} {
		if !rules[want] {
			t.Fatalf("missing rule %s", want)
		}
	}
}

func TestSemanticPayloadSectionCompiles(t *testing.T) {
	src := Source{Path: semanticPayloadPath, Bytes: readRepoFile(t, semanticPayloadPath)}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compile(CompileInput{Sources: []Source{src}, Ledger: ledger})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "4.1.1" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	rules := map[string]bool{}
	for _, r := range res.EIR.Declarations.Rules {
		rules[r.ID] = true
	}
	for _, want := range []string{
		"payload-mandatory-fields", "semantic-record-version", "semantic-content-fingerprint",
		"certification-envelope", "envelope-mutation-invariance", "confirmation-binding",
	} {
		if !rules[want] {
			t.Fatalf("missing rule %s", want)
		}
	}
}

func TestCanonicalProfileSectionCompiles(t *testing.T) {
	src := Source{Path: canonicalProfilePath, Bytes: readRepoFile(t, canonicalProfilePath)}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compile(CompileInput{Sources: []Source{src}, Ledger: ledger})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "4.1.2" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	counts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		counts[e.ID] = len(e.Values)
	}
	if counts["EVIDENCE-BINDING-KIND"] != 12 || counts["ENVELOPE-KIND"] != 3 ||
		counts["ARTIFACT-TYPE"] != 5 || counts["BINDING-ROW-KIND"] != 2 ||
		counts["HASH-DOMAIN-PREFIX"] != 3 {
		t.Fatalf("unexpected enum counts: %v", counts)
	}
}

func TestAcquisitionLoopSectionCompiles(t *testing.T) {
	src := Source{Path: acquisitionPath, Bytes: readRepoFile(t, acquisitionPath)}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compile(CompileInput{Sources: []Source{src}, Ledger: ledger})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "7.5" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	states := 0
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "ACQUISITION-STATE" {
			states = len(e.Values)
		}
	}
	if states != 5 {
		t.Fatalf("ACQUISITION-STATE has %d values, want 5", states)
	}
}

func TestNormalizedMapsSectionCompiles(t *testing.T) {
	src := Source{Path: normalizedMapsPath, Bytes: readRepoFile(t, normalizedMapsPath)}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compile(CompileInput{Sources: []Source{src}, Ledger: ledger})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "7.6" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	var tableID string
	for _, table := range res.EIR.Declarations.Tables {
		tableID = table.ID
	}
	if tableID != "normalized-map" {
		t.Fatalf("unexpected table %q", tableID)
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
