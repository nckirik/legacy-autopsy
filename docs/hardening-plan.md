# Post-run hardening plan

> Non-authoritative planning document. The CDL sources under [`protocol/`](../protocol/)
> remain normative. This plan sequences the fixes derived from the first full standalone
> run (a small project, 2026-09-28/10-05) and its recorded issues PI-1..14.

## What the run proved

Working well: deterministic coverage arithmetic; a fail-closed §15.5 verifier that
recomputes rather than trusts; honest exclusion/disclosure of runtime scope; the
certification chain (`23`→`27`) as recomputable bytes; discovery, sweeps, tickets, and
SR executed with real per-cell records.

Failing or weak: the handbook was 17 near-empty metadata shells and readability CNFs
were relayed without a human read (PI-11, PI-12); operator surfaces were impractical
(1275-line decision digest, tag-dense reports - PI-6); 91/91 record fingerprints were
hand-stamped and only caught by a newly built verifier; validators had blind spots
(no probe allow-list, enum checks that short-circuited); static-only eligibility was
ambiguous at Exit E (PI-5, PI-10); fail-inclusive was over-applied to recoverable
artifact defects (PI-9); several metadata ambiguities cost agents time (PI-1..4, PI-7,
PI-8); and Markdown-as-database remained the largest source of read/write friction
(M3.5 in [roadmap](roadmap.md), [sqlite-store.md](sqlite-store.md)).

The session export also isolates a more fundamental failure mode: the executor carried
too much of the run in one long-lived model context. Most of the protocol state,
mechanical bookkeeping, semantic work, and accumulated history were repeatedly
materialized together, so the executor became progressively less able to see the run
from outside its own local trajectory. This is a context-management and execution
architecture failure, not a model-quality finding. An independent reviewer does not
need a stronger or different model; the same model in a fresh, deliberately different
context is already an independent semantic viewpoint.

## Run cost and behavior (session export)

The session export (`legacy-deconstruction-protocol-execution.json`) makes the cost
concrete: 1,026 assistant steps and 982 tool calls (735 bash, 83 edit, 71 read, 46
write), 14 context compactions, 136.4M input tokens, and 71 operator messages spread
over 160.5 elapsed hours (2026-09-28 to 2026-10-05, including a power outage) - for a
small repository. Only 5 tool calls errored; the failures were semantic, quality, and
process failures, not tooling failures.

The transcript also confirms the behavioral gaps: the removal of a bulk `[R-SWEPT]`
bookmarking pass after challenge, the decorative-fingerprint discovery, the
under-reaching on a recoverable schema/FDW defect, handler stubs, and relayed
readability. Several operator messages are nudges ("why did you stop?", repeated
continuation prompts, one prompt sent three times) and the decision gate ended in a
wholesale approval rather than per-batch review.

## Architectural conclusion from the run

The post-run direction is now explicit: **DSL + native harness first**. The normative
CDL assembly is the protocol program; the compiler/EIR and Legacy Autopsy runtime own
deterministic execution mechanics; model backends are bounded semantic workers rather
than long-lived interpreters of the whole run.

For managed execution:

- each semantic work item is a separate invocation with an explicit scope, inputs,
  expected result shape, and bounded context projection;
- the authoritative dependency/read closure is validated and fingerprinted by the
  runtime, but it is not synonymous with "put every byte in the model prompt";
- the context builder materializes only the smallest complete semantic projection for
  that work item, with deterministic handles/fetches for additional evidence when
  needed;
- conversation history is disposable; no semantic invocation depends on the executor
  remembering previous chat turns;
- executors submit semantic deltas and rationale, never canonical fingerprints,
  coverage arithmetic, package mechanics, or state-transition bookkeeping;
- review is a second semantic invocation over committed state with an independently
  constructed context. It MAY use the same model as the worker; independence comes
  from context and role separation, not model identity;
- executor uncertainty is not human authority. Protocol mechanics return to the
  runtime/spec, harness defects block, and uncertain semantic work is retried or
  independently reviewed before a Human Hatch is considered.

This matches the target ownership already described in
[architecture.md](architecture.md), [runtime.md](runtime.md), and
[skill-runtime.md](skill-runtime.md): Legacy Autopsy owns orchestration and bounded
context construction; executors own one semantic work unit at a time.

## Acceptance metrics for W6

- no nudges: operator interaction only at declared human gates;
- each reconstruction-affecting decision reviewed in a readable batch - no wholesale
  "all approved";
- at most one compaction, no context-loss interruptions;
- tool-call count for the same project under a quarter of this run;
- zero decorative fingerprints by construction (one recomputing substrate shared by
  emitter and verifier);
- handbook chapters meet the substance floor and are actually read, with complete CNF
  fields;
- recoverable artifact defects resolved without operator prompting (PI-9 applied).
- semantic context size is bounded by the work item's dependency closure, not by run
  age; no invocation requires accumulated chat history or a full-run transcript;
