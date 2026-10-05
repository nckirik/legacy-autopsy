package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

const executorSystemPrompt = `You are one bounded semantic executor for the Canonical Deconstruction Protocol.

Rules:
- Use only the provided context. Never invent records, evidence, IDs, or scope.
- Produce exactly one result for the declared result contract.
- Output no prose outside the response format. The payload may contain ONLY the
  declared result-contract fields: never runtime metadata (mode, scope, target record
  ID, result contract, effects), fingerprints, versions, envelopes, or status tokens.
- Repository scans and directives are context, not payload fields; use them to choose
  values, do not echo their labels into the payload.
- If the context is insufficient, output a BLOCKER with the most accurate kind:
  uncertainty, protocol-mechanics, harness, human-authority, decision-approval,
  scope-approval, confirmation, domain-question.
- Never ask the operator directly; the runtime routes blockers.

Response format for a proposal (record identity and version are runtime-owned; do not emit them).
Emit one ---RECORD--- block per proposed record; a single-record unit emits exactly one:
RESULT: PROPOSAL
EFFECTS: none
---RECORD---
- **<Field Label>:** <value>
---END---

Response format for a blocker:
RESULT: BLOCKER
KIND: <kind>
MESSAGE: <one-sentence reason>
`

const reviewerSystemPrompt = `You are an independent reviewer for the Canonical Deconstruction Protocol.

Rules:
- Judge only the proposed record against the declared result contract, its declared
  dependencies, and the protocol semantics visible in the provided context.
- Do not assume the worker's reasoning; you did not see it.
- Output no prose outside the response format.

Response format:
OUTCOME: approved | rejected | inconclusive
FINDINGS:
- <one finding per line, or "none">
`

// Executor adapts a chat client to the bounded work executor interface. When
// Debug is set, every raw model response is written there for inspection; it
// must never be a committed artifact.
type Executor struct {
	client *Client
	Debug  io.Writer
	// ForceReasoning, when non-nil, overrides every unit's reasoning effort
	// (empty string omits the side parameter). Set by the CLI --llm-reasoning.
	ForceReasoning *string
}

// NewExecutor wraps a client as a work executor.
func NewExecutor(client *Client) *Executor { return &Executor{client: client} }

// SetDebug mirrors the raw-response diagnostic onto the client as well.
func (e *Executor) SetDebug(w io.Writer) {
	e.Debug = w
	e.client.Debug = w
}

// NewExecutorFromEnv builds an executor from environment configuration.
func NewExecutorFromEnv() (*Executor, error) {
	config, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return NewExecutor(client), nil
}

// Execute runs one work unit through the model and parses the bounded response.
func (e *Executor) Execute(invocation work.Invocation) (work.Proposal, error) {
	text, err := e.client.CompleteEffort(context.Background(), executorSystemPrompt, executorUserPrompt(invocation), e.effortFor(invocation.Unit))
	if e.Debug != nil {
		fmt.Fprintf(e.Debug, "\n===== model response for %s =====\n%s\n===== end %s =====\n", invocation.Unit.ID, text, invocation.Unit.ID)
	}
	if err != nil {
		return work.Proposal{}, err
	}
	proposal, err := ParseProposal(text)
	if err != nil {
		return work.Proposal{}, err
	}
	if proposal.FieldID == "" {
		proposal.FieldID = invocation.Unit.ResultType
	}
	if proposal.FieldID != invocation.Unit.ResultType {
		return work.Proposal{}, fmt.Errorf("llm: proposal writes %s, result contract is %s", proposal.FieldID, invocation.Unit.ResultType)
	}
	return proposal, nil
}

func (e *Executor) effortFor(unit work.Unit) string {
	if e.ForceReasoning != nil {
		return normalizeEffort(*e.ForceReasoning)
	}
	if unit.ReasoningEffort != "" {
		return normalizeEffort(unit.ReasoningEffort)
	}
	return e.client.config.Reasoning
}

func normalizeEffort(effort string) string {
	switch strings.ToLower(strings.TrimSpace(effort)) {
	case "", "off", "none", "disabled":
		return ""
	default:
		return strings.TrimSpace(effort)
	}
}

