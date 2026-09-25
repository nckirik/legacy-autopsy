# Execution semantics

> Non-authoritative implementation contract. Possibly provisional and reviewable; it defines execution **below** the frozen source language and adds no language surface.

**Status:** draft for the §7.7 pilot.
**Depends on:** [`cdl.md`](cdl.md) (frozen source-language contract), [`eir.md`](eir.md) (interchange), [`runtime-architecture.md`](runtime-architecture.md) (components).
**Authority:** `protocol.md@4.1.2` remains the sole normative authority until a parity report is accepted. This document defines what it means to execute the frozen language; it must not contradict CDL.

## 1. Two layers: VM semantics versus harness policy

The interpreter boundary only prevents the old coupling if the split is explicit.

| CDL VM semantics (conformance-relevant) | Legacy Autopsy policy (harness) |
|---|---|
| MACHINE step evaluation and staged effect commit | Which worker executes the next step |
| Guard outcome and declared blocking | How many AGENT jobs run concurrently |
| Capability unavailability and cascade | Which model/backend handles an AGENT task |
| State/value lifetime realization | Retry, backoff, and scheduling |
| `BLOCK`, `REOPEN`, `RETURN`, `LOOP`, `REQUEST` meaning | UI presentation of blocked work |
| Checkpoint contents and resume validation | Multi-autopsy registry and isolation |
| Artifact invalidation versus deletion | Process/service lifecycle |
| Canonical execution trace | Internal metrics and telemetry |

Policy may change behavior that is not protocol-observable. If a policy choice changes
a protocol-visible result, it is a semantics defect, not a policy setting.

## 2. Execution units

- The **program** is the EIR produced by CDL from the source.
- The **transaction unit** is one `STEP`. A step either commits protocol-visible effects
  according to these rules or it does not.
- An **invocation** is a harness grouping of one or more steps with identity, scope,
  and `0G` bookkeeping. The pilot executes steps in order in one process without
  invocations; invocation semantics are defined here only as far as they constrain the
  transaction model.
- AGENT steps are not transactions of the VM: they are dispatched, awaited, and
  recorded; their outputs are drafts.

## 3. MACHINE step transaction model

Each MACHINE step executes against a **staged view** of the runtime state. The staged
view is a copy-on-write overlay; reads observe the step's own staged writes.

```
evaluate capability requirements
  capability unavailable -> outcome unsupported, discard staging, apply cascade
evaluate step body against staged view
  evaluation error        -> outcome failed, discard staging, fail-closed record
  guard false             -> outcome blocked, discard staging,
                             commit only the declared OTHERWISE effect set
  otherwise               -> outcome committed, atomically commit staged effects
```

Rules:

1. Exactly one outcome per step from `committed | blocked | unsupported | failed`.
2. `unsupported` and `failed` never write state; they record a diagnostic, mark
   dependent results unsupported, and block dependent effects
   (`RUN NO-SUBSTITUTE-MACHINERY`).
3. A guard failure is not an error: the declared `OTHERWISE BLOCK ...` set is the
   committed outcome.
4. Commit is atomic. Readers never observe partially applied effects.
5. An identity assignment (`SET x := <current value>`) is a stutter: it commits no
   transition and satisfies `ALLOWS` without a self-transition entry (CDL Errata E6).
6. Evaluation errors include unknown call results, type violations that escaped static
   checks, and deterministic capability failures. They are never converted to default
   values.

## 4. MACHINE evaluation details

- **Guards** require a boolean operand; there is no truthiness (CDL E4). A guard on an
  `unknown` predicate value fails closed and yields `blocked`, not `committed`.
- **Capability calls** are resolved through the runtime capability registry by
  declared identifier and version contract. A MACHINE step may bind only
  `deterministic` capabilities (see `runtime-architecture.md` §4). `proposing`
  providers are invisible to the VM; `effect` services are applied by the runtime at
  commit, never as expression values.
- **`HASH` and other reads** evaluate against the staged view at the point of use, so
  a step's own committed staging is what it hashed.
