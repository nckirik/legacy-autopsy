# Promotion Review

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Evaluate promotion locks for one bound component or cohort and request—not execute—an atomic promotion.

## Required identity bindings

Bind exactly one persona, prefix, canonical persona directory, concrete cluster, track, POV, POV file, optional ledger, exact CMP/cohort, source shared path, and target persona/POV.

## Normative protocol sections

- §6.3. `Staging and promotion`; §6.4. `Depromotion`
- §8.1. `Resume identity header`; §8.2. `Strict single-scope modes`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §12.5. `Deterministic validation summary`

## Mandatory read set

Load §8.4 strict-scope inputs, canonical shared block, target POV, complete source/file fingerprints, coverage, tickets, contradictions, dependencies, every-track/snapshot usage, and persona siblings needed for uniqueness.

## Semantic write targets

One `PROMOTION-REQUEST` in the bound persona’s `SHARED-STAGING-BUFFER.md`; no canonical move.

## Allowed assurance/audit side effects

Fixed review evidence/disposition and append-only `0G` only.

## Forbidden mutations

No canonical block move/rewrite, promoted status, partial cohort movement, or direct `_shared`/target synchronization.

## Deterministic validation

Validate all five §6.3 locks, complete current fingerprints, singular combined persona set for cohorts, strict scope, and stale guard.

## Completion condition

Append a fully bound request when every lock passes, or record a deterministic failed/stale disposition; execution remains solely Sequential Reconciliation.
