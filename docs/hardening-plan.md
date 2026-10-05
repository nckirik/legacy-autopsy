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

## Acceptance metrics for W5

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

### W1 - Protocol text revision (v4.2, minor)

Turn accepted PIs into CDL changes with fixtures, drift tests, a fresh oracle snapshot,
and a regenerated edition.

- **Clarifications:** PI-1 commented-out/dormant units are in-scope Disabled evidence;
  PI-2 Preflight bootstrap read set; PI-3 invocation-ID grammar or explicit delegation;
  PI-4 single generated version literal; PI-5 static-only eligibility through both
  gates; PI-7 cold-read scope (`0H` + cited `0G` range); PI-8 tooling location and
  non-membership; PI-9 fail-inclusive stops at absent authority; PI-10 Exit E inherits
  the disclosed scope and test doubles are scaffolding only.
- **Quality gates:** PI-11 handbook substance criteria; PI-12 readability confirmation
  schema and relay ban; PI-6 minimum operator-presentation rules for questions and
  reports (projections only, no schema change).
- **Execution/authority boundaries:** PI-13 Human Hatch may carry only genuine
  human-authority/domain uncertainty, never executor/protocol uncertainty; PI-14
  mandatory read/dependency closure is distinct from prompt materialization so bounded
  semantic invocations remain conforming.

### W2 - Native runtime and certification engine

Make the run's ad-hoc mechanics native DSL/runtime behavior rather than work delegated
to the semantic executor.

- Execute the deterministic MACHINE slice from compiled CDL/EIR wherever implemented;
  unsupported mechanics fail explicitly instead of being improvised by the model.
- Canonical fingerprint substrate for every record class (semantic content, DEC
  content, dependency-set, evidence-set, envelope, package-member, transport), computed
  once and reused by emitter and verifier.
- Registry-driven validity: enums/statuses/kinds from the §5.1 registry and `§4.1`
  type grammar, checked at write time and rechecked at verify time.
- Schema-conformant emitters for `23`/`24`/`25`/`26`/`27` and the §15.5 verifier;
  golden tests for emit-twice byte equality, tamper detection, and fail-closed results.
- Validator allow-lists (for example `.extracted/probes/[TICKET-*].log`) and explicit
  post-conditions for append-only logs.
- The semantic executor returns a typed semantic delta plus evidence/rationale. IDs,
  canonical serialization, fingerprints, dependency/evidence-set hashes, coverage
  arithmetic, transitions, packaging, and certification envelopes are runtime-owned
  effects and cannot be hand-stamped by a model.

### W3 - Semantic work units, context builder, and independent review

Turn the existing bounded-invocation target architecture into the primary execution
path before another acceptance run.

- Scheduler claims exactly one semantic work item at a time. A work item declares mode,
  semantic scope, source/evidence dependencies, expected result contract, and allowed
  effects; it never silently rolls into another cluster/POV/mode.
- Context construction separates **validated dependency closure** from **materialized
  model context**. The runtime fingerprints the full required closure while rendering
  only the smallest complete projection needed for the task, plus deterministic
  on-demand access to omitted evidence.
- Every invocation is cold-replayable from committed state. Run age, prior chat turns,
  compaction summaries, and model memory are not inputs.
- Add a reviewer/adjudicator invocation for uncertain or high-impact semantic results.
  The reviewer receives committed evidence and the worker result in a fresh,
  independently constructed context. The same model MAY fill both roles; context
  independence is the required separation.
- Before Human Hatch escalation, classify the blocker as human authority/domain fact,
  unresolved source semantics, protocol mechanics, harness defect, or executor
  uncertainty. Only the first category, and genuine external/domain facts where the
  protocol requires a human, reach the operator.
- Instrument context bytes/tokens, dependency count, fetch expansion, retries, review
  outcomes, and operator escalations per work item so context pathologies are visible
  rather than inferred after a week-long run.

### W4 - Records and storage binding

Reduce read/write friction with the SQLite binding from [sqlite-store.md](sqlite-store.md):
one table per record class, constraints from the registries, append-only `0G`,
canonical dump for hashing/packaging, and `la` helpers. Markdown remains an export for
review, generated, never hand-edited.

### W5 - Handbook and human surfaces

- A handbook renderer that produces narrative chapters from `90`-`96` synthesis:
  audience purpose, observed legacy versus target decisions, representative
  examples/diagrams, navigation, and volume tied to the underlying records.
- A deterministic substance lint per PI-11, plus a real human readability workflow with
  complete CNF fields per PI-12.
- An operator renderer per PI-6: plain-language questions with options and
  consequences; one-line-per-decision indexes; reports whose first screen answers what
  changed, what is blocked, and what is needed.

### W6 - Acceptance rerun

Run the same small project again with v4.2 through the bounded semantic-work
harness, not as one long-lived protocol-execution conversation.

- Acceptance criteria: zero illegal tokens by construction; every fingerprint
  recomputed (no decorative values); handbook chapters meet the substance floor and a
  human actually reviews them; a decision batch is reviewable in the target time; a
  static-only run reaches both gates with disclosed residuals and no runtime
  requirement; the §15.5 verifier passes on the signed bundle and fails on tamper.
- Compare against the exported baseline: token/context growth, compactions, tool calls,
  semantic retries, reviewer catches, Human Hatch questions, operator nudges, and wall
  clock. A run that produces a better handbook but still needs full-run contexts and
  repeated human protocol debugging is not accepted.
- Sample reconstruction claims independently for evidence fidelity, and perform a
  separate developer-usefulness review of the handbook. Structural certification,
  evidence correctness, and reconstruction usefulness are distinct acceptance layers.

## Sequencing and gates

1. Owner decisions below; then W1 changes accepted PI-by-PI with fixtures.
2. W2 and W3 are the critical path and may start against accepted W1 semantics in
   parallel: native mechanics plus bounded semantic execution must exist before W6.
3. W4 storage binding may follow the runtime interfaces. The first bounded-context
   acceptance pass does not need to wait for SQLite if the harness already owns
   canonical reads/writes and can prove the same contracts over Markdown.
4. W5 must land before W6 so the rerun measures both execution quality and handoff
   quality.
5. W6 produces the comparative metrics against the exported first-run baseline.

## Decisions needed from the owner

1. Revision mechanics: bump to **v4.2**, and snapshot each accepted revision as a new
   frozen oracle (keeping previous oracles; drift tests target the current one).
2. Handbook substance floor baseline (for example: required chapter sections plus a
   minimum claim-backed assertion count per chapter, with reason-bearing N/A allowed).
3. Operator-surface minimums and a target review time for a 50-decision batch.
4. Which semantic result classes require independent review by default versus review
   only on uncertainty/risk. Model identity is deliberately not part of this policy.

## Immediate next step

Accept or reject PI-1..14 one by one (updating their status), then start W1 with the
accepted set. The storage and handbook workstreams are already sketched in
[sqlite-store.md](sqlite-store.md) and roadmap M3.5.
