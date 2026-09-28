# Roadmap

> Non-authoritative planning projection. Milestone completion never overrides [`protocol.md`](../protocol.md), and it does not imply Protocol v4 conformance unless the required Part 17 fixtures pass. The source-language contract is frozen in [`cdl.md`](cdl.md); execution contracts are in [`eir.md`](eir.md), [`execution-semantics.md`](execution-semantics.md), and [`runtime-architecture.md`](runtime-architecture.md).

## Two tracks

The project runs two dependency-ordered tracks:

```
S — Spec track     protocol ownership: language, EIR, native VM, prompt backends
M — Harness track  runtime services: workspace, service, agents, UI, operations
```

Ownership rule: protocol semantics — modes, states, gates, hashing, coverage,
packaging, confirmation transitions, iteration, reconciliation — are defined in
`protocol.cdl` and executed through EIR. Harness milestones implement runtime
services, orchestration, persistence, UI, and capability implementations; they never
restate protocol rules in Go. Prompt editions are render backends over EIR, not
maintained prompt documents.

Each implemented normative rule requires deterministic tests, positive/negative
fixtures, stable expected results, and cited protocol sections. The §§17.2 families are
not abandoned; they become Spec-track migration targets with harness support behind
them.

## Spec track

### S0 — Contracts and freeze — implemented (documents only)

Deliver the frozen source-language contract, execution contracts, and migration
discipline before any runtime code:

- `cdl.md` frozen with Errata 1;
- `eir.md`, `execution-semantics.md`, `runtime-architecture.md`;
- `migration-audit.md` with oracle freeze, transitional-code freeze, and deletion policy;
- the 24 bootstrap fixtures and their runner declared the independent parity oracle.

No runtime code exists yet. No semantic, gate, or conformance claim is implied.

### S1 — §7.7 end-to-end pilot — implemented

Additive pilot for one section, proving the contracts are implementable:

- `protocol/export-reconciliation.cdl` as frozen-grammar source (`.cdl` sources stay
  under `protocol/` until crown so they imply no authority; promoted to `protocol/`
  at S4);
- `cdl/` compiler subset: parser, resolver/typechecker, capability closure,
  omission checks, ledger verification, deterministic EIR with canonical hash, and a
  template-total prompt renderer (no runtime imports; enforced by a boundary test);
- native VM subset: staged-effect transaction model, `SET`/`GUARD`/`BLOCK`/`IF`/
  `FOR EACH`/`APPEND`, deterministic `hash.sha256`, in-memory state;
- independent reference VM consuming the same EIR;
- one prompt backend rendering the §7.7 light projection;
- synthetic execution fixtures, EIR/prompt/ledger goldens, and eight compiler negative
  fixtures including all six from CDL §14;
- test A (backend trace-hash equality) and test B (protocol parity against
  `protocol.md@4.1.2` plus the unchanged 24-fixture oracle);
- additive CLI verbs `spec compile`, `spec render`, and `spec run` over the same
  packages (no new semantic surface).

Gate **G1** and gate **G2** pass via `go test ./cdl/ ./internal/parity/`. Stop
conditions held: no grammar growth. One pilot finding corrected example data only: a
`BLOCK` target must resolve to a declared gate, so `fingerprint` became the declared
gate `artifact-fingerprint` in the example and contract.

### S2 — Coverage profile, capability contracts, context design — implemented

- classified all 148 headings (20 Parts + 128 sections) in
  `analysis/protocol-classification.json`; deterministic coverage, vocabulary,
  fingerprint, and metrics checks in `internal/analysis`, with
  `analysis/metrics.golden.json`;
- profile findings in [coverage-profile.md](coverage-profile.md): MACHINE 45.0% and
  MIXED 47.0% of section lines, pure AGENT prose 8.0%; 98 artifact sections, 33 gate and
  47 transition sections, 36 capability-dependent, 17 reuse-defining;
