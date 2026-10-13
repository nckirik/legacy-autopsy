package handbook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexMarkdownIsTagFreeAndLinked(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "chapter-01.md"), []byte("## 01. System Overview\n\n### Audience Purpose\n\nThe operator reads a map of the system (CLM-000000000001).\n\n### Observed Legacy Behavior\n\nText.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "chapter-05.md"), []byte("## 05. Use Cases\n\n### Audience Purpose\n\nThe developer studies request flows.\n\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, err := LoadIndexEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := string(IndexMarkdown(entries))
	for _, marker := range []string{
		"- [01. System Overview](chapter-01.md) - The operator reads a map of the system.",
		"- [05. Use Cases](chapter-05.md) - The developer studies request flows.",
	} {
		if !strings.Contains(out, marker) {
			t.Fatalf("index missing %q:\n%s", marker, out)
		}
	}
}
