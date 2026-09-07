# Sequential Reconciliation

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Deterministically merge concurrent buffers and execute authorized canonical promotion/depromotion transactions.

## Required identity bindings
Use a multi-file header naming every source buffer, destination file/section, transaction scope, snapshot, and exact read/write/forbidden target.

## Normative protocol sections
- §3.1. `` `0A-PREFLIGHT.md` ``; §3.2. `Foundational registries`
- §6.3. `Staging and promotion`; §6.4. `Depromotion`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.6. `Concurrent persona traversal`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §9.1. `Ticket schema and canonical IDs`; §11.1. `Partial versus Final Synthesis`; §12.1. `Traceability — `20-TRACEABILITY.md``; §12.5. `Deterministic validation summary`

## Mandatory read set
Load all relevant canonical `_shared` files, local staging/promotion/question/unmapped buffers, synthesis-gap buffers, destinations, ticket allocation state, promotion locks, indexes, traceability, and frontiers.

## Semantic write targets
Enumerated canonical merges; atomic promotion/depromotion; buffered tickets/frontiers; global `0A` unmapped merges; and buffer consumption/dispositions.

## Allowed assurance/audit side effects
Owned index, `10`–`16`, contradiction, traceability, frontier, stale-propagation, and append-only `0G` updates.

## Forbidden mutations
No nondeterministic ordering, confidence inflation, provenance/persona loss, partial move, duplicate/missing canonical copy, silent incompatible merge, or unrelated semantic invention.

## Deterministic validation
Sort/deduplicate by protocol keys; allocate tickets deterministically; revalidate locks/fingerprints; enforce atomic moves, reciprocal/index/coverage checks, targeted staleness, and rollback on failure.

## Completion condition
Every admitted item is deterministically merged or dispositioned, consumed items are atomically purged, destinations/sidecars agree, no partial transaction remains, and pending unmapped entries are cleared before Discovery.
