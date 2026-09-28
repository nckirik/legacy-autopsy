# Protocol: Legacy System Deconstruction, Assurance, and Reconstruction

**Version:** 4.1.2 (Canonical Reconstruction-Ready Edition)
**Generated-From:** cdl/0.2.0 (sha256:43227ea826c7c733f095ac0aa30add60184e5dc0904f3d9a16bb1529ccd97de5)
**Language:** cdl/0.2 (eir-format 2)
**Authority:** generated render of the canonical CDL sources; do not edit.

---

## Document declarations

### Capabilities

| Capability | Kind | Version |
| :-- | :-- | --: |
| hash.sha256 | deterministic | 1 |
| table.serialize | deterministic | 1 |
| identity.typed-id | deterministic | 1 |
| path.normalize | deterministic | 1 |
| canonical.markdown | deterministic | 1 |

### Registries

- normalized-units
- normalized-map-units

### Artifacts

- approved-evidence
- 0A-PREFLIGHT.md
- 0B-AUTH-MODEL.md
- 0C-SECURITY-PRIVACY.md
- 0D-GLOSSARY.md
- 0E-INDEX.md
- 0F-GLOBAL-STATE.md
- 0G-DECONSTRUCTION-STATE.md
- 0H-CHECKPOINT-SUMMARY.md
- 10-SOURCE-INVENTORY.md
- 11-TRAVERSAL-FRONTIER.md
- 12-CLAIM-EVIDENCE.md
- 13-EXPORT-RECONCILIATION.md
- 14-CONTRADICTIONS.md
- 15-DECISIONS.md
- 16-SOURCE-COVERAGE.md
- 17-ACQUISITION-CANDIDATES.md
- 18-CONFIRMATIONS.md
- 20-TRACEABILITY.md
- 21-COVERAGE-REPORT.md
- 22-GATE-REPORTS.md
- 90-ARCH-BLUEPRINT.md
- 91-DATA-MODEL.md
- 92-BUSINESS-RULES.md
- 93-USE-CASES.md
- 94-INTERFACES.md
- 95-DEPLOYMENT.md
- 96-NON-FUNCTIONAL-SECURITY.md

### Gates

- A-STRUCTURALLY-COMPLETE
- R-SWEPT
- Exit-A
- artifact-fingerprint
- artifact-completeness

### Workflow targets

- handoff
- section-pass

### Base reads

- 0A-PREFLIGHT.md
- 0D-GLOSSARY.md
- 0E-INDEX.md
- 0F-GLOBAL-STATE.md
- 0G-DECONSTRUCTION-STATE.md
- 0H-CHECKPOINT-SUMMARY.md
- 10-SOURCE-INVENTORY.md
- 11-TRAVERSAL-FRONTIER.md
- 12-CLAIM-EVIDENCE.md

### Modes

| Mode | Strict | Reads |
| :-- | :-- | --: |
| Preflight | false | 9 |
| Export Acquisition | false | 12 |
| Discovery | true | 13 |
| Ticket Resolution | true | 13 |
| Promotion Review | true | 13 |
| Sequential Reconciliation | false | 14 |
| Profile Synchronization | false | 13 |
| Cross-Reference-Reconciliation | false | 14 |
| Partial-Synthesis | false | 17 |
| Final-Synthesis | false | 17 |
| Reconstruction-Handoff | false | 20 |
| Human Hatch | false | 15 |
| Validation/Gate | false | 18 |

## 0.1. Normative language

This protocol is implementation-independent. A conforming implementation MAY use any programming language, runtime architecture, storage engine, agent framework, model provider, user interface, or execution tooling, including ordinary local scripts created by a capable coding harness, provided it preserves every normative behavior, deterministic algorithm, authority boundary, artifact, state transition, and conformance condition defined here. No reference implementation, product, service, CLI, UI, or agent harness is normative or required. The protocol-defined .extracted/ layout is interoperable protocol state, not an implementation source-tree requirement. In this document, conforming runtime means the implementation-neutral deterministic authority role defined here; it need not be a long-lived service or separate installed tool.

