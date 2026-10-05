package work

import (
	"fmt"

	"github.com/nckirik/legacy-autopsy/internal/store"
)

// Proposal is one executor-authored semantic mutation. It is only a proposal:
// the runtime validates it against the result contract and the invocation's
// read set before anything commits. Record identity and version are runtime
// properties (the unit's target record ID, version 1 on create); the executor
// authors payload content only.
type Proposal struct {
	FieldID string
	Payload []byte
	Effects []string
}

// StaleError reports a dependency whose committed fingerprint changed or
// disappeared after the invocation read it. The mutation is rejected; a fresh
// invocation is required.
type StaleError struct {
	Ref Ref
}

func (e *StaleError) Error() string {
	return fmt.Sprintf("work: stale dependency %s", e.Ref)
}

// Commit validates a proposal against the unit's result contract and allowed
// effects, stale-checks every observed dependency except the proposal's own
// target, and only then ingests it.
func Commit(caps *Caps, unit Unit, p Proposal) error {
	if p.FieldID != unit.ResultType {
		return fmt.Errorf("work: unit %s proposal writes %s, result contract is %s", unit.ID, p.FieldID, unit.ResultType)
	}
	for _, effect := range p.Effects {
		allowed := false
		for _, candidate := range unit.AllowedEffects {
			if candidate == effect {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("work: unit %s effect %q is not allowed", unit.ID, effect)
		}
	}
	target := Ref{FieldID: unit.ResultType, RecordID: unit.RecordID}
	for ref, want := range caps.ReadSet() {
		if ref == target {
			continue
		}
		got, err := store.FingerprintOf(caps.db, caps.doc, ref.FieldID, ref.RecordID)
		if err != nil || got != want {
			return &StaleError{Ref: ref}
		}
	}
	return store.Ingest(caps.db, caps.doc, p.FieldID, unit.RecordID, 1, p.Payload)
}
