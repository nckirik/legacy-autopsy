package handbook

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

const (
	mintPayloadMarker = "<!-- HANDBOOK-SEMANTIC-PAYLOAD:"
	mintEnvelopeStart = "<!-- CERTIFICATION-ENVELOPE: excluded from semantic version/fingerprint -->"
	mintEnvelopeEnd   = "<!-- END-CERTIFICATION-ENVELOPE -->"
	hbkCarrier        = "SEMANTIC-CONTENT-FINGERPRINT"
)

var (
	anchorPattern   = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	upstreamPattern = regexp.MustCompile(`summary>([A-Z][A-Z0-9]{1,7}-[A-F0-9]{12,64})<`)
	synthesisTypes  = map[string]bool{
		"MOD": true, "ER": true, "REL": true, "SM": true, "DR": true, "BR": true,
		"UC": true, "IF": true, "DEP": true, "CFG": true, "SCHED": true,
		"NFR": true, "SEC": true, "FLT": true, "PRF": true,
	}
)

// UpstreamRef binds one synthesis record to its exact version and semantic
// fingerprint at minting time.
type UpstreamRef struct {
	ID          string
	Version     int
	Fingerprint string
}

// MintedSection is one projection-bound HBK handbook section.
type MintedSection struct {
	Path        string // handbook-relative path (chapter-NN.md)
	Anchor      string
	ID          string
	Title       string
	Version     int
	Fingerprint string
	Upstream    []UpstreamRef
	Bytes       []byte
}

// ResolveUpstream resolves one synthesis record ID to its exact version and
// semantic fingerprint. Unresolvable upstream IDs fail the mint.
type ResolveUpstream func(id string) (UpstreamRef, bool)

// MintSections wraps every chapter file in its HBK semantic payload and
// certification envelope. Re-minting an already minted file preserves its
// anchor, identity, and version, and recomputes the semantic fingerprint.
func MintSections(dir, namespace string, resolve ResolveUpstream) ([]MintedSection, error) {
	names, err := chapterFiles(dir)
	if err != nil {
		return nil, err
	}
	sections := make([]MintedSection, 0, len(names))
	for _, name := range names {
		path := filepath.ToSlash(name)
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		original := string(raw)
		anchor := parseField(original, "STABLE-SECTION-ANCHOR")
		if anchor == "" {
			anchor = strings.TrimSuffix(strings.TrimSuffix(name, ".md"), ".chapter")
		}
		if !anchorPattern.MatchString(anchor) {
			return nil, fmt.Errorf("handbook mint: %s anchor %q is not a stable ASCII anchor", name, anchor)
		}
		version := 1
		if raw := parseField(original, "SEMANTIC-RECORD-VERSION"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 1 {
				return nil, fmt.Errorf("handbook mint: %s has invalid semantic record version %q", name, raw)
			}
			version = parsed
		}
		upstream, err := upstreamRefs(original, resolve)
		if err != nil {
			return nil, fmt.Errorf("handbook mint: %s: %w", name, err)
		}
		id, err := capabilities.TypedID("HBK", namespace, path+"#"+anchor, "handbook-section")
		if err != nil {
			return nil, err
		}
		content := unmint(original)
		title := chapterTitle(content, name)
		var payload strings.Builder
		fmt.Fprintf(&payload, "### %s %s\n\n", id, title)
		fmt.Fprintf(&payload, "- **HBK-ID:** %s\n", id)
		fmt.Fprintf(&payload, "- **HANDBOOK-RELATIVE-PATH:** %s\n", path)
		fmt.Fprintf(&payload, "- **STABLE-SECTION-ANCHOR:** %s\n", anchor)
		fmt.Fprintf(&payload, "- **SEMANTIC-RECORD-VERSION:** %d\n", version)
		fmt.Fprintf(&payload, "- **%s:** PENDING\n", hbkCarrier)
		fmt.Fprintf(&payload, "- **UPSTREAM:** %s\n\n", formatUpstream(upstream))
		fmt.Fprintf(&payload, "%s %s -->\n\n", mintPayloadMarker, id)
		payload.WriteString(strings.Trim(content, "\n") + "\n")
		region := payload.String()
		fingerprint, err := capabilities.SemanticContentFingerprint([]byte(region), hbkCarrier)
		if err != nil {
			return nil, err
		}
		region = strings.Replace(region, "- **"+hbkCarrier+":** PENDING", "- **"+hbkCarrier+":** "+fingerprint, 1)
		var b strings.Builder
		b.WriteString(region)
		if !strings.HasSuffix(region, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(mintEnvelopeStart)
		b.WriteString("\n\n")
		b.WriteString("- **HANDBOOK-STATUS:** B-POPULATED\n")
		b.WriteString("- **Confirmation:** None\n")
		b.WriteString("- **Confirmation Envelope Version:** None\n")
		b.WriteString("- **Confirmation Envelope Audit:** None\n")
		b.WriteString(mintEnvelopeEnd)
		b.WriteString("\n")
		sections = append(sections, MintedSection{
			Path:        path,
			Anchor:      anchor,
			ID:          id,
			Title:       title,
			Version:     version,
			Fingerprint: fingerprint,
			Upstream:    upstream,
			Bytes:       []byte(b.String()),
		})
	}
	return sections, nil
}

// LoadMintedSections reads minted chapter files back and fails closed when a
// declared semantic fingerprint does not recompute.
func LoadMintedSections(dir string) ([]MintedSection, error) {
	names, err := chapterFiles(dir)
	if err != nil {
		return nil, err
	}
	sections := make([]MintedSection, 0, len(names))
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		content := string(raw)
		id := parseField(content, "HBK-ID")
		if id == "" {
			return nil, fmt.Errorf("handbook: %s is not a minted HBK section", name)
		}
		path := parseField(content, "HANDBOOK-RELATIVE-PATH")
		anchor := parseField(content, "STABLE-SECTION-ANCHOR")
		versionRaw := parseField(content, "SEMANTIC-RECORD-VERSION")
		version, err := strconv.Atoi(versionRaw)
		if err != nil || version < 1 {
			return nil, fmt.Errorf("handbook: %s has invalid semantic record version %q", name, versionRaw)
		}
		fingerprint := parseField(content, hbkCarrier)
		if fingerprint == "" {
			return nil, fmt.Errorf("handbook: %s has no semantic content fingerprint", name)
		}
		upstream, err := parseUpstream(parseField(content, "UPSTREAM"))
		if err != nil {
			return nil, fmt.Errorf("handbook: %s: %w", name, err)
		}
		region, err := mintedRegion(content)
		if err != nil {
			return nil, fmt.Errorf("handbook: %s: %w", name, err)
		}
		recomputed, err := capabilities.SemanticContentFingerprint([]byte(region), hbkCarrier)
		if err != nil {
			return nil, err
		}
		if recomputed != fingerprint {
			return nil, fmt.Errorf("handbook: %s semantic fingerprint does not recompute (stale payload)", name)
		}
		sections = append(sections, MintedSection{
			Path:        path,
			Anchor:      anchor,
			ID:          id,
			Title:       chapterTitle(unmint(content), name),
			Version:     version,
			Fingerprint: fingerprint,
			Upstream:    upstream,
			Bytes:       raw,
		})
	}
	return sections, nil
}

