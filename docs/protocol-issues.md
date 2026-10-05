# Proposed protocol issues

> Non-authoritative backlog. The CDL sources under [`protocol/`](../protocol/) are
> normative; nothing here overrides them. Issues are recorded before any protocol
> change, evaluated, and only then implemented in CDL with the required version bump
> and regenerated artifacts. Protocol sources are never edited while a standalone run
> is in progress.

Post-crown revision process (to settle before the first accepted change): accepted
issues land in `protocol/*.cdl`, regenerate `protocol.md` and goldens, and bump the
protocol version. The frozen `protocol/legacy/protocol-4.1.3.md` remains the bootstrap
parity oracle; whether revisions keep drift-checking against it or require a new
oracle snapshot is an open process decision.

## PI-1: Commented-out code and dormant units are not named as in-scope evidence

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (three loaded `Q_*.sql` units whose
  only call sites are commented out were left without inventory/coverage rows and
  parked as pending-mirror discoveries)
- **Affected sections:** §2.4 traversal tracks, §2.5 capability and runtime state,
  §4.3 source inventory, §6.3 staging/promotion (`S-DEAD-CODE`), §10/§16 coverage
- **Observed text:** "`Disabled`: explicitly dormant, disabled, disconnected,
  legacy-retained, or configured off behavior." (§2.4); "In-scope disabled units
  require the same atomic static coverage as active units and block Exit A when
  undisposed." (§2.5); "`[S-DEAD-CODE]` requires an inventory-backed zero-caller proof
  across every traversal track and every environment/snapshot in the component's bound
  scope denominator." (§6.3)
- **Issue:** The protocol defines `Disabled` behavior and requires coverage for it, but
  it never names commented-out code, commented call sites, or loaded-but-never-executed
  units as forensic evidence to extract. Implementations can therefore read commented
  wiring as dead text or out-of-scope, leaving referenced units unmapped and out of the
  denominator.
- **Proposed direction:** state explicitly that commented-out wiring is evidence of
  dormant/legacy-retained behavior; the units it references are in scope and MUST be
  extracted with `Traversal Track: Disabled` (or `Shadow/Conditional` when an
  authorized human identifies an intended paused/re-enableable feature); comment status
  never proves dead code, an exclusion, or a zero-caller result; the unit keeps the same
  atomic static coverage as active units and its disposition defers to a human target
  decision (`Retain`, `Preserve Dormant`, `Redesign`, `Retire`).
- **Impact:** changes the effective denominator and Exit A expectations for any target
  with commented-out wiring; no semantic redesign of tracks or coverage.

## PI-2: Bootstrap read set is unsatisfiable before `0A` exists

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (bootstrap ambiguity report)
- **Affected sections:** §8.4 mandatory read sets, §3.1 `0A-PREFLIGHT.md`
- **Observed text:** "All modes load protocol constants/statuses, `0A`, `0D`, `0E`,
  `0F`, `0G` or `0H`, relevant source inventory/frontier/claim records, and the active
  snapshot metadata." (§8.4)
- **Issue:** The first Preflight invocation creates `0A` and the other containers; the
  literal read set cannot be satisfied, and the protocol defines no bootstrap
  exemption.
- **Proposed direction:** define the Preflight bootstrap minimum read set (protocol
  constants/statuses plus operator inputs and the target snapshot); the full base set
  becomes mandatory from the next invocation onward.

## PI-3: Invocation ID grammar is undefined

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (convention adopted and documented by
  the harness)
- **Affected sections:** §8.1 resume identity header
- **Observed text:** "Invocation ID: [globally unique]" (§8.1)
- **Issue:** The protocol requires global uniqueness but defines no grammar, allocation
  scope, or non-reuse rule for invocation IDs, unlike typed record IDs (§4.1).
- **Proposed direction:** either define a grammar and allocation rule (for example
  iteration-persona-POV-sequence with ledger-checked non-reuse) or explicitly delegate
  the format to the conforming implementation while requiring uniqueness and
  non-reuse against recorded state.

## PI-4: Protocol version literal differs between the 0A template and §8.1

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (version contradiction report)
- **Affected sections:** §3.1 `0A-PREFLIGHT.md` template, §8.1 resume identity header
- **Observed text:** "`- **Protocol Version:** 4.1`" (§3.1) versus
  "`Protocol Version: v4.1.3`" (§8.1)
- **Issue:** Two literals for the same field invite disagreement about which value is
  authoritative and how it is formatted (`4.1` versus `v4.1.3`).
