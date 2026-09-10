# Roadmap

> Non-authoritative planning projection. Milestone completion never overrides [`protocol.md`](../protocol.md), and it does not imply Protocol v4 conformance unless the required Part 19 fixtures pass.

Milestones are incremental and dependency-ordered completion criteria. M0 intentionally includes narrow bootstrap slices of M1 (structural parsing/protocol discovery), M2 (foundational IDs/paths), M3 (initialization/checks), and M6 (bounded context assembly); those later milestones are not complete. Each implemented normative rule requires deterministic tests, positive/negative fixtures, stable expected results, and cited protocol sections.

The application should become useful and observable while deeper protocol automation is still being implemented. The early runtime milestone therefore follows the foundational protocol, identity, and workspace mechanics; it does not claim semantic execution or gate support before their validators exist.

## M0 — Repository and skill foundation

Establish authority messaging, contributor rules, architecture/development/workspace/conformance docs, the thin skill and 13 mode projections, provider-neutral Go module/CLI foundation, routing validation, basic deterministic primitives, workspace scaffolding, bootstrap fixtures, and honest CI reporting. This milestone is implemented; later milestones deepen protocol behavior and introduce the service/workbench repository shape.

## M1 — Markdown AST + protocol model

Implement the small structural Markdown AST abstraction, source spans, heading/record and payload/envelope boundaries, field/table paths, protocol version discovery, official invocation-mode registry, and protocol-to-skill routing validation. No whole-document regex substitutes.

## M2 — Deterministic IDs/path normalization

Implement normalized relative paths, traversal rejection, exact Unicode/case preservation, canonical ID keys, SHA-256 prefix and deterministic collision extension, aliases/tombstones, and foundational identity fixtures. Include PRF/HBK coordinate foundations without claiming complete PRF/HBK conformance.

## M3 — Workspace initialization and structural validators

Implement `doctor`, non-fabricating `init`, workspace discovery, four-plane skeletons, persona/shared ownership boundaries, atomic file writes, `0G`/derived `0H` mechanics, and initial structural validation. Initialization must not create covered, confirmed, or passed semantic state.

## M4 — Local autopsy runtime foundation

Introduce the first observable application shell without claiming semantic mode execution:

- installable user-level OS service with a registry in the user's application-data directory and startup reconciliation;
- multi-autopsy registry with a persisted opaque random `autopsy_id`, repository/snapshot binding, rooted workspace isolation, explicit relocation reattach, and duplicate-root rejection;
- initial scheduler with one active invocation per autopsy, two globally, FIFO within an autopsy, round-robin across autopsies, and simple invocation-ID idempotency/cancellation/recovery;
- common invocation/result boundary, exposed first through the generic skill in any compatible coding harness;
- skill-driven project attachment and a reserved non-authoritative `.legacy-autopsy/` configuration boundary;
- explicit prohibition on UI/CLI manual semantic-work start or claim during this stage;
- workspace transaction/post-commit hook that records every started invocation in `0G` and publishes committed revisions without waiting for projections;
- optional `.extracted/` watcher for out-of-band change detection, revalidation, and startup/full-rescan recovery—not as authority for service-owned writes;
- lower-priority derivative workers whose Atlas lag/failure cannot consume semantic-worker capacity or block extraction;
- loopback-only local API with no account/client authentication or external serving, plus strict `Host`/`Origin`/CSRF controls and secret-safe logs/events;
- browser shell using Angular with Taiga UI, Analog's Vite plugin (not the full Analog framework), multi-autopsy tabs, lifecycle/status indicators, and activity stream;
- initial Cytoscape Atlas surface backed by lower-priority regenerable JSON under `.legacy-autopsy/atlas/*`, never participating in extraction closure;
- questions inbox shell that cannot bypass protocol prerequisites or establish truth.

At this milestone, the skill is the only semantic-work ingress; the UI and CLI may observe, configure, administer, and present authorized human actions but cannot manually start or claim work. Protocol-dependent operations without complete validators report `unsupported`. Empty scaffolds and operational metadata must not fabricate graph semantics, coverage, Human Hatch outcomes, Exit A, or conformance.

## M5 — Claims/source/frontier primitives

Implement source inventory, concrete denominator rows, claim-level evidence, frontiers including immediately terminal boundaries, status/path-qualified enum dispatch, acquisition register primitives, atomic record constraints, and reproducible zero-domain proofs. Feed only committed validated records into runtime projections.

