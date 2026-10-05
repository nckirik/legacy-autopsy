package work

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

// Context is the bounded, cold-replayable projection of one work unit: the
// declared dependency closure, materialized explicitly and nothing else.
type Context struct {
	UnitID  string
	Bytes   []byte
	Closure ReadSet
}

// BuildContext fetches exactly the declared dependency closure through the
// capability surface, so the read set and the materialized projection stay in
// sync. Undeclared records never enter the context.
func BuildContext(caps *Caps, unit Unit) (Context, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# Work unit %s\n\n", unit.ID)
	fmt.Fprintf(&b, "- **Mode:** %s\n- **Scope:** %s\n- **Result contract:** %s\n- **Target record ID:** %s\n\n", unit.Mode, unit.Scope, unit.ResultType, unit.RecordID)
	fmt.Fprintf(&b, "## Declared dependency closure\n\n")
	for _, ref := range unit.Reads {
		payload, err := caps.Get(ref)
		if err != nil {
			return Context{}, fmt.Errorf("work: unit %s dependency %s: %w", unit.ID, ref, err)
		}
		fmt.Fprintf(&b, "### %s\n\n- **Fingerprint:** %s\n\n```markdown\n%s```\n\n", ref, caps.reads[ref], ensureTrailingNewline(string(payload)))
	}
	if len(unit.Sections) > 0 {
		fmt.Fprintf(&b, "\n## Normative protocol sections\n\n")
		for _, reference := range unit.Sections {
			text, err := SectionText(caps.doc, reference)
			if err != nil {
				return Context{}, fmt.Errorf("work: unit %s section %s: %w", unit.ID, reference, err)
			}
			fmt.Fprintf(&b, "<!-- section %s -->\n%s\n\n", reference, text)
		}
	}
	for i, input := range unit.Inputs {
		name := input.Name
		if name == "" {
			name = fmt.Sprintf("input-%d", i+1)
		}
		fmt.Fprintf(&b, "### runtime input: %s\n\n- **Fingerprint:** %s\n\n%s\n", name, capabilities.Hash(input.Bytes), ensureTrailingNewline(string(input.Bytes)))
	}
	return Context{UnitID: unit.ID, Bytes: []byte(b.String()), Closure: caps.ReadSet()}, nil
}

// BuildReviewContext reconstructs a review view independently of any worker
// context: the proposed records plus the unit's declared dependencies, read
// fresh from the store through a separate capability surface.
func BuildReviewContext(caps *Caps, unit Unit, proposed []ResolvedRecord) (Context, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# Independent review: %s\n\n", unit.ID)
	fmt.Fprintf(&b, "- **Result contract:** %s\n- **Proposed records:** %d\n\n", unit.ResultType, len(proposed))
	targets := map[Ref]bool{}
	for _, record := range proposed {
		targets[Ref{FieldID: unit.ResultType, RecordID: record.RecordID}] = true
		fmt.Fprintf(&b, "## Proposed record %s\n\n```markdown\n%s```\n\n", record.RecordID, ensureTrailingNewline(string(record.Payload)))
	}
	if len(unit.Reads) > 0 {
		fmt.Fprintf(&b, "## Declared dependency closure\n\n")
		for _, ref := range unit.Reads {
			if targets[ref] {
				continue
			}
			payload, err := caps.Get(ref)
			if err != nil {
				return Context{}, fmt.Errorf("work: review dependency %s: %w", ref, err)
			}
			fmt.Fprintf(&b, "### %s\n\n- **Fingerprint:** %s\n\n```markdown\n%s```\n\n", ref, caps.reads[ref], ensureTrailingNewline(string(payload)))
		}
	}
	return Context{UnitID: unit.ID, Bytes: []byte(b.String()), Closure: caps.ReadSet()}, nil
}

func ensureTrailingNewline(value string) string {
	if strings.HasSuffix(value, "\n") {
		return value
	}
	return value + "\n"
}
