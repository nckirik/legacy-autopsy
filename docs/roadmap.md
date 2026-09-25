# Roadmap

> Non-authoritative planning projection. Milestone completion never overrides [`protocol.md`](../protocol.md), and it does not imply Protocol v4 conformance unless the required Part 19 fixtures pass. The source-language contract is frozen in [`cdl.md`](cdl.md); execution contracts are in [`eir.md`](eir.md), [`execution-semantics.md`](execution-semantics.md), and [`runtime-architecture.md`](runtime-architecture.md).

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
fixtures, stable expected results, and cited protocol sections. The §§19.2 families are
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

- `examples/spec/export-reconciliation.cdl` as frozen-grammar source (`.cdl` sources stay
  under `examples/spec/` until crown so they imply no authority; promoted to `protocol/`
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

### S2 — Coverage profile, capability contracts, context design — next

- classify `protocol.md@4.1.2` in the CDL vocabulary (rule count, token share, gate and
  transition ownership, `USES` density, capability-dependent rules, AGENT contracts,
  free prose);
- audit context-assembly leakage: can `EVIDENCE`/`USES` determine the mandatory
  read/scope set, or is procedural harness logic required?
- define deterministic capability contracts (identity, path normalization,
  canonicalization, hashing, tables) with kind, version, and failure modes;
- define proposing-provider provenance and the effect-service boundary.

Gate **G3**: coverage profile and leakage audit reviewed; capability contracts accepted.

### S3 — Section-by-section protocol migration

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

### S4 — Crown

- machine-generated parity report over source, compiler, stdlib, renderer, fixtures,
  and outputs, accepted by a human;
- `examples/spec/*.cdl` is promoted to `protocol/` and becomes the normative source;
  `protocol.md` and prompt editions become generated artifacts;
- Part 19.2 conformance is re-expressed as EIR/VM conformance plus protocol parity;
- Prompt editions are released as backend renders with fingerprints.

No crown before G1–G3 and G4 (S2/S3 review) pass.

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

### M15 — Full Part 19 conformance — Spec-track

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

## Deferred Protocol §19.2 families and their targets

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

Proceed to S2: classify `protocol.md@4.1.2` in the CDL vocabulary, audit
context-assembly leakage against EIR `EVIDENCE`/`USES`, and define the deterministic
capability contracts (identity, path normalization, canonicalization, hashing, tables)
per this roadmap and the [migration audit](migration-audit.md). Do not begin S3 or
harness work until the S2 gates are reviewed.
