// Package workspace creates and checks the non-fabricated Protocol v4
// Markdown workspace skeleton described by Parts 1 and 3.
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

var forensicFiles = []string{
	"0A-PREFLIGHT.md", "0B-AUTH-MODEL.md", "0C-SECURITY-PRIVACY.md", "0D-GLOSSARY.md",
	"0E-INDEX.md", "0F-GLOBAL-STATE.md", "0G-DECONSTRUCTION-STATE.md", "0H-CHECKPOINT-SUMMARY.md",
}
var assuranceFiles = []string{
	"10-SOURCE-INVENTORY.md", "11-TRAVERSAL-FRONTIER.md", "12-CLAIM-EVIDENCE.md", "13-EXPORT-RECONCILIATION.md",
	"14-CONTRADICTIONS.md", "15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "17-ACQUISITION-CANDIDATES.md",
	"18-CONFIRMATIONS.md", "20-TRACEABILITY.md", "21-COVERAGE-REPORT.md", "22-GATE-REPORTS.md",
}
var synthesisFiles = []string{
	"90-ARCH-BLUEPRINT.md", "91-DATA-MODEL.md", "92-BUSINESS-RULES.md", "93-USE-CASES.md",
	"94-INTERFACES.md", "95-DEPLOYMENT.md", "96-NON-FUNCTIONAL-SECURITY.md",
}
var handbookFiles = []string{
	"START-HERE.md", "01-SYSTEM-OVERVIEW.md", "02-ARCHITECTURE.md", "03-DOMAIN-AND-DATA.md",
	"04-PERSONAS-AUTH-AND-VISIBILITY.md", "05-USE-CASES.md", "06-BUSINESS-RULES.md",
	"07-INTERFACES-AND-INTEGRATIONS.md", "08-BACKGROUND-PROCESSING.md", "09-DEPLOYMENT-AND-CONFIGURATION.md",
	"10-SECURITY-PRIVACY-AND-AUDIT.md", "11-NON-FUNCTIONAL-PROFILE.md",
	"12-DISABLED-AND-DORMANT-CAPABILITIES.md", "13-MODERNIZATION-DECISIONS.md",
	"14-EQUIVALENCE-AND-ACCEPTANCE.md", "15-KNOWN-GAPS-AND-RISKS.md", "16-GLOSSARY-AND-REFERENCE.md",
}

var povNames = []string{"FRONTEND", "BACKEND", "BACKGROUND", "DATA", "WIRING"}

var personaDirectoryPattern = regexp.MustCompile(`^([A-Z][A-Z0-9]{1,7})-([a-z0-9]+(?:-[a-z0-9]+)*)$`)

type PersonaDirectory struct {
	Prefix string
	Slug   string
}

func ParsePersonaDirectory(name string) (PersonaDirectory, error) {
	if name == "_shared" {
		return PersonaDirectory{}, fmt.Errorf("%q is the reserved non-persona directory", name)
	}
	matches := personaDirectoryPattern.FindStringSubmatch(name)
	if matches == nil {
		return PersonaDirectory{}, fmt.Errorf("persona directory %q must match <persona-prefix>-<persona-slug>", name)
	}
	return PersonaDirectory{Prefix: matches[1], Slug: matches[2]}, nil
}

func RequiredFiles() []string {
	var files []string
	files = append(files, forensicFiles...)
	files = append(files, assuranceFiles...)
	files = append(files, synthesisFiles...)
	for i, name := range povNames {
		prefix := fmt.Sprintf("%02d-%s", i+1, name)
		files = append(files, filepath.ToSlash(filepath.Join("personas", "_shared", prefix+"-POV.md")))
		files = append(files, filepath.ToSlash(filepath.Join("personas", "_shared", prefix+"-QUESTIONS.md")))
	}
	for _, name := range handbookFiles {
		files = append(files, filepath.ToSlash(filepath.Join("handbook", name)))
	}
	sort.Strings(files)
	return files
}

func Init(root string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if _, err := os.Stat(root); err == nil {
		return fmt.Errorf("workspace already exists: %s", root)
	} else if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(root)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(parent, ".legacy-autopsy-init-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	for _, dir := range []string{"personas/_shared", "normalized", "probes", "handbook"} {
		if err := os.MkdirAll(filepath.Join(stage, filepath.FromSlash(dir)), 0o755); err != nil {
			return err
		}
	}
	for _, rel := range RequiredFiles() {
		if err := writeScaffold(stage, rel); err != nil {
			return err
		}
	}
	if err := os.Rename(stage, root); err != nil {
		return fmt.Errorf("publish workspace atomically: %w", err)
	}
	return nil
}

func writeScaffold(root, rel string) error {
	title := strings.TrimSuffix(filepath.Base(rel), ".md")
	plane := "Forensic Plane"
	switch {
	case strings.HasPrefix(rel, "handbook/"):
		plane = "Human Reconstruction Handbook"
	case strings.HasPrefix(filepath.Base(rel), "1"), strings.HasPrefix(filepath.Base(rel), "2"):
		plane = "Assurance Plane"
	case strings.HasPrefix(filepath.Base(rel), "9"):
		plane = "Synthesis Plane"
	}
	content := fmt.Sprintf("# %s\n\n> Initialization scaffold — %s. No source facts, evidence claims, coverage, confirmations, or gate outcomes have been asserted.\n", title, plane)
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func Check(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("workspace root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("workspace root is not a directory: %s", root)
	}
	var problems []string
	personaRoot := filepath.Join(root, "personas")
	personaEntries, err := os.ReadDir(personaRoot)
	if err != nil {
		problems = append(problems, fmt.Sprintf("personas: %v", err))
	} else {
		for _, entry := range personaEntries {
			if entry.Name() == "_shared" {
				if !entry.IsDir() {
					problems = append(problems, "personas/_shared: must be a directory")
				}
				continue
			}
			if !entry.IsDir() {
				problems = append(problems, fmt.Sprintf("personas/%s: persona container entries must be directories", entry.Name()))
				continue
			}
			if _, err := ParsePersonaDirectory(entry.Name()); err != nil {
				problems = append(problems, fmt.Sprintf("personas/%s: %v", entry.Name(), err))
			}
		}
	}
	for _, rel := range RequiredFiles() {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", rel, err))
			continue
		}
		doc, err := md.Parse(string(data))
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: invalid Markdown: %v", rel, err))
			continue
		}
		expectedHeading := strings.TrimSuffix(filepath.Base(rel), ".md")
		headings := 0
		for _, node := range doc.Nodes {
			if node.Kind == md.Heading && node.Level == 1 {
				headings++
				if node.Text != expectedHeading {
					problems = append(problems, fmt.Sprintf("%s: H1 must be %q, found %q", rel, expectedHeading, node.Text))
				}
			}
		}
		if headings != 1 {
			problems = append(problems, fmt.Sprintf("%s: expected exactly one H1, found %d", rel, headings))
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("workspace validation failed:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
}
