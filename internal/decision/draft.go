// Package decision drafts model-authored decision content from observed inputs
// and commits human-approved decisions. Drafts are candidate content for human
// decision review; only the reviewed commit path creates canonical DECISION
// records, and its approval fields come from the recorded human review.
package decision

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Completer is one bounded model call (satisfied by the llm adapter).
type Completer interface {
	Complete(ctx context.Context, system, user string) (string, error)
}

// Target is one committed subject whose decision reference is missing.
type Target struct {
	SubjectID     string
	SubjectField  string // FIELD ID of the subject record
	DecisionType  string // scope | contradiction | target | migration | capability
	SubjectLine   string // "<FIELD>/<ID> version <v> fingerprint <fp>"
	ObservedLines []string
}

// Targets enumerates the decision subjects deterministically: synthesis and
// capability records whose target decision is unbound, and contradictions whose
// resolution requires a decision.
func Targets(db string, doc *cdl.EIRDoc) ([]Target, error) {
	var targets []Target
	collect := func(fieldID, decisionField string, decisionType string, claimsField string, skip func(values map[string]string) bool) error {
		names, _, _, err := work.FieldContract(doc, fieldID)
		if err != nil {
			return err
		}
		ids, err := store.List(db, doc, fieldID)
		if err != nil {
			return err
		}
		sort.Strings(ids)
		for _, id := range ids {
			payload, err := store.Get(db, doc, fieldID, id)
			if err != nil {
				return err
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return err
			}
			if skip != nil && skip(values) {
				continue
			}
			version, err := store.VersionOf(db, doc, fieldID, id)
			if err != nil {
				return err
			}
			fingerprint, err := store.FingerprintOf(db, doc, fieldID, id)
			if err != nil {
				return err
			}
			target := Target{
				SubjectID:    id,
				SubjectField: fieldID,
				DecisionType: decisionType,
				SubjectLine:  fmt.Sprintf("%s/%s version %d %s", fieldID, id, version, fingerprint),
			}
			target.ObservedLines = append(target.ObservedLines, target.SubjectLine)
			for _, token := range splitTokens(values[claimsField]) {
				if !strings.HasPrefix(token, "CLM-") {
					continue
				}
				claimPayload, err := store.Get(db, doc, "CLAIM", token)
				if err != nil {
					continue
				}
				claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
				if err != nil {
					return err
				}
				claimValues, err := work.PayloadText(claimNames, claimPayload)
				if err != nil {
					return err
				}
				claimVersion, err := store.VersionOf(db, doc, "CLAIM", token)
				if err != nil {
					return err
				}
				claimFingerprint, err := store.FingerprintOf(db, doc, "CLAIM", token)
				if err != nil {
					return err
				}
				target.ObservedLines = append(target.ObservedLines, fmt.Sprintf("CLAIM/%s version %d %s", token, claimVersion, claimFingerprint))
				_ = claimValues
			}
			targets = append(targets, target)
		}
		return nil
	}
	unbound := func(values map[string]string) bool {
		value := strings.TrimSpace(values["target-decision-dec"])
		return value == "" || strings.EqualFold(value, "None")
	}
	if err := collect("ARCHITECTURE-MODULE", "target-decision", "target", "source-bindings", func(values map[string]string) bool {
		value := strings.TrimSpace(values["target-decision"])
		return value != "" && !strings.EqualFold(value, "None")
	}); err != nil {
		return nil, err
	}
	if err := collect("USE-CASE", "target-decision-dec", "target", "source-bindings", unbound); err != nil {
		return nil, err
	}
	if err := collect("CAPABILITY-RECORD", "target-decision-dec", "capability", "observed-claims", func(values map[string]string) bool {
		return !unbound(values)
	}); err != nil {
		return nil, err
	}
	if err := collect("CONTRADICTION", "resolution-decision-id", "contradiction", "claims-in-conflict", func(values map[string]string) bool {
		state := strings.TrimSpace(values["resolution-state"])
		decided := strings.TrimSpace(values["resolution-decision-id"])
		return !strings.EqualFold(state, "Decision-Required") || (decided != "" && !strings.EqualFold(decided, "None"))
	}); err != nil {
		return nil, err
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].DecisionType != targets[j].DecisionType {
			return targets[i].DecisionType < targets[j].DecisionType
		}
		return targets[i].SubjectID < targets[j].SubjectID
	})
	return targets, nil
}

