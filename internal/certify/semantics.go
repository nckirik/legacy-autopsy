package certify

import (
	"fmt"
	"regexp"
)

var (
	prefixPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,7}$`)
	slugPattern   = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// ValidationHints returns the deterministic field-level rules as short
// human-readable lines, so executor prompts and validators share one source.
func ValidationHints(fieldID string) []string {
	switch fieldID {
	case "PERSONA-REGISTRY":
		return []string{
			`prefix must match [A-Z][A-Z0-9]{1,7} (for example BE, WEB, OPS)`,
			`persona-slug must be lowercase words joined by hyphens`,
			`persona-directory must equal prefix + "-" + persona-slug`,
		}
	}
	return nil
}

// ValidateSemantics applies protocol field-level semantics that go beyond the
// declared types. Unknown FIELD IDs pass through; only the deterministic
// identity contracts implemented so far are enforced.
func ValidateSemantics(fieldID string, values map[string]any) error {
	switch fieldID {
	case "PERSONA-REGISTRY":
		prefix := stringValue(values["prefix"])
		slug := stringValue(values["persona-slug"])
		directory := stringValue(values["persona-directory"])
		if !prefixPattern.MatchString(prefix) {
			return fmt.Errorf("%s.prefix %q must match [A-Z][A-Z0-9]{1,7}", fieldID, prefix)
		}
		if !slugPattern.MatchString(slug) {
			return fmt.Errorf("%s.persona-slug %q must be lowercase alphanumeric words joined by hyphens", fieldID, slug)
		}
		if directory != prefix+"-"+slug {
			return fmt.Errorf("%s.persona-directory %q must equal %s-%s", fieldID, directory, prefix, slug)
		}
	}
	return nil
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}
