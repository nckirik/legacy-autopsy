package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

const executorSystemPrompt = `You are one bounded semantic executor for the Canonical Deconstruction Protocol.

Rules:
- Use only the provided context. Never invent records, evidence, IDs, or scope.
- Produce exactly one result for the declared result contract.
- Output no prose outside the response format. No fingerprints, versions, envelopes,
  or status metadata inside the payload; the runtime owns all mechanical fields.
- If the context is insufficient, output a BLOCKER with the most accurate kind:
  uncertainty, protocol-mechanics, harness, human-authority, decision-approval,
  scope-approval, confirmation, domain-question.
- Never ask the operator directly; the runtime routes blockers.

Response format for a proposal (record identity and version are runtime-owned; do not emit them):
RESULT: PROPOSAL
EFFECTS: none
---PAYLOAD---
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

// Executor adapts a chat client to the bounded work executor interface.
type Executor struct {
	client *Client
}

// NewExecutor wraps a client as a work executor.
func NewExecutor(client *Client) *Executor { return &Executor{client: client} }

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
	text, err := e.client.Complete(context.Background(), executorSystemPrompt, executorUserPrompt(invocation))
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

func executorUserPrompt(invocation work.Invocation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "- **Mode:** %s\n", invocation.Unit.Mode)
	fmt.Fprintf(&b, "- **Scope:** %s\n", invocation.Unit.Scope)
	fmt.Fprintf(&b, "- **Result contract:** %s\n", invocation.Unit.ResultType)
	if len(invocation.Unit.AllowedEffects) > 0 {
		fmt.Fprintf(&b, "- **Allowed effects:** %s\n", strings.Join(invocation.Unit.AllowedEffects, ", "))
	}
	fmt.Fprintf(&b, "\n%s", invocation.Context.Bytes)
	return b.String()
}

// ParseProposal parses the strict executor response format.
func ParseProposal(text string) (work.Proposal, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	fields := map[string]string{}
	var payload []string
	inPayload := false
	payloadStarted := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---PAYLOAD---" {
			inPayload = true
			payloadStarted = true
			continue
		}
		if inPayload {
			if trimmed == "---END---" {
				inPayload = false
				continue
			}
			payload = append(payload, line)
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
	if !payloadStarted {
		return work.Proposal{}, fmt.Errorf("llm: proposal response has no payload block")
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
	return work.Proposal{
		FieldID: fields["FIELD"],
		Payload: []byte(strings.Join(payload, "\n") + "\n"),
		Effects: effects,
	}, nil
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
