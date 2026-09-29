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
  - require the run to declare its posture (probe env available versus static-only) so
    gate results are interpreted against it;
  - apply the identical rule to databases and external platforms: sanitized exports and
    probes are nice-to-have; their absence is a disposal decision, not a protocol
    defect.
- **Impact:** No evidence-rule weakening; changes the completion posture and its
  documentation only. Exit E follows the same disposal rules.

## PI-6: Operator-facing interaction surface is unspecified and tag-dense

- **Status:** proposed
- **Source:** standalone protocol run 2026-09-28 (operator feedback: tags, aliases, and
  status tokens are hard to decode; questions and reports arrive as protocol objects
  rather than human requests)
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
  harness-side deterministic store in roadmap M3.5.

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
