# Protocol: Legacy System Deconstruction, Assurance, and Reconstruction

**Version:** 4.1.1 (Canonical Reconstruction-Ready Edition)
**Status:** Normative  
**Purpose:** Produce an evidence-grounded, complete, framework-agnostic description of a legacy system and a separately reviewed reconstruction package without requiring downstream readers to reopen the legacy source.

---

## 0. Purpose, Guarantees, and Conceptual Flow

### 0.1. Normative language

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHOULD**, **SHOULD NOT**, and **MAY** are normative. MUST and MUST NOT are absolute protocol requirements. SHOULD and SHOULD NOT require a recorded reason when not followed. MAY is optional.

This protocol is implementation-independent. A conforming implementation MAY use any programming language, runtime architecture, storage engine, agent framework, model provider, user interface, or execution tooling, including ordinary local scripts created by a capable coding harness, provided it preserves every normative behavior, deterministic algorithm, authority boundary, artifact, state transition, and conformance condition defined here. No reference implementation, product, service, CLI, UI, or agent harness is normative or required. The protocol-defined `.extracted/` layout is interoperable protocol state, not an implementation source-tree requirement. In this document, **conforming runtime** means the implementation-neutral deterministic authority role defined below; it need not be a long-lived service or separate installed tool.

A **source unit** is one concrete, independently inventoryable unit in or affecting the system: a file, symbol, route, query, job, table, view, database routine, constraint, configuration key, workflow, node, edge, page, widget, action, binding, or equivalent. A **component** is one independently evidenced behavioral or structural record extracted from one source unit or from a precisely bounded portion of one source unit. A **material claim** is a fact whose falsity would alter behavior, scope, safety, compatibility, architecture, data semantics, or a reconstruction decision.

Protocol work has three authority classes:

1. **Executor-owned semantic work:** source traversal, behavior and business-rule interpretation, semantic classification, ambiguity recognition, claim formulation, and reconstruction prose. An executor MAY be a coding harness, model, remote worker, or human, but every semantic classification that affects scope, closure, mutation, or a gate MUST record its actor, invocation or human authority, reason, evidence/claim bindings, and pinned snapshot.
2. **Conforming-runtime deterministic work:** identity, path normalization, structural parsing, canonicalization, ordering, fingerprints, scope and stale-input checks, finite-state transitions, coverage arithmetic, completeness, gates, package integrity, and other protocol-defined mechanical results. A conforming runtime MUST compute these values or independently reproduce them from authoritative inputs before commit. Executor-supplied deterministic values are proposals only and MUST NOT be accepted as authoritative unless deterministic validation reproduces them exactly.
3. **Human-authorized decisions:** explicit exclusions, confirmations, policy or business intent, irreducible domain meaning, and acceptance of external facts only where this protocol permits it. A model inference or runtime computation MUST NOT substitute for the required human identity, authority basis, review scope, and exact input bindings.

Every normative predicate that affects conformance, mutation authorization, coverage, candidacy, or a gate MUST be one of: a **deterministic predicate** computed by the conforming runtime; a **semantic predicate** classified by an executor or authorized human with the provenance above; or an **unresolved semantic predicate** recorded as `Unknown` with its affected scope and required ticket/question. Unless a section defines a stricter result, `Unknown` is fail-inclusive: the item remains applicable/in scope for closure, cannot justify omission or `Not-Applicable`, and blocks any success that depends on the predicate being false. Natural-language guidance that does not affect those outcomes need not create a protocol record.

### 0.2. Guarantees

When followed completely, this protocol guarantees that:

1. every in-scope concrete source unit has a stable identity and a terminal disposition;
2. traversal closure is proven from an explicit source denominator and frontier, not inferred from ticket silence;
3. every material fact is represented by one claim with an evidence level, coordinate, and snapshot fingerprint;
4. active, conditional, shadow, disabled, retired, and unknown capabilities remain distinct across environments and snapshots;
5. serialized exports are acquired, exploded, reconciled, and traced at logical-unit level without exposing raw secrets;
6. personas and POVs retain strict ownership and isolation during discovery;
7. cross-layer ambiguity is resolved through a finite-state ticket lifecycle, runtime probes, or an explicit human hatch—never by guessing;
8. synthesis artifacts are structured semantic projections of closed forensic and assurance records, with derivation provenance and targeted stale propagation;
9. a separate human reconstruction handbook is navigable without using dense identifier catalogs as its primary interface;
10. reconstruction may treat the delivered bundle as its sole input only after the signed outer bundle manifest authoritatively attests `EXIT-E-STATUS: Passed` under the strict non-cyclic Exit E sequence.

### 0.3. Non-guarantees

The protocol does not:

- determine legal, privacy, security, or regulatory compliance;
- prove that unavailable production state equals repository state;
- make low-evidence facts true by repeating them in synthesis;
- convert correlation into causation;
- decide modernization architecture without accountable human approval;
- authorize execution against production, enable dormant production features, or move production data;
- make generated machine sidecars or helper output authoritative semantic truth.

All outputs remain a draft until the applicable human confirmation gates pass. Human architects and domain owners remain responsible for modernization decisions and final validation.

### 0.4. Conceptual flow

```text
Scope and snapshots
  -> source inventory and persona/cluster ownership
  -> acquisition and logical explosion where required
  -> isolated persona + cluster + POV traversal
  -> atomic forensic records + claim evidence + frontier updates
  -> tickets, probes, reconciliation, staging, and promotion
  -> Composite Exit A: Deconstruction Closure
  -> final structured synthesis
  -> human reconstruction handbook and decisions
  -> Exit E candidate validation and human confirmation
  -> Exit E Content-Readiness Report
  -> scope certificate
  -> signed outer bundle manifest with EXIT-E-STATUS: Passed
  -> signed, fingerprinted reconstruction bundle
```

Partial synthesis MAY run before Exit A for bounded review. It is provisional and cannot authorize reconstruction.

---

# Part 1. Four-Plane Information Model

## 1.1. Plane 1 — Forensic Plane

**Purpose:** Preserve what was directly found, where it was found, how it is invoked, and what remains uncertain.

**Authoritative records:**

```text
.extracted/
  0A-PREFLIGHT.md
  0B-AUTH-MODEL.md
  0C-SECURITY-PRIVACY.md
  0D-GLOSSARY.md
  0E-INDEX.md
  0F-GLOBAL-STATE.md
  0G-DECONSTRUCTION-STATE.md
  0H-CHECKPOINT-SUMMARY.md
  personas/
    _shared/
      01-FRONTEND-POV.md
      01-FRONTEND-QUESTIONS.md
      ... 05-WIRING-POV.md
      05-WIRING-QUESTIONS.md
    {persona-prefix}-{persona-slug}/
      PERSONA-PROFILE.md
      01-FRONTEND-POV.md
      ... 05-WIRING-QUESTIONS.md
      SHARED-STAGING-BUFFER.md
      SHARED-QUESTIONS-BUFFER.md
      UNMAPPED-DISCOVERY-BUFFER.md
  normalized/
    *.export-map.md
  probes/
    [TICKET-ID].log
```

**Sources of truth:** original source snapshots and approved sanitized evidence projections are evidence sources; the structured forensic records are the canonical extracted representation. `0G` is append-only process history. `0H` is derived and never replaces `0G`. Normalized maps are navigation indexes, not behavioral truth.

**Writers:** isolated POV invocations have exactly one semantic POV target and at most one semantic ledger target; their fixed transaction bundle MAY also append directly produced source, frontier, claim, contradiction, coverage, persona-local unmapped-discovery-buffer, and audit-log records explicitly allowed by the identity header. They cannot mutate unrelated semantic files. Acquisition writes the acquisition register, normalized maps, reconciliation records, approved registry fragments, and its directly allowed discovery placeholders. Profile Synchronization exclusively writes persona semantic payloads and synchronized `0A` persona rows. Authorized `Human Hatch: Confirmation` or `Reconstruction-Handoff` confirmation actions MAY update only the separately delimited profile certification envelope; they MUST NOT change profile semantic content, `SEMANTIC-RECORD-VERSION`, or `SEMANTIC-CONTENT-FINGERPRINT`. Sequential Reconciliation owns cross-persona/shared merges, atomic promotion and depromotion, buffered tickets, and merging persona-local unmapped discoveries into `0A`.

## 1.2. Plane 2 — Assurance Plane

**Purpose:** Prove denominator completeness, traversal closure, evidence support, export reconciliation, reference integrity, and gate outcomes.

**Authoritative records:**

```text
.extracted/
  10-SOURCE-INVENTORY.md
  11-TRAVERSAL-FRONTIER.md
  12-CLAIM-EVIDENCE.md
  13-EXPORT-RECONCILIATION.md
  14-CONTRADICTIONS.md
  15-DECISIONS.md
  16-SOURCE-COVERAGE.md
  17-ACQUISITION-CANDIDATES.md
  18-CONFIRMATIONS.md
  20-TRACEABILITY.md
  21-COVERAGE-REPORT.md
  22-GATE-REPORTS.md
```

**Sources of truth:** assurance records are authoritative for scope, denominator membership, evidence claims, coverage arithmetic, reconciliation, contradiction status, and gate results. They cite forensic evidence but do not replace it.

**Writers:** Acquisition, isolated traversal invocations, Sequential Reconciliation, Cross-Reference-Reconciliation, Human Hatch, and deterministic validators write only their explicitly owned assurance sections. Gate reports MUST be generated from the authoritative records and MUST include deterministic input fingerprints.

## 1.3. Plane 3 — Synthesis Plane

**Purpose:** Transform closed forensic facts and assurance proofs into coherent reconstruction semantics.

**Authoritative records:**

```text
.extracted/
  90-ARCH-BLUEPRINT.md
  91-DATA-MODEL.md
  92-BUSINESS-RULES.md
  93-USE-CASES.md
  94-INTERFACES.md
  95-DEPLOYMENT.md
  96-NON-FUNCTIONAL-SECURITY.md
```

**Sources of truth:** these structured catalogs are authoritative for synthesized architecture, entities, rules, use cases, interfaces, deployment/configuration, and NFR/security semantics. Every block MUST cite `DERIVED-FROM` IDs and fingerprints. Synthesis MUST NOT silently create a legacy fact absent from the Forensic or Assurance Plane.

**Writers:** Partial-Synthesis and Final-Synthesis create or update semantic payloads and their semantic versions/fingerprints. Cross-Reference-Reconciliation MAY update only reference fields after IDs exist; any reference field included in the semantic payload is a semantic change and MUST increment the semantic record version and fingerprint. Authorized confirmation attaches only the certification envelope and transitions eligible blocks from `[B-POPULATED]` to `[B-CONFIRMED]` without regenerating or changing semantic content. Synthesis writes discovered gaps only to a synthesis gap buffer; it never opens tickets directly.

## 1.4. Plane 4 — Human Reconstruction Handbook

**Purpose:** Present the confirmed bundle through progressive disclosure for architects, developers, domain experts, security reviewers, data engineers, QA, and operators.

**Canonical set:**

```text
.extracted/handbook/
  START-HERE.md
  01-SYSTEM-OVERVIEW.md
  02-ARCHITECTURE.md
  03-DOMAIN-AND-DATA.md
  04-PERSONAS-AUTH-AND-VISIBILITY.md
  05-USE-CASES.md
  06-BUSINESS-RULES.md
  07-INTERFACES-AND-INTEGRATIONS.md
  08-BACKGROUND-PROCESSING.md
  09-DEPLOYMENT-AND-CONFIGURATION.md
  10-SECURITY-PRIVACY-AND-AUDIT.md
  11-NON-FUNCTIONAL-PROFILE.md
  12-DISABLED-AND-DORMANT-CAPABILITIES.md
  13-MODERNIZATION-DECISIONS.md
  14-EQUIVALENCE-AND-ACCEPTANCE.md
  15-KNOWN-GAPS-AND-RISKS.md
  16-GLOSSARY-AND-REFERENCE.md
```

**Source of truth:** the handbook is a synchronized, human-readable projection of B-populated or B-confirmed Synthesis Plane records, assurance scope, and approved decisions. It is not an independent semantic authority. Its complete semantic payload MUST be generated as `[B-POPULATED]`, assigned legal `HBK` identities, and fingerprinted before the Exit E Candidate Report and human review. Semantic corrections MUST be made in the upstream authoritative record first and then re-projected before candidacy. Confirmation only attaches the handbook certification envelope; it MUST NOT regenerate or alter semantic content. Broken derivation marks the affected handbook section `[B-STALE]`.

**Writers:** Reconstruction-Handoff generates and updates handbook semantic payloads before candidacy. Human editors MAY improve wording, examples, navigation, and diagrams before candidate fingerprinting only when semantic meaning remains unchanged; after candidacy, any content edit is a semantic payload change that invalidates the candidate and restarts from populated generation. Authorized confirmation actions MAY update only certification envelopes. Semantic edits require upstream changes and reprojection.

## 1.5. Non-authoritative sidecars

JSON, SQLite, graph, search, code-generation, or validator sidecars MAY be generated for automation. They MUST carry source fingerprints, schema version, generation time, and a prominent `NON-AUTHORITATIVE-DERIVATIVE` marker. Conflicts are resolved in favor of the canonical Markdown records and original evidence, in that order for semantics.

---

# Part 2. Personas, POVs, Entry Ownership, and Traversal Tracks

## 2.1. Persona definition and purity

A **Persona** is a distinct human actor, service actor, external system, or autonomous runtime class entering through a concrete execution surface. Personas are selected by distinct physical or virtual entry paths, not merely business titles. Roles sharing identical routes, triggers, and paths SHOULD be grouped; materially different data visibility, authorization, or execution paths require distinct personas.

**Persona Purity Invariant:** a `personas/{persona-directory}/` directory MUST document only that persona's triggers, invocation paths, local UI/runtime state, and usage references. Shared components are linked, not copied as another persona's behavior. Purity does not imply exclusive component ownership.

`personas/` is the sole container for persona-scoped and canonical shared forensic files. `personas/_shared/` is a reserved non-persona directory containing canonical shared POV and question records; the leading underscore provides deterministic visual/sort separation and does not create a persona. Persona discovery, counting, purity, and profile validation MUST ignore that reserved directory as a persona while still validating its shared-record ownership rules.

Every active persona MUST have one globally unique uppercase prefix matching `[A-Z][A-Z0-9]{1,7}` and one explicit lowercase persona slug matching `[a-z0-9]+(?:-[a-z0-9]+)*`. Prefix uniqueness is validated in `0A-PREFLIGHT.md`; reuse is forbidden, including retired personas. The canonical persona-directory basename is exactly `<persona-prefix>-<persona-slug>` and therefore matches `[A-Z][A-Z0-9]{1,7}-[a-z0-9]+(?:-[a-z0-9]+)*`. The `0A` registry stores the persona slug explicitly; it MUST NOT be inferred from, transliterated from, or normalized from the display name. The basename prefix, registry prefix, invocation `Persona Prefix`, ticket prefix, and PRF owner coordinate MUST agree exactly. A non-reserved directory directly beneath `personas/` that is unregistered, malformed, or prefix-mismatched is invalid; `_shared` MUST NOT appear as a persona prefix, slug, basename, or registry row.

## 2.2. Canonical entry ownership

Every concrete entry coordinate MUST have exactly one canonical owner persona. Intentional shared entry surfaces MUST use a `SHARED-ENTRY` record that:

- owns the coordinate once;
- lists participating personas and each local invocation condition;
- identifies one coordination owner for registry maintenance;
- does not weaken persona purity;
- maps each persona to its own invocation reference.

Duplicate ownership without a shared-entry record is a gate failure.

### DB-autonomous ownership

Database-side autonomous behavior uses the reserved persona `db-autonomous` with prefix `DBR` when it fires because of database mutation, timer, internal rule, or database event rather than an app entry. An app routine that causes a trigger to fire remains owned by its app persona; the trigger body is owned by DBR and cross-linked through a `DR` record. Procedures explicitly invoked only through one app entry remain app-owned but are still extracted by POV-4. Mixed invocation requires a shared-entry record. Ownership is determined from firing mechanism, never convenience.

## 2.3. Five isolated POVs

1. **POV-1 Frontend Workflows:** client UI, server-rendered views, widget behavior, presentation state, validation before transport, templates, and interaction flow.
2. **POV-2 Backend Workflows:** domain decisions, business logic, service pipelines, validation engines, and transactional intent independent of transport and physical schema.
3. **POV-3 Background Workflows:** schedules, queues, event loops, delayed execution, retries, async coordination, and long-running workers.
4. **POV-4 Data Movement Workflows:** reads, writes, queries, transactions, ORM behavior, caches, files, tables, views, constraints, triggers, stored routines, generated fields, and persistence lifecycle.
5. **POV-5 System Wiring Workflows:** routes, middleware, transport contracts, graph edges, callbacks, integrations, orchestration, and inter-system coupling.

A POV MUST NOT fill another POV's semantic fields directly. Cross-POV findings become a ticket or persona-local shared buffer entry.

## 2.4. Traversal tracks

Every entry cluster and invocation MUST declare exactly one **TRAVERSAL-TRACK**:

- `Normal`: currently active, normally exposed behavior for the specified environment/snapshot;
- `Shadow/Conditional`: flag-gated, environment-dependent, alternate, shadow, indirectly exposed, or not normally selected behavior;
- `Disabled`: explicitly dormant, disabled, disconnected, legacy-retained, or configured off behavior.

Tracks require separate invocations even when they share source units. They MAY cross-link the same `SRC` IDs. `Disabled` is not equivalent to dead code, retired code, irrelevant code, or out-of-scope code.

## 2.5. Capability and runtime state

Capability state is semantic metadata, not a prefix-group tag:

```markdown
- **CAPABILITY-STATE:** Active | Conditional | Shadow | Disabled | Retired | Unknown
- **RUNTIME-STATE-MATRIX:**
  | Environment | Snapshot | State | Enable/Disable Mechanism | CLM Evidence | Last Known At |
- **Historical / Intended Personas:** [...]
- **Dependencies:** [SRC/CMP/CFG/IF IDs]
- **Last-Known State:** (value, environment, date, or Unknown)
```

State MAY differ by environment or snapshot. Absence of production evidence MUST NOT be converted into deployed behavior. In-scope disabled units require the same atomic static coverage as active units and block Exit A when undisposed.

A disabled capability MUST NOT be enabled in production for a probe. Runtime observation is permitted only in an isolated sandbox with explicit human authorization, masked/non-production data, and a probe specification that states the disabled behavior and safety boundaries.

## 2.6. Capability grouping and target decision

A capability may group active and dormant source units for navigation, but the `CAP` is non-promotable and inherits no evidence, status, or coverage from members. Every effective-denominator Disabled `SRC`, `CMP`, and `UC` MUST map to exactly one canonical `CAP`; the only alternative is an exact approved exclusion that removes that unit from the effective denominator. Shadow/Conditional units MAY map to a CAP when they represent a governed capability. CAP membership is mandatory in disabled traceability, synthesis, handbook output, and target-decision validation.

Multiple or conflicting CAP memberships are invalid unless an explicit typed relationship and approved `DEC` define why the memberships overlap, identify one canonical target-decision owner, and prevent duplicate denominator or decision effects. Such a relationship does not make both memberships canonical and MUST NOT silently duplicate a unit. Every disabled/dormant capability MUST preserve:

- exact member IDs, including all effective-denominator Disabled `SRC`/`CMP`/`UC` units;
- status by environment/snapshot;
- enable/disable mechanism and evidence;
- dependencies and affected personas;
- historical/intended behavior;
- safety and reconstruction risk;
- one explicit target decision: `Restore`, `Preserve Dormant`, `Redesign`, or `Retire`;
- decision owner, rationale, and approval state.

Observed legacy facts and target decisions MUST appear in separate fields and handbook sections.

---

# Part 3. Workspace Registries and Schemas

## 3.1. `0A-PREFLIGHT.md`

```markdown
# 0A-PREFLIGHT: System Inventory and Discovery Registry

## 1. Metadata and Boundaries

- **System Namespace:** [stable slug]
- **Target Repository / Evidence Roots:** [...]
- **Primary Technology Stack:** [...]
- **Protocol Version:** 4.1
- **Current Iteration:** [one canonical token from §10.6: ALFA ... ZULU]
- **Default Max Traversal Depth:** 3
- **Included Environments / Snapshots:** [...]
- **Explicit Scope Policy:** [...]
- **Protocol-Migration Metadata:** None | [v3.17 migration state]

## 2. Persona Registry

| Persona | Prefix | Persona Slug | Persona Directory | Actor Class | Canonical Entry IDs | Scope Paths | Registry State |

## 3. Entry-Point Clusters

### [CLUSTER-ID] [Title]

- **Concrete Entry Coordinates:** [no unexpanded wildcard at closure]
- **Canonical Owner Persona:** [name/prefix or SHARED-ENTRY-ID]
- **Primary POV:** POV-1 | POV-2 | POV-3 | POV-4 | POV-5
- **Traversal Track:** Normal | Shadow/Conditional | Disabled
- **Capability State:** Active | Conditional | Shadow | Disabled | Retired | Unknown
- **Environment/Snapshot:** [...]
- **Source IDs:** [SRC IDs]
- **Max Depth:** [integer]
- **Traversal Status:** [R-ACTIVE] | [R-UNSWEPT] | [R-IN-PROGRESS] | [R-SWEPT] | [R-PENDING-HUMAN-REVIEW] | [R-PENDING-PERSONA-ASSIGNMENT] | [R-DEPROMOTION-PENDING]
- **Required Traversal Matrix:**
  | Persona/Owner | Concrete Entry | Track | Environment | Snapshot | POV | Applicability | Status | SWP ID | Claims |
  | ... | ... | ... | ... | ... | POV-1 | Required or Not-Applicable | [R-UNSWEPT]/[R-IN-PROGRESS]/[R-SWEPT] or N/A | ... | ... |

Every persona × concrete entry × track × environment × snapshot × POV cell MUST exist. `Not-Applicable` requires a claim-backed reason. `Primary POV` selects the first routing pass but never reduces this closure denominator.

### Database-Side Autonomous Clusters

[Same schema; DBR ownership rule applies.]

## 4. Unmapped Discovery Queue

### [DISCOVERY-ID]

- **Coordinate / SRC:** [...]
- **Discovered By Invocation:** [...]
- **Suggested POV / Persona:** [...]
- **Traversal Status:** [R-PENDING-PERSONA-ASSIGNMENT] | [R-PENDING-HUMAN-REVIEW]

## 5. Serialized Export Acquisition Register

[Section 7.2 schema]
```

Agents MAY append acquisition reference placeholders directly where Acquisition permits. Bound Discovery invocations MUST NOT write the global Unmapped Discovery Queue; they append the persona-local buffer defined below. Only Preflight/Profile Synchronization may assign personas or alter boundaries. Sequential Reconciliation is the sole merger of persona-local unmapped discoveries into `0A`.

### Persona-local unmapped discovery buffer

Each `personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md` entry has this schema:

```markdown
### [DSC-ID]

- **Concrete Coordinate / Proposed SRC / Fingerprint:** [...]
- **Discovering Invocation / Persona / Cluster / Track / POV:** [...]
- **Discovery Method / Boundary FRT:** [...]
- **Suggested Persona / Cluster / Track / POV:** [...]
- **Reason Unmapped / Materiality:** [...]
- **DEDUP-KEY:** [sha256(system namespace | normalized coordinate | snapshot | environment | track)]
- **State:** Pending-Merge | Merged | Superseded
```

A bound Discovery invocation MUST append this record in the same transaction that records the discovery and its `FRT`; delayed or conversational-only discovery is forbidden. Sequential Reconciliation sorts pending entries by `DEDUP-KEY`, then `DSC-ID`, merges equal keys into one canonical `0A` queue item while preserving all origin invocations and evidence, records contradictions for incompatible proposals, marks source entries consumed, and completes before any later Discovery invocation begins. A later Discovery invocation MUST reject a checkpoint showing unmerged pending entries. Acquisition MAY continue to create its allowed direct placeholders and global discoveries because Acquisition owns the applicable register sections.

## 3.2. Foundational registries

- **`0B-AUTH-MODEL.md`:** authentication mechanisms, token/session lifecycle, role/capability bindings, route/middleware/inline/DB enforcement, termination, and audit references. All POVs may propose buffered facts; Sequential Reconciliation merges them.
- **`0C-SECURITY-PRIVACY.md`:** sensitive-data evidence and review points for human assessment. It records risks and safeguards without declaring compliance.
- **`0D-GLOSSARY.md`:** canonical domain terms, aliases, conflicts, and evidence.
- **`0E-INDEX.md`:** directory-to-POV and source/component/persona cross-reference index.
- **`0F-GLOBAL-STATE.md`:** architectural invariants, shared runtime state, mutators/readers, lifecycle, and mutation history.
- **`0G-DECONSTRUCTION-STATE.md`:** chronological append-only invocation log.
- **`0H-CHECKPOINT-SUMMARY.md`:** derived resume accelerator.

### Shared-entry, authentication, privacy, and capability records

```markdown
### [SHR-ID] [Shared Entry] <!-- stored in 0A -->

- **Concrete Entry / SRC:** [...]
- **Participating Personas and Invocation Conditions:** [...]
- **Coordination Owner:** [...]
- **Persona Invocation Pointers:** [...]
- **Environment / Snapshot / Track:** [...]
- **Material Claims / Version:** [...]

### [AUTH-ID] [Authentication or Session Mechanism] <!-- stored in 0B -->

- **Mechanism / Credential Form:** [...]
- **Issuer / Validator / Storage:** [...]
- **Creation, Refresh, Expiry, Revocation, Termination:** [...]
- **Authorization Bindings / Enforcement Points:** [...]
- **Personas / Interfaces / Environments:** [...]
- **Secrets Handling / Audit:** [references only, never values]
- **Material Claims / Version:** [...]

### [PRV-ID] [Security or Privacy Finding] <!-- stored in 0C -->

- **Data / Threat / Review Category:** [...]
- **Affected SRC/ER/IF/Persona/Environment:** [...]
- **Observed Handling / Risk / Unknowns:** [...]
- **Human Review Required:** [...]
- **Material Claims / Version:** [...]
- **Synthesis SEC Mirror:** [SEC-ID or Pending-Synthesis-Reference]

### [CAP-ID] [Capability] <!-- stored in 0F -->

- **RECORD-KIND:** GOVERNED-CAPABILITY-CONTAINER
- **Canonical Member SRC/CMP/UC/IF/CFG/DR IDs:** [...]
- **Overlapping CAP Relationships / Decisions:** [typed relation + DEC IDs or None]
- **Runtime-State Matrix / Mechanism / Last Known State:** [...]
- **Historical and Intended Personas / Behavior:** [...]
- **Dependencies / Safety / Risk:** [...]
- **Observed Claims:** [CLM IDs]
- **Target Decision State:** Bound | Pending
- **Target Decision DEC:** [DEC-ID or None]
- **Version / Supersession:** [...]
```

A `CAP` is a non-promotable governance container and inherits no child staging, evidence, or coverage. Its own material governance metadata remains claim-backed. Every effective-denominator Disabled `SRC`, `CMP`, and `UC` has exactly one canonical CAP membership unless exactly approved excluded; overlap requires the explicit relationship and decision defined in §2.6.

### `0E` routine cross-reference

```markdown
### [CMP-ID] [Atomic Record Name]

- **SRC:** [SRC IDs]
- **Source Coordinate:** [...]
- **Canonical Extract:** [file -> heading]
- **Personas:** [...]
- **Cross-Layer Consumers:** [CMP IDs or None]
```

### `0F` invariant/state records

```markdown
### [INV-ID] [Invariant]

- **Constraint:** [...]
- **Verification Vectors:** [CLM IDs]
- **Discovered In:** [Invocation ID / POV]

### [STATE-ID] [Shared State]

- **Medium / Shape:** [...]
- **Initializer / Mutators / Readers:** [CMP IDs]
- **Lifecycle / Invalidation:** [...]
- **Claims:** [CLM IDs]
```

