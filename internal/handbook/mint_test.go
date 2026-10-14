package handbook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMintSectionsRoundTrip(t *testing.T) {
	dir := t.TempDir()
	chapter := "## 01. Overview\n\nA developer reads this chapter first. The system starts on port 3000.\n\n<details><summary>MOD-ABCDEF012345</summary>\n\n#### Module\n\n- **Field:** value\n\n</details>\n"
	if err := os.WriteFile(filepath.Join(dir, "chapter-01.md"), []byte(chapter), 0o644); err != nil {
		t.Fatal(err)
	}
	resolve := func(id string) (UpstreamRef, bool) {
		if id != "MOD-ABCDEF012345" {
			return UpstreamRef{}, false
		}
		return UpstreamRef{ID: id, Version: 2, Fingerprint: "sha256:abc"}, true
	}
	sections, err := MintSections(dir, "demo", resolve)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("sections = %d, want 1", len(sections))
	}
	section := sections[0]
	if !strings.HasPrefix(section.ID, "HBK-") || len(section.ID) != 16 {
		t.Fatalf("HBK id = %q", section.ID)
	}
	if section.Anchor != "chapter-01" || section.Path != "chapter-01.md" {
		t.Fatalf("section coordinates = %+v", section)
	}
	if len(section.Upstream) != 1 || section.Upstream[0].Version != 2 {
		t.Fatalf("upstream = %+v", section.Upstream)
	}
	if err := WriteSections(dir, sections); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadMintedSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || loaded[0].ID != section.ID || loaded[0].Fingerprint != section.Fingerprint {
		t.Fatalf("loaded = %+v", loaded)
	}
	again, err := MintSections(dir, "demo", resolve)
	if err != nil {
		t.Fatal(err)
	}
	if string(again[0].Bytes) != string(section.Bytes) {
		t.Fatalf("re-mint changed bytes")
	}
	restored, err := RestoreDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if restored != 1 {
		t.Fatalf("restored = %d, want 1", restored)
	}
	clean, err := os.ReadFile(filepath.Join(dir, "chapter-01.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(clean) != chapter {
		t.Fatalf("restore did not reproduce the clean chapter:\n%s", clean)
	}
	if _, ok := sections[0].Upstream[0], true; !ok {
		t.Fatal("unreachable")
	}
}

func TestMintSectionsFailsOnUnresolvedUpstream(t *testing.T) {
	dir := t.TempDir()
	chapter := "## 01. Overview\n\n<details><summary>SEC-FFFFFFFFFFFF</summary>\n\n</details>\n"
	if err := os.WriteFile(filepath.Join(dir, "chapter-01.md"), []byte(chapter), 0o644); err != nil {
		t.Fatal(err)
	}
	resolve := func(id string) (UpstreamRef, bool) { return UpstreamRef{}, false }
	if _, err := MintSections(dir, "demo", resolve); err == nil {
		t.Fatalf("unresolved upstream accepted")
	}
}
