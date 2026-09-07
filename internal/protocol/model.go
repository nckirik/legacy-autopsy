// Package protocol loads implementation-facing facts from the normative
// protocol Markdown without making this package a second specification.
package protocol

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

const InvocationModeHeading = "8.3. Invocation-mode enum and explicit multi-file modes"

type Model struct {
	Path        string
	Version     string
	Fingerprint string
	Document    *md.Document
	Modes       []string
}

func Load(path string) (*Model, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read protocol: %w", err)
	}
	doc, err := md.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse protocol Markdown: %w", err)
	}
	version := ""
	for _, node := range doc.Nodes {
		if node.Kind == md.Field && node.Label == "Version" {
			version = strings.Fields(node.Value)[0]
			break
		}
	}
	if version == "" {
		return nil, fmt.Errorf("protocol Version field not found structurally")
	}
	section, err := doc.Section(InvocationModeHeading)
	if err != nil {
		return nil, err
	}
	var modes []string
	for _, node := range section {
		if node.Kind != md.CodeBlock || !strings.Contains(node.Raw, "Preflight") {
			continue
		}
		for _, value := range strings.Split(strings.ReplaceAll(node.Raw, "\n", " "), "|") {
			value = strings.TrimSpace(value)
			if value != "" {
				modes = append(modes, value)
			}
		}
		break
	}
	if len(modes) == 0 {
		return nil, fmt.Errorf("invocation-mode registry not found in section %q", InvocationModeHeading)
	}
	sum := sha256.Sum256(data)
	return &Model{Path: path, Version: version, Fingerprint: fmt.Sprintf("sha256:%x", sum), Document: doc, Modes: modes}, nil
}

func HeadingFromReference(reference string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(reference), "§"))
}

func (m *Model) HasSection(reference string) bool {
	return m.Document.HeadingCount(HeadingFromReference(reference)) == 1
}

func (m *Model) SectionText(reference string) (string, error) {
	return m.Document.SectionText(HeadingFromReference(reference))
}
