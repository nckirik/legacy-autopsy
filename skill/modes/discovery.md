# Discovery

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Perform one bounded, persona-driven forensic traversal without crossing the bound root, persona, track, or POV.

## Required identity bindings

Bind exactly one persona, prefix, canonical persona directory, concrete cluster, track, POV, POV file, environment/snapshot, root, and depth; bind zero or one ledger and enumerate the fixed sidecar bundle.

## Normative protocol sections

- §2.3. `Five isolated POVs`; §2.4. `Traversal tracks`; §2.5. `Capability and runtime state`
- §4.2. `Atomic unit rule`; §4.3. `Source inventory denominator — `10-SOURCE-INVENTORY.md``; §4.4. `Traversal frontier — `11-TRAVERSAL-FRONTIER.md``; §4.5. `Scope-specific source coverage — `16-SOURCE-COVERAGE.md``
- §5.2. `Claim-level evidence — `12-CLAIM-EVIDENCE.md``; §6.1. `Atomic component schema`; §6.2. `Dual-entry discovery`
- §8.1. `Resume identity header`; §8.2. `Strict single-scope modes`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §10.1. `Evidence-backed `[R-SWEPT]``; §12.5. `Deterministic validation summary`

## Mandatory read set

Load §8.4 common and strict-scope sets: bound POV/ledger, matching shared records, persona unmapped buffer, concrete source, normalized map, conditional `0B`/`0C`, relevant dependencies, and current proof that no `Pending-Merge` unmapped entry exists.

## Semantic write targets

One bound POV and optional ledger; persona invocation pointers and local staging/unmapped/question buffers required by dual-entry discovery.

## Allowed assurance/audit side effects

Directly produced `10`, `11`, `12`, `14`, `16`, persona-local buffers, and append-only `0G`.

## Forbidden mutations

No root/persona/track/POV switch, other POV semantics, direct `_shared` merge, direct global `0A` unmapped write, traversal through unmapped boundaries, silent depth stop, or invented synthesis IDs.

## Deterministic validation

Validate strict identity isolation, stale guard, atomic schemas, evidence levels, exactly one FRT per boundary, tuple uniqueness, snapshot/track separation, transactional buffers, and §10.1 before `[R-SWEPT]`.

## Completion condition

The bound root is processed to depth; reached units/claims are recorded or ticketed, every boundary has an FRT, all local side effects commit atomically, and remaining work is represented truthfully.
