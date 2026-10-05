# Post-run hardening plan

> Non-authoritative planning document. The CDL sources under [`protocol/`](../protocol/)
> remain normative. This plan sequences the fixes derived from the first full standalone
> run (a small project, 2026-09-28/10-05) and its recorded issues PI-1..12.

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

## Principles (unchanged)

- `Unknown` stays fail-inclusive; nothing here weakens evidence or authority rules.
- Deterministic mechanics live in code; semantics stay with agents and humans.
- Runtime, probes, and live databases remain optional; static-only is a first-class
  completion path through both gates.
- Human confirmations, exclusions, and signatures remain human acts.

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

### W2 - Certification and verification engine

Make the run's ad-hoc machinery a maintained, reusable component (the light harness
core).

- Canonical fingerprint substrate for every record class (semantic content, DEC
  content, dependency-set, evidence-set, envelope, package-member, transport), computed
  once and reused by emitter and verifier.
- Registry-driven validity: enums/statuses/kinds from the §5.1 registry and `§4.1`
  type grammar, checked at write time and rechecked at verify time.
- Schema-conformant emitters for `23`/`24`/`25`/`26`/`27` and the §15.5 verifier;
  golden tests for emit-twice byte equality, tamper detection, and fail-closed results.
- Validator allow-lists (for example `.extracted/probes/[TICKET-*].log`) and explicit
  post-conditions for append-only logs.

### W3 - Records and storage binding

Reduce read/write friction with the SQLite binding from [sqlite-store.md](sqlite-store.md):
one table per record class, constraints from the registries, append-only `0G`,
canonical dump for hashing/packaging, and `la` helpers. Markdown remains an export for
review, generated, never hand-edited.

### W4 - Handbook and human surfaces

- A handbook renderer that produces narrative chapters from `90`-`96` synthesis:
  audience purpose, observed legacy versus target decisions, representative
  examples/diagrams, navigation, and volume tied to the underlying records.
- A deterministic substance lint per PI-11, plus a real human readability workflow with
  complete CNF fields per PI-12.
- An operator renderer per PI-6: plain-language questions with options and
  consequences; one-line-per-decision indexes; reports whose first screen answers what
  changed, what is blocked, and what is needed.

### W5 - Acceptance rerun

Run the same small project again with v4.2 and the light harness.

- Acceptance criteria: zero illegal tokens by construction; every fingerprint
  recomputed (no decorative values); handbook chapters meet the substance floor and a
  human actually reviews them; a decision batch is reviewable in the target time; a
  static-only run reaches both gates with disclosed residuals and no runtime
  requirement; the §15.5 verifier passes on the signed bundle and fails on tamper.

## Sequencing and gates

1. Owner decisions below; then W1 changes accepted PI-by-PI with fixtures.
2. W2 designs against W1 wording; can start in parallel, must land before W5.
3. W3 may follow W2; the W5 rerun may be Markdown-first to isolate protocol and
   handbook fixes before the storage binding is mandatory.
4. W4 before W5; W5 produces the improvement metrics.

## Decisions needed from the owner

1. Revision mechanics: bump to **v4.2**, and snapshot each accepted revision as a new
   frozen oracle (keeping previous oracles; drift tests target the current one).
2. W5 shape: Markdown-first rerun, or wait for the SQLite binding and rerun once on it.
3. Handbook substance floor baseline (for example: required chapter sections plus a
   minimum claim-backed assertion count per chapter, with reason-bearing N/A allowed).
4. Operator-surface minimums and a target review time for a 50-decision batch.

## Immediate next step

Accept or reject PI-1..12 one by one (updating their status), then start W1 with the
accepted set. The storage and handbook workstreams are already sketched in
[sqlite-store.md](sqlite-store.md) and roadmap M3.5.