## M6 — Invocation context, cold resume, and executor rollout

Implement all official mode routing, Protocol §8.1 identity headers, mandatory read sets, minimum-sufficient authoritative section packets, strict-scope verification, stale-checkpoint guards, §8.5 structured cold-resume checks, and append-only invocation history. The executor proposes semantic resume/actions; the runtime recomputes identity, scope, loaded fingerprints, write rights, staleness, and the cold-resume fingerprint before commit, including audited no-mutation outcomes. Integrate service claim/context/submit with fingerprint revalidation, one-invocation leases, idempotent submission, and crash-safe recovery. Conversation memory remains disposable.

Deliver execution channels in order:

1. complete the generic service-backed skill flow so it works from any compatible coding harness, with no UI/CLI manual semantic start;
2. add thin managed adapters for named harnesses such as OpenCode, Codex, Claude Code, and Kiro, backed by versioned per-project `.legacy-autopsy/` configuration and secret-free runtime metadata;
3. add the minimal internal direct-model runner only after the skill and adapter contracts are stable.

All three use the same invocation/result boundary. Adapter configuration is operational, not protocol authority; credentials remain outside the repository and `.extracted/`.

## M7 — Ticket FSM, reconciliation, and asynchronous human dependencies

Implement ticket identities/transitions, closed escalation reasons, runtime and non-runtime Human Hatch prerequisites, honest human-required gap/partial outcomes distinct from exclusion, probe ingestion boundaries, persona-local buffers, deterministic Sequential Reconciliation, promotion/depromotion atomicity, unmapped discovery and synthesis-gap merging, concurrency locks, and reciprocal audit effects. Back the questions inbox with validated reason/evidence records and dependency impact; continue independent frontiers while human input waits.

## M8 — Coverage and Exit A

Implement per-kind/track/environment/snapshot coverage arithmetic, evidence-backed sweeps, source/frontier/ticket/export closure, dynamic-caller-aware dead-code proof, provenance-bound lifecycle/applicability predicates with fail-inclusive `Unknown`, ownership and referential integrity, disabled capability governance, contradiction checks, deterministic reports, and all Composite Exit A conditions. Project validated status into the workbench; UI state cannot establish a gate. Claim support only after validators and fixtures pass.

## M9 — Export acquisition helper

Implement the independently runnable operator-controlled Appsmith/n8n helper: dry-run default, safe credential channels, private state outside the repository, field-aware sanitization, secret scanning, fail-closed schema/lineage review, normalized logical explosion, reconciliation candidates, and atomic approved publication. No production mutation or third-party transmission.

## M10 — Structured synthesis

Implement Partial-Synthesis and Final-Synthesis boundaries, common semantic headers, `DERIVED-FROM` graphs, semantic versions/fingerprints, gap buffering, evidence profiles, atomic catalogs `90`–`96`, and Cross-Reference-Reconciliation with targeted stale propagation. Expand Atlas nodes and edges only from committed records with inspectable provenance.

## M11 — Confirmation/staleness

Implement candidate-bound `CNF` records, decision approvals, exact payload-versus-envelope separation, human-only confirmation transitions, dependency impact graphs, semantic staleness, and invalidation/re-entry behavior without mutating confirmed payload identity.

## M12 — Handbook and Atlas linking

Implement HBK identities and stable anchors, synchronized handbook generation, progressive-disclosure navigation, deterministic non-vacuity/link/diagram checks, candidate-bound CNF readability review coverage, upstream-first correction, handbook staleness, confirmation envelopes, and stable bidirectional links between handbook sections and Atlas records. Both views remain projections over authoritative workspace state.

## M13 — Canonical packaging/hashing

Implement the normative Markdown canonical hash profile, independent typed-record hashes, file-transport fingerprints, exact carrier/envelope exclusions and legal prior-artifact envelope references, canonical envelope bindings, row evidence-set fingerprints, complete package-member fingerprints, artifact hash domains, post-hash instance bindings, and deterministic schemas/order for all five packaging artifacts.

## M14 — Exit E

Implement content and decision readiness, exact stage-specific check registries and evidence bindings, equivalence-suite validation, exact confirmations, candidate report/manifest, content-readiness report, scope certificate, directly snapshot-bound outer manifest, cross-chain snapshot equality, signature validation, reproducible final-bundle verification/non-authoritative receipts, and the strict acyclic six-step sequence. Only the validated signed outer payload may state `EXIT-E-STATUS: Passed`.