- context-assembly leakage quantified (9 common + 50 mode-specific read targets, 3
  strict modes, 26 iteration tokens, 221 sections cited by hand-maintained routing) with
  the required §8 migration shape;
- deterministic capability contracts defined in
  [capability-contracts.md](capability-contracts.md) with kind, determinism,
  inputs/outputs, failure modes, and current implementation status.

Gate **G3** passes via `go test ./internal/analysis/` plus review of the two documents.

### S3 — Section-by-section protocol migration — implemented (crown complete)

Migrate `protocol.cdl` in dependency order, each with compiler checks, VM support,
reference-VM comparison, and dual-source drift against `protocol.md`:

1. §4 identity, paths, canonical hashing;
2. §7.5–7.7 acquisition and export reconciliation;
3. §8 invocation identity, context, cold resume;
4. §10 iteration accounting;
5. §3 workspace registries and schema gradients;
6. §5 decisions, confirmations, staleness;
7. §15 packaging schemas, gates, Exit A/E sequencing;
8. §2 personas, entries, traversal tracks; remaining Parts.

Harness surfaces are rebased as their protocol backing migrates: context packets
resolve from EIR, workspace checks derive from declared artifacts, skill projections
become generated reference/routing views.

S3 progress:

- **Document assembly implemented:** `protocol/globals.cdl` owns document-global
  declarations once; `protocol/assembly.json` orders the sources; the compiler
  emits one EIR with cross-section resolution; `internal/spec` compiles it in-process
  for the runtime; assembly EIR and prompt goldens replace the per-section goldens.
- **§7.1–7.4 implemented (authored schema + drift):**
  `protocol/acquisition-trust.cdl` declares the helper trust-boundary rules, the
  export acquisition register schema, the virtual-coordinate rule, and the n8n/Appsmith
  inventory-kind enums; field and enum lists drift-check against protocol.md.
- **§3.1 implemented (authored schema + drift):** `protocol/preflight-registry.cdl`
  declares the preflight header, entry-cluster, unmapped-discovery, and persona-local
  discovery-buffer schemas plus the boundary and DB-autonomous-cluster rules; field
  lists drift-check against protocol.md.
- **§3.2 implemented (authored schema + drift):** `protocol/workspace-registries.cdl`
  declares the shared-entry, auth, privacy, capability, routine cross-reference,
  invariant, and shared-state record schemas plus the append-only 0G and derived 0H
  rule; `internal/parity` drift-checks every block field label against protocol.md.
- **§5.2–5.7 implemented (authored schemas + drift):**
  `protocol/evidence-and-decisions.cdl` declares the claim and block-summary,
  contradiction, decision, and confirmation record schemas plus the sanitized-projection
  and candidate-reconciliation rules; `internal/parity` drift-checks every record field
  label against protocol.md and requires §5.3/§5.6 token coverage.
- **§5.1 implemented (generated closed lists):** `internal/migration` deterministically
  extracts the seven prefix groups and the exhaustive finite-enum registry from
  protocol.md and renders `protocol/status-taxonomy.cdl` (109 enums) with a
  byte-parity regeneration test; the assembly compiles it and the ledger covers it.
- **§4.1 implemented (capability-backed):** `protocol/typed-id.cdl` declares the
  ID-prefix and iteration registries plus the normative key, path, PRF/HBK, collision,
  and registry rules; `internal/capabilities.TypedID` implements base generation and is
  tested against an independent recomputation of the §4.1 formula; `internal/parity`
  drift-checks both registries against `protocol.md`. The identity ledger is now
  multi-source (format 2) with global uniqueness across sources.
- **§4.1.1 implemented (capability-backed):**
  `protocol/semantic-payload-identity.cdl` declares the payload mandatory fields,
  semantic record version, semantic-content fingerprint, certification envelope,
  envelope mutation invariance, and confirmation binding rules; `internal/parity`
  token-drift-checks every hyphenated normative token in protocol.md §4.1.1 against the
  CDL text. The `canonical.markdown` engine remains unimplemented (see §4.1.2).