## 0.2. Guarantees

## 0.3. Non-guarantees

## 0.4. Conceptual flow

## 1.1. Plane 1 — Forensic Plane

The Forensic Plane preserves what was directly found, where it was found, how it is invoked, and what remains uncertain.

## 1.2. Plane 2 — Assurance Plane

The Assurance Plane proves denominator completeness, traversal closure, evidence support, export reconciliation, reference integrity, and gate outcomes.

## 1.3. Plane 3 — Synthesis Plane

The Synthesis Plane transforms closed forensic facts and assurance proofs into coherent reconstruction semantics.

## 1.4. Plane 4 — Human Reconstruction Handbook

The Human Reconstruction Handbook presents the confirmed bundle through progressive disclosure for architects, developers, domain experts, security reviewers, data engineers, QA, and operators.

## 1.5. Non-authoritative sidecars

## 2.1. Persona definition and purity

## 2.2. Canonical entry ownership

## 2.3. Five isolated POVs

## 2.4. Traversal tracks

## 2.5. Capability and runtime state

Capability state is semantic metadata, not a prefix-group tag.

## 2.6. Capability grouping and target decision

## 3.1. `0A-PREFLIGHT.md`

Maintain metadata and boundaries, the persona registry, concrete and DB-autonomous clusters, every persona-by-entry-by-track-by-environment-by-snapshot-by-POV cell, claim-backed applicability, the unmapped queue, and the acquisition register.

## 3.2. Foundational registries

0B owns auth and session; 0C owns privacy and security risk without compliance claims; 0D owns glossary; 0E owns index; 0F owns invariants and global state; 0G is the append-only invocation history; 0H is the derived resume accelerator.

## 4.1. ID generation

Every canonical typed record that participates in identity, reference, traceability, candidacy, or confirmation uses a legal protocol ID generated from its canonical key.

## 4.1.1. Semantic payload identity and certification envelopes

Every synthesis block, persona profile, and semantic handbook section has a canonical semantic payload and a separately delimited certification envelope.

## 4.1.2. Normative canonical hash profile and packaging artifacts

Unless a schema explicitly requests a raw-source byte fingerprint, all semantic records, DEC, CNF, Exit E Candidate Reports, candidate payload manifests, Exit E Content-Readiness Reports, scope certificates, and outer bundle manifests use this canonical hash profile.

## 4.2. Atomic unit rule

## 4.3. Source inventory denominator — `10-SOURCE-INVENTORY.md`

One record exists per concrete source unit.

## 4.4. Traversal frontier — `11-TRAVERSAL-FRONTIER.md`

Every discovered traversal boundary MUST create exactly one FRT record, including a boundary recognized as terminal, duplicate, excluded, or already mapped at discovery time. Immediate recognition MAY create the record directly in its terminal state, but MUST NOT omit it. Every discovered but not yet terminal traversal boundary remains open until terminally disposed.

## 4.5. Scope-specific source coverage — `16-SOURCE-COVERAGE.md`

One row exists per immutable SRC version by persona/owner by track by environment by snapshot scope. Rows are disjoint within a subgroup; a source participating in several tracks has separate rows and is counted once in each declared track report, while the aggregate unique-SRC report deduplicates by SRC version.

## 4.6. Stable supersession and tombstones

Records are never silently deleted. A renamed record retains its ID when identity is proven.

## 5.1. Dedicated prefix groups

Finite protocol enums MUST be validated from the normative Markdown AST using the complete schema/context-qualified field path, never an unanchored text match, a label-only lookup, or a token-prefix guess. Prefix-family tokens must be bracketed where declared, and every finite-enum path is declared here.

## 5.2. Claim-level evidence — `12-CLAIM-EVIDENCE.md`