const draftSystemPrompt = `You draft one modernization or scope decision for a legacy-system reconstruction. A human reviewer will approve, change, or reject your draft, so state the observed inputs, the realistic options, the recommended decision, the rationale, and the risks precisely. Do not invent facts beyond the observed inputs. Where the evidence is silent, say Unknown instead of guessing.

Output exactly these six sections and nothing else:

## Options
<the realistic alternatives, including keeping the status quo>

## Decision
<one recommended choice; when the subject type is a use case, the choice must be one of: Preserve, Change-By-Approved-Migration, Redesign, Retire>

## Rationale
<why this choice follows from the observed inputs>

## Risks
<residual risks, unknowns, and what could invalidate the decision>

## Exact Scope and Denominator Effect
<which source units, coverage rows, or gate conditions this decision affects>

## Affected Records / Rollout / Rollback
<affected record IDs, staged rollout, and how to reverse the decision>`

func draftUserPrompt(target Target) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Subject: %s/%s\n", target.SubjectField, target.SubjectID)
	fmt.Fprintf(&b, "Decision type: %s\n\n", target.DecisionType)
	b.WriteString("Observed inputs (authoritative IDs, versions, fingerprints):\n\n")
	for _, line := range target.ObservedLines {
		fmt.Fprintf(&b, "- %s\n", line)
	}
	b.WriteString("\nThe subject record's committed fields are the observed legacy state. Draft the decision now.\n")
	return b.String()
}

// Draft is one model-authored decision draft for a target.
type Draft struct {
	Target  Target
	Content string
}

// DraftAll drafts every target with the bounded completer, one call per target.
func DraftAll(ctx context.Context, completer Completer, db string, doc *cdl.EIRDoc) ([]Draft, error) {
	targets, err := Targets(db, doc)
	if err != nil {
		return nil, err
	}
	var drafts []Draft
	for _, target := range targets {
		content, err := completer.Complete(ctx, draftSystemPrompt, draftUserPrompt(target))
		if err != nil {
			return nil, fmt.Errorf("decision: draft %s: %w", target.SubjectID, err)
		}
		drafts = append(drafts, Draft{Target: target, Content: strings.TrimSpace(content)})
	}
	return drafts, nil
}

// DraftContent extracts the six sections from the model output, failing closed
// when a section is missing.
func DraftContent(content string) (map[string]string, error) {
	sections := map[string]string{}
	current := ""
	var body []string
	flush := func() {
		if current != "" {
			sections[current] = strings.TrimSpace(strings.Join(body, "\n"))
		}
		body = nil
	}
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			flush()
			current = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
			continue
		}
		if current != "" {
			body = append(body, line)
		}
	}
	flush()
	for _, want := range []string{"Options", "Decision", "Rationale", "Risks", "Exact Scope and Denominator Effect", "Affected Records / Rollout / Rollback"} {
		if strings.TrimSpace(sections[want]) == "" {
			return nil, fmt.Errorf("decision draft is missing the %q section", want)
		}
	}
	return sections, nil
}

func splitTokens(raw string) []string {
	var tokens []string
	for _, token := range strings.Split(raw, ",") {
		for _, piece := range strings.Split(token, ";") {
			piece = strings.TrimSpace(piece)
			if piece == "" || strings.EqualFold(piece, "none") {
				continue
			}
			tokens = append(tokens, strings.Fields(piece)[0])
		}
	}
	return tokens
}
