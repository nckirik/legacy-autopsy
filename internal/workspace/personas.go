package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PersonaRegistration is one registered persona binding used to validate the
// workspace layout: the 0A registry prefix, explicit slug, and canonical
// directory basename `<persona-prefix>-<persona-slug>`.
type PersonaRegistration struct {
	Prefix    string
	Slug      string
	Directory string
}

// ValidatePersonaOwnership cross-checks the persona directories under root
// against the registered personas (§2.3 persona identity grammar and the
// canonical basename): every registered persona has exactly one directory with
// an agreeing prefix and basename, every non-reserved directory is registered,
// prefixes are unique, and `_shared` is never a persona. Problems are returned
// together, sorted, and fail closed.
func ValidatePersonaOwnership(root string, registrations []PersonaRegistration) error {
	var problems []string
	personaRoot := filepath.Join(root, "personas")
	entries, err := os.ReadDir(personaRoot)
	if err != nil {
		return fmt.Errorf("personas: %w", err)
	}
	directories := map[string]bool{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		directories[entry.Name()] = true
	}
	registered := map[string]PersonaRegistration{}
	byDirectory := map[string]string{}
	for _, registration := range registrations {
		prefix := strings.ToUpper(strings.TrimSpace(registration.Prefix))
		slug := strings.TrimSpace(registration.Slug)
		directory := strings.TrimSpace(registration.Directory)
		if prefix == "" || slug == "" || directory == "" {
			problems = append(problems, "a registered persona is missing its prefix, slug, or directory")
			continue
		}
		if strings.EqualFold(prefix, "_shared") {
			problems = append(problems, "the reserved _shared directory must not be registered as a persona")
			continue
		}
		if expected := prefix + "-" + slug; directory != expected {
			problems = append(problems, fmt.Sprintf("persona %s registers directory %q; the canonical basename is %q", prefix, directory, expected))
			continue
		}
		if prior, ok := registered[prefix]; ok {
			problems = append(problems, fmt.Sprintf("prefix %s is registered twice (%q and %q); reuse is forbidden", prefix, prior.Directory, directory))
			continue
		}
		if prior, ok := byDirectory[directory]; ok {
			problems = append(problems, fmt.Sprintf("directory %s is claimed by prefixes %s and %s", directory, prior, prefix))
			continue
		}
		registered[prefix] = registration
		byDirectory[directory] = prefix
		if !directories[directory] {
			problems = append(problems, fmt.Sprintf("persona %s has no workspace directory personas/%s", prefix, directory))
		}
	}
	for name := range directories {
		if name == "_shared" {
			continue
		}
		parsed, err := ParsePersonaDirectory(name)
		if err != nil {
			problems = append(problems, fmt.Sprintf("personas/%s: %v", name, err))
			continue
		}
		if _, ok := registered[strings.ToUpper(parsed.Prefix)]; !ok {
			problems = append(problems, fmt.Sprintf("personas/%s is not registered in 0A", name))
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("persona ownership validation failed:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
}