- each semantic work item can be replayed cold from committed state and its declared
  inputs, with no hidden conversational prerequisites;
- reviewer invocations receive an independently built context and can challenge the
  worker without inheriting its reasoning trajectory;
- Human Hatch traffic contains only genuine human-authority/domain questions; protocol
  interpretation, record mechanics, and executor uncertainty never become operator
  decisions.

## Principles (unchanged)

- `Unknown` stays fail-inclusive; nothing here weakens evidence or authority rules.
- Deterministic mechanics live in code; semantics stay with agents and humans.
- Runtime, probes, and live databases remain optional; static-only is a first-class
  completion path through both gates.
- Human confirmations, exclusions, and signatures remain human acts.
- Model context is a per-work-item projection, not a persistence layer. Authoritative
  state lives in validated records/runtime state, never conversation history.
- Semantic independence is context independence: worker and reviewer MAY use the same
  model, but never the same reasoning context for the same adjudication.
- Human authority is not a fallback for protocol or harness uncertainty.

## Workstreams

Status markers: `[x]` done, `[~]` started, `[ ]` todo.

Overall: **W1 done; W2/W3 in progress; W4-W6 pending.**

### W1 - Protocol text revision (v4.2, minor) - done

- [x] PI-1..PI-14 accepted and recorded with decisions.
- [x] Clarifications batch 1: dormant/commented-out units, bootstrap read set,
      cold-read scope, closure vs prompt materialization, invocation-ID grammar,
      fail-inclusive boundary, tooling location, static-only eligibility, Human Hatch
      blocker classification, operator presentation.
- [x] Quality-gate batch 2: handbook substance and readability rules, CNF readability
      field and relay ban, static-only Exit E eligibility.
- [x] PI-4 single generated version literal: one `{{PROTOCOL-VERSION}}` render token
      substituted from `ProtocolVersion` at EIR build; renderer, 0A template, §8.1
      example, and footer share the one `4.2` form (no `v` prefix, no duplicate footer).
- [x] Current-edition oracle snapshot `protocol/legacy/protocol-4.2.md` with freshness
      test; 4.1.3 kept as history.
- [x] Revision-marker guard test keeps accepted clarifications in the generated edition.
- [x] Routing/skill manifest aligned to 4.2; parity report regenerated.
- [ ] Harness-side enforcement of the PI-11 substance lint and PI-12 CNF schema
      (lands with W2/W5).
- [ ] Storage-binding abstraction (declared canonical binding, per-binding canonical
      serialization/fingerprint/package rules): proposed PI-15, candidate v4.3; the
      SQLite store stays a development artifact until then.
- [ ] Relational declaration semantics in CDL/EIR (`STORAGE`, `PARENT`, `KEY`, `REF`;
      language `cdl/0.4`, EIR format 4) per [relational-mapping.md](relational-mapping.md);
      table-per-`FIELD` remains explicitly provisional until this lands.

### W2 - Native runtime and certification engine

- [ ] Execute the deterministic MACHINE slice from compiled CDL/EIR wherever
      implemented; unsupported mechanics fail explicitly instead of being improvised by
      the model.
- [x] Record validation: `internal/certify.ValidateRecord` (declared fields only,
      required presence, integer and string-set types, closed enum domains).
- [~] Canonical fingerprint substrate: content, envelope, evidence-set,
  package-member, and transport fingerprints exist in `internal/capabilities`; not yet
  unified into one emitter/verifier authority.
- [ ] Registry-driven validity at write time and verify time (enums, statuses, kinds).
- [ ] Schema-conformant emitters for `23`/`24`/`25`/`26`/`27` and the §15.5 verifier;
      golden tests for emit-twice byte equality, tamper detection, fail-closed results.
- [ ] Validator allow-lists (for example `.extracted/probes/[TICKET-*].log`) and
      append-only post-conditions.
- [ ] Typed semantic deltas: IDs, canonical serialization, fingerprints, coverage
      arithmetic, transitions, packaging, and envelopes are runtime-owned.

### W3 - SQLite canonical store, semantic work units, context builder, review

- [~] Schema generation from EIR (`internal/store.Schema`): one table per FIELD block,
  enum `CHECK`s, provenance, append-only `0G` triggers. Provisional bootstrap: the
  FIELD-to-table semantics item in W1 blocks treating this as normative.
- [x] Apply without Go dependencies (`store.Apply`, `spec schema --apply`) and
      `spec schema --out`; atomic: DDL runs in one transaction and a mid-script failure
      leaves zero tables (negative test).
- [x] Canonical dump and fingerprint (`store.CanonicalDump`, `spec store`): carries the
      generated schema plus rows, verifies the live schema against EIR before emitting,
      content-only fingerprint (provenance excluded), restorable and drift-rejecting
      tests.
- [x] `set<string>` canonical form: ascending unique; duplicates and non-string
      elements rejected (`certify.CanonicalStringSet`), order owned by the runtime.
- [ ] Import/export between Markdown records and the store (pilot one record class,
      then generalize; round-trip fixtures).
