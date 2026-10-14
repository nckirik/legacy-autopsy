package handbook

import (
	"fmt"
	"strings"
)

// Record is one synthesis record embedded into a handbook chapter.
type Record struct {
	ID      string
	Payload []byte
}

// Chapter is the model-authored narrative plus the deterministic bindings for
// one handbook chapter. Narrative sections are content; the assembler owns
// structure, ordering, and the required-section layout.
type Chapter struct {
	Number          string
	Title           string
	AudiencePurpose string
	Observed        []Record
	TargetDecisions string
	Navigation      []string
}

// Assemble renders chapters in the normative required-section layout (audience
// purpose, observed legacy behavior, target decisions, navigation) with the
// embedded synthesis records. The output is validated by LintChapter; the
// assembler never invents narrative content.
func Assemble(chapters []Chapter) []byte {
	var b strings.Builder
	b.WriteString("# Legacy System Reconstruction Handbook\n\n")
	b.WriteString("> Generated reconstruction projection. Source of truth is the canonical protocol state.\n\n")
	for _, chapter := range chapters {
		fmt.Fprintf(&b, "## %s. %s\n\n", chapter.Number, chapter.Title)
		b.WriteString("### Summary\n\n")
		writeProse(&b, chapter.AudiencePurpose)
		b.WriteString("\n### Observed Legacy Behavior\n\n")
		if len(chapter.Observed) == 0 {
			b.WriteString("None\n\n")
		}
		for _, record := range chapter.Observed {
			fmt.Fprintf(&b, "#### %s\n\n", record.ID)
			writeProse(&b, strings.TrimSpace(string(record.Payload)))
			b.WriteString("\n")
		}
		b.WriteString("\n### Target Decisions\n\n")
		writeProse(&b, chapter.TargetDecisions)
		b.WriteString("\n### Navigation\n\n")
		if len(chapter.Navigation) == 0 {
			b.WriteString("- None\n\n")
		}
		for _, target := range chapter.Navigation {
			fmt.Fprintf(&b, "- %s\n", target)
		}
		b.WriteString("\n---\n\n")
	}
	return []byte(b.String())
}

func writeProse(b *strings.Builder, prose string) {
	prose = strings.TrimSpace(prose)
	if prose == "" {
		b.WriteString("Not-Applicable: no content for this chapter class.\n")
		return
	}
	b.WriteString(prose)
	b.WriteString("\n")
}