- **Proposed direction:** keep one authoritative version source (the edition header),
  generate both occurrences from it, and define the exact literal format to record in
  artifacts.

## PI-5: Live/sandbox access must stay optional; static-only disposal is implicit

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28, gate `GATE-A-001` (NO-SANDBOX reported
  as the dominant irreducible blocker; runtime/deployed units stayed `[C-PARTIAL]` /
  `Unknown` and Exit A was NOT MET)
- **Affected sections:** §7.1 trust boundary, §9.3-9.6 probes and no-mock integrity,
  §10 coverage, §15.1 Exit E conditions, §0.1 `Unknown` fail-inclusive
- **Observed text:** "When static analysis is insufficient and a sandbox is
  unavailable: ... use the Human Hatch path and preserve the coverage gap/partial
  state." (§9.5); "Every effective in-scope source unit is `[C-COVERED]`; all approved
  exclusions are exact, risk-assessed, disclosed, and removed from denominator
  arithmetic; no `[C-GAP]` or `[C-PARTIAL]` remains." (§15.1); "Live acquisition is an
  operator-side activity. ... Agents MUST NOT invoke acquisition using production
  credentials." (§7.1)
- **Issue:** Nothing in the protocol requires execution, sandbox, probe, live database,
  or deployed-state access - all are stated or implied as operator-side and optional.
  But the protocol never states that static-only completion is a first-class posture,
  and the only legitimate Exit A path for runtime-unverifiable units (exact
  human-approved exclusions) is implicit. Runs therefore treat NO-SANDBOX as an
  irreducible prerequisite rather than a disposal choice, and cannot tell whether
  static-only completion is compliant.
- **Proposed direction:**
  - state explicitly that execution, sandbox, probe, live database, and deployed-state
    access are optional throughout; no mode, gate, or conformance condition requires
    them;
  - define the static-only disposal: runtime-unverifiable units become exact
    `Approved-Excluded` units with unit IDs, rationale, residual risk, authority, and a
    reason such as `Runtime-Observation-Unavailable` (or the existing `NO-SANDBOX` /
    `EXTERNAL-FACT-UNAVAILABLE`), or remain `[C-GAP]`/`[C-PARTIAL]` and honestly block
    Exit A;
  - keep `Unknown` fail-inclusive and forbid fabrication or mocks as evidence; exclude
    exact units, never blanket categories;
  - state that a static-only run is eligible through **both** gates: no mode, Exit A
    condition, or Exit E condition may require execution, sandbox, probe, live database,
    or deployed-state access; runtime-only aspects are disposed by exact exclusions with
    disclosed residual risk, and the remaining conditions (confirmations, handbook,
    equivalence, packaging) are satisfied on the static scope;
  - require the run to declare its posture (probe env available versus static-only) so
    gate results are interpreted against it;
  - apply the identical rule to databases and external platforms: sanitized exports and
    probes are nice-to-have; their absence is a disposal decision, not a protocol
    defect.
- **Impact:** No evidence-rule weakening; changes the completion posture and its
  documentation only. Exit E follows the same disposal rules.

## PI-6: Operator-facing interaction surface is unspecified and tag-dense

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28/30 (operator feedback: tags, aliases,
  and status tokens are hard to decode; questions and reports arrive as protocol
  objects rather than human requests; the 85-decision Exit E approval digest rendered
  to 1275 lines, and the operator chose wholesale approval because per-decision review
  was impractical)
- **Affected sections:** §8.6 concurrent persona buffers, §9 ticket and Human Hatch
  interactions, §13.1 audience paths (existing rule for reconstruction output), §15
  report schemas
- **Issue:** The protocol specifies the handbook as progressive, human-first output,
  but says nothing about the run-time interaction surface. Agents therefore present
  raw record tags, IDs, and status tokens when asking humans for decisions or
  reporting progress, which is reviewable only by decoding the record schema.
- **Proposed direction:** define minimum presentation rules for operator-facing
  questions and reports: plain-language summary first (what is needed and why), options
  with consequences and authority, explicit "what happens next", with IDs and status
  tokens secondary and hyperlinked. Machine records and schemas are unchanged; the rule
  governs projections only. Record schemas, not prose style, remain normative.
- **Impact:** presentation-only; no new statuses, fields, or authority. Pairs with the
  harness-side deterministic store in roadmap M3.5. Confirmed by the Exit E approval
  pass: an approval surface that a maintainer cannot practically read weakens the human
  authority it is supposed to exercise.