---

# Part 4. Deterministic Identity and Universal Atomicity

## 4.1. ID generation

The system namespace is fixed at preflight. Every canonical typed record that participates in identity, reference, traceability, candidacy, or confirmation MUST use a legal protocol ID. Ticket IDs use §9.1; all other canonical typed record IDs use:

```text
canonical-key = protocol-id-type + "|" + system-namespace + "|" + normalized-owner-coordinate + "|" + semantic-discriminator
hash = uppercase(first 12 hexadecimal characters of SHA-256(canonical-key UTF-8))
id = TYPE + "-" + hash
```

Packaging artifacts are not typed records and do not receive protocol record IDs. Their normative `HASH-DOMAIN` identity is exactly `artifact-type + "|" + normalized-relative-path`; that identity is included in the artifact payload hash preimage. Only after the payload fingerprint is computed, the artifact instance receives the separate `POST-HASH-ARTIFACT-INSTANCE-BINDING` `artifact-type + "|" + normalized-relative-path + "|" + payload-fingerprint`. The post-hash binding is not part of its own fingerprint preimage and MUST NOT be substituted for the hash-domain identity.

For `SRC`, include kind: `SRC-KIND-HASH`, where `KIND` is one of `FILE`, `SYM`, `ROUTE`, `RPC`, `JOB`, `QUERY`, `TABLE`, `VIEW`, `DR`, `CONSTRAINT`, `GENERATED`, `CFG`, `WORKFLOW`, `NODE`, `EDGE`, `PAGE`, `WIDGET`, `ACTION`, `BINDING`, `QUEUE`, `EVENT`, `STORAGE`, `IFACE`, or `OTHER`. Classification uses the most specific executable/data kind; a unit with several roles has one primary SRC kind and explicit typed relationship records to the others, preventing denominator double-counting. Stable platform IDs take precedence in normalized coordinates. Unicode is preserved exactly; path separators and harmless whitespace are normalized, but identifiers are never transliterated. A proven move or rename retains its original identity canonical key and records the new coordinate as versioned location metadata and an alias; the path is not rehashed.

Required forms:

```regex
SRC-(FILE|SYM|ROUTE|RPC|JOB|QUERY|TABLE|VIEW|DR|CONSTRAINT|GENERATED|CFG|WORKFLOW|NODE|EDGE|PAGE|WIDGET|ACTION|BINDING|QUEUE|EVENT|STORAGE|IFACE|OTHER)-[A-F0-9]{12,64}
(CMP|CLM|ER|REL|BR|UC|IF|DEP|CFG|SCHED|NFR|SEC|SM|DR|CAP|PRV|FLT|MOD|PRF|HBK)-[A-F0-9]{12,64}
(INV|AUTH|STATE|DEC|CON|FRT|SWP|COV|CND|CLU|SHR|DSC|EXP|GAP|CNF)-[A-F0-9]{12,64}
(ALFA|BRAVO|CHARLIE|DELTA|ECHO|FOXTROT|GOLF|HOTEL|INDIA|JULIETT|KILO|LIMA|MIKE|NOVEMBER|OSCAR|PAPA|QUEBEC|ROMEO|SIERRA|TANGO|UNIFORM|VICTOR|WHISKEY|X-RAY|YANKEE|ZULU)-[A-Z][A-Z0-9]{1,7}-(01|02|03|04|05)-[0-9]{3,6}
```

The auxiliary types are `INV` (invariant), `AUTH`, `STATE`, `DEC` (decision), `CON` (contradiction), `FRT` (frontier), `SWP` (sweep), `COV` (scope-specific source coverage row), `CND` (acquisition candidate or operator action), `CLU` (cluster), `SHR` (shared entry), `DSC` (discovery), `EXP` (export), `GAP` (synthesis gap), and `CNF` (human confirmation). `MOD` identifies an atomic semantic architecture module or bounded context. `PRF` identifies one persona profile; its normalized owner coordinate is the globally unique persona prefix, and its semantic discriminator is the fixed literal `persona-profile`, so its canonical key is determined only by protocol type, system namespace, and persona prefix. `HBK` identifies one semantic handbook section. Its normalized owner coordinate is the normalized handbook-relative path plus `#` plus its explicitly assigned stable section anchor, and its semantic discriminator is the fixed literal `handbook-section`, so its canonical key is determined only by protocol type, system namespace, normalized path, and stable anchor. The heading text is not an identity source. A heading/title edit does not alter the HBK ID when the assigned anchor and semantic identity remain unchanged. A proven move retains its HBK identity under the rename rule and records the prior coordinate as an alias; a semantic split creates new HBK IDs and tombstones or supersedes the old record. `CAP` identifies a governed capability container, `PRV` an atomic security/privacy finding in `0C`, and `FLT` a fault/resilience record. Schema placeholders such as `[FRONTIER-ID]`, `[COVERAGE-ROW-ID]`, `[CANDIDATE/ACTION-ID]`, `[MODULE-ID]`, `[DECISION-ID]`, or `[EXPORT-ID]` resolve to these forms. `COV` owner coordinates are the immutable SRC-version/persona/track/environment/snapshot tuple; `CND` owner coordinates are the helper-published candidate or operator-action stable identity. A semantic catalog without a physical owner coordinate uses a stable domain path plus normalized semantic title as its owner coordinate; its discriminator includes record type and authoritative domain owner.

A **normalized relative path** is relative to the applicable workspace root or package root, serialized with POSIX `/`, has no leading `./`, removes repeated separators and `.` path segments, and rejects every `..` segment. Path case and Unicode code points are preserved exactly; implicit Unicode normalization, transliteration, and case folding are forbidden. The root used for normalization MUST be declared by the enclosing schema or package operation and remain fixed for the artifact or record set.

An HBK **stable section anchor** is an explicitly assigned opaque ASCII lowercase token matching `[a-z0-9]+(?:-[a-z0-9]+)*`. It is never renderer-generated or inferred from heading text. For every HBK, the sole normalization root is the literal `.extracted/handbook/` directory declared in §1.4; `HANDBOOK-RELATIVE-PATH` omits that root (for example `START-HERE.md`), and package-relative encodings such as `handbook/START-HERE.md` or `.extracted/handbook/START-HERE.md` are invalid identity coordinates. The exact HBK coordinate is `normalized-handbook-relative-path + "#" + stable-section-anchor`. Punctuation, uppercase characters, non-ASCII characters, empty segments, leading/trailing hyphens, and repeated hyphens are invalid anchors.

A collision MUST be resolved by extending both colliding hashes to 16, then 20, up to 64 characters; existing shorter IDs receive aliases and controlled supersession, never silent mutation.

The ID registry in `20-TRACEABILITY.md` MUST verify global uniqueness, type correctness, canonical-key uniqueness, ticket-tuple uniqueness, coverage-tuple uniqueness, candidate/action identity uniqueness, PRF-per-persona-prefix uniqueness, HBK path-and-anchor uniqueness, and tombstone reservation. It MUST also verify reciprocal `PRF <-> 0A persona row`, `HBK <-> handbook path/anchor and upstream synthesis records`, and exact candidate/CNF bindings to PRF/HBK IDs. IDs are stable across movement when semantic identity and stable platform identity remain unchanged. A semantic split creates new IDs and tombstones/supersedes the old record.

## 4.1.1. Semantic payload identity and certification envelopes

Every synthesis block, persona profile, and semantic handbook section has a canonical semantic payload and a separately delimited certification envelope. A persona profile payload begins with its mandatory `PRF-ID`; a handbook-section payload begins with its mandatory `HBK-ID` and includes separately ordered `HANDBOOK-RELATIVE-PATH` and `STABLE-SECTION-ANCHOR` fields before its semantic record version and fingerprint carrier. These typed IDs and HBK coordinate fields are part of the semantic payload and candidate/CNF binding; headings are not identity parser inputs.

- **`SEMANTIC-RECORD-VERSION`:** a positive integer changed only when semantic payload meaning or payload bytes under the canonicalization rules change.
- **`SEMANTIC-CONTENT-FINGERPRINT`:** `sha256` of the record's canonical semantic payload under §4.1.2. The hash input includes record ID, semantic record version, schema/protocol version, and all semantic fields, prose, tables, diagrams, examples, and resolved semantic references with their bound versions/fingerprints. The carrier field that stores `SEMANTIC-CONTENT-FINGERPRINT` is excluded from its own hash input.
- **Certification envelope:** status and approval metadata attached after payload generation. It contains `BLUEPRINT-STATUS` or `HANDBOOK-STATUS`, `CNF` references, confirmation-envelope version, approval/signature identities and timestamps, and envelope audit fields.

The semantic-content hash and semantic record version MUST exclude the complete certification envelope, including B/handbook status, CNF reference, confirmation-envelope version, approval/signature identities and timestamps, candidate report or candidate payload manifest references explicitly permitted by the envelope schema, and envelope digests. A certification envelope MUST NOT reference the later Exit E Content-Readiness Report, scope certificate, outer bundle manifest, or any other artifact that did not exist when the envelope was attached. Attaching a confirmation or changing `[B-POPULATED]` to `[B-CONFIRMED]` therefore changes only the envelope and MUST NOT change the semantic record version or semantic-content fingerprint. Any edit to semantic prose, diagrams, examples, claims, derivation references, semantic cross-references, PRF identity, or HBK identity increments `SEMANTIC-RECORD-VERSION`, recomputes `SEMANTIC-CONTENT-FINGERPRINT`, invalidates the old envelope, and triggers targeted stale propagation.

A `CNF` confirms exact typed record IDs—including `PRF` and `HBK`—`SEMANTIC-RECORD-VERSION` values, and `SEMANTIC-CONTENT-FINGERPRINT` values, never mutable whole-file bytes or status-bearing envelopes. Validators MUST parse the payload/envelope boundary from the Markdown AST, recompute payload hashes, and reject envelope fields inside the semantic payload or semantic fields inside the envelope.

## 4.1.2. Normative canonical hash profile and packaging artifacts

Unless a schema explicitly requests a raw-source byte fingerprint, all semantic records, `DEC`, `CNF`, Exit E Candidate Reports, candidate payload manifests, Exit E Content-Readiness Reports, scope certificates, and outer bundle manifests MUST use this canonical hash profile:

1. parse the bounded payload as a normative Markdown AST and reject malformed or overlapping boundaries;
2. normalize line endings to LF while preserving Unicode code points exactly without transliteration or implicit Unicode normalization;
3. serialize fields in the deterministic order declared by the applicable schema; serialize tables with declared column order, normalized delimiter/alignment syntax and insignificant cell-edge whitespace, preserving semantic row order unless the schema declares a deterministic row sort key;
4. remove insignificant trailing whitespace, AST-separator blank-line variance, and transport-only indentation while preserving whitespace in code spans, fenced blocks, diagrams, and other literal nodes; emit exactly one terminal LF;
5. include the protocol/schema version, record or artifact identity, and every in-boundary semantic/content node; exclude only the exact carrier and envelope regions authorized by the applicable schema.

For a typed record in a multi-record Markdown file, the hash boundary begins at that record's typed-ID heading and ends immediately before the next typed-record heading of the same or higher level, or at end of file. Each record is canonicalized and hashed independently after applying its exact envelope exclusion. The containing file has a separate transport fingerprint over its normalized relative path and ordered sequence of `(record type, record ID, record version where applicable, record payload fingerprint)` bindings. That file fingerprint detects ordering or transport changes but is not a semantic-record identity, is not substituted for a record fingerprint, and MUST NOT be bound by a `CNF`.

`DEC` uses this profile over its heading-delimited decision content, excludes exactly `DECISION-CONTENT-FINGERPRINT` as its own carrier and the existing `<!-- DECISION-APPROVAL-ENVELOPE: ... -->` through `<!-- END-DECISION-APPROVAL-ENVELOPE -->` region. `CNF` uses this profile over its heading-delimited confirmation payload, excludes exactly `CNF Payload Fingerprint` as its own carrier and the existing `<!-- CNF-DIGEST-SIGNATURE-ENVELOPE: ... -->` through `<!-- END-CNF-DIGEST-SIGNATURE-ENVELOPE -->` region. Semantic synthesis, PRF, and HBK records use §4.1.1 and exclude exactly their existing certification envelope.

Semantic record fingerprints and the ordered file-transport fingerprint intentionally exclude mutable certification/approval/signature envelopes and therefore MUST NOT be used as final package-member integrity fingerprints. Every populated certification envelope, DEC approval envelope, and CNF digest/signature envelope receives an externally stored **canonical envelope fingerprint**. Its preimage is the UTF-8 domain prefix `ENVELOPE|` plus normalized containing path, target typed record ID, envelope kind (`CERTIFICATION`, `DECISION-APPROVAL`, or `CNF-DIGEST-SIGNATURE`), positive envelope version, one LF, and the complete canonically serialized Markdown AST region between the exact envelope markers. No field in that region is excluded: status, target-content binding, approver/signatory identity, authority, timestamps, audit data, signatures, and any payload-fingerprint carrier inside the envelope are all included. Missing, duplicate, malformed, overlapping, or versionless final envelopes are invalid.

An **envelope binding** is the exact tuple `(target typed record ID, envelope kind, envelope version, canonical envelope fingerprint)`. `Certification Envelope Binding`, `Approval Envelope Binding`, and `CNF Signature Envelope Binding` fields MUST serialize that tuple and MUST be recomputed rather than trusted. Attaching or changing an envelope leaves the underlying semantic/decision/CNF payload fingerprint unchanged as intended, but changes its envelope fingerprint and binding.

Every final package member also receives an external **package-member file fingerprint**. For Markdown, its preimage is UTF-8 `PACKAGE-MEMBER|` plus normalized package-relative path, one LF, and the complete canonicalized file AST, including every record fingerprint carrier and every certification, approval, digest, and signature envelope with no excluded node. For non-Markdown members, the same path-bound domain prefixes the exact immutable bytes. This fingerprint is stored only in later manifests/reports, never inside the member whose bytes it covers, so it cannot self-reference. `Containing File Fingerprint`, `Package-Member File Fingerprint`, and outer-manifest member fingerprints mean this full finalized-file fingerprint, never the envelope-excluding transport fingerprint. Content-readiness and step-6 validation MUST recompute every envelope binding and package-member file fingerprint; removal or mutation of an envelope after review invalidates the package while leaving the separately scoped semantic fingerprint unchanged.

An **evidence-set fingerprint** is the value serialized in every `Evidence Fingerprint` check or blocker cell in §15.1. Its semantic purpose is to bind that row to the exact authoritative evidence and deterministic validation inputs used for its result; it is not itself a semantic-record, containing-file, validator-output, envelope, or package-member fingerprint. The corresponding §12.5 validation summary MUST record the row's complete evidence-binding set. Each binding is the tuple `(normalized authoritative input path, binding kind, typed record or artifact ID or None, version or None, authoritative fingerprint)` using exactly one closed kind below:

| Binding Kind        | Required Fingerprint Domain                                                                              | Version Field                                     |
| :------------------ | :------------------------------------------------------------------------------------------------------- | :------------------------------------------------ |
| RAW-SOURCE-BYTES    | the schema-requested immutable raw-source byte fingerprint                                               | immutable source version or `None`                |
| RECORD-PAYLOAD      | the §4.1.2 heading-delimited canonical payload fingerprint for a non-semantic authoritative typed record | record version or `None` when the schema has none |
| SEMANTIC-CONTENT    | `SEMANTIC-CONTENT-FINGERPRINT` under §4.1.1                                                              | semantic record version                           |
| DECISION-CONTENT    | `DECISION-CONTENT-FINGERPRINT` under §5.5                                                                | decision-content version                          |
| CNF-PAYLOAD         | `CNF Payload Fingerprint` under §5.7                                                                     | `None`                                            |
| ARTIFACT-PAYLOAD    | `ARTIFACT-PAYLOAD-FINGERPRINT` under the generic §4.1.2 artifact schema                                  | artifact-schema version                           |
| CANONICAL-ENVELOPE  | canonical envelope fingerprint under §4.1.2                                                              | positive envelope version                         |
| FILE-TRANSPORT      | the ordered containing-file transport fingerprint under §4.1.2                                           | `None`                                            |
| PACKAGE-MEMBER-FILE | the complete finalized path-bound package-member fingerprint under §4.1.2                                | `None`                                            |
| VALIDATION-SUMMARY  | the §12.5 deterministic validation-summary output fingerprint                                            | validator/check-registry version                  |
| DENOMINATOR-QUERY   | the frozen denominator query/result fingerprint required by the applicable source/sweep/coverage schema  | denominator schema version                        |
| PROTOCOL-REGISTRY   | the canonical fingerprint of the identified protocol-defined registry                                    | registry version                                  |

`Binding Kind` selects the only legal fingerprint domain for that row input; a different valid fingerprint over the same file or record is a mismatch, not an alternative. The conforming runtime MUST reject an unknown kind, a version inconsistent with the table, or duplicate tuples; sort bindings bytewise by the complete tuple after §4.1 normalization; and compute `sha256` over UTF-8 `EVIDENCE-SET|` plus the packaging artifact type, one LF, row kind `CHECK` or `BLOCKER`, one LF, stable row ID, one LF, the canonically serialized `System / Protocol / Snapshot` value, one LF, and the §4.1.2 canonical table serialization of the sorted tuples with columns in tuple order. The set MUST contain every input that can change the row result and at least one `DENOMINATOR-QUERY`, `PROTOCOL-REGISTRY`, or authoritative record/artifact binding even for an evidence-backed zero domain. It MUST NOT include the containing report, its own carrier, or a later artifact. A conforming runtime MUST recompute this fingerprint; an executor-supplied digest is non-authoritative and a mismatch is a structural failure.

Each packaging artifact is one bounded payload. Its `HASH-DOMAIN` identity is exactly `artifact-type + "|" + normalized-relative-path`, using the path normalization rules in §4.1, and that identity is included in the canonical payload hash preimage. The payload fingerprint MUST NOT appear anywhere in its own preimage. After the payload fingerprint is computed, the separate `POST-HASH-ARTIFACT-INSTANCE-BINDING` is `artifact-type + "|" + normalized-relative-path + "|" + payload-fingerprint`; this binding identifies the resulting artifact instance, is computed only post-hash, and is not part of its own preimage. Packaging artifacts MUST use these exact generic boundaries:

```markdown
<!-- ARTIFACT-PAYLOAD: [ARTIFACT-TYPE] -->

- **ARTIFACT-TYPE:** [EXIT-E-CANDIDATE-REPORT | CANDIDATE-PAYLOAD-MANIFEST | EXIT-E-CONTENT-READINESS-REPORT | SCOPE-CERTIFICATE | OUTER-BUNDLE-MANIFEST]
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
  [all artifact payload fields]
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256; this carrier field excluded from its own hash]
<!-- END-ARTIFACT-PAYLOAD: [ARTIFACT-TYPE] -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: [ARTIFACT-TYPE] -->

- **Signature / Signatories / Timestamps:** [as required by the artifact stage]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: [ARTIFACT-TYPE] -->
```

The payload fingerprint covers exactly the AST content between the artifact-payload markers under this profile, excluding only `ARTIFACT-PAYLOAD-FINGERPRINT`; it excludes the complete matching artifact digest/signature envelope. The canonical preimage is prefixed or otherwise structurally bound to the exact `HASH-DOMAIN` identity `artifact-type + "|" + normalized-relative-path`; neither the payload fingerprint nor the post-hash artifact-instance binding may occur in that preimage. The `POST-HASH-ARTIFACT-INSTANCE-BINDING` is computed only after the payload fingerprint. The Exit E Candidate Report, candidate payload manifest, Exit E Content-Readiness Report, scope certificate, and outer bundle manifest MUST use this schema with their literal artifact type. Their type, path, member lists, status fields, references, and content timestamps remain hashed payload. No artifact may specialize or substitute boundary marker names; this generic schema supersedes certificate- or report-specific digest markers.

## 4.2. Atomic unit rule

Each promotable forensic record and each independently meaningful synthesis record MUST describe exactly one independently evidenced unit: one function, method, route, query, job, workflow node, graph edge, action, binding, DB routine, constraint, widget behavior, interface operation, architecture module/bounded context, or equivalent.

A **purely navigational grouping container**—for example a family index, class overview, workflow overview, capability index, module index, or use-case index—MUST be marked `RECORD-KIND: GROUPING-CONTAINER`. It has only a title, navigation metadata, and an explicit member-ID list; it MUST NOT carry independent semantic claims, `STAGING-STATUS`, `COVERAGE-STATUS`, blueprint/handbook status, `BLOCK-CONFIDENCE-RANK`, `BLOCK-EVIDENCE-PROFILE`, semantic record version/fingerprint, or confirmation, and MUST NOT inherit any such value from children.

A module or bounded context that states responsibility, public interfaces, encapsulation boundary, dependencies, invariants, or data passed is not a grouping container. It is an atomic semantic `MOD` block governed by the common synthesis header, material CLMs, reciprocal references, coverage, semantic version/fingerprint, and confirmation. Module indexes MAY remain purely navigational grouping containers. A `CAP` is the governed capability container defined in §2.6 and §3.2: it may carry its own claim-backed governance metadata, but it remains non-promotable and never inherits member evidence or coverage.

Serialized files are containers. A whole export MUST NOT be represented as one component. A container SRC is `[C-COVERED]` only when its child-enumeration claim is complete for the pinned snapshot and every child coverage row is terminally covered or exactly excluded; it never inherits a child evidence rank/profile and never requires an aggregate promotable CMP. A container with no children requires an evidence-backed zero-child enumeration claim.

## 4.3. Source inventory denominator — `10-SOURCE-INVENTORY.md`

One record exists per concrete source unit:

```markdown
### [SRC-ID] [Name]

- **KIND:** [enumerated kind]
- **Concrete Coordinate:** [exact path:line/symbol or virtual coordinate]
- **Parent / Container SRC:** [SRC-ID or None]
- **Source Fingerprint / Snapshot / SRC Version:** [sha256, snapshot ID, immutable version]
- **Discovery Method:** filesystem | AST | route registry | DB catalog | export explosion | runtime | operator-approved manifest | other
- **Scope Summary:** [derived from scope-specific coverage rows]
- **Environment / Runtime State:** [matrix or reference]
- **Traversal Track(s):** [Normal, Shadow/Conditional, Disabled]
- **Canonical Owner:** [persona or SHARED-ENTRY ID]
- **Disposition Summary:** [derived; never substitutes for scope rows]
- **Mappings:** [CMP, ER, DR, IF, CFG, SCHED, SEC, or other IDs]
- **Frontier References:** [FRONTIER IDs]
- **Coverage Rows:** [references into 16-SOURCE-COVERAGE.md]
- **Coverage Summary:** [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED] | Mixed
- **Claims:** [CLM IDs]
```

Each source fingerprint/snapshot is an immutable SRC version. A changed source creates a new version linked by supersession; it does not overwrite prior snapshot state. Scope, disposition, and coverage authority live in the scope-specific rows in `16-SOURCE-COVERAGE.md`; inventory summaries are derived.

A `Candidate` scope row is temporary and included in closure until changed by an approved scope decision. `Approved-Excluded` requires exact units, rationale, risk, approver, date, and `DEC` ID. Wildcards alone cannot identify exclusions at closure.

### Required inventory kinds

Inventory MUST cover, where applicable: files; symbols; routes; RPC handlers; jobs; queries; tables; views; triggers; stored functions/procedures; generated columns; constraints; configuration keys and feature flags; workflow containers; nodes; edges; pages; widgets; actions; bindings; queues/events; files/storage contracts; and external interfaces.

Absence of a kind MUST be proven by an inventory-backed enumerator record specifying inspected roots/catalogs, method, snapshot, result count zero, limitations, and CLM evidence. An empty section or zero rows without this proof is invalid.

## 4.4. Traversal frontier — `11-TRAVERSAL-FRONTIER.md`

Every discovered traversal boundary MUST create exactly one `FRT` record, including a boundary recognized as terminal, duplicate, excluded, or already mapped at discovery time. Immediate recognition MAY create the record directly in its terminal state, but MUST NOT omit it. Every discovered but not yet terminal traversal boundary remains open until terminally disposed:

```markdown
### [FRONTIER-ID] [Boundary]

- **Origin Invocation / Cluster / POV / Track:** [...]
- **From SRC / CMP:** [...]
- **Target Concrete Coordinate or Resolution Rule:** [...]
- **Target SRC:** [SRC-ID or Pending Inventory]
- **Discovery Method:** call | import | route | graph edge | query reference | DB firing | config gate | runtime | other
- **Depth:** [integer]
- **Snapshot / Fingerprint:** [...]
- **State:** Open | Traversed | Terminal-Mapped | Terminal-Excluded | Terminal-Duplicate | Terminal-Human-Blocked | Superseded
- **Disposition Evidence:** [CLM and mapping/decision IDs]
- **Child Frontier IDs:** [...]
- **Closing Invocation:** [...]
```

Depth bounds create open frontier records; they do not silently terminate traversal. Every frontier MUST be terminal for Exit A. `Terminal-Human-Blocked` remains a coverage gap or partial unless the underlying unit is separately approved as excluded.

## 4.5. Scope-specific source coverage — `16-SOURCE-COVERAGE.md`

One row exists per immutable SRC version × persona/owner × track × environment × snapshot scope. Rows are disjoint within a subgroup; a source participating in several tracks has separate rows and is counted once in each declared track report, while the aggregate unique-SRC report deduplicates by SRC version.

```markdown
| COV ID | SRC Version | Persona/Owner | Track | Environment | Snapshot | Scope | Disposition | C Status | Exclusion DEC | Mappings | Frontier/Ticket/Claim Evidence |
```

The `COV` ID is deterministically derived from the complete immutable scope tuple and MUST be unique for it; duplicate tuple rows are invalid.

`Scope` is `In-Scope`, `Candidate`, or `Approved-Excluded`; disposition is `Unmapped`, `Mapped`, `Superseded`, `Retired`, `Approved-Excluded`, or `Human-Blocked`. The row, not the inventory summary, is the arithmetic input.

## 4.6. Stable supersession and tombstones

Records are never silently deleted. A renamed record retains its ID when identity is proven. A split, merge, or semantic replacement uses:

```markdown
- **RECORD-VERSION:** [integer]
- **SUPERSEDES:** [IDs or None]
- **SUPERSEDED-BY:** [IDs or None]
- **TOMBSTONE:** yes | no
- **Tombstone Reason / Decision:** [...]
```

Tombstones preserve old IDs, incoming references, and snapshot lineage. Validators reject references to tombstones unless the referencing block explicitly records historical use.

---

# Part 5. Status Taxonomy and Evidence

## 5.1. Dedicated prefix groups

Semantic fields MUST NOT be overloaded.

### Staging (`S-`)

- `[S-DISCOVERED]`: atomic component found but not fully verified.
- `[S-STAGING-VERIFIED]`: required component claims and metadata validated; eligible for promotion review.
- `[S-PROMOTED]`: moved to one persona because all promotion locks passed.
- `[S-DEAD-CODE]`: no active, conditional, shadow, disabled, historical, or intended caller is found after complete inventory-backed proof. Disabled does not qualify.
- `[S-RETIRED]`: confirmed removed/superseded from the relevant snapshot, with evidence and disposition.

### Tickets (`T-`)

- Active: `[T-OPEN]`, `[T-PENDING]`, `[T-FOLLOW-UP]`, `[T-PROBE-REQUIRED]`.
- Terminal: `[T-COMPLETE]`, `[T-INVALID]`, `[T-SUPERSEDED]`, `[T-HUMAN-REQUIRED]`, `[T-OUT-OF-SCOPE]`.

