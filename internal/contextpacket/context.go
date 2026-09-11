// Package contextpacket builds bounded, non-authoritative invocation packets
// containing exact routed protocol sections and confined workspace inputs.
package contextpacket

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/identity"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
	"github.com/nckirik/legacy-autopsy/internal/routing"
	"github.com/nckirik/legacy-autopsy/internal/workspace"
)

type Options struct {
	Mode, SystemNamespace, Iteration, InvocationID, InvocationScope        string
	Persona, PersonaPrefix, PersonaDirectory, Cluster, Track, POV, POVFile string
	QuestionLedger, EnvironmentSnapshot, HumanHatchAction, Workspace       string
	MaxTraversalDepth                                                      int
	AdditionalReadTargets, AllowedWriteTargets, ForbiddenWriteTargets      []string
}

type FileBinding struct {
	Path, Fingerprint, Content string
}

type Packet struct {
	GeneratedAt           time.Time
	Projection            routing.Projection
	Protocol              *protocol.Model
	Options               Options
	Files                 []FileBinding
	ProjectionText        string
	CheckpointFingerprint string
}

var iterations = stringSet("ALFA", "BRAVO", "CHARLIE", "DELTA", "ECHO", "FOXTROT", "GOLF", "HOTEL", "INDIA", "JULIETT", "KILO", "LIMA", "MIKE", "NOVEMBER", "OSCAR", "PAPA", "QUEBEC", "ROMEO", "SIERRA", "TANGO", "UNIFORM", "VICTOR", "WHISKEY", "X-RAY", "YANKEE", "ZULU")

func stringSet(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func Build(repoRoot string, model *protocol.Model, manifest *routing.Manifest, options Options) (*Packet, error) {
	projection, ok := manifest.Projection(options.Mode)
	if !ok {
		return nil, fmt.Errorf("unknown invocation mode %q", options.Mode)
	}
	if options.SystemNamespace == "" || options.Iteration == "" || options.InvocationID == "" || options.InvocationScope == "" || options.EnvironmentSnapshot == "" {
		return nil, fmt.Errorf("system namespace, iteration, invocation ID, invocation scope, and environment/snapshot are required")
	}
	if _, ok := iterations[options.Iteration]; !ok {
		return nil, fmt.Errorf("iteration %q is not a canonical §10.6 token", options.Iteration)
	}
	if len(options.AllowedWriteTargets) == 0 || len(options.ForbiddenWriteTargets) == 0 {
		return nil, fmt.Errorf("at least one allowed and one forbidden write target are required")
	}
	if strictMode(options.Mode) {
		if options.Persona == "" || options.PersonaPrefix == "" || options.PersonaDirectory == "" || options.Cluster == "" || options.Track == "" || options.POV == "" || options.POVFile == "" || options.MaxTraversalDepth <= 0 {
			return nil, fmt.Errorf("%s requires persona, persona prefix, persona directory, cluster, track, POV, POV file, and positive traversal-depth bindings", options.Mode)
		}
		personaDirectory, err := workspace.ParsePersonaDirectory(options.PersonaDirectory)
		if err != nil {
			return nil, err
		}
		if personaDirectory.Prefix != options.PersonaPrefix {
			return nil, fmt.Errorf("persona directory prefix %q does not match bound persona prefix %q", personaDirectory.Prefix, options.PersonaPrefix)
		}
	}
	if options.Mode == "Human Hatch" && options.HumanHatchAction == "" {
		return nil, fmt.Errorf("Human Hatch requires an action subtype")
	}
	if options.Mode != "Preflight" && len(options.AdditionalReadTargets) == 0 {
		return nil, fmt.Errorf("%s requires explicit --read-target bindings for its mode-specific or transitive read closure", options.Mode)
	}

	projectionData, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(projection.SkillProjection)))
	if err != nil {
		return nil, err
	}
	readTargets := commonReadTargets()
	readTargets = append(readTargets, staticModeReadTargets(options.Mode)...)
	if options.POVFile != "" {
		readTargets = append(readTargets, options.POVFile)
		readTargets = append(readTargets, filepath.ToSlash(filepath.Join("personas", "_shared", filepath.Base(options.POVFile))))
	}
	if options.QuestionLedger != "" {
		readTargets = append(readTargets, options.QuestionLedger)
	}
	if options.PersonaDirectory != "" && strictMode(options.Mode) {
		readTargets = append(readTargets, filepath.ToSlash(filepath.Join("personas", options.PersonaDirectory, "UNMAPPED-DISCOVERY-BUFFER.md")))
	}
	readTargets = append(readTargets, options.AdditionalReadTargets...)
	readTargets, err = normalizedUnique(readTargets)
	if err != nil {
		return nil, err
	}
	bindings := make([]FileBinding, 0, len(readTargets))
	for _, rel := range readTargets {
		data, err := confinedRead(options.Workspace, rel)
		if err != nil {
			return nil, fmt.Errorf("mandatory read target %s: %w", rel, err)
		}
		sum := sha256.Sum256(data)
		bindings = append(bindings, FileBinding{Path: rel, Content: string(data), Fingerprint: fmt.Sprintf("sha256:%x", sum)})
	}
	checkpoint := ""
	for _, binding := range bindings {
		if binding.Path == "0H-CHECKPOINT-SUMMARY.md" {
			checkpoint = binding.Fingerprint
			break
		}
	}
	return &Packet{GeneratedAt: time.Now().UTC(), Projection: projection, Protocol: model, Options: options, Files: bindings, ProjectionText: string(projectionData), CheckpointFingerprint: checkpoint}, nil
}

