# Roadmap

> Non-authoritative planning projection. Milestone completion never overrides `protocol.md`, and it does not imply Protocol v4 conformance unless the required Part 19 fixtures pass.

Milestones are incremental and dependency-ordered completion criteria. M0 intentionally includes narrow bootstrap slices of M1 (structural parsing/protocol discovery), M2 (foundational IDs/paths), M3 (initialization/checks), and M5 (bounded context assembly); those later milestones are not complete. Each implemented normative rule requires deterministic tests, positive/negative fixtures, stable expected results, and cited protocol sections.

## M0 — Repository and skill foundation

Establish authority messaging, contributor rules, architecture/development/workspace/conformance docs, the thin skill and 13 mode projections, provider-neutral Go module/CLI foundation, routing validation, basic deterministic primitives, workspace scaffolding, bootstrap fixtures, and honest CI reporting. This milestone is implemented; later milestones extend protocol depth without changing the repository shape.

## M1 — Markdown AST + protocol model

Implement the small structural Markdown AST abstraction, source spans, heading/record and payload/envelope boundaries, field/table paths, protocol version discovery, official invocation-mode registry, and protocol-to-skill routing validation. No whole-document regex substitutes.

## M2 — Deterministic IDs/path normalization

Implement normalized relative paths, traversal rejection, exact Unicode/case preservation, canonical ID keys, SHA-256 prefix and deterministic collision extension, aliases/tombstones, and foundational identity fixtures. Include PRF/HBK coordinate foundations without claiming complete PRF/HBK conformance.

## M3 — Workspace initialization and structural validators

Implement `doctor`, non-fabricating `init`, workspace discovery, four-plane skeletons, persona/shared ownership boundaries, atomic file writes, `0G`/derived `0H` mechanics, and initial structural validation. Initialization must not create covered, confirmed, or passed semantic state.

## M4 — Claims/source/frontier primitives

Implement source inventory, concrete denominator rows, claim-level evidence, frontiers including immediately terminal boundaries, status/path-qualified enum dispatch, acquisition register primitives, atomic record constraints, and reproducible zero-domain proofs.

## M5 — Invocation context builder + cold resume

Implement all official mode routing, Protocol §8.1 identity headers, mandatory read sets, minimum-sufficient authoritative section packets, strict-scope verification, stale-checkpoint guards, §8.5 cold-resume summaries, and append-only invocation history. Conversation memory remains disposable.

## M6 — Ticket FSM and reconciliation

Implement ticket identities/transitions, probe ingestion boundaries, Human Hatch prerequisites, persona-local buffers, deterministic Sequential Reconciliation, promotion/depromotion atomicity, unmapped discovery and synthesis-gap merging, concurrency locks, and reciprocal audit effects.

## M7 — Coverage and Exit A

Implement per-kind/track/environment/snapshot coverage arithmetic, evidence-backed sweeps, source/frontier/ticket/export closure, ownership and referential integrity, disabled capability governance, contradiction checks, deterministic reports, and all Composite Exit A conditions. Claim support only after validators and fixtures pass.

## M8 — Export acquisition helper

Implement the independently runnable operator-controlled Appsmith/n8n helper: dry-run default, safe credential channels, private state outside the repository, field-aware sanitization, secret scanning, fail-closed schema/lineage review, normalized logical explosion, reconciliation candidates, and atomic approved publication. No production mutation or third-party transmission.

## M9 — Structured synthesis

Implement Partial-Synthesis and Final-Synthesis boundaries, common semantic headers, `DERIVED-FROM` graphs, semantic versions/fingerprints, gap buffering, evidence profiles, atomic catalogs `90`–`96`, and Cross-Reference-Reconciliation with targeted stale propagation.

## M10 — Confirmation/staleness

Implement candidate-bound `CNF` records, decision approvals, exact payload-versus-envelope separation, human-only confirmation transitions, dependency impact graphs, semantic staleness, and invalidation/re-entry behavior without mutating confirmed payload identity.