## PI-7: Cold-read scope is undefined between `0G` and `0H`

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (ambiguity report; the run loaded the
  full append-only `0G` history because "0G or 0H" does not say which)
- **Affected sections:** §8.4 mandatory read sets, §8.5 cold resume, §8.7 stale
  checkpoint guard, §8.8 `0G` invocation log
- **Observed text:** "All modes load protocol constants/statuses, `0A`, `0D`, `0E`,
  `0F`, `0G` or `0H`, relevant source inventory/frontier/claim records, and the active
  snapshot metadata." (§8.4); "`0H` is derived and never replaces `0G`." (§1.1)
- **Issue:** "`0G` or `0H`" does not define the cold-read scope. Reading all of `0G`
  grows without bound across a long run; reading only `0H` risks missing state when the
  summary is stale. The stale-checkpoint guard does not say how much history must be
  reloaded when the summary is invalid.
- **Proposed direction:** state the rule explicitly: an invocation loads a current,
  fingerprinted `0H` plus the `0G` range it cites; the full `0G` is required only when
  `0H` is absent, stale, or fails its fingerprint/stale check, and then the tail from
  the last validated checkpoint suffices. This is a read-scope clarification; `0G`
  stays append-only and `0H` stays derived and non-authoritative.

## PI-8: Harness tooling location and package membership are unspecified

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (ambiguity report; tooling written to
  `.tmp/` and described as non-canonical, with no rule to cite)
- **Affected sections:** §0.1 authority classes, §1.1 forensic-plane layout, §1.5
  non-authoritative sidecars, §15.3 package contents
- **Observed text:** "No reference implementation, product, service, CLI, UI, or agent
  harness is normative or required." (§0.1); "JSON, SQLite, graph, search,
  code-generation, or validator sidecars MAY be generated for automation. They MUST
  carry source fingerprints, schema version, generation time, and a prominent
  `NON-AUTHORITATIVE-DERIVATIVE` marker." (§1.5)
- **Issue:** The protocol permits sidecars and disclaims reference implementations, but
  never states where an executor's implementation tooling lives relative to
  `.extracted/`, whether it may be listed as a package member, or that the prompt path
  must remain executable without it. Runs therefore document the question instead of
  following a rule.
- **Proposed direction:** state that implementation tooling (runtimes, validators,
  helper scripts) is non-authoritative, lives outside `.extracted/` in operator-local
  scratch, is never a package member or certification input, and is never required by
  the protocol text: `.extracted/` contains protocol records only, and the
  Markdown-record path must remain executable with ordinary local tools.

## PI-9: Fail-inclusive stops at absent authority, not at recoverable artifact defects

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (sandbox probe pass; `grafana` schema
  recreated and a foreign table materialized from delivered columns instead of being
  parked as Unknown; the earlier classification had over-applied fail-inclusive)
- **Affected sections:** §0.1 class-3 predicates and fail-inclusive `Unknown`, §9.5
  no-mock integrity, §4.4 frontier terminality, §10 sweep validity
- **Observed text:** "When static analysis is insufficient and a sandbox is
  unavailable: ... use the Human Hatch path and preserve the coverage gap/partial
  state." (§9.5); "`Unknown` applicability is fail-inclusive: the dependency remains in
  the required read/validation closure and blocks mutation if it cannot be loaded and
  validated." (§0.1); "A test double used to exercise already-proven local logic MAY be
  a testing mechanism but MUST NOT be treated as observation of an opaque dependency."
  (§9.5)
- **Issue:** the line between "reconstruct what a delivered artifact already defines"
  and "fabricate missing authority" is implicit. Runs can therefore over-apply
  fail-inclusive and leave cells unresolved when the actual defect is recoverable (a
  dump missing its `CREATE SCHEMA`/role, a foreign table whose columns are delivered, a
  loadable artifact), while the opposite error - authoring a schema that was never
  delivered - is equally possible.
- **Proposed direction:** state the rule explicitly: fail-inclusive applies at absent
  authority (undelivered DDL, remote or production data, missing credentials,
  unreachable external systems), not at recoverable artifact defects. Reconstruct
  objects only from delivered DDL/columns; label synthetic data as such; probe the
  app-side logic those objects enable; keep the external/remote aspect as an explicit
  residual or exclusion. Authoring schema, columns, or values that were never delivered
  remains forbidden.
- **Impact:** no evidence weakening; reduces false `Unknown`s and makes the §9.5 line
  testable.

## PI-10: Exit E scope and test doubles for external dependencies

