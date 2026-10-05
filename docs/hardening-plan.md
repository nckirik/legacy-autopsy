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

Overall: **W1 and W3 done; W2 in progress; W4-W6 pending.**

### W1 - Protocol text revision (v4.2 + v4.2.1 patch) - done

- [x] PI-1..PI-15 accepted and recorded with decisions.
- [x] Clarifications batch 1: dormant/commented-out units, bootstrap read set,
      cold-read scope, closure vs prompt materialization, invocation-ID grammar,
      fail-inclusive boundary, tooling location, static-only eligibility, Human Hatch
      blocker classification, operator presentation.
- [x] Quality-gate batch 2: handbook substance and readability rules, CNF readability
      field and relay ban, static-only Exit E eligibility.
- [x] PI-4 single generated version literal: one `{{PROTOCOL-VERSION}}` render token
      substituted from `ProtocolVersion` at EIR build; renderer, 0A template, §8.1
      example, and footer share one edition literal (no `v` prefix, no duplicate footer).
- [x] Current-edition oracle snapshot `protocol/legacy/protocol-4.2.1.md` with freshness
      test; `protocol-4.2.md` and 4.1.3 kept as history.
- [x] Revision-marker guard test keeps accepted clarifications in the generated edition.
- [x] Routing/skill manifest aligned to 4.2; parity report regenerated.
- [x] Harness-side enforcement of the PI-11 substance lint and PI-12 CNF schema
      (`internal/handbook`: required sections, claim-backed assertion floor with
      reason-bearing N/A, readability-review schema with scope/limitations/timestamp/HBK
      bindings).
- [x] Storage-binding abstraction (declared canonical binding, per-binding canonical
      serialization/fingerprint/package rules) per PI-15, implemented in the v4.2.1
      patch: §1.5, §3.1 `0A` metadata, §4.1.2, §15.3, §17.1-17.2; harness enforcement
      follows with the store work. The SQLite store stays a development artifact until
      the relational declaration semantics land.
- [x] Relational declaration semantics in CDL/EIR (`STORAGE`, `PARENT`, `KEY`, `REF`;
      language `cdl/0.4`, EIR format 4) per [relational-mapping.md](relational-mapping.md),
      with all 86 FIELD blocks classified and the store generator rewired to schema
      version 2. Deferred: `child` table projection and `REF` foreign keys, with the
      import/export work.

### W2 - Native runtime and certification engine

- [x] Deterministic MACHINE slice executes from compiled CDL/EIR through
      `internal/runtime` with fixture-supplied agent outputs; unsupported mechanics
      produce explicit diagnostics. Semantic (agent) work routes through `work` units.
- [x] Record validation: `internal/certify.ValidateRecord` (declared fields only,
      required presence, integer and string-set types, closed enum domains).
- [x] Canonical fingerprint substrate unified in `internal/capabilities` (content,
      envelope, evidence-set, package-member, transport, hash-domain, instance binding);
      `internal/artifact` emit/verify paths and `work` commits use it as the single
      recomputation authority.
- [x] Registry-driven validity at write time: enum-typed values validate against the
      compiled §5.1 enums in `certify`, enum columns get schema `CHECK`s, and ingest
      rejects unknown values (`TestIngestRejectsUnknownEnum`).
- [~] Artifact validation and packaging (`internal/artifact`): declared field sequence
  and order validation, canonical payload fingerprints with carrier exclusion and
  path binding, artifact-instance bindings, sorted package-member manifests, and
  fail-closed member verification (tamper, missing, extra, order-independent).
  Bundle verification (`artifact.VerifyBundle`, `spec artifact validate|verify`):
  outer payload validity, declared member set vs actual set, fingerprint
  recomputation, snapshot equality, tamper/missing/extra rejection. Remaining: the
  full §15.4 chain references and final check-registry fingerprints, plus assembly of
  the `23`-`27` human deliverables. Architecture (decision 8): payloads are
  executor-authored Markdown validated through the label-to-column mapping; no label
  table is introduced.
- [x] Append-only runtime post-conditions: `store.LogState` plus
      `store.VerifyMonotonicLogs` fail closed on any disappearing log record; the
      database triggers remain as a second guard. Validator allow-list enforcement
      (for example `.extracted/probes/[TICKET-*].log`) follows with the probe workflow
      validator.
- [~] Typed semantic deltas: IDs (runtime-assigned from the unit target), canonical
  payload serialization, fingerprints, and artifact member bindings are
  runtime-owned. Coverage arithmetic, FSM transitions, and envelope emission follow
  with the gate/verifier wiring.

