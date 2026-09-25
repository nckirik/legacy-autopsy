# Runtime architecture

> Non-authoritative architecture projection. Possibly provisional and reviewable; it defines component boundaries below the frozen source language and adds no language surface.

**Status:** draft for the §7.7 pilot.
**Depends on:** [`cdl.md`](cdl.md) (frozen language), [`eir.md`](eir.md) (interchange), [`execution-semantics.md`](execution-semantics.md) (execution model).
**Authority:** `protocol.md@4.1.2` remains sole normative authority until parity. Nothing here may contradict CDL or the execution semantics.

## 1. Layers

```
protocol source (protocol.cdl)
        │
        ▼
CDL language
  grammar + static semantics + stdlib + projection machinery
        │
        ▼
Execution IR (EIR)
  versioned canonical interchange contract
        │
        ├────────────────────────┐
        ▼                        ▼
Native runtime              Prompt backend
(Legacy Autopsy)            (full / light)
        │                        │
        ├─ MACHINE VM             └─ partial interpretation: MACHINE -> instructions,
        ├─ AGENT scheduler +         AGENT -> normative clauses, authority -> proposed/draft
        │  structural validator
        ├─ state/artifact runtime
        └─ effect services
```

Layering rule: each arrow is a contract. The source language does not know the
runtime. The runtime does not re-interpret source; it consumes EIR. Prompt backends
consume the same EIR and are alternate execution targets, not document generators that
happen to resemble the protocol.

## 2. Dependency direction

```
cdl/       must not import legacy-autopsy runtime, capabilities, storage, or UI
legacy-autopsy    depends on cdl
```

- `cdl` depends on the Go standard library and its own stdlib definitions only.
- The dependency is enforced mechanically by an import-boundary test, not convention.
- Same repository initially. A future extraction of `cdl` into its own module or
  repository must require no code change beyond module paths.
- `protocol.cdl` contains no host-language constructs; it binds only declared
  capabilities and effect services.

Independence test: a third party can implement the EIR instruction set and capability
contracts and obtain equivalent protocol semantics. The pilot proves the contract is
implementable by shipping a tiny second runtime (test A below).

## 3. Components

| Component                         | Owns                                                                          | Must not                                           |
| --------------------------------- | ----------------------------------------------------------------------------- | -------------------------------------------------- |
| `cdl` parser/resolver/typechecker | source structure, references, types, capability closure, IDs                  | know about workspaces, agents, filesystem          |
| `cdl` compiler                    | deterministic EIR + hashes                                                    | embed runtime policy                               |
| `cdl` projection machinery        | deterministic renderers, templates, omission rules                            | generate wording freely                            |
| MACHINE VM                        | step transactions, staged effects, guards, workflow instructions              | validate semantic truth                            |
| AGENT scheduler                   | context assembly, dispatch, structural contract validation, draft attestation | judge semantic correctness, mark confirmed         |
| state/artifact runtime            | tables, values, states, checkpoint persistence, invalidation                  | reinterpret protocol rules                         |
| effect services                   | atomic writes, transactions, append-only logs                                 | behave as expression values                        |
| capability registry               | deterministic capabilities and proposing providers, versioned                 | let proposing providers feed authoritative results |
| backends                          | native, prompt, reference                                                     | diverge from EIR semantics                         |

## 4. Capability taxonomy (hard boundary)

| Kind                     | Examples                                                                                                               | Eligibility                                                                       |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| deterministic capability | `hash.sha256`, canonical serialization, exact parser/query over declared structures, deterministic filesystem metadata | eligible for MACHINE-authoritative computation                                    |
| proposing provider       | source analyzer, AST explorer, heuristic matcher, search, LLM-assisted extraction                                      | evidence/proposals only; may feed AGENT steps, never authoritative MACHINE values |
| effect service           | artifact write, state transaction, checkpoint, append log, filesystem mutation                                         | runtime infrastructure applied at commit; never an expression value               |

Enforcement: **an authoritative MACHINE result requires deterministic provenance.** A
proposing provider used in a MACHINE step is a registration/compile error, not a
runtime downgrade. This keeps the CDL authority model true without special cases.

Registration contract per binding:

