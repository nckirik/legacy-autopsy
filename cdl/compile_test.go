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
	examplePath           = "examples/spec/export-reconciliation.cdl"
	typedIDPath           = "examples/spec/typed-id.cdl"
	semanticPayloadPath   = "examples/spec/semantic-payload-identity.cdl"
	canonicalProfilePath  = "examples/spec/canonical-hash-profile.cdl"
	acquisitionPath       = "examples/spec/export-acquisition-loop.cdl"
	normalizedMapsPath    = "examples/spec/normalized-maps.cdl"
	invocationModesPath   = "examples/spec/invocation-modes.cdl"
	coldResumePath        = "examples/spec/cold-resume.cdl"
	ticketFSMPath         = "examples/spec/ticket-fsm.cdl"
	invocationContextPath = "examples/spec/invocation-context.cdl"
	coverageExitsPath     = "examples/spec/coverage-and-exits.cdl"
	globalsPath           = "examples/spec/globals.cdl"
	assemblyPath          = "examples/spec/assembly.json"
	ledgerPath            = "examples/spec/identity-ledger.json"
	goldenEIRPath         = "examples/spec/golden/protocol.eir.json"
	goldenPromptPath      = "examples/spec/golden/protocol.light.prompt.md"
	generatedAt           = "2026-09-25T00:00:00Z"
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

func globalsSource(t *testing.T) Source {
	t.Helper()
	return Source{Path: globalsPath, Bytes: readRepoFile(t, globalsPath)}
}

// compileSources compiles the shared globals plus the named section sources.
func compileSources(t *testing.T, rels ...string) *Result {
	t.Helper()
	sources := []Source{globalsSource(t)}
	for _, rel := range rels {
		sources = append(sources, Source{Path: rel, Bytes: readRepoFile(t, rel)})
	}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatalf("load ledger (run: go test ./cdl -run TestUpdateAssets -update): %v", err)
	}
	res, err := Compile(CompileInput{Sources: sources, Ledger: ledger})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return res
}

// compileAssembly compiles the full ordered document assembly.
func compileAssembly(t *testing.T) *Result {
	t.Helper()
	sources, err := LoadAssembly(repoPath("."), assemblyPath)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatalf("load ledger (run: go test ./cdl -run TestUpdateAssets -update): %v", err)
	}
	res, err := Compile(CompileInput{Sources: sources, Ledger: ledger})
	if err != nil {
		t.Fatalf("compile assembly: %v", err)
	}
	return res
}

func compileExample(t *testing.T) *Result {
	t.Helper()
	return compileSources(t, examplePath)
}

// mergeInto merges one parsed program into another for test-side assembly.
func mergeInto(dst, src *Program) {
	dst.Globals.Capabilities = append(dst.Globals.Capabilities, src.Globals.Capabilities...)
	dst.Globals.Registries = append(dst.Globals.Registries, src.Globals.Registries...)
	dst.Globals.Artifacts = append(dst.Globals.Artifacts, src.Globals.Artifacts...)
	dst.Globals.Rules = append(dst.Globals.Rules, src.Globals.Rules...)
	dst.Globals.Gates = append(dst.Globals.Gates, src.Globals.Gates...)
	dst.Globals.WorkflowTargets = append(dst.Globals.WorkflowTargets, src.Globals.WorkflowTargets...)
	dst.Globals.BaseReads = append(dst.Globals.BaseReads, src.Globals.BaseReads...)
	dst.Globals.Modes = append(dst.Globals.Modes, src.Globals.Modes...)
	dst.Sections = append(dst.Sections, src.Sections...)
	dst.Projections = append(dst.Projections, src.Projections...)
}