- [ ] Semantic capability surface for workers (ID lookup, search, relation/evidence
      traversal, semantic proposals, control actions) over relational transactions.
- [ ] Automatic read-set provenance: every read extends the invocation's exact
      fingerprinted dependency set; commit stale-checks it.
- [ ] One-work-item scheduler/claim with mode, scope, dependency, result-contract, and
      allowed-effect declarations.
- [ ] Context builder separating validated dependency closure from materialized model
      projection (smallest complete projection plus deterministic fetches).
- [ ] Cold replayability of every invocation from committed state; no chat-history
      prerequisites.
- [ ] Reviewer/adjudicator invocation with an independently constructed context (same
      model allowed).
- [ ] Blocker classification before Human Hatch (only genuine human/domain questions
      reach the operator).
- [ ] Per-work-item instrumentation: context bytes/tokens, dependency count, fetch
      expansion, retries, review outcomes, operator escalations.
- [ ] Runtime question before each run: review policy (which classes require
      independent review) and operator-presentation minimums.

### W4 - Storage conformance and projections

- [ ] Registry constraints beyond enums: unique scope tuples, references as keys,
      canonical registry iteration order.
- [ ] Append-only `0G` enforced as a runtime post-condition, not only a trigger.
- [ ] Canonical dump/hash/package rules wired into the packaging emitters.
- [ ] Round-trip migration fixtures (Markdown → store → canonical dump → compare).
- [ ] Generated human/audit projections from the store.
- [ ] Administrative `la` helpers kept outside the semantic model contract.

### W5 - Handbook and human surfaces

- [ ] Narrative handbook renderer from `90`-`96` synthesis (audience purpose, observed
      legacy vs target decisions, examples/diagrams, navigation, volume tied to records).
- [ ] PI-11 substance lint: required sections plus at least eight claim-backed
      assertions per chapter, reason-bearing `Not-Applicable` allowed.
- [ ] Real readability workflow with complete CNF fields per PI-12 (no relay-only
      attestation).
- [ ] Operator renderer per PI-6: plain-language questions with options and
      consequences; one-line-per-decision indexes; first screen says what changed, what is
      blocked, and what is needed.

### W6 - Acceptance rerun

- [ ] Rerun the same small project on v4.2 through the bounded semantic-work harness.
- [ ] Acceptance criteria: zero illegal tokens by construction; every fingerprint
      recomputed; handbook meets the substance floor and is actually read; decision batch
      reviewable in the target time; static-only run reaches both gates with disclosed
      residuals; §15.5 verifier passes on the signed bundle and fails on tamper.
- [ ] Comparative metrics against the exported baseline: token/context growth,
      compactions, tool calls, semantic retries, reviewer catches, Human Hatch questions,
      operator nudges, wall clock.
- [ ] Independent evidence sampling plus a separate developer-usefulness review.
- [ ] Remove the README `TEMP` notice only when these metrics pass.

## Sequencing and gates

1. W1 done; W1 semantics feed W2/W3.
2. W2 and W3 are the critical path and run in parallel: native mechanics, the
   canonical store, and bounded semantic execution must exist before W6.
3. W4 hardens migration/conformance/projections after the core SQLite/runtime contract
   exists; the rerun does not use Markdown as working state.
4. W5 must land before W6 so the rerun measures execution quality and handoff quality.
5. W6 produces the comparative metrics against the exported first-run baseline.

## Decisions

1. [x] Revision mechanics: v4.2, snapshot each accepted revision as the current oracle
       (old oracles kept).
2. [x] Handbook substance floor: required sections plus at least eight claim-backed
       assertions per chapter, reason-bearing `N/A` allowed.
3. [x] Invocation-ID grammar: `ITERATION-PERSONA-POV-SEQUENCE`.
4. [x] Review policy and operator-presentation minimums: asked as a runtime
       configuration question before each run (implementation in W3).
5. [x] `set<string>` canonical form: ascending unique; duplicates rejected, order
       runtime-owned. `list<string>` preserves order; structured collections use owned
       child rows, not list columns.
6. [~] Canonical state versus physical form: declare a canonical storage binding rather
   than mandating SQLite or fixing Markdown (`markdown/1` stays the default). The
   normative revision is PI-15; the store remains provisional until PI-15 and the
   relational declaration semantics land (W1).

## Temporary usage notice

Until W2/W3 land and the W6 acceptance metrics pass, the README carries a temporary
"not for real autopsies yet" warning (marked `TEMP`). Remove it only when the
acceptance metrics in this plan pass on the rerun.

## Next up

First settle the W1 language/normative items that unblock mechanical storage: the
relational declaration semantics (`STORAGE`/`PARENT`/`KEY`/`REF`) and the PI-15 storage
binding. Then W3 import/export (Markdown ↔ store pilot, frozen until the binding
question is answered) and W2 evidence-set recomputation plus the package-artifact
emitters; then the semantic capability surface and the runtime configuration question.
