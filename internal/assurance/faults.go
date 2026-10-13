package assurance

import (
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// FaultStats reports fault/resilience topology integrity.
type FaultStats struct {
	FLTs           int
	Unbound        []string
	Dangling       []string
	UnresolvedRefs []string
}

var (
	typedTokenPattern = regexp.MustCompile(`\b(CLM|BR|UC|IF|SCHED|DEP|CMP|ER|MOD|SM|FLT|SEC|NFR|REL|CFG|PRF|PRV|CAP|DR)-[A-F0-9]{12,64}\b`)
	srcTokenPattern   = regexp.MustCompile(`\bSRC-[A-Z]+-[A-F0-9]{12,64}\b`)
	fltTokenPattern   = regexp.MustCompile(`\bFLT-[A-F0-9]{12,64}\b`)
)

var tokenField = map[string]string{
	"CLM":   "CLAIM",
	"BR":    "BUSINESS-RULE-RECORD",
	"UC":    "USE-CASE",
	"IF":    "INTERFACE-RECORD",
	"SCHED": "SCHEDULE-RECORD",
	"DEP":   "DEPLOYMENT-ELEMENT",
	"CMP":   "ATOMIC-COMPONENT",
	"ER":    "ENTITY-RECORD",
	"MOD":   "ARCHITECTURE-MODULE",
	"SM":    "STATE-MACHINE",
	"FLT":   "FAULT-RECORD",
	"SEC":   "SECURITY-FINDING",
	"NFR":   "NFR-RECORD",
	"REL":   "RELATIONSHIP-RECORD",
	"CFG":   "CONFIGURATION-RECORD",
	"PRF":   "PERSONA-PROFILE",
	"PRV":   "PRIVACY-FINDING",
	"CAP":   "CAPABILITY-RECORD",
	"DR":    "DATABASE-ROUTINE",
}

// FaultCoverage evaluates the deterministic portions of §15.1 condition 11:
// every committed FLT record carries claim evidence and every fault-class,
// propagation-path, and cross-record reference resolves. A workspace with no
// FLT records fails closed rather than claiming an unproven zero domain.
func FaultCoverage(db string, doc *cdl.EIRDoc) (FaultStats, error) {
	var stats FaultStats
	names, _, _, err := work.FieldContract(doc, "FAULT-RECORD")
	if err != nil {
		return stats, err
	}
	ids, err := store.List(db, doc, "FAULT-RECORD")
	if err != nil {
		return stats, err
	}
	sort.Strings(ids)
	stats.FLTs = len(ids)
	bound, err := claimBoundRecords(db, doc, "FAULT-RECORD")
	if err != nil {
		return stats, err
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if !bound[id] {
			stats.Unbound = append(stats.Unbound, id)
		}
		payload, err := store.Get(db, doc, "FAULT-RECORD", id)
		if err != nil {
			return stats, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return stats, err
		}
		for _, field := range []string{"failure-class-trigger-detection-point", "propagation-path", "related-rules-claims"} {
			for _, token := range typedTokens(values[field]) {
				if seen[token] {
					continue
				}
				seen[token] = true
				resolved, err := tokenResolves(db, doc, token)
				if err != nil {
					return stats, err
				}
				if !resolved {
					stats.Dangling = append(stats.Dangling, token)
				}
			}
		}
	}
	modNames, _, _, err := work.FieldContract(doc, "ARCHITECTURE-MODULE")
	if err != nil {
		return stats, err
	}
	modIDs, err := store.List(db, doc, "ARCHITECTURE-MODULE")
	if err != nil {
		return stats, err
	}
	for _, id := range modIDs {
		payload, err := store.Get(db, doc, "ARCHITECTURE-MODULE", id)
		if err != nil {
			return stats, err
		}
		values, err := work.PayloadText(modNames, payload)
		if err != nil {
			return stats, err
		}
		for _, token := range fltTokenPattern.FindAllString(values["failure-propagation"], -1) {
			if seen[token] {
				continue
			}
			seen[token] = true
			exists, err := store.Exists(db, doc, "FAULT-RECORD", token)
			if err != nil {
				return stats, err
			}
			if !exists {
				stats.UnresolvedRefs = append(stats.UnresolvedRefs, token)
			}
		}
	}
	sort.Strings(stats.Unbound)
	sort.Strings(stats.Dangling)
	sort.Strings(stats.UnresolvedRefs)
	return stats, nil
}

// claimBoundRecords returns the derived record IDs that have at least one
// resolving CLAIM binding edge.
func claimBoundRecords(db string, doc *cdl.EIRDoc, fieldID string) (map[string]bool, error) {
	names, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	bound := map[string]bool{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		if !strings.EqualFold(strings.TrimSpace(values["binding-kind"]), "CLAIM") {
			continue
		}
		record := firstTokenField(values["derived-record-id"])
		claim := firstTokenField(values["claim-id"])
		if record == "" || claim == "" || strings.EqualFold(claim, "none") {
			continue
		}
		exists, err := store.Exists(db, doc, "CLAIM", claim)
		if err != nil {
			return nil, err
		}
		if exists {
			bound[record] = true
		}
	}
	return bound, nil
}

func typedTokens(value string) []string {
	tokens := typedTokenPattern.FindAllString(value, -1)
	tokens = append(tokens, srcTokenPattern.FindAllString(value, -1)...)
	return tokens
}

func tokenResolves(db string, doc *cdl.EIRDoc, token string) (bool, error) {
	if strings.HasPrefix(token, "SRC-") {
		return store.Exists(db, doc, "SOURCE-INVENTORY-RECORD", token)
	}
	prefix, _, ok := strings.Cut(token, "-")
	if !ok {
		return true, nil
	}
	fieldID, ok := tokenField[prefix]
	if !ok {
		return true, nil
	}
	return store.Exists(db, doc, fieldID, token)
}
