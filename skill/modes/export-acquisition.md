# Export Acquisition

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Register, sanitize, explode, and reconcile one approved serialized-export acquisition transaction.

## Required identity bindings

Bind Protocol v4.1, namespace, iteration/invocation, `Invocation Mode: Export Acquisition`, pinned snapshot/environment, exact approved artifact/projection set, and every multi-file read/write/forbidden target.

## Normative protocol sections

- §4.2. `Atomic unit rule`; §4.3. `Source inventory denominator — `10-SOURCE-INVENTORY.md``
- §5.3. `Sanitized projections and helper limitations`; §5.6. `Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md``
- §7.1. `Trust boundary`; §7.2. `Acquisition register`; §7.3. `Logical explosion and virtual coordinates`; §7.4. `Required explosion coverage`; §7.5. `n8n human acquisition loop`; §7.6. `Normalized maps`; §7.7. `Export reconciliation — `13-EXPORT-RECONCILIATION.md``
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §12.5. `Deterministic validation summary`

## Mandatory read set

Load §8.4 acquisition inputs: `0A`, `0C`, `0G`/valid `0H`, approved projections, lineage/sanitization reports, normalized maps, source inventory, export reconciliation, and snapshot metadata.

## Semantic write targets

Enumerated acquisition-register sections, normalized maps, acquisition-owned inventory/reconciliation/candidate records, approved placeholders/fragments, and security findings.

## Allowed assurance/audit side effects

Directly produced `10`, `13`, `17`, applicable `0C`, and append-only `0G` entries.

## Forbidden mutations

No live production acquisition, credentials, production mutation/execution, raw secrets/private exports/hints, correlation-as-evidence, behavior synthesis, or mixed snapshots.

## Deterministic validation

Validate helper schema/version, lineage, sanitization and secret scan, fingerprints, complete logical explosion, one reconciliation row per coordinate, terminal candidate/action dispositions, ordering, and snapshot consistency.

## Completion condition

The bounded transaction is atomically registered and truthfully dispositioned; `[A-STRUCTURALLY-COMPLETE]` requires static reference closure, 100% normalized reconciliation, and terminal candidate/action/reference state.
