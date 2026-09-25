# Repository migration and disposition audit

> Non-authoritative implementation planning. Possibly provisional and reviewable. [`protocol.md`](../protocol.md) remains the sole normative authority until a parity report is accepted; [`cdl.md`](cdl.md) is the frozen language contract.

**Status:** proposal for review before the §7.7 pilot.
**Scope:** every committed code, fixture, skill, doc, and tool artifact, plus the untracked session archive.
**Companion:** [`runtime-architecture.md`](runtime-architecture.md) §7 defines the target shape; this document decides what happens to what already exists.

## 1. Principles

1. **The oracle is sacred.** Anything used to verify the new system must not be rewritten
   by the new system. `internal/fixtures` + `fixtures/` stay frozen and independent
   through bootstrap (parity test B).
2. **Proven replacements, not beliefs.** Nothing is deleted because the new design
   looks cleaner. Deletion requires a replacement that passes parity plus a usage audit.
3. **No new protocol logic in Go.** Transitional packages receive bug fixes only; every
   new protocol rule enters through source language → EIR → runtime.
4. **Additive pilot.** The §7.7 pilot adds directories and tests; it does not restructure
   existing packages.
5. **Explicit authority.** Every surviving artifact carries a status: normative authority,
   frozen language contract, design contract, implementation doc, design history, or
   frozen oracle.

## 2. Code and fixture disposition

| Artifact                                  | Role today                                  | Disposition                                                                                                       | Gate to change                               |
| ----------------------------------------- | ------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| `cmd/legacy-autopsy`                      | thin CLI entry                              | Keep; add spec/pilot wiring later                                                                                 | pilot                                        |
| `internal/cli`                            | command wiring                              | Rework (additive): spec compile/run commands; no protocol logic                                                   | pilot                                        |
| `internal/protocol`                       | protocol.md version/mode/section reader     | Freeze; parity oracle and routing input                                                                           | replaced by EIR when routing is generated    |
| `internal/markdown`                       | structural Markdown model                   | Keep as oracle parser; `cdl` gets its own parser                                                                  | never shared implicitly; revisit after pilot |
| `internal/routing`                        | skill/projection drift checks               | Freeze; candidate retire when skill views are generated from EIR                                                  | expressibility check                         |
| `internal/identity`                       | foundational typed IDs                      | Transitional freeze; becomes a deterministic capability with a declared contract                                  | pilot capability work                        |
| `internal/canonical`                      | basic Markdown hash (explicitly not §4.1.2) | Transitional freeze; superseded by stdlib canonicalization capability                                             | never use for protocol-significant hashing   |
| `internal/workspace`                      | non-fabricating skeleton init/check         | Keep; rebase onto runtime effect services and declared artifact schemas later                                     | pilot review                                 |
| `internal/contextpacket`                  | bounded context from protocol sections      | Keep with caution; **named leakage risk** (runtime-architecture §5); rework to resolve `EVIDENCE`/`USES` from EIR | pilot review                                 |
| `internal/fixtures` + `fixtures/`         | 24 bootstrap cases                          | **Freeze as independent oracle**; only additive fixtures with unchanged existing cases                            | never rewritten by the pilot                 |
| `skill/` (SKILL.md, modes.json, 13 modes) | hand-maintained routing projections         | Freeze pending expressibility check; candidate for generation from EIR reference views                            | coverage profile + pilot                     |
| `schemas/README.md`                       | placeholder                                 | Keep; fold into docs rework                                                                                       | docs phase                                   |
| `tools/export-helper/`                    | operator acquisition tool                   | Keep; out of protocol core; review separately                                                                     | later milestone                              |
| `examples/minimal`                        | synthetic example                           | Keep; reuse as prompt-backend acceptance fixture                                                                  | pilot                                        |
| `.github/workflows/ci.yml`                | M0 checks                                   | Update additively when pilot tests are stable (spec compile, test A/B)                                            | pilot completion                             |
| `AGENTS.md`                               | contribution rules                          | Update after pilot review with `cdl`/runtime boundaries and oracle-freeze rules                                   | post-pilot                                   |

## 3. Documentation disposition