`[T-HUMAN-REQUIRED]` can be reached only through the Human Hatch in Part 9. It terminalizes orchestration, not coverage.

### Evidence (`E-`)

- `[E-DIRECT]`: direct static evidence from a source or approved semantics-preserving projection.
- `[E-RUNTIME-OBSERVED]`: captured isolated-runtime telemetry linked to a probe log.
- `[E-CROSS-VALIDATED]`: compatible support from at least two independent evidence routes or POV layers.
- `[E-INFERRED]`: reasoned from naming, patterns, correlation, or incomplete static context.
- `[E-UNKNOWN]`: unsupported, opaque, redacted, removed, dynamic, or unavailable.

### Acquisition (`A-`)

`[A-REQUESTED]`, `[A-RECEIVED]`, `[A-EXPLODED]`, `[A-BLOCKED-MISSING-REFERENCE]`, `[A-STRUCTURALLY-COMPLETE]`.

### Registry/traversal (`R-`)

`[R-ACTIVE]`, `[R-UNSWEPT]`, `[R-IN-PROGRESS]`, `[R-SWEPT]`, `[R-PENDING-HUMAN-REVIEW]`, `[R-PENDING-PERSONA-ASSIGNMENT]`, `[R-DEPROMOTION-PENDING]`.

### Blueprint (`B-`)

- `[B-DRAFT]`: bounded synthesis is incomplete or provisional.
- `[B-POPULATED]`: mandatory content and provenance are present; human confirmation pending.
- `[B-CONFIRMED]`: human-confirmed against a pinned closed snapshot.
- `[B-STALE]`: one or more upstream dependencies changed or became invalid.

### Coverage (`C-`)

- `[C-COVERED]`: the record satisfies the coverage contract for its plane. Forensic/Assurance coverage requires atomic SRC/CMP/CLM/frontier/ticket/export closure and existing internal references; Synthesis/Handbook coverage additionally requires complete higher-level semantics and reciprocal references.
- `[C-PARTIAL]`: some required mapping or evidence exists but that plane's coverage criteria are incomplete.
- `[C-GAP]`: known unit lacks required mapping/evidence/disposition.
- `[C-EXCLUDED]`: removed from the in-scope denominator through an approved exact exclusion.

Capability/runtime state is represented only by the non-prefix fields in Part 2.5.

### Finite-enum field registry and validation

Finite protocol enums MUST be validated from a normative Markdown AST using the complete schema/context-qualified field path, never an unanchored text match, a label-only lookup, or a token-prefix guess. For a bullet field, the parser identifies the parent typed schema, removes Markdown emphasis from the literal field label, requires the field's declared scalar or list shape, and validates each value against the enum assigned to that exact AST path. For a table cell, the parser resolves the parent typed schema, literal table heading, and literal column label before validating the cell. Aliases, invented labels, inherited enums, and applying one schema's label map to another schema are forbidden unless explicitly introduced by a protocol/schema version.

Prefix-family tokens (`S-`, `T-`, `E-`, `A-`, `R-`, `B-`, and `C-`) MUST be bracketed everywhere they are declared below. All other registry literals are exact unbracketed values. The only exceptions involving prefix-family contexts are the exact unbracketed packaging/prepackage literals and the exact unbracketed traversal-matrix literal `N/A` declared below. Multiple scalar tokens, prose mixed into a finite-enum value, an unbracketed prefix-family token, an unknown literal, and an undeclared finite-enum path are validation errors.

The exhaustive registry of all normative finite-choice scalar and table-cell paths defined by this protocol is:

```text
Preflight.Current Iteration                              -> ALFA | BRAVO | CHARLIE | DELTA | ECHO | FOXTROT | GOLF | HOTEL | INDIA | JULIETT | KILO | LIMA | MIKE | NOVEMBER | OSCAR | PAPA | QUEBEC | ROMEO | SIERRA | TANGO | UNIFORM | VICTOR | WHISKEY | X-RAY | YANKEE | ZULU
PersonaRegistry[].Registry State                        -> Active | Conditional | Disabled | Historical | Retired
CapabilityRegistry.CAPABILITY-STATE                     -> Active | Conditional | Shadow | Disabled | Retired | Unknown
CapabilityRegistry.RUNTIME-STATE-MATRIX[].State          -> Active | Conditional | Shadow | Disabled | Retired | Unknown
EntryCluster.Primary POV                                -> POV-1 | POV-2 | POV-3 | POV-4 | POV-5
EntryCluster.Traversal Track                            -> Normal | Shadow/Conditional | Disabled
EntryCluster.Capability State                           -> Active | Conditional | Shadow | Disabled | Retired | Unknown
EntryCluster.Traversal Status                           -> [R-ACTIVE] | [R-UNSWEPT] | [R-IN-PROGRESS] | [R-SWEPT] | [R-PENDING-HUMAN-REVIEW] | [R-PENDING-PERSONA-ASSIGNMENT] | [R-DEPROMOTION-PENDING]
EntryCluster.Required Traversal Matrix[].Track          -> Normal | Shadow/Conditional | Disabled
EntryCluster.Required Traversal Matrix[].POV            -> POV-1 | POV-2 | POV-3 | POV-4 | POV-5
EntryCluster.Required Traversal Matrix[].Applicability  -> Required | Not-Applicable
EntryCluster.Required Traversal Matrix[].Status         -> [R-UNSWEPT] | [R-IN-PROGRESS] | [R-SWEPT] | N/A
UnmappedDiscovery.Traversal Status                      -> [R-PENDING-HUMAN-REVIEW] | [R-PENDING-PERSONA-ASSIGNMENT]
UnmappedDiscoveryBuffer.State                           -> Pending-Merge | Merged | Superseded
Export.Platform                                           -> n8n | Appsmith | Other
Export.Acquisition Status                               -> [A-REQUESTED] | [A-RECEIVED] | [A-EXPLODED] | [A-BLOCKED-MISSING-REFERENCE] | [A-STRUCTURALLY-COMPLETE]
PackagingArtifact.ARTIFACT-TYPE                         -> EXIT-E-CANDIDATE-REPORT | CANDIDATE-PAYLOAD-MANIFEST | EXIT-E-CONTENT-READINESS-REPORT | SCOPE-CERTIFICATE | OUTER-BUNDLE-MANIFEST
SourceInventory.KIND                                    -> FILE | SYM | ROUTE | RPC | JOB | QUERY | TABLE | VIEW | DR | CONSTRAINT | GENERATED | CFG | WORKFLOW | NODE | EDGE | PAGE | WIDGET | ACTION | BINDING | QUEUE | EVENT | STORAGE | IFACE | OTHER
SourceInventory.Discovery Method                        -> filesystem | AST | route registry | DB catalog | export explosion | runtime | operator-approved manifest | other
SourceInventory.Traversal Track(s)[]                    -> Normal | Shadow/Conditional | Disabled
SourceInventory.Coverage Summary                        -> [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED] | Mixed
Frontier.Discovery Method                               -> call | import | route | graph edge | query reference | DB firing | config gate | runtime | other
Frontier.State                                          -> Open | Traversed | Terminal-Mapped | Terminal-Excluded | Terminal-Duplicate | Terminal-Human-Blocked | Superseded
SourceCoverage[].Track                                  -> Normal | Shadow/Conditional | Disabled
SourceCoverage[].Scope                                  -> In-Scope | Candidate | Approved-Excluded
SourceCoverage[].Disposition                            -> Unmapped | Mapped | Superseded | Retired | Approved-Excluded | Human-Blocked
SourceCoverage[].C Status                               -> [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED]
RecordSupersession.TOMBSTONE                            -> yes | no
GroupingContainer.RECORD-KIND                           -> GROUPING-CONTAINER
CapabilityRecord.RECORD-KIND                            -> GOVERNED-CAPABILITY-CONTAINER
AtomicComponent.RECORD-KIND                             -> ATOMIC
AtomicComponent.POV Owner                               -> POV-1 | POV-2 | POV-3 | POV-4 | POV-5
AtomicComponent.STAGING-STATUS                          -> [S-DISCOVERED] | [S-STAGING-VERIFIED] | [S-PROMOTED] | [S-DEAD-CODE] | [S-RETIRED]
AtomicComponent.COVERAGE-STATUS                         -> [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED]
AtomicComponent.TRAVERSAL-TRACKS[]                      -> Normal | Shadow/Conditional | Disabled
AtomicComponent.CAPABILITY-STATE                        -> Active | Conditional | Shadow | Disabled | Retired | Unknown
AtomicComponent.RUNTIME-STATE-MATRIX[].State             -> Active | Conditional | Shadow | Disabled | Retired | Unknown
AtomicComponent.BLOCK-CONFIDENCE-RANK                    -> None | Low | Medium | High
AtomicComponent.BLOCK-EVIDENCE-PROFILE[].Evidence Level  -> [E-DIRECT] | [E-RUNTIME-OBSERVED] | [E-CROSS-VALIDATED] | [E-INFERRED] | [E-UNKNOWN]
SharedReference.Track                                    -> Normal | Shadow/Conditional | Disabled
SharedReference.Capability State                         -> Active | Conditional | Shadow | Disabled | Retired | Unknown
SharedReference.Canonical Staging Status                -> [S-DISCOVERED] | [S-STAGING-VERIFIED]
Claim.Materiality                                       -> Material | Supporting
Claim.Evidence Level                                    -> [E-DIRECT] | [E-RUNTIME-OBSERVED] | [E-CROSS-VALIDATED] | [E-INFERRED] | [E-UNKNOWN]
Claim.Projection Classification                         -> raw-visible | preserved | normalized-no-semantic-change | aliased | redacted | removed | unclassified
Contradiction.Resolution State                          -> Open | Explained-By-Scope | Resolved-By-Evidence | Decision-Required | Superseded
Decision.Type                                           -> scope | contradiction | target | migration | capability
Decision.Approval Status                                -> Proposed | Approved | Superseded
CapabilityRecord.Target Decision State                  -> Bound | Pending
Confirmation.Human Hatch Action                         -> Confirmation
Confirmation.Handbook Readability Review                 -> Pass | Fail | Not-Applicable
Confirmation.Invalidated By                             -> semantic change | dependency change | decision supersession | None
AcquisitionCandidates[].State                           -> Included | Rejected-With-Operator-Decision | Approved-Excluded | Pending
ExportReconciliation[].Reconciliation State             -> Mapped-Atomic | Mapped-Container-Only | Approved-Excluded | Gap | Superseded
InvocationHeader.Current Iteration                      -> ALFA | BRAVO | CHARLIE | DELTA | ECHO | FOXTROT | GOLF | HOTEL | INDIA | JULIETT | KILO | LIMA | MIKE | NOVEMBER | OSCAR | PAPA | QUEBEC | ROMEO | SIERRA | TANGO | UNIFORM | VICTOR | WHISKEY | X-RAY | YANKEE | ZULU
InvocationHeader.Invocation Mode                        -> Preflight | Export Acquisition | Discovery | Ticket Resolution | Promotion Review | Sequential Reconciliation | Profile Synchronization | Cross-Reference-Reconciliation | Partial-Synthesis | Final-Synthesis | Reconstruction-Handoff | Human Hatch | Validation/Gate
InvocationHeader.Human Hatch Action                     -> Ticket Escalation | Decision/Scope Approval | Confirmation | None
InvocationHeader.Traversal Track                        -> Normal | Shadow/Conditional | Disabled | None
InvocationHeader.Active POV                             -> POV-1 | POV-2 | POV-3 | POV-4 | POV-5 | None
ColdResume.Runtime Scope / Mode / Write-Right Check     -> Pass | Fail
ColdResume.Runtime Stale Check                          -> Pass | Fail
ColdResume.Cold-Resume Result                           -> Pass | Fail
Ticket.Status                                           -> [T-OPEN] | [T-PENDING] | [T-FOLLOW-UP] | [T-PROBE-REQUIRED] | [T-COMPLETE] | [T-INVALID] | [T-SUPERSEDED] | [T-HUMAN-REQUIRED] | [T-OUT-OF-SCOPE]
Ticket.Escalation Reason                                -> None | RUNTIME-REQUIRED | THIRD-PARTY-OPAQUE | UNAUTHORIZED-ACCESS | UNSAFE | DESTRUCTIVE | NO-SANDBOX | NO-AUTHORITY-TO-DECIDE | DOMAIN-MEANING-REQUIRED | EXTERNAL-FACT-UNAVAILABLE
Probe.Safety Classification                             -> non-destructive | controlled mutation
SynthesisGap.Materiality                                -> Material | Supporting
SynthesisGap.Disposition                                -> Pending | Ticket-Created | Frontier-Created | Both-Created | Merged | Invalid | Approved-Excluded
SynthesisBlock.COVERAGE-STATUS                          -> [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED]
SynthesisBlock.BLOCK-CONFIDENCE-RANK                    -> None | Low | Medium | High
SynthesisBlock.BLOCK-EVIDENCE-PROFILE[].Evidence Level  -> [E-DIRECT] | [E-RUNTIME-OBSERVED] | [E-CROSS-VALIDATED] | [E-INFERRED] | [E-UNKNOWN]
SynthesisBlock.BLUEPRINT-STATUS                         -> [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
PersonaProfile.BLOCK-CONFIDENCE-RANK                    -> None | Low | Medium | High
PersonaProfile.BLOCK-EVIDENCE-PROFILE[].Evidence Level  -> [E-DIRECT] | [E-RUNTIME-OBSERVED] | [E-CROSS-VALIDATED] | [E-INFERRED] | [E-UNKNOWN]
PersonaProfile.BLUEPRINT-STATUS                         -> [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
HandbookSection.HANDBOOK-STATUS                         -> [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
Entity.Storage                                           -> table | view | collection | cache-key-space | file | in-memory
Relationship.Cardinality                                -> 1:1 | 1:N | M:N
Relationship.Enforcement                                -> DB-FK | application | mixed | none
Relationship.Delete / Update Semantics                  -> cascade | restrict | nullify | custom
DatabaseRoutine.Kind                                    -> trigger | function | procedure | rule | generated-column | materialized-refresh | constraint-helper
DatabaseRoutine.Target Decision                         -> Retain | Lift | Redesign | Retire | Undecided
BusinessRule.Violation Policy                           -> reject | warn | clamp | default | compensate | rollback
UseCase.Capability Classification                       -> Active | Conditional | Shadow | Disabled | Retired
UseCase.Target Decision                                 -> Preserve | Change-By-Approved-Migration | Redesign | Retire
BusinessRule.Capability State                           -> Active | Disabled
ScheduledJob.Capability State                           -> Active | Conditional | Shadow | Disabled | Retired | Unknown
Deployment.Single Point of Failure                      -> yes | no | unknown
Configuration.Medium                                    -> env | DB | flag service | file | build-time
NFR.Attribute                                            -> latency | throughput | availability | RPO | RTO | scalability | operability | other
Security.Category                                       -> auth | authorization | transport | at-rest | in-transit | PII | secrets | audit | threat | privacy
Security.Risk Level                                     -> high | medium | low | unassessed
Traceability[].C Status                                 -> [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED]
Traceability.Lifecycle Eligibility[].Classification     -> Eligible | Ineligible | Unknown
ValidationSummary.Evidence Bindings[].Binding Kind       -> RAW-SOURCE-BYTES | RECORD-PAYLOAD | SEMANTIC-CONTENT | DECISION-CONTENT | CNF-PAYLOAD | ARTIFACT-PAYLOAD | CANONICAL-ENVELOPE | FILE-TRANSPORT | PACKAGE-MEMBER-FILE | VALIDATION-SUMMARY | DENOMINATOR-QUERY | PROTOCOL-REGISTRY
ExitECandidateReport.Check Results[].Check ID             -> §15.1.1 `EXIT-E-CHECKS-v1` rows with Stage `CANDIDATE`
ExitECandidateReport.Check Results[].Result             -> Pass | Fail
ExitECandidateReport.Candidate Result                   -> Ready-For-Human-Review | Blocked
ExitEContentReadinessReport.Check Results[].Check ID     -> §15.1.1 `EXIT-E-CHECKS-v1` rows with Stage `CONTENT-READINESS`
ExitEContentReadinessReport.Check Results[].Result      -> Pass | Fail
ExitEContentReadinessReport.EXIT-E-PREPACKAGE-STATE     -> Pending
ExitEContentReadinessReport.Content Readiness Result    -> Ready-For-Certificate | Blocked
OuterBundleManifest.EXIT-E-STATUS                       -> Passed
FinalVerificationReceipt.Check Results[].Check ID       -> §15.5 `EXIT-E-FINAL-CHECKS-v1` Check ID rows
FinalVerificationReceipt.Check Results[].Result         -> Pass | Fail
FinalVerificationReceipt.Overall Result                 -> Pass | Fail
```

`EntryCluster.Required Traversal Matrix[].Status` accepts `N/A` if and only if the same row's `Applicability` is exactly `Not-Applicable`; such a row MUST also carry the claim-backed reason required by §3.1. A `Required` row MUST use one bracketed allowed `R-` token and MUST NOT use `N/A`. A `Not-Applicable` row MUST use exactly `N/A` and MUST NOT use an `R-` token.

Fields that merely contain descriptive prose or domain-defined values are not finite protocol enums. In particular, `PersonaProfile.UI Surfaces and State`, state names inside an eligible `SM`, free-form environment names, and narrative state summaries MUST NOT be dispatched through this registry. Conversely, any field intended to carry a finite protocol enum must have a declared path above; encountering a finite-enum-like value at an undeclared path is an error rather than permission to infer an enum from its label or token.

## 5.2. Claim-level evidence — `12-CLAIM-EVIDENCE.md`

Each claim MUST contain exactly one fact and exactly one evidence level. Mixed evidence in one claim is forbidden; split it into multiple claims.

```markdown
### [CLM-ID]

- **Subject ID / Record Version:** [SRC/CMP/ER/etc.]
- **Field / Step:** [exact field name, table cell, UC step, or handbook assertion]
- **Claim:** [one falsifiable statement]
- **Materiality:** Material | Supporting
- **Evidence Level:** [E-DIRECT] | [E-RUNTIME-OBSERVED] | [E-CROSS-VALIDATED] | [E-INFERRED] | [E-UNKNOWN]
- **Anchor:** [physical path:line, virtual coordinate, probe log range, catalog query result, or approved projection pointer]
- **Source Fingerprint / Snapshot:** [...]
- **Evidence Route / Lineage ID:** [independent acquisition/static/runtime route]
- **Projection Classification:** raw-visible | preserved | normalized-no-semantic-change | aliased | redacted | removed | unclassified
- **Independent Supporting Claims:** [CLM IDs or None]
- **Independence Basis:** [different raw lineage, different runtime observation, different physical enforcement, or None]
- **Contradiction IDs:** [CONTRADICTION IDs or None]
- **Created / Last Revalidated By:** [Invocation IDs]
```

### Block evidence summary

Every atomic forensic or synthesis block lists its material claims and carries both fields below:

```markdown
- **BLOCK-CONFIDENCE-RANK:** None | Low | Medium | High
- **BLOCK-EVIDENCE-PROFILE:**
  | Evidence Level | Count | Claim IDs | Route/Lineage IDs |
```

`BLOCK-CONFIDENCE-RANK` is the minimum material-claim rank under this total mapping: `[E-UNKNOWN] -> None`, `[E-INFERRED] -> Low`, `[E-CROSS-VALIDATED] -> Medium`, and both `[E-DIRECT]` and `[E-RUNTIME-OBSERVED] -> High`. The profile MUST retain separate rows, counts, claim IDs, and route/lineage IDs for all five evidence levels; direct and runtime-observed evidence MUST NOT be collapsed into an invented evidence tag. `[E-CROSS-VALIDATED]` requires at least two supporting claims whose route/lineage IDs are independent. Claims sharing a raw fingerprint or projection lineage, duplicated from the same runtime trace, or derived from one another are not independent even when written by different POVs. A block with no material claims is invalid. Purely navigational grouping containers have neither field.

## 5.3. Sanitized projections and helper limitations

A sanitized export produced by an approved local helper is an agent-visible evidence projection. It MUST include raw-source fingerprint lineage and a sanitization report classifying preserved, normalized, aliased, redacted, removed, and unclassified content.

`[E-DIRECT]` is permitted only for preserved or semantics-preservingly normalized behavior. Aliases prove reference existence, not secret value, connectivity, authorization, or runtime validity. Behavior dependent on redacted, removed, or unclassified content remains `[E-INFERRED]` or `[E-UNKNOWN]` and creates a ticket when material.

Entity-correlation output is a candidate-discovery aid only. It MUST NOT become evidence, a call edge, ownership proof, or automatic inclusion without a direct structural signal or explicit human scope decision. Private matching hints and private entity-usage indexes MUST NOT enter agent-visible inventory.

## 5.4. Contradictions — `14-CONTRADICTIONS.md`

Conflicting material claims MUST NOT be silently resolved by selecting one. Record:

```markdown
### [CONTRADICTION-ID]

- **Claims in Conflict:** [CLM IDs]
- **Scope / Environments / Snapshots:** [...]
- **Conflict Description:** [...]
- **Resolution State:** Open | Explained-By-Scope | Resolved-By-Evidence | Decision-Required | Superseded
- **Resolution / Decision ID:** [...]
- **Affected Records and Stale Impact:** [...]
```

An open material contradiction blocks applicable coverage, Exit A, and confirmation.

## 5.5. Authoritative decisions — `15-DECISIONS.md`

`15-DECISIONS.md` exists from Preflight onward and is the source of truth for scope exclusions, contradiction resolutions, disabled-capability choices, migration choices, and reconstruction decisions. `Human Hatch: Decision/Scope Approval` may append decision content or attach approval envelopes; Reconstruction-Handoff projects approved records into the handbook but does not create a second semantic log.

```markdown
### [DEC-ID] [Decision]

- **DECISION-CONTENT-VERSION:** [positive integer]
- **DECISION-CONTENT-FINGERPRINT:** [sha256 of canonical decision content under §4.1.2; this carrier field excluded]
- **Type:** scope | contradiction | target | migration | capability
- **Observed Inputs:** [SRC/CLM/ticket/contradiction IDs and fingerprints]
- **Options / Decision / Rationale / Risks:** [...]
- **Exact Scope and Denominator Effect:** [...]
- **Affected Records / Rollout / Rollback:** [...]
- **Content Supersession:** [...]

<!-- DECISION-APPROVAL-ENVELOPE: excluded from decision-content version/fingerprint -->

- **Approval Status:** Proposed | Approved | Superseded
- **Approval Envelope Version:** [positive integer]
- **Owner / Approver / Authority Basis / Approval Timestamp:** [...]
- **Approved Content Binding:** [DEC-ID + exact decision-content version/fingerprint]
- **Approval Signature / Audit:** [...]
<!-- END-DECISION-APPROVAL-ENVELOPE -->
```

The candidate payload manifest binds immutable decision-content versions/fingerprints. Each `DEC` is canonicalized as an individual heading-delimited record under §4.1.2 and excludes exactly its fingerprint carrier and the literal decision-approval envelope shown above. Changing options, decision, rationale, risks, scope effect, rollout, rollback, or observed inputs increments the decision-content version and changes its fingerprint. Human approval after candidacy attaches only the approval envelope and does not change decision content or its fingerprint; approving content other than the candidate binding is invalid. Every decision reference includes the content version/fingerprint and, when approval is required, the matching approval-envelope status.

## 5.6. Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md`

Every helper-published candidate and operator action receives a row:

```markdown
| CND ID | Relationship Class | Sanitized Signals | Dynamic Reference | Decision | Operator/Date | Resulting EXP/SRC/DSC/DEC IDs | State |
```

Each helper-published candidate or operator action receives a deterministic `CND` ID; its canonical key binds the stable candidate/action identity, helper policy version, and pinned snapshot. Duplicate candidate/action identities or reused CND IDs are invalid.

States are `Included`, `Rejected-With-Operator-Decision`, `Approved-Excluded`, or `Pending`. The helper manifest MUST attest that all accessible workflow metadata was scanned for the pinned snapshot, all direct structural matches were included, all high-confidence entity-only matches were human-reviewed, and all configured medium-confidence candidates received an operator decision; scoring thresholds and policy versions are recorded. Lower-confidence omissions follow the configured, approved scope policy. Helper publication completeness means safe approved outputs were produced; protocol `[A-STRUCTURALLY-COMPLETE]` additionally requires every candidate/operator action/dynamic reference row terminal. Correlation remains discovery/scope input only, never behavioral evidence or an automatic edge.

## 5.7. Human confirmations — `18-CONFIRMATIONS.md`

```markdown
### [CNF-ID] [Confirmation]

- **Human Hatch Action:** Confirmation
- **Approver / Role / Authority Basis:** [...]
- **Review Timestamp / Review Scope / Limitations:** [...]
- **Handbook Readability Review:** Pass | Fail | Not-Applicable
- **Confirmed Record Bindings:**
  | Typed Record ID (including PRF/HBK) | SEMANTIC-RECORD-VERSION | SEMANTIC-CONTENT-FINGERPRINT | Candidate Manifest Entry |
- **Pinned Exit A / Exit E Candidate Report / Candidate Payload Manifest Fingerprints:** [...]
- **Required DEC Content Versions / Fingerprints and Approval Envelopes Reviewed:** [...]
- **Invalidated By:** semantic change | dependency change | decision supersession | None

<!-- CNF-DIGEST-SIGNATURE-ENVELOPE: excluded from CNF payload hash -->

- **CNF Signature Envelope Version:** [positive integer]
- **CNF Payload Fingerprint:** [sha256 of the CNF payload under §4.1.2; this carrier field excluded]
- **Signature or Approval Reference / Signature Timestamp:** [...]
<!-- END-CNF-DIGEST-SIGNATURE-ENVELOPE -->
```

A B-confirmed synthesis block, `PRF` profile, or `HBK` handbook section MUST reference a valid `CNF` whose binding matches its exact typed record ID, semantic record version, and semantic-content fingerprint as listed in the candidate manifest. Confirmation validation MUST load and validate the same direct-and-transitive dependency closure required for Profile Synchronization in §8.4 when the confirmed record is a PRF, and the corresponding complete `DERIVED-FROM` closure for synthesis/HBK records. Confirmation metadata is an envelope attachment and MUST NOT alter the confirmed semantic payload/version/fingerprint. Dependency or semantic change invalidates the confirmation through targeted stale propagation; status-only and confirmation-envelope changes do not. Each `CNF` is canonicalized as an individual heading-delimited record under §4.1.2. Its payload fingerprint excludes exactly the literal CNF digest/signature envelope shown above and its own fingerprint carrier; review identity, authority, scope, limitations, typed record bindings, and review timestamp remain hashed.

---

# Part 6. Atomic Forensic Extraction and Evolution

## 6.1. Atomic component schema

```markdown
### [CMP-ID] [Atomic Unit Name]

- **RECORD-KIND:** ATOMIC
- **Source IDs:** [SRC IDs]
- **File Path & Line Range:** [...]
- **Virtual Coordinate:** [required for serialized units]
- **POV Owner:** POV-1 | POV-2 | POV-3 | POV-4 | POV-5
- **General Intent:** [...]
- **Inputs / Outputs:** [...]
- **AUTH-BINDING:** [AUTH IDs or None]
- **STAGING-STATUS:** [S-*]
- **COVERAGE-STATUS:** [C-*]
- **PERSONAS:** [persona prefixes]
- **TRAVERSAL-TRACKS:** [Normal, Shadow/Conditional, Disabled]
- **CAPABILITY-STATE:** Active | Conditional | Shadow | Disabled | Retired | Unknown
- **RUNTIME-STATE-MATRIX:** [environment/snapshot state rows under §2.5]
- **DEPENDS-ON / DEPENDED-ON-BY:** [CMP/SRC IDs]
- **Synthesis References:** [IF/UC/BR/ER/DR/CFG/SCHED IDs; populated by reconciliation]
- **Side Effects:** [...]
- **Exceptional Behaviors:** [...]
- **Edge and Odd Cases:** [...]
- **Behavioral Notes for Reconstruction:** [deterministic observed contract, not target design]
- **Material Claims:** [CLM IDs]
- **BLOCK-CONFIDENCE-RANK:** None | Low | Medium | High
- **BLOCK-EVIDENCE-PROFILE:** [counts, claim IDs, and route/lineage IDs by each E level]
- **Evidence Anchors:** [convenience list; CLMs remain authoritative]
- **Record Version / Supersession:** [...]
```