- **§4.1.2 implemented (capability-backed, bounded v1):**
  `protocol/canonical-hash-profile.cdl` declares the canonical profile steps, the
  envelope/package-member/evidence-set/transport rules, and the closed registries
  (binding kinds, envelope kinds, artifact types, row kinds, domain prefixes);
  `internal/capabilities/canonical.go` implements canonical Markdown normalization and
  the envelope, package-member, transport, and evidence-set fingerprints;
  `internal/parity` drift-checks the registries and every normative token against
  `protocol.md`.
- **§15 implemented (authored schema + drift):** `protocol/packaging.cdl`
  declares the Exit E strict conditions, universal packaging rules, all five artifact
  payload/envelope schemas, the candidate/readiness/final check registries, the
  scope-certificate tables, the acyclic sequence and hash-domain rules, and the final
  verification receipt; every schema, table shape, and check ID drift-checks against
  protocol.md. §12 implemented (authored schema + drift): `protocol/traceability.cdl`
  declares the traceability row, lifecycle-eligibility row, classification enum, the
  reciprocal-reference/ID-timing/provenance rules, and the §12.5 evidence-binding
  table; row shapes drift-check against protocol.md.
- **§2 implemented (authored schema + drift):** `protocol/personas-povs.cdl`
  declares the persona definition/purity/grammar rules, entry ownership including
  DB-autonomous, the five POVs, traversal tracks, the capability-state block and
  runtime-state matrix, and the target-decision enum; labels, matrix columns, and enum
  drift-check against protocol.md. The §3.1 persona-registry and traversal-matrix table
  shapes were added to `preflight-registry.cdl` and drift-checked.
- **§4.2–4.6 implemented (authored schema + drift):**
  `protocol/inventory-and-frontier.cdl` declares the atomic-unit and
  container-semantics rules, the source-inventory record with discovery-method and
  coverage-summary enums plus the 23 required inventory kinds, the frontier record with
  state/method enums, the coverage row with scope/disposition enums, and the
  supersession/tombstone block; schemas and enums drift-check against protocol.md.
- **§6.1–6.5 implemented (authored schema + drift):**
  `protocol/extraction-evolution.cdl` declares the atomic-component record with
  the block-confidence-rank enum, the shared-reference schema, and the dual-entry,
  promotion-lock/execution, dead-code, depromotion, and surgical-patching rules;
  schemas and enum drift-check against protocol.md.
- **§8.3–8.4, §9.3–9.6 implemented (authored schema + drift):** the invocation-mode
  enum drift-checks against the generated registry, the multi-file mode targets and
  mandatory read sets are declared as rules, and `protocol/ticket-fsm.cdl` gains
  the Human Hatch prerequisites, probe specification with safety enum, no-mock
  integrity rule, and no-mock fallback payload.
- **§13.1–13.4, §14.1–14.6 implemented (authored schema + drift):**
  `protocol/handbook-and-decisions.cdl` declares the audience roles,
  chapter/disabled-chapter/quality-gate rules, readability-review binding, decision-log
  projection, modernization prompts, rule ownership, interface baseline, equivalence
  test kinds, and the sole-input invariant; audience and test-kind enums drift-check
  against protocol.md.
- **§11.1–11.13 implemented (authored schema + drift):**
  `protocol/synthesis-catalogs.cdl` declares the synthesis gap buffer, common
  block header and certification envelope, and all ten synthesis catalogs (MOD, ER,
  REL, SM, DR, BR, UC, IF, DEP/CFG/SCHED, NFR/SEC/FLT) plus the persona profile;
  every record field list, nested table shape, and closed enum drift-checks against
  protocol.md.
- **§0.1–0.4, §1.1–1.5 implemented (authored schema + drift):**
  `protocol/foundations.cdl` declares the normative keywords, authority classes,
  predicate kinds, guarantees/non-guarantees, conceptual flow stages, the four-plane
  model with its sources/writers rules, and the non-authoritative sidecar rule; enums
  drift-check against protocol.md.