func TestUpdateAssets(t *testing.T) {
	if !*update {
		t.Skip("run with -update to rewrite golden assets and identity ledger")
	}
	sources, err := LoadAssembly(repoPath("."), assemblyPath)
	if err != nil {
		t.Fatal(err)
	}
	var ledgerSources []LedgerSource
	for _, src := range sources {
		globals := globalsSource(t)
		prog, diags := Parse(globals.Path, globals.Bytes)
		if len(diags) > 0 {
			t.Fatalf("globals parse: %v", diags)
		}
		if src.Path != globalsPath {
			parsed, diags := Parse(src.Path, src.Bytes)
			if len(diags) > 0 {
				t.Fatalf("%s parse: %v", src.Path, diags)
			}
			mergeInto(prog, parsed)
		}
		if semDiags, _ := resolve(prog); len(semDiags) > 0 {
			t.Fatalf("%s resolve: %v", src.Path, semDiags)
		}
		ledgerSources = append(ledgerSources, LedgerSource{
			Path:              src.Path,
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
	doc, eirBytes, _, err := BuildEIRFromAssembly(t, sources)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoPath(goldenEIRPath), append(eirBytes, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	fingerprint := sourceFingerprint(sources)
	prompt, err := RenderPrompt(doc, "LIGHT", "prompt", Provenance{
		SourcePath:        assemblyPath,
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

// BuildEIRFromAssembly parses and resolves the ordered sources, then builds EIR.
func BuildEIRFromAssembly(t *testing.T, sources []Source) (*EIRDoc, []byte, string, error) {
	t.Helper()
	var prog Program
	var diags Diagnostics
	for _, src := range sources {
		parsed, d := Parse(src.Path, src.Bytes)
		diags = append(diags, d...)
		if parsed != nil {
			mergeInto(&prog, parsed)
		}
	}
	if len(diags) > 0 {
		t.Fatalf("parse: %v", diags)
	}
	if semDiags, _ := resolve(&prog); len(semDiags) > 0 {
		t.Fatalf("resolve: %v", semDiags)
	}
	return BuildEIR(&prog, Versions{}, sourceFingerprint(sources))
}

func TestGoldenEIR(t *testing.T) {
	res := compileAssembly(t)
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
	res := compileAssembly(t)
	prompt, err := RenderPrompt(res.EIR, "LIGHT", "prompt", Provenance{
		SourcePath:        assemblyPath,
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
	res := compileSources(t, typedIDPath)
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
	res := compileSources(t, semanticPayloadPath)
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
	res := compileSources(t, canonicalProfilePath)
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
	res := compileSources(t, acquisitionPath)
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
	res := compileSources(t, normalizedMapsPath)
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

func TestInvocationModesCompile(t *testing.T) {
	res := compileSources(t, invocationModesPath)
	base := res.EIR.Declarations.BaseReads
	if len(base) != 9 {
		t.Fatalf("base reads has %d entries, want 9", len(base))
	}
	modes := res.EIR.Declarations.Modes
	if len(modes) != 13 {
		t.Fatalf("%d modes, want 13", len(modes))
	}
	strict := 0
	for _, mode := range modes {
		if mode.Strict {
			strict++
		}
		if len(mode.Reads) < len(base) {
			t.Fatalf("mode %s reads fewer than the base set", mode.ID)
		}
	}
	if strict != 3 {
		t.Fatalf("%d strict modes, want 3", strict)
	}
}

func TestModeDeclarationChecks(t *testing.T) {
	valid := "GLOBAL DECLARATIONS\n  ARTIFACT 0A-PREFLIGHT.md\nEND\n" +
		"BASE-READS\n  0A-PREFLIGHT.md\nEND\n" +
		"MODE Preflight\nEND\n"
	if _, err := Compile(CompileInput{
		Sources: []Source{{Path: "inline.cdl", Bytes: []byte(valid)}},
		Ledger:  map[string]IDRecord{},
	}); err != nil {
		t.Fatalf("valid mode source rejected: %v", err)
	}
	cases := []struct{ name, in, code string }{
		{
			name: "base read unresolved",
			in:   strings.Replace(valid, "  0A-PREFLIGHT.md\nEND", "  missing.md\nEND", 1),
			code: "CDL_UNRESOLVED_REFERENCE",
		},
		{
			name: "mode read unresolved",
			in:   strings.Replace(valid, "MODE Preflight\nEND", "MODE Preflight\n  READS missing.md\nEND", 1),
			code: "CDL_UNRESOLVED_REFERENCE",
		},
		{
			name: "duplicate mode",
			in:   valid + "\nMODE Preflight\nEND\n",
			code: "CDL_DUPLICATE_ID",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Compile(CompileInput{
				Sources: []Source{{Path: "inline.cdl", Bytes: []byte(tc.in)}},
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

func TestColdResumeSectionCompiles(t *testing.T) {
	res := compileSources(t, coldResumePath)
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "8.5" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	fields := 0
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "COLD-RESUME-FIELD" {
			fields = len(e.Values)
		}
	}
	if fields != 12 {
		t.Fatalf("COLD-RESUME-FIELD has %d values, want 12", fields)
	}
}

func TestTicketFSMCompiles(t *testing.T) {
	res := compileSources(t, ticketFSMPath)
	if len(res.EIR.Sections) != 2 || res.EIR.Sections[0].Number != "9.1" || res.EIR.Sections[1].Number != "9.2" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	transitions := 0
	states := 0
	for _, st := range res.EIR.Declarations.States {
		if st.ID == "TICKET-STATE" {
			transitions = len(st.Allows)
			states = len(st.Values)
		}
	}
	if states != 9 || transitions != 10 {
		t.Fatalf("TICKET-STATE has %d values and %d transitions, want 9 and 10", states, transitions)
	}
}

func TestInvocationContextCompiles(t *testing.T) {
	res := compileSources(t, invocationContextPath)
	numbers := []string{}
	for _, sec := range res.EIR.Sections {
		numbers = append(numbers, sec.Number)
	}
	if strings.Join(numbers, ",") != "8.1,8.2,8.6,8.7,8.8" {
		t.Fatalf("unexpected sections: %v", numbers)
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["INVOCATION-HEADER"] != 21 || counts["INVOCATION-LOG"] != 15 {
		t.Fatalf("unexpected field counts: %v", counts)
	}
}

func TestCoverageAndExitsCompile(t *testing.T) {
	res := compileSources(t, coverageExitsPath)
	numbers := []string{}
	for _, sec := range res.EIR.Sections {
		numbers = append(numbers, sec.Number)
	}
	if strings.Join(numbers, ",") != "10.1,10.2,10.3,10.4,10.5,10.6" {
		t.Fatalf("unexpected sections: %v", numbers)
	}
	counts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		counts[e.ID] = len(e.Values)
	}
	if counts["EXIT-A-CONDITION"] != 12 || counts["TRAVERSAL-CELL-STATE"] != 7 {
		t.Fatalf("unexpected enum counts: %v", counts)
	}
	fieldCount := 0
	for _, f := range res.EIR.Declarations.Fields {
		if f.ID == "SWEEP-RECORD" {
			fieldCount = len(f.Fields)
		}
	}
	if fieldCount != 11 {
		t.Fatalf("SWEEP-RECORD has %d fields, want 11", fieldCount)
	}
}

func TestLedgerRejectsUnknownIdentity(t *testing.T) {
	ledger, err := LoadLedger(repoPath(ledgerPath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = Compile(CompileInput{Sources: []Source{globalsSource(t), exampleSource(t)}, Ledger: ledger})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	trimmed := map[string]IDRecord{}
	for k, v := range ledger {
		if k != "reconcile.guard" {
			trimmed[k] = v
		}
	}
	_, err = Compile(CompileInput{Sources: []Source{globalsSource(t), exampleSource(t)}, Ledger: trimmed})
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
				Sources: []Source{globalsSource(t), {Path: examplePath, Bytes: []byte(tc.mutate(base))}},
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