- **Status:** proposed
- **Source:** standalone run follow-up 2026-09-30 (Exit E was modelled as
  sandbox-blocked; operator directed local mock SAML/IdP and SMTP to continue)
- **Affected sections:** §9.5 no-mock integrity, §14.5 equivalence suite, §15.1 Exit E
  conditions, §15.4 certification sequence
- **Observed text:** "A test double used to exercise already-proven local logic MAY be
  a testing mechanism but MUST NOT be treated as observation of an opaque dependency."
  (§9.5); "Every effective in-scope source unit is `[C-COVERED]`; all approved
  exclusions are exact, risk-assessed, disclosed, and removed from denominator
  arithmetic." (§15.1); "Generated tests are reviewable derivatives." (§14.5)
- **Issue:** Exit E can be read as requiring runtime observation of external
  dependencies (IdP, SMTP relay), while §15.1 condition 2 already disposes such units
  through exact exclusions and §9.5 permits test doubles for proven local logic. The
  interaction is implicit, so runs may stall Exit E on sandbox availability
  unnecessarily, or over-claim results from mocks.
- **Proposed direction:** state that Exit E inherits the Exit A disclosed exclusion
  scope and does not re-require runtime-only units; the absence of a sandbox is never an
  Exit E blocker, and a static-only run is eligible for the full pipeline. Local mocks
  may scaffold equivalence tests and handbook walkthroughs only as labeled test doubles,
  never as external observations or coverage evidence; confirmations and approvals
  remain human acts. Packaging and signatures still follow §15.4 unchanged.
- **Impact:** no evidence weakening; makes Exit E reachable on a static-only scope
  without misrepresenting mocks.

## PI-11: Handbook quality gate admits metadata shells

- **Status:** proposed
- **Source:** standalone protocol run, Exit E verification 2026-10-05 (the run passed the
  handbook checks with a 17-file handbook whose chapters are roughly 29 lines each,
  mostly HBK envelope metadata plus one summary bullet)
- **Affected sections:** §13.2 chapter behavior, §13.4 handbook quality gate, §15.1
  condition 4, §5.7 CNF/bindings
- **Observed text:** "all required chapters exist and each contains at least one
  claim-backed semantic handbook assertion or an evidence-backed zero-domain statement;
  headings, navigation links, placeholders, or grouping prose alone are vacuous"
  (§13.4)
- **Issue:** the non-vacuity rule is satisfied by a single sentence per chapter, so a
  handbook can certify while delivering no usable reconstruction prose: no audience
  paths (§13.1), no observed-versus-target separation, no examples or diagrams, no
  navigation, and no assertion density tied to the underlying synthesis records. The
  run's `02-ARCHITECTURE.md`, for example, is a 29-line envelope with one bullet for 12
  modules.
- **Proposed direction:** define testable handbook substance: per-chapter required
  sections (audience purpose, observed legacy, target decisions, navigation), an
  assertion-density floor tied to claims and synthesis records, representative
  examples/diagrams where the chapter class requires them, and an explicit N/A with
  reason where a chapter is legitimately empty. The deterministic gate fails chapters
  that are shells; a human usefulness test remains separate from the boolean
  readability field.
- **Impact:** raises handoff work and makes Exit E's handbook claim meaningful; no new
  statuses or authority changes.

## PI-12: Human readability confirmation can be relay-only

- **Status:** proposed
- **Source:** standalone protocol run 2026-10-05 (90 HBK CNFs record readability "Pass"
  while stating "no independent line-by-line human read attested"; envelope audit
  fields are `None`)
- **Affected sections:** §5.7 CNF schema, §13.4 post-candidacy review, §15.1 conditions
  13 and 19
- **Observed text:** "After candidacy, authorized human readability review is recorded
  through one or more §5.7 `CNF` records ... and whose `Handbook Readability Review` is
  `Pass`" (§13.4); "`[B-CONFIRMED]` is never machine-assigned" (§15.1 condition 19)
- **Issue:** readability CNFs were relayed under operator standing authorization with
  empty envelope version/audit fields, and the declared `Handbook Readability Review`
  field was not used; the run explicitly disclosed that no line-by-line human read
  occurred, yet the gate passed. Relay may be acceptable for some approvals, but a
  handbook readability claim without a human read undermines the confirmation gate.
- **Proposed direction:** require readability CNFs to use the declared schema
  (`Handbook Readability Review` in Pass/Fail/Not-Applicable), record non-empty review
  scope, limitations, and timestamp, and bind the exact reviewed HBK payload; reject
  relay-only attestation for handbook readability, or require the operator to declare
  the handbook unread and block content readiness until a real review is recorded.
