# W6 comparative metrics - veri-yonetimi rehearsal

Method and sources: the runtime SQLite store (invocation log `0G`, semantic
record counts), the emitted packaging chain and review records under `preview/`,
the per-chapter narration attempt logs, and the session history for human acts.
The first-run baseline is the exported 2026 first run recorded in
[hardening-plan.md](hardening-plan.md) ("Acceptance metrics"). Values that were
not persisted by the runtime are marked **not recorded** instead of estimated.

## Runtime execution (bounded semantic work)

| Metric                     | First-run baseline | This rehearsal                             |
| :------------------------- | :----------------- | :----------------------------------------- |
| Wall clock                 | 160.5 h            | hours, not fully logged (see gaps)         |
| Operator messages / nudges | 71                 | operator acted at designed gates only      |
| Context compactions        | 14                 | 0 by construction (fresh unit contexts)    |
| Tool calls                 | 982                | not recorded as a counter                  |
| Total model input tokens   | 136.4 M            | not recorded                               |
| Recorded store invocations | n/a                | 49                                         |
| Invocation retries         | n/a                | 44 zero, 4 one, 1 four                     |
| Zero-domain declarations   | n/a                | 3 (DR, export reconciliation, one catalog) |

Invocations by mode (from `0G`): Partial-Synthesis 30, Discovery 10, Profile
Synchronization 5, Preflight 3, Export Acquisition 1.

## Committed canonical state on the fresh workspace

| Artifact family                 | Count |
| :------------------------------ | ----: |
| Source inventory records        |    24 |
| Claims                          |    29 |
| Coverage rows                   |    24 |
| Traversal matrix rows           |   150 |
| Sweep records                   |   150 |
| Derived-from edges              |   247 |
| Lifecycle eligibility records   |    24 |
| Synthesis blocks (all catalogs) |    60 |
| Persona profiles (PRF)          |     5 |
| Privacy findings (PRV)          |     4 |
| Confirmations (CNF)             |    22 |

## Handbook and human gates

- 16 chapters narrated; one chapter (02) revised from recorded reviewer feedback
  ("the summary isn't about the architecture") and re-accepted; chapter 09 and 10
  needed retries during generation.
- Readability review: 16 verdicts, final result Pass, fingerprints bound.
- Acceptance suite: 23 mappings, all reviewed (Pass), coverage complete.
- Confirmations: 22 CNFs (1 synthesis batch, 5 PRF, 16 HBK) and 16 B-CONFIRMED
  envelopes; signed outer manifest and Pass final verification receipt.
- Operator messages in this rehearsal were the designed human acts plus a handful
  of UX corrections; no completion nudges were needed after each stage was built.

## Trust and reproducibility evidence

- Candidate report 14/14 Pass with recomputed evidence-set fingerprints; manifest
  81 semantic bindings; readiness Ready-For-Certificate; chain and 51-member
  bundle verified from bytes; receipt fingerprint
  `sha256:c0ae9739...`.

## Gaps and follow-up instrumentation

- Per-unit wall clock, token usage, retry causes, and tool-call counts are not
  persisted by the runtime; only retry counts and outcomes reach `0G`. Adding
  duration/usage telemetry to the invocation log (or a sidecar metrics table)
  would make the next comparison mechanical.
- Operator-message counting needs a session-level counter; the harness currently
  records only store-visible acts.