### W3 - SQLite canonical store, semantic work units, context builder, review

- [x] Schema generation from EIR (`internal/store.Schema`, version 2): tables from
      declared `STORAGE` classes (48 tables today), enum `CHECK`s, declared `KEY`s as
      `UNIQUE` natural keys, provenance, append-only `log` triggers; unknown or missing
      classes fail closed.
- [x] Apply without Go dependencies (`store.Apply`, `spec schema --apply`) and
      `spec schema --out`; atomic: DDL runs in one transaction and a mid-script failure
      leaves zero tables (negative test).
- [x] Canonical dump and fingerprint (`store.CanonicalDump`, `spec store`): carries the
      generated schema plus rows, verifies the live schema against EIR before emitting,
      content-only fingerprint (provenance excluded), restorable and drift-rejecting
      tests.
- [x] `set<string>` canonical form: ascending unique; duplicates and non-string
      elements rejected (`certify.CanonicalStringSet`), order owned by the runtime.
- [x] Record ingestion and access (`store.Ingest`, `Get`, `List`, `Render`,
      `Inspect`, `spec record put|get|list|render|inspect`): canonical payloads stored
      verbatim, columns derived by the deterministic label-to-column mapping, fingerprints
      recomputed, human renders clean of mechanical metadata, debug renders marked
      non-authoritative, all failures closed. Bootstrap `0A` import (`ImportPreflight`)
      ingests the preflight header (including the declared storage binding, defaulting to
      `markdown/1`) and every persona registry row, with a fixture round-trip.
- [x] Semantic capability surface (`internal/work.Caps`): payload get, record list,
      payload search, with automatic read-set binding; no raw SQL write path. Semantic
      proposals are the executor return value; review/escalation are runtime-routed
      control outcomes, not model-called setters.
- [x] Automatic read-set provenance: every `Caps.Get` binds the observed
      fingerprint; `work.Commit` stale-checks the accumulated set (excluding the
      proposal target) before ingest and rejects with `StaleError`.
- [x] One-work-item scheduler (`work.Session.Run`): sequential single-writer claim
      with per-unit mode, scope, declared reads, result contract, allowed effects, and
      `RequiresReview`; `Unit.Validate` fails closed against the schema.
- [x] Context builder (`work.BuildContext`): the declared closure is validated and
      fingerprinted, and only that closure is materialized; undeclared records never
      enter the context.
- [x] Cold replayability: invocations are replayable from the committed store plus
      the declared unit (payload, closure, read set are all store-derived; no
      conversation state).
- [x] Reviewer/adjudicator (`work.Reviewer`): independently constructed context
      (fresh capability surface, proposed payload plus dependencies only), policy-driven
      by `Config.ReviewMode`; reject retries the unit.
- [x] Blocker classification (`work.ClassifyBlocker`): closed classes route human
      authority to the operator; uncertainty retries; protocol mechanics and harness
      defects fail closed.
- [x] Per-work-item instrumentation (`work.UnitStats`/`Report`): context bytes,
      dependency count, fetch count, retries, review outcome, escalation, blocker class.
- [x] Runtime question before each run (`work.RunQuestions`/`AskConfig`): review
      mode plus operator-presentation minimums, defaults per PI-6/PI-13.

### W4 - Storage conformance and projections

- [~] Registry constraints: canonical registry iteration order is deterministic
  (tables sorted by name, rows by record ID) and declared `KEY`s become `UNIQUE`
  natural keys; scope-tuple uniqueness and reference-as-key constraints land when
  the scope-tuple fields are modeled as declared keys.
- [x] Append-only `0G` as a runtime post-condition: `store.LogState` snapshots log
      record IDs and `store.VerifyMonotonicLogs` fails closed on any disappearance, in
      addition to the database triggers.
- [~] Canonical dump and package-member fingerprints exist and are verified against
  EIR-derived schemas; wiring the dump into the handoff bundle and outer-manifest
  chain follows with the §15.5 verifier.
- [x] Round-trip fixture: bootstrap `0A` import → canonical dump is identical for
      different import orders (`TestBootstrapRoundTripDeterministic`).
- [x] Generated human projections (`store.Render`) and diagnostic renders
      (`store.Inspect`, canonical dump for audit) from the store.
- [x] Administrative helpers stay outside the semantic contract: `spec schema`,
      `spec record put|get|list|render|inspect`, `spec store`, `spec work` (no raw SQL
      path for executors).