All fields are atomic-unit-specific. A routine block MUST NOT aggregate unrelated methods, branches, graph nodes, or queries merely because they share a module.

## 6.2. Dual-entry discovery

Persona-driven discovery is the only behavioral discovery path. There is no independent shared-code sweep.

When a persona entry discovers a component:

1. create the full canonical atomic record in the persona's `SHARED-STAGING-BUFFER.md` for later merge to `personas/_shared/XX-*-POV.md` with `[S-DISCOVERED]`;
2. immediately create a persona invocation pointer in `personas/{persona-directory}/XX-*-POV.md`.

```markdown
### [CMP-ID] [Name] — Shared Reference

- **Canonical Location:** `personas/_shared/XX-*-POV.md` -> [CMP-ID]
- **Local Trigger SRC / Coordinate:** [...]
- **Local Invocation Signature:** [...]
- **Track:** Normal | Shadow/Conditional | Disabled
- **Capability State:** Active | Conditional | Shadow | Disabled | Retired | Unknown
- **Canonical Staging Status:** [S-DISCOVERED] | [S-STAGING-VERIFIED]
```

This ensures a persona manifest represents its full operational footprint without duplicating canonical behavior. If discovery reaches a coordinate that cannot yet be assigned to the bound persona/cluster/track/POV, the same transaction MUST create its `FRT` and append a `DSC` entry to `personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md`; the invocation MUST NOT write `0A` directly or continue through the unmapped boundary. Sequential Reconciliation MUST merge all such pending entries before any later Discovery invocation begins, as specified in §3.1.

## 6.3. Staging and promotion

A shared atomic component advances from `[S-DISCOVERED]` to `[S-STAGING-VERIFIED]` only when mandatory fields, claim evidence, source mappings, references, and ticket state validate.

Promotion eligibility requires all locks:

1. **Usage uniqueness:** exactly one persona, considering Normal, Shadow/Conditional, Disabled, and historical/intended usage within scope.
2. **Ambiguity:** no active ticket or open material contradiction references the component.
3. **Dependency:** no canonical shared component depends on it.
4. **Coverage:** source units are `[C-COVERED]` for the promoted component.
5. **Evidence:** no material `[E-UNKNOWN]` claim and no material unresolved redaction dependency.

`Promotion Review` validates these locks but MUST NOT move or rewrite canonical blocks. It appends one `PROMOTION-REQUEST` to the bound persona's `personas/{persona-directory}/SHARED-STAGING-BUFFER.md`, containing the CMP/cohort IDs, target persona/POV, source `_shared` paths and complete file/record fingerprints, lock-evidence IDs, review invocation, and requested `[S-PROMOTED]` state. Connected atomic components MAY be requested as one cohort only when their combined persona set is singular and all locks pass. Grouping containers never promote.

Sequential Reconciliation is the sole promotion executor. It loads the complete request, canonical `personas/_shared/` source blocks, target persona POV files, indexes, traceability, and current lock inputs; applies the stale-check guard; revalidates every lock; and atomically moves each full canonical block from `_shared` to the target persona file, changes it to `[S-PROMOTED]`, updates or removes shared-reference pointers as applicable, updates indexes/traceability, consumes the request, and logs one transaction. A failed or stale revalidation performs no partial move and leaves a dispositioned request for renewed review. At no point may both source and target contain canonical copies, nor may neither contain one.

`[S-DEAD-CODE]` requires an inventory-backed zero-caller proof across every traversal track and every environment/snapshot in the component's bound scope denominator. Any linked nonterminal dynamic-reference or acquisition-candidate row, active ticket, open frontier, unresolved possible-caller edge, or discovered active/conditional/shadow/disabled/historical/intended caller makes `[S-DEAD-CODE]` invalid. A dynamically observed caller disproves dead-code classification. An unavailable caller domain may be removed from the proof denominator only by an exact approved exclusion; it is never equivalent to proven no caller. These effects are local to the linked component and its dependency closure and MUST NOT downgrade unrelated coverage.

## 6.4. Depromotion

If a promoted component gains another persona caller, a shared dependency, or expanded persona usage:

1. during parallel work, record `[R-DEPROMOTION-PENDING]` in the discovering persona's staging buffer;
2. Sequential Reconciliation moves the full block to `personas/_shared/`, resets it to `[S-STAGING-VERIFIED]`, merges personas and claims, and replaces persona copies with shared references;
3. log the event in `0G` and update traceability and stale dependencies.

## 6.5. Surgical evolution

POV files MUST be patched surgically. Unrelated blocks MUST NOT be deleted, reordered, or reformatted. Ticket resolution updates affected assumptions and cites the ticket. Superseded facts remain traceable through versions/tombstones. Clean-sweep rewrites are prohibited except an explicitly approved protocol migration that preserves all IDs and history.

---

# Part 7. Serialized Export Acquisition and Local Helper Trust Model

## 7.1. Trust boundary

Live acquisition is an operator-side activity. The local helper MAY read raw Appsmith exports and approved read-only n8n interfaces, build private entity indexes, sanitize, fingerprint, and publish approved projections. Agents MUST NOT invoke acquisition using production credentials.

The helper MUST:

- default to dry-run;
- accept credentials only through environment, stdin, or OS credential storage, never command arguments;
- keep raw exports and private matching state outside the target repository and authoritative `.extracted/` workspace, in operator-controlled storage not published to executor-visible projections;
- use restrictive temporary storage and best-effort cleanup;
- perform field-aware sanitization and secret scanning;
- fail closed for unsupported schema versions, unknown sensitive fields, unresolved high-confidence secret findings, missing lineage, mixed snapshots, or rejected operator review;
- make no mutation to Appsmith, n8n, or production;
- make no third-party transmission of exports, secrets, private hints, or repository content;
- publish atomically only approved sanitized inventory and bridge fragments.

Agents MAY inspect helper source and approved outputs. Bridge fragments are proposed patches requiring authorized approval before merge.

## 7.2. Acquisition register

```markdown
### [EXPORT-ID] [Artifact]

- **Platform:** n8n | Appsmith | Other
- **Stable Artifact ID / Display Name:** [...]
- **Raw Source Alias:** [non-sensitive alias]
- **Agent-Visible Export Path:** [...]
- **Raw Source Fingerprint / Projection Fingerprint:** [sha256]
- **Snapshot:** [...]
- **Acquisition Status:** [A-*]
- **Referenced By:** [virtual coordinates]
- **Outstanding Referenced Artifacts:** [...]
- **Normalized Map:** [...]
- **Reconciliation Summary:** [mapped/excluded/total]
- **Sanitization / Lineage Reports:** [...]
- **Operator Action Required:** [...]
```

## 7.3. Logical explosion and virtual coordinates

Before POV extraction, each export is decomposed into atomic logical units. Coordinates use:

```text
[source-export]#[logical-kind]/[stable-id]--[normalized-name]
```

Stable internal IDs take precedence. Name-only identities receive deterministic disambiguators and remain low-confidence until verified.

Examples:

```text
inventory/n8n/WF-17.json#node/42--Webhook_PaymentCallback
inventory/n8n/WF-17.json#edge/42->58--success
inventory/appsmith/app.json#page/P1--PaymentPage
inventory/appsmith/app.json#widget/W7--PaymentPage.PayButton
inventory/appsmith/app.json#action/A9--PaymentPage.PayButton.onClick
inventory/appsmith/app.json#query/Q3--GetDebtList
inventory/appsmith/app.json#binding/B6--PaymentPage.Table1.data
```

Virtual coordinates are canonical for entry coordinates, evidence, dependencies, frontiers, and tickets. Physical container fingerprint and virtual coordinate MUST both be preserved.

## 7.4. Required explosion coverage

For n8n, inventory workflow metadata, triggers, webhooks, schedules, nodes, code/function nodes, transformations, conditions, merges, retries, error branches, credentials references, database/cache/file I/O, external calls, sub-workflows, configured error workflows, and every graph edge.

For Appsmith, inventory pages, widgets, bindings, JSObjects, actions, queries, datasources, APIs, event handlers, action chains, validations, navigation, permissions, environment references, timers, polling, and deferred behavior. A complete Appsmith application export may be complete as a physical container, but that does not prove runtime datasource values, plugin behavior, secrets, environment configuration, or external dependencies. Those remain separately evidenced, ticketed, acquired, or exactly excluded.

Primary ownership follows the POV definitions. Every graph edge is an atomic POV-5 record, even when its payload meaning creates tickets for another POV.

## 7.5. n8n human acquisition loop

1. **Seed:** register available inventory/exports as `[A-RECEIVED]` with stable IDs and fingerprints.
2. **Explosion:** generate source inventory, normalized units, entry candidates, and edges; transition to `[A-EXPLODED]`.
3. **Reference discovery:** inspect static sub-workflows, error workflows, IDs/names in expressions/code, queue/event links, and webhook relationships.
4. **Immediate placeholder:** every missing referenced workflow receives an `[A-REQUESTED]` register entry immediately with exact known identity, referring virtual coordinate, and operator action. The referrer becomes `[A-BLOCKED-MISSING-REFERENCE]`. Unknown ownership also creates an unmapped discovery. Traversal stops at the boundary; missing behavior is not hypothesized.
5. **Operator export:** stage a sanitized projection without raw secrets.
6. **Cold resume:** fingerprint, explode, link, restore the referrer to `[A-EXPLODED]`, and resume its open frontier.
7. **Closure:** transition to `[A-STRUCTURALLY-COMPLETE]` only after static reference closure, 100% normalized reconciliation, and terminal disposition of every published helper candidate, operator action, inaccessible artifact, and dynamic reference in `17-ACQUISITION-CANDIDATES.md`. A dynamic reference may close only through inclusion/mapping or an exact approved scope decision; recording it as unresolved is not protocol closure.

Dynamic references are recorded as unresolved. After persona/POV assignment, static investigation creates a ticket; they are never guessed.

## 7.6. Normalized maps

Normalized maps list metadata, each virtual unit, stable ID, type, primary POV, source fingerprint, edges, external systems, credential aliases, and ambiguities. They are navigation-only indexes. They MUST NOT summarize an export as one routine or bypass persona-driven discovery.

## 7.7. Export reconciliation — `13-EXPORT-RECONCILIATION.md`

Every normalized coordinate has one row:

```markdown
| Export ID | Normalized Coordinate | SRC ID | Atomic CMP/POV-5 Edge ID | Exclusion Decision | Reconciliation State | CLM |
```

Allowed states: `Mapped-Atomic`, `Mapped-Container-Only` (temporary and non-closing), `Approved-Excluded`, `Gap`, `Superseded`. A container mapping does not reconcile its children. `Mapped-Atomic + Approved-Excluded == Total Normalized Units` is required for 100% reconciliation.

100% reconciliation blocks `[A-STRUCTURALLY-COMPLETE]`, `[R-SWEPT]`, and Exit A. Entity correlations never satisfy reconciliation.

---

# Part 8. Invocation Isolation, Concurrency, and Reconciliation

## 8.1. Resume identity header

Every invocation begins with:

```yaml
Protocol Version: v4.1.1
System Namespace: [slug]
Current Iteration: [canonical §10.6 token ALFA..ZULU]
Invocation ID: [globally unique]
Invocation Mode: [enum below]
Human Hatch Action:
  [
    Ticket Escalation | Decision/Scope Approval | Confirmation | None; non-None only when Invocation Mode is Human Hatch,
  ]
Invocation Scope: [exact scope]
Persona: [exactly one or None where mode permits]
Persona Prefix: [unique prefix or None]
Persona Directory:
  [exact `<persona-prefix>-<persona-slug>` basename or None where mode permits]
Entry Cluster: [exactly one concrete cluster or None where mode permits]
Traversal Track: [Normal | Shadow/Conditional | Disabled | None]
Active POV: [POV-1..POV-5 | None]
Active POV File: [exactly one or None]
Active Question Ledger: [zero or one]
Environment/Snapshot: [...]
Max Traversal Depth: [integer]
Allowed Read Targets: [...]
Allowed Write Targets: [...]
Forbidden Write Targets: [...]
Loaded Checkpoint Fingerprint: [...]
```

The executor MUST declare its intended semantic actions against this header before writing. The conforming runtime MUST independently validate the declared actions, exact read/write scope, loaded fingerprints, checkpoint freshness, and mode rights before commit; an executor's self-assessment cannot authorize mutation.

## 8.2. Strict single-scope modes

`Discovery`, `Ticket Resolution`, and `Promotion Review` MUST bind exactly one persona, its persona prefix, its canonical persona directory, one concrete entry cluster, one traversal track, one POV, one POV file, and at most one question ledger. Each has one semantic POV target and optional one semantic ledger target. Its fixed transactional sidecar bundle MAY append only directly produced rows/blocks in `10`, `11`, `12`, `14`, `16`, persona-local shared buffers, `personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md`, and `0G`; these are assurance/audit side effects, not permission to change scope or another POV's semantics. Cross-POV outputs go to the bound ledger or persona-local buffers. An unmapped discovery MUST use the local buffer and requires Sequential Reconciliation before any later Discovery invocation.

A cluster containing wildcards MUST be expanded into concrete coordinates before a closing invocation. One invocation may traverse multiple atomic units reachable from its one concrete root, subject to depth/frontier rules, but may not switch roots, personas, POVs, or tracks.

### ALFA start behavior

ALFA is the first bounded discovery iteration. Every POV may log questions immediately; no maturity phase may suppress cross-boundary ambiguity. On large repositories, ALFA MUST still obey concrete cluster, track, persona, POV, and depth isolation. Branches beyond the bound become frontier records and tickets rather than an attempted full-tree sweep.

## 8.3. Invocation-mode enum and explicit multi-file modes

The exhaustive invocation-mode enum is:

```text
Preflight | Export Acquisition | Discovery | Ticket Resolution | Promotion Review |
Sequential Reconciliation | Profile Synchronization | Cross-Reference-Reconciliation |
Partial-Synthesis | Final-Synthesis | Reconstruction-Handoff | Human Hatch | Validation/Gate
```

`Preflight` has one semantic target (`0A`); `Discovery`, `Ticket Resolution`, and `Promotion Review` obey §8.2; `Validation/Gate` writes one report target. All MAY append `0G` and fixed assurance sidecars allowed by their contracts. Probe ingestion executes as `Ticket Resolution`. Neutral arbitration executes as `Ticket Resolution` with the original persona/cluster/track and a recorded neutral POV assignment. Protocol migration is metadata on a sequence of existing modes, not a separate mode. Human confirmation executes only as `Human Hatch: Confirmation` or the equivalent explicitly authorized Reconstruction-Handoff confirmation-envelope action and writes `18-CONFIRMATIONS.md` plus exact envelope transitions without semantic mutation.

Only these modes may intentionally span multiple semantic targets:

- `Export Acquisition`: acquisition register, inventory/reconciliation/candidate sections, normalized maps, security findings, and `0G`.
- `Sequential Reconciliation`: canonical `personas/_shared/` merges plus persona-local shared-staging, promotion-request, question, unmapped-discovery, and synthesis-gap buffer merges; atomic promotion and depromotion; buffered ticket/frontier creation; deterministic deduplication; indexes; traceability; frontier dispositions; and `0G`.
- `Profile Synchronization`: exactly one persona profile semantic payload plus its atomic `0A` registry/ownership updates and `0G`; an authorized confirmation action may later update only its certification envelope.
- `Cross-Reference-Reconciliation`: only already-existing ID reference fields, reciprocal links, traceability, and `0G`; no new legacy facts. Changed semantic reference bindings increment semantic versions/fingerprints.
- `Partial-Synthesis`: bounded provisional `90`–`96` blocks and a synthesis gap buffer.
- `Final-Synthesis`: closed-corpus `90`–`96` blocks and synthesis gap buffer.
- `Reconstruction-Handoff`: populated handbook payloads, approved-decision projection, authorized confirmation envelopes, equivalence suite, candidate/outer payload manifests, scope certificate, and package assembly in the exact §15.4 stages; it does not write Validation/Gate reports.
- `Human Hatch` has three mandatory action subtypes in its identity header:
  - `Ticket Escalation`: after §9.3 prerequisites, terminalizes the ticket, records exact coverage/scope/risk effects, and MAY create an exact exclusion decision when authorized;
  - `Decision/Scope Approval`: creates, approves, or supersedes `DEC` records only after the reviewer receives exact observed inputs, affected denominator tuples, risks, alternatives, authority basis, and pinned fingerprints; it has no probe prerequisite unless the decision itself resolves a ticket through Ticket Escalation;
  - `Confirmation`: creates `CNF` records and attaches certification envelopes only after candidate validation succeeds and the reviewer receives exact semantic versions/fingerprints, candidate manifest, required decision-content bindings and approval envelopes, limitations, and authority basis; it has no probe prerequisite and cannot edit semantic payloads.
- `Human Hatch` writes only the records authorized by its declared subtype, plus exact coverage/traceability effects and `0G`; a subtype cannot borrow another subtype's rights.

Each multi-file mode MUST enumerate exact files and record sections in its identity header; the mode name alone grants no broad write access. `Validation/Gate` reads authoritative inputs and writes only its named report block plus `0G`. Human confirmation is an authorized `Human Hatch: Confirmation` or `Reconstruction-Handoff` confirmation-envelope action: it writes `18-CONFIRMATIONS.md` and exact envelope transitions only, never semantic profile/synthesis/handbook payloads.

## 8.4. Mandatory read sets

All modes load protocol constants/statuses, `0A`, `0D`, `0E`, `0F`, `0G` or `0H`, relevant source inventory/frontier/claim records, and the active snapshot metadata. Semantic decisions that a dependency or scope dimension is applicable or not applicable MUST carry the §0.1 classifier provenance and claim/evidence bindings. `Unknown` applicability is fail-inclusive: the dependency remains in the required read/validation closure and blocks mutation if it cannot be loaded and validated.

Strict single-scope modes additionally load the active POV file, at most one active ledger, the matching shared POV/ledger for persona scope, the persona's `UNMAPPED-DISCOVERY-BUFFER.md`, concrete source excerpts, and relevant normalized map. Before Discovery, the invocation MUST verify from a current `0H` global buffer summary or a deterministic read of every persona-local unmapped buffer that no `Pending-Merge` entry exists; an absent/stale summary without the full fallback read blocks the invocation. Load `0B` for auth; `0C` for sensitive data; dependency source excerpts for cross-boundary references; probe logs for matching probe tickets; persona siblings only for promotion uniqueness.

Acquisition loads `0A`, `0C`, `0G/0H`, supplied approved projections, lineage/sanitization reports, existing normalized maps, source inventory, and export reconciliation.

Sequential Reconciliation loads all relevant canonical `personas/_shared/` files, persona-local shared-staging/promotion-request/question/unmapped-discovery buffers, synthesis-gap buffers, and destination records. Profile Synchronization MUST load `0A`, the target `PRF` profile, and every record directly referenced by or applicable to any profile field, then recursively load the complete transitive `DERIVED-FROM` closure of those records. This closure explicitly includes every applicable `AUTH`, `CLM`, `SRC`, `CMP`, `ER`, `REL`, `BR`, `UC`, `IF`, `CFG`, `SCHED`, `DR`, `SM`, `STATE`, `DEP`, `MOD`, `CAP`, `DEC`, `SEC`, `NFR`, and `FLT` record, plus relevant source inventory, coverage, traceability, contradiction, decision-approval, and snapshot records needed to prove currency. Every loaded dependency MUST have its exact record ID, version where applicable, semantic/content fingerprint, source snapshot, and reciprocal binding verified. Profile Synchronization MUST reject mutation if any directly applicable record or any transitive dependency is absent, unresolved, stale, version-mismatched, fingerprint-mismatched, or read from a mixed snapshot. Confirmation validation of a PRF MUST load and validate this identical closure against the candidate-bound versions and fingerprints before attaching or accepting its envelope. Partial/Final Synthesis loads required closed or bounded forensic and assurance inputs plus existing synthesis blocks. Final Synthesis MUST load all active personas, all POV/ledger files, all relevant claims, reconciliation, coverage, contradictions, and gate inputs.

## 8.5. Cold resume check

Before mutation, the executor MUST return the semantic assessment fields below as part of its invocation result, and the conforming runtime MUST complete the deterministic validation fields before deciding whether to commit:

```markdown
- **Iteration / Invocation / Mode:** [...]
- **Persona / Cluster / Track / POV:** [...]
- **Relevant Active Tickets and Frontiers:** [...]
- **Last Mutation Affecting Targets:** [...]
- **Semantic Blockers and Contradictions:** [...]
- **Proposed Allowed Next Actions:** [...]
- **Explicit Context Exclusions and Reasons:** [...]
- **Runtime-Validated Loaded Files / Fingerprints:** [...]
- **Runtime Scope / Mode / Write-Right Check:** Pass | Fail
- **Runtime Stale Check:** Pass | Fail
- **Cold-Resume Result:** Pass | Fail
- **Cold-Resume Check Fingerprint:** [...]
```

The cold-resume check fingerprint is `sha256` over UTF-8 `COLD-RESUME|` plus system namespace, one `|`, invocation ID, one LF, and the §4.1.2 canonical serialization of every field above except its own carrier. The runtime MUST reproduce identity, scope, loaded-file fingerprints, stale state, and write rights from authoritative filesystem state; it validates but does not invent the executor's semantic blockers, actions, or exclusion reasons. A missing field, mismatched fingerprint, failed runtime check, or semantic exclusion that would omit required scope yields `Cold-Resume Result: Fail`, rejects semantic mutation, and still appends the required no-mutation `0G` entry with the check fingerprint and exact failure. The successful check fingerprint/result is bound into the same transaction's `0G` entry, making the pre-mutation authorization externally observable.

Filesystem state is authoritative over conversation. Contradictions in files create tickets/contradiction records; conversational memory cannot resolve them.

## 8.6. Concurrent persona traversal

Personas MAY run concurrently in one named iteration. They MUST NOT write canonical shared POV or shared question files directly.

- shared component discoveries go to `personas/{persona-directory}/SHARED-STAGING-BUFFER.md`;
- outbound shared questions go to `personas/{persona-directory}/SHARED-QUESTIONS-BUFFER.md`;
- unmapped coordinates go to `personas/{persona-directory}/UNMAPPED-DISCOVERY-BUFFER.md` with a same-transaction `FRT`;
- persona-owned invocation pointers may be written immediately;
- depromotion is annotated `[R-DEPROMOTION-PENDING]` and deferred.

Sequential Reconciliation collects, deterministically sorts and deduplicates all buffer classes, promotion requests, and synthesis gaps; merges compatible claims without improperly raising `BLOCK-CONFIDENCE-RANK` or losing profile detail; records contradictions; merges persona sets; commits shared records/tickets; merges unmapped discoveries into `0A`; creates required ticket/frontier work from `GAP` records; executes promotion and depromotion atomically under §§6.3–6.4; updates assurance links; and purges consumed buffers atomically. No later Discovery invocation may begin while an unmapped-discovery buffer contains `Pending-Merge`.

## 8.7. Stale checkpoint guard

Before writing, compare loaded checkpoint and source fingerprints with target record versions. If any target changed after the checkpoint, stop and request a refreshed read set. Changed serialized fingerprints invalidate affected normalized maps, source records, claims, coverage, and downstream blocks before further mutation.

## 8.8. `0G` invocation log

Append after every invocation:

```markdown
### [ITERATION] - [INVOCATION-ID] - [Scope] - [Mode]

- **Identity Header Summary:** [...]
- **Cold-Resume Result / Check Fingerprint:** [...]
- **Declared Semantic Actions / Runtime Authorization Result:** [...]
- **Loaded Files / Fingerprints:** [...]
- **Mutation Summary:** [...]
- **Tickets Created / Updated:** [...]
- **Frontiers Created / Closed:** [...]
- **Unmapped Discoveries / GAP Records Buffered, Merged, or Disposed:** [...]
- **Exports Acquired / Requested:** [...]
- **SRC / CMP / CLM IDs Updated:** [...]
- **Evidence Changes:** [...]
- **Promotion / Depromotion:** [...]
- **Coverage / Reconciliation Changes:** [...]
- **Stale Propagation:** [...]
- **Next Required Loads:** [...]
```

`0H` MAY compress current active states but MUST include its source `0G` range and fingerprint. When used to authorize a later Discovery invocation, it MUST also include a fingerprinted global count/list of every persona-local `Pending-Merge` unmapped discovery; otherwise the full-buffer fallback read in §8.4 is required.

---

# Part 9. Ticket FSM, Static Investigation, Probes, and Human Hatch

## 9.1. Ticket schema and canonical IDs

Ticket IDs use `MONIKER-PERSONA-POV-SEQUENCE`. Sequence allocation resets per iteration/persona/target POV and is reserved atomically. For buffered concurrent tickets, Sequential Reconciliation sorts unassigned candidates by target POV, normalized source coordinate, source persona prefix, and SHA-256 of the normalized description, then assigns the next free sequence. A directly authorized single-ledger invocation reserves the next free sequence before writing. The tuple and a stored creation fingerprint MUST be unique; an allocated ID is never reused, even when invalidated or superseded.

```markdown
### [TICKET-ID] [Context]

- **Source POV / File / Invocation:** [...]
- **Target POV / Ledger:** [...]
- **Status:** [T-*]
- **Trace History:** [actor:iteration transitions]
- **Affected SRC/CMP/CLM/Frontier IDs:** [...]
- **Memory Anchor and Creation Context:** [complete stateless assumptions]
- **Description:** [one technical ambiguity]
- **Static Investigation:** [bounded paths/methods attempted, evidence, and result]
- **Escalation Reason:** None | RUNTIME-REQUIRED | THIRD-PARTY-OPAQUE | UNAUTHORIZED-ACCESS | UNSAFE | DESTRUCTIVE | NO-SANDBOX | NO-AUTHORITY-TO-DECIDE | DOMAIN-MEANING-REQUIRED | EXTERNAL-FACT-UNAVAILABLE
- **Probe Feasibility / Inapplicability:** [assessment, proposed safe probe for RUNTIME-REQUIRED, or why a probe cannot safely/legally/meaningfully resolve the ambiguity]
- **Coverage / Promotion Effect:** [...]
```

## 9.2. FSM and write rights

```text
T-OPEN --target investigates--> T-PENDING
T-OPEN/T-FOLLOW-UP --target requires runtime--> T-PROBE-REQUIRED
T-PENDING --origin accepts--> T-COMPLETE
T-PENDING --origin rejects--> T-FOLLOW-UP -> T-OPEN investigation
Active state --valid proof--> T-INVALID or T-SUPERSEDED
Active state --mutual scope decision--> T-OUT-OF-SCOPE
T-PROBE-REQUIRED --target ingests successful probe--> T-PENDING
T-OPEN/T-FOLLOW-UP/T-PROBE-REQUIRED --Human Hatch after §9.3 prerequisites--> T-HUMAN-REQUIRED
```

Any POV may originate `[T-OPEN]` through the correct direct ledger or buffer. Only the target POV updates to `[T-PENDING]` or `[T-PROBE-REQUIRED]`. Only the origin accepts/rejects pending answers. Origin or target may propose out-of-scope; terminalization requires recorded mutual agreement or Human Hatch approval. A conforming runtime validates every transition, escalation reason, required evidence field, actor right, and source state; an executor cannot assign a terminal or probe state merely by writing its token.