- **S3 section migration complete:** every substantive protocol section now has a CDL
  source (26 sources, 542 identities). The only remaining `## N.` matches without a CDL
  section are the §0 wrapper heading and the internal `0A-PREFLIGHT.md` template
  headings (`## 1.`–`## 5.`), whose content is covered by §3.1/§3.2 and §7.2.
- **Historical remaining inventory (superseded by completion above):** §2.1–2.6, §4.2–4.6, §6.1–6.5,
  none beyond the internal template headings.
- **Skill projections decision:** the 13 projections stay hand-maintained with the
  routing validator through S3; generation from EIR reference views is deferred to S4.
- **Collision extension implemented:** `internal/capabilities.ExtendCollision` extends
  colliding hashes to the declared 16/20/.../64 lengths and fails closed on identical
  or non-colliding input; aliases and supersession remain registry concerns (§12.1).
- **§8.5 implemented (capability-backed):** `protocol/cold-resume.cdl` declares the
  12-field assessment block and the fingerprint/validation rules;
  `internal/capabilities.ColdResumeFingerprint` implements the `COLD-RESUME|` preimage
  with carrier exclusion and fail-closed field validation; `internal/parity`
  drift-checks the field list against protocol.md §8.5.
- **§8 declarations complete:** `protocol/invocation-modes.cdl` declares the
  13-mode enum, the three strict modes, the document baseline read set, and every
  mode's specific read set (language 0.2 `BASE-READS`/`MODE`, EIR format 2);
  `cold-resume.cdl` declares and implements the §8.5 check fingerprint;
  `invocation-context.cdl` declares §8.1 identity-header fields, §8.2 strict bindings
  and ALFA behavior, §8.6 concurrency buffers and reconciliation, §8.7 stale guard,
  and §8.8 invocation-log fields. `internal/parity` drift-checks modes, strictness,
  baseline, cold-resume fields, header/log labels, and normative tokens, and proves
  read-set equivalence with the frozen Go implementation. `contextpacket` now consumes
  EIR `modes[]` and the declared iteration tokens through `internal/spec`; the hardcoded
  read sets, strict-mode switch, and iteration list are retired, closing the
  context-assembly leakage.
- **§10.1–10.6 implemented (capability-backed):**
  `protocol/coverage-and-exits.cdl` declares sweep records and traversal-cell
  states, coverage scope/status domains and arithmetic rules, the 12 Exit A
  conditions, Exit B/C rules, and the iteration-accounting rules;
  `internal/capabilities.NextIteration` advances the declared token budget
  fail-closed; `internal/parity` drift-checks the condition list, sweep fields, and
  the 26 iteration tokens against protocol.md.
- **§9.1–9.2 implemented (capability-backed):** `protocol/ticket-fsm.cdl`
  declares the nine ticket states, the ten allowed transitions, the eleven-field
  ticket schema, and the ten escalation reasons; `internal/parity` drift-checks the
  transition table, state domain, and reason enum against protocol.md.
- **§7.5–7.6 implemented (capability-backed):**
  `protocol/export-acquisition-loop.cdl` declares the five-state acquisition
  domain and the eight loop rules, with the state domain drift-checked against the
  bracket tokens in §7.5; `protocol/normalized-maps.cdl` declares the
  normalized-map record and table schema plus the navigation-only rule. §7.7 was the
  S1 pilot.

### S4 — Crown — implemented

- **Numbering removed from CDL (language 0.3):** `NUMBER` is gone; sections declare an
  ordered `PART <id> "<title>"` and optional `SUBSECTION OF <parent-id>`, and the
  compiler derives section/part numbers from declaration order and part boundaries at
  EIR/render time. Parts formerly numbered 19–20 are renumbered 17–18. Parser,
  resolver, renderer, generator, docs, and goldens are updated; numbering is now purely
  a rendering artifact.