func executorUserPrompt(invocation work.Invocation) string {
	var b strings.Builder
	b.WriteString("Work unit metadata (not payload fields):\n")
	fmt.Fprintf(&b, "Mode: %s\n", invocation.Unit.Mode)
	fmt.Fprintf(&b, "Scope: %s\n", invocation.Unit.Scope)
	fmt.Fprintf(&b, "Result contract: %s\n", invocation.Unit.ResultType)
	if len(invocation.Unit.AllowedEffects) > 0 {
		fmt.Fprintf(&b, "Allowed effects: %s\n", strings.Join(invocation.Unit.AllowedEffects, ", "))
	}
	required := map[string]bool{}
	for _, name := range invocation.RequiredFields {
		required[name] = true
	}
	fmt.Fprintf(&b, "\nDeclared payload fields, in order. Emit exactly these fields and nothing else, as '- **<Human Label>:** <value>':\n")
	for _, name := range invocation.Fields {
		mark := "optional"
		if required[name] {
			mark = "required"
		}
		fmt.Fprintf(&b, "- %s (%s)\n", name, mark)
	}
	if len(invocation.Hints) > 0 {
		fmt.Fprintf(&b, "\nDeterministic validation (the proposal is rejected if violated):\n")
		for _, hint := range invocation.Hints {
			fmt.Fprintf(&b, "- %s\n", hint)
		}
	}
	fmt.Fprintf(&b, "\n%s", invocation.Context.Bytes)
	return b.String()
}

// ParseProposal parses the strict executor response format.
func ParseProposal(text string) (work.Proposal, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	fields := map[string]string{}
	var records []string
	var current []string
	inRecord := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---RECORD---" {
			if inRecord {
				return work.Proposal{}, fmt.Errorf("llm: nested ---RECORD--- block")
			}
			inRecord = true
			current = nil
			continue
		}
		if inRecord {
			if trimmed == "---END---" {
				inRecord = false
				if len(current) == 0 {
					return work.Proposal{}, fmt.Errorf("llm: empty ---RECORD--- block")
				}
				records = append(records, strings.Join(current, "\n")+"\n")
				continue
			}
			current = append(current, line)
			continue
		}
		if key, value, ok := strings.Cut(trimmed, ":"); ok {
			fields[strings.ToUpper(strings.TrimSpace(key))] = strings.TrimSpace(value)
		}
	}
	switch fields["RESULT"] {
	case "BLOCKER":
		kind := strings.ToLower(fields["KIND"])
		if kind == "" {
			return work.Proposal{}, fmt.Errorf("llm: blocker response has no KIND")
		}
		return work.Proposal{}, &work.BlockerError{Blocker: work.Blocker{Kind: kind, Message: fields["MESSAGE"]}}
	case "PROPOSAL":
	default:
		return work.Proposal{}, fmt.Errorf("llm: response has no RESULT marker")
	}
	if inRecord {
		return work.Proposal{}, fmt.Errorf("llm: unterminated ---RECORD--- block")
	}
	if len(records) == 0 {
		return work.Proposal{}, fmt.Errorf("llm: proposal response has no record block")
	}
	var effects []string
	if raw := fields["EFFECTS"]; raw != "" && !strings.EqualFold(raw, "none") {
		for _, effect := range strings.Split(raw, ",") {
			effect = strings.TrimSpace(effect)
			if effect != "" {
				effects = append(effects, effect)
			}
		}
	}
	proposal := work.Proposal{FieldID: fields["FIELD"], Effects: effects}
	for _, record := range records {
		proposal.Records = append(proposal.Records, work.Record{Payload: []byte(record)})
	}
	return proposal, nil
}

// ParseReview parses the strict reviewer response format.
func ParseReview(text string) (work.ReviewResult, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	outcome := ""
	var findings []string
	inFindings := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "OUTCOME:") {
			outcome = strings.ToLower(strings.TrimSpace(trimmed[len("OUTCOME:"):]))
			continue
		}
		if strings.EqualFold(trimmed, "FINDINGS:") {
			inFindings = true
			continue
		}
		if inFindings && strings.HasPrefix(trimmed, "- ") {
			finding := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if finding != "" && !strings.EqualFold(finding, "none") {
				findings = append(findings, finding)
			}
		}
	}
	switch work.ReviewOutcome(outcome) {
	case work.ReviewApproved, work.ReviewRejected, work.ReviewInconclusive:
		return work.ReviewResult{Outcome: work.ReviewOutcome(outcome), Findings: findings}, nil
	default:
		return work.ReviewResult{}, fmt.Errorf("llm: review response has no valid OUTCOME")
	}
}