**Direct creation of `[T-HUMAN-REQUIRED]` is forbidden.**

## 9.3. Human Hatch action prerequisites

### Human Hatch Action: Ticket Escalation

Only this subtype may terminalize an unresolved ticket as human-required. Before `Human Hatch: Ticket Escalation`:

1. perform and record bounded static investigation with exact methods, evidence, and remaining ambiguity;
2. select exactly one closed `Escalation Reason` and record the probe feasibility/inapplicability assessment;
3. for `RUNTIME-REQUIRED`, transition to `[T-PROBE-REQUIRED]`, propose a safe observable probe, and record the attempted probe or why it cannot be attempted safely, legally, or technically;
4. for every non-runtime reason, remain `[T-OPEN]` or `[T-FOLLOW-UP]` and record why runtime observation is inapplicable, unauthorized, unsafe, destructive, unavailable, or incapable of resolving the required domain meaning/external fact; a fake probe or transition through `[T-PROBE-REQUIRED]` is forbidden;
5. invoke the subtype with exact ticket, scope, denominator, risk, and coverage consequences plus reviewer authority.

Direct creation of `[T-HUMAN-REQUIRED]`, an escalation with `Escalation Reason: None`, or a reason inconsistent with the recorded feasibility assessment is invalid. Ticket Escalation may terminalize the ticket as `[T-HUMAN-REQUIRED]`, but the affected source remains `[C-GAP]` or `[C-PARTIAL]` unless an authorized human separately or concurrently approves an exact `[C-EXCLUDED]` decision under `Decision/Scope Approval`. No terminal ticket silently implies coverage. Human-required is an honest unresolved orchestration outcome distinct from exclusion: it blocks Exit A for any scope claiming the affected unit's closure, remains visible in partial synthesis and final risk reporting, and permits independent unaffected work to continue.

### Human Hatch Action: Decision/Scope Approval

This subtype does not require a fake static investigation or probe. It requires: exact proposed `DEC` content; observed `SRC`/`CLM`/ticket/contradiction inputs and immutable fingerprints; exact affected `COV` tuples and denominator arithmetic; alternatives, rationale, residual risks, rollout/rollback where applicable; reviewer identity, authority basis, and timestamp. If used to resolve a ticket, the Ticket Escalation prerequisites still apply to that ticket.

### Human Hatch Action: Confirmation

This subtype does not require a static investigation or probe. It requires: a passing Exit E Candidate Report; candidate payload manifest; exact typed record IDs including `PRF`/`HBK` where applicable; exact `SEMANTIC-RECORD-VERSION` and `SEMANTIC-CONTENT-FINGERPRINT` bindings; required approved `DECISION-CONTENT-VERSION`/`DECISION-CONTENT-FINGERPRINT` bindings and matching approval envelopes; review scope and limitations; reviewer identity and authority basis. For a PRF it MUST validate the complete §8.4 dependency closure; for other semantic records it MUST validate the complete applicable `DERIVED-FROM` closure. It may create `CNF` records and attach certification envelopes only; semantic payload edits are forbidden.

## 9.4. Probe specification and asynchronous handoff

```markdown
- **PROBE-ID:** [equal to TICKET-ID]
- **Target SRC/CMP:** [...]
- **Question / Expected Observable:** [...]
- **Exact Inputs / Preconditions:** [...]
- **Required Environment / Headers / Config:** [no secret values]
- **Safety Classification:** non-destructive | controlled mutation
- **Data Masking / Isolation Approval:** [...]
- **Disabled-Capability Authorization:** [required when applicable]
- **Logs / DB Diffs / Payloads to Capture:** [...]
- **Success / Failure Interpretation:** [...]
```

The operator executes in an isolated sandbox and writes `.extracted/probes/[PROBE-ID].log`. In a bound `Ticket Resolution` invocation, the target POV ingests the log, creates runtime claims, updates affected records, and moves the ticket to `[T-PENDING]`. The originating POV then accepts it as `[T-COMPLETE]` or rejects it as `[T-FOLLOW-UP]`; probe ingestion never bypasses origin validation. Probe output MUST be fingerprinted and scrubbed of secrets/PII before agent exposure.

## 9.5. No-mock integrity

When static analysis is insufficient and a sandbox is unavailable:

- do not fabricate paths, values, payloads, mocks, or specifications as evidence;
- keep claims `[E-INFERRED]` or `[E-UNKNOWN]`;
- keep affected components `[S-DISCOVERED]` and ineligible for promotion;
- use the Human Hatch path and preserve the coverage gap/partial state.

A test double used to exercise already-proven local logic MAY be a testing mechanism but MUST NOT be treated as observation of an opaque dependency.

## 9.6. No-mock fallback payload

```markdown
- **Suspected Behavior:** [...]
- **Evidence Claims / Anchors:** [...]
- **Missing Runtime Observation:** [...]
- **Static Methods Attempted:** [...]
- **Probe Feasibility / Attempt Result:** [...]
- **Proposed Human Verification:** [...]
- **Affected Promotion / Coverage / Exit Scope:** [...]
```

---

# Part 10. Evidence-Backed Sweep and Composite Deconstruction Closure

## 10.1. Evidence-backed `[R-SWEPT]`

A required traversal-matrix cell may be `[R-SWEPT]` only with a sweep record. Every required cell has its own sweep; claim-backed `Not-Applicable` cells have no sweep and remain visible in the denominator.

```markdown
### [SWEEP-ID] [Cluster / Track]

- **Snapshot / Environment:** [...]
- **Persona / Cluster / POV / Track:** [...]
- **Expanded Concrete Entry Coordinates:** [complete list; no wildcard-only closure]
- **Source Denominator:** [SRC IDs or inventory query + frozen result fingerprint]
- **Expected / Mapped / Excluded Counts by Kind:** [...]
- **Open / Terminal Frontier Counts:** [...]
- **Active / Terminal Ticket Counts:** [...]
- **Export Reconciliation:** [mapped/excluded/total and report fingerprint]
- **Coverage Result:** [...]
- **Closing Invocation:** [...]
- **Evidence Claims:** [CLM IDs proving enumeration and closure]
```

`R-SWEPT` is invalid if expected counts cannot be reproduced, any discovered boundary lacks exactly one `FRT`, any frontier is open, any normalized unit is unreconciled, any material contradiction is open, or a wildcard is the only denominator. Terminal-at-discovery frontiers are included in terminal counts and disposition evidence; omission is not an optimization.

## 10.2. Per-kind coverage arithmetic

For each kind `k` and each environment/snapshot/track scope:

```text
Gk = count(scope-specific rows in 16-SOURCE-COVERAGE where SRC kind=k and Scope in {In-Scope, Candidate, Approved-Excluded})
Xk = count(rows in Gk where Scope=Approved-Excluded and C Status=C-EXCLUDED with an approved exact DEC)
EffectiveDenominatork = Gk - Xk
Coveredk = count(rows in EffectiveDenominatork where C Status=C-COVERED)
Partialk = count(rows in EffectiveDenominatork where C Status=C-PARTIAL)
Gapk = count(rows in EffectiveDenominatork where C Status=C-GAP or Disposition in {Unmapped, Human-Blocked})
Candidatek = count(rows in EffectiveDenominatork where Scope=Candidate)

Valid partition: EffectiveDenominatork = Coveredk + Partialk + Gapk
Closed for kind k: Candidatek = 0 and Partialk = 0 and Gapk = 0 and Coveredk = EffectiveDenominatork
```

An excluded unit remains in the gross discovered population `Gk` and exclusion/risk reporting but is removed exactly once from the effective denominator. Candidate units cannot disappear; they must become In-Scope or Approved-Excluded. Coverage is computed separately for every applicable kind, track, snapshot, and environment, then aggregated. A global percentage MUST NOT hide a failed subgroup.

`[C-COVERED]` for an effective-denominator unit requires: a terminal in-scope source disposition; atomic mapping; material CLMs; no unresolved contradiction; required reciprocal references; terminal frontier; applicable export reconciliation; no active ticket affecting the unit; and, for a Disabled SRC/CMP/UC, exactly one canonical CAP membership with any overlap explicitly decided. An approved exclusion yields `[C-EXCLUDED]`, not `[C-COVERED]`, and is governed by the subtraction above.

## 10.3. Composite Exit A — Deconstruction Closure

Exit A succeeds only when all conditions hold for the pinned scope:

1. **Ticket closure:** zero `[T-OPEN]`, `[T-PENDING]`, `[T-FOLLOW-UP]`, or `[T-PROBE-REQUIRED]` entries.
2. **Acquisition closure:** every in-scope artifact is `[A-STRUCTURALLY-COMPLETE]`; requested or blocked references are resolved or exactly approved excluded.
3. **Registry closure:** no `[R-PENDING-HUMAN-REVIEW]`, `[R-PENDING-PERSONA-ASSIGNMENT]`, `[R-DEPROMOTION-PENDING]`, required `[R-UNSWEPT]`, or required `[R-IN-PROGRESS]` cell remains; every required traversal-matrix cell is `[R-SWEPT]` and every `Not-Applicable` cell is claim-backed.
4. **Source disposition:** every in-scope concrete source unit is terminally disposed and no Candidate remains.
5. **Frontier closure:** every traversal frontier is terminal.
6. **Export reconciliation:** every normalized logical unit is atomically mapped or approved excluded.
7. **Sweep validity:** every required traversal-matrix cell has a valid evidence-backed `[R-SWEPT]` record. Applicable POVs are exactly rows whose `Applicability` is `Required`; independent evaluation means each such row has its own valid sweep and MUST NOT inherit another row's result.
8. **Atomic coverage:** all effective denominator units are `[C-COVERED]`; no `[C-PARTIAL]` or `[C-GAP]` remains for Exit A scope.
9. **Forensic referential integrity:** ID uniqueness, existing forensic/assurance reciprocal links, no invalid tombstone references, no orphan claims, no dangling source/component/persona/entry references, and canonical ownership all validate. `Pending-Synthesis-Reference` placeholders are permitted until synthesis and do not count as dangling at Exit A.
10. **Contradiction closure:** no open material contradiction remains.
11. **Disabled equality and capability governance:** all in-scope disabled/shadow/conditional units meet the same static atomic closure rules; every effective-denominator Disabled `SRC`, `CMP`, and `UC` maps to exactly one canonical `CAP`, and every multiple/conflicting membership has the explicit relationship and approved decision required by §2.6.
12. **Human-required effect:** every `[T-HUMAN-REQUIRED]` is paired with either an exact approved exclusion or a scope that does not claim closure. It never silently counts as covered.

Exit A writes a deterministic report to `22-GATE-REPORTS.md` containing pass/fail per condition, counts by kind/track/environment, input fingerprints, validator version, and pinned snapshot.

## 10.4. Exit B — Stagnation

If a complete named iteration has zero valid mutations—no new atomic source/component/claim/frontier records, no ticket transition, no acquisition/reconciliation update, and no justified coverage change—the pipeline halts as stagnant and reports remaining state. Cosmetic edits do not count.

## 10.5. Exit C — Oscillation

If an actor-transition sequence of length 1–4 repeats for two consecutive completed cycles in one ticket trace, run one neutral Arbitration invocation, normally POV-5. If it fails to break the cycle, freeze the ticket and escalate through Human Hatch.

## 10.6. Exit D — Limit exhaustion and deterministic iteration accounting

Named iterations use exactly these canonical tokens wherever stored, compared, or consumed by an ID or state machine: `ALFA`, `BRAVO`, `CHARLIE`, `DELTA`, `ECHO`, `FOXTROT`, `GOLF`, `HOTEL`, `INDIA`, `JULIETT`, `KILO`, `LIMA`, `MIKE`, `NOVEMBER`, `OSCAR`, `PAPA`, `QUEBEC`, `ROMEO`, `SIERRA`, `TANGO`, `UNIFORM`, `VICTOR`, `WHISKEY`, `X-RAY`, `YANKEE`, `ZULU`. These 26 tokens are a total per-system budget across the complete protocol run, including every deconstruction pass reopened by synthesis gaps; the budget never resets after Exit A, synthesis, confirmation failure, or re-entry. Human-facing typography such as “X-ray” is display-only and MUST be normalized to `X-RAY` before persistence or comparison; it is not an accepted stored alias.

A named iteration begins with its first logged invocation and completes only when all invocations admitted to that iteration have terminal invocation log entries, all persona-local shared-staging, question, unmapped-discovery, and synthesis-gap buffers produced in it have been reconciled, and the iteration checkpoint records mutation count, remaining tickets/frontiers/candidates, and the next action. Parallel invocations share the same current token and do not consume extra tokens. A synthesis `GAP` discovered while an iteration is open is reconciled within that iteration when possible; if that iteration is already complete, reopening consumes the next unused token. Re-running synthesis after closure continues under the current open token or next unused token and MUST NOT return to `ALFA` or reuse a completed token. The system MUST NOT proceed beyond `ZULU`. Unclosed work at completion of `ZULU` halts and requires human intervention.

---

# Part 11. Structured Synthesis

## 11.1. Partial versus Final Synthesis

**Partial-Synthesis** runs before Exit A only for an explicitly bounded persona/cluster/track/snapshot. Blocks remain `[B-DRAFT]`, are labeled provisional, and state the incomplete denominator. They MUST NOT be reconstruction inputs.

**Final-Synthesis** runs only after Composite Exit A passes. It writes complete structured catalogs from closed facts. Newly discovered gaps go to `SYNTHESIS-GAP-BUFFER.md`; Sequential Reconciliation opens tickets, creates or updates missing SRC/frontier/coverage records, immediately marks the pinned Exit A report invalid, and marks all dependent B blocks stale. Deconstruction is re-armed and confirmation/handoff is prohibited until Exit A passes again on a new gate-input fingerprint. Final synthesis cannot directly mutate ledgers.

After IDs exist, **Cross-Reference-Reconciliation** may update only fields such as `Bound Use-Cases`, `Rules`, `Entities`, `Interfaces`, and reciprocal references. It MUST NOT change observed behavior, decisions, or evidence; because these references are semantic payload, a changed reference increments the semantic record version/fingerprint and invalidates any certification envelope.

### Synthesis gap buffer

Each `SYNTHESIS-GAP-BUFFER.md` record is normative:

```markdown
### [GAP-ID]

- **Source Synthesis Invocation / Block ID:** [...]
- **Observed Missing Requirement:** [one falsifiable missing semantic or evidence requirement]
- **Affected IDs:** [SRC/CMP/CLM/synthesis IDs]
- **Source Fingerprints / Snapshot:** [...]
- **Proposed Target Persona / Cluster / Track / POV / Ledger:** [...]
- **DEDUP-KEY:** [sha256(system namespace | normalized missing requirement | sorted affected IDs | snapshot)]
- **Materiality:** Material | Supporting
- **Disposition:** Pending | Ticket-Created | Frontier-Created | Both-Created | Merged | Invalid | Approved-Excluded
- **Resulting Ticket / FRT / DEC IDs:** [...]
```

Partial/Final Synthesis writes only `Pending` GAP records and cannot create tickets or frontiers. Sequential Reconciliation sorts by `DEDUP-KEY`, then `GAP-ID`; merges equal keys while retaining every source invocation/block/fingerprint; records a contradiction rather than merging incompatible target proposals; and deterministically creates the minimum required work. A missing cross-POV answer creates one canonical ticket in the proposed ledger; a newly discovered traversal boundary creates one `FRT`; a gap requiring both creates both and cross-links them. Existing equivalent active work is linked rather than duplicated. Sequential Reconciliation updates disposition atomically, invalidates the pinned Exit A report for material gaps, applies targeted staleness, and logs iteration-budget effects under §10.6.

## 11.2. Common synthesis block header

Every atomic synthesis block MUST contain a semantic payload followed by a separately delimited certification envelope:

```markdown
- **SEMANTIC-RECORD-VERSION:** [positive integer]
- **SEMANTIC-CONTENT-FINGERPRINT:** [sha256 per §4.1.1]
- **COVERAGE-STATUS:** [C-COVERED] | [C-PARTIAL] | [C-GAP] | [C-EXCLUDED]
- **DERIVED-FROM:** [SRC/CMP/CLM/INV/etc. IDs with versions and fingerprints]
- **Material Claims:** [CLM IDs]
- **BLOCK-CONFIDENCE-RANK:** None | Low | Medium | High
- **BLOCK-EVIDENCE-PROFILE:** [counts, claim IDs, and route/lineage IDs by each E level]
- **Supersession:** [...]

<!-- CERTIFICATION-ENVELOPE: excluded from semantic version/fingerprint -->

- **BLUEPRINT-STATUS:** [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
- **Confirmation:** [CNF-ID required only for B-CONFIRMED]
- **Confirmation Envelope Version:** [positive integer]
- **Confirmation Envelope Audit:** [...]
<!-- END-CERTIFICATION-ENVELOPE -->
```

Coverage is semantic payload; blueprint status is envelope metadata. `[B-CONFIRMED]` requires `[C-COVERED]`, no stale dependency, a passing candidate binding, and a valid human confirmation of the exact semantic version/fingerprint. Attaching confirmation MUST NOT change semantic bytes. Empty domains require inventory-backed zero claims and may then be confirmed as empty.

## 11.3. Architecture blueprint — `90-ARCH-BLUEPRINT.md`

Contains system boundary and context diagrams; actors/personas and external systems; observed legacy technology separated from target decisions; atomic `MOD` records; inter-module contracts and error propagation; deployment topology; architectural invariants and modernization consequences; capability navigation; and a prominent active-versus-disabled architecture view. Purely navigational module indexes are grouping containers under §4.2 and list `MOD` IDs only.

```markdown
### [MOD-ID] [Module or Bounded Context]

[Common semantic payload and certification envelope]

- **Responsibility / Domain Boundary:** [...]
- **Public Interfaces:** [IF IDs]
- **Encapsulated State / Entities:** [ER/STATE IDs]
- **Dependencies / Dependents:** [MOD/DEP/IF IDs]
- **Data and Events Passed:** [schemas, ER/IF/event IDs]
- **Invariants / Rules Owned:** [INV/BR IDs]
- **Personas / Use Cases / Capabilities:** [persona/UC/CAP IDs]
- **Failure Propagation:** [FLT IDs]
- **Observed Legacy Placement:** [...]
- **Target Decision:** [DEC ID or None]
```

Every semantic architecture statement about responsibility, public interface, encapsulation, dependency, or data passing MUST belong to an atomic `MOD` and have material CLMs. `MOD <-> IF/ER/DEP/BR/UC/CAP/FLT` reciprocal references are required where applicable.

## 11.4. Entity — `91-DATA-MODEL.md`

```markdown
### [ER-ID] [Entity]

[Common header]

- **Domain Context / Purpose:** [...]
- **Storage:** table | view | collection | cache-key-space | file | in-memory
- **Physical SRC IDs:** [...]
- **Primary Key / Secondary Keys / Indexes:** [...]
- **Columns / Fields:**
  | Name | Type | Null | Default | Constraints | Unique | Claims |
- **Relationships:** [REL IDs]
- **Enums / Lookup Domains:** [...]
- **Lifecycle State Fields:** [field -> SM IDs]
- **Derived / Computed Fields:** [formula, dependencies, read/write time]
- **Invariants:** [BR/INV IDs]
- **Audit / Soft Delete:** [...]
- **Retention / Lifecycle:** [...]
- **Permissions:** [persona read/write/delete matrix]
- **Capability State Impact:** [active/disabled usage]
- **Schema Drift Matrix:** [environment/snapshot -> differences + claims]
```

## 11.5. Relationship — `91-DATA-MODEL.md`

```markdown
### [REL-ID] [EntityA.field -> EntityB.field]

[Common header]

- **Cardinality:** 1:1 | 1:N | M:N
- **Enforcement:** DB-FK | application | mixed | none
- **Delete / Update Semantics:** cascade | restrict | nullify | custom
- **Inverse Access:** [...]
- **Business Constraint:** [BR IDs]
- **Reciprocal Entity References:** [ER IDs]
```

## 11.6. State machine — `91-DATA-MODEL.md`

A state machine is required only for an **eligible lifecycle scope**: an entity aggregate or capability with a finite, behaviorally meaningful lifecycle, not every boolean or presentation state.

```markdown
### [SM-ID] [Lifecycle]

[Common header]

- **Scope / Owner Entity or Aggregate:** [ER IDs]
- **State Field(s):** [...]
- **States / Meaning:** [...]
- **Initial / Terminal States:** [...]
- **Transitions:**
  | From | To | Event/UC | Guard BR | Effects | Claims |
- **Invalid / Recovery Transitions:** [...]
- **Environment / Capability Variants:** [...]
```

Eligibility decisions use the §12.1 lifecycle table so absence cannot be assumed. `Unknown` requires a ticket, is treated as eligible for closure, and blocks Exit E until resolved to `Eligible` with a complete `SM` or `Ineligible` with claim-backed provenance.

## 11.7. Database routine — `91-DATA-MODEL.md`

```markdown
### [DR-ID] [Routine]

[Common header]

- **Kind:** trigger | function | procedure | rule | generated-column | materialized-refresh | constraint-helper
- **Owner Table / Entity:** [SRC/ER IDs]
- **Firing Event / Timing / Granularity:** [...]
- **Atomic Body Intent:** [...]
- **Inputs / Outputs:** [...]
- **Entities Read / Written:** [...]
- **Invariants / Rules Enforced:** [INV/BR IDs]
- **Callers / Firers:** [CMP IDs]
- **Failure / Transaction Semantics:** [...]
- **Canonical DB Source / Drift:** [...]
- **Target Decision:** Retain | Lift | Redesign | Retire | Undecided
- **Target Decision DEC:** [DEC ID or None]
```

A production DB routine absent from the canonical DB source creates a material drift contradiction and follows the ticket/Human Hatch path.

## 11.8. Business rule — `92-BUSINESS-RULES.md`

```markdown
### [BR-ID] [Rule]

[Common header]

- **Authoritative Rule Owner:** [one domain/module owner]
- **Defense-in-Depth Mirrors:** [optional consistent enforcement locations]
- **Applies When:** [UC/event/state]
- **Condition:** [deterministic, framework-agnostic predicate with units]
- **Action / Effect:** [...]
- **Precedence / Conflict:** [related BR IDs]
- **Inputs Read / Writable Effects:** [ER fields/events]
- **Violation Policy:** reject | warn | clamp | default | compensate | rollback
- **Boundary Cases / Null / Locale / Overflow:** [...]
- **Idempotency:** [...]
- **Capability State:** Active | Disabled
- **Target Mapping Decision:** [target authoritative owner and optional mirrors]
```

One authoritative rule owner is REQUIRED. Multiple physical enforcement points are allowed only as consistent defense-in-depth mirrors, not competing authorities.

## 11.9. Use case — `93-USE-CASES.md`

```markdown
### [UC-ID] [Use Case]

[Common header]

- **Primary Persona / Secondary Actors:** [...]
- **Goal:** [...]
- **Trigger / Entry:** [IF/SCHED/DR IDs]
- **Preconditions / Postconditions:** [ER/CFG/INV IDs]
- **Invariants Relied Upon:** [INV/BR IDs]
- **Capability Classification:** Active | Conditional | Shadow | Disabled | Retired
- **Observed Legacy Behavior:** [...]

#### Main Flow

| Step | Action | Rules | Data R/W | Interface | Claims |

#### Alternate Flows

[branch, divergence, path, rejoin]

#### Failure and Recovery

[failure, detection, compensation, user/system outcome]

- **Rules / Entities / Interfaces / Side Effects:** [reciprocal ID lists]
- **Behavioral Contract:**
  - Given [pre-state]
  - When [trigger]
  - Then [observable outputs and post-state]
- **Target Decision:** Preserve | Change-By-Approved-Migration | Redesign | Retire
- **Target Decision DEC:** [DEC ID]
```

Disabled use cases remain visible, are not presented as current normal behavior, and require an explicit target decision.

## 11.10. Interface — `94-INTERFACES.md`

```markdown
### [IF-ID] [Interface]

[Common header]

- **Channel / Transport / Direction:** [...]
- **Actor / Persona / AUTH:** [...]
- **Protocol / Content Type:** [...]
- **Request Schema:**
  | Field | Type | Required | Default | Validation BR | Claim |
- **Response Schema:**
  | Field | Type | Notes | Claim |
- **Error Semantics / Retryability:** [...]
- **Guarantees:** idempotency, transactionality, delivery, ordering, timeout
- **AuthN / AuthZ:** [...]
- **Versioning / Compatibility:** [...]
- **Rate Limits / Backpressure:** [...]
- **Bound UC / ER IDs:** [...]
- **Capability State / Environment Matrix:** [...]
- **Compatibility Baseline:** [externally relied-upon semantics]
- **Target Migration Decision:** [versioned approval required for breaking change]
```

External semantics are compatibility baselines, not absolutely frozen. Breaking changes require an approved, versioned migration decision, consumer impact, rollout, and rollback plan.

## 11.11. Deployment, configuration, and scheduling — `95-DEPLOYMENT.md`

```markdown
### [DEP-ID] [Deployment Element]

[Common header]

- **Component / Placement / Replication:** [...]
- **Connections:** [IF/ER/DEP IDs]
- **Single Point of Failure:** yes | no | unknown
- **Environment / Snapshot Matrix:** [...]

### [CFG-ID] [Configuration or Feature Flag]

[Common header]

- **Medium:** env | DB | flag service | file | build-time
- **Exact Key:** [...]
- **Type / Default / Allowed Values:** [...]
- **Source / Precedence / Setter:** [...]
- **Behavioral Effect:** [UC/BR/IF/capability IDs]
- **Enable/Disable Semantics and Environment Matrix:** [...]
- **App-side / DB-side Materialization:** [...]

### [SCHED-ID] [Scheduled or Async Job]

[Common header]

- **Type / Trigger / Frequency / Window:** [...]
- **Owner Persona:** [...]
- **Actions:** [UC/BR/ER IDs]
- **Failure / Retry / Dead-Letter / Idempotency:** [...]
- **Dependencies / Locks:** [...]
- **Capability State:** Active | Conditional | Shadow | Disabled | Retired | Unknown
```

Every configuration flag requires a known effect or a gap. Secret values are never recorded.

## 11.12. NFR and security — `96-NON-FUNCTIONAL-SECURITY.md`

```markdown
### [NFR-ID] [Attribute]

[Common header]

- **Attribute:** latency | throughput | availability | RPO | RTO | scalability | operability | other
- **Observed / Inferred Value:** [...]
- **Scope / Environment:** [...]
- **Driver / Constraint:** [DEP/IF/ER/CLM IDs]
- **Target Decision / SLO:** [separate from observation]

### [SEC-ID] [Finding]

[Common header]

- **Category:** auth | authorization | transport | at-rest | in-transit | PII | secrets | audit | threat | privacy
- **Finding / Data Classification:** [...]
- **Risk Level:** high | medium | low | unassessed
- **Affected IDs / Environments:** [...]
- **Observed Safeguard / Gap:** [...]
- **Human Review / Target Mitigation:** [...]

### [FLT-ID] [Fault or Resilience Behavior]

[Common header]

- **Failure Class / Trigger / Detection Point:** [...]
- **Propagation Path:** [UC/IF/SCHED/DEP/CMP IDs]
- **Retry / Backoff / Timeout / Circuit Behavior:** [...]
- **Rollback / Compensation / Transaction Boundary:** [...]
- **Dead-Letter / Quarantine / Operator Action:** [...]
- **User/System Outcome / Observability:** [...]
- **Environment / Capability Variants:** [...]
- **Related Rules / Claims:** [BR/CLM IDs]
```

Every relevant `0C` finding MUST have a reciprocal `SEC` mirror. Neither file declares legal compliance. Every distinct material failure class discovered in forensic records MUST map to an `FLT`; `96` MUST include a system-level failure taxonomy and retry/rollback/compensation/dead-letter topology, with an evidence-backed zero statement for any genuinely absent mechanism.