- **Part 17/18 removed by owner decision:** the v3.17 migration plan
  (`§17.1–§17.3`) and the Migration Checklist (`§18`) were deleted from both the CDL
  sources and the oracle text; the `Protocol-Migration Metadata` preflight field and
  the §8.3 migration note were removed with them. Part numbering keeps the gap (Parts
  19–20 unchanged) so every existing cross-reference stays valid. The protocol version
  remains 4.1.2 because the frozen 24-case fixture oracle requires that exact version.
- **Parity report generator implemented:** `internal/parityreport` deterministically
  builds `analysis/parity-report.json` and `analysis/parity-report.md` over the protocol
  fingerprint, assembly sources/ledger, compiler, renderer goldens, drift-test
  inventory, and the 24-case fixture oracle; a freshness test keeps both artifacts
  regenerated (`go test ./internal/parityreport -update`). Report shows 104/104
  substantive protocol sections implemented (the §0 wrapper is a `wrapper`), all six
  checks pass, 24/24 fixtures pass, and four declared deferrals.
- **Prompt renders are per-projection:** a `PROJECTION` MUST declare `COVERS SECTION`
  (the renderer now fails closed otherwise); repository goldens carry the §7.7 pilot
  prompt as `protocol/golden/protocol.pilot.prompt.md`. The full standalone edition is
  `protocol.md`; a whole-document prompt render is not yet declared.
- **S4 follow-up — generated-edition prompt usability:** verbatim enrichment is done —
  section GOALs now carry the oracle body verbatim (prose, lists, tables, fences,
  field bullets), Parts 16/18/20 were added as sections, `PART` declarations render
  the 20 Part headings, and text blocks preserve newlines/blank lines/comments.
  `TestGeneratedProtocolCoversOracle` enforces content completeness. The standalone
  protocol-quality test on the generated edition is the remaining gate.
- **Generated protocol.md completed:** the renderer now keys declarations by section ID
  (previously dropped every rule/field/enum), renders global declarations under a
  `Document Declarations` preamble, and every one of the 108 sections carries a GOAL
  description; `TestEverySectionHasGoal` guards the invariant. Render counts: 164
  rules, 86 fields, 164 enums, 3 states, 2 tables.
- **Crown executed (S4):** the parity report was accepted by the human gate; the CDL
  sources moved to `protocol/`, the legacy text is frozen at
  `protocol/legacy/protocol-4.1.2.md` as the bootstrap parity oracle, and `protocol.md`
  plus `protocol/golden/protocol.generated.md` are generated by `cdl.RenderProtocol`
  with a freshness test. Generated section titles match the legacy headings so routing,
  context-packet, and fixture consumers keep working; Part 17.2 conformance is
  re-expressed as EIR/VM conformance plus protocol parity in `docs/conformance.md`.
- machine-generated parity report over source, compiler, stdlib, renderer, fixtures,
  and outputs, accepted by a human;
- `protocol/*.cdl` is promoted to `protocol/` and becomes the normative source;
  `protocol.md` and prompt editions become generated artifacts;
- Part 17.2 conformance is re-expressed as EIR/VM conformance plus protocol parity;
- Prompt editions are released as backend renders with fingerprints.

No crown before G1–G3 and G4 (S2/S3 review) pass.

### S5 — CDL editor tooling — planned after crown

Planned as a separate package/repository that depends only on the language toolchain
(`cdl`), never on Legacy Autopsy runtime code:

- **Syntax highlighting first:** TextMate grammar plus VS Code language configuration
  for `.cdl` (keywords, comments, `REQUIRE """..."""`, identifiers, string/scalar
  literals). This can begin during S3/S4 because the keyword surface is frozen; it must
  track grammar additions through S3.
- **Canonical printer:** expose `cdl fmt` as a deterministic, template-total AST printer
  that preserves normative text (`REQUIRE """..."""`, `TEXT`, `REASON`, `GOAL`) verbatim,
  with golden fixtures per AST form. This requires the parser to cover every migrated
  section, so it lands after S3.
