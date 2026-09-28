// Package editing validates non-authoritative editor tooling against the CDL
// language surface.
package editing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const grammarPath = "../../editors/vscode/syntaxes/cdl.tmLanguage.json"

// keywordSurface lists every keyword the CDL parser dispatches on. The grammar
// must mention each one so highlighting cannot silently drift from the parser.
var keywordSurface = []string{
	"GLOBAL DECLARATIONS", "SECTION", "PART", "SUBSECTION OF", "PROJECTION",
	"BASE-READS", "MODE", "TITLE", "ARTIFACT", "GOAL", "USES",
	"REQUIRES CAPABILITY", "TYPE", "STATE", "VALUE", "ENUM", "FIELD", "TABLE",
	"RULE", "STEP", "END", "CAPABILITY", "REGISTRY", "GATES", "WORKFLOW-TARGET",
	"STRICT", "READS", "VALUES", "INITIAL", "ALLOWS", "LIFETIME", "ROWS",
	"ROW-TYPE", "KEY", "PREDICATE", "MACHINE:", "AGENT:", "EVIDENCE", "PRODUCES",
	"RESULT", "REQUIRE", "FOR EACH", "APPEND", "LET", "SET", "GUARD",
	"OTHERWISE BLOCK", "IF", "ELSE", "BLOCK", "CHANNELS", "CHANNEL",
	"COVERS SECTION", "TEXT", "OMIT STEP", "REASON", "SUPPLIED-BY",
	"AND", "OR", "NOT", "COUNT", "WHERE", "true", "false",
}

func TestGrammarCoversKeywordSurface(t *testing.T) {
	raw, err := os.ReadFile(filepath.Clean(grammarPath))
	if err != nil {
		t.Fatal(err)
	}
	var grammar struct {
		ScopeName string `json:"scopeName"`
		Patterns  []any  `json:"patterns"`
	}
	if err := json.Unmarshal(raw, &grammar); err != nil {
		t.Fatalf("grammar is not valid JSON: %v", err)
	}
	if grammar.ScopeName != "source.cdl" {
		t.Fatalf("scopeName = %q, want source.cdl", grammar.ScopeName)
	}
	if len(grammar.Patterns) == 0 {
		t.Fatal("grammar has no patterns")
	}
	text := string(raw)
	for _, keyword := range keywordSurface {
		if !strings.Contains(text, keyword) {
			t.Errorf("grammar does not mention keyword %q", keyword)
		}
	}
	if !strings.Contains(text, "text.html.markdown") {
		t.Error("grammar does not embed Markdown for opaque blocks")
	}
}

func TestLanguageContribution(t *testing.T) {
	raw, err := os.ReadFile("../../editors/vscode/package.json")
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct {
		Contributes struct {
			Languages []struct {
				ID         string   `json:"id"`
				Extensions []string `json:"extensions"`
			} `json:"languages"`
			Grammars []struct {
				ScopeName string `json:"scopeName"`
				Path      string `json:"path"`
			} `json:"grammars"`
		} `json:"contributes"`
	}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatalf("package.json is not valid JSON: %v", err)
	}
	if len(pkg.Contributes.Languages) != 1 || pkg.Contributes.Languages[0].ID != "cdl" {
		t.Fatalf("unexpected language contribution: %+v", pkg.Contributes.Languages)
	}
	if ext := pkg.Contributes.Languages[0].Extensions; len(ext) != 1 || ext[0] != ".cdl" {
		t.Fatalf("unexpected extensions: %v", ext)
	}
	if len(pkg.Contributes.Grammars) != 1 || pkg.Contributes.Grammars[0].ScopeName != "source.cdl" {
		t.Fatalf("unexpected grammar contribution: %+v", pkg.Contributes.Grammars)
	}
}
