// Package handoff evaluates the §15.1 candidate conditions that the current
// canonical state can prove deterministically. Conditions whose models do not
// exist yet are reported Unsupported, never silently passed.
package handoff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/handbook"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// CandidateCheck is one §15.1 candidate condition result.
type CandidateCheck struct {
	ID     string
	Result string // Pass | Fail | Unsupported
	Detail string
}

// Options configures candidate evaluation inputs that are projections rather
// than store state.
type Options struct {
	// HandbookDir, when set, contains chapter-NN.md files to lint for
	// CANDIDATE-04.
	HandbookDir string
}

// EvaluateCandidateChecks evaluates the candidate report check registry against
// committed state.
func EvaluateCandidateChecks(db string, doc *cdl.EIRDoc, opts Options) ([]CandidateCheck, error) {
	var checks []CandidateCheck
	add := func(id, result, detail string) {
		checks = append(checks, CandidateCheck{ID: id, Result: result, Detail: detail})
	}

	// CANDIDATE-01: Exit A green with a pinned composite assertion.
	gate, err := assurance.EvaluateExitA(db, doc)
	if err != nil {
		return nil, err
	}
	if gate.Passed {
		add("CANDIDATE-01-EXIT-A", "Pass", "all deterministic Exit A conditions pass")
	} else {
		add("CANDIDATE-01-EXIT-A", "Fail", "Exit A is not green")
	}

	// CANDIDATE-02: source closure - every source unit covered and dispositioned.
	inventoryIDs, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	summary, err := assurance.CoverageSummary(db, doc)
	if err != nil {
		return nil, err
	}
	unmapped := summary.CoverageByDisposition["Unmapped"]
	if len(inventoryIDs) == 0 {
		add("CANDIDATE-02-SOURCE-CLOSURE", "Fail", "no source inventory records exist")
	} else if len(coverageIDs) != len(inventoryIDs) {
		add("CANDIDATE-02-SOURCE-CLOSURE", "Fail", fmt.Sprintf("%d source units but %d coverage rows", len(inventoryIDs), len(coverageIDs)))
	} else if unmapped > 0 {
		add("CANDIDATE-02-SOURCE-CLOSURE", "Fail", fmt.Sprintf("%d source units remain Unmapped", unmapped))
	} else {
		add("CANDIDATE-02-SOURCE-CLOSURE", "Pass", fmt.Sprintf("%d source units covered and dispositioned", len(inventoryIDs)))
	}

	// CANDIDATE-03: synthesis domains - every mapped unit has provenance.
	mappedUnits, unboundUnits, err := synthesisCoverage(db, doc, coverageIDs)
	if err != nil {
		return nil, err
	}
	if mappedUnits == 0 {
		add("CANDIDATE-03-SYNTHESIS-DOMAINS", "Fail", "no mapped source units")
	} else if unboundUnits > 0 {
		add("CANDIDATE-03-SYNTHESIS-DOMAINS", "Fail", fmt.Sprintf("%d mapped units lack synthesis provenance", unboundUnits))
	} else {
		add("CANDIDATE-03-SYNTHESIS-DOMAINS", "Pass", fmt.Sprintf("%d mapped units have synthesis provenance", mappedUnits))
	}

	if opts.HandbookDir == "" {
		add("CANDIDATE-04-HANDBOOK-QUALITY", "Unsupported", "no rendered handbook directory was provided")
	} else {
		chapters, err := handbook.LoadIndexEntries(opts.HandbookDir)
		if err != nil {
			return nil, err
		}
		expected := len(handbook.DefaultChapters())
		if len(chapters) < expected {
			add("CANDIDATE-04-HANDBOOK-QUALITY", "Fail", fmt.Sprintf("handbook has %d of %d chapters", len(chapters), expected))
		} else if len(chapters) == 0 {
			add("CANDIDATE-04-HANDBOOK-QUALITY", "Fail", "handbook directory has no chapter files")
		} else {
			var findings []string
			for _, chapter := range chapters {
				content, err := os.ReadFile(filepath.Join(opts.HandbookDir, chapter.Path))
				if err != nil {
					return nil, err
				}
				for _, finding := range handbook.LintChapter(chapter.Number, string(content)) {
					findings = append(findings, finding.String())
				}
			}
			if len(findings) > 0 {
				add("CANDIDATE-04-HANDBOOK-QUALITY", "Fail", fmt.Sprintf("%d handbook findings; first: %s", len(findings), findings[0]))
			} else {
				add("CANDIDATE-04-HANDBOOK-QUALITY", "Pass", fmt.Sprintf("%d handbook chapters pass substance and STE gates", len(chapters)))
			}
		}
	}

	// CANDIDATE-05: referential integrity of committed supporting records.
	if err := referentialIntegrity(db, doc); err != nil {
		add("CANDIDATE-05-REFERENTIAL-INTEGRITY", "Fail", err.Error())
	} else {
		add("CANDIDATE-05-REFERENTIAL-INTEGRITY", "Pass", "edges, claims, and matrix rows reference committed records")
	}

	add("CANDIDATE-06-LIFECYCLE", "Unsupported", "lifecycle eligibility records are not generated yet")

	// CANDIDATE-07: use cases carry a resolved capability classification.
	useCases, detail, err := useCaseCheck(db, doc)
	if err != nil {
		return nil, err
	}
	add("CANDIDATE-07-USE-CASES", useCases, detail)

	// CANDIDATE-08: configurations have a behavioral effect or exclusion.
	configs, detail, err := configurationCheck(db, doc)
	if err != nil {
		return nil, err
	}
	add("CANDIDATE-08-CONFIGURATION", configs, detail)

	add("CANDIDATE-09-PRIVACY-SECURITY", "Unsupported", "privacy finding mirrors are not generated yet")
	add("CANDIDATE-10-DISABLED-CAPABILITIES", "Unsupported", "capability containers for disabled units are not generated yet")
	add("CANDIDATE-11-FAULTS", "Unsupported", "failure-class registry coverage is not generated yet")
	add("CANDIDATE-12-PROFILES", "Unsupported", "persona profile payloads are not generated yet")
	pending, contradictions, err := pendingAndContradictions(db, doc)
	if err != nil {
		return nil, err
	}
	if pending > 0 {
		add("CANDIDATE-13-PENDING-AND-CONTRADICTIONS", "Fail", fmt.Sprintf("%d synthesis records carry pending or provisional status", pending))
	} else {
		add("CANDIDATE-13-PENDING-AND-CONTRADICTIONS", "Pass", fmt.Sprintf("no pending synthesis records; %d contradictions recorded", contradictions))
	}

	// CANDIDATE-14: snapshot consistency of committed coverage state.
	snapshot, err := headerSnapshot(db, doc)
	if err != nil {
		return nil, err
	}
	mismatched, err := snapshotMismatches(db, doc, coverageIDs, snapshot)
	if err != nil {
		return nil, err
	}
	if snapshot == "" {
		add("CANDIDATE-14-SNAPSHOT-INPUTS", "Fail", "preflight header declares no environment/snapshot")
	} else if mismatched > 0 {
		add("CANDIDATE-14-SNAPSHOT-INPUTS", "Fail", fmt.Sprintf("%d coverage rows use another snapshot", mismatched))
	} else {
		add("CANDIDATE-14-SNAPSHOT-INPUTS", "Pass", "committed coverage rows match the pinned snapshot")
	}

	sort.Slice(checks, func(i, j int) bool { return checks[i].ID < checks[j].ID })
	return checks, nil
}

