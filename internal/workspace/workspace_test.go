package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitAndCheck(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".extracted")
	if err := Init(root); err != nil {
		t.Fatal(err)
	}
	if err := Check(root); err != nil {
		t.Fatal(err)
	}
	for _, rel := range RequiredFiles() {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		for _, forbidden := range []string{"[C-COVERED]", "[B-CONFIRMED]", "EXIT-E-STATUS: Passed"} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("%s fabricates %s", rel, forbidden)
			}
		}
	}
}

func TestCheckRejectsWrongHeading(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".extracted")
	if err := Init(root); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "10-SOURCE-INVENTORY.md"), []byte("# Wrong\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Check(root); err == nil {
		t.Fatal("expected heading validation error")
	}
}

func TestCheckRejectsMissingRequiredFile(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".extracted")
	if err := Init(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "10-SOURCE-INVENTORY.md")); err != nil {
		t.Fatal(err)
	}
	if err := Check(root); err == nil {
		t.Fatal("expected validation error")
	}
}
