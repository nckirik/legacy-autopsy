package cdl

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFormatRoundTripProtocol(t *testing.T) {
	files, err := filepath.Glob("../protocol/*.cdl")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no protocol sources found")
	}
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before, diags := Parse(path, src)
		if len(diags) > 0 {
			t.Fatalf("%s: %v", path, diags)
		}
		formatted, diags := FormatSource(path, src)
		if len(diags) > 0 {
			t.Fatalf("%s: %v", path, diags)
		}
		after, diags := Parse(path, formatted)
		if len(diags) > 0 {
			t.Fatalf("%s formatted: %v", path, diags)
		}
		if !reflect.DeepEqual(before, after) {
			t.Fatalf("%s: formatting changed the AST", path)
		}
		if second := Format(after); string(second) != string(formatted) {
			t.Fatalf("%s: formatting is not idempotent", path)
		}
	}
}

func TestFormatRejectsUnpreservedComments(t *testing.T) {
	src := []byte("SECTION s\nTITLE \"t\"\nGOAL\n  body\nEND\n# trailing nested is fine at top\n")
	if _, diags := FormatSource("t.cdl", src); len(diags) > 0 {
		t.Fatalf("top-level trailing comment rejected: %v", diags)
	}
	inline := []byte("SECTION s\nTITLE \"t\" # note\nGOAL\n  body\nEND\n")
	if _, diags := FormatSource("t.cdl", inline); len(diags) == 0 {
		t.Fatal("inline comment accepted")
	}
	nested := []byte("SECTION s\nTITLE \"t\"\n# comment inside section\nGOAL\n  body\nEND\n")
	if _, diags := FormatSource("t.cdl", nested); len(diags) == 0 {
		t.Fatal("nested comment accepted")
	}
}

func TestFormatPreservesOpaqueText(t *testing.T) {
	src := []byte("SECTION s\nTITLE \"t\"\nGOAL\n  # not a comment\n  | a | b |\n\n  ```\n  code # still code\n  ```\nEND\n")
	prog, diags := Parse("t.cdl", src)
	if len(diags) > 0 {
		t.Fatal(diags)
	}
	out, diags := FormatSource("t.cdl", src)
	if len(diags) > 0 {
		t.Fatal(diags)
	}
	after, diags := Parse("t.cdl", out)
	if len(diags) > 0 {
		t.Fatal(diags)
	}
	if after.Sections[0].Goal != prog.Sections[0].Goal {
		t.Fatalf("opaque text changed:\n%q\n%q", prog.Sections[0].Goal, after.Sections[0].Goal)
	}
	if !strings.Contains(string(out), "# not a comment") || !strings.Contains(string(out), "code # still code") {
		t.Fatalf("opaque content missing:\n%s", out)
	}
}

func TestFormatExpressions(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want string
	}{
		{"a == b", "a == b"},
		{"a - (b - c)", "a - (b - c)"},
		{"COUNT(t WHERE f == 1 + 2)", "COUNT(t WHERE f == 1 + 2)"},
		{"NOT (a OR b)", "NOT (a OR b)"},
		{"F(a, b)", "F(a, b)"},
	} {
		expr, err := parseExprSource(tc.src)
		if err != nil {
			t.Fatalf("%s: %v", tc.src, err)
		}
		if got := exprString(expr); got != tc.want {
			t.Fatalf("%s: got %q want %q", tc.src, got, tc.want)
		}
		reparsed, err := parseExprSource(exprString(expr))
		if err != nil {
			t.Fatalf("%s reparse: %v", tc.src, err)
		}
		if !reflect.DeepEqual(expr, reparsed) {
			t.Fatalf("%s: expression round-trip changed the AST", tc.src)
		}
	}
}