- **Formatter hosting:** once the printer is total, dprint may host it through a wasm
  plugin or the Exec plugin calling `cdl fmt`. dprint is the runner, never the formatter.
- **No source authority:** the extension and formatter are non-authoritative; they never
  change semantics, and generated artifacts are still regenerated, not reformatted.

No editor-tooling work starts before S3 unless it is highlighting-only and
additive.

## Harness track

### M0 — Repository and skill foundation — implemented

Establish authority messaging, contributor rules, architecture/development/workspace/
conformance docs, the thin skill and 13 mode projections, provider-neutral Go
module/CLI foundation, routing validation, basic deterministic primitives, workspace
scaffolding, bootstrap fixtures, and honest CI reporting. Frozen as the parity oracle
per the [migration audit](migration-audit.md).

### M1 — Protocol model and routing — partially implemented, remainder in Spec track

The M0 slice already provides structural Markdown parsing, protocol version/mode
discovery, and routing validation. Remaining protocol-model depth is folded into
S2/S3; no new protocol logic is added to Go. Skill projections stay hand-maintained
until S3 decides generation.

### M2 — Deterministic IDs/path normalization — Spec-track capability contract (S2/S3)

Normalized relative paths, traversal rejection, exact Unicode/case preservation,
canonical ID keys, SHA-256 prefix and deterministic collision extension,
aliases/tombstones, and PRF/HBK coordinate foundations become declared
capabilities+spec sections in S2/S3. The harness consumes them; `internal/identity`
remains a frozen transitional primitive until then.

### M3 — Workspace initialization and structural validators — harness

Implement `doctor`, non-fabricating `init`, workspace discovery, four-plane skeletons,
persona/shared ownership boundaries, atomic file writes, `0G`/derived `0H` mechanics,
and initial structural validation. Initialization must not create covered, confirmed,
or passed semantic state. Structural checks migrate to declarations generated from EIR
as sections land.

### M4 — Local autopsy runtime foundation — harness (depends S1/S3)

Introduce the first observable application shell without claiming semantic mode
execution. The scheduler and workspace transactions execute EIR steps through the
native VM; policy remains harness-owned per `execution-semantics.md` §1:

- installable user-level OS service with a registry and startup reconciliation;
- multi-autopsy registry with persisted opaque random `autopsy_id`, repository/snapshot
  binding, rooted workspace isolation, relocation reattach, duplicate-root rejection;
- initial scheduler with one active invocation per autopsy, two globally, FIFO within
  an autopsy, round-robin across autopsies, invocation-ID idempotency/cancellation/
  recovery;
- common invocation/result boundary, exposed first through the generic skill;
- skill-driven attachment and reserved non-authoritative `.legacy-autopsy/` config;
- prohibition on UI/CLI manual semantic-work start or claim during this stage;
- transaction/post-commit hook recording invocations in `0G` and publishing revisions;
- optional `.extracted/` watcher for out-of-band changes and recovery, never authority;
- lower-priority derivative workers that cannot block extraction;
- loopback-only local API with strict `Host`/`Origin`/CSRF controls and secret-safe logs;
- browser shell (Angular, Taiga UI, Analog Vite plugin) with tabs, lifecycle/status,
  and activity stream;
- initial Cytoscape Atlas surface backed by regenerable JSON;
- questions inbox shell that cannot bypass prerequisites or establish truth.

Empty scaffolds and operational metadata must not fabricate graph semantics, coverage,
Human Hatch outcomes, Exit A, or conformance.

### M5 — Claims/source/frontier primitives — Spec-track migration + runtime

Source inventory, denominator rows, claim-level evidence, frontiers, status/path
enum dispatch, acquisition register primitives, atomic record constraints, and
zero-domain proofs are protocol semantics; they migrate in S3 with harness projections
afterward.

### M6 — Invocation context, cold resume, and executor rollout — harness (depends S2/S3)

