# Protocol: Legacy System Deconstruction, Assurance, and Reconstruction

**Version:** 4.1.2 (Canonical Reconstruction-Ready Edition)
**Generated-From:** cdl/0.2.0 (sha256:ac555923b97a096db32c6f57d926409290aab72443b00b90048ec7b8b4c8b021)
**Language:** cdl/0.2 (eir-format 2)
**Authority:** generated render of the canonical CDL sources; do not edit.

---

# Document Declarations

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

### Rule: normative-keywords

The declared key words are normative. MUST and MUST NOT are absolute protocol requirements. SHOULD and SHOULD NOT require a recorded reason when not followed. MAY is optional.

### Rule: source-units-and-claims

A source unit is one concrete, independently inventoryable unit in or affecting the system: a file, symbol, route, query, job, table, view, database routine, constraint, configuration key, workflow, node, edge, page, widget, action, binding, or equivalent. A component is one independently evidenced behavioral or structural record extracted from one source unit or from a precisely bounded portion of one source unit. A material claim is a fact whose falsity would alter behavior, scope, safety, compatibility, architecture, data semantics, or a reconstruction decision.

### Rule: authority-classes

Protocol work has three authority classes. Executor-owned semantic work covers source traversal, behavior and business-rule interpretation, semantic classification, ambiguity recognition, claim formulation, and reconstruction prose; an executor MAY be a coding harness, model, remote worker, or human, but every semantic classification that affects scope, closure, mutation, or a gate MUST record its actor, invocation or human authority, reason, evidence/claim bindings, and pinned snapshot. Conforming-runtime deterministic work covers identity, path normalization, structural parsing, canonicalization, ordering, fingerprints, scope and stale-input checks, finite-state transitions, coverage arithmetic, completeness, gates, package integrity, and other protocol-defined mechanical results; a conforming runtime MUST compute these values or independently reproduce them from authoritative inputs before commit, and executor-supplied deterministic values are proposals only and MUST NOT be accepted as authoritative unless deterministic validation reproduces them exactly. Human-authorized decisions cover explicit exclusions, confirmations, policy or business intent, irreducible domain meaning, and acceptance of external facts only where this protocol permits it; a model inference or runtime computation MUST NOT substitute for the required human identity, authority basis, review scope, and exact input bindings.

### Rule: predicate-classification

Every normative predicate that affects conformance, mutation authorization, coverage, candidacy, or a gate MUST be one of: a deterministic predicate computed by the conforming runtime; a semantic predicate classified by an executor or authorized human with the required provenance; or an unresolved semantic predicate recorded as Unknown with its affected scope and required ticket/question. Unless a section defines a stricter result, Unknown is fail-inclusive: the item remains applicable/in scope for closure, cannot justify omission or Not-Applicable, and blocks any success that depends on the predicate being false. Natural-language guidance that does not affect those outcomes need not create a protocol record.

### Enum: NORMATIVE-KEYWORD

- MUST
- MUST NOT
- REQUIRED
- SHOULD
- SHOULD NOT
- MAY

### Enum: AUTHORITY-CLASS

- Executor-owned semantic work
- Conforming-runtime deterministic work
- Human-authorized decisions

### Enum: PREDICATE-KIND

- deterministic predicate
- semantic predicate
- unresolved semantic predicate

## 0.2. Guarantees

When followed completely, the protocol guarantees identity and terminal disposition for every in-scope source unit, denominator-backed traversal closure, one evidence-bound claim per material fact, distinct capability states across environments and snapshots, logical-unit export acquisition without raw secrets, persona and POV isolation, finite-state ambiguity resolution, derivation-proven synthesis with targeted staleness, a navigable handbook, and sole-input reconstruction only under a signed outer manifest.

### Rule: protocol-guarantees

When followed completely, this protocol guarantees that every in-scope concrete source unit has a stable identity and a terminal disposition; traversal closure is proven from an explicit source denominator and frontier, not inferred from ticket silence; every material fact is represented by one claim with an evidence level, coordinate, and snapshot fingerprint; active, conditional, shadow, disabled, retired, and unknown capabilities remain distinct across environments and snapshots; serialized exports are acquired, exploded, reconciled, and traced at logical-unit level without exposing raw secrets; personas and POVs retain strict ownership and isolation during discovery; cross-layer ambiguity is resolved through a finite-state ticket lifecycle, runtime probes, or an explicit human hatch, never by guessing; synthesis artifacts are structured semantic projections of closed forensic and assurance records with derivation provenance and targeted stale propagation; a separate human reconstruction handbook is navigable without using dense identifier catalogs as its primary interface; and reconstruction may treat the delivered bundle as its sole input only after the signed outer bundle manifest authoritatively attests EXIT-E-STATUS: Passed under the strict non-cyclic Exit E sequence.

## 0.3. Non-guarantees

The protocol does not determine compliance, prove repository state equals production state, make low-evidence facts true by repetition, convert correlation to causation, decide modernization without accountable human approval, authorize production execution or dormant enablement, or make sidecars authoritative.

### Rule: protocol-non-guarantees

The protocol does not determine legal, privacy, security, or regulatory compliance; prove that unavailable production state equals repository state; make low-evidence facts true by repeating them in synthesis; convert correlation into causation; decide modernization architecture without accountable human approval; authorize execution against production, enable dormant production features, or move production data; or make generated machine sidecars or helper output authoritative semantic truth. All outputs remain a draft until the applicable human confirmation gates pass. Human architects and domain owners remain responsible for modernization decisions and final validation.

## 0.4. Conceptual flow

The conceptual flow runs from scope and snapshots through inventory, acquisition, isolated traversal, forensic records, reconciliation, Exit A, synthesis, handbook, Exit E candidacy and confirmation, content readiness, certificate, and the signed outer bundle manifest.

### Rule: conceptual-flow-stages

The conceptual flow runs from scope and snapshots through source inventory and persona/cluster ownership, acquisition and logical explosion where required, isolated persona + cluster + POV traversal, atomic forensic records + claim evidence + frontier updates, tickets/probes/reconciliation/staging/promotion, Composite Exit A: Deconstruction Closure, final structured synthesis, human reconstruction handbook and decisions, Exit E candidate validation and human confirmation, Exit E Content-Readiness Report, scope certificate, and the signed outer bundle manifest with EXIT-E-STATUS: Passed, to the signed, fingerprinted reconstruction bundle. Partial synthesis MAY run before Exit A for bounded review. It is provisional and cannot authorize reconstruction.

### Enum: FLOW-STAGE

- Scope and snapshots
- source inventory and persona/cluster ownership
- acquisition and logical explosion where required
- isolated persona + cluster + POV traversal
- atomic forensic records + claim evidence + frontier updates
- tickets, probes, reconciliation, staging, and promotion
- Composite Exit A: Deconstruction Closure
- final structured synthesis
- human reconstruction handbook and decisions
- Exit E candidate validation and human confirmation
- Exit E Content-Readiness Report
- scope certificate
- signed outer bundle manifest with EXIT-E-STATUS: Passed
- signed, fingerprinted reconstruction bundle

## 1.1. Plane 1 — Forensic Plane

The Forensic Plane preserves what was directly found, where it was found, how it is invoked, and what remains uncertain.

### Rule: forensic-plane-sources

Original source snapshots and approved sanitized evidence projections are evidence sources; the structured forensic records are the canonical extracted representation. 0G is append-only process history. 0H is derived and never replaces 0G. Normalized maps are navigation indexes, not behavioral truth.

### Rule: forensic-plane-writers

Isolated POV invocations have exactly one semantic POV target and at most one semantic ledger target; their fixed transaction bundle MAY also append directly produced source, frontier, claim, contradiction, coverage, persona-local unmapped-discovery-buffer, and audit-log records explicitly allowed by the identity header. They cannot mutate unrelated semantic files. Acquisition writes the acquisition register, normalized maps, reconciliation records, approved registry fragments, and its directly allowed discovery placeholders. Profile Synchronization exclusively writes persona semantic payloads and synchronized 0A persona rows. Authorized Human Hatch: Confirmation or Reconstruction-Handoff confirmation actions MAY update only the separately delimited profile certification envelope; they MUST NOT change profile semantic content, SEMANTIC-RECORD-VERSION, or SEMANTIC-CONTENT-FINGERPRINT. Sequential Reconciliation owns cross-persona/shared merges, atomic promotion and depromotion, buffered tickets, and merging persona-local unmapped discoveries into 0A.

### Enum: PLANE

- Forensic Plane
- Assurance Plane
- Synthesis Plane
- Human Reconstruction Handbook

## 1.2. Plane 2 — Assurance Plane

The Assurance Plane proves denominator completeness, traversal closure, evidence support, export reconciliation, reference integrity, and gate outcomes.

### Rule: assurance-plane-sources

Assurance records are authoritative for scope, denominator membership, evidence claims, coverage arithmetic, reconciliation, contradiction status, and gate results. They cite forensic evidence but do not replace it.

### Rule: assurance-plane-writers

Acquisition, isolated traversal invocations, Sequential Reconciliation, Cross-Reference-Reconciliation, Human Hatch, and deterministic validators write only their explicitly owned assurance sections. Gate reports MUST be generated from the authoritative records and MUST include deterministic input fingerprints.

## 1.3. Plane 3 — Synthesis Plane

The Synthesis Plane transforms closed forensic facts and assurance proofs into coherent reconstruction semantics.

### Rule: synthesis-plane-sources

The structured catalogs are authoritative for synthesized architecture, entities, rules, use cases, interfaces, deployment/configuration, and NFR/security semantics. Every block MUST cite DERIVED-FROM IDs and fingerprints. Synthesis MUST NOT silently create a legacy fact absent from the Forensic or Assurance Plane.

### Rule: synthesis-plane-writers

Partial-Synthesis and Final-Synthesis create or update semantic payloads and their semantic versions/fingerprints. Cross-Reference-Reconciliation MAY update only reference fields after IDs exist; any reference field included in the semantic payload is a semantic change and MUST increment the semantic record version and fingerprint. Authorized confirmation attaches only the certification envelope and transitions eligible blocks from B-POPULATED to B-CONFIRMED without regenerating or changing semantic content. Synthesis writes discovered gaps only to a synthesis gap buffer; it never opens tickets directly.

## 1.4. Plane 4 — Human Reconstruction Handbook

The Human Reconstruction Handbook presents the confirmed bundle through progressive disclosure for architects, developers, domain experts, security reviewers, data engineers, QA, and operators.

### Rule: handbook-plane-sources

The handbook is a synchronized, human-readable projection of B-populated or B-confirmed Synthesis Plane records, assurance scope, and approved decisions. It is not an independent semantic authority. Its complete semantic payload MUST be generated as B-POPULATED, assigned legal HBK identities, and fingerprinted before the Exit E Candidate Report and human review. Semantic corrections MUST be made in the upstream authoritative record first and then re-projected before candidacy. Confirmation only attaches the handbook certification envelope; it MUST NOT regenerate or alter semantic content. Broken derivation marks the affected handbook section B-STALE.

### Rule: handbook-plane-writers

Reconstruction-Handoff generates and updates handbook semantic payloads before candidacy. Human editors MAY improve wording, examples, navigation, and diagrams before candidate fingerprinting only when semantic meaning remains unchanged; after candidacy, any content edit is a semantic payload change that invalidates the candidate and restarts from populated generation. Authorized confirmation actions MAY update only certification envelopes. Semantic edits require upstream changes and reprojection.

## 1.5. Non-authoritative sidecars

JSON, SQLite, graph, search, code-generation, or validator sidecars may be generated for automation; they carry source fingerprints, schema version, generation time, and a prominent NON-AUTHORITATIVE-DERIVATIVE marker, and conflicts resolve in favor of canonical Markdown records and original evidence.

### Rule: sidecar-derivatives

JSON, SQLite, graph, search, code-generation, or validator sidecars MAY be generated for automation. They MUST carry source fingerprints, schema version, generation time, and a prominent NON-AUTHORITATIVE-DERIVATIVE marker. Conflicts are resolved in favor of the canonical Markdown records and original evidence, in that order for semantics.

## 2.1. Persona definition and purity

A Persona is a distinct human actor, service actor, external system, or autonomous runtime class entering through a concrete execution surface. Personas are selected by distinct physical or virtual entry paths, and persona directories are scoped by the purity invariant.

### Rule: persona-purity

A Persona is a distinct human actor, service actor, external system, or autonomous runtime class entering through a concrete execution surface. Personas are selected by distinct physical or virtual entry paths, not merely business titles. Roles sharing identical routes, triggers, and paths SHOULD be grouped; materially different data visibility, authorization, or execution paths require distinct personas. Persona Purity Invariant: a personas/{persona-directory}/ directory MUST document only that persona's triggers, invocation paths, local UI/runtime state, and usage references. Shared components are linked, not copied as another persona's behavior. Purity does not imply exclusive component ownership.

### Rule: persona-container

personas/ is the sole container for persona-scoped and canonical shared forensic files. personas/_shared/ is a reserved non-persona directory containing canonical shared POV and question records; the leading underscore provides deterministic visual/sort separation and does not create a persona. Persona discovery, counting, purity, and profile validation MUST ignore that reserved directory as a persona while still validating its shared-record ownership rules.

### Rule: persona-identity-grammar

Every active persona MUST have one globally unique uppercase prefix matching [A-Z][A-Z0-9]{1,7} and one explicit lowercase persona slug matching [a-z0-9]+(?:-[a-z0-9]+)*. Prefix uniqueness is validated in 0A-PREFLIGHT.md; reuse is forbidden, including retired personas. The canonical persona-directory basename is exactly <persona-prefix>-<persona-slug> and therefore matches [A-Z][A-Z0-9]{1,7}-[a-z0-9]+(?:-[a-z0-9]+)*. The 0A registry stores the persona slug explicitly; it MUST NOT be inferred from, transliterated from, or normalized from the display name. The basename prefix, registry prefix, invocation Persona Prefix, ticket prefix, and PRF owner coordinate MUST agree exactly. A non-reserved directory directly beneath personas/ that is unregistered, malformed, or prefix-mismatched is invalid; _shared MUST NOT appear as a persona prefix, slug, basename, or registry row.

## 2.2. Canonical entry ownership

Every concrete entry coordinate has exactly one canonical owner persona; intentional shared surfaces use a SHARED-ENTRY record, and database-side autonomous behavior is owned by the reserved db-autonomous persona.

### Rule: entry-ownership

Every concrete entry coordinate MUST have exactly one canonical owner persona. Intentional shared entry surfaces MUST use a SHARED-ENTRY record that owns the coordinate once; lists participating personas and each local invocation condition; identifies one coordination owner for registry maintenance; does not weaken persona purity; and maps each persona to its own invocation reference. Duplicate ownership without a shared-entry record is a gate failure.

### Rule: db-autonomous-ownership

Database-side autonomous behavior uses the reserved persona db-autonomous with prefix DBR when it fires because of database mutation, timer, internal rule, or database event rather than an app entry. An app routine that causes a trigger to fire remains owned by its app persona; the trigger body is owned by DBR and cross-linked through a DR record. Procedures explicitly invoked only through one app entry remain app-owned but are still extracted by POV-4. Mixed invocation requires a shared-entry record. Ownership is determined from firing mechanism, never convenience.

## 2.3. Five isolated POVs

Five isolated POVs partition semantic extraction: frontend, backend, background, data movement, and system wiring. A POV MUST NOT fill another POV's semantic fields directly.

### Rule: pov-isolation

POV-1 Frontend Workflows: client UI, server-rendered views, widget behavior, presentation state, validation before transport, templates, and interaction flow. POV-2 Backend Workflows: domain decisions, business logic, service pipelines, validation engines, and transactional intent independent of transport and physical schema. POV-3 Background Workflows: schedules, queues, event loops, delayed execution, retries, async coordination, and long-running workers. POV-4 Data Movement Workflows: reads, writes, queries, transactions, ORM behavior, caches, files, tables, views, constraints, triggers, stored routines, generated fields, and persistence lifecycle. POV-5 System Wiring Workflows: routes, middleware, transport contracts, graph edges, callbacks, integrations, orchestration, and inter-system coupling. A POV MUST NOT fill another POV's semantic fields directly; cross-POV findings become a ticket or persona-local shared buffer entry.

## 2.4. Traversal tracks

Every entry cluster and invocation declares exactly one traversal track: Normal, Shadow/Conditional, or Disabled. Tracks require separate invocations even when they share source units.

### Rule: track-declarations

Every entry cluster and invocation MUST declare exactly one TRAVERSAL-TRACK. Normal: currently active, normally exposed behavior for the specified environment/snapshot. Shadow/Conditional: flag-gated, environment-dependent, alternate, shadow, indirectly exposed, or not normally selected behavior. Disabled: explicitly dormant, disabled, disconnected, legacy-retained, or configured off behavior. Tracks require separate invocations even when they share source units. They MAY cross-link the same SRC IDs. Disabled is not equivalent to dead code, retired code, irrelevant code, or out-of-scope code.

## 2.5. Capability and runtime state

Capability state is semantic metadata, not a prefix-group tag.

### Rule: capability-state-rules

State MAY differ by environment or snapshot. Absence of production evidence MUST NOT be converted into deployed behavior. In-scope disabled units require the same atomic static coverage as active units and block Exit A when undisposed. A disabled capability MUST NOT be enabled in production for a probe. Runtime observation is permitted only in an isolated sandbox with explicit human authorization, masked/non-production data, and a probe specification that states the disabled behavior and safety boundaries.

### Field: CAPABILITY-STATE-BLOCK

| Field | Type | Required |
| :-- | :-- | :-- |
| capability-state | REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE | true |
| runtime-state-matrix | string | true |
| historical-intended-personas | string | true |
| dependencies | string | true |
| last-known-state | string | true |

### Field: RUNTIME-STATE-MATRIX

| Field | Type | Required |
| :-- | :-- | :-- |
| environment | string | true |
| snapshot | string | true |
| state | string | true |
| enable-disable-mechanism | string | true |
| clm-evidence | string | true |
| last-known-at | string | true |

## 2.6. Capability grouping and target decision

A capability groups active and dormant source units for navigation, remains non-promotable, and never inherits member evidence or coverage. Every disabled or dormant capability preserves exact members, state, mechanism, dependencies, history, risk, and one explicit target decision.

### Rule: capability-target-decision

A capability may group active and dormant source units for navigation, but the CAP is non-promotable and inherits no evidence, status, or coverage from members. Every effective-denominator Disabled SRC, CMP, and UC MUST map to exactly one canonical CAP; the only alternative is an exact approved exclusion that removes that unit from the effective denominator. Shadow/Conditional units MAY map to a CAP when they represent a governed capability. CAP membership is mandatory in disabled traceability, synthesis, handbook output, and target-decision validation. Multiple or conflicting CAP memberships are invalid unless an explicit typed relationship and approved DEC define why the memberships overlap, identify one canonical target-decision owner, and prevent duplicate denominator or decision effects. Such a relationship does not make both memberships canonical and MUST NOT silently duplicate a unit. Every disabled/dormant capability MUST preserve exact member IDs including all effective-denominator Disabled SRC/CMP/UC units; status by environment/snapshot; enable/disable mechanism and evidence; dependencies and affected personas; historical/intended behavior; safety and reconstruction risk; one explicit target decision Restore, Preserve Dormant, Redesign, or Retire; and decision owner, rationale, and approval state. Observed legacy facts and target decisions MUST appear in separate fields and handbook sections.

### Enum: TARGET-DECISION

- Restore
- Preserve Dormant
- Redesign
- Retire

## 3.1. `0A-PREFLIGHT.md`

Maintain metadata and boundaries, the persona registry, concrete and DB-autonomous clusters, every persona-by-entry-by-track-by-environment-by-snapshot-by-POV cell, claim-backed applicability, the unmapped queue, and the acquisition register.

### Rule: preflight-boundaries

Agents MAY append acquisition reference placeholders directly where Acquisition permits. Bound Discovery invocations MUST NOT write the global Unmapped Discovery Queue; they append the persona-local buffer. Only Preflight or Profile Synchronization may assign personas or alter boundaries. Sequential Reconciliation alone merges persona-local unmapped discoveries into 0A. Wildcards cannot close and Candidate remains in the denominator.

### Rule: database-autonomous-clusters

Database-side autonomous clusters use the same schema as other clusters and apply the DBR ownership rule.

### Field: PERSONA-REGISTRY

| Field | Type | Required |
| :-- | :-- | :-- |
| persona | string | true |
| prefix | string | true |
| persona-slug | string | true |
| persona-directory | string | true |
| actor-class | string | true |
| canonical-entry-ids | string | true |
| scope-paths | string | true |
| registry-state | string | true |

### Field: REQUIRED-TRAVERSAL-MATRIX

| Field | Type | Required |
| :-- | :-- | :-- |
| persona-owner | string | true |
| concrete-entry | string | true |
| track | string | true |
| environment | string | true |
| snapshot | string | true |
| pov | string | true |
| applicability | string | true |
| status | string | true |
| swp-id | string | true |
| claims | string | true |

### Field: PREFLIGHT-HEADER

| Field | Type | Required |
| :-- | :-- | :-- |
| system-namespace | string | true |
| target-repository-evidence-roots | string | true |
| primary-technology-stack | string | true |
| protocol-version | string | true |
| current-iteration | REGISTRY-PREFLIGHT-CURRENT-ITERATION | true |
| default-max-traversal-depth | int | true |
| included-environments-snapshots | string | true |
| explicit-scope-policy | string | true |
| protocol-migration-metadata | string | true |

### Field: ENTRY-CLUSTER

| Field | Type | Required |
| :-- | :-- | :-- |
| concrete-entry-coordinates | set<string> | true |
| canonical-owner-persona | string | true |
| primary-pov | REGISTRY-ENTRYCLUSTER-PRIMARY-POV | true |
| traversal-track | REGISTRY-ENTRYCLUSTER-TRAVERSAL-TRACK | true |
| capability-state | REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE | true |
| environment-snapshot | string | true |
| source-ids | set<string> | true |
| max-depth | int | true |
| traversal-status | REGISTRY-ENTRYCLUSTER-TRAVERSAL-STATUS | true |
| required-traversal-matrix | string | true |

### Field: UNMAPPED-DISCOVERY

| Field | Type | Required |
| :-- | :-- | :-- |
| coordinate-src | string | true |
| discovered-by-invocation | string | true |
| suggested-pov-persona | string | true |
| traversal-status | REGISTRY-UNMAPPEDDISCOVERY-TRAVERSAL-STATUS | true |

### Field: DISCOVERY-BUFFER

| Field | Type | Required |
| :-- | :-- | :-- |
| concrete-coordinate-proposed-src-fingerprint | string | true |
| discovering-invocation-persona-cluster-track-pov | string | true |
| discovery-method-boundary-frt | string | true |
| suggested-persona-cluster-track-pov | string | true |
| reason-unmapped-materiality | string | true |
| dedup-key | hash | true |
| state | REGISTRY-UNMAPPEDDISCOVERYBUFFER-STATE | true |

## 3.2. Foundational registries

0B owns auth and session; 0C owns privacy and security risk without compliance claims; 0D owns glossary; 0E owns index; 0F owns invariants and global state; 0G is the append-only invocation history; 0H is the derived resume accelerator.

### Rule: registry-append-only

0G is chronological and append-only. 0H is a derived accelerator: it MAY compress current active states but MUST include its source 0G range and fingerprint, and a stale 0H forces authoritative fallback.

### Field: SHARED-ENTRY

| Field | Type | Required |
| :-- | :-- | :-- |
| concrete-entry-src | string | true |
| participating-personas-and-invocation-conditions | string | true |
| coordination-owner | string | true |
| persona-invocation-pointers | set<string> | true |
| environment-snapshot-track | string | true |
| material-claims-version | string | true |

### Field: AUTH-MECHANISM

