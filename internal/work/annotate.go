package work

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// typedIDPattern matches a well-formed typed protocol ID.
var typedIDPattern = regexp.MustCompile(`^[A-Z][A-Z0-9-]*-[A-F0-9]{12,64}$`)

// annotateRuntimeFields rewrites the named payload fields so every record ID
// token carries its exact committed version and fingerprint
// (`ID@version N sha256:...`). Fingerprints are runtime-owned data: executors
// list bare IDs (optionally with prose suffixes), and any model-authored
// annotation is replaced with the committed truth. Unresolvable IDs fail
// closed.
func annotateRuntimeFields(caps *Caps, fields []string, payload []byte) ([]byte, error) {
	for _, field := range fields {
		rewritten, err := annotateField(caps, field, payload)
		if err != nil {
			return nil, err
		}
		payload = rewritten
	}
	return payload, nil
}

func annotateField(caps *Caps, field string, payload []byte) ([]byte, error) {
	lines := strings.Split(string(payload), "\n")
	column := capabilities.ColumnName(field)
	found := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- **") {
			continue
		}
		end := strings.Index(trimmed[2:], ":**")
		if end < 0 {
			continue
		}
		label := trimmed[4 : end+2]
		if capabilities.ColumnName(label) != column {
			continue
		}
		found++
		value := strings.TrimSpace(trimmed[end+6:])
		if value == "" || strings.EqualFold(value, "None") {
			continue
		}
		var annotated []string
		for _, token := range splitAnnotationTokens(value) {
			id, prose := splitAnnotationToken(token)
			if id == "" {
				continue
			}
			annotation, err := resolveAnnotation(caps, id)
			if err != nil {
				return nil, err
			}
			if prose != "" {
				annotation += " " + prose
			}
			annotated = append(annotated, annotation)
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		lines[i] = indent + "- **" + label + ":** " + strings.Join(annotated, ", ")
	}
	if found == 0 {
		return nil, fmt.Errorf("runtime annotated field %q is missing; emit it with the listed record IDs", field)
	}
	return []byte(strings.Join(lines, "\n")), nil
}

// splitAnnotationTokens splits a list value on comma and semicolon boundaries
// that are outside parentheses, so prose suffixes keep their punctuation.
func splitAnnotationTokens(value string) []string {
	var tokens []string
	depth := 0
	var b strings.Builder
	for _, r := range value {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',', ';':
			if depth == 0 {
				if token := strings.TrimSpace(b.String()); token != "" {
					tokens = append(tokens, token)
				}
				b.Reset()
				continue
			}
		}
		b.WriteRune(r)
	}
	if token := strings.TrimSpace(b.String()); token != "" {
		tokens = append(tokens, token)
	}
	return tokens
}

// splitAnnotationToken separates the leading record ID from any model-authored
// annotation or prose suffix and drops the annotation; the prose is preserved.
func splitAnnotationToken(token string) (string, string) {
	token = strings.TrimSpace(token)
	id := token
	prose := ""
	if idx := strings.IndexAny(token, " ("); idx > 0 {
		id = strings.TrimSpace(token[:idx])
		prose = strings.TrimSpace(token[idx:])
	}
	if at := strings.Index(id, "@"); at > 0 {
		id = strings.TrimSpace(id[:at])
	}
	return strings.Trim(id, "[]`"), prose
}

// resolveAnnotation resolves one record ID to `ID@version N fingerprint`.
func resolveAnnotation(caps *Caps, id string) (string, error) {
	if field, record, ok := strings.Cut(id, "/"); ok && strings.Contains(field, "-") {
		version, err := store.VersionOf(caps.db, caps.doc, field, record)
		if err != nil {
			return "", fmt.Errorf("runtime annotation %s does not resolve", id)
		}
		fingerprint, err := store.FingerprintOf(caps.db, caps.doc, field, record)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s@version %d %s", id, version, fingerprint), nil
	}
	if !typedIDPattern.MatchString(id) {
		return "", fmt.Errorf("runtime annotation %s is not a typed record reference", id)
	}
	// Typed IDs are globally unique; resolve the owning table by lookup order.
	for _, field := range caps.doc.Declarations.Fields {
		switch field.Storage {
		case "record", "registry", "singleton", "log":
		default:
			continue
		}
		if _, err := store.Get(caps.db, caps.doc, field.ID, id); err != nil {
			continue
		}
		version, err := store.VersionOf(caps.db, caps.doc, field.ID, id)
		if err != nil {
			return "", err
		}
		fingerprint, err := store.FingerprintOf(caps.db, caps.doc, field.ID, id)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s@version %d %s", id, version, fingerprint), nil
	}
	return "", fmt.Errorf("runtime annotation %s does not resolve to a committed record", id)
}
