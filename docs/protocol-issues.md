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