## M11 — Handbook

Implement HBK identities and stable anchors, synchronized handbook generation, progressive-disclosure navigation, quality/link/diagram checks, upstream-first correction, handbook staleness, and confirmation envelopes. The handbook remains a projection, not independent authority.

## M12 — Canonical packaging/hashing

Implement the normative Markdown canonical hash profile, independent typed-record hashes, file-transport fingerprints, exact carrier/envelope exclusions, canonical envelope bindings, complete package-member fingerprints, artifact hash domains, post-hash instance bindings, and deterministic schemas/order for all five packaging artifacts.

## M13 — Exit E

Implement content and decision readiness, equivalence-suite validation, exact confirmations, candidate report/manifest, content-readiness report, scope certificate, outer manifest, signature validation, and the strict acyclic six-step sequence. Only the validated signed outer payload may state `EXIT-E-STATUS: Passed`.

## M14 — Full Part 19 conformance suite

Complete every required deterministic positive and negative fixture family, stabilize diagnostics, run all groups in CI, and block applicable gates on failures. Perform a requirement-to-fixture audit against the current normative protocol before making any conformance claim.

## M15 — First real legacy-system dry run

Run a controlled, non-production pilot with pinned snapshots and approved sanitized evidence. Exercise acquisition, deconstruction, Exit A, synthesis, human confirmation, handbook, packaging, and Exit E; document gaps and feed fixes back through earlier milestones without weakening protocol rules.

## Explicitly deferred Protocol §19.2 families

All are currently deferred and remain roadmap-visible:

1. **Persona workspace layout and ownership** — M3/M6, including `_shared`, purity, buffers, and atomic Sequential Reconciliation moves.
2. **Context-qualified finite enums and registry completeness** — M1/M4, including bidirectional schema registry and traversal `N/A` coupling.
3. **PRF identity and semantic hashing** — M2/M9/M10, including `0A`, candidate, traceability, and CNF bindings.
4. **HBK identity and path/anchor normalization** — M2/M11, including the fixed handbook root, exact Unicode/path handling, stable anchors, title edits, and proven moves.
5. **Profile Synchronization closure** — M5/M9/M10, including complete direct/applicable and transitive dependency closure and identical confirmation closure.
6. **Record versus artifact hashing** — M1/M12, including independent record hashes, ordered file transport, artifact hash-domain identity, and post-hash binding.
7. **Canonicalization and exact exclusions** — M1/M10/M12, including LF/Unicode/order/whitespace handling and exact certification, DEC, CNF, and artifact boundaries.
8. **Final envelope and package-member integrity** — M10/M12/M13, including identity/version-bound envelope fingerprints and complete finalized-file fingerprints.
9. **Packaging schemas and deterministic ordering** — M12/M13, covering all five schemas, field order, empty values/tables, sort keys, counts, and cycle prevention.
10. **Acyclic Exit E sequence** — M13, from candidate report through the signed outer manifest, with `Pending` through steps 1–5 and no post-step-6 completion artifact.

## Other deferred protocol families

The roadmap also explicitly carries source/claim/frontier evidence (M4), stateless invocation and cold resume (M5), tickets/probes/Human Hatch and reconciliation (M6), coverage and Exit A (M7), export acquisition/reconciliation (M8), synthesis and gaps (M9), decisions/confirmations/staleness (M10), handbook/equivalence usability (M11), canonical packaging (M12), Exit E (M13), full acceptance fixtures (M14), and corpus acceptance in a real dry run (M15).

## Current next milestone

Proceed to M1: deepen the structural Markdown AST and protocol model with typed-record, field-path, table, payload, and envelope boundaries. Then complete M2 collision extension and identity coordinates before expanding M3 workspace ownership validation. Do not start gate or packaging claims before those structures can locate normative boundaries reliably.