func commonReadTargets() []string {
	return []string{"0A-PREFLIGHT.md", "0D-GLOSSARY.md", "0E-INDEX.md", "0F-GLOBAL-STATE.md", "0G-DECONSTRUCTION-STATE.md", "0H-CHECKPOINT-SUMMARY.md", "10-SOURCE-INVENTORY.md", "11-TRAVERSAL-FRONTIER.md", "12-CLAIM-EVIDENCE.md"}
}

func staticModeReadTargets(mode string) []string {
	switch mode {
	case "Export Acquisition":
		return []string{"0C-SECURITY-PRIVACY.md", "13-EXPORT-RECONCILIATION.md", "17-ACQUISITION-CANDIDATES.md"}
	case "Discovery", "Ticket Resolution", "Promotion Review":
		return []string{"0B-AUTH-MODEL.md", "0C-SECURITY-PRIVACY.md", "14-CONTRADICTIONS.md", "16-SOURCE-COVERAGE.md"}
	case "Sequential Reconciliation", "Cross-Reference-Reconciliation":
		return []string{"13-EXPORT-RECONCILIATION.md", "14-CONTRADICTIONS.md", "15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "20-TRACEABILITY.md"}
	case "Profile Synchronization":
		return []string{"15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "18-CONFIRMATIONS.md", "20-TRACEABILITY.md"}
	case "Partial-Synthesis", "Final-Synthesis":
		return []string{"13-EXPORT-RECONCILIATION.md", "14-CONTRADICTIONS.md", "15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "18-CONFIRMATIONS.md", "20-TRACEABILITY.md", "21-COVERAGE-REPORT.md", "22-GATE-REPORTS.md"}
	case "Reconstruction-Handoff":
		return []string{"15-DECISIONS.md", "18-CONFIRMATIONS.md", "20-TRACEABILITY.md", "22-GATE-REPORTS.md", "90-ARCH-BLUEPRINT.md", "91-DATA-MODEL.md", "92-BUSINESS-RULES.md", "93-USE-CASES.md", "94-INTERFACES.md", "95-DEPLOYMENT.md", "96-NON-FUNCTIONAL-SECURITY.md"}
	case "Human Hatch":
		return []string{"14-CONTRADICTIONS.md", "15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "18-CONFIRMATIONS.md", "20-TRACEABILITY.md", "22-GATE-REPORTS.md"}
	case "Validation/Gate":
		return []string{"13-EXPORT-RECONCILIATION.md", "14-CONTRADICTIONS.md", "15-DECISIONS.md", "16-SOURCE-COVERAGE.md", "17-ACQUISITION-CANDIDATES.md", "18-CONFIRMATIONS.md", "20-TRACEABILITY.md", "21-COVERAGE-REPORT.md", "22-GATE-REPORTS.md"}
	default:
		return nil
	}
}

func strictMode(mode string) bool {
	return mode == "Discovery" || mode == "Ticket Resolution" || mode == "Promotion Review"
}

func normalizedUnique(values []string) ([]string, error) {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		normalized, err := identity.NormalizeRelativePath(value)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	sort.Strings(result)
	return result, nil
}

func confinedRead(root, relative string) ([]byte, error) {
	normalized, err := identity.NormalizeRelativePath(relative)
	if err != nil {
		return nil, err
	}
	rootPath, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	rootPath, err = filepath.EvalSymlinks(rootPath)
	if err != nil {
		return nil, err
	}
	candidate, err := filepath.EvalSymlinks(filepath.Join(rootPath, filepath.FromSlash(normalized)))
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(rootPath, candidate)
	if err != nil {
		return nil, err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return nil, fmt.Errorf("target escapes workspace root")
	}
	return os.ReadFile(candidate)
}

func (p *Packet) Markdown() (string, error) {
	var out strings.Builder
	fmt.Fprintln(&out, "# Legacy Autopsy Invocation Context")
	fmt.Fprintln(&out, "\n> NON-AUTHORITATIVE-DERIVATIVE")
	fmt.Fprintln(&out, "> Schema: legacy-autopsy-context/v1")
	fmt.Fprintf(&out, "> Generated: %s\n", p.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&out, "> Protocol source fingerprint: %s\n", p.Protocol.Fingerprint)
	fmt.Fprintln(&out, "\n## Resume identity header")
	fields := [][2]string{
		{"Protocol Version", "v" + p.Protocol.Version}, {"System Namespace", p.Options.SystemNamespace}, {"Current Iteration", p.Options.Iteration},
		{"Invocation ID", p.Options.InvocationID}, {"Invocation Mode", p.Options.Mode}, {"Human Hatch Action", valueOrNone(p.Options.HumanHatchAction)},
		{"Invocation Scope", p.Options.InvocationScope}, {"Persona", valueOrNone(p.Options.Persona)}, {"Persona Prefix", valueOrNone(p.Options.PersonaPrefix)},
		{"Persona Directory", valueOrNone(p.Options.PersonaDirectory)}, {"Entry Cluster", valueOrNone(p.Options.Cluster)}, {"Traversal Track", valueOrNone(p.Options.Track)}, {"Active POV", valueOrNone(p.Options.POV)},
		{"Active POV File", valueOrNone(p.Options.POVFile)}, {"Active Question Ledger", valueOrNone(p.Options.QuestionLedger)},
		{"Environment/Snapshot", p.Options.EnvironmentSnapshot}, {"Max Traversal Depth", fmt.Sprintf("%d", p.Options.MaxTraversalDepth)},
		{"Allowed Read Targets", joinOrNone(filePaths(p.Files))}, {"Allowed Write Targets", joinOrNone(p.Options.AllowedWriteTargets)},
		{"Forbidden Write Targets", joinOrNone(p.Options.ForbiddenWriteTargets)}, {"Loaded Checkpoint Fingerprint", p.CheckpointFingerprint},
	}
	for _, field := range fields {
		fmt.Fprintf(&out, "- **%s:** %s\n", field[0], field[1])
	}
	fmt.Fprintln(&out, "\n- **Cold Resume Status:** inputs assembled; stale comparison and mutation authorization are not performed by this bootstrap command")
	fmt.Fprintln(&out, "\n## Workspace read set")
	for _, file := range p.Files {
		fmt.Fprintf(&out, "\n### `%s`\n\n- **Fingerprint:** `%s`\n\n<!-- BEGIN WORKSPACE-INPUT: %s -->\n%s", file.Path, file.Fingerprint, file.Path, file.Content)
		if !strings.HasSuffix(file.Content, "\n") {
			fmt.Fprintln(&out)
		}
		fmt.Fprintf(&out, "<!-- END WORKSPACE-INPUT: %s -->\n", file.Path)
	}
	fmt.Fprintln(&out, "\n## Operational projection")
	fmt.Fprintln(&out, p.ProjectionText)
	fmt.Fprintln(&out, "\n## Exact normative sections")
	for _, reference := range p.Projection.ProtocolSections {
		text, err := p.Protocol.SectionText(reference)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&out, "\n<!-- routed from %s -->\n%s", reference, text)
	}
	return out.String(), nil
}

func filePaths(files []FileBinding) []string {
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Path
	}
	return paths
}

func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "None"
	}
	return strings.Join(values, ", ")
}

func valueOrNone(value string) string {
	if value == "" {
		return "None"
	}
	return value
}
