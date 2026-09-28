package cdl

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden assets and identity ledger")

const (
	examplePath           = "protocol/export-reconciliation.cdl"
	typedIDPath           = "protocol/typed-id.cdl"
	semanticPayloadPath   = "protocol/semantic-payload-identity.cdl"
	canonicalProfilePath  = "protocol/canonical-hash-profile.cdl"
	acquisitionPath       = "protocol/export-acquisition-loop.cdl"
	normalizedMapsPath    = "protocol/normalized-maps.cdl"
	invocationModesPath   = "protocol/invocation-modes.cdl"
	coldResumePath        = "protocol/cold-resume.cdl"
	ticketFSMPath         = "protocol/ticket-fsm.cdl"
	invocationContextPath = "protocol/invocation-context.cdl"
	coverageExitsPath     = "protocol/coverage-and-exits.cdl"
	statusTaxonomyPath    = "protocol/status-taxonomy.cdl"
	evidenceDecisionsPath = "protocol/evidence-and-decisions.cdl"
	workspaceRegistryPath = "protocol/workspace-registries.cdl"
	preflightRegistryPath = "protocol/preflight-registry.cdl"
	acquisitionTrustPath  = "protocol/acquisition-trust.cdl"
	traceabilityPath      = "protocol/traceability.cdl"
	packagingPath         = "protocol/packaging.cdl"
	foundationsPath       = "protocol/foundations.cdl"
	personasPovsPath      = "protocol/personas-povs.cdl"
	inventoryFrontierPath = "protocol/inventory-and-frontier.cdl"
	extractionEvoPath     = "protocol/extraction-evolution.cdl"
	handbookDecisionsPath = "protocol/handbook-and-decisions.cdl"
	goldenProtocolPath    = "protocol/golden/protocol.generated.md"
	synthesisCatalogsPath = "protocol/synthesis-catalogs.cdl"
	migrationConfPath     = "protocol/migration-and-conformance.cdl"
	globalsPath           = "protocol/globals.cdl"
	assemblyPath          = "protocol/assembly.json"
	ledgerPath            = "protocol/identity-ledger.json"
	goldenEIRPath         = "protocol/golden/protocol.eir.json"
	goldenPromptPath      = "protocol/golden/protocol.light.prompt.md"
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
	seen := map[string]bool{globalsPath: true}
	for _, rel := range rels {
		if seen[rel] {
			continue
		}
		seen[rel] = true
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
	dst.Globals.Parts = append(dst.Globals.Parts, src.Globals.Parts...)
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
	sharedProg, sharedDiags := Parse(statusTaxonomyPath, readRepoFile(t, statusTaxonomyPath))
	if len(sharedDiags) > 0 {
		t.Fatalf("status taxonomy parse: %v", sharedDiags)
	}
	sharedIDs := map[string]bool{}
	for _, id := range collectIDs(sharedProg) {
		sharedIDs[id.ID] = true
	}
	var ledgerSources []LedgerSource
	for _, src := range sources {
		globals := globalsSource(t)
		prog, diags := Parse(globals.Path, globals.Bytes)
		if len(diags) > 0 {
			t.Fatalf("globals parse: %v", diags)
		}
		if src.Path != globalsPath && src.Path != statusTaxonomyPath {
			shared, diags := Parse(statusTaxonomyPath, readRepoFile(t, statusTaxonomyPath))
			if len(diags) > 0 {
				t.Fatalf("status taxonomy parse: %v", diags)
			}
			mergeInto(prog, shared)
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
		ids := collectIDs(prog)
		if src.Path != statusTaxonomyPath {
			var filtered []IDRecord
			for _, id := range ids {
				if !sharedIDs[id.ID] {
					filtered = append(filtered, id)
				}
			}
			ids = filtered
		}
		ledgerSources = append(ledgerSources, LedgerSource{
			Path:              src.Path,
			SourceFingerprint: sourceFingerprint([]Source{src}),
			IDs:               ids,
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
	rendered := RenderProtocol(doc)
	if err := os.WriteFile(repoPath(goldenProtocolPath), rendered, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoPath("protocol.md"), rendered, 0o644); err != nil {
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

func TestEverySectionHasGoal(t *testing.T) {
	res := compileAssembly(t)
	for _, section := range res.EIR.Sections {
		if strings.TrimSpace(section.Goal) == "" {
			t.Fatalf("section %s (%s) has no GOAL description", section.Number, section.Title)
		}
	}
}

func TestGoldenProtocolRender(t *testing.T) {
	res := compileAssembly(t)
	got := RenderProtocol(res.EIR)
	for _, path := range []string{goldenProtocolPath, "protocol.md"} {
		want, err := os.ReadFile(repoPath(path))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("%s is stale (run: go test ./cdl -run TestUpdateAssets -update)", path)
		}
	}
	for _, marker := range []string{"**Version:** 4.1.2", "## 8.3. Invocation-mode enum and explicit multi-file modes", "Preflight | Export Acquisition"} {
		if !strings.Contains(string(got), marker) {
			t.Fatalf("generated protocol missing %q", marker)
		}
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
	numbers := []string{}
	for _, sec := range res.EIR.Sections {
		numbers = append(numbers, sec.Number)
	}
	if strings.Join(numbers, ",") != "9.1,9.2,9.3,9.4,9.5,9.6" {
		t.Fatalf("unexpected sections: %v", numbers)
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
	if strings.Join(numbers, ",") != "8.1,8.2,8.3,8.4,8.6,8.7,8.8" {
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

func TestStatusTaxonomyCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath)
	if len(res.EIR.Sections) != 1 || res.EIR.Sections[0].Number != "5.1" {
		t.Fatalf("unexpected sections: %+v", res.EIR.Sections)
	}
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}
	if len(enums) != 109 {
		t.Fatalf("%d enums, want 109", len(enums))
	}
	if len(enums["PREFIX-STAGING"]) != 5 {
		t.Fatalf("PREFIX-STAGING has %d values", len(enums["PREFIX-STAGING"]))
	}
	if len(enums["REGISTRY-TICKET-STATUS"]) != 9 {
		t.Fatalf("REGISTRY-TICKET-STATUS has %d values", len(enums["REGISTRY-TICKET-STATUS"]))
	}
	if len(enums["REGISTRY-CONFIGURATION-MEDIUM"]) != 5 {
		t.Fatalf("REGISTRY-CONFIGURATION-MEDIUM has %d values", len(enums["REGISTRY-CONFIGURATION-MEDIUM"]))
	}
}

func TestEvidenceAndDecisionsCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, evidenceDecisionsPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"5.2", "5.3", "5.4", "5.5", "5.6", "5.7"} {
		if !numbers[want] {
			t.Fatalf("section %s missing from %v", want, numbers)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["CLAIM"] != 13 || counts["BLOCK-EVIDENCE-SUMMARY"] != 2 ||
		counts["CONTRADICTION"] != 6 || counts["DECISION"] != 13 || counts["CONFIRMATION"] != 11 {
		t.Fatalf("unexpected field counts: %v", counts)
	}
}

func TestWorkspaceRegistriesCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, workspaceRegistryPath)
	found := false
	for _, sec := range res.EIR.Sections {
		if sec.Number == "3.2" {
			found = true
		}
	}
	if !found {
		t.Fatal("section 3.2 missing")
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{
		"SHARED-ENTRY": 6, "AUTH-MECHANISM": 7, "PRIVACY-FINDING": 6,
		"CAPABILITY-RECORD": 10, "ROUTINE-CROSS-REFERENCE": 5, "INVARIANT": 3, "SHARED-STATE": 4,
	}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
}

func TestPreflightRegistryCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, preflightRegistryPath)
	found := false
	for _, sec := range res.EIR.Sections {
		if sec.Number == "3.1" {
			found = true
		}
	}
	if !found {
		t.Fatal("section 3.1 missing")
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{"PREFLIGHT-HEADER": 9, "ENTRY-CLUSTER": 10, "UNMAPPED-DISCOVERY": 4, "DISCOVERY-BUFFER": 7, "PERSONA-REGISTRY": 8, "REQUIRED-TRAVERSAL-MATRIX": 10}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
}

func TestAcquisitionTrustCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, acquisitionTrustPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"7.1", "7.2", "7.3", "7.4"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["EXPORT"] != 13 {
		t.Fatalf("EXPORT has %d fields, want 13", counts["EXPORT"])
	}
	enums := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = len(e.Values)
	}
	if enums["N8N-INVENTORY-KIND"] != 17 || enums["APPSMITH-INVENTORY-KIND"] != 17 {
		t.Fatalf("unexpected inventory enum counts: %v", enums)
	}
}

func TestTraceabilityCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, traceabilityPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"12.1", "12.2", "12.3", "12.4", "12.5"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{"TRACEABILITY-ROW": 14, "LIFECYCLE-ELIGIBILITY": 6, "EVIDENCE-BINDING": 7}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "CLASSIFICATION" && len(e.Values) != 3 {
			t.Fatalf("CLASSIFICATION has %d values, want 3", len(e.Values))
		}
	}
}

func TestPackagingCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, packagingPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"15.1", "15.1.1", "15.1.2", "15.1.3", "15.1.4", "15.2", "15.3", "15.4", "15.5"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{
		"EXIT-E-CANDIDATE-REPORT":         16,
		"CANDIDATE-PAYLOAD-MANIFEST":      9,
		"EXIT-E-CONTENT-READINESS-REPORT": 16,
		"SCOPE-CERTIFICATE":               15,
		"OUTER-BUNDLE-MANIFEST":           11,
		"FINAL-VERIFICATION-RECEIPT":      9,
		"ARTIFACT-SIGNATURE-ENVELOPE":     2,
		"OUTER-SIGNATURE-ENVELOPE":        2,
	}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
	enumCounts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enumCounts[e.ID] = len(e.Values)
	}
	if enumCounts["EXIT-E-CANDIDATE-CHECKS"] != 14 || enumCounts["EXIT-E-CONTENT-READINESS-CHECKS"] != 11 || enumCounts["FINAL-CHECK"] != 9 {
		t.Fatalf("unexpected check registry counts: %v", enumCounts)
	}
}

func TestPersonasPovsCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, personasPovsPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"2.1", "2.2", "2.3", "2.4", "2.5", "2.6"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["CAPABILITY-STATE-BLOCK"] != 5 || counts["RUNTIME-STATE-MATRIX"] != 6 {
		t.Fatalf("unexpected capability field counts: %v", counts)
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "TARGET-DECISION" && len(e.Values) != 4 {
			t.Fatalf("TARGET-DECISION has %d values, want 4", len(e.Values))
		}
	}
}

func TestInventoryAndFrontierCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, inventoryFrontierPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"4.2", "4.3", "4.4", "4.5", "4.6"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{
		"SOURCE-INVENTORY-RECORD": 15,
		"FRONTIER-RECORD":         11,
		"SOURCE-COVERAGE-ROW":     12,
		"SUPERSESSION-TOMBSTONE":  5,
	}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
	enumCounts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enumCounts[e.ID] = len(e.Values)
	}
	wantEnums := map[string]int{
		"DISCOVERY-METHOD":          8,
		"COVERAGE-SUMMARY":          5,
		"INVENTORY-KIND":            23,
		"FRONTIER-DISCOVERY-METHOD": 9,
		"FRONTIER-STATE":            7,
		"COVERAGE-SCOPE":            3,
		"COVERAGE-DISPOSITION":      6,
	}
	for id, size := range wantEnums {
		if enumCounts[id] != size {
			t.Fatalf("%s has %d values, want %d", id, enumCounts[id], size)
		}
	}
}

func TestExtractionEvolutionCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, extractionEvoPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"6.1", "6.2", "6.3", "6.4", "6.5"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["ATOMIC-COMPONENT"] != 25 || counts["SHARED-REFERENCE"] != 6 {
		t.Fatalf("unexpected extraction field counts: %v", counts)
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "BLOCK-CONFIDENCE-RANK" && len(e.Values) != 4 {
			t.Fatalf("BLOCK-CONFIDENCE-RANK has %d values, want 4", len(e.Values))
		}
	}
}

func TestModeReadSetsCompiles(t *testing.T) {
	res := compileSources(t, invocationContextPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"8.3", "8.4"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
}

func TestHumanHatchProbeCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, ticketFSMPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"9.3", "9.4", "9.5", "9.6"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["PROBE-SPECIFICATION"] != 10 || counts["NO-MOCK-FALLBACK"] != 7 {
		t.Fatalf("unexpected probe field counts: %v", counts)
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "SAFETY-CLASSIFICATION" && len(e.Values) != 2 {
			t.Fatalf("SAFETY-CLASSIFICATION has %d values, want 2", len(e.Values))
		}
	}
}

func TestHandbookAndDecisionsCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, handbookDecisionsPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"13.1", "13.2", "13.3", "13.4", "14.1", "14.2", "14.3", "14.4", "14.5", "14.6"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	enumCounts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enumCounts[e.ID] = len(e.Values)
	}
	if enumCounts["AUDIENCE-ROLE"] != 6 || enumCounts["EQUIVALENCE-TEST-KIND"] != 6 {
		t.Fatalf("unexpected handbook enum counts: %v", enumCounts)
	}
}

