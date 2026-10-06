// Package discovery runs the Discovery pipeline: deterministic bounded source
// enumeration feeding model-authored source-inventory records with
// runtime-computed SRC identities.
package discovery

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/preflight"
)

// Candidate is one deterministically enumerated source unit.
type Candidate struct {
	Path        string
	Size        int64
	Fingerprint string
}

// Enumerate walks the target repository and returns bounded, sorted candidates.
// scopePaths restrict the walk to those normalized relative prefixes; an empty
// scope means the whole repository. Fingerprints are streamed SHA-256 over the
// exact file bytes.
func Enumerate(root string, scopePaths []string, maxFiles int) ([]Candidate, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("discovery: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("discovery: %s is not a directory", root)
	}
	scopes := make([]string, 0, len(scopePaths))
	for _, scope := range scopePaths {
		normalized, err := capabilities.NormalizeRelativePath(scope)
		if err != nil {
			return nil, fmt.Errorf("discovery: scope %q: %w", scope, err)
		}
		scopes = append(scopes, normalized)
	}
	inScope := func(relative string) bool {
		if len(scopes) == 0 {
			return true
		}
		for _, scope := range scopes {
			if relative == scope || strings.HasPrefix(relative, scope+"/") {
				return true
			}
		}
		return false
	}
	candidates := make([]Candidate, 0, maxFiles)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != root && preflight.SkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.Mode().IsRegular() || len(candidates) >= maxFiles {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if !inScope(relative) {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			file.Close()
			return err
		}
		file.Close()
		candidates = append(candidates, Candidate{
			Path:        relative,
			Size:        info.Size(),
			Fingerprint: fmt.Sprintf("sha256:%x", hasher.Sum(nil)),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discovery: enumerate %s: %w", root, err)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
	return candidates, nil
}