| Field | Type | Required |
| :-- | :-- | :-- |
| mechanism-credential-form | string | true |
| issuer-validator-storage | string | true |
| creation-refresh-expiry-revocation-termination | string | true |
| authorization-bindings-enforcement-points | string | true |
| personas-interfaces-environments | string | true |
| secrets-handling-audit | string | true |
| material-claims-version | string | true |

### Field: PRIVACY-FINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| data-threat-review-category | string | true |
| affected-src-er-if-persona-environment | string | true |
| observed-handling-risk-unknowns | string | true |
| human-review-required | string | true |
| material-claims-version | string | true |
| synthesis-sec-mirror | string | true |

### Field: CAPABILITY-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| record-kind | string | true |
| canonical-member-src-cmp-uc-if-cfg-dr-ids | set<string> | true |
| overlapping-cap-relationships-decisions | string | true |
| runtime-state-matrix-mechanism-last-known-state | string | true |
| historical-and-intended-personas-behavior | string | true |
| dependencies-safety-risk | string | true |
| observed-claims | set<string> | true |
| target-decision-state | REGISTRY-CAPABILITYRECORD-TARGET-DECISION-STATE | true |
| target-decision-dec | string | false |
| version-supersession | string | true |

### Field: ROUTINE-CROSS-REFERENCE

| Field | Type | Required |
| :-- | :-- | :-- |
| src | set<string> | true |
| source-coordinate | string | true |
| canonical-extract | string | true |
| personas | string | true |
| cross-layer-consumers | set<string> | false |

### Field: INVARIANT

| Field | Type | Required |
| :-- | :-- | :-- |
| constraint | string | true |
| verification-vectors | set<string> | true |
| discovered-in | string | true |

### Field: SHARED-STATE

| Field | Type | Required |
| :-- | :-- | :-- |
| medium-shape | string | true |
| initializer-mutators-readers | set<string> | true |
| lifecycle-invalidation | string | true |
| claims | set<string> | true |

## 4.1. ID generation

Every canonical typed record that participates in identity, reference, traceability, candidacy, or confirmation uses a legal protocol ID generated from its canonical key.

### Rule: canonical-key

canonical-key = protocol-id-type + "|" + system-namespace + "|" + normalized-owner-coordinate + "|" + semantic-discriminator. hash = uppercase first 12 hexadecimal characters of SHA-256 of canonical-key UTF-8. id = TYPE + "-" + hash. Packaging artifacts are not typed records and do not receive protocol record IDs.

### Rule: src-kind-and-coordinates

For SRC, include kind: SRC-KIND-HASH, where KIND is one of the declared SRC kinds. Classification uses the most specific executable/data kind; a unit with several roles has one primary SRC kind and explicit typed relationship records to the others, preventing denominator double-counting. Stable platform IDs take precedence in normalized coordinates. Unicode is preserved exactly; path separators and harmless whitespace are normalized, but identifiers are never transliterated. A proven move or rename retains its original identity canonical key and records the new coordinate as versioned location metadata and an alias; the path is not rehashed.

### Rule: normalized-relative-path

A normalized relative path is relative to the applicable workspace root or package root, serialized with POSIX "/", has no leading "./", removes repeated separators and "." path segments, and rejects every ".." segment. Path case and Unicode code points are preserved exactly; implicit Unicode normalization, transliteration, and case folding are forbidden. The root used for normalization MUST be declared by the enclosing schema or package operation and remain fixed for the artifact or record set.

### Rule: prf-hbk-identity

PRF identifies one persona profile; its normalized owner coordinate is the globally unique persona prefix, and its semantic discriminator is the fixed literal persona-profile. HBK identifies one semantic handbook section; its normalized owner coordinate is the normalized handbook-relative path plus a hash separator plus its explicitly assigned stable section anchor, and its semantic discriminator is the fixed literal handbook-section. The heading text is not an identity source. An HBK stable section anchor is an explicitly assigned opaque ASCII lowercase token matching [a-z0-9]+(?:-[a-z0-9]+)*; it is never renderer-generated or inferred from heading text.

### Rule: collision-extension

A collision MUST be resolved by extending both colliding hashes to 16, then 20, up to 64 characters; existing shorter IDs receive aliases and controlled supersession, never silent mutation.

### Rule: id-registry-checks

The ID registry in 20-TRACEABILITY.md MUST verify global uniqueness, type correctness, canonical-key uniqueness, ticket-tuple uniqueness, coverage-tuple uniqueness, candidate/action identity uniqueness, PRF-per-persona-prefix uniqueness, HBK path-and-anchor uniqueness, and tombstone reservation. It MUST also verify reciprocal PRF to 0A persona row, HBK to handbook path/anchor and upstream synthesis records, and exact candidate/CNF bindings to PRF/HBK IDs. IDs are stable across movement when semantic identity and stable platform identity remain unchanged. A semantic split creates new IDs and tombstones or supersedes the old record.

### Enum: ID-PREFIX

- SRC-FILE
- SRC-SYM
- SRC-ROUTE
- SRC-RPC
- SRC-JOB
- SRC-QUERY
- SRC-TABLE
- SRC-VIEW
- SRC-DR
- SRC-CONSTRAINT
- SRC-GENERATED
- SRC-CFG
- SRC-WORKFLOW
- SRC-NODE
- SRC-EDGE
- SRC-PAGE
- SRC-WIDGET
- SRC-ACTION
- SRC-BINDING
- SRC-QUEUE
- SRC-EVENT
- SRC-STORAGE
- SRC-IFACE
- SRC-OTHER
- CMP
- CLM
- ER
- REL
- BR
- UC
- IF
- DEP
- CFG
- SCHED
- NFR
- SEC
- SM
- DR
- CAP
- PRV
- FLT
- MOD
- PRF
- HBK
- INV
- AUTH
- STATE
- DEC
- CON
- FRT
- SWP
- COV
- CND
- CLU
- SHR
- DSC
- EXP
- GAP
- CNF

### Enum: ID-ITERATION

- ALFA
- BRAVO
- CHARLIE
- DELTA
- ECHO
- FOXTROT
- GOLF
- HOTEL
- INDIA
- JULIETT
- KILO
- LIMA
- MIKE
- NOVEMBER
- OSCAR
- PAPA
- QUEBEC
- ROMEO
- SIERRA
- TANGO
- UNIFORM
- VICTOR
- WHISKEY
- X-RAY
- YANKEE
- ZULU

## 4.1.1. Semantic payload identity and certification envelopes

Every synthesis block, persona profile, and semantic handbook section has a canonical semantic payload and a separately delimited certification envelope.

### Rule: payload-mandatory-fields

A persona profile payload begins with its mandatory PRF-ID; a handbook-section payload begins with its mandatory HBK-ID and includes separately ordered HANDBOOK-RELATIVE-PATH and STABLE-SECTION-ANCHOR fields before its semantic record version and fingerprint carrier. These typed IDs and HBK coordinate fields are part of the semantic payload and candidate/CNF binding; headings are not identity parser inputs.

### Rule: semantic-record-version

SEMANTIC-RECORD-VERSION is a positive integer changed only when semantic payload meaning or payload bytes under the canonicalization rules change.

### Rule: semantic-content-fingerprint

SEMANTIC-CONTENT-FINGERPRINT is the sha256 of the record's canonical semantic payload under §4.1.2. The hash input includes record ID, semantic record version, schema/protocol version, and all semantic fields, prose, tables, diagrams, examples, and resolved semantic references with their bound versions/fingerprints. The carrier field that stores SEMANTIC-CONTENT-FINGERPRINT is excluded from its own hash input.

### Rule: certification-envelope

A certification envelope is status and approval metadata attached after payload generation: BLUEPRINT-STATUS or HANDBOOK-STATUS, CNF references, confirmation-envelope version, approval/signature identities and timestamps, and envelope audit fields. The semantic-content hash and semantic record version MUST exclude the complete certification envelope, including B/handbook status, CNF reference, candidate report or candidate payload manifest references explicitly permitted by the envelope schema, and envelope digests. A certification envelope MUST NOT reference the later Exit E Content-Readiness Report, scope certificate, outer bundle manifest, or any other artifact that did not exist when the envelope was attached.

### Rule: envelope-mutation-invariance

Attaching a confirmation or changing B-POPULATED to B-CONFIRMED changes only the envelope and MUST NOT change the semantic record version or semantic-content fingerprint. Any edit to semantic prose, diagrams, examples, claims, derivation references, semantic cross-references, PRF identity, or HBK identity increments SEMANTIC-RECORD-VERSION, recomputes SEMANTIC-CONTENT-FINGERPRINT, invalidates the old envelope, and triggers targeted stale propagation.

### Rule: confirmation-binding

A CNF confirms exact typed record IDs including PRF and HBK, SEMANTIC-RECORD-VERSION values, and SEMANTIC-CONTENT-FINGERPRINT values, never mutable whole-file bytes or status-bearing envelopes. Validators MUST parse the payload/envelope boundary from the Markdown AST, recompute payload hashes, and reject envelope fields inside the semantic payload or semantic fields inside the envelope.

## 4.1.2. Normative canonical hash profile and packaging artifacts

Unless a schema explicitly requests a raw-source byte fingerprint, all semantic records, DEC, CNF, Exit E Candidate Reports, candidate payload manifests, Exit E Content-Readiness Reports, scope certificates, and outer bundle manifests use this canonical hash profile.

### Rule: canonical-profile

The canonical hash profile parses the bounded payload as a normative Markdown AST and rejects malformed or overlapping boundaries; normalizes line endings to LF while preserving Unicode code points exactly without transliteration or implicit Unicode normalization; serializes fields in the deterministic order declared by the applicable schema and serializes tables with declared column order, normalized delimiter/alignment syntax, and insignificant cell-edge whitespace while preserving semantic row order unless the schema declares a deterministic row sort key; removes insignificant trailing whitespace, AST-separator blank-line variance, and transport-only indentation while preserving whitespace in code spans, fenced blocks, diagrams, and other literal nodes, and emits exactly one terminal LF; includes the protocol/schema version, record or artifact identity, and every in-boundary semantic/content node, and excludes only the exact carrier and envelope regions authorized by the applicable schema.

### Rule: record-boundary-and-file-transport

For a typed record in a multi-record Markdown file, the hash boundary begins at that record's typed-ID heading and ends immediately before the next typed-record heading of the same or higher level, or at end of file. Each record is canonicalized and hashed independently after applying its exact envelope exclusion. The containing file has a separate FILE-TRANSPORT fingerprint over its normalized relative path and ordered sequence of record type, record ID, record version where applicable, and record payload fingerprint bindings. That file fingerprint detects ordering or transport changes but is not a semantic-record identity, is not substituted for a record fingerprint, and MUST NOT be bound by a CNF.

### Rule: record-envelope-exclusions

DEC uses this profile over its heading-delimited decision content and excludes exactly DECISION-CONTENT-FINGERPRINT as its own carrier and the existing DECISION-APPROVAL-ENVELOPE through END-DECISION-APPROVAL-ENVELOPE region. CNF uses this profile over its heading-delimited confirmation payload and excludes exactly CNF Payload Fingerprint as its own carrier and the existing CNF-DIGEST-SIGNATURE-ENVELOPE through END-CNF-DIGEST-SIGNATURE-ENVELOPE region. Semantic synthesis, PRF, and HBK records use §4.1.1 and exclude exactly their existing certification envelope. SEMANTIC-CONTENT-FINGERPRINT and DECISION-CONTENT-FINGERPRINT carriers are excluded from their own preimages.

### Rule: envelope-fingerprint-and-binding

Every populated certification envelope, DEC approval envelope, and CNF digest/signature envelope receives an externally stored canonical envelope fingerprint. Its preimage is the UTF-8 domain prefix ENVELOPE followed by the normalized containing path, target typed record ID, envelope kind (CERTIFICATION, DECISION-APPROVAL, or CNF-DIGEST-SIGNATURE), positive envelope version, one LF, and the complete canonically serialized Markdown AST region between the exact envelope markers; fields are separated by a vertical bar. No field in that region is excluded: status, target-content binding, approver/signatory identity, authority, timestamps, audit data, signatures, and any payload-fingerprint carrier inside the envelope are all included. Missing, duplicate, malformed, overlapping, or versionless final envelopes are invalid. An envelope binding is the exact tuple of target typed record ID, envelope kind, envelope version, and canonical envelope fingerprint; the CERTIFICATION, DECISION-APPROVAL, and CNF-DIGEST-SIGNATURE envelope binding fields MUST serialize that tuple and MUST be recomputed rather than trusted.

### Rule: package-member-fingerprint

Every final package member also receives an external PACKAGE-MEMBER file fingerprint. For Markdown, its preimage is UTF-8 PACKAGE-MEMBER followed by the normalized package-relative path, one LF, and the complete canonicalized file AST, including every record fingerprint carrier and every certification, approval, digest, and signature envelope with no excluded node. For non-Markdown members, the same path-bound domain prefixes the exact immutable bytes. This fingerprint is stored only in later manifests or reports, never inside the member whose bytes it covers, so it cannot self-reference. Containing File Fingerprint, PACKAGE-MEMBER-FILE, and outer-manifest member fingerprints mean this full finalized-file fingerprint, never the envelope-excluding FILE-TRANSPORT fingerprint.

### Rule: evidence-set-fingerprint

An EVIDENCE-SET fingerprint is the value serialized in every Evidence Fingerprint check or blocker cell in §15.1 and binds that row to the exact authoritative evidence and deterministic validation inputs used for its result. Each binding is the tuple of normalized authoritative input path, binding kind, typed record or artifact ID or None, version or None, and authoritative fingerprint, using exactly one declared kind: RAW-SOURCE-BYTES, RECORD-PAYLOAD, SEMANTIC-CONTENT, DECISION-CONTENT, CNF-PAYLOAD, ARTIFACT-PAYLOAD, CANONICAL-ENVELOPE, FILE-TRANSPORT, PACKAGE-MEMBER-FILE, VALIDATION-SUMMARY, DENOMINATOR-QUERY, or PROTOCOL-REGISTRY. Binding kind selects the only legal fingerprint domain for that row input; a different valid fingerprint over the same file or record is a mismatch, not an alternative. The conforming runtime MUST reject an unknown kind, a version inconsistent with the table, or duplicate tuples; sort bindings bytewise by the complete tuple after §4.1 normalization; and compute sha256 over UTF-8 EVIDENCE-SET plus the packaging artifact type, one LF, row kind CHECK or BLOCKER, one LF, stable row ID, one LF, the canonically serialized System / Protocol / Snapshot value, one LF, and the canonical table serialization of the sorted tuples with columns in tuple order. The set MUST contain every input that can change the row result and at least one DENOMINATOR-QUERY, PROTOCOL-REGISTRY, or authoritative record/artifact binding even for an evidence-backed zero domain. It MUST NOT include the containing report, its own carrier, or a later artifact. A conforming runtime MUST recompute this fingerprint; an executor-supplied digest is non-authoritative and a mismatch is a structural failure.

### Rule: artifact-payload-boundaries

Each packaging artifact is one bounded payload. Its HASH-DOMAIN identity is exactly artifact type plus a vertical bar plus normalized relative path using §4.1 normalization, and that identity is included in the canonical payload hash preimage. The payload fingerprint MUST NOT appear anywhere in its own preimage. After the payload fingerprint is computed, the separate POST-HASH-ARTIFACT-INSTANCE-BINDING is artifact type, normalized relative path, and payload fingerprint; this binding identifies the resulting artifact instance, is computed only post-hash, and is not part of its own preimage. Packaging artifacts MUST use the exact generic ARTIFACT-PAYLOAD boundaries with ARTIFACT-TYPE, ARTIFACT-PATH, PROTOCOL / ARTIFACT-SCHEMA VERSION, and ARTIFACT-PAYLOAD-FINGERPRINT fields and a matching ARTIFACT-DIGEST-SIGNATURE-ENVELOPE; the payload fingerprint covers exactly the AST content between the ARTIFACT-PAYLOAD and END-ARTIFACT-PAYLOAD markers, excluding only ARTIFACT-PAYLOAD-FINGERPRINT and the complete matching digest/signature envelope, which runs to END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE. The EXIT-E-CANDIDATE-REPORT, CANDIDATE-PAYLOAD-MANIFEST, EXIT-E-CONTENT-READINESS-REPORT, SCOPE-CERTIFICATE, and OUTER-BUNDLE-MANIFEST artifacts MUST use this schema with their literal artifact type. Their type, path, member lists, status fields, references, and content timestamps remain hashed payload. No artifact may specialize or substitute boundary marker names; this generic schema supersedes certificate- or report-specific digest markers.

### Enum: ENVELOPE-KIND

- CERTIFICATION
- DECISION-APPROVAL
- CNF-DIGEST-SIGNATURE

### Enum: EVIDENCE-BINDING-KIND

- RAW-SOURCE-BYTES
- RECORD-PAYLOAD
- SEMANTIC-CONTENT
- DECISION-CONTENT
- CNF-PAYLOAD
- ARTIFACT-PAYLOAD
- CANONICAL-ENVELOPE
- FILE-TRANSPORT
- PACKAGE-MEMBER-FILE
- VALIDATION-SUMMARY
- DENOMINATOR-QUERY
- PROTOCOL-REGISTRY

### Enum: BINDING-ROW-KIND

- CHECK
- BLOCKER

### Enum: ARTIFACT-TYPE

- EXIT-E-CANDIDATE-REPORT
- CANDIDATE-PAYLOAD-MANIFEST
- EXIT-E-CONTENT-READINESS-REPORT
- SCOPE-CERTIFICATE
- OUTER-BUNDLE-MANIFEST

### Enum: HASH-DOMAIN-PREFIX

- ENVELOPE
- PACKAGE-MEMBER
- EVIDENCE-SET

## 4.2. Atomic unit rule

Each promotable forensic record and each independently meaningful synthesis record describes exactly one independently evidenced unit; purely navigational grouping containers are marked GROUPING-CONTAINER, and serialized containers never represent a whole export as one component.

### Rule: atomic-unit-principle

Each promotable forensic record and each independently meaningful synthesis record MUST describe exactly one independently evidenced unit: one function, method, route, query, job, workflow node, graph edge, action, binding, DB routine, constraint, widget behavior, interface operation, architecture module/bounded context, or equivalent. A purely navigational grouping container, for example a family index, class overview, workflow overview, capability index, module index, or use-case index, MUST be marked RECORD-KIND: GROUPING-CONTAINER. It has only a title, navigation metadata, and an explicit member-ID list; it MUST NOT carry independent semantic claims, STAGING-STATUS, COVERAGE-STATUS, blueprint/handbook status, BLOCK-CONFIDENCE-RANK, BLOCK-EVIDENCE-PROFILE, semantic record version/fingerprint, or confirmation, and MUST NOT inherit any such value from children.

### Rule: module-and-container-semantics

A module or bounded context that states responsibility, public interfaces, encapsulation boundary, dependencies, invariants, or data passed is not a grouping container. It is an atomic semantic MOD block governed by the common synthesis header, material CLMs, reciprocal references, coverage, semantic version/fingerprint, and confirmation. Module indexes MAY remain purely navigational grouping containers. A CAP is the governed capability container defined in §2.6 and §3.2: it may carry its own claim-backed governance metadata, but it remains non-promotable and never inherits member evidence or coverage. Serialized files are containers. A whole export MUST NOT be represented as one component. A container SRC is C-COVERED only when its child-enumeration claim is complete for the pinned snapshot and every child coverage row is terminally covered or exactly excluded; it never inherits a child evidence rank/profile and never requires an aggregate promotable CMP. A container with no children requires an evidence-backed zero-child enumeration claim.

## 4.3. Source inventory denominator — `10-SOURCE-INVENTORY.md`

One record exists per concrete source unit.

### Rule: source-versioning

Each source fingerprint/snapshot is an immutable SRC version. A changed source creates a new version linked by supersession; it does not overwrite prior snapshot state. Scope, disposition, and coverage authority live in the scope-specific rows in 16-SOURCE-COVERAGE.md; inventory summaries are derived. A Candidate scope row is temporary and included in closure until changed by an approved scope decision. Approved-Excluded requires exact units, rationale, risk, approver, date, and DEC ID. Wildcards alone cannot identify exclusions at closure.

### Rule: inventory-absence-proof

Inventory MUST cover, where applicable, the declared inventory kinds. Absence of a kind MUST be proven by an inventory-backed enumerator record specifying inspected roots/catalogs, method, snapshot, result count zero, limitations, and CLM evidence. An empty section or zero rows without this proof is invalid.

### Enum: DISCOVERY-METHOD

- filesystem
- AST
- route registry
- DB catalog
- export explosion
- runtime
- operator-approved manifest
- other

### Enum: COVERAGE-SUMMARY

- C-COVERED
- C-PARTIAL
- C-GAP
- C-EXCLUDED
- Mixed

### Enum: INVENTORY-KIND

- files
- symbols
- routes
- RPC handlers
- jobs
- queries
- tables
- views
- triggers
- stored functions/procedures
- generated columns
- constraints
- configuration keys and feature flags
- workflow containers
- nodes
- edges
- pages
- widgets
- actions
- bindings
- queues/events
- files/storage contracts
- external interfaces

### Field: SOURCE-INVENTORY-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| kind | string | true |
| concrete-coordinate | string | true |
| parent-container-src | string | true |
| source-fingerprint-snapshot-src-version | string | true |
| discovery-method | DISCOVERY-METHOD | true |
| scope-summary | string | true |
| environment-runtime-state | string | true |
| traversal-tracks | string | true |
| canonical-owner | string | true |
| disposition-summary | string | true |
| mappings | string | true |
| frontier-references | string | true |
| coverage-rows | string | true |
| coverage-summary | COVERAGE-SUMMARY | true |
| claims | string | true |

## 4.4. Traversal frontier — `11-TRAVERSAL-FRONTIER.md`

Every discovered traversal boundary MUST create exactly one FRT record, including a boundary recognized as terminal, duplicate, excluded, or already mapped at discovery time. Immediate recognition MAY create the record directly in its terminal state, but MUST NOT omit it. Every discovered but not yet terminal traversal boundary remains open until terminally disposed.

### Rule: frontier-terminality

Depth bounds create open frontier records; they do not silently terminate traversal. Every frontier MUST be terminal for Exit A. Terminal-Human-Blocked remains a coverage gap or partial unless the underlying unit is separately approved as excluded.

### Enum: FRONTIER-DISCOVERY-METHOD

- call
- import
- route
- graph edge
- query reference
- DB firing
- config gate
- runtime
- other

### Enum: FRONTIER-STATE

- Open
- Traversed
- Terminal-Mapped
- Terminal-Excluded
- Terminal-Duplicate
- Terminal-Human-Blocked
- Superseded

### Field: FRONTIER-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| origin-invocation-cluster-pov-track | string | true |
| from-src-cmp | string | true |
| target-concrete-coordinate-or-resolution-rule | string | true |
| target-src | string | true |
| discovery-method | FRONTIER-DISCOVERY-METHOD | true |
| depth | int | true |
| snapshot-fingerprint | string | true |
| state | FRONTIER-STATE | true |
| disposition-evidence | string | true |
| child-frontier-ids | string | true |
| closing-invocation | string | true |

## 4.5. Scope-specific source coverage — `16-SOURCE-COVERAGE.md`

One row exists per immutable SRC version by persona/owner by track by environment by snapshot scope. Rows are disjoint within a subgroup; a source participating in several tracks has separate rows and is counted once in each declared track report, while the aggregate unique-SRC report deduplicates by SRC version.

### Rule: coverage-row-identity

The COV ID is deterministically derived from the complete immutable scope tuple and MUST be unique for it; duplicate tuple rows are invalid. The row, not the inventory summary, is the arithmetic input.

### Enum: COVERAGE-SCOPE

- In-Scope
- Candidate
- Approved-Excluded

### Enum: COVERAGE-DISPOSITION

- Unmapped
- Mapped
- Superseded
- Retired
- Approved-Excluded
- Human-Blocked

