# Cross-Reference-Reconciliation

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Reconcile reference fields among existing records without creating new legacy facts.

## Required identity bindings

Bind exact existing record IDs, fields, reciprocal destinations, traceability/index sections, snapshot, and every enumerated write target.

## Normative protocol sections

- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §11.1. `Partial versus Final Synthesis`
- §12.1. `Traceability — `20-TRACEABILITY.md``; §12.2. `Required reciprocal references`; §12.3. `Synthesis-ID timing`; §12.4. `Dependency and impact provenance`; §12.5. `Deterministic validation summary`

## Mandatory read set

Load all affected existing IDs and placeholders, reciprocal records, `20-TRACEABILITY.md`, dependency graph, current semantic versions/fingerprints, and certification envelopes.

## Semantic write targets

Only existing ID-reference fields and their reciprocal links.

## Allowed assurance/audit side effects

Required `0E`, `20`, targeted stale propagation, and append-only `0G` updates.

## Forbidden mutations

No new legacy facts, invented IDs, changed observed behavior/decision/evidence, unresolved dangling links, or unchanged semantic version after changing a semantic reference binding.

## Deterministic validation

Validate type, existence, reciprocity, orphan/dangling/duplicate-owner/tombstone rules, deterministic version/fingerprint increment, envelope invalidation, and reverse-dependency staleness.

## Completion condition

Every declared target is reconciled or explicitly blocking; reciprocal/traceability checks agree, and changed payloads are re-versioned and re-fingerprinted.