Every claim record binds its subject, field or step, evidence level, anchor, and lineage so that materiality and projection classification are explicit.

## 5.3. Sanitized projections and helper limitations

A sanitized export produced by an approved local helper is an agent-visible evidence projection.

## 5.4. Contradictions — `14-CONTRADICTIONS.md`

Every contradiction record names the conflicting claims, scope, and stale impact.

## 5.5. Authoritative decisions — `15-DECISIONS.md`

Every decision record carries its content version and fingerprint, exact scope and denominator effect, approval status, and approval envelope bindings.

## 5.6. Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md`

Every helper-published candidate and operator action receives a row in 17-ACQUISITION-CANDIDATES.md and a deterministic CND ID whose canonical key binds the stable candidate or action identity, helper policy version, and pinned snapshot.

## 5.7. Human confirmations — `18-CONFIRMATIONS.md`

A CNF record confirms exact typed record IDs, semantic versions, and fingerprints against a pinned closed snapshot, never mutable whole-file bytes.

## 6.1. Atomic component schema

The atomic component record schema is declared below.

## 6.2. Dual-entry discovery

Persona-driven discovery is the only behavioral discovery path. There is no independent shared-code sweep.

## 6.3. Staging and promotion

A shared atomic component advances from S-DISCOVERED to S-STAGING-VERIFIED only when mandatory fields, claim evidence, source mappings, references, and ticket state validate.

## 6.4. Depromotion

## 6.5. Surgical evolution

## 7.1. Trust boundary

Live acquisition is an operator-side activity; agents MUST NOT invoke acquisition using production credentials.

## 7.2. Acquisition register

The acquisition register records every export artifact, its projection, and its acquisition status.

## 7.3. Logical explosion and virtual coordinates

Before POV extraction, each export is decomposed into atomic logical units whose virtual coordinates are canonical for entry coordinates, evidence, dependencies, frontiers, and tickets.

## 7.4. Required explosion coverage

Explosion coverage is enumerated per platform and must include every graph edge.

## 7.5. n8n human acquisition loop

Register, explode, trace, and reconcile one approved serialized export through a human acquisition loop, stopping at every missing reference boundary.

## 7.6. Normalized maps

Normalized maps are navigation-only indexes over one exploded export.

## 7.7. Export reconciliation — `13-EXPORT-RECONCILIATION.md`

**Artifact:** 13-EXPORT-RECONCILIATION.md

Every normalized export coordinate has exactly one dispositioned row.

## 8.1. Resume identity header

Every invocation begins with the identity header, and the conforming runtime independently validates declared actions, exact read/write scope, loaded fingerprints, checkpoint freshness, and mode rights before commit.

## 8.2. Strict single-scope modes

Discovery, Ticket Resolution, and Promotion Review bind exactly one persona, prefix, directory, cluster, track, POV, POV file, and at most one question ledger.

## 8.3. Invocation-mode enum and explicit multi-file modes

The exhaustive invocation-mode enum is the declared registry; multi-file modes enumerate their exact targets.

```text
Preflight | Export Acquisition | Discovery | Ticket Resolution | Promotion Review | Sequential Reconciliation | Profile Synchronization | Cross-Reference-Reconciliation | Partial-Synthesis | Final-Synthesis | Reconstruction-Handoff | Human Hatch | Validation/Gate
```

## 8.4. Mandatory read sets

All modes load protocol constants/statuses, 0A, 0D, 0E, 0F, 0G or 0H, relevant source inventory/frontier/claim records, and the active snapshot metadata.

## 8.5. Cold resume check

Before mutation, the executor returns the semantic assessment fields and the conforming runtime completes the deterministic validation fields before deciding whether to commit.

## 8.6. Concurrent persona traversal

Personas MAY run concurrently in one named iteration but MUST NOT write canonical shared POV or shared question files directly.

## 8.7. Stale checkpoint guard