func synthesisCoverage(db string, doc *cdl.EIRDoc, coverageIDs []string) (int, int, error) {
	edgeNames, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return 0, 0, err
	}
	edgeIDs, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return 0, 0, err
	}
	bound := map[string]bool{}
	for _, id := range edgeIDs {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return 0, 0, err
		}
		values, err := work.PayloadText(edgeNames, payload)
		if err != nil {
			return 0, 0, err
		}
		bound[firstToken(values["source-record-id"])] = true
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, 0, err
	}
	mapped, unbound := 0, 0
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return 0, 0, err
		}
		values, err := work.PayloadText(coverageNames, payload)
		if err != nil {
			return 0, 0, err
		}
		if values["disposition"] != "Mapped" {
			continue
		}
		mapped++
		if !bound[firstToken(values["src-version"])] {
			unbound++
		}
	}
	return mapped, unbound, nil
}

func referentialIntegrity(db string, doc *cdl.EIRDoc) error {
	edgeNames, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return err
	}
	edgeIDs, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return err
	}
	for _, id := range edgeIDs {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return err
		}
		values, err := work.PayloadText(edgeNames, payload)
		if err != nil {
			return err
		}
		exists, err := store.Exists(db, doc, "SOURCE-INVENTORY-RECORD", firstToken(values["source-record-id"]))
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("edge %s references missing source unit %s", id, values["source-record-id"])
		}
	}
	matrixNames, _, _, err := work.FieldContract(doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return err
	}
	matrixIDs, err := store.List(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW")
	if err != nil {
		return err
	}
	for _, id := range matrixIDs {
		payload, err := store.Get(db, doc, "REQUIRED-TRAVERSAL-MATRIX-ROW", id)
		if err != nil {
			return err
		}
		values, err := work.PayloadText(matrixNames, payload)
		if err != nil {
			return err
		}
		exists, err := store.Exists(db, doc, "ENTRY-CLUSTER", values["cluster-id"])
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("matrix row %s references missing cluster %s", id, values["cluster-id"])
		}
	}
	claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return err
	}
	claimIDs, err := store.List(db, doc, "CLAIM")
	if err != nil {
		return err
	}
	for _, id := range claimIDs {
		payload, err := store.Get(db, doc, "CLAIM", id)
		if err != nil {
			return err
		}
		values, err := work.PayloadText(claimNames, payload)
		if err != nil {
			return err
		}
		exists, err := store.Exists(db, doc, "SOURCE-INVENTORY-RECORD", firstToken(values["subject-id-record-version"]))
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("claim %s references missing source unit %s", id, values["subject-id-record-version"])
		}
	}
	return nil
}

