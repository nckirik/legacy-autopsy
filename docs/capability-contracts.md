# Deterministic capability contracts

> Non-authoritative design contract. Below the frozen CDL surface; adds no language surface. `protocol.md` remains the sole normative authority until a parity report is accepted.

**Status:** S2 deliverable (contracts defined; implementation only where noted).
**Related:** [`runtime-architecture.md`](runtime-architecture.md) §4 (taxonomy),
[`execution-semantics.md`](execution-semantics.md) §3–4 (availability and fail-closed),
[`coverage-profile.md`](coverage-profile.md) (36 capability-dependent sections).

## 1. Registration model

Every runtime binding declares:

```
id, version
kind: deterministic | proposing | effect
determinism: reproducible-exact | reproducible-class | heuristic
inputs, output schema
failure modes and diagnostics
profile/availability (OS, tool, runtime version) where applicable
```

Enforcement rule: **an authoritative MACHINE result requires deterministic provenance.**
A `proposing` provider used in a MACHINE step is a registration/compile error, never a
runtime downgrade. Effect services are applied at commit and never appear as expression
values.

Availability: `REQUIRES CAPABILITY` plus `NO-SUBSTITUTE-MACHINERY` — an unavailable
capability marks dependent results `unsupported`, blocks dependent effects, leaves
independent work untouched, and never licenses an executor-authored substitute.

## 2. Deterministic capability contracts

### 2.1 `hash.sha256` v1 — implemented (pilot)

- **Kind/determinism:** deterministic, reproducible-exact.
- **Inputs:** exact byte sequence, or a declared table for table hashing.
- **Output:** `sha256:<lowercase-hex>`.
- **Rules:** no normalization, no exclusion logic, no host-dependent input. The pilot
  hashes a declared table by canonical serialization (see `table.serialize`).
- **Failures:** unavailable → `unsupported`; empty input is valid (hash of empty bytes);
  any I/O failure is fail-closed.
- **Current:** `internal/capabilities` (pilot), used by `hash` EIR expression.

### 2.2 `table.serialize` v1 — implemented (pilot)

- **Kind/determinism:** deterministic, reproducible-exact.
- **Inputs:** a declared table with `ROWS`, `ROW-TYPE`, `KEY`.
- **Output:** canonical bytes usable by `hash.sha256`.
- **Rules:** rows sorted by key (stable), object keys sorted by the EIR canonical
  serialization, no timestamps or host data, only declared fields.
- **Failures:** unknown table → evaluation error (fail-closed); missing key field →
  evaluation error.
- **Current:** `internal/runtime.tableBytes`; canonical form via `cdl.CanonicalBytes`.

### 2.3 `path.normalize` v1 — transitional implementation

- **Kind/determinism:** deterministic, reproducible-exact.
- **Inputs:** a declared root and a candidate relative path.
- **Output:** normalized slash-separated path; read resolution confined to the root
  with traversal and symlink-escape rejection.
- **Rules:** reject absolute paths, `..` escapes, and symlink targets outside the root;
  preserve exact Unicode and case; no NFC/NFD normalization in v1.
- **Failures:** escape/invalid path → typed rejection (`PATH_INVALID` class); never a
  best-effort path.
- **Current:** `internal/identity.NormalizeRelativePath` and
  `contextpacket.confinedRead` (transitional). To re-express as a declared capability
  in S3; `cdl` must not depend on it.

### 2.4 `identity.typed-id` v1 — base generation implemented (S3 §4.1)

- **Kind/determinism:** deterministic, reproducible-exact.
- **Inputs:** type prefix, namespace, normalized owner coordinate, semantic
  discriminator, snapshot-independent fields per §4.1.
- **Output:** `TYPE-` plus the uppercase first 12 hex characters of
  `SHA-256(canonical-key UTF-8)`.
- **Rules:** §4.1 ID generation and §5.1 prefix groups; alias/tombstone handling per
  §4.6. The prefix and iteration registries are authored CDL declarations
  (`protocol/typed-id.cdl`), drift-checked against `protocol.md@4.1.2`; PRF/HBK
  fixed-discriminator rules are declared as normative rules.
- **Failures:** unknown prefix, malformed namespace/owner/discriminator → fail-closed
  diagnostics.
- **Current:** `internal/capabilities.TypedID` implements base generation and is tested
  against an independent recomputation of the §4.1 formula. Collision extension to
  16/20/…/64 characters and the `20-TRACEABILITY.md` registry checks are **not** yet
  implemented. §4.1.1/§4.1.2 remain unclaimed.

### 2.5 `canonical.markdown` v1 — bounded implementation (S3 §4.1.2)

- **Kind/determinism:** deterministic, reproducible-exact.
- **Inputs:** Markdown bytes; a named carrier field; envelope regions; artifact type and
  normalized path; evidence-set tuples.
- **Output:** canonical Markdown bytes plus the §4.1.2 fingerprints: semantic content,
  envelope, package member, file transport, and evidence set.
- **Rules and limits:** LF normalization, exactly one terminal LF, trailing-whitespace
  removal, blank-run collapse, and table pipe/separator normalization outside fenced
  code; fenced content is preserved verbatim. This bounded v1 does not implement a full
  Markdown AST boundary parser and does not normalize indented literal blocks.
- **Declared semantics:** §4.1.1 and §4.1.2 rules are declared in
  `protocol/semantic-payload-identity.cdl` and
  `protocol/canonical-hash-profile.cdl`, drift-checked against `protocol.md`
  (normative tokens, binding-kind table, envelope kinds, artifact types, row kinds,
  domain prefixes).
- **Failures:** unknown envelope kind, non-positive envelope version, unknown evidence
  binding kind, duplicate evidence tuples, missing anchor binding, and invalid row kind
  fail closed.
- **Current:** `internal/capabilities/canonical.go`. `internal/canonical` remains the
  deliberately limited bootstrap hash and must not be used for protocol-significant
  hashing. The same package implements the §8.5 `COLD-RESUME|` check fingerprint with
  carrier exclusion and fail-closed field validation.

## 3. Proposing providers

Source analyzers, AST explorers, heuristic matchers, search, and LLM-assisted
extraction are evidence producers, not capabilities in the MACHINE sense.

- Output carries provenance: provider id/version, inputs hash, method, reproducibility
  class.
- Usable only as AGENT `EVIDENCE` or registry input; never bound by a MACHINE step.
- No proposing provider is implemented in the pilot.

## 4. Effect services

Artifact write, state transaction, checkpoint persistence, append-only logs, and
filesystem mutation are runtime infrastructure:

- applied atomically at step commit (`execution-semantics.md` §3);
- never expression values and never named in `REQUIRES CAPABILITY`;
- invalidation is a status change, not deletion; attestations are append-only.

## 5. Contract status

| Contract             | Kind          | Status            | Blocks                  |
| -------------------- | ------------- | ----------------- | ----------------------- |
| `hash.sha256`        | deterministic | pilot implemented | §4.1.2, §12.5, §15.x    |
| `table.serialize`    | deterministic | pilot implemented | schema sections (98)    |
| `path.normalize`     | deterministic | transitional      | §8 roots, persona paths |
| `identity.typed-id`  | deterministic | base implemented  | §4.1, §5.1, §6–§12, §15 |
| `canonical.markdown` | deterministic | bounded v1        | §4.1.2, §14.5, §15.x    |

## 6. What this does not prove

These are design contracts for reference implementations. They prove nothing about
`protocol.md` conformance, Exit A/E, or autopsy success, and no unreferenced capability
may be assumed available.
