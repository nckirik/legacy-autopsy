// Package routing validates the non-authoritative mode index against the
// structurally parsed normative protocol.
package routing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	md "github.com/nckirik/legacy-autopsy/internal/markdown"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

const ProjectionMarker = "NON-AUTHORITATIVE OPERATIONAL PROJECTION"

type Manifest struct {
	Authority               string       `json:"authority"`
	Description             string       `json:"description"`
	ExpectedProtocolVersion string       `json:"expected_protocol_version"`
	Projections             []Projection `json:"projections"`
}

type Projection struct {
	Mode             string   `json:"mode"`
	SkillProjection  string   `json:"skill_projection"`
	ProtocolSections []string `json:"protocol_sections"`
}

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mode manifest: %w", err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse mode manifest as JSON: %w", err)
	}
	return &manifest, nil
}

func Validate(repoRoot string, manifest *Manifest, model *protocol.Model) error {
	var problems []string
	if !strings.Contains(manifest.Authority, ProjectionMarker) {
		problems = append(problems, "manifest lacks non-authoritative marker")
	}
	if manifest.ExpectedProtocolVersion != model.Version {
		problems = append(problems, fmt.Sprintf("protocol version mismatch: manifest=%q protocol=%q", manifest.ExpectedProtocolVersion, model.Version))
	}

	protocolModes := make(map[string]struct{}, len(model.Modes))
	for _, mode := range model.Modes {
		protocolModes[mode] = struct{}{}
	}
	seenModes := map[string]struct{}{}
	seenPaths := map[string]struct{}{}
	for _, projection := range manifest.Projections {
		if _, ok := protocolModes[projection.Mode]; !ok {
			problems = append(problems, fmt.Sprintf("unknown manifest mode %q", projection.Mode))
		}
		if _, duplicate := seenModes[projection.Mode]; duplicate {
			problems = append(problems, fmt.Sprintf("duplicate projection for mode %q", projection.Mode))
		}
		seenModes[projection.Mode] = struct{}{}
		if _, duplicate := seenPaths[projection.SkillProjection]; duplicate {
			problems = append(problems, fmt.Sprintf("projection path reused: %s", projection.SkillProjection))
		}
		seenPaths[projection.SkillProjection] = struct{}{}
		for _, section := range projection.ProtocolSections {
			if !model.HasSection(section) {
				problems = append(problems, fmt.Sprintf("%s references missing protocol heading %q", projection.Mode, section))
			}
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(projection.SkillProjection)))
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s projection: %v", projection.Mode, err))
		} else {
			text := string(data)
			if !strings.Contains(text, ProjectionMarker) {
				problems = append(problems, fmt.Sprintf("%s projection lacks non-authoritative marker", projection.Mode))
			}
			doc, parseErr := md.Parse(text)
			if parseErr != nil {
				problems = append(problems, fmt.Sprintf("%s projection Markdown: %v", projection.Mode, parseErr))
			} else {
				if doc.HeadingCount(projection.Mode) != 1 {
					problems = append(problems, fmt.Sprintf("%s projection must contain exactly one H1 named %q", projection.Mode, projection.Mode))
				}
				for _, required := range []string{"Purpose", "Required identity bindings", "Normative protocol sections", "Mandatory read set", "Semantic write targets", "Allowed assurance/audit side effects", "Forbidden mutations", "Deterministic validation", "Completion condition"} {
					if doc.HeadingCount(required) != 1 {
						problems = append(problems, fmt.Sprintf("%s projection requires exactly one %q section", projection.Mode, required))
					}
				}
			}
		}
	}
	for _, mode := range model.Modes {
		if _, ok := seenModes[mode]; !ok {
			problems = append(problems, fmt.Sprintf("protocol mode %q has no projection", mode))
		}
	}

	modeDir := filepath.Join(repoRoot, "skill", "modes")
	entries, err := os.ReadDir(modeDir)
	if err != nil {
		problems = append(problems, fmt.Sprintf("read projection directory: %v", err))
	} else {
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}
			rel := filepath.ToSlash(filepath.Join("skill", "modes", entry.Name()))
			if _, ok := seenPaths[rel]; !ok {
				problems = append(problems, fmt.Sprintf("unknown projection file %s", rel))
			}
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("routing validation failed:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
}

func (m *Manifest) Projection(mode string) (Projection, bool) {
	for _, projection := range m.Projections {
		if projection.Mode == mode {
			return projection, true
		}
	}
	return Projection{}, false
}