### Field: SOURCE-COVERAGE-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| cov-id | string | true |
| src-version | string | true |
| persona-owner | string | true |
| track | string | true |
| environment | string | true |
| snapshot | string | true |
| scope | COVERAGE-SCOPE | true |
| disposition | COVERAGE-DISPOSITION | true |
| c-status | string | true |
| exclusion-dec | string | true |
| mappings | string | true |
| frontier-ticket-claim-evidence | string | true |

## 4.6. Stable supersession and tombstones

Records are never silently deleted. A renamed record retains its ID when identity is proven.

### Rule: tombstone-rules

A split, merge, or semantic replacement uses the declared supersession block. Tombstones preserve old IDs, incoming references, and snapshot lineage. Validators reject references to tombstones unless the referencing block explicitly records historical use.

### Field: SUPERSESSION-TOMBSTONE

| Field | Type | Required |
| :-- | :-- | :-- |
| record-version | int | true |
| supersedes | string | true |
| superseded-by | string | true |
| tombstone | string | true |
| tombstone-reason-decision | string | true |

## 5.1. Dedicated prefix groups

Finite protocol enums MUST be validated from the normative Markdown AST using the complete schema/context-qualified field path, never an unanchored text match, a label-only lookup, or a token-prefix guess. Prefix-family tokens must be bracketed where declared, and every finite-enum path is declared here.

### Rule: finite-enum-registry

The exhaustive registry of all normative finite-choice scalar and table-cell paths is declared as ENUM domains in this section. Aliases, invented labels, inherited enums, and applying one schema's label map to another schema are forbidden unless explicitly introduced by a protocol/schema version. An undeclared finite-enum path is a validation error rather than permission to infer an enum from its label or token.

### Rule: traversal-matrix-status

EntryCluster.Required Traversal Matrix[].Status accepts N/A if and only if the same row's Applicability is exactly Not-Applicable; such a row MUST carry the claim-backed reason required by §3.1. A Required row MUST use one bracketed allowed R- token and MUST NOT use N/A. A Not-Applicable row MUST use exactly N/A and MUST NOT use an R- token.

### Enum: PREFIX-STAGING

- [S-DISCOVERED]
- [S-STAGING-VERIFIED]
- [S-PROMOTED]
- [S-DEAD-CODE]
- [S-RETIRED]

### Enum: PREFIX-TICKETS

- [T-OPEN]
- [T-PENDING]
- [T-FOLLOW-UP]
- [T-PROBE-REQUIRED]
- [T-COMPLETE]
- [T-INVALID]
- [T-SUPERSEDED]
- [T-HUMAN-REQUIRED]
- [T-OUT-OF-SCOPE]

### Enum: PREFIX-EVIDENCE

- [E-DIRECT]
- [E-RUNTIME-OBSERVED]
- [E-CROSS-VALIDATED]
- [E-INFERRED]
- [E-UNKNOWN]

### Enum: PREFIX-ACQUISITION

- [A-REQUESTED]
- [A-RECEIVED]
- [A-EXPLODED]
- [A-BLOCKED-MISSING-REFERENCE]
- [A-STRUCTURALLY-COMPLETE]

### Enum: PREFIX-REGISTRY-TRAVERSAL

- [R-ACTIVE]
- [R-UNSWEPT]
- [R-IN-PROGRESS]
- [R-SWEPT]
- [R-PENDING-HUMAN-REVIEW]
- [R-PENDING-PERSONA-ASSIGNMENT]
- [R-DEPROMOTION-PENDING]

### Enum: PREFIX-BLUEPRINT

- [B-DRAFT]
- [B-POPULATED]
- [B-CONFIRMED]
- [B-STALE]

### Enum: PREFIX-COVERAGE

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]

### Enum: REGISTRY-PREFLIGHT-CURRENT-ITERATION

- ALFA
- BRAVO
- CHARLIE
- DELTA
- ECHO
- FOXTROT
- GOLF
- HOTEL
- INDIA
- JULIETT
- KILO
- LIMA
- MIKE
- NOVEMBER
- OSCAR
- PAPA
- QUEBEC
- ROMEO
- SIERRA
- TANGO
- UNIFORM
- VICTOR
- WHISKEY
- X-RAY
- YANKEE
- ZULU

### Enum: REGISTRY-PERSONAREGISTRY-REGISTRY-STATE

- Active
- Conditional
- Disabled
- Historical
- Retired

### Enum: REGISTRY-CAPABILITYREGISTRY-CAPABILITY-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-CAPABILITYREGISTRY-RUNTIME-STATE-MATRIX-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-ENTRYCLUSTER-PRIMARY-POV

- POV-1
- POV-2
- POV-3
- POV-4
- POV-5

### Enum: REGISTRY-ENTRYCLUSTER-TRAVERSAL-TRACK

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-ENTRYCLUSTER-TRAVERSAL-STATUS

- [R-ACTIVE]
- [R-UNSWEPT]
- [R-IN-PROGRESS]
- [R-SWEPT]
- [R-PENDING-HUMAN-REVIEW]
- [R-PENDING-PERSONA-ASSIGNMENT]
- [R-DEPROMOTION-PENDING]

### Enum: REGISTRY-ENTRYCLUSTER-REQUIRED-TRAVERSAL-MATRIX-TRACK

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-ENTRYCLUSTER-REQUIRED-TRAVERSAL-MATRIX-POV

- POV-1
- POV-2
- POV-3
- POV-4
- POV-5

### Enum: REGISTRY-ENTRYCLUSTER-REQUIRED-TRAVERSAL-MATRIX-APPLICABILITY

- Required
- Not-Applicable

### Enum: REGISTRY-ENTRYCLUSTER-REQUIRED-TRAVERSAL-MATRIX-STATUS

- [R-UNSWEPT]
- [R-IN-PROGRESS]
- [R-SWEPT]
- N/A

### Enum: REGISTRY-UNMAPPEDDISCOVERY-TRAVERSAL-STATUS

- [R-PENDING-HUMAN-REVIEW]
- [R-PENDING-PERSONA-ASSIGNMENT]

### Enum: REGISTRY-UNMAPPEDDISCOVERYBUFFER-STATE

- Pending-Merge
- Merged
- Superseded

### Enum: REGISTRY-EXPORT-PLATFORM

- n8n
- Appsmith
- Other

### Enum: REGISTRY-EXPORT-ACQUISITION-STATUS

- [A-REQUESTED]
- [A-RECEIVED]
- [A-EXPLODED]
- [A-BLOCKED-MISSING-REFERENCE]
- [A-STRUCTURALLY-COMPLETE]

### Enum: REGISTRY-PACKAGINGARTIFACT-ARTIFACT-TYPE

- EXIT-E-CANDIDATE-REPORT
- CANDIDATE-PAYLOAD-MANIFEST
- EXIT-E-CONTENT-READINESS-REPORT
- SCOPE-CERTIFICATE
- OUTER-BUNDLE-MANIFEST

### Enum: REGISTRY-SOURCEINVENTORY-KIND

- FILE
- SYM
- ROUTE
- RPC
- JOB
- QUERY
- TABLE
- VIEW
- DR
- CONSTRAINT
- GENERATED
- CFG
- WORKFLOW
- NODE
- EDGE
- PAGE
- WIDGET
- ACTION
- BINDING
- QUEUE
- EVENT
- STORAGE
- IFACE
- OTHER

### Enum: REGISTRY-SOURCEINVENTORY-DISCOVERY-METHOD

- filesystem
- AST
- route registry
- DB catalog
- export explosion
- runtime
- operator-approved manifest
- other

### Enum: REGISTRY-SOURCEINVENTORY-TRAVERSAL-TRACK-S

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-SOURCEINVENTORY-COVERAGE-SUMMARY

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]
- Mixed

### Enum: REGISTRY-FRONTIER-DISCOVERY-METHOD

- call
- import
- route
- graph edge
- query reference
- DB firing
- config gate
- runtime
- other

### Enum: REGISTRY-FRONTIER-STATE

- Open
- Traversed
- Terminal-Mapped
- Terminal-Excluded
- Terminal-Duplicate
- Terminal-Human-Blocked
- Superseded

### Enum: REGISTRY-SOURCECOVERAGE-TRACK

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-SOURCECOVERAGE-SCOPE

- In-Scope
- Candidate
- Approved-Excluded

### Enum: REGISTRY-SOURCECOVERAGE-DISPOSITION

- Unmapped
- Mapped
- Superseded
- Retired
- Approved-Excluded
- Human-Blocked

### Enum: REGISTRY-SOURCECOVERAGE-C-STATUS

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]

### Enum: REGISTRY-RECORDSUPERSESSION-TOMBSTONE

- yes
- no

### Enum: REGISTRY-GROUPINGCONTAINER-RECORD-KIND

- GROUPING-CONTAINER

### Enum: REGISTRY-CAPABILITYRECORD-RECORD-KIND

- GOVERNED-CAPABILITY-CONTAINER

### Enum: REGISTRY-ATOMICCOMPONENT-RECORD-KIND

- ATOMIC

### Enum: REGISTRY-ATOMICCOMPONENT-POV-OWNER

- POV-1
- POV-2
- POV-3
- POV-4
- POV-5

### Enum: REGISTRY-ATOMICCOMPONENT-STAGING-STATUS

- [S-DISCOVERED]
- [S-STAGING-VERIFIED]
- [S-PROMOTED]
- [S-DEAD-CODE]
- [S-RETIRED]

### Enum: REGISTRY-ATOMICCOMPONENT-COVERAGE-STATUS

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]

### Enum: REGISTRY-ATOMICCOMPONENT-TRAVERSAL-TRACKS

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-ATOMICCOMPONENT-CAPABILITY-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-ATOMICCOMPONENT-RUNTIME-STATE-MATRIX-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-ATOMICCOMPONENT-BLOCK-CONFIDENCE-RANK

- None
- Low
- Medium
- High

### Enum: REGISTRY-ATOMICCOMPONENT-BLOCK-EVIDENCE-PROFILE-EVIDENCE-LEVEL

- [E-DIRECT]
- [E-RUNTIME-OBSERVED]
- [E-CROSS-VALIDATED]
- [E-INFERRED]
- [E-UNKNOWN]

### Enum: REGISTRY-SHAREDREFERENCE-TRACK

- Normal
- Shadow/Conditional
- Disabled

### Enum: REGISTRY-SHAREDREFERENCE-CAPABILITY-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-SHAREDREFERENCE-CANONICAL-STAGING-STATUS

- [S-DISCOVERED]
- [S-STAGING-VERIFIED]

### Enum: REGISTRY-CLAIM-MATERIALITY

- Material
- Supporting

### Enum: REGISTRY-CLAIM-EVIDENCE-LEVEL

- [E-DIRECT]
- [E-RUNTIME-OBSERVED]
- [E-CROSS-VALIDATED]
- [E-INFERRED]
- [E-UNKNOWN]

### Enum: REGISTRY-CLAIM-PROJECTION-CLASSIFICATION

- raw-visible
- preserved
- normalized-no-semantic-change
- aliased
- redacted
- removed
- unclassified

### Enum: REGISTRY-CONTRADICTION-RESOLUTION-STATE

- Open
- Explained-By-Scope
- Resolved-By-Evidence
- Decision-Required
- Superseded

### Enum: REGISTRY-DECISION-TYPE

- scope
- contradiction
- target
- migration
- capability

### Enum: REGISTRY-DECISION-APPROVAL-STATUS

- Proposed
- Approved
- Superseded

### Enum: REGISTRY-CAPABILITYRECORD-TARGET-DECISION-STATE

- Bound
- Pending

### Enum: REGISTRY-CONFIRMATION-HUMAN-HATCH-ACTION

- Confirmation

### Enum: REGISTRY-CONFIRMATION-HANDBOOK-READABILITY-REVIEW

- Pass
- Fail
- Not-Applicable

### Enum: REGISTRY-CONFIRMATION-INVALIDATED-BY

- semantic change
- dependency change
- decision supersession
- None

### Enum: REGISTRY-ACQUISITIONCANDIDATES-STATE

- Included
- Rejected-With-Operator-Decision
- Approved-Excluded
- Pending

### Enum: REGISTRY-EXPORTRECONCILIATION-RECONCILIATION-STATE

- Mapped-Atomic
- Mapped-Container-Only
- Approved-Excluded
- Gap
- Superseded

### Enum: REGISTRY-INVOCATIONHEADER-CURRENT-ITERATION

- ALFA
- BRAVO
- CHARLIE
- DELTA
- ECHO
- FOXTROT
- GOLF
- HOTEL
- INDIA
- JULIETT
- KILO
- LIMA
- MIKE
- NOVEMBER
- OSCAR
- PAPA
- QUEBEC
- ROMEO
- SIERRA
- TANGO
- UNIFORM
- VICTOR
- WHISKEY
- X-RAY
- YANKEE
- ZULU

### Enum: REGISTRY-INVOCATIONHEADER-INVOCATION-MODE

- Preflight
- Export Acquisition
- Discovery
- Ticket Resolution
- Promotion Review
- Sequential Reconciliation
- Profile Synchronization
- Cross-Reference-Reconciliation
- Partial-Synthesis
- Final-Synthesis
- Reconstruction-Handoff
- Human Hatch
- Validation/Gate

### Enum: REGISTRY-INVOCATIONHEADER-HUMAN-HATCH-ACTION

- Ticket Escalation
- Decision/Scope Approval
- Confirmation
- None

### Enum: REGISTRY-INVOCATIONHEADER-TRAVERSAL-TRACK

- Normal
- Shadow/Conditional
- Disabled
- None

### Enum: REGISTRY-INVOCATIONHEADER-ACTIVE-POV

- POV-1
- POV-2
- POV-3
- POV-4
- POV-5
- None

### Enum: REGISTRY-COLDRESUME-RUNTIME-SCOPE-MODE-WRITE-RIGHT-CHECK

- Pass
- Fail

### Enum: REGISTRY-COLDRESUME-RUNTIME-STALE-CHECK

- Pass
- Fail

### Enum: REGISTRY-COLDRESUME-COLD-RESUME-RESULT

- Pass
- Fail

### Enum: REGISTRY-TICKET-STATUS

- [T-OPEN]
- [T-PENDING]
- [T-FOLLOW-UP]
- [T-PROBE-REQUIRED]
- [T-COMPLETE]
- [T-INVALID]
- [T-SUPERSEDED]
- [T-HUMAN-REQUIRED]
- [T-OUT-OF-SCOPE]

### Enum: REGISTRY-TICKET-ESCALATION-REASON

- None
- RUNTIME-REQUIRED
- THIRD-PARTY-OPAQUE
- UNAUTHORIZED-ACCESS
- UNSAFE
- DESTRUCTIVE
- NO-SANDBOX
- NO-AUTHORITY-TO-DECIDE
- DOMAIN-MEANING-REQUIRED
- EXTERNAL-FACT-UNAVAILABLE

### Enum: REGISTRY-PROBE-SAFETY-CLASSIFICATION

- non-destructive
- controlled mutation

### Enum: REGISTRY-SYNTHESISGAP-MATERIALITY

- Material
- Supporting

### Enum: REGISTRY-SYNTHESISGAP-DISPOSITION

- Pending
- Ticket-Created
- Frontier-Created
- Both-Created
- Merged
- Invalid
- Approved-Excluded

### Enum: REGISTRY-SYNTHESISBLOCK-COVERAGE-STATUS

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]

### Enum: REGISTRY-SYNTHESISBLOCK-BLOCK-CONFIDENCE-RANK

- None
- Low
- Medium
- High

### Enum: REGISTRY-SYNTHESISBLOCK-BLOCK-EVIDENCE-PROFILE-EVIDENCE-LEVEL

- [E-DIRECT]
- [E-RUNTIME-OBSERVED]
- [E-CROSS-VALIDATED]
- [E-INFERRED]
- [E-UNKNOWN]

### Enum: REGISTRY-SYNTHESISBLOCK-BLUEPRINT-STATUS

- [B-DRAFT]
- [B-POPULATED]
- [B-CONFIRMED]
- [B-STALE]

### Enum: REGISTRY-PERSONAPROFILE-BLOCK-CONFIDENCE-RANK

- None
- Low
- Medium
- High

### Enum: REGISTRY-PERSONAPROFILE-BLOCK-EVIDENCE-PROFILE-EVIDENCE-LEVEL

- [E-DIRECT]
- [E-RUNTIME-OBSERVED]
- [E-CROSS-VALIDATED]
- [E-INFERRED]
- [E-UNKNOWN]

### Enum: REGISTRY-PERSONAPROFILE-BLUEPRINT-STATUS

- [B-DRAFT]
- [B-POPULATED]
- [B-CONFIRMED]
- [B-STALE]

### Enum: REGISTRY-HANDBOOKSECTION-HANDBOOK-STATUS

- [B-DRAFT]
- [B-POPULATED]
- [B-CONFIRMED]
- [B-STALE]

### Enum: REGISTRY-ENTITY-STORAGE

- table
- view
- collection
- cache-key-space
- file
- in-memory

### Enum: REGISTRY-RELATIONSHIP-CARDINALITY

- 1:1
- 1:N
- M:N

### Enum: REGISTRY-RELATIONSHIP-ENFORCEMENT

- DB-FK
- application
- mixed
- none

### Enum: REGISTRY-RELATIONSHIP-DELETE-UPDATE-SEMANTICS

- cascade
- restrict
- nullify
- custom

### Enum: REGISTRY-DATABASEROUTINE-KIND

- trigger
- function
- procedure
- rule
- generated-column
- materialized-refresh
- constraint-helper

### Enum: REGISTRY-DATABASEROUTINE-TARGET-DECISION

- Retain
- Lift
- Redesign
- Retire
- Undecided

### Enum: REGISTRY-BUSINESSRULE-VIOLATION-POLICY

- reject
- warn
- clamp
- default
- compensate
- rollback

### Enum: REGISTRY-USECASE-CAPABILITY-CLASSIFICATION

- Active
- Conditional
- Shadow
- Disabled
- Retired

### Enum: REGISTRY-USECASE-TARGET-DECISION

- Preserve
- Change-By-Approved-Migration
- Redesign
- Retire

### Enum: REGISTRY-BUSINESSRULE-CAPABILITY-STATE

- Active
- Disabled

### Enum: REGISTRY-SCHEDULEDJOB-CAPABILITY-STATE

- Active
- Conditional
- Shadow
- Disabled
- Retired
- Unknown

### Enum: REGISTRY-DEPLOYMENT-SINGLE-POINT-OF-FAILURE

- yes
- no
- unknown

### Enum: REGISTRY-CONFIGURATION-MEDIUM

- env
- DB
- flag service
- file
- build-time

### Enum: REGISTRY-NFR-ATTRIBUTE

- latency
- throughput
- availability
- RPO
- RTO
- scalability
- operability
- other

### Enum: REGISTRY-SECURITY-CATEGORY

- auth
- authorization
- transport
- at-rest
- in-transit
- PII
- secrets
- audit
- threat
- privacy

### Enum: REGISTRY-SECURITY-RISK-LEVEL

- high
- medium
- low
- unassessed

### Enum: REGISTRY-TRACEABILITY-C-STATUS

- [C-COVERED]
- [C-PARTIAL]
- [C-GAP]
- [C-EXCLUDED]

### Enum: REGISTRY-TRACEABILITY-LIFECYCLE-ELIGIBILITY-CLASSIFICATION

- Eligible
- Ineligible
- Unknown

### Enum: REGISTRY-VALIDATIONSUMMARY-EVIDENCE-BINDINGS-BINDING-KIND

- RAW-SOURCE-BYTES
- RECORD-PAYLOAD
- SEMANTIC-CONTENT
- DECISION-CONTENT
- CNF-PAYLOAD
- ARTIFACT-PAYLOAD
- CANONICAL-ENVELOPE
- FILE-TRANSPORT
- PACKAGE-MEMBER-FILE
- VALIDATION-SUMMARY
- DENOMINATOR-QUERY
- PROTOCOL-REGISTRY

### Enum: REGISTRY-EXITECANDIDATEREPORT-CHECK-RESULTS-RESULT

- Pass
- Fail

### Enum: REGISTRY-EXITECANDIDATEREPORT-CANDIDATE-RESULT

- Ready-For-Human-Review
- Blocked

### Enum: REGISTRY-EXITECONTENTREADINESSREPORT-CHECK-RESULTS-RESULT

- Pass
- Fail

### Enum: REGISTRY-EXITECONTENTREADINESSREPORT-EXIT-E-PREPACKAGE-STATE

- Pending

### Enum: REGISTRY-EXITECONTENTREADINESSREPORT-CONTENT-READINESS-RESULT

- Ready-For-Certificate
- Blocked

### Enum: REGISTRY-OUTERBUNDLEMANIFEST-EXIT-E-STATUS

- Passed

### Enum: REGISTRY-FINALVERIFICATIONRECEIPT-CHECK-RESULTS-RESULT

- Pass
- Fail

### Enum: REGISTRY-FINALVERIFICATIONRECEIPT-OVERALL-RESULT

- Pass
- Fail

## 5.2. Claim-level evidence — `12-CLAIM-EVIDENCE.md`

Every claim record binds its subject, field or step, evidence level, anchor, and lineage so that materiality and projection classification are explicit.

### Field: CLAIM

| Field | Type | Required |
| :-- | :-- | :-- |
| subject-id-record-version | string | true |
| field-step | string | true |
| claim | string | true |
| materiality | REGISTRY-CLAIM-MATERIALITY | true |
| evidence-level | REGISTRY-CLAIM-EVIDENCE-LEVEL | true |
| anchor | string | true |
| source-fingerprint-snapshot | string | true |
| evidence-route-lineage-id | string | true |
| projection-classification | REGISTRY-CLAIM-PROJECTION-CLASSIFICATION | true |
| independent-supporting-claims | set<string> | false |
| independence-basis | string | true |
| contradiction-ids | set<string> | false |
| created-last-revalidated-by | string | true |

### Field: BLOCK-EVIDENCE-SUMMARY

| Field | Type | Required |
| :-- | :-- | :-- |
| block-confidence-rank | REGISTRY-ATOMICCOMPONENT-BLOCK-CONFIDENCE-RANK | true |
| block-evidence-profile | string | true |

## 5.3. Sanitized projections and helper limitations

A sanitized export produced by an approved local helper is an agent-visible evidence projection.

### Rule: sanitized-projections-rules

The projection MUST include raw-source fingerprint lineage and a sanitization report classifying preserved, normalized, aliased, redacted, removed, and unclassified content. E-DIRECT is permitted only for preserved or semantics-preservingly normalized behavior. Aliases prove reference existence, not secret value, connectivity, authorization, or runtime validity. Behavior dependent on redacted, removed, or unclassified content remains E-INFERRED or E-UNKNOWN and creates a ticket when material. Entity-correlation output is a candidate-discovery aid only. It MUST NOT become evidence, a call edge, ownership proof, or automatic inclusion without a direct structural signal or explicit human scope decision. Private matching hints and private entity-usage indexes MUST NOT enter agent-visible inventory.

## 5.4. Contradictions — `14-CONTRADICTIONS.md`

Every contradiction record names the conflicting claims, scope, and stale impact.

### Field: CONTRADICTION

| Field | Type | Required |
| :-- | :-- | :-- |
| claims-in-conflict | set<string> | true |
| scope-environments-snapshots | string | true |
| conflict-description | string | true |
| resolution-state | REGISTRY-CONTRADICTION-RESOLUTION-STATE | true |
| resolution-decision-id | string | false |
| affected-records-and-stale-impact | string | true |

## 5.5. Authoritative decisions — `15-DECISIONS.md`

Every decision record carries its content version and fingerprint, exact scope and denominator effect, approval status, and approval envelope bindings.

### Field: DECISION