## 11.13. Persona profile

Profile Synchronization writes one profile and its `0A` row atomically:

```markdown
## [PRF-ID] Persona: [Name]

- **PRF-ID:** [same PRF-ID derived from system namespace + persona prefix]
- **Prefix / Actor Class:** [...]
- **Canonical and Shared Entries:** [SRC/IF/cluster IDs]
- **Authentication and Session:** [AUTH refs, lifecycle]
- **Authorization Matrix:**
  | Resource / Capability | Allowed | Enforcement | Claims |
- **Visible Data Scope:** [ER read/write/delete]
- **Key Active Workflows:** [UC IDs]
- **Disabled / Historical Workflows:** [UC/capability IDs]
- **Triggers / Events / Schedules:** [...]
- **UI Surfaces and State:** [...]
- **Failure / Edge Behavior:** [...]
- **Dependencies / Couplings:** [...]
- **Observed Legacy Facts:** [CLM IDs]
- **Target Decisions:** [decision IDs]
- **SEMANTIC-RECORD-VERSION:** [positive integer; included in semantic hash preimage]
- **SEMANTIC-CONTENT-FINGERPRINT:** [sha256 per §4.1.1; this carrier field alone is excluded from its own hash preimage]
- **DERIVED-FROM:** [dependency IDs with exact versions/fingerprints]
- **Material Claims:** [CLM IDs]
- **BLOCK-CONFIDENCE-RANK:** None | Low | Medium | High
- **BLOCK-EVIDENCE-PROFILE:** [counts, claim IDs, and route/lineage IDs by each E level]

<!-- CERTIFICATION-ENVELOPE: confirmation writers may update only this region -->

- **BLUEPRINT-STATUS:** [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
- **Confirmation:** [CNF-ID required for B-CONFIRMED]
- **Confirmation Envelope Version:** [positive integer]
- **Confirmation Envelope Audit:** [...]
<!-- END-CERTIFICATION-ENVELOPE -->
```

Profile Synchronization is the only writer of profile semantic payload and synchronizes its mandatory `PRF-ID` atomically with the matching `0A` persona-prefix row. Authorized `Human Hatch: Confirmation` or `Reconstruction-Handoff` confirmation actions may update only the delimited certification envelope after loading and validating the complete §8.4 direct-and-transitive dependency closure and verifying the exact candidate-bound PRF ID, semantic version, and semantic fingerprint; they MUST NOT edit profile semantic content/version. A profile without a registry row, a registry row without a profile, a duplicate PRF for one prefix, or a PRF whose canonical prefix coordinate disagrees with `0A` is invalid after Profile Synchronization begins.

---

# Part 12. Traceability, Referential Integrity, and Staleness

## 12.1. Traceability — `20-TRACEABILITY.md`

One atomic row per source-to-semantic mapping:

```markdown
| SRC | Persona/Owner/PRF | Track/CAP | POV/CMP | MOD | UC | BR | ER/REL/SM/DR | IF | DEP/CFG/SCHED | PRV/SEC/NFR/FLT | HBK | CLM/DEC/CNF | C Status |
```

Multiple rows MAY represent multiple mappings, but source coverage is evaluated once per `SRC` and requires all applicable mappings. The ID registry, source denominator query fingerprints, lifecycle eligibility decisions, exclusions, and reciprocal reference check results are included in this file.

Lifecycle eligibility uses one row per candidate entity aggregate or capability:

```markdown
| Scope ID | Classification | Reason | CLM / Evidence Bindings | Classified By / Invocation | Snapshot |
| :------- | :------------- | :----- | :---------------------- | :------------------------- | :------- |
```

`Classification` is `Eligible | Ineligible | Unknown`. It is a semantic predicate: the classifier MUST record provenance and claim/evidence bindings. `Unknown` is fail-inclusive and is treated as eligible for completeness and gate purposes until resolved; it requires a ticket and cannot justify omission of an `SM`. The conforming runtime deterministically verifies row uniqueness, snapshot consistency, required provenance, and the corresponding `SM` presence or explicit ineligibility result.

Higher-level traceability invariants are:

1. every source unit in the inventory has at least one traceability row, including exact exclusions;
2. every `ER` is read or written by at least one `UC`, or is explicitly classified as structural/reference-only with a claim-backed reason;
3. every `IF` is consumed or produced by at least one `UC`, or is exactly excluded/retired;
4. every `UC` step cites at least one `BR`, entity read/write, interface effect, state transition, or explicit claim-backed no-op;
5. every `SCHED` and `DR` maps to at least one `UC`, `BR`, or invariant;
6. every `CFG` maps to a behavioral effect or exact approved exclusion;
7. every semantic `MOD` has reciprocal links to each applicable `IF`, `ER`, `DEP`, `BR`, `UC`, `CAP`, and `FLT`, while module indexes carry only member links;
8. every effective-denominator Disabled `SRC`, `CMP`, and `UC` maps to exactly one canonical `CAP`, or is exactly approved excluded; any overlap has an explicit typed relationship and approved `DEC` with one target-decision owner;
9. every synthesis ID is reachable from a source/claim lineage and from at least one handbook navigation path when confirmed.

## 12.2. Required reciprocal references

Validators enforce at least:

- `ER <-> REL`; `ER state field <-> SM` when eligible;
- `DR caller/firer <-> CMP DB-ROUTINES-FIRED`;
- `BR authoritative owner <-> MOD`; `BR <-> UC steps`; `BR <-> ER effects`;
- `MOD <-> IF/ER/DEP/BR/UC/CAP/FLT` where applicable; purely navigational module indexes <-> member MOD IDs only;
- `UC <-> persona`, `UC <-> IF`, `UC <-> ER/BR/SCHED/DR`;
- `IF <-> UC` and interface entity mappings;
- `CFG <-> affected UC/BR/IF/CAP`;
- `SCHED <-> UC/ER/persona`;
- `PRV <-> SEC`; `NFR <-> driver`; `FLT <-> UC/IF/SCHED/DEP/BR`;
- `AUTH <-> persona/IF/enforcement point`; `SHR <-> participating persona pointers`; `CAP <-> canonical member IDs/CFG/DEC`, with exact-one disabled membership and explicit overlap relations;
- every `COV <->` its unique immutable SRC-version/persona/track/environment/snapshot tuple and evidence; every `CND <->` its helper candidate/action identity and resulting disposition IDs;
- every `[B-CONFIRMED]` block, `PRF`, and `HBK` `<-> valid CNF` with exact candidate-bound typed ID/version/fingerprint; every exclusion/target choice `<-> DEC`;
- every persona-prefix row in `0A <->` exactly one `PRF`, and every `PRF <->` that row plus its complete §8.4 dependency closure;
- every `HBK <->` one normalized handbook path/stable anchor, its upstream synthesis IDs/versions/fingerprints, candidate-manifest entry, and at least one handbook navigation path;
- `CMP <-> SRC`, claims, persona references, and dependencies;
- normalized coordinate <-> SRC <-> atomic CMP/edge or exclusion;
- handbook `HBK` section <-> normalized relative path/stable anchor <-> synthesis IDs and exact versions/fingerprints.

No dangling, orphan, duplicate-owner, invalid-type, or unapproved tombstone reference is allowed at Exit A or E.

## 12.3. Synthesis-ID timing

Discovery records MAY contain `Pending-Synthesis-Reference` placeholders, never invented IDs. Partial/Final Synthesis allocates deterministic semantic IDs from canonical keys. Cross-Reference-Reconciliation then replaces placeholders and adds reciprocal references without inventing new legacy facts; because reference bindings are part of semantic payload, it increments the affected semantic version/fingerprint and invalidates prior envelopes. Pending placeholders block confirmation and Exit E, but do not block early forensic discovery.

## 12.4. Dependency and impact provenance

Each synthesis/handbook block maintains `DERIVED-FROM` edges to exact record versions and fingerprints. When an upstream source, claim, mapping, decision, or exclusion changes:

1. compute reverse dependency closure;
2. mark only affected synthesis blocks and handbook sections `[B-STALE]`;
3. log cause and impact in `0G` and the block;
4. re-run the relevant bounded synthesis and validation;
5. preserve unaffected confirmed blocks.

Broad whole-file staleness is allowed only when dependency granularity is unavailable; that limitation is itself reported.

## 12.5. Deterministic validation summary

Every validator run records:

- schema/protocol/check-registry/validator versions;
- authoritative input file fingerprints and a deterministic input-set fingerprint;
- checks executed and stable check IDs;
- for every check and blocker row, the complete sorted evidence-binding tuples and recomputed §4.1.2 evidence-set fingerprint;
- pass/fail counts and exact offending IDs;
- coverage arithmetic by kind/track/environment;
- broken links, invalid diagrams, duplicate IDs, orphan/dangling references, stale blocks, unsupported statuses, and unresolved contradictions;
- any authorized human review input used by a check, including reviewer identity/authority, review scope, timestamp, exact candidate-bound fingerprints, limitations, and result;
- output fingerprint and timestamp.

Evidence bindings are externally observable in the validation summary using this exact table shape, with one or more rows per Check/Blocker ID as required:

```markdown
- **Evidence Bindings:**
  | Row Kind | Check / Blocker ID | Normalized Authoritative Input Path | Binding Kind | Typed Record / Artifact ID or None | Version or None | Authoritative Fingerprint |
  | :------- | :----------------- | :---------------------------------- | :----------- | :--------------------------------- | :-------------- | :------------------------ |
```

Rows sort by `(Row Kind, Check / Blocker ID, Normalized Authoritative Input Path, Binding Kind, Typed Record / Artifact ID or None, Version or None, Authoritative Fingerprint)`. `Binding Kind` uses the closed §4.1.2 mapping. A changed input invalidates the prior summary. Deterministic values supplied by an executor MUST be independently reproduced. A missing required check, unknown or duplicate Check ID, incomplete evidence-binding set, or mismatch against the applicable §15.1.1 registry makes the summary and dependent gate result invalid.

---

# Part 13. Human Reconstruction Handbook

Every semantic handbook section consists of a complete semantic payload with a mandatory typed HBK heading followed by a compact machine-validated certification envelope:

```markdown
### [HBK-ID] [Handbook Section Title]

- **HBK-ID:** [HBK-ID derived from system namespace + HANDBOOK-RELATIVE-PATH + STABLE-SECTION-ANCHOR]
- **HANDBOOK-RELATIVE-PATH:** [normalized under §4.1 and always relative to the fixed `.extracted/handbook/` root; the root itself is omitted]
- **STABLE-SECTION-ANCHOR:** [explicit opaque ASCII lowercase token matching `[a-z0-9]+(?:-[a-z0-9]+)*`]
- **SEMANTIC-RECORD-VERSION:** [positive integer]
- **SEMANTIC-CONTENT-FINGERPRINT:** [sha256 per §4.1.1]
- **UPSTREAM:** [synthesis IDs, exact semantic versions, semantic fingerprints]
- **Payload Content:** [all semantic prose, diagrams, examples, and navigation assertions]

<!-- CERTIFICATION-ENVELOPE: excluded from semantic version/fingerprint -->

- **HANDBOOK-STATUS:** [B-DRAFT] | [B-POPULATED] | [B-CONFIRMED] | [B-STALE]
- **Confirmation:** [CNF-ID required for B-CONFIRMED]
- **Confirmation Envelope Version:** [positive integer]
- **Confirmation Envelope Audit:** [...]
<!-- END-CERTIFICATION-ENVELOPE -->
```

The handbook payload MUST be fully generated, link/diagram validated, `[B-POPULATED]`, assigned globally unique path-and-anchor-derived `HBK` IDs, and included by exact HBK ID/version/fingerprint in the Exit E candidate payload manifest before human review. Validators MUST read HBK identity coordinates only from `HANDBOOK-RELATIVE-PATH` and `STABLE-SECTION-ANCHOR`; the Markdown heading/title is display content and MUST NOT be parsed as an anchor or identity coordinate. The path MUST satisfy §4.1 normalization and the anchor MUST satisfy the assigned opaque-token grammar. A heading/title edit preserves identity when the semantic identity, path, and assigned anchor remain unchanged. A proven path move retains identity only under the existing rename/alias rule and records the prior coordinate; an unproven move or semantic split creates new HBK IDs, with splits tombstoning or superseding the old record. Confirmation attaches only the envelope and its `CNF` MUST bind that exact HBK identity; it MUST NOT regenerate, rewrite, reflow, or otherwise change payload bytes. Navigation-only prose may be part of a chapter payload without independent claims, but any post-candidate edit changes the semantic fingerprint and restarts certification.

## 13.1. Audience paths and progressive disclosure

`START-HERE.md` MUST provide short paths for:

- product/domain owner: overview -> personas -> use cases -> rules -> decisions/gaps;
- architect/developer: architecture -> data -> interfaces -> background -> deployment -> decisions;
- security/privacy reviewer: personas/auth -> security/privacy -> data -> interfaces -> risks;
- data engineer: domain/data -> DB routines -> retention -> migration decisions;
- QA: use cases -> equivalence/acceptance -> rules -> known gaps;
- operator/SRE: deployment/config -> background -> NFR -> security -> run risks.

Each chapter begins with prose, diagrams, and representative examples before dense references. IDs appear as hyperlinks/reference anchors, not the primary user interface. Evidence rank/profile is visible but unobtrusive: e.g. a short badge or note linking to evidence, not repeated raw claim tables in the main narrative.

## 13.2. Required chapter behavior

Every chapter MUST:

- distinguish **Observed Legacy** from **Target Decision**;
- separate Normal behavior from Shadow/Conditional and Disabled behavior;
- link to upstream synthesis records and relevant source/evidence on demand;
- include navigation to previous/next, START-HERE, glossary, decisions, and gaps;
- avoid unexplained acronyms and identifier-only prose;
- show examples that are evidence-backed or explicitly labeled illustrative;
- pass link and diagram validation.

`08-BACKGROUND-PROCESSING.md` MUST include the consolidated fault/resilience view: failure taxonomy, propagation, retries/backoff, transaction rollback, compensation, dead-letter/quarantine paths, observability, and operator/user outcomes, cross-linked to `FLT` records.

## 13.3. Disabled and Dormant Features chapter

`12-DISABLED-AND-DORMANT-CAPABILITIES.md` is mandatory and prominent. It provides an index and one human-oriented section per canonical `CAP` containing exact disabled member IDs, status/environment matrix, mechanism, evidence rank/profile, dependencies, historical/intended personas and behavior, safety/risk, related architecture/data/rules/use-cases/interfaces, and target decision (`Restore`, `Preserve Dormant`, `Redesign`, `Retire`). It MUST demonstrate that every effective-denominator Disabled `SRC`, `CMP`, and `UC` appears under exactly one canonical CAP or an exact approved exclusion, and disclose any approved overlap relationship/decision without duplicating decision ownership.

Disabled features MUST also be cross-linked from every affected chapter. They MUST NOT be blended into current active behavior or omitted because they are not executable in production.

## 13.4. Handbook quality gate

Before candidacy, the conforming runtime validates:

- all required chapters exist and each contains at least one claim-backed semantic handbook assertion or an evidence-backed zero-domain statement; headings, navigation links, placeholders, or grouping prose alone are vacuous;
- all internal links and anchor targets resolve;
- every Mermaid/diagram block parses with a conforming deterministic parser; unavailable parsing capability is `unsupported` and blocks candidacy rather than silently falling back to invisible self-assessment;
- each candidate synthesis block is discoverable from at least one handbook path;
- active versus disabled presentation is unambiguous;
- target decisions are not presented as observed facts;
- no `[B-STALE]`, provisional, or pending-reference section remains.

After candidacy, authorized human readability review is recorded through one or more §5.7 `CNF` records whose combined exact HBK bindings cover every candidate handbook section and whose `Handbook Readability Review` is `Pass`; their hashed review scope/limitations state how core user goals were followed without reading raw catalogs. A CNF that does not review HBK content uses `Not-Applicable`; `Fail` or incomplete HBK coverage blocks content readiness. The Exit E Content-Readiness Report MUST validate reviewer authority, candidate manifest/report bindings, complete HBK coverage, timestamps, results, and unchanged handbook fingerprints. Human review cannot substitute for deterministic link, diagram, or non-vacuity checks, and an unrecorded review has no gate effect.

---

# Part 14. Decisions, Reconstruction, Equivalence, and Compatibility

## 14.1. Decision log

`15-DECISIONS.md` is the authoritative decision register from Preflight onward. Reconstruction-Handoff produces a human-readable decision-log projection containing every applicable approved record and its source link; it MUST NOT fork or reinterpret the decision. The canonical schema is §5.5. Target decisions never alter observed facts. Superseded decisions remain tombstoned and linked.

## 14.2. Modernization mapping templates

The following are decision prompts, not mandates:

- global/out-of-band state -> request scope, explicit service state, or managed store;
- mixed server/client presentation ownership -> deliberate SSR/SPA/hybrid ownership with synchronization rules;
- scattered auth checks -> centralized policy boundary with documented defense-in-depth checks;
- DB-autonomous rules -> intentionally retain, lift, redesign, or retire per `DR` and invariant;
- synchronous cross-availability side effects -> evaluate durable commands/events;
- string-composed persistence -> evaluate a controlled data-access boundary;
- filesystem side effects -> evaluate durable storage retaining lifecycle semantics;
- application-only relationships -> enforce physically or document compensating consistency;
- derived fields -> retain, materialize, recompute, or remove with a decision;
- soft-delete/audit mechanics -> preserve the requirement, not necessarily legacy columns.

The target architecture owner chooses the result and records it.

## 14.3. Business-rule ownership

Every `BR` maps to exactly one authoritative target owner. Optional validation, gateway, UI, or DB mirrors are permitted only if they remain semantically consistent and cite the authoritative rule. Conflicting mirrors are prohibited.

## 14.4. Interface compatibility

External interface schemas, errors, ordering, idempotency, delivery, and timing guarantees form a compatibility baseline. They MAY change only through an approved versioned migration decision identifying consumers, compatibility period, data migration, rollout, observability, and rollback. Internal transport MAY change freely if externally observable and use-case contracts remain satisfied or have approved changes.

## 14.5. Equivalence and acceptance suite

Generate:

- example acceptance tests from every UC Given/When/Then contract;
- property assertions from invariants and BR conditions;
- state-transition tests from every eligible SM;
- contract tests for external IF compatibility baselines;
- disabled-capability tests that verify intended dormant state and prevent accidental enablement, plus restoration tests only when target decision is Restore;
- migration tests for approved behavior changes.

Generated tests are reviewable derivatives. Their expected semantics come from B-confirmed synthesis and approved decisions.

## 14.6. Sole-input invariant

Legacy source may be treated as unavailable during reconstruction only when the bundle:

- has a validated signed outer bundle manifest whose hashed payload contains `EXIT-E-STATUS: Passed`;
- contains only applicable `[B-CONFIRMED]` semantic blocks;
- is covered by the required human signatures on that outer manifest;
- pins all authoritative fingerprints;
- includes assurance, scope certificate, decisions, equivalence suite, and known residual risks.

Anything absent from that certified scope is out of reconstruction scope and requires a new decision or a reopened deconstruction cycle.

---

# Part 15. Exit E — Reconstruction Readiness

## 15.1. Strict conditions

Exit E remains `Pending` through steps 1–5 of §15.4. The stage-4 Exit E Content-Readiness Report validates all content, decision, confirmation, and package-member conditions available at that stage, but it is not the Exit E pass artifact and MUST NOT state or imply `Passed`. Exit E transitions to `Passed` only at step 6 when all conditions below are true and the signed outer bundle manifest's hashed payload validates with the literal field `EXIT-E-STATUS: Passed`:

1. Composite Exit A passed for the same pinned snapshot, or a documented later non-semantic packaging change preserves its inputs.
2. Every effective in-scope source unit is `[C-COVERED]`; all approved exclusions are exact, risk-assessed, disclosed, and removed from denominator arithmetic; no `[C-GAP]` or `[C-PARTIAL]` remains.
3. All `90`–`96` mandatory domains are complete with evidence-backed zero-domain records where genuinely empty; every semantic block has a valid semantic version/fingerprint.
4. The complete handbook semantic payload exists as `[B-POPULATED]` before candidacy, passes quality gates, is synchronized to synthesis, and is included in the candidate payload manifest without later semantic regeneration.
5. All reciprocal references, legal typed IDs including `COV`/`CND`/`MOD`/`PRF`/`HBK`, claims, ownership, persona purity, diagrams, and links validate.
6. Every eligible lifecycle scope has a complete `SM`; each ineligible decision is explicit.
7. Every UC has a complete behavioral contract and reciprocal BR/ER/IF references.
8. Every CFG has an effect, environment/state matrix, or an approved exclusion.
9. Every `PRV` finding in `0C` has a reciprocal `SEC` record.
10. Every effective-denominator Disabled `SRC`, `CMP`, and `UC` has exactly one canonical `CAP`; disabled/dormant capabilities are indexed, cross-linked, risk-described, and have an **Approved** target decision; overlaps have explicit relationships and approved decisions.
11. Every material fault class has an `FLT` record and the consolidated fault/resilience topology is complete.
12. Every persona profile has a unique legal `PRF` ID, is atomically synchronized with its `0A` persona-prefix row, has the complete §8.4 dependency closure validated against exact versions/fingerprints without mixed snapshots, is semantically fingerprinted, `[B-CONFIRMED]`, and is covered by a valid exact-PRF-binding `CNF`.
13. Every synthesis, PRF profile, and HBK handbook semantic payload exactly matches its typed-ID entry in the candidate payload manifest; every certification envelope is `[B-CONFIRMED]`, references a valid `CNF` for its exact typed record ID, `SEMANTIC-RECORD-VERSION`, and `SEMANTIC-CONTENT-FINGERPRINT`, and is neither draft, populated-only, nor stale.
14. Every reconstruction-affecting scope, target, capability, and migration decision has an approval envelope with `Approval Status: Approved` bound to the exact candidate decision-content version/fingerprint; decisions with `Approval Status: Proposed` cannot pass Exit E.
15. The reviewed equivalence suite maps every UC contract, BR/invariant property, eligible SM transition, external IF baseline, and disabled-capability safeguard to at least one test or an exact approved exclusion.
16. All pending synthesis references and GAP records are reconciled; no orphan/dangling/tombstone violation remains.
17. All material contradictions and decision-required items are resolved or represented by approved scope exclusions that preserve and disclose exact residual risk.
18. The system namespace, protocol version, environments, and snapshot are present and exactly equal across the Exit A report, Exit E Candidate Report, Exit E Content-Readiness Report, scope certificate, every snapshot-bearing package member, and outer bundle manifest; the source inventory denominator, export projections, claims, candidate payload manifest, confirmations, and remaining package members are fingerprint-pinned to that identity under §4.1.2 and the acyclic exclusions in §15.4. Missing identity or snapshot mixing fails closed.
19. Human architect/domain confirmation is recorded in `18-CONFIRMATIONS.md`; `[B-CONFIRMED]` is never machine-assigned, and confirmation-envelope attachment did not alter any semantic version/fingerprint.
20. At step 6, the outer bundle manifest payload directly declares the exact system/protocol/environment/snapshot identity, final check-registry fingerprint, and pre-signature input-set fingerprint; includes every package member other than itself, including the Exit E Content-Readiness Report and scope certificate; and contains `EXIT-E-STATUS: Passed`. All listed package-member file fingerprints, required certification/approval/CNF envelope bindings, cross-artifact snapshot identities, and required final checks recompute and validate, and its required human signature envelope validates under §15.5. The signed outer manifest is the authoritative non-cyclic Exit E completion attestation. No post-step-6 authoritative Exit E artifact is created.

There are no row-count shortcuts. “File exists,” “section exists,” or “zero rows” does not satisfy completeness without denominator-backed proof.

## 15.1.1. Universal deterministic packaging rules

The five packaging artifact types are `EXIT-E-CANDIDATE-REPORT`, `CANDIDATE-PAYLOAD-MANIFEST`, `EXIT-E-CONTENT-READINESS-REPORT`, `SCOPE-CERTIFICATE`, and `OUTER-BUNDLE-MANIFEST`. Each MUST instantiate the generic §4.1.2 payload and envelope boundaries and its exact ordered schema below or in §15.2/§15.4. For every packaging payload:

1. payload fields appear exactly once and in the declared order; an unknown, duplicate, reordered, or omitted required payload field is a structural error;
2. a declared optional or empty scalar is serialized as the exact literal `None`, never omitted, blank, null-like, or represented by prose;
3. an empty table retains its declared header and separator rows and contains no data row;
4. every normalized path uses §4.1; path normalization occurs before sorting, hashing, or binding;
5. every table is sorted by its schema-declared tuple after normalization, using ascending bytewise comparison of the exact preserved Unicode code points serialized as UTF-8; duplicate sort-key tuples are invalid unless the table schema explicitly permits them;
6. an ID list not represented as a table is sorted by `(record type, typed ID)`; a path/member list is sorted by `(normalized path, artifact type)`; no source discovery order is semantically significant;
7. semantically equal inputs presented in different source orders MUST serialize to one canonical payload and one payload fingerprint;
8. the `HASH-DOMAIN` identity is exactly `artifact-type + "|" + normalized-relative-path`, is included in the preimage, and cannot contain the payload fingerprint; the `POST-HASH-ARTIFACT-INSTANCE-BINDING` is computed only after hashing and is never part of its own preimage.

The declared table sort keys are: candidate validation or readiness checks by `(Check ID)`; blockers by `(Blocker ID)`; semantic counts by `(Record Type)`; decision counts by `(Decision Type)`; semantic bindings by `(Record Type, Typed Record ID)`; decision bindings by `(DEC ID)`; confirmed bindings by `(Record Type, Typed Record ID)`; approved decision bindings by `(DEC ID)`; included scope dimensions by `(Dimension, Value)`; source counts by `(SRC Kind)`; exclusions by `(DEC ID, Excluded Unit ID)`; residual risks by `(Risk ID)`; and every package-member table, including the outer manifest, by `(Normalized Path, Artifact Type)`. Fields described as a deterministic set or tuple use the corresponding rule above.

The required Exit E report check registry is `EXIT-E-CHECKS-v1`. Its canonical registry fingerprint is `sha256` over UTF-8 `EXIT-E-CHECK-REGISTRY|EXIT-E-CHECKS-v1` plus one LF and the §4.1.2 canonical serialization of this exact ordered table:

