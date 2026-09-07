# Human Hatch

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Record exactly one authorized human action without borrowing rights from another Human Hatch subtype.

## Required identity bindings
Bind `Invocation Mode: Human Hatch`, exactly one non-`None` action (`Ticket Escalation`, `Decision/Scope Approval`, or `Confirmation`), reviewer identity/authority, exact scope, affected records, snapshot, and fingerprints.

## Normative protocol sections
- §5.5. `Authoritative decisions — `15-DECISIONS.md``; §5.7. `Human confirmations — `18-CONFIRMATIONS.md``
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §9.1. `Ticket schema and canonical IDs`; §9.2. `FSM and write rights`; §9.3. `Human Hatch action prerequisites`; §9.4. `Probe specification and asynchronous handoff`
- §10.2. `Per-kind coverage arithmetic`; §12.5. `Deterministic validation summary`

## Mandatory read set
Load subtype prerequisites: escalation ticket/static/probe history; decision inputs, exact COV arithmetic, alternatives/risks; or passing candidate/manifest, semantic and DEC bindings, dependency closure, limitations, and authority for confirmation.

## Semantic write targets
Subtype-only: escalation terminal ticket/effects; decision content or approval envelope; or CNF plus exact certification envelopes.

## Allowed assurance/audit side effects
Only exact authorized coverage, traceability, risk effects, and append-only `0G`.

## Forbidden mutations
No cross-subtype rights, direct human-required creation without prerequisites, terminal-ticket-as-coverage inference, semantic PRF/synthesis/HBK edits during confirmation, or mismatched candidate approval/confirmation.

## Deterministic validation
Validate subtype/action and authority, §9.3 prerequisites, DEC/CNF hashes and envelope bindings, denominator arithmetic, candidate/closure identity, signatures/timestamps, and no semantic-byte change for envelope-only actions.

## Completion condition
One authorized action and all exact effects commit atomically; escalation alone preserves the gap unless an exact approved exclusion exists, and confirmation changes only envelopes/CNF state.