Implement mode routing, §8.1 identity headers, mandatory read sets, minimum-sufficient
authoritative section packets, strict-scope verification, stale-checkpoint guards,
§8.5 structured cold-resume checks, and append-only invocation history. Context
assembly must resolve from EIR `EVIDENCE`/`USES`; the S2 leakage audit is a
prerequisite and remaining procedural logic is a defect. The executor proposes
semantic resume/actions; the runtime recomputes identity, scope, fingerprints, write
rights, staleness, and the cold-resume fingerprint before commit.

Execution channels in order:

1. complete the generic service-backed skill flow for any compatible coding harness;
2. thin managed adapters for named harnesses such as OpenCode, Codex, Claude Code, and
   Kiro, backed by versioned per-project configuration and secret-free metadata;
3. the minimal internal direct-model runner only after skill/adapter contracts stabilize.

Adapter configuration is operational, not protocol authority; credentials stay outside
the repository and `.extracted/`.

### M7 — Ticket FSM, reconciliation, and human dependencies — Spec-track migration + runtime

Ticket identities/transitions, closed escalation reasons, Human Hatch prerequisites,
honest gap/partial outcomes, probe ingestion, persona-local buffers, deterministic
Sequential Reconciliation, promotion/depromotion atomicity, unmapped discovery merging,
concurrency locks, and reciprocal audit effects migrate as protocol semantics; the
questions inbox is the runtime surface and dependency impact comes from EIR state.

### M8 — Coverage and Exit A — Spec-track migration + runtime

Per-kind/track/environment/snapshot coverage arithmetic, evidence-backed sweeps,
source/frontier/ticket/export closure, dynamic-caller dead-code proof,
provenance-bound applicability with fail-inclusive `Unknown`, ownership/referential
integrity, disabled-capability governance, contradiction checks, deterministic
reports, and Composite Exit A conditions migrate to `protocol.cdl`; the workbench
projects validated status and never establishes a gate.

### M9 — Export acquisition helper — operator tool

Implement the independently runnable operator-controlled Appsmith/n8n helper: dry-run
default, safe credential channels, private state outside the repository, field-aware
sanitization, secret scanning, fail-closed schema/lineage review, normalized logical
explosion, reconciliation candidates, and atomic approved publication. No production
mutation or third-party transmission. The reconciliation behavior itself is §7.7 spec
from S1/S3.

### M10 — Structured synthesis — Spec-track migration + runtime

Partial/Final-Synthesis boundaries, common semantic headers, `DERIVED-FROM` graphs,
semantic versions/fingerprints, gap buffering, evidence profiles, catalogs `90`–`96`,
and Cross-Reference-Reconciliation with targeted stale propagation migrate as spec;
Atlas expands only from committed records with inspectable provenance.

### M11 — Confirmation/staleness — Spec-track migration + runtime

Candidate-bound `CNF` records, decision approvals, payload/envelope separation,
human-only confirmation transitions, dependency impact, semantic staleness, and
invalidation/re-entry behavior. Attestations are append-only and bound to immutable
payload hashes (`execution-semantics.md` §5, §10).

### M12 — Handbook and Atlas linking — backend + runtime

HBK identities and anchors, synchronized handbook generation, progressive-disclosure
navigation, non-vacuity/link/diagram checks, readability-review coverage, upstream-first
correction, staleness, confirmation envelopes, and stable bidirectional links. The
handbook is a backend projection, not a hand-maintained document.

### M13 — Canonical packaging/hashing — Spec-track capabilities (S2/S3)

The §4.1.2 canonical hash profile, typed-record hashes, file-transport fingerprints,
carrier/envelope exclusions, canonical envelope bindings, evidence-set fingerprints,
package-member fingerprints, hash domains, post-hash instance bindings, and the five
packaging schemas become stdlib/capability contracts plus spec sections. This replaces
the deliberately limited `internal/canonical` primitive, which is never used for
protocol-significant hashing.