- **`FOR EACH`** iterates a declared registry/table in its canonical order (CDL: stable
  ordering is part of the collection's declaration, not of iteration).
- **`REQUIRES CAPABILITY`** without an in-source `AVAILABLE(...)` branch yields
  `unsupported` when unavailable. The compiler warns on this shape; it is legal.

## 5. AGENT step execution

The runtime does not evaluate semantic clauses. For an AGENT step it:

1. resolves the declared `EVIDENCE` bindings into a bounded context (see leakage audit
   in `runtime-architecture.md` §5);
2. dispatches the semantic task to the configured agent backend;
3. validates the structural contract: produced fields exist, result values belong to
   the declared domain, required rows are covered, evidence bindings are recorded;
4. records the output as `draft` with provenance (backend, model, inputs hash,
   timestamp, EIR hash);
5. hands the draft to the confirmation flow.

The runtime never records an AGENT output as `confirmed`; only an append-only human
attestation bound to the immutable payload hash does that (CDL §1).

Structural validation failures are `failed` for the step and do not write drafts.

## 6. State, values, and lifetimes

- `LET` bindings are step-local and vanish at step end. A cross-step read is a compile
  error and never reaches the VM.
- `VALUE` and `STATE` lifetimes are realized as follows:

| LIFETIME | Realization | Survives process restart |
|---|---|---|
| `step` | reserved for future use; same as `LET` | no |
| `section` | in-memory store keyed by `(section, identifier)` | no; recomputed or unavailable |
| `run` | in-memory store for the run/invocation sequence | no; recomputed or unavailable |
| `checkpoint` | durable workspace record with producing step, inputs fingerprint, EIR hash | yes, after validation |

- Tables are canonical ordered collections; `APPEND`/`REMOVE` are staged like any other
  effect and commit atomically with the step.
- AGENT-authored table fields are drafts. A MACHINE guard over draft state computes an
  authoritative *computation* over unconfirmed semantic input; confirmation status is
  tracked separately and may block completion attestation. The pilot records
  confirmation status but does not implement the confirmation flow.
- State transitions derive from committed `SET` effects only; staged-then-discarded
  assignments never appear in the transition history.

## 7. Workflow semantics

Workflow operations are VM instructions, not harness methods:

- `BLOCK <gates>` marks the named gates blocked with the step's diagnostic and reason.
  Blocking commits even on a guard-failure transaction.
- `REOPEN <records>` invalidates dependent records (status change), never deletes data.
  The pilot defines the instruction and stubs the dependency graph.
- `RETURN <target>` ends the current section execution with the declared workflow
  target outcome.
- `LOOP <target>` transfers control to the declared workflow target; loop accounting
  and iteration tokens are protocol state, not VM counters.
- `REQUEST human-input` records a pending human dependency and stops the affected
  branch; independent work continues.

The pilot does not execute `RETURN`, `LOOP`, `REOPEN`, or `REQUEST`; their semantics
are fixed here so the pilot cannot accidentally invent them later.

## 8. Authority propagation

- MACHINE outcomes are `authoritative` in the native channel because all VM bindings
  are deterministic or effect services; there is no promotion path for `proposed`
  values (CDL §1).
- Staged effects that fail to commit leave no authority trace.
- AGENT outputs are `draft` until confirmed; `proposed` never appears in native VM
  outcomes — it belongs to the prompt-emulation channel outside the VM.
- The canonical execution trace records authority per committed effect.

## 9. Checkpoints and resume

A checkpoint records, at minimum:

- EIR hash, language/stdlib versions, protocol version;
- completed step ids and their outcomes;
- `checkpoint`-lifetime values with producing step and inputs fingerprint;
- gate block state and pending human dependencies;
- append-only attestation positions.

Resume validates the EIR hash and versions, validates/recomputes `checkpoint` values
whose inputs changed, and marks stale dependent records rather than reusing them. Any
mismatch is fail-closed. The pilot implements a serialized in-memory checkpoint for
the §7.7 slice only.

## 10. Artifacts, invalidation, attestations

- Artifact invalidation is a status/provenance change, not deletion.
- Deletion is limited to staging/temporary state that has never been committed.
- Attestations (including confirmations) are append-only records bound to immutable
  payload hashes; they are never rewritten.

## 11. Invocation boundaries

An invocation records identity, scope, and the step outcomes it committed. `0G`-style
append history, leases, and staleness checks are protocol-visible and therefore
CDL semantics; the pilot does not implement them and must not fake them.

## 12. Concurrency

Deferred. The pilot executes single-threaded. `state` and `run` stores may assume one
writer. `checkpoint` records are written by a single VM at a time; concurrent execution
requires lease semantics that are not yet defined.

## 13. Canonical execution trace

Backend conformance requires a comparable trace:

```
trace := ordered list of step records
step record := {
  step, owner, outcome,                 # committed | blocked | unsupported | failed
  authority,                            # authoritative | draft | none
  effects: [ {target, operation, value-or-ref} ... ],   # canonical order
  diagnostics: [ {code, subject} ... ],                 # canonical order
  eir-hash, channel
}
trace-hash := sha256(canonical serialization of trace)   # see eir.md §3
```

Two runtimes consuming the same EIR and the same inputs must produce the same
`trace-hash`. This is test A (backend conformance) in `runtime-architecture.md` §6. It
is deliberately narrower than correctness: traces record what the VM did, not whether
`protocol.md` semantics were correctly migrated.

## 14. Pilot subset (exactly what is implemented for §7.7)

| Instruction/feature | Pilot |
|---|---|
| `SET` to `STATE`/`VALUE`, stutter rule | yes |
| `GUARD` + `OTHERWISE BLOCK` | yes |
| `IF` / `ELSE` | yes |
| `FOR EACH` / `APPEND` on a table | yes |
| `REQUIRES`/`AVAILABLE` capability check | yes |
| `HASH` deterministic capability | yes |
| expression operators `== != < <= > >= + - AND OR NOT` | yes |
| `COUNT` | yes |
| AGENT step structural contract + draft record | stub dispatch, real validation |
| `REOPEN`, `RETURN`, `LOOP`, `REQUEST` | defined, not executed |
| invocations, leases, `0G`, confirmation flow | not implemented |
| checkpoint/resume | in-memory serialization only |