| Doc                                                                          | Status                                     | Disposition                                                                               |
| ---------------------------------------------------------------------------- | ------------------------------------------ | ----------------------------------------------------------------------------------------- |
| `protocol.md`                                                                | sole normative authority                   | Keep frozen except demonstrated defects; generated later from `protocol.cdl`              |
| `docs/cdl.md`                                                                | frozen source-language contract            | Canonical; change only via errata or a new version                                        |
| `docs/execution-semantics.md`, `docs/eir.md`, `docs/runtime-architecture.md` | design contracts (below language)          | Maintain through pilot findings                                                           |
| `docs/spec-language-substrate-v0/v1/v2.md`                                   | superseded drafts                          | Deleted; [`cdl.md`](cdl.md) is the frozen contract                                        |
| `docs/architecture.md`                                                       | stale M0 boundary description              | Rework to the layered architecture after pilot                                            |
| `docs/conformance.md`                                                        | fixture contract + standalone quality test | Rework additively: add spec/backend parity (tests A/B), keep the oracle description       |
| `docs/development.md`                                                        | build/check workflow                       | Update additively with `cdl`, runtime, pilot commands                                     |
| `docs/roadmap.md`                                                            | milestone plan predating the DSL           | Rework: add an explicit Spec track and re-sequence around pilot evidence                  |
| `docs/runtime.md`                                                            | target service/workbench design            | Align with runtime-architecture (MACHINE VM, AGENT scheduler); keep as product design     |
| `docs/skill-runtime.md`                                                      | skill-first stage                          | Review against the interpreter model; likely partially stale                              |
| `docs/workspace.md`                                                          | workspace semantics                        | Keep; align with effect services and invalidation rules                                   |
| `README.md`                                                                  | visitor entry point                        | Rework after pilot to describe the spec track and current status honestly                 |
| `kiro-session-*.zip`                                                         | untracked private session archive          | Gitignored; never commit. File still sits in the tree until an archive location is chosen |

### Docs authority map (target)

```
normative           protocol.md (until crown) -> protocol.cdl (after)
frozen language     cdl.md
design contracts    eir.md, execution-semantics.md, runtime-architecture.md
implementation      conformance.md, development.md, workspace.md, runtime.md
planning            roadmap.md, migration-audit.md (this file)
```

No document may claim normative value in a domain owned by another row.

## 4. Transitional freeze policy

While the pilot runs, `internal/protocol`, `internal/markdown`, `internal/routing`,
`internal/identity`, `internal/canonical`, `internal/contextpacket`, `internal/fixtures`,
and `skill/` accept only:

- bug fixes that preserve all 24 fixture outcomes;
- additive changes explicitly required by the pilot;
- no new protocol behavior, no renames, no refactors.

`internal/workspace` and `internal/cli` may receive additive pilot wiring. Every change
to a frozen artifact must state which fixture outcomes were re-verified.

## 5. Migration phases and gates

```
Phase 0  audit accepted; session archive gitignored; freeze policy in effect   [done]
Phase 1  pilot (additive): cdl subset, runtime VM, reference VM,
         prompt backend, tests A/B for §7.7                               [done]
Phase 2  pilot review: leakage audit (context assembly), capability extraction
         decisions (identity/canonical), skill-generation decision, docs
         consolidation into history + authority map
         [mostly done: skill-generation decision still open]              [next]
Phase 3  migration: section-by-section protocol.cdl growth, dual-source drift
         checks, roadmap/README/CI rewrite
Phase 4  crown: parity report accepted by human; DSL becomes normative source;
         Markdown becomes generated
```

Gates:

- **G1** test A passes (native and reference VMs agree on trace hash for §7.7). [pass]
- **G2** test B passes (EIR behavior matches `protocol.md@4.1.2`; 24 fixtures unchanged). [pass]
- **G3** coverage profile produced with the context-assembly leakage category counted. [pass]
- **G4** pilot review accepted before any frozen artifact is reworked or retired. [S2]

## 6. Deletion policy

An artifact may be deleted only when all hold:

1. its replacement passes G1 and G2 for the affected scope;
2. a usage audit shows no remaining references in code, docs, or CI;
3. the deletion is recorded here with the replacing artifact and gate evidence.

Already-deleted, for the record: the rejected prompt-edition pipeline (maps, reviews,
ledger, generator) was untracked and removed before commit; it never entered history.
It is not a migration concern.

## 7. Open questions

1. Timing of the `cdl` extraction into its own module or repository.
2. Whether `internal/markdown` ever feeds `cdl` or remains oracle-only.
3. Whether the 13 skill projections are generated from EIR or kept hand-maintained
   with a routing validator.
4. Whether the Spec track is a new milestone series or a re-scoped M1.
5. When to regenerate `README.md`/`roadmap.md` — after G4, per Phase 2.
