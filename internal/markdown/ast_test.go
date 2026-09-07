package markdown

import "testing"

func TestParseAndSection(t *testing.T) {
	doc, err := Parse("**Version:** 4.0\n\n## 1. One\n\n- **State:** Ready\n\n```text\nA | B\n```\n\n## 2. Two\n")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := doc.Heading("1. One"); !ok {
		t.Fatal("heading not found")
	}
	section, err := doc.Section("1. One")
	if err != nil {
		t.Fatal(err)
	}
	if len(section) == 0 || section[0].Kind != Heading {
		t.Fatalf("unexpected section: %#v", section)
	}
	text, err := doc.SectionText("1. One")
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("empty section text")
	}
}

func TestParseRejectsUnclosedFence(t *testing.T) {
	if _, err := Parse("# Title\n\n```text\nmissing close\n"); err == nil {
		t.Fatal("expected error")
	}
}
