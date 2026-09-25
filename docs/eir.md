# Execution IR (EIR)

> Non-authoritative implementation contract. Possibly provisional and reviewable; it defines the interchange format below the frozen source language and adds no language surface.

**Status:** draft for the §7.7 pilot.
**Depends on:** [`cdl.md`](cdl.md) (frozen language), [`execution-semantics.md`](execution-semantics.md) (what execution means).
**Authority:** `protocol.md@4.1.2` remains sole normative authority until parity.

## 1. Purpose

EIR is the stable contract between the language and every backend:

```
source
  │
  ▼
parser AST              ─┐
  ▼                      │ compiler-private; free to evolve
typed/HIR/MIR           ─┘
  │
  ▼
Execution IR (EIR)      ← external contract: native VM, prompt backends,
                          reference VM, generated exports
```

The compiler's internal representations are implementation details. EIR is not: it is
versioned, canonically serialized, hashed, and consumed identically by all backends.
If native execution and prompt rendering consume the same EIR hash, backend
consistency becomes checkable (test A) instead of aspirational.

EIR carries semantics, never host behavior. It contains no scheduling, retry, storage,
or UI concerns.

## 2. Envelope

Every EIR document starts with an envelope:

```
eir-format:        1
language:          cdl/<version>
stdlib:            cdl-stdlib/<version>
protocol:          canonical-deconstruction/<version>
source-fingerprint: sha256:<hash over ordered source bytes>
generator:         cdl/<version>
eir-sha256:        sha256:<hash over canonical EIR body>
```

- The source fingerprint binds the exact input bytes and ordered file list.
- The `eir-sha256` covers the canonical serialization of the body (§3).
- Provenance headers on generated artifacts and execution traces cite both fingerprints.
- `protocol` is carried now so parity reports and projections can pin it; it becomes
  self-referential once the protocol is compiled from source.

## 3. Canonical serialization and hashing

Canonical form (pilot: canonical JSON):

1. UTF-8, LF line endings, no BOM, no insignificant whitespace.
2. Object keys sorted by Unicode code point; arrays preserve semantic order.
3. Strings use only required escaping; no locale-dependent formatting.
4. Integers only for numbers; no floats, exponents, or negative zero.
5. No timestamps, absolute paths, hostnames, process IDs, or randomness.
6. Identifiers are the declared stable IDs, never source line numbers.
7. Hashes are lowercase hex of `sha256` over the canonical byte sequence.

Determinism rule: compiling identical source bytes with identical language, stdlib,
and generator versions MUST yield byte-identical EIR. A nondeterministic EIR is a
compiler defect, not a formatting difference.

Execution traces use the same canonicalization; `trace-hash` is defined in
`execution-semantics.md` §13.

## 4. Semantic content

EIR body sections:

| Section | Contains |
|---|---|
| `declarations` | types, values (with lifetime), states (`VALUES`/`INITIAL`/`ALLOWS`), fields, enums, artifacts, registries, tables, gates, workflow targets, capabilities (kind/version), rules and predicates |
| `sections` | section identity, number, title, goal, uses/requires edges, ordered steps |
| `steps` | id, owner, typed contract (evidence/produces/result/require for AGENT), ordered operations (MACHINE) |
| `projections` | logical views, channels, overrides, omissions, supplied-by |

Resolution requirements: every reference is resolved and bound in EIR; unresolved
references never reach EIR. Capability kinds are resolved so a MACHINE operation can
never bind a `proposing` provider (`runtime-architecture.md` §4).

## 5. Instruction set v1 (pilot subset)

Expressions:

```
literal        integer | string | boolean | enum-value | set-literal
path           resolved declaration reference (qualified across sections)
call           COUNT | HASH | NEXT | ALL | ANY | AVAILABLE | UNAVAILABLE
operators      == != < <= > >= + - AND OR NOT IN
```

Effects and control:

| Source (CDL) | EIR operation | Outcome impact |
|---|---|---|
| `SET x := expr` | `set` | staged effect; stutter-aware |
| `GUARD expr` + `OTHERWISE BLOCK g,...` | `guard` + `block` | `blocked` commits block set |
| `IF`/`ELSE` | `branch` | no step-level outcome |
| `FOR EACH x IN reg` | `foreach` | iterates canonical order |
| `APPEND row FOR x` | `append` | staged table effect |
| `REQUIRES CAPABILITY c` (step declaration) | `requires` field | unavailable ⇒ `unsupported` |
| `BLOCK g,...` | `block` | commits gate block |
| `REOPEN`, `RETURN`, `LOOP`, `REQUEST` | `reopen`, `return`, `loop`, `request` | defined; not executed in pilot |

AGENT steps serialize as their contract plus the verbatim `REQUIRE` text (which is
normative content, not a comment).

Diagnostics use stable codes; the pilot catalog is minimal
(`CAPABILITY_UNAVAILABLE`, `EIR_EVAL_ERROR`, `GUARD_BLOCKED`, `CONTRACT_INVALID`).

## 6. Projections in EIR

Projections are part of the contract because prompt backends consume them:

```
projection := {
  id, channels[], covers[],
  step-overrides: { step -> text },
  rule-overrides: { rule -> text },
  omissions: [ { step, channels[]?, reason, supplied-by } ]
}
```

The compiler resolves omission legality (CDL §7) before emitting EIR. Prompt renderers
perform no semantic filtering of their own; they render what EIR declares.

## 7. Versioning

