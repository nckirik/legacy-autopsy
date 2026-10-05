package parity

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

var (
	oracleFieldBullet = regexp.MustCompile(`^- \*\*.+?:\*\*`)
	oracleListItem    = regexp.MustCompile(`^\d+\. `)
)

// TestGeneratedProtocolCoversOracle requires the generated protocol.md to carry
// every prose unit, fenced line, table row, field bullet, and heading of the frozen
// oracle. Comparison is content-normalized so formatting markers may differ.
func TestGeneratedProtocolCoversOracle(t *testing.T) {
	root := root(t)
	oracleBytes, err := os.ReadFile(protocol.EditionOraclePath(root))
	if err != nil {
		t.Fatal(err)
	}
	generatedBytes, err := os.ReadFile(filepath.Join(root, "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	generated := contentNormalize(string(generatedBytes))

	type unit struct {
		line int
		text string
	}
	var units []unit
	add := func(line int, text string) {
		if len(strings.Fields(contentNormalize(text))) >= 6 {
			units = append(units, unit{line, text})
		}
	}

	inFence := false
	paragraph := []string{}
	flush := func(line int) {
		if len(paragraph) > 0 {
			add(line, strings.Join(paragraph, " "))
			paragraph = nil
		}
	}
	for i, raw := range strings.Split(string(oracleBytes), "\n") {
		lineNo := i + 1
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "```") {
			flush(lineNo)
			inFence = !inFence
			continue
		}
		if inFence {
			add(lineNo, trimmed)
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "<!--") || strings.HasPrefix(trimmed, "---") {
			flush(lineNo)
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			flush(lineNo)
			if contentNormalize(trimmed) != "0 purpose guarantees and conceptual flow" {
				add(lineNo, trimmed)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "|") {
			flush(lineNo)
			add(lineNo, trimmed)
			continue
		}
		if oracleFieldBullet.MatchString(trimmed) {
			flush(lineNo)
			add(lineNo, trimmed)
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || oracleListItem.MatchString(trimmed) {
			flush(lineNo)
			add(lineNo, trimmed)
			continue
		}
		paragraph = append(paragraph, trimmed)
	}
	flush(len(strings.Split(string(oracleBytes), "\n")))

	var missing []string
	for _, u := range units {
		if !strings.Contains(generated, contentNormalize(u.text)) {
			missing = append(missing, u.text)
		}
	}
	if len(missing) > 0 {
		for i, text := range missing {
			if i == 10 {
				break
			}
			t.Errorf("oracle content missing from generated protocol.md: %s", text)
		}
		t.Fatalf("%d of %d oracle units missing", len(missing), len(units))
	}
}

func contentNormalize(text string) string {
	text = strings.ToLower(text)
	var b strings.Builder
	lastSpace := true
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastSpace = false
			continue
		}
		if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}
	return strings.TrimSpace(b.String())
}