Before writing, compare loaded checkpoint and source fingerprints with target record versions.

## 8.8. `0G` invocation log

Append the invocation log record after every invocation; 0H MAY compress current active states but MUST include its source 0G range and fingerprint.

## 9.1. Ticket schema and canonical IDs

Ticket IDs use MONIKER-PERSONA-POV-SEQUENCE; allocation resets per iteration/persona/target POV and is reserved atomically.

## 9.2. FSM and write rights

The conforming runtime validates every transition, escalation reason, required evidence field, actor right, and source state; an executor cannot assign a terminal or probe state merely by writing its token.

## 9.3. Human Hatch action prerequisites

Each Human Hatch subtype has distinct prerequisites; a fake probe or transition through T-PROBE-REQUIRED is forbidden.

## 9.4. Probe specification and asynchronous handoff

## 9.5. No-mock integrity

## 9.6. No-mock fallback payload

The no-mock fallback payload schema is declared below.

## 10.1. Evidence-backed `[R-SWEPT]`

A sweep record carries SWEEP-ID. A required traversal-matrix cell may be R-SWEPT only with a sweep record; every required cell has its own sweep, and claim-backed Not-Applicable cells have no sweep and remain visible in the denominator.

## 10.2. Per-kind coverage arithmetic

Coverage is computed separately for every applicable kind, track, snapshot, and environment, then aggregated; a global percentage MUST NOT hide a failed subgroup.

## 10.3. Composite Exit A — Deconstruction Closure

Exit A succeeds only when all declared conditions hold for the pinned scope.

## 10.4. Exit B — Stagnation

If a complete named iteration has zero valid mutations the pipeline halts as stagnant and reports remaining state; cosmetic edits do not count.

## 10.5. Exit C — Oscillation

If an actor-transition sequence of length 1 to 4 repeats for two consecutive completed cycles in one ticket trace, run one neutral Arbitration invocation, normally POV-5.

## 10.6. Exit D — Limit exhaustion and deterministic iteration accounting

The 26 canonical iteration tokens are a total per-system budget across the complete protocol run; the budget never resets, and the system MUST NOT proceed beyond ZULU.

## 11.1. Partial versus Final Synthesis

Partial-Synthesis runs before Exit A only for an explicitly bounded persona/cluster/track/snapshot; Final-Synthesis runs only after Composite Exit A passes.

## 11.2. Common synthesis block header

Every atomic synthesis block MUST contain a semantic payload followed by a separately delimited certification envelope.

## 11.3. Architecture blueprint — `90-ARCH-BLUEPRINT.md`

The architecture blueprint contains system boundary and context diagrams; actors/personas and external systems; observed legacy technology separated from target decisions; atomic MOD records; inter-module contracts and error propagation; deployment topology; architectural invariants and modernization consequences; capability navigation; and a prominent active-versus-disabled architecture view. Purely navigational module indexes are grouping containers under §4.2 and list MOD IDs only.

## 11.4. Entity — `91-DATA-MODEL.md`

## 11.5. Relationship — `91-DATA-MODEL.md`

## 11.6. State machine — `91-DATA-MODEL.md`

A state machine is required only for an eligible lifecycle scope: an entity aggregate or capability with a finite, behaviorally meaningful lifecycle, not every boolean or presentation state.

## 11.7. Database routine — `91-DATA-MODEL.md`

## 11.8. Business rule — `92-BUSINESS-RULES.md`

## 11.9. Use case — `93-USE-CASES.md`

## 11.10. Interface — `94-INTERFACES.md`

## 11.11. Deployment, configuration, and scheduling — `95-DEPLOYMENT.md`

## 11.12. NFR and security — `96-NON-FUNCTIONAL-SECURITY.md`

## 11.13. Persona profile

Profile Synchronization writes one profile and its 0A row atomically.

## 12.1. Traceability — `20-TRACEABILITY.md`

