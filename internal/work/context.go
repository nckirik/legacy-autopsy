package work

import (
	"fmt"
	"strings"
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
	return Context{UnitID: unit.ID, Bytes: []byte(b.String()), Closure: caps.ReadSet()}, nil
}

// BuildReviewContext reconstructs a review view independently of any worker
// context: the proposed payload plus the unit's declared dependencies, read
// fresh from the store through a separate capability surface.
func BuildReviewContext(caps *Caps, unit Unit, target Ref, proposed []byte) (Context, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# Independent review: %s\n\n", unit.ID)
	fmt.Fprintf(&b, "- **Result contract:** %s\n- **Target record:** %s\n\n", unit.ResultType, target)
	fmt.Fprintf(&b, "## Proposed record\n\n```markdown\n%s```\n\n", ensureTrailingNewline(string(proposed)))
	if len(unit.Reads) > 0 {
		fmt.Fprintf(&b, "## Declared dependency closure\n\n")
		for _, ref := range unit.Reads {
			if ref == target {
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
