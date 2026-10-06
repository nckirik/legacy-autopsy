package work

import (
	"fmt"

	"github.com/nckirik/legacy-autopsy/internal/certify"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// Record is one executor-authored record payload.
type Record struct {
	Payload []byte
}

// Proposal is one executor-authored semantic mutation. It is only a proposal:
// the runtime validates it against the result contract and the invocation's
// read set before anything commits. Record identity and version are runtime
// properties (single-record target or ID derived from a declared payload field,
// version 1 on create); the executor authors payload content only.
type Proposal struct {
	FieldID string
	Records []Record
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
	_, err := CommitWithIDs(caps, unit, p)
	return err
}

// CommitWithIDs is Commit returning the runtime-assigned record IDs.
func CommitWithIDs(caps *Caps, unit Unit, p Proposal) ([]string, error) {
	if p.FieldID != unit.ResultType {
		return nil, fmt.Errorf("work: unit %s proposal writes %s, result contract is %s", unit.ID, p.FieldID, unit.ResultType)
	}
	if len(p.Records) == 0 {
		return nil, fmt.Errorf("work: unit %s proposal has no records", unit.ID)
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
			return nil, fmt.Errorf("work: unit %s effect %q is not allowed", unit.ID, effect)
		}
	}
	resolved, err := ResolveRecords(caps.doc, unit, p)
	if err != nil {
		return nil, err
	}
	names, _, _, err := FieldContract(caps.doc, p.FieldID)
	if err != nil {
		return nil, err
	}
	lookup := func(fieldID, recordID string) (bool, error) {
		return store.Exists(caps.db, caps.doc, fieldID, recordID)
	}
	reader := func(fieldID, recordID string) (map[string]string, error) {
		payload, err := store.Get(caps.db, caps.doc, fieldID, recordID)
		if err != nil {
			return nil, err
		}
		names, _, _, err := FieldContract(caps.doc, fieldID)
		if err != nil {
			return nil, err
		}
		return PayloadText(names, payload)
	}
	for _, record := range resolved {
		values, err := PayloadText(names, record.Payload)
		if err != nil {
			return nil, fmt.Errorf("work: unit %s record %s: %w", unit.ID, record.RecordID, err)
		}
		if err := certify.ValidateReferences(p.FieldID, values, lookup); err != nil {
			return nil, fmt.Errorf("work: unit %s record %s: %w", unit.ID, record.RecordID, err)
		}
		if err := certify.ValidateCrossReferences(p.FieldID, values, reader); err != nil {
			return nil, fmt.Errorf("work: unit %s record %s: %w", unit.ID, record.RecordID, err)
		}
	}
	targets := map[Ref]bool{}
	for _, record := range resolved {
		targets[Ref{FieldID: p.FieldID, RecordID: record.RecordID}] = true
	}
	for ref, want := range caps.ReadSet() {
		if targets[ref] {
			continue
		}
		got, err := store.FingerprintOf(caps.db, caps.doc, ref.FieldID, ref.RecordID)
		if err != nil || got != want {
			return nil, &StaleError{Ref: ref}
		}
	}
	inputs := make([]store.RecordInput, 0, len(resolved))
	for _, record := range resolved {
		inputs = append(inputs, store.RecordInput{RecordID: record.RecordID, Version: 1, Payload: record.Payload})
	}
	if err := store.IngestMany(caps.db, caps.doc, p.FieldID, inputs); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(resolved))
	for _, record := range resolved {
		ids = append(ids, record.RecordID)
	}
	return ids, nil
}
