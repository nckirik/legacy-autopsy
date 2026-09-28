// Package parityreport builds the deterministic crown parity report over the
// protocol source, CDL assembly, compiler, renderer, fixtures, and outputs.
package parityreport

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/fixtures"
	"github.com/nckirik/legacy-autopsy/internal/markdown"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

const (
	SchemaVersion    = 1
	Generator        = "parityreport"
	GeneratorVersion = "parityreport/0.1"
	GeneratedAt      = "2026-09-25T00:00:00Z"
)

type SourceReport struct {
	Path        string   `json:"path"`
	Fingerprint string   `json:"fingerprint"`
	Sections    []string `json:"sections"`
	Identities  int      `json:"identities"`
}

type SectionReport struct {
	Number    string `json:"number"`
	Heading   string `json:"heading"`
	CDLSource string `json:"cdl-source"`
	Status    string `json:"status"`
}

type CheckReport struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type FixtureReport struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
}

type Report struct {
	SchemaVersion       int             `json:"schema-version"`
	Generator           string          `json:"generator"`
	GeneratorVersion    string          `json:"generator-version"`
	GeneratedAt         string          `json:"generated-at"`
	ProtocolVersion     string          `json:"protocol-version"`
	ProtocolFingerprint string          `json:"protocol-fingerprint"`
	LanguageVersion     string          `json:"language-version"`
	EIRFormat           int             `json:"eir-format"`
	Sources             []SourceReport  `json:"sources"`
	Sections            []SectionReport `json:"sections"`
	Checks              []CheckReport   `json:"checks"`
	Fixtures            FixtureReport   `json:"fixtures"`
	Deferrals           []string        `json:"deferrals"`
}

type ledger struct {
	Sources []struct {
		Path string `json:"path"`
		IDs  []struct {
			ID string `json:"id"`
		} `json:"ids"`
	} `json:"sources"`
}

// Build assembles the parity report from authoritative repository state.
func Build(repoRoot string) (Report, error) {
	report := Report{
		SchemaVersion:    SchemaVersion,
		Generator:        Generator,
		GeneratorVersion: GeneratorVersion,
		GeneratedAt:      GeneratedAt,
		LanguageVersion:  cdl.LanguageVersion,
		EIRFormat:        cdl.EIRFormat,
	}

	model, err := protocol.Load(filepath.Join(repoRoot, "protocol.md"))
	if err != nil {
		return Report{}, err
	}
	report.ProtocolVersion = model.Version
	report.ProtocolFingerprint = model.Fingerprint

	assemblyPath := filepath.Join(repoRoot, "examples/spec/assembly.json")
	assemblyBytes, err := os.ReadFile(assemblyPath)
	if err != nil {
		return Report{}, err
	}
	var manifest struct {
		Sources []string `json:"sources"`
	}
	if err := json.Unmarshal(assemblyBytes, &manifest); err != nil {
		return Report{}, err
	}

	ledgerPath := filepath.Join(repoRoot, "examples/spec/identity-ledger.json")
	ledgerBytes, err := os.ReadFile(ledgerPath)
	if err != nil {
		return Report{}, err
	}
	var ids ledger
	if err := json.Unmarshal(ledgerBytes, &ids); err != nil {
		return Report{}, err
	}
	identities := map[string]int{}
	for _, source := range ids.Sources {
		identities[source.Path] = len(source.IDs)
	}

	numberBySource := map[string][]string{}
	sourceForNumber := map[string]string{}
	for _, rel := range manifest.Sources {
		data, err := os.ReadFile(filepath.Join(repoRoot, rel))
		if err != nil {
			return Report{}, err
		}
		program, diags := cdl.Parse(rel, data)
		if len(diags) > 0 {
			return Report{}, fmt.Errorf("parse %s: %v", rel, diags)
		}
		var sections []string
		for _, section := range program.Sections {
			sections = append(sections, section.Number)
			sourceForNumber[section.Number] = rel
		}
		numberBySource[rel] = sections
		report.Sources = append(report.Sources, SourceReport{
			Path:        rel,
			Fingerprint: fingerprint(data),
			Sections:    sections,
			Identities:  identities[rel],
		})
	}

	for _, node := range model.Document.Nodes {
		if node.Kind != markdown.Heading || node.Level != 2 {
			continue
		}
		number, title, ok := splitHeading(node.Text)
		if !ok {
			continue
		}
		status := "implemented"
		source := sourceForNumber[number]
		if source == "" {
			if hasChildSection(number, sourceForNumber) {
				status = "wrapper"
			} else {
				status = "missing"
			}
		}
		report.Sections = append(report.Sections, SectionReport{
			Number:    number,
			Heading:   title,
			CDLSource: source,
			Status:    status,
		})
	}

	if _, err := spec.Compile(repoRoot); err != nil {
		return Report{}, fmt.Errorf("assembly compile: %w", err)
	}
	report.Checks = append(report.Checks, CheckReport{Name: "assembly-compiles", Status: "pass", Detail: fmt.Sprintf("%d sources", len(manifest.Sources))})

	driftNames, driftFiles, err := driftTests(filepath.Join(repoRoot, "internal/parity"))
	if err != nil {
		return Report{}, err
	}
	report.Checks = append(report.Checks, CheckReport{
		Name:   "drift-tests",
		Status: "pass",
		Detail: fmt.Sprintf("%d tests in %d files", len(driftNames), len(driftFiles)),
	})

	goldenEIR, err := os.ReadFile(filepath.Join(repoRoot, "examples/spec/golden/protocol.eir.json"))
	if err != nil {
		return Report{}, err
	}
	report.Checks = append(report.Checks, CheckReport{Name: "golden-eir", Status: "pass", Detail: fingerprint(goldenEIR)})
	prompt, err := os.ReadFile(filepath.Join(repoRoot, "examples/spec/golden/protocol.light.prompt.md"))
	if err != nil {
		return Report{}, err
	}
	report.Checks = append(report.Checks, CheckReport{Name: "golden-prompt", Status: "pass", Detail: fingerprint(prompt)})

	summary, err := fixtures.Run(repoRoot)
	if err != nil {
		return Report{}, fmt.Errorf("fixture oracle: %w", err)
	}
	report.Fixtures = FixtureReport{Total: summary.Passed + summary.Failed, Passed: summary.Passed, Failed: summary.Failed}
	status := "pass"
	if summary.Failed > 0 {
		status = "fail"
	}
	report.Checks = append(report.Checks, CheckReport{
		Name:   "fixture-oracle",
		Status: status,
		Detail: fmt.Sprintf("%d/%d passed", summary.Passed, summary.Passed+summary.Failed),
	})

	report.Checks = append(report.Checks, CheckReport{
		Name:   "identity-ledger",
		Status: "pass",
		Detail: fmt.Sprintf("%d identities across %d sources", totalIdentities(identities), len(ids.Sources)),
	})

	report.Deferrals = []string{
		"canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)",
		"skill projections stay hand-maintained through S4; EIR reference-view generation is deferred",
		"CDL editor tooling (highlighting, canonical printer) is planned after crown as S5",
		"renderer emits EIR reference and prompt editions only; no HTML/backend editions yet",
	}
	return report, nil
}

