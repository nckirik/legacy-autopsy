# Coverage profile — protocol.md@4.1.2 in CDL vocabulary

> Non-authoritative analysis. `protocol.md` remains the sole normative authority; this profile classifies its sections to plan migration and proves nothing about conformance.

**Status:** S2 deliverable (implemented).
**Inputs:** [`protocol.md`](../protocol.md)@4.1.2, [`analysis/protocol-classification.json`](../analysis/protocol-classification.json), [`analysis/metrics.golden.json`](../analysis/metrics.golden.json).
**Companion:** [`capability-contracts.md`](capability-contracts.md).

## Method and judgement boundary

Each of the 148 headings (20 Part containers + 128 sections) is classified in the CDL
vocabulary with one driver and zero or more tags.

- **MACHINE** — predominantly deterministic rules, algorithms, registries, arithmetic,
  or state transitions.
- **AGENT** — predominantly semantic judgement or prose.
- **MIXED** — both are load-bearing.
- **CONTAINER** — a level-1 Part with subsections; excluded from section metrics.

Tags: `gate`, `transition`, `artifact`, `capability`, `reuse`, `contract`
(AGENT/MIXED outputs that are structurally checkable), `human`, `policy`.

Judgement lives in the authored classification ([`cdl.md`](cdl.md) rule
`SEMANTIC-IS-NORMATIVE`: meaning is authored, mechanics are computed). Deterministic
checks (`internal/analysis`) verify only structure: every heading covered exactly once
in order, closed vocabularies, source fingerprint match, and arithmetic over authored
data. Regenerate metrics with:

```sh
go test ./internal/analysis/ -update
go test ./internal/analysis/
```

## Headline metrics

From [`analysis/metrics.golden.json`](../analysis/metrics.golden.json):

| Driver                          | Sections | Section lines | Line share |
| ------------------------------- | -------: | ------------: | ---------: |
| MACHINE                         |       59 |         1,294 |      45.0% |
| MIXED                           |       64 |         1,352 |      47.0% |
| AGENT                           |        8 |           229 |       8.0% |
| CONTAINER (Parts, unclassified) |       17 |             — |          — |

Tag frequency across the 128 sections: `artifact` 98, `contract` 50, `transition` 47,
`capability` 36, `gate` 33, `reuse` 17, `human` 14, `policy` 5. Structured outputs:
42 of 72 AGENT/MIXED sections are contract-shaped; 28 remain free prose.

## Findings

1. **Mechanics dominate.** 92% of classified section lines are MACHINE or MIXED. The
   protocol is not an undifferentiated prose document; the CDL bet is supported by the
   text itself.
2. **Artifacts are the center of gravity.** 98 of 128 sections define artifact, record,
   or table schemas. Declared tables plus structural validation are the largest
   deterministic surface, and the pilot's `TABLE`/`FIELD`/`ENUM` forms are the right
   foundation.
3. **Gates and transitions are the skeleton.** 33 gate and 47 transition sections,
   concentrated in Parts 6 (promotion/depromotion), 8 (invocation and resume), 9
   (ticket FSM), 10 (exits), 15 (Exit E), and 20 (artifact write matrix). These carry
   most completion semantics and should migrate early in S3.
4. **Reuse is concentrated.** §4.1, §4.1.1, §4.1.2, §5.1, §8.3, §8.4, §12.1, and
   §15.1.1 define profiles consumed across the document. One correct declaration buys
   correctness in many sections; migrate reuse first.
5. **Capabilities gate large scope.** 36 sections depend on supplied deterministic
   capabilities. The five contracts in [`capability-contracts.md`](capability-contracts.md)
   are therefore prerequisites for most of S3, not an optimization.
6. **Most AGENT work is contract-shaped.** 42 of 72 AGENT/MIXED sections produce
   enumerated outputs against declared evidence; only 28 are free prose, mostly Parts
   0, 2, 13, 16, and 19 guidance and policy. Semantic sections are not the migration
   bottleneck.

## Context-assembly leakage audit (named category)

Question: can EIR `EVIDENCE`/`USES` determine the mandatory read/scope set, or does the
harness need procedural logic?

Current state is a leak, quantified in `internal/contextpacket` and `skill/modes.json`:

- `commonReadTargets()`: 9 hardcoded workspace files.
- `staticModeReadTargets()`: 50 target entries across 8 mode cases covering 11 of the
  13 invocation modes.
- `strictMode()`: 3 modes hardcoded; persona/POV/shared-file pairing and the
  persona-local buffer path hardcoded.
- `iterations`: the 26 §10.6 tokens hardcoded.
- `skill/modes.json`: 306 hand-maintained lines, 13 projections citing 221 protocol
  sections as routing metadata.

Verdict: §8.4 read sets are protocol content currently living in Go, so `EVIDENCE`/
`USES` alone cannot reproduce them. The leak is bounded and declarative — lists, flags,
and path templates, not algorithms — and roughly 60 lines retire when §8 migrates.

Required S3 shape (no surface added now):

1. `MODE` declarations in CDL (mode id, strictness, mandatory read set, write-target
   class) compiled into EIR `modes[]`, replacing `staticModeReadTargets`/`strictMode`.
2. The §10.6 iteration sequence as a declared enum, replacing the hardcoded set.
3. Persona/POV path rules as declared path templates plus the existing path
   capabilities, not ad-hoc joins in the packet builder.
4. `skill/modes.json` routing regenerated or validated against EIR reference views.

Acceptable harness policy that is not leakage: packet rendering layout, file-read
mechanics, transport, and scheduling.

## Recommended S3 order

1. **§4.1–4.1.2** — identity, payload envelopes, canonical hash profile; highest reuse.
2. **§7.5–7.7** — acquisition and export reconciliation; pilot already done.
3. **§8** — invocation identity, read sets, cold resume; retires the leakage above.
4. **§10** — sweep, coverage arithmetic, Exit A–D.
5. **§3 + §5** — registries, prefix groups, finite-enum validation, evidence/decisions.
6. **§15 + §12** — packaging, evidence bindings, validation summary.
7. Remaining MIXED schema sections, then AGENT-only guidance (Parts 0, 13).

## What this does not prove

The classification is authored judgement; the checks prove coverage, vocabulary, and
arithmetic only. This profile does not prove semantic fidelity, that MACHINE sections
are fully expressible in the frozen CDL surface, that AGENT contracts are sufficient,
or anything about conformance, Exit A/E, package acceptance, or autopsy success.
