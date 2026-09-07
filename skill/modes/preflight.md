# Preflight

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Establish one structurally valid, resumable `0A-PREFLIGHT.md` scope without claiming discovery or closure.

## Required identity bindings
Bind Protocol v4.0, namespace, iteration, invocation ID, `Invocation Mode: Preflight`, snapshot/environment, exact scope, sole semantic target `0A`, and explicit read/write/forbidden targets. Persona, cluster, track, and POV may be `None` only where §8.1 permits.

## Normative protocol sections
- §2.1. `Persona definition and purity`
- §2.2. `Canonical entry ownership`
- §2.4. `Traversal tracks`
- §3.1. `` `0A-PREFLIGHT.md` ``
- §4.1. `ID generation`
- §5.1. `Dedicated prefix groups`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`
- §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §12.5. `Deterministic validation summary`

## Mandatory read set
Load the §8.4 common set, active snapshot metadata, existing `0A` registries, evidence roots, persona/ownership context, and inventory evidence needed for declared boundaries.

## Semantic write targets
Only `0A-PREFLIGHT.md`: namespace, boundaries, iteration, personas/prefixes, clusters, tracks, traversal matrix, and acquisition-register scaffolding.

## Allowed assurance/audit side effects
Contractually required fixed assurance rows and append-only `0G` history.

## Forbidden mutations
No POV, synthesis, handbook, profile, decision, confirmation, or unrelated registry semantics; no fabricated facts or closure states.

## Deterministic validation
Validate version/mode, namespace and IDs, prefix uniqueness/non-reuse, `_shared` exclusion, canonical entry ownership, concrete track/matrix structure, and absence of wildcard-only closure.

## Completion condition
The declared preflight scope and sole `0A` target are valid and resumable; no discovery, coverage, Exit A, or confirmation is implied.