| Field | Type | Required |
| :-- | :-- | :-- |
| decision-content-version | int | true |
| decision-content-fingerprint | hash | true |
| type | REGISTRY-DECISION-TYPE | true |
| observed-inputs | string | true |
| options-decision-rationale-risks | string | true |
| exact-scope-and-denominator-effect | string | true |
| affected-records-rollout-rollback | string | true |
| content-supersession | string | true |
| approval-status | REGISTRY-DECISION-APPROVAL-STATUS | true |
| approval-envelope-version | int | true |
| owner-approver-authority-basis-approval-timestamp | string | true |
| approved-content-binding | string | true |
| approval-signature-audit | string | true |

## 5.6. Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md`

Every helper-published candidate and operator action receives a row in 17-ACQUISITION-CANDIDATES.md and a deterministic CND ID whose canonical key binds the stable candidate or action identity, helper policy version, and pinned snapshot.

### Rule: candidate-reconciliation

Duplicate candidate or action identities and reused CND IDs are invalid. States are Included, Rejected-With-Operator-Decision, Approved-Excluded, or Pending. The helper manifest MUST attest that all accessible workflow metadata was scanned for the pinned snapshot, all direct structural matches were included, all high-confidence entity-only matches were human-reviewed, and all configured medium-confidence candidates received an operator decision; scoring thresholds and policy versions are recorded. Lower-confidence omissions follow the configured, approved scope policy. Helper publication completeness means safe approved outputs were produced; protocol A-STRUCTURALLY-COMPLETE additionally requires every candidate, operator action, and dynamic reference row terminal. Correlation remains discovery and scope input only, never behavioral evidence or an automatic edge.

## 5.7. Human confirmations — `18-CONFIRMATIONS.md`

A CNF record confirms exact typed record IDs, semantic versions, and fingerprints against a pinned closed snapshot, never mutable whole-file bytes.

### Field: CONFIRMATION

| Field | Type | Required |
| :-- | :-- | :-- |
| human-hatch-action | REGISTRY-CONFIRMATION-HUMAN-HATCH-ACTION | true |
| approver-role-authority-basis | string | true |
| review-timestamp-review-scope-limitations | string | true |
| handbook-readability-review | REGISTRY-CONFIRMATION-HANDBOOK-READABILITY-REVIEW | true |
| confirmed-record-bindings | set<string> | true |
| pinned-exit-a-exit-e-candidate-report-candidate-payload-manifest-fingerprints | string | true |
| required-dec-content-versions-fingerprints-and-approval-envelopes-reviewed | string | true |
| invalidated-by | REGISTRY-CONFIRMATION-INVALIDATED-BY | true |
| cnf-signature-envelope-version | int | true |
| cnf-payload-fingerprint | hash | true |
| signature-or-approval-reference-signature-timestamp | string | true |

## 6.1. Atomic component schema

The atomic component record schema is declared below.

### Rule: atomic-unit-specificity

All fields are atomic-unit-specific. A routine block MUST NOT aggregate unrelated methods, branches, graph nodes, or queries merely because they share a module.

### Enum: BLOCK-CONFIDENCE-RANK

- None
- Low
- Medium
- High

### Field: ATOMIC-COMPONENT

| Field | Type | Required |
| :-- | :-- | :-- |
| record-kind | string | true |
| source-ids | string | true |
| file-path-line-range | string | true |
| virtual-coordinate | string | true |
| pov-owner | REGISTRY-ENTRYCLUSTER-PRIMARY-POV | true |
| general-intent | string | true |
| inputs-outputs | string | true |
| auth-binding | string | true |
| staging-status | string | true |
| coverage-status | string | true |
| personas | string | true |
| traversal-tracks | string | true |
| capability-state | REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE | true |
| runtime-state-matrix | string | true |
| depends-on-depended-on-by | string | true |
| synthesis-references | string | true |
| side-effects | string | true |
| exceptional-behaviors | string | true |
| edge-and-odd-cases | string | true |
| behavioral-notes-for-reconstruction | string | true |
| material-claims | string | true |
| block-confidence-rank | BLOCK-CONFIDENCE-RANK | true |
| block-evidence-profile | string | true |
| evidence-anchors | string | true |
| record-version-supersession | string | true |

## 6.2. Dual-entry discovery

Persona-driven discovery is the only behavioral discovery path. There is no independent shared-code sweep.

### Rule: persona-driven-discovery

When a persona entry discovers a component: create the full canonical atomic record in the persona's SHARED-STAGING-BUFFER.md for later merge to personas/_shared/XX-*-POV.md with S-DISCOVERED; then immediately create a persona invocation pointer in personas/{persona-directory}/XX-*-POV.md using the shared reference schema. This ensures a persona manifest represents its full operational footprint without duplicating canonical behavior. If discovery reaches a coordinate that cannot yet be assigned to the bound persona/cluster/track/POV, the same transaction MUST create its FRT and append a DSC entry to personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md; the invocation MUST NOT write 0A directly or continue through the unmapped boundary. Sequential Reconciliation MUST merge all such pending entries before any later Discovery invocation begins, as specified in §3.1.

### Field: SHARED-REFERENCE

| Field | Type | Required |
| :-- | :-- | :-- |
| canonical-location | string | true |
| local-trigger-src-coordinate | string | true |
| local-invocation-signature | string | true |
| track | string | true |
| capability-state | REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE | true |
| canonical-staging-status | string | true |

## 6.3. Staging and promotion

A shared atomic component advances from S-DISCOVERED to S-STAGING-VERIFIED only when mandatory fields, claim evidence, source mappings, references, and ticket state validate.

### Rule: promotion-locks

Promotion eligibility requires all locks: usage uniqueness, meaning exactly one persona considering Normal, Shadow/Conditional, Disabled, and historical/intended usage within scope; ambiguity, meaning no active ticket or open material contradiction references the component; dependency, meaning no canonical shared component depends on it; coverage, meaning source units are C-COVERED for the promoted component; and evidence, meaning no material E-UNKNOWN claim and no material unresolved redaction dependency.

### Rule: promotion-execution

Promotion Review validates these locks but MUST NOT move or rewrite canonical blocks. It appends one PROMOTION-REQUEST to the bound persona's personas/{persona-directory}/SHARED-STAGING-BUFFER.md, containing the CMP/cohort IDs, target persona/POV, source _shared paths and complete file/record fingerprints, lock-evidence IDs, review invocation, and requested S-PROMOTED state. Connected atomic components MAY be requested as one cohort only when their combined persona set is singular and all locks pass. Grouping containers never promote. Sequential Reconciliation is the sole promotion executor. It loads the complete request, canonical personas/_shared/ source blocks, target persona POV files, indexes, traceability, and current lock inputs; applies the stale-check guard; revalidates every lock; and atomically moves each full canonical block from _shared to the target persona file, changes it to S-PROMOTED, updates or removes shared-reference pointers as applicable, updates indexes/traceability, consumes the request, and logs one transaction. A failed or stale revalidation performs no partial move and leaves a dispositioned request for renewed review. At no point may both source and target contain canonical copies, nor may neither contain one.

### Rule: dead-code-proof

S-DEAD-CODE requires an inventory-backed zero-caller proof across every traversal track and every environment/snapshot in the component's bound scope denominator. Any linked nonterminal dynamic-reference or acquisition-candidate row, active ticket, open frontier, unresolved possible-caller edge, or discovered active/conditional/shadow/disabled/historical/intended caller makes S-DEAD-CODE invalid. A dynamically observed caller disproves dead-code classification. An unavailable caller domain may be removed from the proof denominator only by an exact approved exclusion; it is never equivalent to proven no caller. These effects are local to the linked component and its dependency closure and MUST NOT downgrade unrelated coverage.

## 6.4. Depromotion

If a promoted component gains another persona caller, a shared dependency, or expanded persona usage, it is depromoted back to personas/_shared/ through Sequential Reconciliation with history preserved.

### Rule: depromotion-steps

If a promoted component gains another persona caller, a shared dependency, or expanded persona usage: during parallel work, record R-DEPROMOTION-PENDING in the discovering persona's staging buffer; Sequential Reconciliation moves the full block to personas/_shared/, resets it to S-STAGING-VERIFIED, merges personas and claims, and replaces persona copies with shared references; and log the event in 0G and update traceability and stale dependencies.

## 6.5. Surgical evolution

POV files are patched surgically: unrelated blocks are never deleted, reordered, or reformatted, ticket resolution cites the ticket, and superseded facts remain traceable through versions and tombstones.

### Rule: surgical-patching

POV files MUST be patched surgically. Unrelated blocks MUST NOT be deleted, reordered, or reformatted. Ticket resolution updates affected assumptions and cites the ticket. Superseded facts remain traceable through versions/tombstones. Clean-sweep rewrites are prohibited except an explicitly approved protocol migration that preserves all IDs and history.

## 7.1. Trust boundary

Live acquisition is an operator-side activity; agents MUST NOT invoke acquisition using production credentials.

### Rule: helper-trust-boundary

The local helper MAY read raw Appsmith exports and approved read-only n8n interfaces, build private entity indexes, sanitize, fingerprint, and publish approved projections. The helper MUST default to dry-run; accept credentials only through environment, stdin, or OS credential storage, never command arguments; keep raw exports and private matching state outside the target repository and authoritative .extracted workspace in operator-controlled storage not published to executor-visible projections; use restrictive temporary storage and best-effort cleanup; perform field-aware sanitization and secret scanning; fail closed for unsupported schema versions, unknown sensitive fields, unresolved high-confidence secret findings, missing lineage, mixed snapshots, or rejected operator review; make no mutation to Appsmith, n8n, or production; make no third-party transmission of exports, secrets, private hints, or repository content; and publish atomically only approved sanitized inventory and bridge fragments. Agents MAY inspect helper source and approved outputs; bridge fragments are proposed patches requiring authorized approval before merge.

## 7.2. Acquisition register

The acquisition register records every export artifact, its projection, and its acquisition status.

### Field: EXPORT

| Field | Type | Required |
| :-- | :-- | :-- |
| platform | REGISTRY-EXPORT-PLATFORM | true |
| stable-artifact-id-display-name | string | true |
| raw-source-alias | string | true |
| agent-visible-export-path | string | true |
| raw-source-fingerprint-projection-fingerprint | string | true |
| snapshot | string | true |
| acquisition-status | REGISTRY-EXPORT-ACQUISITION-STATUS | true |
| referenced-by | set<string> | false |
| outstanding-referenced-artifacts | set<string> | false |
| normalized-map | string | false |
| reconciliation-summary | string | true |
| sanitization-lineage-reports | string | true |
| operator-action-required | string | true |

## 7.3. Logical explosion and virtual coordinates

Before POV extraction, each export is decomposed into atomic logical units whose virtual coordinates are canonical for entry coordinates, evidence, dependencies, frontiers, and tickets.

### Rule: virtual-coordinates

Coordinates name the source export, then the logical kind and stable ID, then the normalized name, using the declared hash and slash separators. Stable internal IDs take precedence. Name-only identities receive deterministic disambiguators and remain low-confidence until verified. Physical container fingerprint and virtual coordinate MUST both be preserved.

## 7.4. Required explosion coverage

Explosion coverage is enumerated per platform and must include every graph edge.

### Rule: explosion-coverage

A complete Appsmith application export may be complete as a physical container, but that does not prove runtime datasource values, plugin behavior, secrets, environment configuration, or external dependencies. Those remain separately evidenced, ticketed, acquired, or exactly excluded. Primary ownership follows the POV definitions. Every graph edge is an atomic POV-5 record, even when its payload meaning creates tickets for another POV.

### Enum: N8N-INVENTORY-KIND

- workflow metadata
- triggers
- webhooks
- schedules
- nodes
- code/function nodes
- transformations
- conditions
- merges
- retries
- error branches
- credentials references
- database/cache/file I/O
- external calls
- sub-workflows
- configured error workflows
- every graph edge

### Enum: APPSMITH-INVENTORY-KIND

- pages
- widgets
- bindings
- JSObjects
- actions
- queries
- datasources
- APIs
- event handlers
- action chains
- validations
- navigation
- permissions
- environment references
- timers
- polling
- deferred behavior

## 7.5. n8n human acquisition loop

Register, explode, trace, and reconcile one approved serialized export through a human acquisition loop, stopping at every missing reference boundary.

### Rule: acquisition-seed

Register available inventory or exports as A-RECEIVED with stable IDs and fingerprints.

### Rule: acquisition-explosion

Generate source inventory, normalized units, entry candidates, and edges, then transition to A-EXPLODED.

### Rule: reference-discovery

Inspect static sub-workflows, error workflows, IDs or names in expressions or code, queue or event links, and webhook relationships.

### Rule: immediate-placeholder

Every missing referenced workflow receives an A-REQUESTED register entry immediately with exact known identity, referring virtual coordinate, and operator action. The referrer becomes A-BLOCKED-MISSING-REFERENCE. Unknown ownership also creates an unmapped discovery. Traversal stops at the boundary; missing behavior is not hypothesized.

### Rule: operator-export

Stage a sanitized projection without raw secrets.

### Rule: cold-resume

Fingerprint, explode, link, restore the referrer to A-EXPLODED, and resume its open frontier.

### Rule: acquisition-closure

Transition to A-STRUCTURALLY-COMPLETE only after static reference closure, 100 percent normalized reconciliation, and terminal disposition of every published helper candidate, operator action, inaccessible artifact, and dynamic reference in 17-ACQUISITION-CANDIDATES.md. A dynamic reference may close only through inclusion or mapping or an exact approved scope decision; recording it as unresolved is not protocol closure.

### Rule: dynamic-references

Dynamic references are recorded as unresolved. After persona or POV assignment, static investigation creates a ticket; they are never guessed.

### Enum: ACQUISITION-STATE

- A-RECEIVED
- A-EXPLODED
- A-REQUESTED
- A-BLOCKED-MISSING-REFERENCE
- A-STRUCTURALLY-COMPLETE

## 7.6. Normalized maps

Normalized maps are navigation-only indexes over one exploded export.

### Rule: normalized-maps-navigation

Normalized maps list metadata, each virtual unit, stable ID, type, primary POV, source fingerprint, edges, external systems, credential aliases, and ambiguities. They are navigation-only indexes. They MUST NOT summarize an export as one routine or bypass persona-driven discovery.

### Field: NORMALIZED-MAP-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| metadata | string | true |
| virtual-unit | id | true |
| stable-id | id | true |
| type | string | true |
| primary-pov | string | true |
| source-fingerprint | hash | true |
| edges | set<string> | false |
| external-systems | set<string> | false |
| credential-aliases | set<string> | false |
| ambiguities | set<string> | false |

### Table: normalized-map

Rows: NORMALIZED-MAP-ROW, key: stable-id.

## 7.7. Export reconciliation — `13-EXPORT-RECONCILIATION.md`

**Artifact:** 13-EXPORT-RECONCILIATION.md

Every normalized export coordinate has exactly one dispositioned row.

### Rule: reconciliation-completeness

Define the exact completion predicate for export reconciliation.

### Enum: RECONCILIATION-STATE

- Mapped-Atomic
- Mapped-Container-Only
- Approved-Excluded
- Gap
- Superseded

### State: reconciliation-complete

- Values: unknown, true, false
- Initial: unknown
- unknown -> true
- unknown -> false
- true -> false
- false -> true

### State: fingerprint-status

- Values: unknown, computed, unsupported
- Initial: unknown
- unknown -> computed
- unknown -> unsupported
- computed -> unsupported
- unsupported -> computed

### Field: RECONCILIATION-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| export-id | EXPORT-ID | true |
| normalized-coordinate | COORDINATE | true |
| src-id | SRC-ID | true |
| atomic-edge-id | CMP-OR-POV5-EDGE-ID | true |
| exclusion-decision | EXCLUSION-DECISION | true |
| reconciliation-state | RECONCILIATION-STATE | true |
| clm | CLM-LIST | false |

### Table: 13-EXPORT-RECONCILIATION

Rows: RECONCILIATION-ROW, key: normalized-coordinate.

## 8.1. Resume identity header

Every invocation begins with the identity header, and the conforming runtime independently validates declared actions, exact read/write scope, loaded fingerprints, checkpoint freshness, and mode rights before commit.

### Rule: header-authorization

The executor MUST declare its intended semantic actions against the identity header before writing. Active POV is POV-1 through POV-5. The conforming runtime MUST independently validate the declared actions, exact read/write scope, loaded fingerprints, checkpoint freshness, and mode rights before commit; an executor's self-assessment cannot authorize mutation.

### Field: INVOCATION-HEADER

| Field | Type | Required |
| :-- | :-- | :-- |
| protocol-version | string | true |
| system-namespace | string | true |
| current-iteration | string | true |
| invocation-id | id | true |
| invocation-mode | string | true |
| human-hatch-action | string | false |
| invocation-scope | string | true |
| persona | string | false |
| persona-prefix | string | false |
| persona-directory | string | false |
| entry-cluster | string | false |
| traversal-track | string | false |
| active-pov | string | false |
| active-pov-file | string | false |
| active-question-ledger | string | false |
| environment-snapshot | string | true |
| max-traversal-depth | int | true |
| allowed-read-targets | set<string> | true |
| allowed-write-targets | set<string> | true |
| forbidden-write-targets | set<string> | true |
| loaded-checkpoint-fingerprint | hash | false |

## 8.2. Strict single-scope modes

Discovery, Ticket Resolution, and Promotion Review bind exactly one persona, prefix, directory, cluster, track, POV, POV file, and at most one question ledger.

### Rule: strict-bindings

Each strict mode has one semantic POV target and optional one semantic ledger target. Its fixed transactional sidecar bundle MAY append only directly produced rows or blocks in 10, 11, 12, 14, 16, persona-local shared buffers, personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md, and 0G; these are assurance and audit side effects, not permission to change scope or another POV's semantics. Cross-POV outputs go to the bound ledger or persona-local buffers. An unmapped discovery MUST use the local buffer and requires Sequential Reconciliation before any later Discovery invocation.

### Rule: wildcard-expansion

A cluster containing wildcards MUST be expanded into concrete coordinates before a closing invocation. One invocation may traverse multiple atomic units reachable from its one concrete root, subject to depth and frontier rules, but may not switch roots, personas, POVs, or tracks.

### Rule: alfa-start-behavior

ALFA is the first bounded discovery iteration. Every POV may log questions immediately; no maturity phase may suppress cross-boundary ambiguity. On large repositories, ALFA MUST still obey concrete cluster, track, persona, POV, and depth isolation. Branches beyond the bound become frontier records and tickets rather than an attempted full-tree sweep.

## 8.3. Invocation-mode enum and explicit multi-file modes

The exhaustive invocation-mode enum is the declared registry; multi-file modes enumerate their exact targets.

```text
Preflight | Export Acquisition | Discovery | Ticket Resolution | Promotion Review | Sequential Reconciliation | Profile Synchronization | Cross-Reference-Reconciliation | Partial-Synthesis | Final-Synthesis | Reconstruction-Handoff | Human Hatch | Validation/Gate
```

### Rule: mode-targets

Preflight has one semantic target (0A); Discovery, Ticket Resolution, and Promotion Review obey §8.2; Validation/Gate writes one report target. All MAY append 0G and fixed assurance sidecars allowed by their contracts. Probe ingestion executes as Ticket Resolution. Neutral arbitration executes as Ticket Resolution with the original persona/cluster/track and a recorded neutral POV assignment. Protocol migration is metadata on a sequence of existing modes, not a separate mode. Human confirmation executes only as Human Hatch: Confirmation or the equivalent explicitly authorized Reconstruction-Handoff confirmation-envelope action and writes 18-CONFIRMATIONS.md plus exact envelope transitions without semantic mutation.

### Rule: multi-file-modes

Only these modes may intentionally span multiple semantic targets. Export Acquisition: acquisition register, inventory/reconciliation/candidate sections, normalized maps, security findings, and 0G. Sequential Reconciliation: canonical personas/_shared/ merges plus persona-local shared-staging, promotion-request, question, unmapped-discovery, and synthesis-gap buffer merges; atomic promotion and depromotion; buffered ticket/frontier creation; deterministic deduplication; indexes; traceability; frontier dispositions; and 0G. Profile Synchronization: exactly one persona profile semantic payload plus its atomic 0A registry/ownership updates and 0G, with an authorized confirmation action later updating only its certification envelope. Cross-Reference-Reconciliation: only already-existing ID reference fields, reciprocal links, traceability, and 0G; no new legacy facts, and changed semantic reference bindings increment semantic versions/fingerprints. Partial-Synthesis: bounded provisional 90-96 blocks and a synthesis gap buffer. Final-Synthesis: closed-corpus 90-96 blocks and synthesis gap buffer. Reconstruction-Handoff: populated handbook payloads, approved-decision projection, authorized confirmation envelopes, equivalence suite, candidate/outer payload manifests, scope certificate, and package assembly in the exact §15.4 stages; it does not write Validation/Gate reports. Human Hatch has three mandatory action subtypes in its identity header: Ticket Escalation, Decision/Scope Approval, and Confirmation, each limited to its declared prerequisites and rights; a subtype cannot borrow another subtype's rights, and Human Hatch writes only the records authorized by its declared subtype plus exact coverage/traceability effects and 0G. Each multi-file mode MUST enumerate exact files and record sections in its identity header; the mode name alone grants no broad write access. Validation/Gate reads authoritative inputs and writes only its named report block plus 0G. Human confirmation writes 18-CONFIRMATIONS.md and exact envelope transitions only, never semantic profile/synthesis/handbook payloads.

## 8.4. Mandatory read sets

All modes load protocol constants/statuses, 0A, 0D, 0E, 0F, 0G or 0H, relevant source inventory/frontier/claim records, and the active snapshot metadata.

### Rule: applicability-provenance

Semantic decisions that a dependency or scope dimension is applicable or not applicable MUST carry the §0.1 classifier provenance and claim/evidence bindings. Unknown applicability is fail-inclusive: the dependency remains in the required read/validation closure and blocks mutation if it cannot be loaded and validated.

### Rule: strict-mode-read-sets

Strict single-scope modes additionally load the active POV file, at most one active ledger, the matching shared POV/ledger for persona scope, the persona's UNMAPPED-DISCOVERY-BUFFER.md, concrete source excerpts, and relevant normalized map. Before Discovery, the invocation MUST verify from a current 0H global buffer summary or a deterministic read of every persona-local unmapped buffer that no Pending-Merge entry exists; an absent/stale summary without the full fallback read blocks the invocation. Load 0B for auth; 0C for sensitive data; dependency source excerpts for cross-boundary references; probe logs for matching probe tickets; persona siblings only for promotion uniqueness.

### Rule: acquisition-read-sets

Acquisition loads 0A, 0C, 0G/0H, supplied approved projections, lineage/sanitization reports, existing normalized maps, source inventory, and export reconciliation.

### Rule: reconciliation-read-sets

Sequential Reconciliation loads all relevant canonical personas/_shared/ files, persona-local shared-staging/promotion-request/question/unmapped-discovery buffers, synthesis-gap buffers, and destination records.

### Rule: profile-closure-read-sets

Profile Synchronization MUST load 0A, the target PRF profile, and every record directly referenced by or applicable to any profile field, then recursively load the complete transitive DERIVED-FROM closure of those records. This closure explicitly includes every applicable AUTH, CLM, SRC, CMP, ER, REL, BR, UC, IF, CFG, SCHED, DR, SM, STATE, DEP, MOD, CAP, DEC, SEC, NFR, and FLT record, plus relevant source inventory, coverage, traceability, contradiction, decision-approval, and snapshot records needed to prove currency. Every loaded dependency MUST have its exact record ID, version where applicable, semantic/content fingerprint, source snapshot, and reciprocal binding verified. Profile Synchronization MUST reject mutation if any directly applicable record or any transitive dependency is absent, unresolved, stale, version-mismatched, fingerprint-mismatched, or read from a mixed snapshot. Confirmation validation of a PRF MUST load and validate this identical closure against the candidate-bound versions and fingerprints before attaching or accepting its envelope.

### Rule: synthesis-read-sets