- **Impact:** prevents certifying unread handbooks; other confirmation types keep their
  existing relay policy.

## PI-13: Human Hatch conflates human authority with executor uncertainty

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28/10-05 and session export review
  2026-10-05 (operator was repeatedly asked to resolve protocol mechanics and executor
  uncertainty; recommended options were accepted as human choices, including one bulk
  sweep disposition later reverted as protocol-invalid)
- **Affected sections:** §8 invocation/resume semantics, §9 ticket and Human Hatch
  behavior, §10 sweep execution, §15 approvals and confirmations
- **Issue:** The protocol defines how unresolved work can reach a Human Hatch, but does
  not sharply distinguish a fact/choice that belongs to human authority from an
  executor that is uncertain about protocol interpretation, record mechanics, or its
  own reasoning. In the run, questions about sweep legality, gate semantics, record
  encoding, and certification mechanics were surfaced to the operator. Selecting the
  executor's recommended option then made model uncertainty appear as human authority.
  One such recommended bulk `[R-SWEPT]` operation was later discovered to violate the
  protocol and had to be reverted.
- **Proposed direction:**
  - require every escalation to classify its blocker before Human Hatch presentation:
    human authority/domain fact, unavailable external fact, unresolved source
    semantics, protocol mechanics, implementation/harness defect, or executor
    uncertainty;
  - permit Human Hatch only for genuine human-authority/domain decisions and external
    facts that the protocol explicitly assigns to a human; unresolved source semantics
    MAY reach a human only when the protocol requires domain interpretation rather than
    more executor/reviewer work;
  - protocol mechanics are resolved from the normative specification/runtime;
    implementation defects block the implementation; executor uncertainty is retried
    or independently reviewed, never converted into a human decision merely because a
    model asks;
  - a recommended option is presentation only and carries no additional authority;
    human approval cannot legalize a transition or artifact that violates protocol
    mechanics;
  - keep the resulting operator question linked to the semantic blocker and authority
    basis so later audits can distinguish decision, delegation, and attestation.
- **Impact:** strengthens rather than weakens human authority. It prevents the operator
  from becoming a fallback protocol interpreter and prevents executor mistakes from
  being laundered through human approval.

## PI-14: Mandatory read closure is conflated with model-context materialization

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28/10-05 and session export review
  2026-10-05 (1,026 assistant steps, 14 compactions, 136.4M reported input tokens;
  hundreds of invocations carried 100k+ token contexts while the run repeatedly loaded
  growing protocol/workspace state)
- **Affected sections:** §8.4 mandatory read sets, §8.5 cold resume, §8.7 stale
  checkpoint guard, §8.8 invocation log, context/read terminology throughout §8
- **Observed text:** "All modes load protocol constants/statuses, `0A`, `0D`, `0E`,
  `0F`, `0G` or `0H`, relevant source inventory/frontier/claim records, and the active
  snapshot metadata." (§8.4)
- **Issue:** The protocol correctly requires a complete authoritative read/dependency
  closure, but "load" can be interpreted by prompt-based executors as "materialize all
  required bytes into one model context." That couples semantic work to accumulated run
  history, drives context growth and compaction, and makes an executor deep inside a
  long-running context less able to independently reassess assumptions. The required
  dependency closure and the model's working projection are separate concepts.
- **Proposed direction:**
  - define the mandatory read set as the authoritative dependency closure that a
    conforming executor/runtime MUST validate, fingerprint, and make addressable before
    mutation;
  - explicitly allow a conforming runtime or harness to materialize only a bounded,
    task-complete projection of that closure into a semantic executor context, provided
    omitted dependencies remain deterministically retrievable and their identities /
    fingerprints are bound to the invocation;
  - require semantic invocations to be cold-replayable from committed state and
    declared inputs; conversation history, prior model reasoning, and compaction
    summaries are never authoritative prerequisites;
  - require independently constructed contexts for independent review/adjudication;
    the same underlying model MAY perform worker and reviewer roles because independence
    is established by context and role separation, not model identity;
  - failure to fit a semantic work item into a bounded projection is an execution /
    decomposition problem, not permission to depend on an ever-growing chat transcript.
- **Impact:** no reduction in evidence or read-set obligations. It separates protocol
  dependency correctness from prompt rendering and makes native/runtime, skill, and
  prompt backends capable of equivalent semantics without equivalent context size.
