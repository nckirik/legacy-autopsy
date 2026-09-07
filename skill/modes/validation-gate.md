# Validation/Gate

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Run one named deterministic check or gate and emit only its exact report target.

## Required identity bindings
Bind one gate/check/stage, authoritative input fingerprint set, pinned snapshot, named report target, and explicit forbidden semantic targets.

## Normative protocol sections
- §4.1. `ID generation`; §4.1.1. `Semantic payload identity and certification envelopes`; §4.1.2. `Normative canonical hash profile and packaging artifacts`; §5.1. `Dedicated prefix groups`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §10.1. `Evidence-backed `[R-SWEPT]``; §10.2. `Per-kind coverage arithmetic`; §10.3. `Composite Exit A — Deconstruction Closure`
- §12.1. `Traceability — `20-TRACEABILITY.md``; §12.2. `Required reciprocal references`; §12.5. `Deterministic validation summary`; §13.4. `Handbook quality gate`
- §15.1. `Strict conditions`; §15.2. `Scope certificate`; §15.3. `Package contents`; §15.4. `Acyclic certification sequence and hash domains`; §19.2. `Normative conformance fixtures`

## Mandatory read set
Load every authoritative input for the named check: Exit A denominator/frontier/ticket/acquisition/coverage data, or stage-specific candidate/content-readiness/package/outer-manifest records and fingerprints.

## Semantic write targets
Only the named report block/artifact. Stage 2 owns candidate reports/manifests; stage 4 owns content-readiness; stage 6 validates the outer manifest in process.

## Allowed assurance/audit side effects
Deterministic `21`/`22` or exact stage reports with stable check IDs, input/output fingerprints, and append-only `0G`.

## Forbidden mutations
No source, forensic, synthesis, handbook, decision, or confirmation semantics; no hand-edited arithmetic, unsupported pass, pre-step-6 Passed state, or post-step-6 completion report.

## Deterministic validation
Run applicable structural Markdown, enum, ID, hash/envelope, reference, arithmetic, schema/order, handbook, package-member/signature, §15.1, and conformance-fixture checks; report unsupported checks explicitly.

## Completion condition
The named report deterministically records result, offending IDs, fingerprints, validator version, and stage state. Exit E passes only after in-process validation of the completed signed step-6 outer manifest.