| Stage             | Check ID                                  | Requirement                                                                                          |
| :---------------- | :---------------------------------------- | :--------------------------------------------------------------------------------------------------- |
| CANDIDATE         | CANDIDATE-01-EXIT-A                       | §15.1 condition 1 and a valid pinned Composite Exit A                                                |
| CANDIDATE         | CANDIDATE-02-SOURCE-CLOSURE               | §15.1 condition 2                                                                                    |
| CANDIDATE         | CANDIDATE-03-SYNTHESIS-DOMAINS            | §15.1 condition 3                                                                                    |
| CANDIDATE         | CANDIDATE-04-HANDBOOK-QUALITY             | §15.1 condition 4 and the deterministic pre-candidate §13.4 checks                                   |
| CANDIDATE         | CANDIDATE-05-REFERENTIAL-INTEGRITY        | §15.1 condition 5                                                                                    |
| CANDIDATE         | CANDIDATE-06-LIFECYCLE                    | §15.1 condition 6, including no `Unknown` eligibility                                                |
| CANDIDATE         | CANDIDATE-07-USE-CASES                    | §15.1 condition 7                                                                                    |
| CANDIDATE         | CANDIDATE-08-CONFIGURATION                | §15.1 condition 8                                                                                    |
| CANDIDATE         | CANDIDATE-09-PRIVACY-SECURITY             | §15.1 condition 9                                                                                    |
| CANDIDATE         | CANDIDATE-10-DISABLED-CAPABILITIES        | §15.1 condition 10                                                                                   |
| CANDIDATE         | CANDIDATE-11-FAULTS                       | §15.1 condition 11                                                                                   |
| CANDIDATE         | CANDIDATE-12-PROFILES                     | populated PRF identity, dependency closure, and candidate eligibility portions of §15.1 condition 12 |
| CANDIDATE         | CANDIDATE-13-PENDING-AND-CONTRADICTIONS   | pre-confirmation portions of §15.1 conditions 16–17                                                  |
| CANDIDATE         | CANDIDATE-14-SNAPSHOT-INPUTS              | candidate-stage portions of §15.1 condition 18                                                       |
| CONTENT-READINESS | READINESS-01-CANDIDATE-IMMUTABILITY       | unchanged candidate report, manifest, semantic payloads, decisions, and dependencies                 |
| CONTENT-READINESS | READINESS-02-CONFIRMATIONS                | confirmed semantic bindings and §15.1 conditions 12–13 and 19                                        |
| CONTENT-READINESS | READINESS-03-DECISION-APPROVALS           | §15.1 condition 14                                                                                   |
| CONTENT-READINESS | READINESS-04-EQUIVALENCE                  | §15.1 condition 15                                                                                   |
| CONTENT-READINESS | READINESS-05-HANDBOOK-READABILITY         | complete fingerprint-bound human review under §13.4                                                  |
| CONTENT-READINESS | READINESS-06-NO-PENDING-OR-CONTRADICTIONS | §15.1 conditions 16–17 after review                                                                  |
| CONTENT-READINESS | READINESS-07-ENVELOPE-INTEGRITY           | every certification, approval, and CNF envelope binding available at stage 4                         |
| CONTENT-READINESS | READINESS-08-PREPACKAGE-MEMBERS           | complete stage-4 member set and package-member file fingerprints                                     |
| CONTENT-READINESS | READINESS-09-SNAPSHOT-CONSISTENCY         | exact system/protocol/environment/snapshot equality across the chain and members                     |
| CONTENT-READINESS | READINESS-10-RESIDUAL-RISK                | every approved exclusion and residual risk is disclosed and bound                                    |
| CONTENT-READINESS | READINESS-11-SEQUENCE                     | all stage-4 §15.4 dependency edges point only to earlier artifacts and Exit E remains `Pending`      |

A report's `Check Results.Check ID` set MUST equal exactly the rows for its stage: no missing, duplicate, unknown, or other-stage ID is permitted. Every required row is serialized and has `Result: Pass | Fail`; this registry defines no `Not-Applicable` result because each check validates its own applicability denominator and evidence-backed zero domains. A missing or malformed row is a structural report failure, not a pass or an omitted/unknown result. `Ready-For-Human-Review` or `Ready-For-Certificate` requires exact-set equality, every row `Pass`, every row's evidence-set fingerprint to recompute, and an empty blockers table. Therefore a header-only empty check table can never satisfy either report. Any failed check requires at least one corresponding blocker; a blocker cannot compensate for a missing check.

## 15.1.2. Exit E Candidate Report ordered payload schema

The Exit E Candidate Report is emitted before the candidate payload manifest. Its exact payload field order is:

```markdown
<!-- ARTIFACT-PAYLOAD: EXIT-E-CANDIDATE-REPORT -->

# Exit E Candidate Report

- **ARTIFACT-TYPE:** EXIT-E-CANDIDATE-REPORT
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
- **System / Protocol / Snapshot:** [pinned system namespace, protocol version, environments, and snapshot]
- **Required Check Registry / Fingerprint:** [EXIT-E-CHECKS-v1 and its canonical registry fingerprint]
- **Exit A Report Fingerprint:** [...]
- **Deterministic Candidate Input-Set Fingerprint:** [sha256 of the canonically sorted authoritative input path/fingerprint bindings]
- **Coverage Validation Summary:** [passed count, failed count, deterministic summary fingerprint]
- **Traceability Validation Summary:** [passed count, failed count, deterministic summary fingerprint]
- **Handbook Validation Summary:** [passed count, failed count, deterministic summary fingerprint]
- **Semantic Record Counts:**
  | Record Type | Count |
  | :---------- | ----: |
- **Decision Content Counts:**
  | Decision Type | Count |
  | :------------ | ----: |
- **Check Results:**
  | Check ID | Requirement | Result | Evidence Fingerprint |
  | :------- | :---------- | :----- | :------------------- |
- **Blockers:**
  | Blocker ID | Affected IDs | Description | Evidence Fingerprint |
  | :--------- | :----------- | :---------- | :------------------- |
- **Candidate Result:** Ready-For-Human-Review | Blocked
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256 under §4.1.2; this carrier field excluded]
<!-- END-ARTIFACT-PAYLOAD: EXIT-E-CANDIDATE-REPORT -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: EXIT-E-CANDIDATE-REPORT -->

- **Signature / Signatories / Timestamps:** [value or None]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: EXIT-E-CANDIDATE-REPORT -->
```

`Check Results` sorts by `(Check ID)` and uses `Pass | Fail`; `Blockers` sorts by `(Blocker ID)`. `Ready-For-Human-Review` requires exact `CANDIDATE` registry-set equality, every required check passing with a valid evidence-set fingerprint, and an empty blockers table. A failed check requires a blocker; a missing/duplicate/unknown check makes the report structurally invalid. The report MUST NOT reference or include the later candidate payload manifest, confirmations or CNFs, confirmation envelopes, Exit E Content-Readiness Report, scope certificate, or outer bundle manifest. This ordering prevents a candidate report/manifest hash cycle.

## 15.1.3. Candidate Payload Manifest ordered payload schema

The candidate payload manifest is emitted only after the immutable candidate report and references it one-way. Its exact payload field order is:

```markdown
<!-- ARTIFACT-PAYLOAD: CANDIDATE-PAYLOAD-MANIFEST -->

# Candidate Payload Manifest

- **ARTIFACT-TYPE:** CANDIDATE-PAYLOAD-MANIFEST
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
- **Exit E Candidate Report Fingerprint:** [...]
- **Semantic Bindings:**
  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | Dependency-Set Fingerprint | Normalized Path |
  | :---------- | :-------------- | ----------------------: | :--------------------------- | :------------------------- | :-------------- |
- **DEC Content Bindings:**
  | DEC ID | Decision-Content Version | Decision-Content Fingerprint |
  | :----- | -----------------------: | :--------------------------- |
- **Semantic Binding Count:** [non-negative integer]
- **DEC Content Binding Count:** [non-negative integer]
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256 under §4.1.2; this carrier field excluded]
<!-- END-ARTIFACT-PAYLOAD: CANDIDATE-PAYLOAD-MANIFEST -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: CANDIDATE-PAYLOAD-MANIFEST -->

- **Signature / Signatories / Timestamps:** [value or None]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: CANDIDATE-PAYLOAD-MANIFEST -->
```

`Semantic Bindings` sorts by `(Record Type, Typed Record ID)` and `DEC Content Bindings` by `(DEC ID)`. Each count MUST equal its table's data-row count. The manifest MUST NOT include future DEC approval envelopes, CNFs, confirmation envelopes, the Exit E Content-Readiness Report, scope certificate, or outer bundle manifest. The candidate report does not reference this manifest; only this manifest references the earlier report.

## 15.1.4. Exit E Content-Readiness Report ordered payload schema

The stage-4 Exit E Content-Readiness Report is a pre-certificate artifact. Its exact payload field order is:

```markdown
<!-- ARTIFACT-PAYLOAD: EXIT-E-CONTENT-READINESS-REPORT -->

# Exit E Content-Readiness Report

- **ARTIFACT-TYPE:** EXIT-E-CONTENT-READINESS-REPORT
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
- **System / Protocol / Snapshot:** [pinned system namespace, protocol version, environments, and snapshot]
- **Required Check Registry / Fingerprint:** [EXIT-E-CHECKS-v1 and its canonical registry fingerprint]
- **Exit A Report Fingerprint:** [...]
- **Exit E Candidate Report Fingerprint:** [...]
- **Candidate Payload Manifest Fingerprint:** [...]
- **Confirmed Semantic Bindings:**
  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | Certification Envelope Binding | Semantic Containing Path | Semantic Package-Member File Fingerprint | CNF ID | CNF Payload Fingerprint | CNF Signature Envelope Binding | CNF Containing Path | CNF Package-Member File Fingerprint |
  | :---------- | :-------------- | ----------------------: | :--------------------------- | :----------------------------- | :----------------------- | :--------------------------------------- | :----- | :---------------------- | :----------------------------- | :------------------ | :---------------------------------- |
- **Approved DEC Bindings:**
  | DEC ID | Decision-Content Version | Decision-Content Fingerprint | Approval Envelope Binding | DEC Containing Path | DEC Package-Member File Fingerprint |
  | :----- | -----------------------: | :--------------------------- | :------------------------ | :------------------ | :---------------------------------- |
- **Pre-Certificate Package Members:**
  | Artifact Type | Normalized Path | Package-Member File Fingerprint |
  | :------------ | :-------------- | :------------------------------ |
- **Check Results:**
  | Check ID | Requirement | Result | Evidence Fingerprint |
  | :------- | :---------- | :----- | :------------------- |
- **EXIT-E-PREPACKAGE-STATE:** Pending
- **Content Readiness Result:** Ready-For-Certificate | Blocked
- **Blockers:**
  | Blocker ID | Affected IDs | Description | Evidence Fingerprint |
  | :--------- | :----------- | :---------- | :------------------- |
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256 under §4.1.2; this carrier field excluded]
<!-- END-ARTIFACT-PAYLOAD: EXIT-E-CONTENT-READINESS-REPORT -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: EXIT-E-CONTENT-READINESS-REPORT -->

- **Signature / Signatories / Timestamps:** [value or None]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: EXIT-E-CONTENT-READINESS-REPORT -->
```

Confirmed bindings sort by `(Record Type, Typed Record ID)`, approved decisions by `(DEC ID)`, pre-certificate members by `(Normalized Path, Artifact Type)`, checks by `(Check ID)`, and blockers by `(Blocker ID)`. Check results use `Pass | Fail`. `Ready-For-Certificate` requires exact `CONTENT-READINESS` registry-set equality, every required check passing with a valid evidence-set fingerprint, and an empty blockers table. A failed check requires a blocker; a missing/duplicate/unknown check makes the report structurally invalid. The payload MUST omit the outer-manifest Exit E completion field and completion literal and MUST NOT reference or depend on the later scope certificate or outer bundle manifest. It cannot pass Exit E; Exit E remains `Pending` until the signed outer-manifest step 6.

## 15.2. Scope certificate

```markdown
<!-- ARTIFACT-PAYLOAD: SCOPE-CERTIFICATE -->

# Scope Certificate

- **ARTIFACT-TYPE:** SCOPE-CERTIFICATE
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
- **System / Protocol / Snapshot:** [pinned system namespace, protocol version, environments, and snapshot]
- **Included Scope Dimensions:**
  | Dimension | Value |
  | :-------- | :---- |
- **Source Counts by Kind:**
  | SRC Kind | Gross Discovered | Approved Excluded | Effective Denominator | Covered |
  | :------- | ---------------: | ----------------: | --------------------: | ------: |
- **Approved Exclusions and Residual Risks:**
  | DEC ID | Excluded Unit ID | Decision-Content Version | Decision-Content Fingerprint | Approval Envelope Binding | Risk ID | Exact Residual Risk |
  | :----- | :--------------- | -----------------------: | :--------------------------- | :------------------------ | :------ | :------------------ |
- **Export Reconciliation:** [mapped count, approved-excluded count, total normalized units, reconciliation report fingerprint]
- **Exit A Report Fingerprint:** [...]
- **Exit E Candidate Report Fingerprint:** [...]
- **Candidate Payload Manifest Fingerprint:** [...]
- **Exit E Content-Readiness Report Fingerprint:** [...]
- **Confirmed Semantic Payload Bindings:**
  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | CNF ID | CNF Payload Fingerprint |
  | :---------- | :-------------- | ----------------------: | :--------------------------- | :----- | :---------------------- |
- **Known Residual Risks:**
  | Risk ID | Affected IDs | Source DEC ID | Exact Residual Risk | Disclosure Path |
  | :------ | :----------- | :------------ | :------------------ | :-------------- |
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256 under §4.1.2; this carrier field excluded]
<!-- END-ARTIFACT-PAYLOAD: SCOPE-CERTIFICATE -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: SCOPE-CERTIFICATE -->

- **Signature / Signatories / Timestamps:** [value or None]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: SCOPE-CERTIFICATE -->
```

`Included Scope Dimensions` sorts by `(Dimension, Value)`; source counts by `(SRC Kind)`; approved exclusions by `(DEC ID, Excluded Unit ID)`; confirmed bindings by `(Record Type, Typed Record ID)`; and residual risks by `(Risk ID)`. Every count and arithmetic relation MUST reproduce the cited denominator/reconciliation reports. The scope certificate is emitted only after the Exit E Content-Readiness Report. It references that immutable report and the candidate/confirmation chain; the content-readiness report does not reference or hash the later certificate. The certificate uses the generic §4.1.2 artifact boundaries and exclusions; all certificate scope, risk, confirmed typed-record bindings, and report-reference fields remain hashed.

## 15.3. Package contents

The certified package contains:

- human reconstruction handbook;
- structured `90`–`96` synthesis catalogs;
- assurance plane records and deterministic validation summaries;
- persona profiles and necessary forensic reference records;
- decision log;
- equivalence/acceptance suite;
- scope certificate;
- Exit E Candidate Report, candidate payload manifest, Exit E Content-Readiness Report, and signed outer bundle manifest;
- source, projection, artifact, and semantic fingerprint manifests;
- Exit A, Exit E Candidate, and content-readiness validation reports;
- known-gap/risk record, which MUST disclose every approved exclusion and its residual risks even though excluded from the effective denominator; it MUST contain no undisclosed or unapproved gap and MUST NOT claim a gap absent from both the gross discovered population and approved risk record.

Raw secrets, raw private exports, private matching hints, and private entity indexes are excluded.

## 15.4. Acyclic certification sequence and hash domains

Certification MUST follow this order; no later artifact is an input to an earlier hash:

1. **Populate semantic payloads.** Final-Synthesis and Profile Synchronization produce complete `[B-POPULATED]` synthesis and `PRF` profile payloads. Reconstruction-Handoff projects the complete handbook as `[B-POPULATED]` `HBK` records, runs the deterministic pre-candidate non-vacuity/link/diagram and synchronization checks in §13.4, and computes every typed record's `SEMANTIC-CONTENT-FINGERPRINT`. Human readability review does not occur in this stage; it is candidate-bound confirmation work in step 3. No confirmation envelope is required yet.
2. **Emit candidate artifacts.** Validation/Gate validates Exit A, all populated typed payloads, proposed decision content, references, coverage, and handbook quality, then emits an immutable Exit E Candidate Report and candidate payload manifest using the generic §4.1.2 artifact schema. The manifest lists every candidate synthesis, PRF, and HBK record ID, semantic record version, semantic-content fingerprint, dependency fingerprint, and required decision-content version/fingerprint binding. Neither candidate artifact includes future DEC approval envelopes, CNFs, confirmation envelopes, the Exit E Content-Readiness Report, certificate, or outer manifest.
3. **Human review and envelope attachment.** Authorized humans approve required `DEC` content by attaching approval envelopes bound to the exact candidate `DECISION-CONTENT-VERSION`/`DECISION-CONTENT-FINGERPRINT`, then create `CNF` records against exact candidate typed-record bindings, including PRF/HBK IDs where applicable. These approval/status attachments do not alter candidate-bound content fingerprints and do not by themselves restart candidacy. Every CNF that confirms candidate HBK content records the §13.4 handbook readability result; the combined passing CNFs MUST cover every candidate HBK binding. The authorized confirmation action attaches `[B-CONFIRMED]`/CNF envelopes without regenerating or changing semantic payload. DEC and CNF hashes use §4.1.2 and their exact record-specific envelope markers.
4. **Emit Exit E Content-Readiness Report.** Validation/Gate recomputes candidate semantic fingerprints, verifies unchanged payloads, recomputes and validates every certification/approval/CNF envelope binding and finalized package-member file fingerprint available at this stage, and validates every content/confirmation/package-member condition available before certificate and outer-manifest creation. It emits the immutable Exit E Content-Readiness Report using artifact type `EXIT-E-CONTENT-READINESS-REPORT`. Its hashed payload may reference Exit A, the Exit E Candidate Report/candidate payload manifest, DEC/CNF records, exact envelope bindings, and all validated pre-certificate package-member file fingerprints available at this stage. It MUST state that Exit E is still `Pending`, MUST NOT contain `EXIT-E-STATUS: Passed`, and MUST NOT reference or depend on the later scope certificate or outer bundle manifest/fingerprint. It is a pre-package content-readiness artifact, not the Exit E pass artifact.
5. **Emit scope certificate.** Reconstruction-Handoff emits the scope certificate using artifact type `SCOPE-CERTIFICATE`, referencing the immutable Exit E Content-Readiness Report fingerprint and prior chain. Creating it MUST NOT change the content-readiness report. Exit E remains `Pending`.
6. **Emit, validate, and sign the authoritative outer bundle manifest; transition Exit E.** Reconstruction-Handoff emits the outer manifest using artifact type `OUTER-BUNDLE-MANIFEST`. Its hashed payload directly declares the pinned `System / Protocol / Snapshot`, the final check-registry identity/fingerprint, and the pre-signature authoritative input-set fingerprint; lists every final package member other than the manifest itself and each external package-member file fingerprint, including the candidate artifacts, DEC/CNF records and their finalized containing files, confirmed semantic records and their finalized containing files, Exit E Content-Readiness Report, and scope certificate; includes the literal field `- **EXIT-E-STATUS:** Passed`; and does not list itself or its own digest as a package member. Validation/Gate requires exact system/protocol/environment/snapshot equality across Exit A, candidate report, content-readiness report, scope certificate, every snapshot-bearing package member, and the outer payload; recomputes every package-member file fingerprint and required certification/approval/CNF envelope binding; verifies the complete member set, generic payload/envelope boundaries, the content-readiness and certificate references, pre-signature input-set fingerprint, final check-registry fingerprint, and all §15.1 conditions. Missing identity, mixed snapshots, or a validly hashed member from another snapshot fails closed. Authorized humans then sign the outer artifact envelope, and the conforming implementation performs the §15.5 deterministic final-bundle verification over the completed signature envelope. Only when every required final check passes does Exit E transition from `Pending` to `Passed`. The signed outer manifest is the authoritative non-cyclic Exit E completion attestation and its payload digest is the bundle fingerprint. No post-step-6 authoritative completion report or attestation is emitted.

The outer payload's `Pre-Signature Validation Input-Set Fingerprint` is computed before outer payload hashing and signing. Its complete preimage is UTF-8 `EXIT-E-PRE-SIGNATURE-INPUTS|`, the canonically serialized outer `System / Protocol / Snapshot`, one LF, the `EXIT-E-CHECKS-v1` identity and fingerprint, one LF, the `EXIT-E-FINAL-CHECKS-v1` identity and fingerprint, one LF, the Exit E Content-Readiness Report artifact-payload fingerprint, one LF, the scope-certificate artifact-payload fingerprint, one LF, and the §4.1.2 canonical serialization of exactly the outer `Package Members` table. Each member row MUST use `PACKAGE-MEMBER-FILE` and the complete finalized path-bound fingerprint for that normalized path/artifact type. The table MUST equal the complete final member set excluding only the outer manifest itself; a missing, extra, duplicate, differently domained, or mixed-snapshot binding changes or invalidates the preimage. The preimage excludes the outer payload, its fingerprint carrier, and its signature envelope, so it is acyclic.

The authoritative step-6 outer manifest MUST instantiate the generic schema as follows; `EXIT-E-STATUS` is inside the hashed payload, while signatures are outside it:

```markdown
<!-- ARTIFACT-PAYLOAD: OUTER-BUNDLE-MANIFEST -->

- **ARTIFACT-TYPE:** OUTER-BUNDLE-MANIFEST
- **ARTIFACT-PATH:** [normalized relative path]
- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** [...]
- **System / Protocol / Snapshot:** [pinned system namespace, protocol version, environments, and snapshot]
- **Final Verification Check Registry / Fingerprint:** [EXIT-E-FINAL-CHECKS-v1 and its canonical registry fingerprint]
- **Pre-Signature Validation Input-Set Fingerprint:** [computed by the exact §15.4 preimage]
- **Exit E Content-Readiness Report Fingerprint:** [...]
- **Scope Certificate Fingerprint:** [...]
- **Package Members:**
  | Artifact Type | Normalized Path | Package-Member File Fingerprint |
  | :------------ | :-------------- | :------------------------------ |
- **EXIT-E-STATUS:** Passed
- **ARTIFACT-PAYLOAD-FINGERPRINT:** [sha256 under §4.1.2; this carrier field excluded]
<!-- END-ARTIFACT-PAYLOAD: OUTER-BUNDLE-MANIFEST -->

<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->

- **Required Human Signatories / Signature Timestamps:** [...]
- **Envelope Audit:** [...]
<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->
```

`Package Members` MUST be sorted by `(Normalized Path, Artifact Type)`, MUST contain every final package member except the outer manifest itself, and MUST reject duplicate normalized paths or duplicate member bindings. The declared outer-manifest field order is exact under §15.1.1.

**Exact exclusion rule:** the normative semantic/payload hash profile and exact boundaries are §4.1.2. Packaging payload fingerprints exclude only their own `ARTIFACT-PAYLOAD-FINGERPRINT` carrier and matching generic artifact digest/signature envelope. Synthesis/PRF/HBK semantic records exclude exactly the certification envelope in §4.1.1; `DEC` excludes exactly the decision-approval envelope in §5.5 and its own carrier; `CNF` excludes exactly the CNF digest/signature envelope in §5.7 and its own carrier. Those exclusions do not apply to canonical envelope fingerprints or external package-member file fingerprints: envelope fingerprints include every field within their exact envelope, and package-member file fingerprints include the complete finalized member. Cross-document references, content timestamps, package-member lists, `EXIT-E-STATUS`, and all other payload fields remain hashed. No artifact may exclude a later artifact reference merely to mask a cycle; such a reference is forbidden by the sequence.

**Restart rule:** a semantic payload, dependency, candidate decision-content binding, or candidate-manifest change restarts at step 1 or 2 as applicable and invalidates downstream CNFs. A DEC review-input or decision-content change restarts at step 2; attaching a matching approval envelope occurs at step 3 and does not restart, while a mismatched envelope is rejected. A CNF or confirmation-envelope defect restarts at step 3. An Exit E Content-Readiness Report payload defect restarts at step 4. A certificate-only payload defect restarts at step 5. An outer-manifest membership, status, fingerprint, order, validation, or signature defect restarts at step 6 and leaves Exit E `Pending`. Replacing any finalized member envelope changes its envelope binding and package-member file fingerprint; it restarts from the stage that owns that envelope and always invalidates/requires regeneration and re-signing of the step-6 outer manifest. Replacing only the outer manifest's own signature envelope with an unchanged outer payload digest restarts step 6 only and requires revalidation/re-signing. Any change that reaches backward into an earlier artifact restarts from the earliest affected step. There is no post-step-6 authoritative artifact whose hash could create a cycle.

## 15.5. Reproducible final-bundle verification

`EXIT-E-FINAL-CHECKS-v1` is the exact final verification registry:

| Check ID                     | Requirement                                                                                                    |
| :--------------------------- | :------------------------------------------------------------------------------------------------------------- |
| FINAL-01-OUTER-PAYLOAD       | exact outer schema, canonical payload fingerprint, field order, and `EXIT-E-STATUS: Passed`                    |
| FINAL-02-SNAPSHOT            | direct exact system/protocol/environment/snapshot equality under §15.1 condition 18                            |
| FINAL-03-MEMBER-SET          | complete package-member set, no self-member, missing member, extra member, or duplicate binding                |
| FINAL-04-MEMBER-FINGERPRINTS | every path-bound finalized package-member file fingerprint recomputes                                          |
| FINAL-05-ENVELOPES           | every required certification, decision-approval, CNF-signature, and outer-signature envelope/binding validates |
| FINAL-06-CHAIN               | candidate, manifest, confirmation, readiness, and certificate references follow the §15.4 order and recompute  |
| FINAL-07-REGISTRY-AND-INPUTS | report/final registry fingerprints and pre-signature authoritative input-set fingerprint recompute             |
| FINAL-08-CONDITIONS          | every applicable §15.1 condition and required report check passed with complete evidence bindings              |
| FINAL-09-SEQUENCE            | Exit E remained `Pending` through step 5 and the bundle contains no post-step-6 completion member or reference |

Its canonical fingerprint is `sha256` over UTF-8 `EXIT-E-FINAL-CHECK-REGISTRY|EXIT-E-FINAL-CHECKS-v1`, one LF, and the §4.1.2 canonical serialization of that table. A conforming runtime MUST expose a deterministic `verify final bundle` operation over the complete signed bundle. The operation recomputes the exact registry and check results from bundle bytes; it does not trust the outer payload's `Passed` literal, stored fingerprints, or executor-supplied values. Missing, duplicate, unknown, or failed checks yield overall `Fail`. `FINAL-09-SEQUENCE` verifies only bundle-observable members, references, and recorded prior-stage states. The prohibition on creating later authoritative completion state is separately enforced by the runtime's post-Exit-E transaction rules and cannot be inferred from standalone bundle bytes.

Every execution records a **final verification receipt** with this schema:

```markdown
- **Protocol / Artifact / Check Registry Versions and Fingerprints:** [...]
- **Normalized Bundle Identity:** [...]
- **Outer Payload / Canonical Outer-Envelope Fingerprints:** [...]
- **System / Protocol / Environment / Snapshot:** [...]
- **Package-Member / Input Bindings and Input-Set Fingerprint:** [exact sorted set]
- **Check Results:**
  | Check ID | Result | Offending Bindings |
  | :------- | :----- | :----------------- |
- **Overall Result:** Pass | Fail
- **Validator Implementation Version / Timestamp:** [diagnostic provenance]
- **Deterministic Result Fingerprint:** [...]
```

The result fingerprint covers every preceding field except `Validator Implementation Version / Timestamp` and its own carrier under §4.1.2 with domain prefix `FINAL-VERIFY-RESULT|`. The receipt is a `NON-AUTHORITATIVE-DERIVATIVE`, MUST remain outside the certified package and authoritative `.extracted/` protocol state, MUST NOT be a package member or input to Exit E, and cannot change the outer payload digest, signature, or gate state. It records which checks the original runtime executed; any third party can discard it and reproduce the same deterministic result from the signed bundle. Creating, deleting, or regenerating a receipt is not a post-step-6 completion artifact or attestation.

---

# Part 16. Security, Privacy, and Operational Guardrails

1. Production data MUST NOT be copied into probes. Operators verify masking and organizational handling approval; agents record evidence but do not declare compliance.
2. Raw credentials, tokens, cookies, keys, signed URLs, connection strings, and secret values MUST NOT be written to `.extracted/`.
3. Sensitive endpoints, topology, and identifiers are redacted or aliased according to policy while preserving safe behavioral shape.
4. External content and generated helper output are untrusted until schema, fingerprint, sanitization, and lineage validation pass.
5. Acquisition is read-only and operator-controlled. Production mutation, workflow execution, or feature enablement requires a separate explicitly authorized operational process outside this protocol.
6. High-risk probes require human authorization, rollback, isolation, and captured scope. Disabled production functionality is never enabled for observation.
7. Security/privacy findings are evidence and risk inputs for qualified human review, not legal conclusions.
8. Secret scanning failures block publication. Unknown field categories fail closed.
9. Snapshot mixing is forbidden. A source changed during acquisition/traversal requires a consistent re-snapshot and stale propagation.