| Change | Bump |
|---|---|
| new EIR field/instruction, backward compatible | `eir-format` |
| changed meaning of an existing instruction | `language` + `eir-format`, new fixture |
| stdlib operation semantics | `stdlib` |
| protocol behavior | `protocol` |
| renderer wording only, golden fixtures updated | generator version |

Compatibility: a backend declares the EIR formats it consumes. Unknown format,
language, or stdlib versions fail closed. EIR documents are reproducible and may be
cached, but caches are never authority; the source fingerprint is.

## 8. Compiler-private layers

AST, HIR, and MIR are explicitly outside the contract. Renaming, splitting, or
replacing them requires no fixture or backend change. If a compiler-internal change
alters emitted EIR, that is an EIR-visible change and requires canonical golden
fixtures.

## 9. Source assembly

The compiler accepts an ordered list of source files. Each contributes declarations
and sections; duplicate stable IDs are errors; declaration references resolve across
the union. This is an input contract, not source syntax: there is no `INCLUDE`
directive in the frozen language. The eventual single `protocol.cdl` is an assembly
of these files, not a required monolithic edit.

## 10. EIR example (abbreviated, §7.7)

```json
{
  "envelope": {
    "eir-format": 1,
    "language": "cdl/0.1",
    "stdlib": "cdl-stdlib/0.1",
    "protocol": "canonical-deconstruction/4.1.2",
    "source-fingerprint": "sha256:...",
    "generator": "cdl/0.1.0",
    "eir-sha256": "sha256:..."
  },
  "declarations": {
    "capabilities": [
      {"id": "hash.sha256", "kind": "deterministic", "version": 1}
    ],
    "states": [
      {"id": "reconciliation-complete", "section": "export-reconciliation",
       "values": ["unknown", "true", "false"], "initial": "unknown",
       "allows": [["unknown","true"],["unknown","false"],
                  ["true","false"],["false","true"]]},
      {"id": "fingerprint-status", "section": "export-reconciliation",
       "values": ["unknown","computed","unsupported"], "initial": "unknown",
       "allows": [["unknown","computed"],["unknown","unsupported"],
                  ["computed","unsupported"],["unsupported","computed"]]}
    ],
    "rules": [
      {"id": "reconciliation-completeness",
       "predicate": "complete",
       "expr": {"call": "eq",
                "args": [{"call": "add",
                          "args": [{"count": {"table": "13-EXPORT-RECONCILIATION",
                                              "where": ["reconciliation-state","==","Mapped-Atomic"]}},
                                   {"count": {"table": "13-EXPORT-RECONCILIATION",
                                              "where": ["reconciliation-state","==","Approved-Excluded"]}}]},
                         {"count": {"table": "13-EXPORT-RECONCILIATION"}}]}}
    ]
  },
  "sections": [
    {"id": "export-reconciliation", "number": "7.7",
     "steps": [
       {"id": "reconcile.populate", "owner": "MACHINE",
        "requires": ["table.serialize"],
        "ops": [{"foreach": {"var": "unit", "in": "normalized-units",
                 "do": [{"append": {"row-type": "RECONCILIATION-ROW", "for": "unit",
                                    "to": "13-EXPORT-RECONCILIATION"}}]}}]},
       {"id": "reconcile.classify", "owner": "AGENT",
        "evidence": ["approved-evidence"],
        "produces": ["13-EXPORT-RECONCILIATION.reconciliation-state",
                     "13-EXPORT-RECONCILIATION.clm"],
        "result": "RECONCILIATION-STATE",
        "require": "Container mappings do not reconcile their children, and entity correlations never satisfy reconciliation."},
       {"id": "reconcile.count", "owner": "MACHINE",
        "uses-rules": ["reconciliation-completeness"],
        "ops": [{"set": {"target": "reconciliation-complete",
                         "value": "reconciliation-completeness.complete"}}]},
       {"id": "reconcile.guard", "owner": "MACHINE",
        "ops": [{"guard": {"expr": ["reconciliation-complete", "==", true]},
                 "otherwise": [{"block": {"gates": ["A-STRUCTURALLY-COMPLETE",
                                                    "R-SWEPT", "Exit-A"]}}]}]},
       {"id": "reconcile.fingerprint", "owner": "MACHINE",
        "requires": ["hash.sha256"],
        "ops": [{"branch": {"if": {"call": "available", "args": ["hash.sha256"]},
                 "then": [{"set": {"target": "fingerprint",
                                   "value": {"call": "hash", "args": ["13-EXPORT-RECONCILIATION"]}}},
                          {"set": {"target": "fingerprint-status", "value": "computed"}}],
                 "else": [{"set": {"target": "fingerprint-status", "value": "unsupported"}},
                          {"block": {"gates": ["fingerprint", "artifact-completeness"]}}]}}]}
     ]}
  ],
  "projections": [
    {"id": "light", "channels": ["prompt", "native"],
     "step-overrides": {
       "reconcile.classify": "Classify every reconciliation row against approved evidence. Container mappings are non-closing until resolved, and entity correlations do not reconcile.",
       "reconcile.count": "Reconciliation is complete only when the mapped-atomic count plus the approved-excluded count equals the total number of normalized units.",
       "reconcile.guard": "Incomplete reconciliation blocks structural completion and Exit A."},
     "omissions": [
       {"step": "reconcile.fingerprint", "channels": ["native"],
        "reason": "The native harness recomputes the artifact fingerprint from canonical bytes.",
        "supplied-by": "native-harness"}]}
  ]
}
```

The excerpt is illustrative; canonical order, exact field names, and hash inputs are
fixed by the compiler's frozen `eir-format` fixtures, not by this example.
