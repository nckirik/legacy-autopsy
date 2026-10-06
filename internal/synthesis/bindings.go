package synthesis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// BindingPlan returns the mechanical unit that materializes DERIVED-FROM-EDGE
// records from each committed synthesis record's declared `source-bindings`.
func BindingPlan(namespace string) work.Unit {
	return work.Unit{
		ID:           "U-SYNTHESIS-BINDINGS",
		Mode:         "Partial-Synthesis",
		Scope:        namespace,
		IDFromFields: []string{"derived-record-id", "source-record-id", "source-version", "binding-kind"},
		ResultType:   "DERIVED-FROM-EDGE",
		Sections:     []string{"dependency-provenance"},
	}
}

// SynthesisRecord is one committed catalog record in bounded form.
type SynthesisRecord struct {
	ID       string
	FieldID  string
	Title    string
	Bindings []string
}

// LoadSynthesisRecords reads every committed catalog record with its declared
// source bindings, sorted for determinism.
func LoadSynthesisRecords(db string, doc *cdl.EIRDoc) ([]SynthesisRecord, error) {
	var records []SynthesisRecord
	for catalogID, catalog := range Catalogs {
		_ = catalogID
		names, _, _, err := work.FieldContract(doc, catalog.FieldID)
		if err != nil {
			return nil, err
		}
		ids, err := store.List(db, doc, catalog.FieldID)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			payload, err := store.Get(db, doc, catalog.FieldID, id)
			if err != nil {
				return nil, err
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return nil, err
			}
			var bindings []string
			for _, token := range strings.Split(values["source-bindings"], ",") {
				token = strings.TrimSpace(token)
				if token == "" || strings.EqualFold(token, "none") {
					continue
				}
				bindings = append(bindings, token)
			}
			sort.Strings(bindings)
			records = append(records, SynthesisRecord{ID: id, FieldID: catalog.FieldID, Title: values["record-title"], Bindings: bindings})
		}
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	return records, nil
}

// LoadInventoryFingerprints maps committed source-unit IDs to their fingerprint
// and version line for edge metadata.
func LoadInventoryFingerprints(db string, doc *cdl.EIRDoc) (map[string]string, error) {
	names, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		out[id] = values["source-fingerprint-snapshot-src-version"]
	}
	return out, nil
}

// MechanicalBindingExecutor materializes one DERIVED-FROM-EDGE per declared
// binding. Source-unit bindings carry the unit's version line; claim bindings
// reference the claim by ID; unknown tokens fail closed.
type MechanicalBindingExecutor struct {
	Records   []SynthesisRecord
	Inventory map[string]string
	// Claims maps each committed claim ID to its subject source unit.
	Claims   map[string]string
	Snapshot string
}

// Execute implements work.Executor.
func (e MechanicalBindingExecutor) Execute(invocation work.Invocation) (work.Proposal, error) {
	proposal := work.Proposal{FieldID: "DERIVED-FROM-EDGE"}
	for _, record := range e.Records {
		for _, token := range record.Bindings {
			sourceID := token
			version := ""
			claimID := "None"
			kind := "SOURCE-UNIT"
			if strings.HasPrefix(token, "CLM-") {
				subject, ok := e.Claims[token]
				if !ok || subject == "" {
					return work.Proposal{}, fmt.Errorf("synthesis: %s binds unknown claim %s", record.ID, token)
				}
				versionLine, ok := e.Inventory[subject]
				if !ok {
					return work.Proposal{}, fmt.Errorf("synthesis: claim %s subject %s is not a committed source unit", token, subject)
				}
				claimID = token
				kind = "CLAIM"
				sourceID = subject
				version = strings.Fields(versionLine)[0] + " version 1"
			} else {
				versionLine, ok := e.Inventory[token]
				if !ok {
					return work.Proposal{}, fmt.Errorf("synthesis: %s binds unknown source unit %s", record.ID, token)
				}
				version = strings.Fields(versionLine)[0] + " version 1"
				sourceID = token
			}
			payload := "- **derived-record-id:** " + record.ID + "\n" +
				"- **derived-field:** " + record.FieldID + "\n" +
				"- **source-record-id:** " + sourceID + "\n" +
				"- **source-version:** " + versionOrNone(version) + "\n" +
				"- **claim-id:** " + claimID + "\n" +
				"- **binding-kind:** " + kind + "\n" +
				"- **evidence-route-lineage-id:** static-filesystem\n" +
				"- **created-by-invocation:** " + invocation.InvocationID + "\n" +
				"- **snapshot-fingerprint:** " + e.Snapshot + "\n" +
				"- **binding-notes:** declared source-bindings on " + record.FieldID + "\n"
			proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
		}
	}
	return proposal, nil
}

func versionOrNone(version string) string {
	if strings.TrimSpace(version) == "" {
		return "None"
	}
	return version
}
