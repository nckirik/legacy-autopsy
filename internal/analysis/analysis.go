// Package analysis computes deterministic coverage metrics over the authored
// protocol classification. It validates structural coverage against protocol.md;
// it never validates the classification judgement itself.
package analysis

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/markdown"
)

// Entry is one authored section classification.
type Entry struct {
	Heading string   `json:"heading"`
	Driver  string   `json:"driver"`
	Tags    []string `json:"tags"`
}

// Classification is the authored classification document.
type Classification struct {
	Protocol          string  `json:"protocol"`
	SourceFingerprint string  `json:"source-fingerprint"`
	Rubric            string  `json:"rubric"`
	Entries           []Entry `json:"entries"`
}

// Metrics is the deterministic coverage profile.
type Metrics struct {
	Generator               string             `json:"generator"`
	SourceFingerprint       string             `json:"source-fingerprint"`
	Headings                int                `json:"headings"`
	Sections                int                `json:"sections"`
	Parts                   int                `json:"parts"`
	Drivers                 map[string]int     `json:"drivers"`
	Tags                    map[string]int     `json:"tags"`
	SectionLines            map[string]int     `json:"section-lines"`
	SectionLineShare        map[string]float64 `json:"section-line-share"`
	StructuredAgentSections int                `json:"structured-agent-sections"`
	FreeProseSections       int                `json:"free-prose-sections"`
}

// ValidDrivers and ValidTags are the closed vocabularies.
var ValidDrivers = map[string]bool{"MACHINE": true, "AGENT": true, "MIXED": true, "CONTAINER": true}

var ValidTags = map[string]bool{
	"gate": true, "transition": true, "artifact": true, "capability": true,
	"reuse": true, "contract": true, "human": true, "policy": true,
}

// Heading is one classified heading location.
type Heading struct {
	Level int
	Text  string
	Start int
	End   int
}

// Load reads an authored classification document.
func Load(path string) (*Classification, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Classification
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("classification: %w", err)
	}
	return &c, nil
}

// Headings returns every heading at levels 1-3 in document order, excluding the
// document title, with non-overlapping line spans.
func Headings(doc *markdown.Document) []Heading {
	var all []Heading
	for _, node := range doc.Nodes {
		if node.Kind != markdown.Heading || node.Level > 3 {
			continue
		}
		all = append(all, Heading{Level: node.Level, Text: node.Text, Start: node.Span.StartLine, End: node.Span.EndLine})
	}
	if len(all) > 0 && strings.HasPrefix(all[0].Text, "Protocol: Legacy System Deconstruction") {
		all = all[1:]
	}
	totalLines := strings.Count(doc.Source, "\n") + 1
	for i := range all {
		end := totalLines
		for j := i + 1; j < len(all); j++ {
			if all[j].Level <= all[i].Level {
				end = all[j].Start - 1
				break
			}
		}
		if end < all[i].Start {
			end = all[i].Start
		}
		all[i].End = end
	}
	return all
}

// Validate checks structural coverage, vocabulary, and fingerprint.
func Validate(doc *markdown.Document, fingerprint string, c *Classification) error {
	if c.SourceFingerprint != fingerprint {
		return fmt.Errorf("classification fingerprint %s does not match protocol.md %s; revise the classification", c.SourceFingerprint, fingerprint)
	}
	headings := Headings(doc)
	if len(headings) != len(c.Entries) {
		return fmt.Errorf("classification covers %d headings but protocol.md has %d", len(c.Entries), len(headings))
	}
	for i, h := range headings {
		e := c.Entries[i]
		if e.Heading != h.Text {
			return fmt.Errorf("entry %d is %q but protocol.md heading is %q", i, e.Heading, h.Text)
		}
		if !ValidDrivers[e.Driver] {
			return fmt.Errorf("entry %q has invalid driver %q", e.Heading, e.Driver)
		}
		for _, tag := range e.Tags {
			if !ValidTags[tag] {
				return fmt.Errorf("entry %q has invalid tag %q", e.Heading, tag)
			}
		}
		if h.Level >= 2 && e.Driver == "CONTAINER" {
			return fmt.Errorf("section %q must have a non-container driver", e.Heading)
		}
	}
	return nil
}

// Compute derives the deterministic coverage metrics.
func Compute(doc *markdown.Document, c *Classification) (Metrics, error) {
	headings := Headings(doc)
	m := Metrics{
		Generator:         "internal/analysis",
		SourceFingerprint: c.SourceFingerprint,
		Headings:          len(c.Entries),
		Drivers:           map[string]int{},
		Tags:              map[string]int{},
		SectionLines:      map[string]int{},
		SectionLineShare:  map[string]float64{},
	}
	totalLines := 0
	for i, e := range c.Entries {
		m.Drivers[e.Driver]++
		for _, tag := range e.Tags {
			m.Tags[tag]++
		}
		if headings[i].Level == 1 {
			m.Parts++
			continue
		}
		m.Sections++
		lines := headings[i].End - headings[i].Start + 1
		m.SectionLines[e.Driver] += lines
		totalLines += lines
		if e.Driver != "MACHINE" && hasTag(e, "contract") {
			m.StructuredAgentSections++
		}
		if !hasTag(e, "contract") && (e.Driver == "AGENT" || e.Driver == "MIXED") {
			m.FreeProseSections++
		}
	}
	if totalLines == 0 {
		return m, nil
	}
	for driver, lines := range m.SectionLines {
		m.SectionLineShare[driver] = math.Round(float64(lines)/float64(totalLines)*10000) / 10000
	}
	return m, nil
}

// Marshal renders metrics canonically with sorted object keys.
func (m Metrics) Marshal() ([]byte, error) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func hasTag(e Entry, tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// SortedKeys is a helper for deterministic reporting.
func SortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