---

# Part 17. Migration from v3.17 Workspaces

Migration is an explicit procedure identified by `Protocol-Migration Metadata` and executed through the existing invocation modes; it is not a separate invocation mode or status prefix and does not claim closure during conversion.

## 17.1. Required order

1. Freeze and fingerprint the v3.17 workspace and source snapshot.
2. Create the v4 source inventory denominator before interpreting old coverage.
3. Allocate deterministic IDs for source, components, claims, and synthesis candidates; preserve old headings as aliases.
4. Expand wildcard clusters and exclusions into concrete coordinates.
5. Split routine families, modules, serialized containers, and aggregated blocks into universal atomic records; retain non-promotable grouping containers for navigation.
6. Inventory and classify Normal, Shadow/Conditional, and Disabled tracks separately by environment/snapshot.
7. Apply deterministic persona-prefix uniqueness and entry ownership; create shared-entry records where intentional.
8. Move each legacy unprefixed persona directory to its registered `<persona-prefix>-<persona-slug>` basename as a proven path move; preserve record IDs and history, record the prior path as an alias, update every path reference and invocation target, and recompute affected file-transport fingerprints without changing prefix-derived PRF identity.
9. Re-run logical explosion and reconcile every normalized export coordinate.
10. Convert block-level evidence into one-fact CLM records; split mixed evidence and derive `BLOCK-CONFIDENCE-RANK` plus `BLOCK-EVIDENCE-PROFILE` mechanically.
11. Build traversal frontier records from depth boundaries, unresolved dependencies, missing exports, and old tickets.
12. Convert direct human-required records into migration findings; reopen them at the static-investigation/probe stage unless an existing authorized exclusion proves the scope disposition.
13. Populate traceability, reciprocal references, contradiction records, and coverage arithmetic.
14. Mark old synthesis-like outputs provisional `[B-DRAFT]`; do not inherit confirmation.
15. Expose all gaps before enabling strict validators. Never backfill evidence merely to satisfy a gate.
16. Run Sequential and Cross-Reference Reconciliation, then evidence-backed sweep validation.
17. Enable Composite Exit A validators only after denominator and frontier stabilization.
18. Run Final Synthesis, handbook projection, human review, and Exit E.

## 17.2. Migration guarantees

Migration MUST preserve ticket history, evidence anchors, fingerprints, old identifiers as aliases, promotion/depromotion history, and uncertainty. It MUST NOT convert disabled to dead, terminal tickets to covered source, inferred evidence to direct, or old synthesis to confirmed. The Draft Principle remains in force; `[B-CONFIRMED]` always requires fresh human confirmation against the pinned v4 bundle.

## 17.3. v3 final-artifact compatibility

v4 supersedes v3 `FINAL-WORKFLOW-ATLAS.md` semantically with `93-USE-CASES.md` plus the corresponding handbook use-case, background-processing, and architecture chapters. v4 supersedes v3 `FINAL-HUMAN-AUDIT-QUEUE.md` semantically with the ticket/contradiction/decision/confirmation/known-gap and gate artifacts in the Forensic, Assurance, and Handbook planes. Neither v3 file is authoritative in a v4 workspace and neither satisfies an Exit A or Exit E requirement.

A migration MAY retain either filename only as a `NON-AUTHORITATIVE-COMPATIBILITY-PROJECTION` for legacy consumers. Such a projection MUST identify its v4 upstream IDs, semantic versions/fingerprints, generation time, projection schema, and staleness state; MUST contain no independent claims, approvals, queue state, or decisions; and MUST be excluded from authoritative denominator and certification inputs except as a listed derivative package member. Conflicts resolve in favor of canonical v4 records.

---

# Part 18. Migration Checklist

- [ ] Snapshot and v3.17 workspace fingerprints pinned.
- [ ] `10-SOURCE-INVENTORY.md` built for every applicable kind.
- [ ] Stable IDs allocated and collision-checked.
- [ ] Wildcards expanded to concrete units.
- [ ] Family/container blocks split into atomic records.
- [ ] Normal, Shadow/Conditional, and Disabled tracks classified separately.
- [ ] Capability-state environment matrices recorded.
- [ ] Persona prefixes, explicit persona slugs, canonical `<persona-prefix>-<persona-slug>` basenames, and canonical entry owners validated.
- [ ] Legacy unprefixed persona directories migrated as proven path moves with aliases, updated references, and stable prefix-derived PRF identity.
- [ ] DB-autonomous ownership applied deterministically.
- [ ] Serialized coordinates fully exploded and reconciled.
- [ ] Material evidence converted to single-fact CLMs.
- [ ] Frontiers reconstructed and terminal dispositions validated.
- [ ] Direct human-required shortcuts removed or justified through Human Hatch/exclusion.
- [ ] Traceability and reciprocal references populated.
- [ ] Coverage gaps exposed before strict gating.
- [ ] Old synthesis marked draft and re-derived.
- [ ] Exit A, Final Synthesis, handbook review, content-readiness validation, scope certificate, and signed outer-manifest Exit E transition run in order.

---

# Part 19. Conformance and Acceptance Checklist

## 19.1. Conforming implementation capabilities

- [ ] Identity headers enforce exact persona/cluster/track/POV/file/ledger isolation.
- [ ] Executors propose semantic content only; the conforming runtime independently computes or reproduces every protocol-defined ID, canonical form, fingerprint, scope/stale check, FSM transition, count, completeness result, gate result, and package verification before commit.
- [ ] Cold-resume results bind executor semantic assessment to runtime-recomputed identity/scope/write-right/stale checks and the mandatory `0G` no-mutation or commit entry.
- [ ] Workspace discovery treats `personas/` as the sole persona container, validates each registered persona basename as exact `<persona-prefix>-<persona-slug>` with registry/header/ticket/PRF prefix agreement, validates canonical shared records under `personas/_shared/`, rejects malformed/unregistered persona directories and `_shared` as a prefix, slug, basename, or persona row, and never counts the reserved directory as a persona.
- [ ] Multi-file writes are restricted to enumerated modes and enforce the `personas/{persona-directory}/` versus `personas/_shared/` ownership boundary.
- [ ] ID generator, collision extension, aliases, versions, tombstones, and `COV`/`CND`/`MOD`/`PRF`/`HBK` types and canonical coordinates are implemented.
- [ ] HBK paths and assigned anchors implement §4.1 exactly; renderer-generated heading anchors are never identity inputs.
- [ ] The §4.1.2 canonical hash profile, per-record and file-transport fingerprints, exact carrier/envelope exclusions, and acyclic package hashes are implemented.
- [ ] Packaging hash-domain identity is exactly `(artifact type, normalized path)`; the payload fingerprint is absent from its own preimage, and the artifact-instance binding is constructed only after hashing.
- [ ] Certification, DEC-approval, and CNF-signature envelopes have identity/version-bound canonical fingerprints; final package-member file fingerprints cover complete finalized files and are distinct from envelope-excluding transport fingerprints.
- [ ] Every finite enum is dispatched by its complete §5.1 Markdown-AST schema path; label-only dispatch, undeclared paths, wrong-family tokens, and invalid traversal `N/A` coupling fail closed.
- [ ] Source enumerators cover every required kind and can prove zero domains.
- [ ] Exactly one FRT is created for every discovered traversal boundary, including immediately terminal boundaries and depth stops.
- [ ] CLM schema rejects mixed evidence and derives non-lossy `BLOCK-CONFIDENCE-RANK` plus `BLOCK-EVIDENCE-PROFILE`.
- [ ] Helper publication fails closed and keeps raw/private state outside the target repository and authoritative `.extracted/` workspace in operator-controlled storage.
- [ ] Export reconciliation is one row per normalized coordinate.
- [ ] Ticket FSM forbids direct human-required creation, validates the closed escalation-reason taxonomy, and permits non-runtime Human Hatch escalation only after bounded static exhaustion and a reason-consistent probe inapplicability assessment.
- [ ] Human-required terminalization preserves gap/partial coverage and remains distinct from exact approved exclusion.
- [ ] Probe workflow fingerprints and sanitizes logs.
- [ ] Concurrent persona buffers, including promotion requests and unmapped discoveries, merge sequentially, deterministically, and atomically before later Discovery; only Sequential Reconciliation moves canonical blocks between `personas/_shared/` and persona POV files.
- [ ] Synthesis GAP records deduplicate deterministically and create required ticket/frontier work only through Sequential Reconciliation.
- [ ] Promotion/depromotion and disabled/dead distinctions validate.
- [ ] Sweep records reproduce expected/mapped/frontier/ticket counts.
- [ ] Coverage arithmetic runs per kind, track, environment, and snapshot.
- [ ] Dead-code validation rejects any linked unresolved possible caller/dynamic reference and confines the effect to the linked dependency closure.
- [ ] Semantic applicability/lifecycle classifications carry provenance; `Unknown` is fail-inclusive and cannot silently justify omission or closure.
- [ ] Reciprocal, orphan, dangling, ownership, and persona-purity checks run.
- [ ] Targeted stale propagation follows `DERIVED-FROM` impact edges.
- [ ] Profile Synchronization loads and validates the complete direct/applicable and transitive §8.4 dependency closure; confirmation reuses the identical candidate-bound closure.
- [ ] Synthesis gap buffering and controlled ticket creation are separated.
- [ ] Human confirmation is required for every B-confirmed block and binds exact semantic versions/fingerprints without changing payloads.
- [ ] All five packaging artifacts enforce exact field order, required/optional/empty representation, schema-declared row sorting, stage-specific exact required-check sets, evidence-set fingerprints, and rejection of unknown, duplicate, reordered, omitted, or stage-ineligible fields/checks.
- [ ] Exit E Candidate Report, candidate payload manifest, Exit E Content-Readiness Report, later scope certificate, and signed outer manifest follow §15.4; direct snapshot equality is enforced across the chain and only the validated signed outer payload can set `EXIT-E-STATUS: Passed`.
- [ ] Final-bundle verification executes the exact final registry and emits only a reproducible non-authoritative receipt outside package/protocol state.
- [ ] Handbook links, diagrams, and evidence-rank/profile display pass the deterministic §13.4 checks before candidacy; candidate-bound human readability review is recorded only through complete passing CNFs after candidacy.
- [ ] Validation summaries are deterministic and fingerprinted.

## 19.2. Normative conformance cases

A conforming implementation MUST satisfy every applicable positive and negative case below and deterministically produce the specified acceptance or rejection. These are protocol-level behavioral cases, not required files or test-runner inputs; their storage, serialization, and execution mechanism are implementation-specific. Failure of an applicable case defeats a claim of protocol conformance. Workspace gate execution is governed independently by the authoritative inputs and gate checks defined elsewhere in this protocol:

- **Persona workspace layout and ownership:** positive cases enumerate two registered persona directories beneath `personas/` whose basenames exactly equal `<persona-prefix>-<persona-slug>`, preserve explicit display names without deriving slugs, prove basename/registry/invocation/ticket/PRF prefix agreement, load canonical shared records only from `personas/_shared/`, route persona-local buffers beneath `personas/{persona-directory}/`, preserve purity/shared-reference behavior, migrate an unprefixed legacy directory as a proven path move without changing prefix-derived PRF identity, and execute an eligible promotion only through one atomic Sequential Reconciliation move; negative cases reject root-level persona directories, a root-level `shared/`, `_shared` as a prefix, slug, basename, or registry row, unprefixed or malformed persona basenames, uppercase or non-ASCII slugs, prefix mismatch, inferred/transliterated display-name slugs, unregistered persona directories, canonical shared records in a persona directory, persona records in `_shared`, writes that cross the ownership boundary, direct Promotion Review moves, non-atomic source/target updates, duplicate canonical copies, and a move that temporarily leaves no canonical copy.
- **Context-qualified finite enums and registry completeness:** a bidirectional schema-to-registry case discovers every normative finite-choice bullet/table field and requires exactly one matching §5.1 path, while a registry-to-schema case rejects orphan or duplicate paths. Positive value cases exercise every declared path, including every status-bearing table column, all block-confidence/evidence-profile contexts, schema-local synthesis choices, and the coupled traversal `Applicability`/`Status` `N/A` case; negative cases reject label-only dispatch, wrong-family tokens at every path, prose mixed with tokens, unbracketed prefix tokens, bracketed unprefixed literals, `N/A` without `Not-Applicable`, an R-token with `Not-Applicable`, invalid table-cell tokens, a schema enum omitted from the registry, and an undeclared/orphan registry path.
- **Invocation ownership and cold resume:** positive cases bind executor-proposed semantic actions/blockers to runtime-recomputed identity, mode, scope, loaded fingerprints, stale state, write rights, cold-resume fingerprint, and the same transaction's `0G` entry; negative cases reject a missing summary, executor-authored deterministic digest that does not reproduce, stale checkpoint, changed target, forbidden intended write, omitted required load, mixed snapshot, mismatched check fingerprint, mutation after a failed check, and interrupted work resumed without a fresh check.
- **Ticket escalation and honest unresolved coverage:** positive cases exercise `RUNTIME-REQUIRED` through `[T-PROBE-REQUIRED]` and every legal non-runtime escalation reason from `[T-OPEN]`/`[T-FOLLOW-UP]`, terminalize only through authorized Human Hatch, preserve linked `[C-GAP]`/`[C-PARTIAL]`, permit unaffected work and bounded partial synthesis, and distinguish a later exact approved exclusion; negative cases reject direct `[T-HUMAN-REQUIRED]`, `None` reason, missing bounded static investigation, runtime escalation without probe feasibility, non-runtime escalation through a fake probe state, reason/assessment mismatch, unauthorized reviewer, and human-required counted as covered or silently excluded.
- **Dead-code and semantic-predicate closure:** positive cases prove a complete no-caller denominator across bound tracks/environments/snapshots, permit a terminal exactly excluded unavailable caller domain, map traversal `Required` rows one-to-one to independent sweep records, and resolve lifecycle eligibility with provenance; negative cases reject `[S-DEAD-CODE]` with an unresolved dynamic reference/candidate/ticket/frontier/possible caller, any observed conditional/disabled/historical/intended caller, missing denominator dimension, an effect on unrelated coverage, lifecycle/applicability without provenance, `Unknown` treated as ineligible/not-applicable, a missing required `SM`, vacuous handbook content, and inherited POV sweep results.
- **Deterministic iteration accounting:** positive cases prove that the first invocation opens `ALFA`, parallel invocations share one token, an open-token synthesis gap consumes no additional token, and a post-completion reopen advances exactly once; negative cases reject completion with unterminated invocations or unreconciled buffers, reuse/reset of a completed token after Exit A/synthesis/confirmation failure, advancement that skips a token, and any work or closure claim beyond completed `ZULU`.
- **PRF identity and semantic hashing:** positive cases derive one PRF ID from system namespace plus persona prefix, bind it reciprocally through `0A`, traceability, candidate manifest, and CNF, and prove that a semantic-record-version change changes its fingerprint while a fingerprint-carrier-only or certification-envelope-only change does not; negative cases reject illegal/duplicate PRFs, prefix-coordinate mismatch, missing separately ordered semantic fields, carrier inclusion in its own preimage, and omitted/substituted candidate or CNF identity.
- **HBK identity and path/anchor normalization:** positive cases fix the sole identity root at `.extracted/handbook/`, preserve path case and exact Unicode code points beneath it, normalize POSIX separators plus repeated and `.` segments, use an explicitly assigned anchor matching `[a-z0-9]+(?:-[a-z0-9]+)*`, preserve identity across title-only edits, and apply the proven-move rule; negative cases reject package-root-prefixed or `.extracted/handbook/`-prefixed coordinate values, alternate roots, `..`, implicit Unicode normalization, transliteration, punctuation/uppercase/Unicode in anchors, renderer-generated anchors, duplicate path-anchor coordinates, unproven moves, missing coordinate fields, and omitted/substituted candidate or CNF identity.
- **Profile Synchronization closure:** a positive case regenerates a PRF from a multi-hop `DERIVED-FROM` graph containing direct/applicable and transitive records, verifies versions/fingerprints/snapshot/reciprocal bindings, and confirms it using the identical candidate-bound closure; negative cases reject a missing direct/applicable record, missing transitive record, stale dependency, version mismatch, fingerprint mismatch, broken reciprocal binding, mixed snapshots, or confirmation closure different from the generation closure.
- **Record versus artifact hashing:** positive cases independently hash two typed records in one Markdown file, produce the separate ordered file-transport fingerprint, hash each packaging artifact as one bounded payload using only `artifact-type + "|" + normalized-relative-path` as hash-domain identity, and construct the post-hash artifact-instance binding afterward; negative cases reject whole-file bytes as a record/CNF identity, record reordering without a changed transport fingerprint, artifact hashes computed from unbounded bytes, any payload fingerprint or post-hash binding in its own preimage, and construction of the instance binding before hashing.
- **Canonicalization and exact exclusions:** positive cases cover normalized LF, exact Unicode preservation, schema field/table order, insignificant whitespace, every semantic certification envelope, the exact DEC approval envelope, the exact CNF digest/signature envelope, every generic artifact payload/envelope boundary, every fingerprint-carrier exclusion, and an allowed candidate binding whose attachment preserves the semantic fingerprint while changing envelope and package-member fingerprints; negative cases mutate one included byte or move content across a boundary and require a hash change or structural rejection, reject a certification-envelope reference to content-readiness/certificate/outer artifacts, and prove that only explicitly authorized envelope changes leave the corresponding semantic/content fingerprint unchanged.
- **Final envelope and package-member integrity:** positive cases compute identity/version-bound fingerprints for certification, DEC-approval, and CNF-signature envelopes, serialize their exact bindings, and compute external path-bound package-member file fingerprints over complete finalized files while preserving unchanged semantic/decision/CNF payload fingerprints; negative cases remove or mutate only each envelope, signature, timestamp, binding, fingerprint carrier, or finalized file byte after content-readiness and require content-readiness/outer validation to fail, and reject use of the envelope-excluding transport fingerprint as a package-member fingerprint or storage of a package-member fingerprint inside its own member.
- **Packaging schemas, gate completeness, and deterministic ordering:** positive cases instantiate all five exact §15 schemas, require exact candidate/content-readiness registry sets, reproduce each evidence-set fingerprint from complete authoritative row bindings independent of source order, canonicalize semantically identical rows supplied in different source orders to one hash, serialize empty scalars as `None`, retain header/separator-only empty non-check tables, and verify declared counts/arithmetic; negative cases reject an empty required check table, one missing required check, unknown/duplicate/wrong-stage check, failed check without blocker, fabricated omission/`Not-Applicable`, wrong evidence domain, omitted/duplicate/stale/self/later evidence binding, executor-authored digest mismatch, unknown/duplicate/reordered/missing/extra fields, duplicate sort keys, noncanonical row order after normalization, omitted empty values, malformed empty tables, report-to-manifest back-reference, and any forbidden later-artifact reference.
- **Acyclic, snapshot-consistent Exit E and final verification:** positive cases execute candidate report -> candidate manifest -> DEC/CNF and confirmation envelopes -> Exit E Content-Readiness Report with Exit E `Pending` -> scope certificate -> outer manifest validation/signature -> exact step-6 transition to `EXIT-E-STATUS: Passed`, bind one direct snapshot identity across the chain, execute every `EXIT-E-FINAL-CHECKS-v1` check, and reproduce a non-authoritative receipt without changing the bundle; negative cases reject `Passed` at steps 1–5, a content-readiness report presented as the pass artifact, a certificate emitted before its content-readiness input, an outer manifest missing the report/certificate/direct identity/registry/input-set field, candidate-readiness-certificate-outer snapshot mismatch, a mixed-snapshot member despite valid individual hashes, payload self-reference, missing/invalid signature, missing/unknown/failed final check, a `Passed` field outside the hashed outer payload, a receipt treated as a package member or authority, and any authoritative post-step-6 completion artifact.

## 19.3. Corpus acceptance

- [ ] Composite Exit A is green with no hidden partial/gap denominator.
- [ ] Every in-scope source unit is atomic, covered, and terminally disposed.
- [ ] Every normalized unit is mapped or exactly excluded.
- [ ] Every required traversal-matrix cell has evidence-backed closure; every N/A cell is claim-backed.
- [ ] Active and disabled behavior are visibly separated but fully cross-linked.
- [ ] Every `MOD`/ER/REL/DR/BR/UC/IF/DEP/CFG/SCHED/NFR/SEC/SM/CAP/PRV/FLT block is claim-backed as applicable; purely navigational indexes carry no semantic status/evidence.
- [ ] Every BR has one authoritative owner and consistent optional mirrors.
- [ ] Every external IF has a compatibility baseline and migration rule.
- [ ] Every eligible lifecycle has a complete SM.
- [ ] Every UC can generate a reviewable acceptance-test skeleton.
- [ ] Every CFG has a behavioral effect or approved exclusion.
- [ ] Every PRV finding has a SEC mirror; every material failure class has an FLT record.
- [ ] Every effective-denominator Disabled SRC/CMP/UC has exactly one canonical CAP; every disabled capability has risk and an approved Restore/Preserve Dormant/Redesign/Retire decision.
- [ ] No synthesis or handbook block is stale, provisional, or pending-reference.
- [ ] Human readers can navigate core behavior without reading raw ID catalogs.
- [ ] Exit E transitioned at step 6 only; the scope certificate is included, and the validated signed outer manifest hash payload authoritatively contains `EXIT-E-STATUS: Passed`.

---

# Part 20. Compact Artifact Ownership Matrix

| Artifact / Record                                       | Semantic Purpose                                                | Authoritative Plane       | Authorized Writer(s)                                                                                                         | Key Non-Writers                                              |
| :------------------------------------------------------ | :-------------------------------------------------------------- | :------------------------ | :--------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------- |
| `0A-PREFLIGHT.md`                                       | scope, personas, clusters, acquisition registry                 | Forensic                  | Preflight/Acquisition; Profile Synchronization for persona rows; Sequential Reconciliation for buffered unmapped discoveries | isolated POVs write only persona-local unmapped buffers      |
| `0B` / `0C` / `0D` / `0F`                               | auth, security/privacy, glossary, global state                  | Forensic                  | scoped POV buffers + Sequential Reconciliation                                                                               | Synthesis/Handoff                                            |
| `0E-INDEX.md`                                           | forensic cross-reference index                                  | Forensic                  | Sequential/Cross-Reference Reconciliation                                                                                    | isolated POV direct global edits                             |
| `0G` / `0H`                                             | append-only history / derived checkpoint                        | Forensic                  | every mode appends `0G`; conforming implementation derives `0H`                                                              | no overwrite of `0G`                                         |
| `personas/_shared/01–05 POV`                            | canonical shared atomic components                              | Forensic                  | Sequential Reconciliation for merge, promotion, and depromotion transactions                                                 | parallel persona agents; Promotion Review direct moves       |
| `personas/{persona-directory}/01–05 POV`                | persona invocation references/promoted atoms                    | Forensic                  | bound isolated POV for references; Sequential Reconciliation for atomic promotion/depromotion                                | other personas; Promotion Review direct moves                |
| persona-local staging/question/unmapped buffers         | concurrent handoff and unmapped discovery                       | Forensic workflow         | bound isolated POV appends; Sequential Reconciliation consumes atomically                                                    | direct `personas/_shared/` or `0A` mutation by isolated POVs |
| `*-QUESTIONS.md`                                        | ticket FSM                                                      | Forensic                  | bound origin/target rights; Sequential Reconciliation for buffered creation                                                  | Synthesis direct writes                                      |
| normalized maps / probes                                | navigation / runtime telemetry                                  | Forensic                  | Acquisition / operator-controlled probe ingestion path                                                                       | behavioral synthesis as source authority                     |
| `10-SOURCE-INVENTORY.md`                                | denominator and dispositions                                    | Assurance                 | Acquisition, bound traversal, Sequential Reconciliation, Human Hatch                                                         | Handbook                                                     |
| `11-TRAVERSAL-FRONTIER.md`                              | traversal boundaries and closure                                | Assurance                 | bound traversal and Sequential Reconciliation                                                                                | Synthesis                                                    |
| `12-CLAIM-EVIDENCE.md`                                  | atomic evidence claims                                          | Assurance                 | evidence-producing bound invocation; reconciliation/validation                                                               | Handbook semantic edits                                      |
| `13-EXPORT-RECONCILIATION.md`                           | normalized-unit closure                                         | Assurance                 | Acquisition and Sequential Reconciliation                                                                                    | entity-correlation helper automation                         |
| `14-CONTRADICTIONS.md`                                  | explicit incompatible claims                                    | Assurance                 | evidence producers and Sequential Reconciliation                                                                             | silent synthesis selection                                   |
| `15-DECISIONS.md`                                       | authoritative scope/target/migration choices                    | Assurance/Governance      | `Human Hatch: Decision/Scope Approval` and authorized humans; Handoff projects only                                          | automatic inference                                          |
| `16-SOURCE-COVERAGE.md`                                 | scope-specific coverage authority                               | Assurance                 | bound traversal, Acquisition, Reconciliation, `Human Hatch: Ticket Escalation`/`Decision/Scope Approval` for exact effects   | synthesis assumptions                                        |
| `17-ACQUISITION-CANDIDATES.md`                          | helper candidate/action dispositions                            | Assurance                 | Acquisition and authorized operator review                                                                                   | automatic correlation-to-edge conversion                     |
| `18-CONFIRMATIONS.md`                                   | exact semantic-version/fingerprint human confirmations          | Assurance/Governance      | authorized humans via `Human Hatch: Confirmation` or Reconstruction-Handoff; envelope only                                   | machine assignment                                           |
| `20-TRACEABILITY.md`                                    | reciprocal mappings and ID registry                             | Assurance                 | Sequential/Cross-Reference Reconciliation and validators                                                                     | isolated semantic invention                                  |
| `21-COVERAGE-REPORT.md`                                 | per-kind arithmetic                                             | Assurance                 | deterministic validator                                                                                                      | manual evidence-rank claims                                  |
| `22-GATE-REPORTS.md`                                    | Exit A, Exit E Candidate, Content-Readiness, validation reports | Assurance                 | deterministic gate runner; content-readiness is pre-package and cannot pass Exit E                                           | isolated POVs                                                |
| `90–96`                                                 | structured reconstruction semantics                             | Synthesis                 | Partial/Final Synthesis semantic payloads; semantic-reference reconciliation; confirmation envelope only                     | Discovery/Ticket modes                                       |
| `PERSONA-PROFILE.md` / `PRF`                            | typed complete persona semantics                                | Forensic/Synthesis bridge | Profile Synchronization semantic payload and PRF/`0A` binding; authorized confirmation envelope writers                      | ordinary POV and Synthesis direct edits                      |
| synthesis gap buffer                                    | proposed gaps                                                   | Synthesis workflow        | Partial/Final Synthesis; consumed by Sequential Reconciliation                                                               | direct ledger mutation                                       |
| handbook / `HBK`                                        | typed browsable reconstruction projection                       | Human Handbook            | Reconstruction-Handoff generates HBK-identified populated payload; authorized confirmation actions attach envelope only      | independent semantic edits                                   |
| decision-log projection                                 | human-readable approved choices                                 | Human Handbook            | Reconstruction-Handoff from `15-DECISIONS.md`                                                                                | independent semantic edits                                   |
| equivalence suite                                       | acceptance derivative                                           | Handoff package           | Reconstruction-Handoff + reviewed generators                                                                                 | legacy source as uncited oracle after Exit E                 |
| candidate/content-readiness/certificate/outer artifacts | acyclic certification and authoritative Exit E completion       | Handoff governance        | ordered §15.4 Handoff/Validation stages; signed outer manifest alone attests `EXIT-E-STATUS: Passed`                         | pre-step-6 artifacts cannot pass Exit E                      |

---

**End of Canonical Deconstruction Protocol v4.1.1**