### M14 — Exit E — Spec-track migration + runtime

Content/decision readiness, stage-specific check registries, equivalence validation,
exact confirmations, candidate report/manifest, content-readiness report, scope
certificate, snapshot-bound signed outer manifest, cross-chain snapshot equality,
signature validation, reproducible verification, and the strict acyclic six-step
sequence. Only the validated signed outer payload may state `EXIT-E-STATUS: Passed`.

### M15 — Full Part 17 conformance — Spec-track

Every required deterministic positive/negative family is expressed as EIR/VM
conformance plus protocol parity. Diagnostics stabilize, all groups run in CI, and
applicable gates block on failures. A requirement-to-fixture audit precedes any
conformance claim. This milestone no longer means hand-written Go validators.

### M16 — First real legacy-system dry run — harness

Run a controlled, non-production pilot with pinned snapshots and approved sanitized
evidence. Exercise the workbench and all three execution stages in order, plus
acquisition, deconstruction, asynchronous questions, Atlas, Exit A, synthesis,
confirmation, handbook, packaging, and Exit E; document gaps and feed fixes back
through S3/M milestones without weakening protocol rules.

## Deferred Protocol §17.2 families and their targets

1. **Persona workspace layout and ownership** — S3 §2 migration, M3/M7 runtime.
2. **Context-qualified finite enums and registry completeness** — S3 §3/§5, M1 remainder.
3. **Invocation ownership and cold resume** — S3 §8, M6 runtime and leakage fix.
4. **Ticket escalation and honest unresolved coverage** — S3, M7/M8 runtime.
5. **Dead-code and semantic-predicate closure** — S3, M5/M8/M12 runtime.
6. **Deterministic iteration accounting** — S3 §10, M6/M7/M10 runtime.
7. **PRF identity and semantic hashing** — S3 §4/§5, M2 capability, M10/M11.
8. **HBK identity and path/anchor normalization** — S3 §4, M12.
9. **Profile Synchronization closure** — S3 §8, M6/M10/M11.
10. **Record versus artifact hashing** — S2 capability, S3 §4/§15, M13.
11. **Canonicalization and exact exclusions** — S2 capability, S3, M11/M13.
12. **Final envelope and package-member integrity** — S3 §15, M11/M13/M14.
13. **Packaging schemas, gate completeness, ordering** — S3 §15, M13/M14.
14. **Acyclic, snapshot-consistent Exit E and final verification** — S3 §15, M14.

## Resolved protocol issue: persona path/prefix drift

The original persona layout used an unconstrained persona slug while separately
assigning each persona a globally unique canonical prefix. That allowed filesystem
ownership paths to drift from invocation, ticket, and profile identities. Protocol §2.1
now binds each registered persona to an explicit persona slug and the canonical
basename `<persona-prefix>-<persona-slug>` while preserving `personas/_shared/`
unchanged. The M0 structural slice validates this basename grammar; complete
registry-to-directory agreement, persona purity, and atomic ownership moves migrate in
S3 with M3/M7 runtime support.

## Resolved protocol defect record: inverted export-reconciliation blocker

Protocol v4.1.1 §7.7 incorrectly said that 100% reconciliation blocks
`[A-STRUCTURALLY-COMPLETE]`, `[R-SWEPT]`, and Exit A, contradicting its own equation
and §7.5. `protocol.md` v4.1.2 corrects §7.7 so less than 100% reconciliation blocks
those states and gate. The corrected rule is the first Spec-track pilot section and
the negative fixture from that defect is retained permanently.

## Current next milestone

Proceed to S3 in the recommended order from the [coverage profile](coverage-profile.md):
§4.1–4.1.2, then §8 (which retires the context-assembly leakage), §10, §3+§5, §15+§12,
then remaining MIXED schema sections and AGENT-only guidance. Each migration section
requires compiler checks, VM support where applicable, a reference-VM comparison, and
dual-source drift against `protocol.md`.