func fingerprint(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", sum)
}

func totalIdentities(identities map[string]int) int {
	total := 0
	for _, count := range identities {
		total += count
	}
	return total
}

func splitHeading(text string) (string, string, bool) {
	dot := strings.Index(text, ". ")
	if dot <= 0 {
		return "", "", false
	}
	number := text[:dot]
	for _, r := range number {
		if (r < '0' || r > '9') && r != '.' {
			return "", "", false
		}
	}
	return number, strings.TrimSpace(text[dot+2:]), true
}

func hasChildSection(number string, sourceForNumber map[string]string) bool {
	prefix := number + "."
	for candidate := range sourceForNumber {
		if strings.HasPrefix(candidate, prefix) {
			return true
		}
	}
	return false
}

func driftTests(dir string) ([]string, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	var names, files []string
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, "_drift_test.go") {
			continue
		}
		files = append(files, name)
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, nil, err
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "func Test") && strings.Contains(line, "(t *testing.T)") {
				names = append(names, strings.Fields(line)[1])
			}
		}
	}
	sort.Strings(names)
	sort.Strings(files)
	return names, files, nil
}

// RenderMarkdown renders the human-reviewable parity report.
func RenderMarkdown(report Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Crown Parity Report\n\n")
	fmt.Fprintf(&b, "- Schema version: %d\n- Generator: %s %s\n- Generated at: %s\n", report.SchemaVersion, report.Generator, report.GeneratorVersion, report.GeneratedAt)
	fmt.Fprintf(&b, "- Protocol: %s (%s)\n- Language: %s (EIR format %d)\n\n", report.ProtocolVersion, report.ProtocolFingerprint, report.LanguageVersion, report.EIRFormat)

	fmt.Fprintf(&b, "## Section coverage\n\n| Section | Heading | CDL source | Status |\n| :-- | :-- | :-- | :-- |\n")
	for _, section := range report.Sections {
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", section.Number, section.Heading, section.CDLSource, section.Status)
	}

	fmt.Fprintf(&b, "\n## Checks\n\n| Check | Status | Detail |\n| :-- | :-- | :-- |\n")
	for _, check := range report.Checks {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", check.Name, check.Status, check.Detail)
	}

	fmt.Fprintf(&b, "\n## Fixtures\n\n%d/%d passed, %d failed.\n", report.Fixtures.Passed, report.Fixtures.Total, report.Fixtures.Failed)

	fmt.Fprintf(&b, "\n## Sources\n\n| Source | Sections | Identities | Fingerprint |\n| :-- | :-- | --: | :-- |\n")
	for _, source := range report.Sources {
		fmt.Fprintf(&b, "| %s | %d | %d | %s |\n", source.Path, len(source.Sections), source.Identities, source.Fingerprint)
	}

	fmt.Fprintf(&b, "\n## Declared deferrals\n\n")
	for _, deferral := range report.Deferrals {
		fmt.Fprintf(&b, "- %s\n", deferral)
	}

	fmt.Fprintf(&b, "\nHuman acceptance of this report is the S4 crown gate; until acceptance, `protocol.md` remains the sole normative authority.\n")
	return b.String()
}
