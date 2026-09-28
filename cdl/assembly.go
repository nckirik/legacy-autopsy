package cdl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Assembly is the ordered source manifest for one document compilation.
type Assembly struct {
	Sources []string `json:"sources"`
}

// LoadAssembly reads an ordered manifest of repository-relative source paths and
// returns the sources for one document compilation.
func LoadAssembly(repoRoot, manifestRel string) ([]Source, error) {
	manifestPath := filepath.Join(repoRoot, filepath.FromSlash(manifestRel))
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var manifest Assembly
	if err := json.Unmarshal(b, &manifest); err != nil {
		return nil, fmt.Errorf("assembly manifest: %w", err)
	}
	if len(manifest.Sources) == 0 {
		return nil, fmt.Errorf("assembly manifest %s lists no sources", manifestRel)
	}
	seen := map[string]bool{}
	out := make([]Source, 0, len(manifest.Sources))
	for _, rel := range manifest.Sources {
		if seen[rel] {
			return nil, fmt.Errorf("assembly manifest %s lists %s twice", manifestRel, rel)
		}
		seen[rel] = true
		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
		if err != nil {
			return nil, err
		}
		out = append(out, Source{Path: rel, Bytes: content})
	}
	return out, nil
}