Partial/Final Synthesis loads required closed or bounded forensic and assurance inputs plus existing synthesis blocks. Final Synthesis MUST load all active personas, all POV/ledger files, all relevant claims, reconciliation, coverage, contradictions, and gate inputs.

## 8.5. Cold resume check

Before mutation, the executor returns the semantic assessment fields and the conforming runtime completes the deterministic validation fields before deciding whether to commit.

### Rule: cold-resume-fingerprint

The cold-resume check fingerprint is sha256 over UTF-8 COLD-RESUME plus the system namespace, one vertical bar, the invocation ID, one LF, and the §4.1.2 canonical serialization of every assessment and validation field except its own carrier Cold-Resume Check Fingerprint.

### Rule: cold-resume-validation

The runtime MUST reproduce identity, scope, loaded-file fingerprints, stale state, and write rights from authoritative filesystem state; it validates but does not invent the executor's semantic blockers, actions, or exclusion reasons. A missing field, mismatched fingerprint, failed runtime check, or semantic exclusion that would omit required scope yields Cold-Resume Result: Fail, rejects semantic mutation, and still appends the required no-mutation 0G entry with the check fingerprint and exact failure. The successful check fingerprint and result are bound into the same transaction's 0G entry, making the pre-mutation authorization externally observable. Filesystem state is authoritative over conversation.

### Enum: COLD-RESUME-FIELD

- Iteration / Invocation / Mode
- Persona / Cluster / Track / POV
- Relevant Active Tickets and Frontiers
- Last Mutation Affecting Targets
- Semantic Blockers and Contradictions
- Proposed Allowed Next Actions
- Explicit Context Exclusions and Reasons
- Runtime-Validated Loaded Files / Fingerprints
- Runtime Scope / Mode / Write-Right Check
- Runtime Stale Check
- Cold-Resume Result
- Cold-Resume Check Fingerprint

## 8.6. Concurrent persona traversal

Personas MAY run concurrently in one named iteration but MUST NOT write canonical shared POV or shared question files directly.

### Rule: concurrency-buffers

Shared component discoveries go to personas/{persona-directory}/SHARED-STAGING-BUFFER.md; outbound shared questions go to personas/{persona-directory}/SHARED-QUESTIONS-BUFFER.md; unmapped coordinates go to personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md with a same-transaction FRT; persona-owned invocation pointers may be written immediately; depromotion is annotated R-DEPROMOTION-PENDING and deferred.

### Rule: sequential-reconciliation-merge

Sequential Reconciliation collects, deterministically sorts, and deduplicates all buffer classes, promotion requests, and synthesis gaps; merges compatible claims without improperly raising BLOCK-CONFIDENCE-RANK or losing profile detail; records contradictions; merges persona sets; commits shared records and tickets; merges unmapped discoveries into 0A; creates required ticket and frontier work from GAP records; executes promotion and depromotion atomically; updates assurance links; and purges consumed buffers atomically. No later Discovery invocation may begin while an unmapped-discovery buffer contains a Pending-Merge entry.

## 8.7. Stale checkpoint guard

Before writing, compare loaded checkpoint and source fingerprints with target record versions.

### Rule: stale-guard

If any target changed after the checkpoint, stop and request a refreshed read set. Changed serialized fingerprints invalidate affected normalized maps, source records, claims, coverage, and downstream blocks before further mutation.

## 8.8. `0G` invocation log

Append the invocation log record after every invocation; 0H MAY compress current active states but MUST include its source 0G range and fingerprint.

### Rule: invocation-log-append

The appended record heading is [ITERATION] - [INVOCATION-ID] - [Scope] - [Mode]. When 0H is used to authorize a later Discovery invocation, it MUST also include a fingerprinted global count and list of every persona-local Pending-Merge unmapped discovery; otherwise the full-buffer fallback read in §8.4 is required.

### Field: INVOCATION-LOG

| Field | Type | Required |
| :-- | :-- | :-- |
| identity-header-summary | string | true |
| cold-resume-result-check-fingerprint | string | true |
| declared-semantic-actions-runtime-authorization-result | string | true |
| loaded-files-fingerprints | string | true |
| mutation-summary | string | true |
| tickets-created-updated | string | true |
| frontiers-created-closed | string | true |
| unmapped-discoveries-gap-records-buffered-merged-or-disposed | string | true |
| exports-acquired-requested | string | true |
| src-cmp-clm-ids-updated | string | true |
| evidence-changes | string | true |
| promotion-depromotion | string | true |
| coverage-reconciliation-changes | string | true |
| stale-propagation | string | true |
| next-required-loads | string | true |

## 9.1. Ticket schema and canonical IDs

Ticket IDs use MONIKER-PERSONA-POV-SEQUENCE; allocation resets per iteration/persona/target POV and is reserved atomically.

### Rule: ticket-identity

A ticket record carries TICKET-ID. Ticket IDs use MONIKER-PERSONA-POV-SEQUENCE. Sequence allocation resets per iteration/persona/target POV and is reserved atomically. For buffered concurrent tickets, Sequential Reconciliation sorts unassigned candidates by target POV, normalized source coordinate, source persona prefix, and SHA-256 of the normalized description, then assigns the next free sequence. A directly authorized single-ledger invocation reserves the next free sequence before writing. The tuple and a stored creation fingerprint MUST be unique; an allocated ID is never reused, even when invalidated or superseded.

### Enum: TICKET-ESCALATION-REASON

- None
- RUNTIME-REQUIRED
- THIRD-PARTY-OPAQUE
- UNAUTHORIZED-ACCESS
- UNSAFE
- DESTRUCTIVE
- NO-SANDBOX
- NO-AUTHORITY-TO-DECIDE
- DOMAIN-MEANING-REQUIRED
- EXTERNAL-FACT-UNAVAILABLE

### State: TICKET-STATE

- Values: T-OPEN, T-PENDING, T-PROBE-REQUIRED, T-COMPLETE, T-FOLLOW-UP, T-INVALID, T-SUPERSEDED, T-OUT-OF-SCOPE, T-HUMAN-REQUIRED
- Initial: T-OPEN
- T-OPEN -> T-PENDING
- T-OPEN -> T-PROBE-REQUIRED
- T-OPEN -> T-HUMAN-REQUIRED
- T-FOLLOW-UP -> T-PROBE-REQUIRED
- T-FOLLOW-UP -> T-OPEN
- T-FOLLOW-UP -> T-HUMAN-REQUIRED
- T-PENDING -> T-COMPLETE
- T-PENDING -> T-FOLLOW-UP
- T-PROBE-REQUIRED -> T-PENDING
- T-PROBE-REQUIRED -> T-HUMAN-REQUIRED

### Field: TICKET

| Field | Type | Required |
| :-- | :-- | :-- |
| source-pov-file-invocation | string | true |
| target-pov-ledger | string | true |
| status | TICKET-STATE | true |
| trace-history | string | true |
| affected-ids | set<string> | false |
| memory-anchor | string | true |
| description | string | true |
| static-investigation | string | true |
| escalation-reason | TICKET-ESCALATION-REASON | true |
| probe-feasibility | string | true |
| coverage-promotion-effect | string | true |

## 9.2. FSM and write rights

The conforming runtime validates every transition, escalation reason, required evidence field, actor right, and source state; an executor cannot assign a terminal or probe state merely by writing its token.

### Rule: ticket-transitions

Any POV may originate T-OPEN through the correct direct ledger or buffer. Only the target POV updates to T-PENDING or T-PROBE-REQUIRED. Only the origin accepts/rejects pending answers. Origin or target may propose out-of-scope; terminalization requires recorded mutual agreement or Human Hatch approval. Direct creation of T-HUMAN-REQUIRED is forbidden.

### Rule: ticket-terminal-transitions

Any active state may transition to T-INVALID or T-SUPERSEDED on valid proof, and to T-OUT-OF-SCOPE on a mutual scope decision.

## 9.3. Human Hatch action prerequisites

Each Human Hatch subtype has distinct prerequisites; a fake probe or transition through T-PROBE-REQUIRED is forbidden.

### Rule: ticket-escalation-prerequisites

Only Ticket Escalation may terminalize an unresolved ticket as human-required. Before Human Hatch: Ticket Escalation: perform and record bounded static investigation with exact methods, evidence, and remaining ambiguity; select exactly one closed Escalation Reason and record the probe feasibility/inapplicability assessment; for RUNTIME-REQUIRED, transition to T-PROBE-REQUIRED, propose a safe observable probe, and record the attempted probe or why it cannot be attempted safely, legally, or technically; for every non-runtime reason, remain T-OPEN or T-FOLLOW-UP and record why runtime observation is inapplicable, unauthorized, unsafe, destructive, unavailable, or incapable of resolving the required domain meaning/external fact; and invoke the subtype with exact ticket, scope, denominator, risk, and coverage consequences plus reviewer authority. Direct creation of T-HUMAN-REQUIRED, an escalation with Escalation Reason: None, or a reason inconsistent with the recorded feasibility assessment is invalid. Ticket Escalation may terminalize the ticket as T-HUMAN-REQUIRED, but the affected source remains C-GAP or C-PARTIAL unless an authorized human separately or concurrently approves an exact C-EXCLUDED decision under Decision/Scope Approval. No terminal ticket silently implies coverage. Human-required is an honest unresolved orchestration outcome distinct from exclusion: it blocks Exit A for any scope claiming the affected unit's closure, remains visible in partial synthesis and final risk reporting, and permits independent unaffected work to continue.

### Rule: decision-scope-approval

Decision/Scope Approval does not require a fake static investigation or probe. It requires exact proposed DEC content; observed SRC/CLM/ticket/contradiction inputs and immutable fingerprints; exact affected COV tuples and denominator arithmetic; alternatives, rationale, residual risks, rollout/rollback where applicable; and reviewer identity, authority basis, and timestamp. If used to resolve a ticket, the Ticket Escalation prerequisites still apply to that ticket.

### Rule: confirmation-prerequisites

Confirmation does not require a static investigation or probe. It requires a passing Exit E Candidate Report; candidate payload manifest; exact typed record IDs including PRF/HBK where applicable; exact SEMANTIC-RECORD-VERSION and SEMANTIC-CONTENT-FINGERPRINT bindings; required approved DECISION-CONTENT-VERSION/DECISION-CONTENT-FINGERPRINT bindings and matching approval envelopes; review scope and limitations; and reviewer identity and authority basis. For a PRF it MUST validate the complete §8.4 dependency closure; for other semantic records it MUST validate the complete applicable DERIVED-FROM closure. It may create CNF records and attach certification envelopes only; semantic payload edits are forbidden.

## 9.4. Probe specification and asynchronous handoff

Probe specifications declare exact target, question, inputs, environment, safety classification, masking and isolation approval, captures, and interpretation; the operator executes in an isolated sandbox and probe output is fingerprinted and sanitized before agent exposure.

### Rule: probe-handoff

The operator executes in an isolated sandbox and writes .extracted/probes/[PROBE-ID].log. In a bound Ticket Resolution invocation, the target POV ingests the log, creates runtime claims, updates affected records, and moves the ticket to T-PENDING. The originating POV then accepts it as T-COMPLETE or rejects it as T-FOLLOW-UP; probe ingestion never bypasses origin validation. Probe output MUST be fingerprinted and scrubbed of secrets/PII before agent exposure.

### Enum: SAFETY-CLASSIFICATION

- non-destructive
- controlled mutation

### Field: PROBE-SPECIFICATION

| Field | Type | Required |
| :-- | :-- | :-- |
| probe-id | string | true |
| target-src-cmp | string | true |
| question-expected-observable | string | true |
| exact-inputs-preconditions | string | true |
| required-environment-headers-config | string | true |
| safety-classification | SAFETY-CLASSIFICATION | true |
| data-masking-isolation-approval | string | true |
| disabled-capability-authorization | string | true |
| logs-db-diffs-payloads-to-capture | string | true |
| success-failure-interpretation | string | true |

## 9.5. No-mock integrity

When static analysis is insufficient and no sandbox is available, no fabricated paths, values, payloads, mocks, or specifications may be used as evidence; claims stay inferred or unknown and the Human Hatch path preserves the coverage gap.

### Rule: no-mock-rules

When static analysis is insufficient and a sandbox is unavailable: do not fabricate paths, values, payloads, mocks, or specifications as evidence; keep claims E-INFERRED or E-UNKNOWN; keep affected components S-DISCOVERED and ineligible for promotion; and use the Human Hatch path and preserve the coverage gap/partial state. A test double used to exercise already-proven local logic MAY be a testing mechanism but MUST NOT be treated as observation of an opaque dependency.

## 9.6. No-mock fallback payload

The no-mock fallback payload schema is declared below.

### Field: NO-MOCK-FALLBACK

| Field | Type | Required |
| :-- | :-- | :-- |
| suspected-behavior | string | true |
| evidence-claims-anchors | string | true |
| missing-runtime-observation | string | true |
| static-methods-attempted | string | true |
| probe-feasibility-attempt-result | string | true |
| proposed-human-verification | string | true |
| affected-promotion-coverage-exit-scope | string | true |

## 10.1. Evidence-backed `[R-SWEPT]`

A sweep record carries SWEEP-ID. A required traversal-matrix cell may be R-SWEPT only with a sweep record; every required cell has its own sweep, and claim-backed Not-Applicable cells have no sweep and remain visible in the denominator.

### Rule: r-swept-validity

R-SWEPT is invalid if expected counts cannot be reproduced, any discovered boundary lacks exactly one FRT, any frontier is open, any normalized unit is unreconciled, any material contradiction is open, or a wildcard is the only denominator. Terminal-at-discovery frontiers are included in terminal counts and disposition evidence; omission is not an optimization.

### Enum: TRAVERSAL-CELL-STATE

- R-SWEPT
- R-UNSWEPT
- R-IN-PROGRESS
- R-PENDING-HUMAN-REVIEW
- R-PENDING-PERSONA-ASSIGNMENT
- R-DEPROMOTION-PENDING
- Not-Applicable

### Field: SWEEP-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| snapshot-environment | string | true |
| persona-cluster-pov-track | string | true |
| expanded-concrete-entry-coordinates | set<string> | true |
| source-denominator | string | true |
| expected-mapped-excluded-counts-by-kind | string | true |
| open-terminal-frontier-counts | string | true |
| active-terminal-ticket-counts | string | true |
| export-reconciliation | string | true |
| coverage-result | string | true |
| closing-invocation | string | true |
| evidence-claims | set<string> | true |

## 10.2. Per-kind coverage arithmetic

Coverage is computed separately for every applicable kind, track, snapshot, and environment, then aggregated; a global percentage MUST NOT hide a failed subgroup.

### Rule: coverage-arithmetic-formulas

For each kind k and each environment/snapshot/track scope: Gk is the count of scope-specific rows in 16-SOURCE-COVERAGE where SRC kind is k and Scope is in In-Scope, Candidate, or Approved-Excluded. Xk is the count of rows in Gk where Scope is Approved-Excluded and C Status is C-EXCLUDED with an approved exact DEC. EffectiveDenominator k is Gk minus Xk. Covered k, Partial k, Gap k, and Candidate k count the effective-denominator rows by C Status C-COVERED, C-PARTIAL, C-GAP or Disposition Unmapped or Human-Blocked, and Scope Candidate respectively. The valid partition is EffectiveDenominator k equals Covered k plus Partial k plus Gap k. Kind k is closed when Candidate k, Partial k, and Gap k are zero and Covered k equals EffectiveDenominator k.

### Rule: coverage-covered-requirements

C-COVERED for an effective-denominator unit requires a terminal in-scope source disposition, atomic mapping, material CLMs, no unresolved contradiction, required reciprocal references, a terminal frontier, applicable export reconciliation, no active ticket affecting the unit, and, for a Disabled SRC, CMP, or UC, exactly one canonical CAP membership with any overlap explicitly decided. An approved exclusion yields C-EXCLUDED, not C-COVERED, and is governed by the subtraction above.

### Enum: SCOPE-CLASS

- In-Scope
- Candidate
- Approved-Excluded

### Enum: COVERAGE-STATUS

- C-COVERED
- C-PARTIAL
- C-GAP
- C-EXCLUDED

## 10.3. Composite Exit A — Deconstruction Closure

Exit A succeeds only when all declared conditions hold for the pinned scope.

### Rule: exit-a-composite

Exit A requires ticket closure with zero T-OPEN, T-PENDING, T-FOLLOW-UP, or T-PROBE-REQUIRED entries; every in-scope artifact A-STRUCTURALLY-COMPLETE; registry closure with no R-PENDING-HUMAN-REVIEW, R-PENDING-PERSONA-ASSIGNMENT, R-DEPROMOTION-PENDING, required R-UNSWEPT, or required R-IN-PROGRESS cell; terminal source disposition with no Candidate; every traversal frontier terminal; every normalized logical unit atomically mapped or approved excluded; a valid evidence-backed R-SWEPT per required cell with independent evaluation per Required applicability row; all effective denominator units C-COVERED with no C-PARTIAL or C-GAP; forensic referential integrity with ID uniqueness, existing reciprocal links, no invalid tombstone references, no orphan claims, no dangling references, and canonical ownership; no open material contradiction; disabled, shadow, and conditional units meeting the same static atomic closure rules with exactly one canonical CAP per effective-denominator Disabled SRC, CMP, or UC; and every T-HUMAN-REQUIRED paired with an exact approved exclusion or a scope that does not claim closure. Exit A writes a deterministic report to 22-GATE-REPORTS.md with pass/fail per condition, counts by kind, track, and environment, input fingerprints, validator version, and pinned snapshot.

### Enum: EXIT-A-CONDITION

- Ticket closure
- Acquisition closure
- Registry closure
- Source disposition
- Frontier closure
- Export reconciliation
- Sweep validity
- Atomic coverage
- Forensic referential integrity
- Contradiction closure
- Disabled equality and capability governance
- Human-required effect

## 10.4. Exit B — Stagnation

If a complete named iteration has zero valid mutations the pipeline halts as stagnant and reports remaining state; cosmetic edits do not count.

### Rule: exit-b-stagnation

Zero valid mutations means no new atomic source, component, claim, or frontier records, no ticket transition, no acquisition or reconciliation update, and no justified coverage change.

## 10.5. Exit C — Oscillation

If an actor-transition sequence of length 1 to 4 repeats for two consecutive completed cycles in one ticket trace, run one neutral Arbitration invocation, normally POV-5.

### Rule: exit-c-oscillation

If the neutral Arbitration invocation fails to break the cycle, freeze the ticket and escalate through Human Hatch.

## 10.6. Exit D — Limit exhaustion and deterministic iteration accounting

The 26 canonical iteration tokens are a total per-system budget across the complete protocol run; the budget never resets, and the system MUST NOT proceed beyond ZULU.

### Rule: iteration-accounting

Named iterations use exactly the canonical tokens ALFA, BRAVO, CHARLIE, DELTA, ECHO, FOXTROT, GOLF, HOTEL, INDIA, JULIETT, KILO, LIMA, MIKE, NOVEMBER, OSCAR, PAPA, QUEBEC, ROMEO, SIERRA, TANGO, UNIFORM, VICTOR, WHISKEY, X-RAY, YANKEE, and ZULU. Human-facing typography such as X-ray is display-only and MUST be normalized to X-RAY before persistence or comparison; it is not an accepted stored alias. A named iteration begins with its first logged invocation and completes only when all admitted invocations have terminal log entries, all persona-local buffers produced in it have been reconciled, and the iteration checkpoint records mutation count, remaining tickets, frontiers, and candidates, and the next action. Parallel invocations share the same current token and do not consume extra tokens.

### Rule: iteration-reopening

A synthesis GAP discovered while an iteration is open is reconciled within that iteration when possible; if that iteration is already complete, reopening consumes the next unused token. Re-running synthesis after closure continues under the current open token or next unused token and MUST NOT return to ALFA or reuse a completed token. Unclosed work at completion of ZULU halts and requires human intervention.

## 11.1. Partial versus Final Synthesis

Partial-Synthesis runs before Exit A only for an explicitly bounded persona/cluster/track/snapshot; Final-Synthesis runs only after Composite Exit A passes.

### Rule: partial-vs-final-synthesis

Partial-Synthesis blocks remain B-DRAFT, are labeled provisional, and state the incomplete denominator. They MUST NOT be reconstruction inputs. Final-Synthesis writes complete structured catalogs from closed facts. Newly discovered gaps go to SYNTHESIS-GAP-BUFFER.md; Sequential Reconciliation opens tickets, creates or updates missing SRC/frontier/coverage records, immediately marks the pinned Exit A report invalid, and marks all dependent B blocks stale. Deconstruction is re-armed and confirmation/handoff is prohibited until Exit A passes again on a new gate-input fingerprint. Final synthesis cannot directly mutate ledgers.

### Rule: cross-reference-updates

After IDs exist, Cross-Reference-Reconciliation may update only fields such as Bound Use-Cases, Rules, Entities, Interfaces, and reciprocal references. It MUST NOT change observed behavior, decisions, or evidence; because these references are semantic payload, a changed reference increments the semantic record version/fingerprint and invalidates any certification envelope.

### Rule: gap-buffer-merge

Partial/Final Synthesis writes only Pending GAP records and cannot create tickets or frontiers. Sequential Reconciliation sorts by DEDUP-KEY, then GAP-ID; merges equal keys while retaining every source invocation/block/fingerprint; records a contradiction rather than merging incompatible target proposals; and deterministically creates the minimum required work. A missing cross-POV answer creates one canonical ticket in the proposed ledger; a newly discovered traversal boundary creates one FRT; a gap requiring both creates both and cross-links them. Existing equivalent active work is linked rather than duplicated. Sequential Reconciliation updates disposition atomically, invalidates the pinned Exit A report for material gaps, applies targeted staleness, and logs iteration-budget effects under §10.6.

### Enum: GAP-MATERIALITY

- Material
- Supporting

### Enum: GAP-DISPOSITION

- Pending
- Ticket-Created
- Frontier-Created
- Both-Created
- Merged
- Invalid
- Approved-Excluded

### Field: SYNTHESIS-GAP

| Field | Type | Required |
| :-- | :-- | :-- |
| source-synthesis-invocation-block-id | string | true |
| observed-missing-requirement | string | true |
| affected-ids | string | true |
| source-fingerprints-snapshot | string | true |
| proposed-target-persona-cluster-track-pov-ledger | string | true |
| dedup-key | hash | true |
| materiality | GAP-MATERIALITY | true |
| disposition | GAP-DISPOSITION | true |
| resulting-ticket-frt-dec-ids | string | true |

## 11.2. Common synthesis block header

Every atomic synthesis block MUST contain a semantic payload followed by a separately delimited certification envelope.

### Rule: common-header-requirements

Coverage is semantic payload; blueprint status is envelope metadata. B-CONFIRMED requires C-COVERED, no stale dependency, a passing candidate binding, and a valid human confirmation of the exact semantic version/fingerprint. Attaching confirmation MUST NOT change semantic bytes. Empty domains require inventory-backed zero claims and may then be confirmed as empty.

### Field: SYNTHESIS-BLOCK-HEADER

| Field | Type | Required |
| :-- | :-- | :-- |
| semantic-record-version | int | true |
| semantic-content-fingerprint | string | true |
| coverage-status | string | true |
| derived-from | string | true |
| material-claims | string | true |
| block-confidence-rank | string | true |
| block-evidence-profile | string | true |
| supersession | string | true |

### Field: CERTIFICATION-ENVELOPE

| Field | Type | Required |
| :-- | :-- | :-- |
| blueprint-status | string | true |
| confirmation | string | true |
| confirmation-envelope-version | int | true |
| confirmation-envelope-audit | string | true |

## 11.3. Architecture blueprint — `90-ARCH-BLUEPRINT.md`

The architecture blueprint contains system boundary and context diagrams; actors/personas and external systems; observed legacy technology separated from target decisions; atomic MOD records; inter-module contracts and error propagation; deployment topology; architectural invariants and modernization consequences; capability navigation; and a prominent active-versus-disabled architecture view. Purely navigational module indexes are grouping containers under §4.2 and list MOD IDs only.