### W5 - Handbook and human surfaces

- [x] Narrative handbook assembler (`internal/handbook.Assemble`): required-section
      layout (audience purpose, observed legacy behavior, target decisions, navigation)
      with embedded synthesis records; narrative prose stays model-authored and the
      assembler never invents content. Volume is tied to bound records.
- [x] PI-11 substance lint (`internal/handbook.LintChapter`): required sections plus
      at least eight claim-backed assertions per chapter, reason-bearing
      `Not-Applicable` allowed; PI-12 CNF schema validator
      (`ValidateReadabilityReview`) requires scope, limitations, timestamp, and exact
      HBK bindings for Pass/Fail.
- [x] Readability workflow (`internal/handbook.VerifyReadabilityCoverage`): combined
      Pass/Fail reviews must bind every candidate chapter exactly once with matching
      candidate fingerprints; schema requires scope, limitations, timestamp, and exact
      HBK bindings, so relay-only attestation fails closed.
- [x] Operator renderer (`internal/render`): decision batches with a
      changed/blocked/needed first screen and one-line-per-item index, plus
      plain-language questions with options and consequences. Runtime supplies
      classifications and payloads; the formatter is display-only.

### Semantic mode pipelines (post-W3)

The harness can now generate and run bounded work units for a real mode. First
pipeline delivered: Preflight (`internal/preflight`, `spec preflight scan|run`),
piloted on `veri-yonetimi` - see [preflight-pilot.md](preflight-pilot.md).

- [x] Preflight: deterministic bounded scan (languages, entries, evidence excerpts),
      `0A` header unit, persona registry batch with runtime-derived IDs, routed
      normative sections, per-unit reasoning policy, retry-with-feedback.
- [ ] Preflight entry clusters and acquisition registry (needs a declared cluster ID
      field).
- [ ] Discovery: source inventory and traversal frontier unit generation.
- [ ] Assurance: claims, coverage, reconciliation unit generation.
- [ ] Synthesis: `90`-`96` semantic records with independent review.
- [ ] Handoff: artifact authoring, packaging, and the §15.5 verifier chain.

### W6 - Acceptance rerun

- [ ] Rerun the same small project on v4.2.1 through the bounded semantic-work harness.
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

1. [x] Revision mechanics: v4.2 then the v4.2.1 patch, snapshot each accepted revision
       as the current oracle (old oracles kept).
2. [x] Handbook substance floor: required sections plus at least eight claim-backed
       assertions per chapter, reason-bearing `N/A` allowed.
3. [x] Invocation-ID grammar: `ITERATION-PERSONA-POV-SEQUENCE`.
4. [x] Review policy and operator-presentation minimums: asked as a runtime
       configuration question before each run (implementation in W3).
5. [x] `set<string>` canonical form: ascending unique; duplicates rejected, order
       runtime-owned. `list<string>` preserves order; structured collections use owned
       child rows, not list columns.
6. [x] Canonical state versus physical form: PI-15 implemented in v4.2.1; declared
       canonical storage binding (`markdown/1` stays the default), bindings agree on
       identity/payload digests, unsupported bindings fail closed. The store remains
       provisional until the relational declaration semantics land (W1).
7. [x] Artifact emission model: executors author artifact payloads in protocol
       Markdown; the runtime validates field sequence and order through the
       label-to-column mapping, computes fingerprints and member bindings, and verifies
       by recomputation. No display-label table is introduced.
8. [x] Markdown surface: under a non-Markdown binding, no operational Markdown is
       maintained; Markdown is only human-readable semantic output rendered on demand
       (handbook, decision batches, review/gate/audit views, operator questions).
       Mechanical records are store-only and render to Markdown only in debug mode with
       a non-authoritative marker. `0A` is the bootstrap exception because it declares
       the binding.

## Temporary usage notice

Until W2/W3 land and the W6 acceptance metrics pass, the README carries a temporary
"not for real autopsies yet" warning (marked `TEMP`). Remove it only when the
acceptance metrics in this plan pass on the rerun.

## Next up

W3 is complete at the deterministic substrate level (`internal/work`), and the
OpenAI-compatible adapter (`internal/llm`, `spec work ping|run`) now drives real
executors/reviewers through it. Remaining W3-family integration: child-table and
foreign-key projection. Next: W2 evidence-set recomputation and the package-artifact
emitters (`23`-`27`) plus the §15.5 verifier, then W4 conformance and W5
handbook/operator surfaces.