## M15 — Full Part 19 conformance suite

Complete every required deterministic positive and negative fixture family, stabilize diagnostics, run all groups in CI, and block applicable gates on failures. Perform a requirement-to-fixture audit against the current normative protocol before making any conformance claim.

## M16 — First real legacy-system dry run

Run a controlled, non-production pilot with pinned snapshots and approved sanitized evidence. Exercise the service/workbench and all three execution stages in order—generic skill, named harness adapters, then minimal direct-model runner—plus acquisition, deconstruction, asynchronous questions, Atlas, Exit A, synthesis, human confirmation, handbook, packaging, and Exit E; document gaps and feed fixes back through earlier milestones without weakening protocol rules.

## Explicitly deferred Protocol §19.2 families

All are currently deferred and remain roadmap-visible:

1. **Persona workspace layout and ownership** — M3/M7, including `_shared`, purity, buffers, and atomic Sequential Reconciliation moves.
2. **Context-qualified finite enums and registry completeness** — M1/M5, including bidirectional schema registry and traversal `N/A` coupling.
3. **Invocation ownership and cold resume** — M6, including semantic proposal versus runtime-computed authorization, stale checks, check fingerprints, and `0G` bindings.
4. **Ticket escalation and honest unresolved coverage** — M7/M8, including closed non-runtime reasons, Human Hatch rights, human-required gap/partial behavior, and exclusion separation.
5. **Dead-code and semantic-predicate closure** — M5/M8/M12, including dynamic possible callers, complete denominators, provenance, fail-inclusive `Unknown`, lifecycle eligibility, and handbook non-vacuity.
6. **Deterministic iteration accounting** — M6/M7/M10, including shared wave tokens, reconciliation before completion, synthesis-gap reopening, and ZULU exhaustion.
7. **PRF identity and semantic hashing** — M2/M10/M11, including `0A`, candidate, traceability, and CNF bindings.
8. **HBK identity and path/anchor normalization** — M2/M12, including the fixed handbook root, exact Unicode/path handling, stable anchors, title edits, and proven moves.
9. **Profile Synchronization closure** — M6/M10/M11, including complete direct/applicable and transitive dependency closure and identical confirmation closure.
10. **Record versus artifact hashing** — M1/M13, including independent record hashes, ordered file transport, artifact hash-domain identity, and post-hash binding.
11. **Canonicalization and exact exclusions** — M1/M11/M13, including LF/Unicode/order/whitespace handling, exact certification/DEC/CNF/artifact boundaries, and legal prior versus forbidden later references.
12. **Final envelope and package-member integrity** — M11/M13/M14, including identity/version-bound envelope fingerprints and complete finalized-file fingerprints.
13. **Packaging schemas, gate completeness, and deterministic ordering** — M13/M14, covering all five schemas, exact report check sets, evidence-set fingerprints, field order, empty values/tables, sort keys, counts, and cycle prevention.
14. **Acyclic, snapshot-consistent Exit E and final verification** — M14, from candidate report through the directly snapshot-bound signed outer manifest, with `Pending` through steps 1–5, reproducible verification receipts, and no authoritative post-step-6 completion artifact.

## Other deferred protocol families

The roadmap also explicitly carries source/claim/frontier evidence (M5), stateless invocation and cold resume (M6), tickets/probes/Human Hatch and reconciliation (M7), coverage and Exit A (M8), export acquisition/reconciliation (M9), synthesis and gaps (M10), decisions/confirmations/staleness (M11), handbook/equivalence usability (M12), canonical packaging (M13), Exit E (M14), full acceptance fixtures (M15), and corpus acceptance in a real dry run (M16).

Runtime/service work does not move those protocol families earlier or weaken their completion criteria. M4 supplies isolation, coordination, observability, and unsupported boundaries; later milestones make its surfaces protocol-aware.

## Current next milestone

Proceed to M1: deepen the structural Markdown AST and protocol model with typed-record, field-path, table, payload, and envelope boundaries. Then complete M2 collision extension and identity coordinates before expanding M3 workspace ownership validation and beginning M4 runtime foundations. Do not start semantic, gate, or packaging claims before their structures and validators can locate and enforce normative boundaries reliably.
