package handbook

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// IndexEntry is one chapter entry in the handbook index.
type IndexEntry struct {
	Number  string
	Title   string
	Path    string
	Summary string
}

// IndexMarkdown renders the top-level handbook index: a tag-free list that
// links each chapter with a one-line summary.
func IndexMarkdown(entries []IndexEntry) []byte {
	sorted := append([]IndexEntry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Number < sorted[j].Number })
	var b strings.Builder
	b.WriteString("# Handbook Index\n\n")
	b.WriteString("Read the chapters in order. Each chapter is a design document with footnoted evidence.\n\n")
	for _, entry := range sorted {
		summary := strings.TrimSpace(entry.Summary)
		if summary == "" {
			summary = "None"
		}
		fmt.Fprintf(&b, "- [%s. %s](%s) - %s\n", entry.Number, entry.Title, entry.Path, summary)
	}
	return []byte(b.String())
}

// LoadIndexEntries reads rendered chapter files (chapter-NN.md) from a directory
// and builds index entries from their headings and audience-purpose sentences.
func LoadIndexEntries(dir string) ([]IndexEntry, error) {
	files, err := filepath.Glob(filepath.Join(dir, "chapter-*.md"))
	if err != nil {
		return nil, err
	}
	entries := make([]IndexEntry, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		entry := IndexEntry{Path: filepath.Base(file)}
		text := string(data)
		for _, line := range strings.Split(text, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "## ") && entry.Number == "" {
				heading := strings.TrimPrefix(trimmed, "## ")
				if dot := strings.Index(heading, ". "); dot > 0 {
					entry.Number = heading[:dot]
					entry.Title = heading[dot+2:]
				}
			}
			if strings.HasPrefix(trimmed, "### Audience Purpose") {
				continue
			}
			if entry.Number != "" && entry.Summary == "" && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				entry.Summary = sanitizeSummary(firstSentence(trimmed))
			}
		}
		if entry.Number == "" {
			continue
		}
		if entry.Title == "" {
			entry.Title = "Chapter " + entry.Number
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// sanitizeSummary removes footnote marks and machine IDs so the index stays a
// clean reading surface.
func sanitizeSummary(text string) string {
	text = regexp.MustCompile(`\[\^\d+\]`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(CLM|MOD|SRC|ER|BR|UC|IF|REL|SM|DR|NFR|SEC|FLT|COV|SWP|FRT)-[A-Z0-9]{6,}`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\(\s*\)`).ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "  ", " ")
	return strings.TrimSpace(strings.ReplaceAll(text, " .", "."))
}

func firstSentence(text string) string {
	if idx := strings.IndexAny(text, ".!?"); idx > 0 {
		return strings.TrimSpace(text[:idx+1])
	}
	return strings.TrimSpace(text)
}