### Rule: module-requirements

Every semantic architecture statement about responsibility, public interface, encapsulation, dependency, or data passing MUST belong to an atomic MOD and have material CLMs. MOD to IF/ER/DEP/BR/UC/CAP/FLT reciprocal references are required where applicable.

### Field: ARCHITECTURE-MODULE

| Field | Type | Required |
| :-- | :-- | :-- |
| responsibility-domain-boundary | string | true |
| public-interfaces | string | true |
| encapsulated-state-entities | string | true |
| dependencies-dependents | string | true |
| data-and-events-passed | string | true |
| invariants-rules-owned | string | true |
| personas-use-cases-capabilities | string | true |
| failure-propagation | string | true |
| observed-legacy-placement | string | true |
| target-decision | string | true |

## 11.4. Entity — `91-DATA-MODEL.md`

Entity records describe one storage-backed entity with domain context, storage kind, physical sources, keys, fields, relationships, lifecycle, invariants, audit, retention, permissions, capability impact, and schema drift.

### Enum: STORAGE-KIND

- table
- view
- collection
- cache-key-space
- file
- in-memory

### Field: ENTITY-COLUMN

| Field | Type | Required |
| :-- | :-- | :-- |
| name | string | true |
| type | string | true |
| null | string | true |
| default | string | true |
| constraints | string | true |
| unique | string | true |
| claims | string | true |

### Field: ENTITY-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| domain-context-purpose | string | true |
| storage | STORAGE-KIND | true |
| physical-src-ids | string | true |
| primary-key-secondary-keys-indexes | string | true |
| columns-fields | string | true |
| relationships | string | true |
| enums-lookup-domains | string | true |
| lifecycle-state-fields | string | true |
| derived-computed-fields | string | true |
| invariants | string | true |
| audit-soft-delete | string | true |
| retention-lifecycle | string | true |
| permissions | string | true |
| capability-state-impact | string | true |
| schema-drift-matrix | string | true |

## 11.5. Relationship — `91-DATA-MODEL.md`

Relationship records describe one entity relationship with cardinality, enforcement, delete and update semantics, inverse access, business constraints, and reciprocal entity references.

### Enum: RELATIONSHIP-CARDINALITY

- 1:1
- 1:N
- M:N

### Enum: RELATIONSHIP-ENFORCEMENT

- DB-FK
- application
- mixed
- none

### Enum: DELETE-UPDATE-SEMANTICS

- cascade
- restrict
- nullify
- custom

### Field: RELATIONSHIP-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| cardinality | RELATIONSHIP-CARDINALITY | true |
| enforcement | RELATIONSHIP-ENFORCEMENT | true |
| delete-update-semantics | DELETE-UPDATE-SEMANTICS | true |
| inverse-access | string | true |
| business-constraint | string | true |
| reciprocal-entity-references | string | true |

## 11.6. State machine — `91-DATA-MODEL.md`

A state machine is required only for an eligible lifecycle scope: an entity aggregate or capability with a finite, behaviorally meaningful lifecycle, not every boolean or presentation state.

### Rule: state-machine-eligibility

Eligibility decisions use the §12.1 lifecycle table so absence cannot be assumed. Unknown requires a ticket, is treated as eligible for closure, and blocks Exit E until resolved to Eligible with a complete SM or Ineligible with claim-backed provenance.

### Field: STATE-TRANSITION

| Field | Type | Required |
| :-- | :-- | :-- |
| from | string | true |
| to | string | true |
| event-uc | string | true |
| guard-br | string | true |
| effects | string | true |
| claims | string | true |

### Field: STATE-MACHINE

| Field | Type | Required |
| :-- | :-- | :-- |
| scope-owner-entity-or-aggregate | string | true |
| state-fields | string | true |
| states-meaning | string | true |
| initial-terminal-states | string | true |
| transitions | string | true |
| invalid-recovery-transitions | string | true |
| environment-capability-variants | string | true |

## 11.7. Database routine — `91-DATA-MODEL.md`

Database routine records describe one trigger, function, procedure, rule, generated column, refresh, or constraint helper with owner, firing semantics, body intent, inputs and outputs, enforced rules, callers, failure semantics, canonical drift, and target decision.

### Rule: database-routine-drift

A production DB routine absent from the canonical DB source creates a material drift contradiction and follows the ticket/Human Hatch path.

### Enum: DB-ROUTINE-KIND

- trigger
- function
- procedure
- rule
- generated-column
- materialized-refresh
- constraint-helper

### Enum: DB-TARGET-DECISION

- Retain
- Lift
- Redesign
- Retire
- Undecided

### Field: DATABASE-ROUTINE

| Field | Type | Required |
| :-- | :-- | :-- |
| kind | DB-ROUTINE-KIND | true |
| owner-table-entity | string | true |
| firing-event-timing-granularity | string | true |
| atomic-body-intent | string | true |
| inputs-outputs | string | true |
| entities-read-written | string | true |
| invariants-rules-enforced | string | true |
| callers-firers | string | true |
| failure-transaction-semantics | string | true |
| canonical-db-source-drift | string | true |
| target-decision | DB-TARGET-DECISION | true |
| target-decision-dec | string | true |

## 11.8. Business rule — `92-BUSINESS-RULES.md`

Business rule records describe one deterministic, framework-agnostic rule with exactly one authoritative owner, optional consistent mirrors, condition, effects, precedence, violation policy, boundary cases, idempotency, and target mapping decision.

### Rule: single-rule-owner

One authoritative rule owner is REQUIRED. Multiple physical enforcement points are allowed only as consistent defense-in-depth mirrors, not competing authorities.

### Enum: VIOLATION-POLICY

- reject
- warn
- clamp
- default
- compensate
- rollback

### Enum: BUSINESS-RULE-CAPABILITY-STATE

- Active
- Disabled

### Field: BUSINESS-RULE-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| authoritative-rule-owner | string | true |
| defense-in-depth-mirrors | string | true |
| applies-when | string | true |
| condition | string | true |
| action-effect | string | true |
| precedence-conflict | string | true |
| inputs-read-writable-effects | string | true |
| violation-policy | VIOLATION-POLICY | true |
| boundary-cases-null-locale-overflow | string | true |
| idempotency | string | true |
| capability-state | BUSINESS-RULE-CAPABILITY-STATE | true |
| target-mapping-decision | string | true |

## 11.9. Use case — `93-USE-CASES.md`

Use case records describe one observable user or system goal with actors, trigger, preconditions and postconditions, main and alternate flows, failure recovery, reciprocal references, a Given/When/Then behavioral contract, and an explicit target decision.

### Rule: use-case-visibility

Disabled use cases remain visible, are not presented as current normal behavior, and require an explicit target decision.

### Enum: USE-CASE-CAPABILITY-CLASSIFICATION

- Active
- Conditional
- Shadow
- Disabled
- Retired

### Enum: USE-CASE-TARGET-DECISION

- Preserve
- Change-By-Approved-Migration
- Redesign
- Retire

### Field: USE-CASE-STEP

| Field | Type | Required |
| :-- | :-- | :-- |
| step | string | true |
| action | string | true |
| rules | string | true |
| data-r-w | string | true |
| interface | string | true |
| claims | string | true |

### Field: USE-CASE

| Field | Type | Required |
| :-- | :-- | :-- |
| primary-persona-secondary-actors | string | true |
| goal | string | true |
| trigger-entry | string | true |
| preconditions-postconditions | string | true |
| invariants-relied-upon | string | true |
| capability-classification | USE-CASE-CAPABILITY-CLASSIFICATION | true |
| observed-legacy-behavior | string | true |
| rules-entities-interfaces-side-effects | string | true |
| behavioral-contract | string | true |
| target-decision | USE-CASE-TARGET-DECISION | true |
| target-decision-dec | string | true |

## 11.10. Interface — `94-INTERFACES.md`

Interface records describe one external or internal interface with transport, actors, request and response schemas, error semantics, guarantees, authentication and authorization, versioning, limits, bound records, capability matrix, compatibility baseline, and target migration decision.

### Rule: interface-compatibility-baseline

External semantics are compatibility baselines, not absolutely frozen. Breaking changes require an approved, versioned migration decision, consumer impact, rollout, and rollback plan.

### Field: INTERFACE-REQUEST-FIELD

| Field | Type | Required |
| :-- | :-- | :-- |
| field | string | true |
| type | string | true |
| required | string | true |
| default | string | true |
| validation-br | string | true |
| claim | string | true |

### Field: INTERFACE-RESPONSE-FIELD

| Field | Type | Required |
| :-- | :-- | :-- |
| field | string | true |
| type | string | true |
| notes | string | true |
| claim | string | true |

### Field: INTERFACE-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| channel-transport-direction | string | true |
| actor-persona-auth | string | true |
| protocol-content-type | string | true |
| request-schema | string | true |
| response-schema | string | true |
| error-semantics-retryability | string | true |
| guarantees | string | true |
| authn-authz | string | true |
| versioning-compatibility | string | true |
| rate-limits-backpressure | string | true |
| bound-uc-er-ids | string | true |
| capability-state-environment-matrix | string | true |
| compatibility-baseline | string | true |
| target-migration-decision | string | true |

## 11.11. Deployment, configuration, and scheduling — `95-DEPLOYMENT.md`

Deployment records describe components, connections, failure exposure, and environment matrices; configuration records describe flags with exact keys, defaults, sources, effects, enable and disable semantics, and materialization; schedule records describe jobs with triggers, owners, actions, failure handling, locks, and capability state.

### Rule: configuration-effect

Every configuration flag requires a known effect or a gap. Secret values are never recorded.

### Enum: CONFIG-MEDIUM

- env
- DB
- flag service
- file
- build-time

### Field: DEPLOYMENT-ELEMENT

| Field | Type | Required |
| :-- | :-- | :-- |
| component-placement-replication | string | true |
| connections | string | true |
| single-point-of-failure | string | true |
| environment-snapshot-matrix | string | true |

### Field: CONFIGURATION-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| medium | CONFIG-MEDIUM | true |
| exact-key | string | true |
| type-default-allowed-values | string | true |
| source-precedence-setter | string | true |
| behavioral-effect | string | true |
| enable-disable-semantics-and-environment-matrix | string | true |
| app-side-db-side-materialization | string | true |

### Field: SCHEDULE-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| type-trigger-frequency-window | string | true |
| owner-persona | string | true |
| actions | string | true |
| failure-retry-dead-letter-idempotency | string | true |
| dependencies-locks | string | true |
| capability-state | REGISTRY-ENTRYCLUSTER-CAPABILITY-STATE | true |

## 11.12. NFR and security — `96-NON-FUNCTIONAL-SECURITY.md`

NFR records describe one attribute with observed value, scope, drivers, and target SLO; security findings describe one category, classification, risk level, affected scope, safeguard or gap, and required review; fault records describe one failure class with propagation, retry and rollback, quarantine, outcomes, variants, and related rules.

### Rule: security-reciprocal

Every relevant 0C finding MUST have a reciprocal SEC mirror. Neither file declares legal compliance. Every distinct material failure class discovered in forensic records MUST map to an FLT; 96 MUST include a system-level failure taxonomy and retry/rollback/compensation/dead-letter topology, with an evidence-backed zero statement for any genuinely absent mechanism.

### Enum: NFR-ATTRIBUTE

- latency
- throughput
- availability
- RPO
- RTO
- scalability
- operability
- other

### Enum: SECURITY-CATEGORY

- auth
- authorization
- transport
- at-rest
- in-transit
- PII
- secrets
- audit
- threat
- privacy

### Enum: RISK-LEVEL

- high
- medium
- low
- unassessed

### Field: NFR-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| attribute | NFR-ATTRIBUTE | true |
| observed-inferred-value | string | true |
| scope-environment | string | true |
| driver-constraint | string | true |
| target-decision-slo | string | true |

### Field: SECURITY-FINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| category | SECURITY-CATEGORY | true |
| finding-data-classification | string | true |
| risk-level | RISK-LEVEL | true |
| affected-ids-environments | string | true |
| observed-safeguard-gap | string | true |
| human-review-target-mitigation | string | true |

### Field: FAULT-RECORD

| Field | Type | Required |
| :-- | :-- | :-- |
| failure-class-trigger-detection-point | string | true |
| propagation-path | string | true |
| retry-backoff-timeout-circuit-behavior | string | true |
| rollback-compensation-transaction-boundary | string | true |
| dead-letter-quarantine-operator-action | string | true |
| user-system-outcome-observability | string | true |
| environment-capability-variants | string | true |
| related-rules-claims | string | true |

## 11.13. Persona profile

Profile Synchronization writes one profile and its 0A row atomically.

### Rule: profile-synchronization

Profile Synchronization is the only writer of profile semantic payload and synchronizes its mandatory PRF-ID atomically with the matching 0A persona-prefix row. Authorized Human Hatch: Confirmation or Reconstruction-Handoff confirmation actions may update only the delimited certification envelope after loading and validating the complete §8.4 direct-and-transitive dependency closure and verifying the exact candidate-bound PRF ID, semantic version, and semantic fingerprint; they MUST NOT edit profile semantic content/version. A profile without a registry row, a registry row without a profile, a duplicate PRF for one prefix, or a PRF whose canonical prefix coordinate disagrees with 0A is invalid after Profile Synchronization begins.

### Field: AUTHORIZATION-MATRIX-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| resource-capability | string | true |
| allowed | string | true |
| enforcement | string | true |
| claims | string | true |

### Field: PERSONA-PROFILE

| Field | Type | Required |
| :-- | :-- | :-- |
| prf-id | string | true |
| prefix-actor-class | string | true |
| canonical-and-shared-entries | string | true |
| authentication-and-session | string | true |
| authorization-matrix | string | true |
| visible-data-scope | string | true |
| key-active-workflows | string | true |
| disabled-historical-workflows | string | true |
| triggers-events-schedules | string | true |
| ui-surfaces-and-state | string | true |
| failure-edge-behavior | string | true |
| dependencies-couplings | string | true |
| observed-legacy-facts | string | true |
| target-decisions | string | true |
| semantic-record-version | int | true |
| semantic-content-fingerprint | string | true |
| derived-from | string | true |
| material-claims | string | true |
| block-confidence-rank | string | true |
| block-evidence-profile | string | true |

## 12.1. Traceability — `20-TRACEABILITY.md`

One atomic row per source-to-semantic mapping; multiple rows MAY represent multiple mappings, but source coverage is evaluated once per SRC and requires all applicable mappings. The ID registry, source denominator query fingerprints, lifecycle eligibility decisions, exclusions, and reciprocal reference check results are included in this file.

### Rule: lifecycle-eligibility

Lifecycle eligibility uses one row per candidate entity aggregate or capability. Classification is a semantic predicate: the classifier MUST record provenance and claim/evidence bindings. Unknown is fail-inclusive and is treated as eligible for completeness and gate purposes until resolved; it requires a ticket and cannot justify omission of an SM. The conforming runtime deterministically verifies row uniqueness, snapshot consistency, required provenance, and the corresponding SM presence or explicit ineligibility result.

### Rule: traceability-invariants

Every source unit in the inventory has at least one traceability row, including exact exclusions; every ER is read or written by at least one UC, or is explicitly classified as structural/reference-only with a claim-backed reason; every IF is consumed or produced by at least one UC, or is exactly excluded/retired; every UC step cites at least one BR, entity read/write, interface effect, state transition, or explicit claim-backed no-op; every SCHED and DR maps to at least one UC, BR, or invariant; every CFG maps to a behavioral effect or exact approved exclusion; every semantic MOD has reciprocal links to each applicable IF, ER, DEP, BR, UC, CAP, and FLT, while module indexes carry only member links; every effective-denominator Disabled SRC, CMP, and UC maps to exactly one canonical CAP, or is exactly approved excluded, and any overlap has an explicit typed relationship and approved DEC with one target-decision owner; every synthesis ID is reachable from a source/claim lineage and from at least one handbook navigation path when confirmed.

### Enum: CLASSIFICATION

- Eligible
- Ineligible
- Unknown

### Field: TRACEABILITY-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| src | string | true |
| persona-owner-prf | string | true |
| track-cap | string | true |
| pov-cmp | string | true |
| mod | string | true |
| uc | string | true |
| br | string | true |
| er-rel-sm-dr | string | true |
| if | string | true |
| dep-cfg-sched | string | true |
| prv-sec-nfr-flt | string | true |
| hbk | string | true |
| clm-dec-cnf | string | true |
| c-status | string | true |

### Field: LIFECYCLE-ELIGIBILITY

| Field | Type | Required |
| :-- | :-- | :-- |
| scope-id | string | true |
| classification | CLASSIFICATION | true |
| reason | string | true |
| clm-evidence-bindings | string | true |
| classified-by-invocation | string | true |
| snapshot | string | true |

## 12.2. Required reciprocal references

Validators enforce the required reciprocal reference pairs, and no dangling, orphan, duplicate-owner, invalid-type, or unapproved tombstone reference is allowed at Exit A or E.

### Rule: required-reciprocal-references

Validators enforce ER with REL, and ER state field with SM when eligible; DR caller/firer with CMP DB-ROUTINES-FIRED; BR authoritative owner with MOD, BR with UC steps, and BR with ER effects; MOD with IF/ER/DEP/BR/UC/CAP/FLT where applicable, while purely navigational module indexes link only to member MOD IDs; UC with persona, UC with IF, and UC with ER/BR/SCHED/DR; IF with UC and interface entity mappings; CFG with affected UC/BR/IF/CAP; SCHED with UC/ER/persona; PRV with SEC; NFR with driver; FLT with UC/IF/SCHED/DEP/BR; AUTH with persona/IF/enforcement point; SHR with participating persona pointers; CAP with canonical member IDs, CFG, and DEC, with exact-one disabled membership and explicit overlap relations; every COV with its unique immutable SRC-version/persona/track/environment/snapshot tuple and evidence; every CND with its helper candidate/action identity and resulting disposition IDs; every B-CONFIRMED block, PRF, and HBK with a valid CNF carrying the exact candidate-bound typed ID/version/fingerprint; every exclusion/target choice with DEC; every persona-prefix row in 0A with exactly one PRF, and every PRF with that row plus its complete §8.4 dependency closure; every HBK with one normalized handbook path/stable anchor, its upstream synthesis IDs/versions/fingerprints, candidate-manifest entry, and at least one handbook navigation path; CMP with SRC, claims, persona references, and dependencies; normalized coordinate with SRC with atomic CMP/edge or exclusion; and handbook HBK section with normalized relative path/stable anchor with synthesis IDs and exact versions/fingerprints. No dangling, orphan, duplicate-owner, invalid-type, or unapproved tombstone reference is allowed at Exit A or E.

## 12.3. Synthesis-ID timing

Discovery records may contain Pending-Synthesis-Reference placeholders but never invented IDs; synthesis allocates deterministic IDs and reconciliation replaces placeholders, incrementing semantic versions and invalidating prior envelopes.

### Rule: synthesis-id-allocation-timing

Discovery records MAY contain Pending-Synthesis-Reference placeholders, never invented IDs. Partial/Final Synthesis allocates deterministic semantic IDs from canonical keys. Cross-Reference-Reconciliation then replaces placeholders and adds reciprocal references without inventing new legacy facts; because reference bindings are part of semantic payload, it increments the affected semantic version/fingerprint and invalidates prior envelopes. Pending placeholders block confirmation and Exit E, but do not block early forensic discovery.

## 12.4. Dependency and impact provenance

Each synthesis and handbook block maintains DERIVED-FROM edges to exact record versions and fingerprints; upstream changes trigger reverse dependency closure, targeted staleness, logged impact, bounded re-derivation, and preservation of unaffected confirmed blocks.

### Rule: dependency-impact-provenance

Each synthesis/handbook block maintains DERIVED-FROM edges to exact record versions and fingerprints. When an upstream source, claim, mapping, decision, or exclusion changes: compute reverse dependency closure; mark only affected synthesis blocks and handbook sections B-STALE; log cause and impact in 0G and the block; re-run the relevant bounded synthesis and validation; and preserve unaffected confirmed blocks. Broad whole-file staleness is allowed only when dependency granularity is unavailable; that limitation is itself reported.

## 12.5. Deterministic validation summary

Every validator run records schema/protocol/check-registry/validator versions; authoritative input file fingerprints and a deterministic input-set fingerprint; checks executed and stable check IDs; for every check and blocker row, the complete sorted evidence-binding tuples and recomputed §4.1.2 evidence-set fingerprint; pass/fail counts and exact offending IDs; coverage arithmetic by kind/track/environment; broken links, invalid diagrams, duplicate IDs, orphan/dangling references, stale blocks, unsupported statuses, and unresolved contradictions; any authorized human review input used by a check, including reviewer identity/authority, review scope, timestamp, exact candidate-bound fingerprints, limitations, and result; and output fingerprint and timestamp.

### Rule: evidence-bindings

Evidence bindings are externally observable in the validation summary with one or more rows per Check/Blocker ID as required, and rows sort by Row Kind, Check / Blocker ID, Normalized Authoritative Input Path, Binding Kind, Typed Record / Artifact ID or None, Version or None, then Authoritative Fingerprint. Binding Kind uses the closed §4.1.2 mapping. A changed input invalidates the prior summary. Deterministic values supplied by an executor MUST be independently reproduced. A missing required check, unknown or duplicate Check ID, incomplete evidence-binding set, or mismatch against the applicable §15.1.1 registry makes the summary and dependent gate result invalid.

### Field: EVIDENCE-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| row-kind | string | true |
| check-blocker-id | string | true |
| normalized-authoritative-input-path | string | true |
| binding-kind | string | true |
| typed-record-artifact-id-or-none | string | true |
| version-or-none | string | true |
| authoritative-fingerprint | string | true |

## 13.1. Audience paths and progressive disclosure

START-HERE.md MUST provide short audience paths through the handbook.

### Rule: audience-path-requirements

Each chapter begins with prose, diagrams, and representative examples before dense references. IDs appear as hyperlinks/reference anchors, not the primary user interface. Evidence rank/profile is visible but unobtrusive: for example a short badge or note linking to evidence, not repeated raw claim tables in the main narrative.

### Enum: AUDIENCE-ROLE

- product/domain owner
- architect/developer
- security/privacy reviewer
- data engineer
- QA
- operator/SRE

## 13.2. Required chapter behavior

Every handbook chapter distinguishes observed legacy from target decisions, separates normal from shadow and disabled behavior, links to upstream records, supports navigation, avoids unexplained acronyms, shows evidence-backed or labeled examples, and passes link and diagram validation.

### Rule: chapter-requirements

Every chapter MUST distinguish Observed Legacy from Target Decision; separate Normal behavior from Shadow/Conditional and Disabled behavior; link to upstream synthesis records and relevant source/evidence on demand; include navigation to previous/next, START-HERE, glossary, decisions, and gaps; avoid unexplained acronyms and identifier-only prose; show examples that are evidence-backed or explicitly labeled illustrative; and pass link and diagram validation. 08-BACKGROUND-PROCESSING.md MUST include the consolidated fault/resilience view: failure taxonomy, propagation, retries/backoff, transaction rollback, compensation, dead-letter/quarantine paths, observability, and operator/user outcomes, cross-linked to FLT records.

## 13.3. Disabled and Dormant Features chapter

12-DISABLED-AND-DORMANT-CAPABILITIES.md is mandatory and prominent, indexing every canonical capability with exact disabled members, state, mechanism, evidence, dependencies, history, risk, cross-links, and target decision.

### Rule: disabled-chapter-requirements