func TestSynthesisCatalogsCompile(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, extractionEvoPath, synthesisCatalogsPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for i := 1; i <= 13; i++ {
		want := fmt.Sprintf("11.%d", i)
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	want := map[string]int{
		"SYNTHESIS-GAP": 9, "SYNTHESIS-BLOCK-HEADER": 8, "CERTIFICATION-ENVELOPE": 4,
		"ARCHITECTURE-MODULE": 10, "ENTITY-RECORD": 15, "ENTITY-COLUMN": 7,
		"RELATIONSHIP-RECORD": 6, "STATE-MACHINE": 7, "STATE-TRANSITION": 6,
		"DATABASE-ROUTINE": 12, "BUSINESS-RULE-RECORD": 12, "USE-CASE": 11,
		"USE-CASE-STEP": 6, "INTERFACE-RECORD": 14, "INTERFACE-REQUEST-FIELD": 6,
		"INTERFACE-RESPONSE-FIELD": 4, "DEPLOYMENT-ELEMENT": 4, "CONFIGURATION-RECORD": 7,
		"SCHEDULE-RECORD": 6, "NFR-RECORD": 5, "SECURITY-FINDING": 6, "FAULT-RECORD": 8,
		"PERSONA-PROFILE": 20, "AUTHORIZATION-MATRIX-ROW": 4,
	}
	for id, size := range want {
		if counts[id] != size {
			t.Fatalf("%s has %d fields, want %d", id, counts[id], size)
		}
	}
	enumCounts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enumCounts[e.ID] = len(e.Values)
	}
	wantEnums := map[string]int{
		"GAP-MATERIALITY": 2, "GAP-DISPOSITION": 7, "STORAGE-KIND": 6,
		"RELATIONSHIP-CARDINALITY": 3, "RELATIONSHIP-ENFORCEMENT": 4,
		"DELETE-UPDATE-SEMANTICS": 4, "DB-ROUTINE-KIND": 7, "DB-TARGET-DECISION": 5,
		"VIOLATION-POLICY": 6, "BUSINESS-RULE-CAPABILITY-STATE": 2,
		"USE-CASE-CAPABILITY-CLASSIFICATION": 5, "USE-CASE-TARGET-DECISION": 4,
		"CONFIG-MEDIUM": 5, "NFR-ATTRIBUTE": 8, "SECURITY-CATEGORY": 10, "RISK-LEVEL": 4,
	}
	for id, size := range wantEnums {
		if enumCounts[id] != size {
			t.Fatalf("%s has %d values, want %d", id, enumCounts[id], size)
		}
	}
}

func TestMigrationAndConformanceCompiles(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, migrationConfPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"17.1", "17.2", "17.3", "19.1", "19.2", "19.3", "20"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	counts := map[string]int{}
	for _, f := range res.EIR.Declarations.Fields {
		counts[f.ID] = len(f.Fields)
	}
	if counts["ARTIFACT-OWNERSHIP-ROW"] != 5 {
		t.Fatalf("ARTIFACT-OWNERSHIP-ROW has %d fields, want 5", counts["ARTIFACT-OWNERSHIP-ROW"])
	}
	for _, e := range res.EIR.Declarations.Enums {
		if e.ID == "CONFORMANCE-CASE" && len(e.Values) != 14 {
			t.Fatalf("CONFORMANCE-CASE has %d values, want 14", len(e.Values))
		}
	}
}

func TestFoundationsCompile(t *testing.T) {
	res := compileSources(t, statusTaxonomyPath, foundationsPath)
	numbers := map[string]bool{}
	for _, sec := range res.EIR.Sections {
		numbers[sec.Number] = true
	}
	for _, want := range []string{"0.1", "0.2", "0.3", "0.4", "1.1", "1.2", "1.3", "1.4", "1.5"} {
		if !numbers[want] {
			t.Fatalf("section %s missing", want)
		}
	}
	enumCounts := map[string]int{}
	for _, e := range res.EIR.Declarations.Enums {
		enumCounts[e.ID] = len(e.Values)
	}
	want := map[string]int{
		"NORMATIVE-KEYWORD": 6,
		"AUTHORITY-CLASS":   3,
		"PREDICATE-KIND":    3,
		"FLOW-STAGE":        14,
		"PLANE":             4,
	}
	for id, size := range want {
		if enumCounts[id] != size {
			t.Fatalf("%s has %d values, want %d", id, enumCounts[id], size)
		}
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