```
id, version
kind: deterministic | proposing | effect
determinism: reproducible-exact | reproducible-class | heuristic
inputs, output schema
failure modes and diagnostics
profile/availability (OS, tool, version) where applicable
```

Proposing outputs carry provenance: provider id/version, inputs hash, method,
reproducibility class. They are evidence, and evidence is falsifiable — never an
authoritative value.

## 5. AGENT scheduling and context assembly

The scheduler resolves declared `EVIDENCE`/`USES` into a bounded context, dispatches,
validates the structural contract (`PRODUCES`, `RESULT`, evidence presence, row
coverage), and records a draft. It never evaluates a `REQUIRE` clause.

**Named leakage audit category (from the architecture review):** can `EVIDENCE` +
`USES` actually determine the mandatory context/read set? §7.7 does not stress this;
cold resume, strict source boundaries, and scoped evidence selection will. If the
harness needs substantial procedural logic to decide _which subset_ of declared
evidence an agent may see, protocol semantics have leaked into the harness. The
coverage profile must count context-assembly logic separately.

## 6. Backends and the two parity tests

Backends consuming the same EIR:

- native — Legacy Autopsy VM;
- prompt — `full.prompt.md` / `light.prompt.md` rendering;
- reference — documentation projection;
- optional generated validators — optimization/export artifacts derived from EIR, never
  the fundamental native execution mechanism, and never the oracle for the interpreter.
  The "native validator projection" named in CDL §14 is therefore realized as EIR
  execution by the native VM, not as bespoke generated Go validator logic.

Two orthogonal tests, neither implying the other:

```
A. Backend conformance
   same EIR + same inputs
     ├─ reference VM
     └─ Legacy Autopsy VM
   -> identical canonical trace hash          (execution-semantics.md §13)

B. Protocol parity
   EIR behavior
        versus
   protocol.md@4.1.2 + the independent 24-fixture oracle
   -> migration correctness, no silent semantic diff
```

Test B must use the existing fixture runner as an external oracle. Do not validate a
generated validator with the compiler that generated it.

## 7. Repository shape (proposed)

```
cdl/                        language: lexer, parser, ast, resolver, types, eir,
                            render, templates; no runtime imports
examples/spec/              bootstrap .cdl sources (no authority implied)
protocol/                   canonical source after crown (S4):
                            protocol.cdl, sections/*.cdl
internal/runtime/           MACHINE VM, step executor, state store, checkpoint
internal/capabilities/      deterministic capabilities, proposing providers, effects
internal/agents/            agent hosts/adapters (M6)
internal/reference/         (or tests/reference-vm) tiny second runtime
tests/parity/               test A and test B harnesses
```

During bootstrap, `.cdl` sources live under `examples/spec/` so they imply no
authority; S4 promotes them to `protocol/`. Multi-file assembly is a compiler input
contract (ordered files contributing declarations/sections), not a language include
directive; no surface is added.

Migration mapping from today's tree:

| Today                                                          | Disposition                                                                            |
| -------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `internal/protocol`, `internal/markdown`, `internal/routing`   | keep as bootstrap/parity oracle; some outputs become generated by `cdl`                |
| `internal/canonical`                                           | reconcile with the deterministic `hash`/canonicalization capabilities when implemented |
| `internal/identity`                                            | rework under declared `TYPE`/identity ledger; candidate for generated validators       |
| `internal/workspace`, `internal/contextpacket`, `internal/cli` | survive as runtime/host surfaces after boundary review                                 |
| `internal/fixtures`                                            | frozen as the independent protocol-parity oracle during bootstrap                      |
| `skill/` + 13 projections                                      | candidate for derivation from EIR reference/routing views                              |
| `fixtures/` 24 cases                                           | frozen parity baseline                                                                 |
| `docs/`                                                        | audit; substrate collapses to one canonical doc + design history                       |

## 8. Non-goals for the §7.7 pilot

No capability host, no agent orchestration, no service, no UI, no multi-autopsy
registry, no leases, no confirmation flow, no full compiler. The pilot implements the
instruction subset in `execution-semantics.md` §14, one deterministic capability, an
in-memory state store, one prompt projection, a stubbed AGENT dispatch, and the two
parity tests for the §7.7 slice.

Growing the frozen language to escape a pilot failure is a stop condition (CDL §14).
