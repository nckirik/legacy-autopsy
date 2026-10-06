// Package preflight builds and runs the Preflight pipeline: a deterministic
// bounded repository scan feeds model-authored `0A` header and persona registry
// records through the bounded work session.
package preflight

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// Scan is a bounded, deterministic repository scan.
type Scan struct {
	Namespace   string
	Files       int
	Languages   []string
	EntryPoints []string
	TopDirs     []string
	Evidence    []Evidence
	Fingerprint string
}

// Evidence is a bounded excerpt of a key repository file.
type Evidence struct {
	Path    string
	Excerpt string
}

// SkipDirectory reports whether a directory name is excluded from scans.
func SkipDirectory(name string) bool {
	return skipDirs[name] || strings.HasPrefix(name, ".")
}

var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true, ".extracted": true,
	".tmp": true, ".idea": true, ".vscode": true, "dist": true, "build": true,
	"target": true, "__pycache__": true,
}

var languageByExt = map[string]string{
	".go": "Go", ".js": "JavaScript", ".mjs": "JavaScript", ".cjs": "JavaScript",
	".ts": "TypeScript", ".tsx": "TypeScript", ".jsx": "JavaScript",
	".py": "Python", ".rb": "Ruby", ".php": "PHP", ".java": "Java", ".kt": "Kotlin",
	".cs": "C#", ".c": "C", ".h": "C", ".cc": "C++", ".cpp": "C++", ".rs": "Rust",
	".sh": "Shell", ".ps1": "PowerShell", ".sql": "SQL", ".yml": "YAML", ".yaml": "YAML",
	".json": "JSON", ".xml": "XML", ".html": "HTML", ".css": "CSS", ".md": "Markdown",
	".tf": "Terraform", ".gradle": "Gradle", ".lua": "Lua", ".pl": "Perl", ".r": "R",
}

// evidenceNames are read (bounded) for preflight evidence.
var evidenceNames = map[string]bool{
	"README.md": true, "readme.md": true, "package.json": true, "go.mod": true,
	"docker-compose.yml": true, "docker-compose.yaml": true, "Dockerfile": true,
	"index.js": true, "server.js": true, "app.js": true, "main.go": true,
	"main.py": true, "app.py": true,
}

var evidenceDirNames = map[string]bool{"routes": true, "passport": true, "api": true}

var entryNames = map[string]bool{
	"main.go": true, "index.js": true, "server.js": true, "app.js": true, "index.ts": true,
	"main.py": true, "manage.py": true, "app.py": true, "main.rb": true, "index.php": true,
	"Dockerfile": true, "docker-compose.yml": true, "docker-compose.yaml": true,
	"package.json": true, "go.mod": true, "Makefile": true, "Procfile": true,
}

// ScanRepo walks the repository with bounded depth and file count. It never
// follows symlinks; it reads only bounded excerpts of key evidence files.
func ScanRepo(root string, maxFiles, maxEntryPoints, maxDepth int) (Scan, error) {
	info, err := os.Stat(root)
	if err != nil {
		return Scan{}, fmt.Errorf("preflight: scan %s: %w", root, err)
	}
	if !info.IsDir() {
		return Scan{}, fmt.Errorf("preflight: %s is not a directory", root)
	}
	scan := Scan{Namespace: capabilities.ColumnName(filepath.Base(filepath.Clean(root)))}
	if scan.Namespace == "" {
		scan.Namespace = "target"
	}
	languages := map[string]int{}
	var entries []string
	var topDirs []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		depth := len(strings.Split(filepath.ToSlash(relative), "/"))
		if info.IsDir() {
			base := info.Name()
			if SkipDirectory(base) {
				return filepath.SkipDir
			}
			if depth == 1 {
				topDirs = append(topDirs, base)
			}
			if depth >= maxDepth {
				return filepath.SkipDir
			}
			return nil
		}
		if scan.Files >= maxFiles {
			return nil
		}
		scan.Files++
		if language, ok := languageByExt[strings.ToLower(filepath.Ext(info.Name()))]; ok {
			languages[language]++
		}
		if entryNames[info.Name()] && depth <= 2 {
			entries = append(entries, filepath.ToSlash(relative))
		}
		parent := filepath.Base(filepath.Dir(relative))
		if (evidenceNames[info.Name()] && depth <= 2) || (evidenceDirNames[parent] && depth <= 3) {
			if len(scan.Evidence) < 12 {
				if data, err := os.ReadFile(path); err == nil {
					excerpt := string(data)
					if len(excerpt) > 1500 {
						excerpt = excerpt[:1500] + "\n...[truncated]"
					}
					scan.Evidence = append(scan.Evidence, Evidence{Path: filepath.ToSlash(relative), Excerpt: excerpt})
				}
			}
		}
		return nil
	})
	if err != nil {
		return Scan{}, fmt.Errorf("preflight: scan %s: %w", root, err)
	}
	for language := range languages {
		scan.Languages = append(scan.Languages, language)
	}
	sort.Strings(scan.Languages)
	sort.Strings(topDirs)
	scan.TopDirs = topDirs
	sort.Strings(entries)
	if len(entries) > maxEntryPoints {
		entries = entries[:maxEntryPoints]
	}
	scan.EntryPoints = entries
	scan.Fingerprint = capabilities.Hash(scan.Markdown())
	return scan, nil
}

// Markdown materializes the scan as a bounded context input.
func (s Scan) Markdown() []byte {
	var b strings.Builder
	b.WriteString("# Repository scan (runtime-generated, bounded)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", s.Namespace)
	fmt.Fprintf(&b, "- **Indexed File Count:** %d\n", s.Files)
	fmt.Fprintf(&b, "- **Languages:** %s\n", joinOrNone(s.Languages))
	fmt.Fprintf(&b, "- **Entry Points:** %s\n", joinOrNone(s.EntryPoints))
	fmt.Fprintf(&b, "- **Top-Level Directories:** %s\n", joinOrNone(s.TopDirs))
	if len(s.Evidence) > 0 {
		b.WriteString("\n## Bounded evidence excerpts\n")
		for _, evidence := range s.Evidence {
			fmt.Fprintf(&b, "\n### %s\n\n```text\n%s\n```\n", evidence.Path, evidence.Excerpt)
		}
	}
	return []byte(b.String())
}

func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "None"
	}
	return strings.Join(values, ", ")
}
