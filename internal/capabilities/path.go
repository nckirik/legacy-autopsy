package capabilities

import (
	"fmt"
	"path"
	"strings"
)

// NormalizeRelativePath normalizes a repository-relative path: POSIX
// separators, no empty or "." segments, no traversal, exact Unicode and case
// preservation, and no transliteration.
func NormalizeRelativePath(value string) (string, error) {
	value = strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	if value == "" {
		return "", fmt.Errorf("relative path is empty")
	}
	if strings.HasPrefix(value, "/") || (len(value) >= 2 && value[1] == ':') {
		return "", fmt.Errorf("path must be relative: %q", value)
	}
	parts := strings.Split(value, "/")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			return "", fmt.Errorf("path traversal is forbidden: %q", value)
		default:
			cleaned = append(cleaned, part)
		}
	}
	if len(cleaned) == 0 {
		return "", fmt.Errorf("relative path resolves to root")
	}
	return path.Join(cleaned...), nil
}