func useCaseCheck(db string, doc *cdl.EIRDoc) (string, string, error) {
	names, _, _, err := work.FieldContract(doc, "USE-CASE")
	if err != nil {
		return "", "", err
	}
	ids, err := store.List(db, doc, "USE-CASE")
	if err != nil {
		return "", "", err
	}
	if len(ids) == 0 {
		return "Fail", "no use cases recorded", nil
	}
	unknown := 0
	for _, id := range ids {
		payload, err := store.Get(db, doc, "USE-CASE", id)
		if err != nil {
			return "", "", err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return "", "", err
		}
		if strings.EqualFold(strings.TrimSpace(values["capability-classification"]), "Unknown") {
			unknown++
		}
	}
	if unknown > 0 {
		return "Fail", fmt.Sprintf("%d use cases have Unknown capability classification", unknown), nil
	}
	return "Pass", fmt.Sprintf("%d use cases carry a resolved capability classification", len(ids)), nil
}

func configurationCheck(db string, doc *cdl.EIRDoc) (string, string, error) {
	names, _, _, err := work.FieldContract(doc, "CONFIGURATION-RECORD")
	if err != nil {
		return "", "", err
	}
	ids, err := store.List(db, doc, "CONFIGURATION-RECORD")
	if err != nil {
		return "", "", err
	}
	if len(ids) == 0 {
		return "Fail", "no configuration records recorded", nil
	}
	missing := 0
	for _, id := range ids {
		payload, err := store.Get(db, doc, "CONFIGURATION-RECORD", id)
		if err != nil {
			return "", "", err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return "", "", err
		}
		effect := strings.TrimSpace(values["behavioral-effect"])
		if effect == "" || strings.EqualFold(effect, "None") || strings.EqualFold(effect, "Unknown") {
			missing++
		}
	}
	if missing > 0 {
		return "Fail", fmt.Sprintf("%d configurations lack a behavioral effect or exclusion", missing), nil
	}
	return "Pass", fmt.Sprintf("%d configurations declare a behavioral effect", len(ids)), nil
}

func headerSnapshot(db string, doc *cdl.EIRDoc) (string, error) {
	payload, err := store.Get(db, doc, "PREFLIGHT-HEADER", "0A")
	if err != nil {
		return "", nil // absent state is reported as a failing check, not an error
	}
	names, _, _, err := work.FieldContract(doc, "PREFLIGHT-HEADER")
	if err != nil {
		return "", err
	}
	values, err := work.PayloadText(names, payload)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(values["included-environments-snapshots"]), nil
}

func snapshotMismatches(db string, doc *cdl.EIRDoc, coverageIDs []string, snapshot string) (int, error) {
	names, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	mismatches := 0
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return 0, err
		}
		if strings.TrimSpace(values["snapshot"]) != snapshot {
			mismatches++
		}
	}
	return mismatches, nil
}

// pendingAndContradictions counts synthesis records with pending/provisional
// status and committed contradiction records.
func pendingAndContradictions(db string, doc *cdl.EIRDoc) (int, int, error) {
	pending := 0
	for _, catalog := range synthesis.Catalogs {
		ids, err := store.List(db, doc, catalog.FieldID)
		if err != nil {
			return 0, 0, err
		}
		if len(ids) == 0 {
			continue
		}
		names, _, _, err := work.FieldContract(doc, catalog.FieldID)
		if err != nil {
			return 0, 0, err
		}
		for _, id := range ids {
			payload, err := store.Get(db, doc, catalog.FieldID, id)
			if err != nil {
				return 0, 0, err
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return 0, 0, err
			}
			for _, value := range values {
				lower := strings.ToLower(value)
				if strings.Contains(lower, "pending") || strings.Contains(lower, "provisional") || strings.Contains(lower, "b-stale") {
					pending++
					break
				}
			}
		}
	}
	contradictions, err := store.List(db, doc, "CONTRADICTION")
	if err != nil {
		return 0, 0, err
	}
	return pending, len(contradictions), nil
}

func firstToken(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return strings.TrimSuffix(fields[0], "/")
}
