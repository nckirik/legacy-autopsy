# Profile Synchronization

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose

Synchronize one persona profile semantic payload with its matching `0A` persona row atomically.

## Required identity bindings

Bind exactly one persona, unique prefix, deterministic PRF ID, target profile, matching `0A` row, snapshot, and explicit read/write/forbidden targets.

## Normative protocol sections

- §4.1.1. `Semantic payload identity and certification envelopes`
- §8.1. `Resume identity header`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume check`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §11.13. `Persona profile`; §12.1. `Traceability — `20-TRACEABILITY.md``; §12.2. `Required reciprocal references`; §12.4. `Dependency and impact provenance`; §12.5. `Deterministic validation summary`

## Mandatory read set

Load `0A`, target PRF, every directly referenced/applicable record, and the complete transitive §8.4 `DERIVED-FROM` closure with versions, fingerprints, reciprocal bindings, assurance, decisions, contradictions, coverage, traceability, and snapshot proof.

## Semantic write targets

Exactly one profile semantic payload and its matching `0A` persona row in one transaction.

## Allowed assurance/audit side effects

Required reciprocal traceability/index bindings and append-only `0G`. A separate authorized confirmation action may later alter only the certification envelope.

## Forbidden mutations

No other profile; no absent/stale/mismatched/mixed-snapshot dependency; no semantic edit during confirmation; no duplicate PRF/prefix; no split profile/registry transaction.

## Deterministic validation

Validate deterministic PRF identity, one PRF per row, complete direct/transitive closure, versions/fingerprints/snapshot, reciprocity, semantic hash, and payload/envelope boundary.

## Completion condition

The single PRF and `0A` row agree atomically, all dependencies are current and pinned, semantic identity is valid, and profile status truthfully reflects its state.