One atomic row per source-to-semantic mapping; multiple rows MAY represent multiple mappings, but source coverage is evaluated once per SRC and requires all applicable mappings. The ID registry, source denominator query fingerprints, lifecycle eligibility decisions, exclusions, and reciprocal reference check results are included in this file.

## 12.2. Required reciprocal references

## 12.3. Synthesis-ID timing

## 12.4. Dependency and impact provenance

## 12.5. Deterministic validation summary

Every validator run records schema/protocol/check-registry/validator versions; authoritative input file fingerprints and a deterministic input-set fingerprint; checks executed and stable check IDs; for every check and blocker row, the complete sorted evidence-binding tuples and recomputed §4.1.2 evidence-set fingerprint; pass/fail counts and exact offending IDs; coverage arithmetic by kind/track/environment; broken links, invalid diagrams, duplicate IDs, orphan/dangling references, stale blocks, unsupported statuses, and unresolved contradictions; any authorized human review input used by a check, including reviewer identity/authority, review scope, timestamp, exact candidate-bound fingerprints, limitations, and result; and output fingerprint and timestamp.

## 13.1. Audience paths and progressive disclosure

START-HERE.md MUST provide short audience paths through the handbook.

## 13.2. Required chapter behavior

## 13.3. Disabled and Dormant Features chapter

## 13.4. Handbook quality gate

## 14.1. Decision log

## 14.2. Modernization mapping templates

## 14.3. Business-rule ownership

## 14.4. Interface compatibility

## 14.5. Equivalence and acceptance suite

## 14.6. Sole-input invariant

## 15.1. Strict conditions

Exit E remains Pending through steps 1-5 of §15.4. The stage-4 Exit E Content-Readiness Report validates all content, decision, confirmation, and package-member conditions available at that stage, but it is not the Exit E pass artifact and MUST NOT state or imply Passed. Exit E transitions to Passed only at step 6 when all conditions below are true and the signed outer bundle manifest's hashed payload validates with the literal field EXIT-E-STATUS: Passed.

## 15.1.1. Universal deterministic packaging rules

The five packaging artifact types are EXIT-E-CANDIDATE-REPORT, CANDIDATE-PAYLOAD-MANIFEST, EXIT-E-CONTENT-READINESS-REPORT, SCOPE-CERTIFICATE, and OUTER-BUNDLE-MANIFEST. Each MUST instantiate the generic §4.1.2 payload and envelope boundaries and its exact ordered schema.

## 15.1.2. Exit E Candidate Report ordered payload schema

The Exit E Candidate Report is emitted before the candidate payload manifest. Its exact payload field order is declared below.

## 15.1.3. Candidate Payload Manifest ordered payload schema

The candidate payload manifest is emitted only after the immutable candidate report and references it one-way. Its exact payload field order is declared below.

## 15.1.4. Exit E Content-Readiness Report ordered payload schema

The stage-4 Exit E Content-Readiness Report is a pre-certificate artifact. Its exact payload field order is declared below.

## 15.2. Scope certificate

The scope certificate binds scope, exclusions, risks, and the prior certification chain. Its exact payload field order is declared below.

## 15.3. Package contents

## 15.4. Acyclic certification sequence and hash domains

Certification MUST follow the declared order; no later artifact is an input to an earlier hash.

## 15.5. Reproducible final-bundle verification

EXIT-E-FINAL-CHECKS-v1 is the exact final verification registry.

## 17.1. Required order

## 17.2. Migration guarantees

## 17.3. v3 final-artifact compatibility

## 19.1. Conforming implementation capabilities

## 19.2. Normative conformance cases

## 19.3. Corpus acceptance

## 20. Compact Artifact Ownership Matrix

The compact artifact ownership matrix records the semantic purpose, authoritative plane, authorized writers, and key non-writers for every artifact and record.

---

**End of Canonical Deconstruction Protocol v4.1.2**