12-DISABLED-AND-DORMANT-CAPABILITIES.md is mandatory and prominent. It provides an index and one human-oriented section per canonical CAP containing exact disabled member IDs, status/environment matrix, mechanism, evidence rank/profile, dependencies, historical/intended personas and behavior, safety/risk, related architecture/data/rules/use-cases/interfaces, and target decision Restore, Preserve Dormant, Redesign, or Retire. It MUST demonstrate that every effective-denominator Disabled SRC, CMP, and UC appears under exactly one canonical CAP or an exact approved exclusion, and disclose any approved overlap relationship/decision without duplicating decision ownership. Disabled features MUST also be cross-linked from every affected chapter. They MUST NOT be blended into current active behavior or omitted because they are not executable in production.

## 13.4. Handbook quality gate

Before candidacy the runtime validates chapter non-vacuity, links and anchors, deterministic diagram parsing, synthesis discoverability, active and disabled clarity, decision separation, and absence of stale or provisional sections; after candidacy, human readability review is recorded only through complete passing CNF bindings.

### Rule: handbook-quality-checks

Before candidacy, the conforming runtime validates: all required chapters exist and each contains at least one claim-backed semantic handbook assertion or an evidence-backed zero-domain statement, while headings, navigation links, placeholders, or grouping prose alone are vacuous; all internal links and anchor targets resolve; every Mermaid/diagram block parses with a conforming deterministic parser, and unavailable parsing capability is unsupported and blocks candidacy rather than silently falling back to invisible self-assessment; each candidate synthesis block is discoverable from at least one handbook path; active versus disabled presentation is unambiguous; target decisions are not presented as observed facts; and no B-STALE, provisional, or pending-reference section remains.

### Rule: handbook-readability-review

After candidacy, authorized human readability review is recorded through one or more §5.7 CNF records whose combined exact HBK bindings cover every candidate handbook section and whose Handbook Readability Review is Pass; their hashed review scope/limitations state how core user goals were followed without reading raw catalogs. A CNF that does not review HBK content uses Not-Applicable; Fail or incomplete HBK coverage blocks content readiness. The Exit E Content-Readiness Report MUST validate reviewer authority, candidate manifest/report bindings, complete HBK coverage, timestamps, results, and unchanged handbook fingerprints. Human review cannot substitute for deterministic link, diagram, or non-vacuity checks, and an unrecorded review has no gate effect.

## 14.1. Decision log

15-DECISIONS.md is the authoritative decision register from Preflight onward; Reconstruction-Handoff projects a human-readable log without forking or reinterpreting decisions, and superseded decisions remain tombstoned and linked.

### Rule: decision-log-projection

15-DECISIONS.md is the authoritative decision register from Preflight onward. Reconstruction-Handoff produces a human-readable decision-log projection containing every applicable approved record and its source link; it MUST NOT fork or reinterpret the decision. The canonical schema is §5.5. Target decisions never alter observed facts. Superseded decisions remain tombstoned and linked.

## 14.2. Modernization mapping templates

The modernization mapping templates are decision prompts, not mandates; the target architecture owner chooses and records each result.

### Rule: mapping-prompts

The following are decision prompts, not mandates: global/out-of-band state maps to request scope, explicit service state, or managed store; mixed server/client presentation ownership maps to deliberate SSR/SPA/hybrid ownership with synchronization rules; scattered auth checks map to a centralized policy boundary with documented defense-in-depth checks; DB-autonomous rules map to intentionally retain, lift, redesign, or retire per DR and invariant; synchronous cross-availability side effects map to evaluating durable commands/events; string-composed persistence maps to evaluating a controlled data-access boundary; filesystem side effects map to evaluating durable storage retaining lifecycle semantics; application-only relationships map to enforcing physically or documenting compensating consistency; derived fields map to retain, materialize, recompute, or remove with a decision; and soft-delete/audit mechanics map to preserving the requirement, not necessarily legacy columns. The target architecture owner chooses the result and records it.

## 14.3. Business-rule ownership

Every business rule maps to exactly one authoritative target owner; optional validation, gateway, UI, or DB mirrors must stay semantically consistent and cite the authoritative rule, and conflicting mirrors are prohibited.

### Rule: business-rule-ownership

Every BR maps to exactly one authoritative target owner. Optional validation, gateway, UI, or DB mirrors are permitted only if they remain semantically consistent and cite the authoritative rule. Conflicting mirrors are prohibited.

## 14.4. Interface compatibility

External interface schemas, errors, ordering, idempotency, delivery, and timing form a compatibility baseline that changes only through an approved versioned migration decision; internal transport may change freely while externally observable contracts remain satisfied.

### Rule: interface-baseline

External interface schemas, errors, ordering, idempotency, delivery, and timing guarantees form a compatibility baseline. They MAY change only through an approved versioned migration decision identifying consumers, compatibility period, data migration, rollout, observability, and rollback. Internal transport MAY change freely if externally observable and use-case contracts remain satisfied or have approved changes.

## 14.5. Equivalence and acceptance suite

The equivalence suite generates reviewable acceptance tests from use-case contracts, property assertions from invariants and rules, state-transition tests, contract tests, disabled-capability and restoration tests, and migration tests, with semantics sourced from confirmed synthesis and approved decisions.

### Rule: equivalence-generation

Generate example acceptance tests from every UC Given/When/Then contract; property assertions from invariants and BR conditions; state-transition tests from every eligible SM; contract tests for external IF compatibility baselines; disabled-capability tests that verify intended dormant state and prevent accidental enablement, plus restoration tests only when target decision is Restore; and migration tests for approved behavior changes. Generated tests are reviewable derivatives. Their expected semantics come from B-confirmed synthesis and approved decisions.

### Enum: EQUIVALENCE-TEST-KIND

- example acceptance tests
- property assertions
- state-transition tests
- contract tests
- disabled-capability tests
- migration tests

## 14.6. Sole-input invariant

Legacy source may be treated as unavailable during reconstruction only when the bundle has a validated signed outer manifest with EXIT-E-STATUS: Passed, contains only confirmed applicable blocks, is signed, pins all authoritative fingerprints, and includes assurance, scope, decisions, equivalence, and residual risks.

### Rule: sole-input-conditions

Legacy source may be treated as unavailable during reconstruction only when the bundle has a validated signed outer bundle manifest whose hashed payload contains EXIT-E-STATUS: Passed; contains only applicable B-CONFIRMED semantic blocks; is covered by the required human signatures on that outer manifest; pins all authoritative fingerprints; and includes assurance, scope certificate, decisions, equivalence suite, and known residual risks. Anything absent from that certified scope is out of reconstruction scope and requires a new decision or a reopened deconstruction cycle.

## 15.1. Strict conditions

Exit E remains Pending through steps 1-5 of §15.4. The stage-4 Exit E Content-Readiness Report validates all content, decision, confirmation, and package-member conditions available at that stage, but it is not the Exit E pass artifact and MUST NOT state or imply Passed. Exit E transitions to Passed only at step 6 when all conditions below are true and the signed outer bundle manifest's hashed payload validates with the literal field EXIT-E-STATUS: Passed.

### Rule: exit-e-conditions

Exit E requires: composite Exit A passed for the same pinned snapshot or a documented later non-semantic packaging change preserving its inputs; every effective in-scope source unit C-COVERED with exact, risk-assessed, disclosed exclusions removed from denominator arithmetic and no C-GAP or C-PARTIAL remaining; all 90-96 mandatory domains complete with evidence-backed zero-domain records where genuinely empty and valid semantic versions/fingerprints; the complete handbook semantic payload existing as B-POPULATED before candidacy, passing quality gates, synchronized to synthesis, and included in the candidate payload manifest without later semantic regeneration; all reciprocal references, legal typed IDs including COV/CND/MOD/PRF/HBK, claims, ownership, persona purity, diagrams, and links validating; every eligible lifecycle scope having a complete SM and every ineligible decision being explicit; every UC having a complete behavioral contract and reciprocal BR/ER/IF references; every CFG having an effect, environment/state matrix, or approved exclusion; every PRV finding in 0C having a reciprocal SEC record; every effective-denominator Disabled SRC, CMP, and UC having exactly one canonical CAP with disabled/dormant capabilities indexed, cross-linked, risk-described, and having an Approved target decision, and overlaps having explicit relationships and approved decisions; every material fault class having an FLT record and the consolidated fault/resilience topology being complete; every persona profile having a unique legal PRF ID, being atomically synchronized with its 0A persona-prefix row, having the complete §8.4 dependency closure validated against exact versions/fingerprints without mixed snapshots, being semantically fingerprinted, B-CONFIRMED, and covered by a valid exact-PRF-binding CNF; every synthesis, PRF profile, and HBK handbook semantic payload exactly matching its typed-ID entry in the candidate payload manifest and every certification envelope being B-CONFIRMED, referencing a valid CNF for its exact typed record ID, SEMANTIC-RECORD-VERSION, and SEMANTIC-CONTENT-FINGERPRINT, and being neither draft, populated-only, nor stale; every reconstruction-affecting scope, target, capability, and migration decision having an approval envelope with Approval Status: Approved bound to the exact candidate decision-content version/fingerprint, while Proposed decisions cannot pass Exit E; the reviewed equivalence suite mapping every UC contract, BR/invariant property, eligible SM transition, external IF baseline, and disabled-capability safeguard to at least one test or exact approved exclusion; all pending synthesis references and GAP records reconciled with no orphan/dangling/tombstone violation remaining; all material contradictions and decision-required items resolved or represented by approved scope exclusions that preserve and disclose exact residual risk; the system namespace, protocol version, environments, and snapshot present and exactly equal across the Exit A report, Exit E Candidate Report, Exit E Content-Readiness Report, scope certificate, every snapshot-bearing package member, and outer bundle manifest, with the source inventory denominator, export projections, claims, candidate payload manifest, confirmations, and remaining package members fingerprint-pinned to that identity under §4.1.2 and the acyclic exclusions in §15.4, failing closed on missing identity or snapshot mixing; human architect/domain confirmation recorded in 18-CONFIRMATIONS.md, with B-CONFIRMED never machine-assigned and confirmation envelope attachment not altering any semantic version/fingerprint; and at step 6 the outer bundle manifest payload directly declaring the exact system/protocol/environment/snapshot identity, final check-registry fingerprint, and pre-signature input-set fingerprint, including every package member other than itself, including the Exit E Content-Readiness Report and scope certificate, and containing EXIT-E-STATUS: Passed, with all listed package-member file fingerprints, required certification/approval/CNF envelope bindings, cross-artifact snapshot identities, and required final checks recomputing and validating and its required human signature envelope validating under §15.5. The signed outer manifest is the authoritative non-cyclic Exit E completion attestation. No post-step-6 authoritative Exit E artifact is created. There are no row-count shortcuts: file exists, section exists, or zero rows does not satisfy completeness without denominator-backed proof.

## 15.1.1. Universal deterministic packaging rules

The five packaging artifact types are EXIT-E-CANDIDATE-REPORT, CANDIDATE-PAYLOAD-MANIFEST, EXIT-E-CONTENT-READINESS-REPORT, SCOPE-CERTIFICATE, and OUTER-BUNDLE-MANIFEST. Each MUST instantiate the generic §4.1.2 payload and envelope boundaries and its exact ordered schema.

### Rule: packaging-payload-rules

For every packaging payload: payload fields appear exactly once and in the declared order, and an unknown, duplicate, reordered, or omitted required payload field is a structural error; a declared optional or empty scalar is serialized as the exact literal None, never omitted, blank, null-like, or represented by prose; an empty table retains its declared header and separator rows and contains no data row; every normalized path uses §4.1 with normalization before sorting, hashing, or binding; every table is sorted by its schema-declared tuple after normalization using ascending bytewise comparison of the exact preserved Unicode code points serialized as UTF-8, and duplicate sort-key tuples are invalid unless the table schema explicitly permits them; an ID list not represented as a table is sorted by (record type, typed ID) and a path/member list by (normalized path, artifact type), with no source discovery order semantically significant; semantically equal inputs presented in different source orders MUST serialize to one canonical payload and one payload fingerprint; the HASH-DOMAIN identity is exactly artifact-type + "|" + normalized-relative-path, is included in the preimage, and cannot contain the payload fingerprint, while the POST-HASH-ARTIFACT-INSTANCE-BINDING is computed only after hashing and is never part of its own preimage. The declared table sort keys are: candidate validation or readiness checks by (Check ID); blockers by (Blocker ID); semantic counts by (Record Type); decision counts by (Decision Type); semantic bindings by (Record Type, Typed Record ID); decision bindings by (DEC ID); confirmed bindings by (Record Type, Typed Record ID); approved decision bindings by (DEC ID); included scope dimensions by (Dimension, Value); source counts by (SRC Kind); exclusions by (DEC ID, Excluded Unit ID); residual risks by (Risk ID); and every package-member table, including the outer manifest, by (Normalized Path, Artifact Type). Fields described as a deterministic set or tuple use the corresponding rule above.

### Rule: exit-e-check-registry

The required Exit E report check registry is EXIT-E-CHECKS-v1. Its canonical registry fingerprint is sha256 over UTF-8 EXIT-E-CHECK-REGISTRY|EXIT-E-CHECKS-v1 plus one LF and the §4.1.2 canonical serialization of the exact ordered registry table. A report's Check Results.Check ID set MUST equal exactly the rows for its stage: no missing, duplicate, unknown, or other-stage ID is permitted. Every required row is serialized and has Result: Pass | Fail; this registry defines no Not-Applicable result because each check validates its own applicability denominator and evidence-backed zero domains. A missing or malformed row is a structural report failure, not a pass or an omitted/unknown result. Ready-For-Human-Review or Ready-For-Certificate requires exact-set equality, every row Pass, every row's evidence-set fingerprint to recompute, and an empty blockers table, so a header-only empty check table can never satisfy either report. Any failed check requires at least one corresponding blocker; a blocker cannot compensate for a missing check.

### Enum: EXIT-E-CANDIDATE-CHECKS

- CANDIDATE-01-EXIT-A
- CANDIDATE-02-SOURCE-CLOSURE
- CANDIDATE-03-SYNTHESIS-DOMAINS
- CANDIDATE-04-HANDBOOK-QUALITY
- CANDIDATE-05-REFERENTIAL-INTEGRITY
- CANDIDATE-06-LIFECYCLE
- CANDIDATE-07-USE-CASES
- CANDIDATE-08-CONFIGURATION
- CANDIDATE-09-PRIVACY-SECURITY
- CANDIDATE-10-DISABLED-CAPABILITIES
- CANDIDATE-11-FAULTS
- CANDIDATE-12-PROFILES
- CANDIDATE-13-PENDING-AND-CONTRADICTIONS
- CANDIDATE-14-SNAPSHOT-INPUTS

### Enum: EXIT-E-CONTENT-READINESS-CHECKS

- READINESS-01-CANDIDATE-IMMUTABILITY
- READINESS-02-CONFIRMATIONS
- READINESS-03-DECISION-APPROVALS
- READINESS-04-EQUIVALENCE
- READINESS-05-HANDBOOK-READABILITY
- READINESS-06-NO-PENDING-OR-CONTRADICTIONS
- READINESS-07-ENVELOPE-INTEGRITY
- READINESS-08-PREPACKAGE-MEMBERS
- READINESS-09-SNAPSHOT-CONSISTENCY
- READINESS-10-RESIDUAL-RISK
- READINESS-11-SEQUENCE

### Enum: FINAL-CHECK

- FINAL-01-OUTER-PAYLOAD
- FINAL-02-SNAPSHOT
- FINAL-03-MEMBER-SET
- FINAL-04-MEMBER-FINGERPRINTS
- FINAL-05-ENVELOPES
- FINAL-06-CHAIN
- FINAL-07-REGISTRY-AND-INPUTS
- FINAL-08-CONDITIONS
- FINAL-09-SEQUENCE

### Field: ARTIFACT-SIGNATURE-ENVELOPE

| Field | Type | Required |
| :-- | :-- | :-- |
| signature-signatories-timestamps | string | true |
| envelope-audit | string | true |

### Field: SEMANTIC-RECORD-COUNT

| Field | Type | Required |
| :-- | :-- | :-- |
| record-type | string | true |
| count | int | true |

### Field: DECISION-CONTENT-COUNT

| Field | Type | Required |
| :-- | :-- | :-- |
| decision-type | string | true |
| count | int | true |

### Field: CHECK-RESULT

| Field | Type | Required |
| :-- | :-- | :-- |
| check-id | string | true |
| requirement | string | true |
| result | string | true |
| evidence-fingerprint | string | true |

### Field: BLOCKER

| Field | Type | Required |
| :-- | :-- | :-- |
| blocker-id | string | true |
| affected-ids | string | true |
| description | string | true |
| evidence-fingerprint | string | true |

## 15.1.2. Exit E Candidate Report ordered payload schema

The Exit E Candidate Report is emitted before the candidate payload manifest. Its exact payload field order is declared below.

### Rule: candidate-report-ordering

Check Results sorts by (Check ID) and uses Pass | Fail; Blockers sorts by (Blocker ID). Ready-For-Human-Review requires exact CANDIDATE registry-set equality, every required check passing with a valid evidence-set fingerprint, and an empty blockers table. A failed check requires a blocker; a missing/duplicate/unknown check makes the report structurally invalid. The report MUST NOT reference or include the later candidate payload manifest, confirmations or CNFs, confirmation envelopes, Exit E Content-Readiness Report, scope certificate, or outer bundle manifest. This ordering prevents a candidate report/manifest hash cycle.

### Field: EXIT-E-CANDIDATE-REPORT

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| artifact-path | string | true |
| protocol-artifact-schema-version | string | true |
| system-protocol-snapshot | string | true |
| required-check-registry-fingerprint | string | true |
| exit-a-report-fingerprint | string | true |
| deterministic-candidate-input-set-fingerprint | string | true |
| coverage-validation-summary | string | true |
| traceability-validation-summary | string | true |
| handbook-validation-summary | string | true |
| semantic-record-counts | string | true |
| decision-content-counts | string | true |
| check-results | string | true |
| blockers | string | true |
| candidate-result | string | true |
| artifact-payload-fingerprint | string | true |

## 15.1.3. Candidate Payload Manifest ordered payload schema

The candidate payload manifest is emitted only after the immutable candidate report and references it one-way. Its exact payload field order is declared below.

### Rule: manifest-ordering

Semantic Bindings sorts by (Record Type, Typed Record ID) and DEC Content Bindings by (DEC ID). Each count MUST equal its table's data-row count. The manifest MUST NOT include future DEC approval envelopes, CNFs, confirmation envelopes, the Exit E Content-Readiness Report, scope certificate, or outer bundle manifest. The candidate report does not reference this manifest; only this manifest references the earlier report.

### Field: SEMANTIC-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| record-type | string | true |
| typed-record-id | string | true |
| semantic-record-version | string | true |
| semantic-content-fingerprint | string | true |
| dependency-set-fingerprint | string | true |
| normalized-path | string | true |

### Field: DEC-CONTENT-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| dec-id | string | true |
| decision-content-version | string | true |
| decision-content-fingerprint | string | true |

### Field: CANDIDATE-PAYLOAD-MANIFEST

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| artifact-path | string | true |
| protocol-artifact-schema-version | string | true |
| exit-e-candidate-report-fingerprint | string | true |
| semantic-bindings | string | true |
| dec-content-bindings | string | true |
| semantic-binding-count | int | true |
| dec-content-binding-count | int | true |
| artifact-payload-fingerprint | string | true |

## 15.1.4. Exit E Content-Readiness Report ordered payload schema

The stage-4 Exit E Content-Readiness Report is a pre-certificate artifact. Its exact payload field order is declared below.

### Rule: content-readiness-ordering

Confirmed bindings sort by (Record Type, Typed Record ID), approved decisions by (DEC ID), pre-certificate members by (Normalized Path, Artifact Type), checks by (Check ID), and blockers by (Blocker ID). Check results use Pass | Fail. Ready-For-Certificate requires exact CONTENT-READINESS registry-set equality, every required check passing with a valid evidence-set fingerprint, and an empty blockers table. A failed check requires a blocker; a missing/duplicate/unknown check makes the report structurally invalid. The payload MUST omit the outer-manifest Exit E completion field and completion literal and MUST NOT reference or depend on the later scope certificate or outer bundle manifest. It cannot pass Exit E; Exit E remains Pending until the signed outer-manifest step 6.

### Field: CONFIRMED-SEMANTIC-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| record-type | string | true |
| typed-record-id | string | true |
| semantic-record-version | string | true |
| semantic-content-fingerprint | string | true |
| certification-envelope-binding | string | true |
| semantic-containing-path | string | true |
| semantic-package-member-file-fingerprint | string | true |
| cnf-id | string | true |
| cnf-payload-fingerprint | string | true |
| cnf-signature-envelope-binding | string | true |
| cnf-containing-path | string | true |
| cnf-package-member-file-fingerprint | string | true |

### Field: APPROVED-DEC-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| dec-id | string | true |
| decision-content-version | string | true |
| decision-content-fingerprint | string | true |
| approval-envelope-binding | string | true |
| dec-containing-path | string | true |
| dec-package-member-file-fingerprint | string | true |

### Field: PACKAGE-MEMBER

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| normalized-path | string | true |
| package-member-file-fingerprint | string | true |

### Field: EXIT-E-CONTENT-READINESS-REPORT

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| artifact-path | string | true |
| protocol-artifact-schema-version | string | true |
| system-protocol-snapshot | string | true |
| required-check-registry-fingerprint | string | true |
| exit-a-report-fingerprint | string | true |
| exit-e-candidate-report-fingerprint | string | true |
| candidate-payload-manifest-fingerprint | string | true |
| confirmed-semantic-bindings | string | true |
| approved-dec-bindings | string | true |
| pre-certificate-package-members | string | true |
| check-results | string | true |
| exit-e-prepackage-state | string | true |
| content-readiness-result | string | true |
| blockers | string | true |
| artifact-payload-fingerprint | string | true |

## 15.2. Scope certificate

The scope certificate binds scope, exclusions, risks, and the prior certification chain. Its exact payload field order is declared below.

### Rule: scope-certificate-ordering

Included Scope Dimensions sorts by (Dimension, Value); source counts by (SRC Kind); approved exclusions by (DEC ID, Excluded Unit ID); confirmed bindings by (Record Type, Typed Record ID); and residual risks by (Risk ID). Every count and arithmetic relation MUST reproduce the cited denominator/reconciliation reports. The scope certificate is emitted only after the Exit E Content-Readiness Report. It references that immutable report and the candidate/confirmation chain; the content-readiness report does not reference or hash the later certificate. The certificate uses the generic §4.1.2 artifact boundaries and exclusions; all certificate scope, risk, confirmed typed-record bindings, and report-reference fields remain hashed.

### Field: INCLUDED-SCOPE-DIMENSION

| Field | Type | Required |
| :-- | :-- | :-- |
| dimension | string | true |
| value | string | true |

### Field: SOURCE-COUNT

| Field | Type | Required |
| :-- | :-- | :-- |
| src-kind | string | true |
| gross-discovered | int | true |
| approved-excluded | int | true |
| effective-denominator | int | true |
| covered | int | true |

### Field: APPROVED-EXCLUSION-RESIDUAL-RISK

| Field | Type | Required |
| :-- | :-- | :-- |
| dec-id | string | true |
| excluded-unit-id | string | true |
| decision-content-version | string | true |
| decision-content-fingerprint | string | true |
| approval-envelope-binding | string | true |
| risk-id | string | true |
| exact-residual-risk | string | true |

### Field: CONFIRMED-SEMANTIC-PAYLOAD-BINDING

| Field | Type | Required |
| :-- | :-- | :-- |
| record-type | string | true |
| typed-record-id | string | true |
| semantic-record-version | string | true |
| semantic-content-fingerprint | string | true |
| cnf-id | string | true |
| cnf-payload-fingerprint | string | true |

### Field: KNOWN-RESIDUAL-RISK

| Field | Type | Required |
| :-- | :-- | :-- |
| risk-id | string | true |
| affected-ids | string | true |
| source-dec-id | string | true |
| exact-residual-risk | string | true |
| disclosure-path | string | true |

