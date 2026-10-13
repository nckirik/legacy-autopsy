package assurance

import (
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// DisabledUnit is one effective-denominator disabled unit that requires exactly
// one canonical CAP membership.
type DisabledUnit struct {
	FieldID string
	ID      string
}

// DisabledStats reports the disabled-capability governance state.
type DisabledStats struct {
	DisabledUnits []DisabledUnit
	SRCTotal      int
	CMPTotal      int
	UCTotal       int
	CAPs          int
	Uncovered     []string
	Overlapping   []string
	Unbound       []string
}

// DisabledCoverage evaluates §15.1 condition 10: every effective-denominator
// Disabled SRC, CMP, and UC maps to exactly one canonical CAP whose target
// decision is bound. A resolved zero denominator is reported, never assumed.
func DisabledCoverage(db string, doc *cdl.EIRDoc) (DisabledStats, error) {
	var stats DisabledStats
	disabled, err := disabledUnits(db, doc)
	if err != nil {
		return stats, err
	}
	stats.DisabledUnits = disabled
	srcAll, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return stats, err
	}
	cmpAll, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		return stats, err
	}
	ucAll, err := store.List(db, doc, "USE-CASE")
	if err != nil {
		return stats, err
	}
	stats.SRCTotal, stats.CMPTotal, stats.UCTotal = len(srcAll), len(cmpAll), len(ucAll)
	capNames, _, _, err := work.FieldContract(doc, "CAPABILITY-RECORD")
	if err != nil {
		return stats, err
	}
	capIDs, err := store.List(db, doc, "CAPABILITY-RECORD")
	if err != nil {
		return stats, err
	}
	sort.Strings(capIDs)
	membership := map[string][]string{}
	bound := map[string]bool{}
	for _, capID := range capIDs {
		payload, err := store.Get(db, doc, "CAPABILITY-RECORD", capID)
		if err != nil {
			return stats, err
		}
		values, err := work.PayloadText(capNames, payload)
		if err != nil {
			return stats, err
		}
		for _, member := range splitSet(values["canonical-member-src-cmp-uc-if-cfg-dr-ids"]) {
			membership[member] = append(membership[member], capID)
		}
		if isBoundCap(db, doc, values) {
			bound[capID] = true
		}
	}
	stats.CAPs = len(capIDs)
	requiredCaps := map[string]bool{}
	for _, unit := range disabled {
		members := membership[unit.ID]
		switch {
		case len(members) == 0:
			stats.Uncovered = append(stats.Uncovered, unit.ID)
		case len(members) > 1:
			stats.Overlapping = append(stats.Overlapping, unit.ID)
			for _, capID := range members {
				requiredCaps[capID] = true
			}
		default:
			requiredCaps[members[0]] = true
		}
	}
	for capID := range requiredCaps {
		if !bound[capID] {
			stats.Unbound = append(stats.Unbound, capID)
		}
	}
	sort.Strings(stats.Uncovered)
	sort.Strings(stats.Overlapping)
	sort.Strings(stats.Unbound)
	return stats, nil
}

// disabledUnits enumerates committed units whose declared state is exactly
// Disabled: SRC traversal tracks, CMP capability state, and UC capability
// classification.
func disabledUnits(db string, doc *cdl.EIRDoc) ([]DisabledUnit, error) {
	var units []DisabledUnit
	srcNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	srcIDs, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	for _, id := range srcIDs {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(srcNames, payload)
		if err != nil {
			return nil, err
		}
		if hasTrack(values["traversal-tracks"], "disabled") {
			units = append(units, DisabledUnit{FieldID: "SOURCE-INVENTORY-RECORD", ID: id})
		}
	}
	cmpNames, _, _, err := work.FieldContract(doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	cmpIDs, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	for _, id := range cmpIDs {
		payload, err := store.Get(db, doc, "ATOMIC-COMPONENT", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(cmpNames, payload)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(strings.TrimSpace(values["capability-state"]), "disabled") {
			units = append(units, DisabledUnit{FieldID: "ATOMIC-COMPONENT", ID: id})
		}
	}
	ucNames, _, _, err := work.FieldContract(doc, "USE-CASE")
	if err != nil {
		return nil, err
	}
	ucIDs, err := store.List(db, doc, "USE-CASE")
	if err != nil {
		return nil, err
	}
	for _, id := range ucIDs {
		payload, err := store.Get(db, doc, "USE-CASE", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(ucNames, payload)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(strings.TrimSpace(values["capability-classification"]), "disabled") {
			units = append(units, DisabledUnit{FieldID: "USE-CASE", ID: id})
		}
	}
	sort.Slice(units, func(i, j int) bool {
		if units[i].FieldID != units[j].FieldID {
			return units[i].FieldID < units[j].FieldID
		}
		return units[i].ID < units[j].ID
	})
	return units, nil
}

// isBoundCap reports whether a CAP declares a Bound target decision with a
// resolving decision record.
func isBoundCap(db string, doc *cdl.EIRDoc, values map[string]string) bool {
	if !strings.EqualFold(strings.TrimSpace(values["target-decision-state"]), "Bound") {
		return false
	}
	dec := firstTokenField(values["target-decision-dec"])
	if dec == "" || strings.EqualFold(dec, "none") {
		return false
	}
	exists, err := store.Exists(db, doc, "DECISION", dec)
	return err == nil && exists
}

func hasTrack(raw, want string) bool {
	for _, token := range strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == '|' }) {
		if strings.EqualFold(strings.TrimSpace(token), want) {
			return true
		}
	}
	return false
}

func splitSet(raw string) []string {
	var out []string
	for _, token := range strings.Split(raw, ",") {
		token = firstTokenField(strings.TrimSpace(token))
		if token == "" || strings.EqualFold(token, "none") {
			continue
		}
		out = append(out, token)
	}
	return out
}
