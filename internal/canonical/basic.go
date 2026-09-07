// Package canonical contains explicitly scoped hashing helpers. BasicMarkdown
// is a bootstrap primitive; it is not the complete Protocol §4.1.2 profile.
package canonical

import (
	"crypto/sha256"
	"fmt"
	"strings"

	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

func BasicMarkdown(source string) (string, error) {
	source = strings.ReplaceAll(strings.ReplaceAll(source, "\r\n", "\n"), "\r", "\n")
	if _, err := md.Parse(source); err != nil {
		return "", err
	}
	lines := strings.Split(source, "\n")
	out := make([]string, 0, len(lines))
	inFence := false
	blank := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
		}
		if !inFence {
			line = strings.TrimRight(line, " \t")
		}
		if line == "" && !inFence {
			if blank {
				continue
			}
			blank = true
		} else {
			blank = false
		}
		out = append(out, line)
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return strings.Join(out, "\n") + "\n", nil
}

func BasicMarkdownFingerprint(source string) (string, error) {
	canonical, err := BasicMarkdown(source)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(canonical))
	return fmt.Sprintf("sha256:%x", sum), nil
}
