package work

import (
	"fmt"
	"strings"
)

// validateInvocationOwnership enforces persona isolation at commit time: a
// persona-bound invocation may only create records that claim its own persona
// as canonical owner, or records whose ownership is SHARED or None (§2.2
// canonical entry ownership; §8.2 strict single-scope bindings: one invocation
// may not switch roots, personas, POVs, or tracks).
func validateInvocationOwnership(unit Unit, values map[string]string) error {
	persona := strings.ToUpper(strings.TrimSpace(unit.Persona))
	if persona == "" {
		return nil
	}
	check := func(field, raw string) error {
		owner := strings.ToUpper(strings.TrimSpace(raw))
		if owner == "" || owner == "NONE" || owner == "SHARED" {
			return nil
		}
		prefix := strings.TrimSpace(strings.SplitN(owner, "/", 2)[0])
		if prefix == persona {
			return nil
		}
		return fmt.Errorf("work: unit %s is bound to persona %s but %s claims owner %q; an invocation may not write another persona's record", unit.ID, persona, field, raw)
	}
	if err := check("canonical-owner", values["canonical-owner"]); err != nil {
		return err
	}
	return check("canonical-owner-persona", values["canonical-owner-persona"])
}