### Field: SCOPE-CERTIFICATE

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| artifact-path | string | true |
| protocol-artifact-schema-version | string | true |
| system-protocol-snapshot | string | true |
| included-scope-dimensions | string | true |
| source-counts-by-kind | string | true |
| approved-exclusions-and-residual-risks | string | true |
| export-reconciliation | string | true |
| exit-a-report-fingerprint | string | true |
| exit-e-candidate-report-fingerprint | string | true |
| candidate-payload-manifest-fingerprint | string | true |
| exit-e-content-readiness-report-fingerprint | string | true |
| confirmed-semantic-payload-bindings | string | true |
| known-residual-risks | string | true |
| artifact-payload-fingerprint | string | true |

## 15.3. Package contents

The certified package contains the handbook, structured synthesis catalogs, assurance records and validation summaries, persona profiles, decisions, equivalence suite, scope certificate, packaging artifacts, fingerprint manifests, validation reports, and a known-gap and risk record; raw secrets and private material are excluded.

### Rule: certified-package-contents

The certified package contains the human reconstruction handbook; structured 90-96 synthesis catalogs; assurance plane records and deterministic validation summaries; persona profiles and necessary forensic reference records; the decision log; the equivalence/acceptance suite; the scope certificate; the Exit E Candidate Report, candidate payload manifest, Exit E Content-Readiness Report, and signed outer bundle manifest; source, projection, artifact, and semantic fingerprint manifests; Exit A, Exit E Candidate, and content-readiness validation reports; and a known-gap/risk record that MUST disclose every approved exclusion and its residual risks even though excluded from the effective denominator, MUST contain no undisclosed or unapproved gap, and MUST NOT claim a gap absent from both the gross discovered population and approved risk record. Raw secrets, raw private exports, private matching hints, and private entity indexes are excluded.

## 15.4. Acyclic certification sequence and hash domains

Certification MUST follow the declared order; no later artifact is an input to an earlier hash.

### Rule: certification-steps

Step 1: Final-Synthesis and Profile Synchronization produce complete B-POPULATED synthesis and PRF profile payloads, Reconstruction-Handoff projects the complete handbook as B-POPULATED HBK records, runs the deterministic pre-candidate non-vacuity/link/diagram and synchronization checks in §13.4, and computes every typed record's SEMANTIC-CONTENT-FINGERPRINT; human readability review does not occur here. Step 2: Validation/Gate validates Exit A, all populated typed payloads, proposed decision content, references, coverage, and handbook quality, then emits an immutable Exit E Candidate Report and candidate payload manifest using the generic §4.1.2 artifact schema; the manifest lists every candidate synthesis, PRF, and HBK record ID, semantic record version, semantic-content fingerprint, dependency fingerprint, and required decision-content version/fingerprint binding; neither candidate artifact includes future DEC approval envelopes, CNFs, confirmation envelopes, the Exit E Content-Readiness Report, certificate, or outer manifest. Step 3: authorized humans approve required DEC content by attaching approval envelopes bound to the exact candidate DECISION-CONTENT-VERSION/DECISION-CONTENT-FINGERPRINT, then create CNF records against exact candidate typed-record bindings, including PRF/HBK IDs where applicable; these attachments do not alter candidate-bound content fingerprints and do not by themselves restart candidacy; every CNF that confirms candidate HBK content records the §13.4 handbook readability result and the combined passing CNFs MUST cover every candidate HBK binding; the authorized confirmation action attaches B-CONFIRMED/CNF envelopes without regenerating or changing semantic payload; DEC and CNF hashes use §4.1.2 and their exact record-specific envelope markers. Step 4: Validation/Gate recomputes candidate semantic fingerprints, verifies unchanged payloads, recomputes and validates every certification/approval/CNF envelope binding and finalized package-member file fingerprint available at this stage, validates every content/confirmation/package-member condition available before certificate and outer-manifest creation, and emits the immutable Exit E Content-Readiness Report; its hashed payload may reference Exit A, the Exit E Candidate Report/candidate payload manifest, DEC/CNF records, exact envelope bindings, and all validated pre-certificate package-member file fingerprints available at this stage; it MUST state that Exit E is still Pending, MUST NOT contain EXIT-E-STATUS: Passed, and MUST NOT reference or depend on the later scope certificate or outer bundle manifest/fingerprint. Step 5: Reconstruction-Handoff emits the scope certificate referencing the immutable Exit E Content-Readiness Report fingerprint and prior chain; creating it MUST NOT change the content-readiness report; Exit E remains Pending. Step 6: Reconstruction-Handoff emits the outer manifest, whose hashed payload directly declares the pinned System / Protocol / Snapshot, the final check-registry identity/fingerprint, and the pre-signature authoritative input-set fingerprint; lists every final package member other than the manifest itself and each external package-member file fingerprint, including the candidate artifacts, DEC/CNF records and their finalized containing files, confirmed semantic records and their finalized containing files, Exit E Content-Readiness Report, and scope certificate; includes the literal field EXIT-E-STATUS: Passed; and does not list itself or its own digest as a package member. Validation/Gate requires exact system/protocol/environment/snapshot equality across Exit A, candidate report, content-readiness report, scope certificate, every snapshot-bearing package member, and the outer payload; recomputes every package-member file fingerprint and required certification/approval/CNF envelope binding; verifies the complete member set, generic payload/envelope boundaries, the content-readiness and certificate references, pre-signature input-set fingerprint, final check-registry fingerprint, and all §15.1 conditions; and fails closed on missing identity, mixed snapshots, or a validly hashed member from another snapshot. Authorized humans then sign the outer artifact envelope, and the conforming implementation performs the §15.5 deterministic final-bundle verification over the completed signature envelope. Only when every required final check passes does Exit E transition from Pending to Passed. The signed outer manifest is the authoritative non-cyclic Exit E completion attestation and its payload digest is the bundle fingerprint. No post-step-6 authoritative completion report or attestation is emitted.

### Rule: pre-signature-input-set

The outer payload's Pre-Signature Validation Input-Set Fingerprint is computed before outer payload hashing and signing. Its complete preimage is UTF-8 EXIT-E-PRE-SIGNATURE-INPUTS|, the canonically serialized outer System / Protocol / Snapshot, one LF, the EXIT-E-CHECKS-v1 identity and fingerprint, one LF, the EXIT-E-FINAL-CHECKS-v1 identity and fingerprint, one LF, the Exit E Content-Readiness Report artifact-payload fingerprint, one LF, the scope-certificate artifact-payload fingerprint, one LF, and the §4.1.2 canonical serialization of exactly the outer Package Members table. Each member row MUST use PACKAGE-MEMBER-FILE and the complete finalized path-bound fingerprint for that normalized path/artifact type. The table MUST equal the complete final member set excluding only the outer manifest itself; a missing, extra, duplicate, differently domained, or mixed-snapshot binding changes or invalidates the preimage. The preimage excludes the outer payload, its fingerprint carrier, and its signature envelope, so it is acyclic.

### Rule: outer-manifest-ordering

Package Members MUST be sorted by (Normalized Path, Artifact Type), MUST contain every final package member except the outer manifest itself, and MUST reject duplicate normalized paths or duplicate member bindings. The declared outer-manifest field order is exact under §15.1.1.

### Rule: exclusion-rule

The normative semantic/payload hash profile and exact boundaries are §4.1.2. Packaging payload fingerprints exclude only their own ARTIFACT-PAYLOAD-FINGERPRINT carrier and matching generic artifact digest/signature envelope. Synthesis/PRF/HBK semantic records exclude exactly the certification envelope in §4.1.1; DEC excludes exactly the decision-approval envelope in §5.5 and its own carrier; CNF excludes exactly the CNF digest/signature envelope in §5.7 and its own carrier. Those exclusions do not apply to canonical envelope fingerprints or external package-member file fingerprints: envelope fingerprints include every field within their exact envelope, and package-member file fingerprints include the complete finalized member. Cross-document references, content timestamps, package-member lists, EXIT-E-STATUS, and all other payload fields remain hashed. No artifact may exclude a later artifact reference merely to mask a cycle; such a reference is forbidden by the sequence.

### Rule: restart-rule

A semantic payload, dependency, candidate decision-content binding, or candidate-manifest change restarts at step 1 or 2 as applicable and invalidates downstream CNFs. A DEC review-input or decision-content change restarts at step 2; attaching a matching approval envelope occurs at step 3 and does not restart, while a mismatched envelope is rejected. A CNF or confirmation-envelope defect restarts at step 3. An Exit E Content-Readiness Report payload defect restarts at step 4. A certificate-only payload defect restarts at step 5. An outer-manifest membership, status, fingerprint, order, validation, or signature defect restarts at step 6 and leaves Exit E Pending. Replacing any finalized member envelope changes its envelope binding and package-member file fingerprint; it restarts from the stage that owns that envelope and always invalidates/requires regeneration and re-signing of the step-6 outer manifest. Replacing only the outer manifest's own signature envelope with an unchanged outer payload digest restarts step 6 only and requires revalidation/re-signing. Any change that reaches backward into an earlier artifact restarts from the earliest affected step. There is no post-step-6 authoritative artifact whose hash could create a cycle.

### Field: OUTER-SIGNATURE-ENVELOPE

| Field | Type | Required |
| :-- | :-- | :-- |
| required-human-signatories-signature-timestamps | string | true |
| envelope-audit | string | true |

### Field: OUTER-BUNDLE-MANIFEST

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-type | string | true |
| artifact-path | string | true |
| protocol-artifact-schema-version | string | true |
| system-protocol-snapshot | string | true |
| final-verification-check-registry-fingerprint | string | true |
| pre-signature-validation-input-set-fingerprint | string | true |
| exit-e-content-readiness-report-fingerprint | string | true |
| scope-certificate-fingerprint | string | true |
| package-members | string | true |
| exit-e-status | string | true |
| artifact-payload-fingerprint | string | true |

## 15.5. Reproducible final-bundle verification

EXIT-E-FINAL-CHECKS-v1 is the exact final verification registry.

### Rule: final-bundle-verification-operation

The canonical fingerprint of EXIT-E-FINAL-CHECKS-v1 is sha256 over UTF-8 EXIT-E-FINAL-CHECK-REGISTRY|EXIT-E-FINAL-CHECKS-v1, one LF, and the §4.1.2 canonical serialization of the final check registry table. A conforming runtime MUST expose a deterministic verify final bundle operation over the complete signed bundle. The operation recomputes the exact registry and check results from bundle bytes; it does not trust the outer payload's Passed literal, stored fingerprints, or executor-supplied values. Missing, duplicate, unknown, or failed checks yield overall Fail. FINAL-09-SEQUENCE verifies only bundle-observable members, references, and recorded prior-stage states. The prohibition on creating later authoritative completion state is separately enforced by the runtime's post-Exit-E transaction rules and cannot be inferred from standalone bundle bytes.

### Rule: final-verification-receipt

Every execution records a final verification receipt with the declared schema. The result fingerprint covers every preceding field except Validator Implementation Version / Timestamp and its own carrier under §4.1.2 with domain prefix FINAL-VERIFY-RESULT|. The receipt is a NON-AUTHORITATIVE-DERIVATIVE, MUST remain outside the certified package and authoritative .extracted protocol state, MUST NOT be a package member or input to Exit E, and cannot change the outer payload digest, signature, or gate state. It records which checks the original runtime executed; any third party can discard it and reproduce the same deterministic result from the signed bundle. Creating, deleting, or regenerating a receipt is not a post-step-6 completion artifact or attestation.

### Field: FINAL-CHECK-RESULT

| Field | Type | Required |
| :-- | :-- | :-- |
| check-id | string | true |
| result | string | true |
| offending-bindings | string | true |

### Field: FINAL-VERIFICATION-RECEIPT

| Field | Type | Required |
| :-- | :-- | :-- |
| protocol-artifact-check-registry-versions-and-fingerprints | string | true |
| normalized-bundle-identity | string | true |
| outer-payload-canonical-outer-envelope-fingerprints | string | true |
| system-protocol-environment-snapshot | string | true |
| package-member-input-bindings-and-input-set-fingerprint | string | true |
| check-results | string | true |
| overall-result | string | true |
| validator-implementation-version-timestamp | string | true |
| deterministic-result-fingerprint | string | true |

## 17.1. Required order

Migration follows a fixed order from snapshot freezing through inventory, ID allocation, wildcard expansion, atomic splitting, track classification, persona normalization, explosion reconciliation, claim conversion, frontier reconstruction, human-required conversion, traceability, draft marking, gap exposure, reconciliation, and only then gated closure and Exit E.

### Rule: required-migration-order

Migration follows this order: freeze and fingerprint the v3.17 workspace and source snapshot; create the v4 source inventory denominator before interpreting old coverage; allocate deterministic IDs for source, components, claims, and synthesis candidates while preserving old headings as aliases; expand wildcard clusters and exclusions into concrete coordinates; split routine families, modules, serialized containers, and aggregated blocks into universal atomic records while retaining non-promotable grouping containers for navigation; inventory and classify Normal, Shadow/Conditional, and Disabled tracks separately by environment/snapshot; apply deterministic persona-prefix uniqueness and entry ownership and create shared-entry records where intentional; move each legacy unprefixed persona directory to its registered persona-prefix-persona-slug basename as a proven path move while preserving record IDs and history, recording the prior path as an alias, updating every path reference and invocation target, and recomputing affected file-transport fingerprints without changing prefix-derived PRF identity; re-run logical explosion and reconcile every normalized export coordinate; convert block-level evidence into one-fact CLM records, splitting mixed evidence and deriving BLOCK-CONFIDENCE-RANK plus BLOCK-EVIDENCE-PROFILE mechanically; build traversal frontier records from depth boundaries, unresolved dependencies, missing exports, and old tickets; convert direct human-required records into migration findings and reopen them at the static-investigation/probe stage unless an existing authorized exclusion proves the scope disposition; populate traceability, reciprocal references, contradiction records, and coverage arithmetic; mark old synthesis-like outputs provisional B-DRAFT without inheriting confirmation; expose all gaps before enabling strict validators and never backfill evidence merely to satisfy a gate; run Sequential and Cross-Reference Reconciliation, then evidence-backed sweep validation; enable Composite Exit A validators only after denominator and frontier stabilization; and run Final Synthesis, handbook projection, human review, and Exit E.

## 17.2. Migration guarantees

Migration preserves ticket history, evidence anchors, fingerprints, old identifiers as aliases, promotion and depromotion history, and uncertainty, and never converts disabled to dead, terminal tickets to covered source, inferred evidence to direct, or old synthesis to confirmed.

### Rule: migration-preservation

Migration MUST preserve ticket history, evidence anchors, fingerprints, old identifiers as aliases, promotion/depromotion history, and uncertainty. It MUST NOT convert disabled to dead, terminal tickets to covered source, inferred evidence to direct, or old synthesis to confirmed. The Draft Principle remains in force; B-CONFIRMED always requires fresh human confirmation against the pinned v4 bundle.

## 17.3. v3 final-artifact compatibility

v4 supersedes the v3 final workflow atlas and human audit queue semantically; legacy filenames may be retained only as non-authoritative compatibility projections that identify their upstream IDs, versions, fingerprints, generation time, schema, and staleness.

### Rule: v3-supersession

v4 supersedes v3 FINAL-WORKFLOW-ATLAS.md semantically with 93-USE-CASES.md plus the corresponding handbook use-case, background-processing, and architecture chapters. v4 supersedes v3 FINAL-HUMAN-AUDIT-QUEUE.md semantically with the ticket/contradiction/decision/confirmation/known-gap and gate artifacts in the Forensic, Assurance, and Handbook planes. Neither v3 file is authoritative in a v4 workspace and neither satisfies an Exit A or Exit E requirement.

### Rule: compatibility-projection

A migration MAY retain either filename only as a NON-AUTHORITATIVE-COMPATIBILITY-PROJECTION for legacy consumers. Such a projection MUST identify its v4 upstream IDs, semantic versions/fingerprints, generation time, projection schema, and staleness state; MUST contain no independent claims, approvals, queue state, or decisions; and MUST be excluded from authoritative denominator and certification inputs except as a listed derivative package member. Conflicts resolve in favor of canonical v4 records.

## 19.1. Conforming implementation capabilities

A conforming implementation provides every deterministic capability and authority boundary declared for identity, scope, hashing, coverage, gates, packaging, handbook validation, and final verification.

### Rule: conforming-implementation-capabilities

A conforming implementation enforces exact persona/cluster/track/POV/file/ledger isolation in identity headers; lets executors propose semantic content only while the runtime independently computes or reproduces every protocol-defined ID, canonical form, fingerprint, scope/stale check, FSM transition, count, completeness result, gate result, and package verification before commit; binds cold-resume semantic assessment to runtime-recomputed identity/scope/write-right/stale checks and the mandatory 0G no-mutation or commit entry; treats personas/ as the sole persona container, validates each registered persona basename as exact persona-prefix-persona-slug with registry/header/ticket/PRF prefix agreement, validates canonical shared records under personas/_shared/, rejects malformed/unregistered persona directories and _shared as a prefix, slug, basename, or persona row, and never counts the reserved directory as a persona; restricts multi-file writes to enumerated modes and enforces the personas/{persona-directory}/ versus personas/_shared/ ownership boundary; implements the ID generator, collision extension, aliases, versions, tombstones, and COV/CND/MOD/PRF/HBK types and canonical coordinates; implements HBK paths and assigned anchors per §4.1 with renderer-generated heading anchors never used as identity inputs; implements the §4.1.2 canonical hash profile, per-record and file-transport fingerprints, exact carrier/envelope exclusions, and acyclic package hashes; uses packaging hash-domain identity exactly (artifact type, normalized path) with the payload fingerprint absent from its own preimage and the artifact-instance binding constructed only after hashing; gives certification, DEC-approval, and CNF-signature envelopes identity/version-bound canonical fingerprints and covers complete finalized files with package-member file fingerprints distinct from envelope-excluding transport fingerprints; dispatches every finite enum by its complete §5.1 Markdown-AST schema path, failing closed on label-only dispatch, undeclared paths, wrong-family tokens, and invalid traversal N/A coupling; covers every required source kind and can prove zero domains; creates exactly one FRT for every discovered traversal boundary including immediately terminal boundaries and depth stops; rejects mixed CLM evidence and derives non-lossy BLOCK-CONFIDENCE-RANK plus BLOCK-EVIDENCE-PROFILE; fails helper publication closed and keeps raw/private state outside the target repository and authoritative .extracted workspace in operator-controlled storage; reconciles exports as one row per normalized coordinate; forbids direct human-required ticket creation, validates the closed escalation-reason taxonomy, and permits non-runtime Human Hatch escalation only after bounded static exhaustion and a reason-consistent probe inapplicability assessment; preserves gap/partial coverage on human-required terminalization distinct from exact approved exclusion; fingerprints and sanitizes probe logs; merges concurrent persona buffers including promotion requests and unmapped discoveries sequentially, deterministically, and atomically before later Discovery, with only Sequential Reconciliation moving canonical blocks between personas/_shared/ and persona POV files; deduplicates synthesis GAP records deterministically and creates required ticket/frontier work only through Sequential Reconciliation; validates promotion/depromotion and disabled/dead distinctions; reproduces expected/mapped/frontier/ticket counts in sweep records; runs coverage arithmetic per kind, track, environment, and snapshot; rejects any linked unresolved possible caller/dynamic reference in dead-code validation while confining the effect to the linked dependency closure; carries provenance on semantic applicability/lifecycle classifications with Unknown fail-inclusive and unable to silently justify omission or closure; runs reciprocal, orphan, dangling, ownership, and persona-purity checks; propagates targeted staleness along DERIVED-FROM impact edges; loads and validates the complete direct/applicable and transitive §8.4 dependency closure in Profile Synchronization and reuses the identical candidate-bound closure for confirmation; separates synthesis gap buffering from controlled ticket creation; requires human confirmation for every B-confirmed block binding exact semantic versions/fingerprints without changing payloads; enforces exact field order, required/optional/empty representation, schema-declared row sorting, stage-specific exact required-check sets, evidence-set fingerprints, and rejection of unknown, duplicate, reordered, omitted, or stage-ineligible fields/checks across all five packaging artifacts; follows §15.4 for the Exit E Candidate Report, candidate payload manifest, Exit E Content-Readiness Report, scope certificate, and signed outer manifest with direct snapshot equality across the chain and only the validated signed outer payload setting EXIT-E-STATUS: Passed; executes the exact final registry and emits only a reproducible non-authoritative receipt outside package/protocol state; passes deterministic §13.4 handbook link/diagram/evidence-display checks before candidacy with candidate-bound human readability review recorded only through complete passing CNFs after candidacy; and produces deterministic fingerprinted validation summaries.

## 19.2. Normative conformance cases

A conforming implementation satisfies every applicable positive and negative normative conformance case and deterministically produces the specified acceptance or rejection.

### Rule: normative-conformance-cases

A conforming implementation MUST satisfy every applicable positive and negative case in the declared case set and deterministically produce the specified acceptance or rejection. These are protocol-level behavioral cases, not required files or test-runner inputs; their storage, serialization, and execution mechanism are implementation-specific. Failure of an applicable case defeats a claim of protocol conformance. Workspace gate execution is governed independently by the authoritative inputs and gate checks defined elsewhere in this protocol.

### Enum: CONFORMANCE-CASE

- Persona workspace layout and ownership
- Context-qualified finite enums and registry completeness
- Invocation ownership and cold resume
- Ticket escalation and honest unresolved coverage
- Dead-code and semantic-predicate closure
- Deterministic iteration accounting
- PRF identity and semantic hashing
- HBK identity and path/anchor normalization
- Profile Synchronization closure
- Record versus artifact hashing
- Canonicalization and exact exclusions
- Final envelope and package-member integrity
- Packaging schemas, gate completeness, and deterministic ordering
- Acyclic, snapshot-consistent Exit E and final verification

## 19.3. Corpus acceptance

Corpus acceptance requires a green Composite Exit A with no hidden gaps, atomic and terminally disposed sources, mapped or exactly excluded units, claim-backed traversal cells, separated but cross-linked active and disabled behavior, claim-backed synthesis records, and a step-6 Exit E transition.

### Rule: corpus-acceptance-requirements

Corpus acceptance requires: Composite Exit A green with no hidden partial/gap denominator; every in-scope source unit atomic, covered, and terminally disposed; every normalized unit mapped or exactly excluded; every required traversal-matrix cell evidence-backed closed and every N/A cell claim-backed; active and disabled behavior visibly separated but fully cross-linked; every MOD/ER/REL/DR/BR/UC/IF/DEP/CFG/SCHED/NFR/SEC/SM/CAP/PRV/FLT block claim-backed as applicable with purely navigational indexes carrying no semantic status/evidence; every BR with one authoritative owner and consistent optional mirrors; every external IF with a compatibility baseline and migration rule; every eligible lifecycle with a complete SM; every UC able to generate a reviewable acceptance-test skeleton; every CFG with a behavioral effect or approved exclusion; every PRV finding with a SEC mirror and every material failure class with an FLT record; every effective-denominator Disabled SRC/CMP/UC with exactly one canonical CAP and every disabled capability with risk and an approved Restore/Preserve Dormant/Redesign/ Retire decision; no synthesis or handbook block stale, provisional, or pending-reference; human readers able to navigate core behavior without reading raw ID catalogs; and Exit E transitioned at step 6 only with the scope certificate included and the validated signed outer manifest hash payload authoritatively containing EXIT-E-STATUS: Passed.

## 20. Compact Artifact Ownership Matrix

The compact artifact ownership matrix records the semantic purpose, authoritative plane, authorized writers, and key non-writers for every artifact and record.

### Field: ARTIFACT-OWNERSHIP-ROW

| Field | Type | Required |
| :-- | :-- | :-- |
| artifact-record | string | true |
| semantic-purpose | string | true |
| authoritative-plane | string | true |
| authorized-writers | string | true |
| key-non-writers | string | true |

---

**End of Canonical Deconstruction Protocol v4.1.2**
