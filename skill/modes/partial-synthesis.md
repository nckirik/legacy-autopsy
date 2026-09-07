# Partial-Synthesis

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Produce a bounded provisional synthesis before Exit A while preserving all gaps and denominator limits.

## Required identity bindings
Bind an explicit persona/cluster/track/snapshot scope and enumerate exact affected `90`–`96` blocks and synthesis-gap buffer targets.

## Normative protocol sections
- §4.1.1. `Semantic payload identity and certification envelopes`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §11.1. `Partial versus Final Synthesis`; §11.2. `Common synthesis block header`
- §11.3. `Architecture blueprint — `90-ARCH-BLUEPRINT.md`` through §11.12. `NFR and security — `96-NON-FUNCTIONAL-SECURITY.md``
- §12.4. `Dependency and impact provenance`; §12.5. `Deterministic validation summary`

## Mandatory read set
Load the §8.4 common set, all bounded forensic/assurance inputs, dependency evidence, and existing affected synthesis blocks.

## Semantic write targets
Only bounded provisional `90`–`96` blocks with `[B-DRAFT]` envelopes and `Pending` synthesis GAP records.

## Allowed assurance/audit side effects
Semantic versions/fingerprints, derivation links, targeted stale records, and append-only `0G`.

## Forbidden mutations
No reconstruction-input status, `[B-CONFIRMED]`, ticket/frontier/ledger writes, direct gap reconciliation, silent legacy facts, or claims beyond the bounded denominator.

## Deterministic validation
Validate legal IDs, domain schemas, exact `DERIVED-FROM` bindings, material claims/evidence profiles, incomplete-denominator labels, hashes, and envelope boundaries.

## Completion condition
The bounded provisional projection and every discovered gap are recorded without claiming closure or reconstruction readiness.
