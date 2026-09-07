# Architecture

> Non-authoritative implementation guide. `protocol.md` is normative and resolves every conflict.

## M0 boundary

M0 provides a Go CLI, thin skill routing, structural Markdown foundations, initial identities/path handling, workspace scaffolding/checks, confined context assembly, bootstrap fixtures, and CI. It does not execute semantic invocation modes or implement complete cold resume, protocol schemas, canonical hashing, gates, acquisition, synthesis, confirmations, handbook generation, or packaging.

## Target boundaries and flow

```text
protocol.md
    ↓
thin agent skill / invocation router
    ↓
context builder
    ↓
harness
    ↓
deterministic core + validators
    ↓
.extracted/ workspace

operator-controlled export acquisition ──approved sanitized projections──▶ workspace
```

- **Normative protocol:** execution law, records, authority, modes, invariants, and gates; never generated from implementation metadata.
- **Thin router:** selects one official mode and points to exact authoritative sections without replacing them.
- **Context builder:** assembles the complete mandatory read set and identity needed for a bounded invocation.
- **Harness:** will enforce one invocation, write scope, stale checks, atomic effects, audit output, and stop conditions.
- **Deterministic core:** owns mechanics unsuitable for model judgment; semantic investigation remains with the LLM.
- **Workspace:** `.extracted/` is persistent state; conversation memory is disposable and generated sidecars are derivative.
- **Acquisition boundary:** a future independent operator tool keeps raw Appsmith/n8n material private and publishes only approved sanitized projections.

## Language decision: Go

The M0 harness is implemented in Go for a cross-platform distributable CLI, strong standard-library filesystem/hashing/serialization support, explicit fail-closed errors, and no external runtime dependency after compilation. Future build-time modules must be justified and pinned behind project-owned, provider-neutral interfaces.

## Current packages

- `internal/markdown`: limited structural nodes and source spans used by M0 readers.
- `internal/protocol`: protocol version, heading, mode registry, and exact section extraction.
- `internal/routing`: manifest/projection drift checks.
- `internal/identity`: foundational typed IDs and relative-path normalization; specialized identities remain unsupported.
- `internal/canonical`: deliberately limited basic Markdown fingerprinting, not the §4.1.2 canonical profile.
- `internal/workspace`: atomic skeleton initialization and early structural checks.
- `internal/contextpacket`: identity validation, confined reads, exact routed text, and packet rendering; automatic closure/stale comparison remain deferred.
- `internal/fixtures`: registered bootstrap cases and stable diagnostic matching.
- `internal/cli`: command wiring.

Planned packages or responsibilities—complete schemas, canonicalization, validation, acquisition adapters, semantic records, gates, and packaging—remain roadmap work. CLI commands should continue to orchestrate provider-neutral packages rather than contain protocol logic.

## Structural Markdown direction

The current parser recognizes headings, emphasized fields, fenced blocks, comments, tables, paragraphs, blanks, and bounded sections. Later milestones must deepen it for typed-record boundaries, complete field paths, table schemas, certification/artifact envelopes, canonical serialization, and stable anchors. Protocol-significant validation must remain structural; regex is appropriate only after locating a protocol-defined regex-constrained value.

## Safety and persistence

Reads and writes use declared roots and path normalization; M0 context reads also reject traversal and symlink escape. Future mutations require restrictive permissions where sensitive, atomic publication, and stale-input checks. Persistent authority comes from validated filesystem records, never chat history.