// WriteSections writes the minted bytes back to the chapter files.
func WriteSections(dir string, sections []MintedSection) error {
	for _, section := range sections {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(section.Path)), section.Bytes, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func chapterFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "chapter-") || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names, nil
}

// unmint strips the previous HBK header and certification envelope, returning
// the original chapter content.
func unmint(content string) string {
	if start := strings.Index(content, mintEnvelopeStart); start >= 0 {
		if end := strings.Index(content[start:], mintEnvelopeEnd); end >= 0 {
			content = content[:start] + content[start+end+len(mintEnvelopeEnd):]
		}
	}
	if start := strings.Index(content, mintPayloadMarker); start >= 0 {
		if end := strings.Index(content[start:], "-->"); end >= 0 {
			content = content[start+end+len("-->"):]
		}
	}
	return strings.Trim(content, "\n")
}

// mintedRegion returns the semantic payload region: from the HBK heading line
// through the line before the certification envelope.
func mintedRegion(content string) (string, error) {
	start := strings.Index(content, "### ")
	if start < 0 {
		return "", fmt.Errorf("minted section has no heading")
	}
	end := strings.Index(content, mintEnvelopeStart)
	if end < 0 || end < start {
		return "", fmt.Errorf("minted section has no certification envelope")
	}
	return content[start:end], nil
}

func parseField(content, label string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		prefix := "- **" + label + ":**"
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
	}
	return ""
}

func chapterTitle(content, name string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
		}
	}
	return strings.TrimSuffix(name, ".md")
}

func upstreamRefs(content string, resolve ResolveUpstream) ([]UpstreamRef, error) {
	seen := map[string]bool{}
	var ids []string
	for _, match := range upstreamPattern.FindAllStringSubmatch(content, -1) {
		id := match[1]
		prefix := id[:strings.Index(id, "-")]
		if !synthesisTypes[prefix] || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Strings(ids)
	refs := make([]UpstreamRef, 0, len(ids))
	for _, id := range ids {
		ref, ok := resolve(id)
		if !ok {
			return nil, fmt.Errorf("upstream record %s is not committed", id)
		}
		ref.ID = id
		refs = append(refs, ref)
	}
	return refs, nil
}

func formatUpstream(refs []UpstreamRef) string {
	if len(refs) == 0 {
		return "None"
	}
	parts := make([]string, 0, len(refs))
	for _, ref := range refs {
		parts = append(parts, fmt.Sprintf("%s@%d@%s", ref.ID, ref.Version, ref.Fingerprint))
	}
	return strings.Join(parts, "; ")
}

func parseUpstream(value string) ([]UpstreamRef, error) {
	if value == "" || strings.EqualFold(value, "None") {
		return nil, nil
	}
	var refs []UpstreamRef
	for _, part := range strings.Split(value, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		pieces := strings.Split(part, "@")
		if len(pieces) != 3 {
			return nil, fmt.Errorf("upstream binding %q is not ID@version@fingerprint", part)
		}
		version, err := strconv.Atoi(pieces[1])
		if err != nil || version < 1 {
			return nil, fmt.Errorf("upstream binding %q has an invalid version", part)
		}
		refs = append(refs, UpstreamRef{ID: pieces[0], Version: version, Fingerprint: pieces[2]})
	}
	return refs, nil
}
