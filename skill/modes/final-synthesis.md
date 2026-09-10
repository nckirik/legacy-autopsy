# Final-Synthesis

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Generate complete closed-corpus synthesis only from a matching passed Composite Exit A input set.

## Required identity bindings

Bind the closed corpus, pinned snapshot and Exit A report/fingerprint, and enumerate exact `90`–`96` blocks and synthesis-gap buffer targets.

## Normative protocol sections

- §4.1.1. `Semantic payload identity and certification envelopes`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §10.3. `Composite Exit A — Deconstruction Closure`
- §11.1. `Partial versus Final Synthesis`; §11.2. `Common synthesis block header`; §11.3. `Architecture blueprint — `90-ARCH-BLUEPRINT.md`` through §11.12. `NFR and security — `96-NON-FUNCTIONAL-SECURITY.md``
- §12.1. `Traceability — `20-TRACEABILITY.md``; §12.2. `Required reciprocal references`; §12.3. `Synthesis-ID timing`; §12.4. `Dependency and impact provenance`; §12.5. `Deterministic validation summary`

## Mandatory read set

Load all active personas and POV/ledger files, relevant claims, acquisition/reconciliation/coverage/contradictions, gate inputs, existing synthesis blocks, and the matching Exit A proof.

## Semantic write targets

Complete closed-corpus `90`–`96` payloads and `Pending` synthesis GAP records.

## Allowed assurance/audit side effects

Semantic versions/fingerprints, derivation graph, targeted staleness, and append-only `0G`.

## Forbidden mutations

No run before matching Exit A, direct ledger/ticket/frontier mutation, confirmation, silent gap suppression, or use of stale Exit A inputs.

## Deterministic validation

Validate Exit A/snapshot match, mandatory and zero domains, schemas, claims/CAP/disabled rules, reciprocal prerequisites, versions/hashes, and stale closure.

## Completion condition

All domains are truthfully populated and fingerprinted, normally `[B-POPULATED]`, or every new gap is buffered; a material gap blocks handoff and re-arms deconstruction through reconciliation.
