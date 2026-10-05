# SQLite canonical record store and semantic capability surface

> Non-authoritative proposal for the next protocol version's storage binding. The
> protocol semantics and record schemas remain normative; the current 4.2 Markdown edition remains valid alongside the legacy
> 4.1.3 record.

## Decision summary

- **Canonical operational store:** `.extracted/legacy-autopsy.db` (SQLite), single
  writer, never hand-edited. Protocol records are persisted here; there is no parallel
  authoritative Markdown state.
- **Markdown is a projection/export:** human review, handbook material, audit bundles,
  and compatibility views may render Markdown on demand, but semantic executors do not
  maintain Markdown record files.
- **Relational semantics:** one table per record class with the existing declared
  fields; status/enum columns are `CHECK`-constrained from the §5.1 registries; typed
  IDs are unique per type; scope tuples (`COV`, cluster cells) are unique; references
  are keys; `0G` is append-only. Illegal tokens and dangling references fail at the
  storage/runtime boundary rather than becoming model cleanup work.
- **Canonical serialization:** a deterministic SQL dump (schema plus rows in declared
  order). Fingerprints and package bindings are computed over canonical bytes; the
  binary database file is rebuildable. A text dump keeps certified state auditable,
  diffable, and hashable without introducing a second semantic authority.
- **Fingerprints unchanged:** each record keeps the canonical payload bytes the
  protocol already hashes next to its queryable columns. Dump-level hashes replace
  per-file transport fingerprints for packaging/manifests only where the new binding
  defines that substitution.
- **Semantic executor surface:** models do not receive raw SQLite write access and do
  not author record envelopes, fingerprints, coverage arithmetic, gate state, or
  Markdown. They use protocol-aware read, semantic-write, and control capabilities.
  The runtime compiles those operations to relational reads/transactions and owns all
  deterministic cascading effects.
- **Automatic read-set provenance:** every model-facing read capability records the
  exact committed IDs/versions/fingerprints observed by the invocation. Submission is
  stale-checked against that accumulated read set before any semantic mutation commits.
- **Independent review:** worker and reviewer invocations use the same semantic
  capability surface but separately constructed contexts/read sets. They MAY use the
  same model; independence is established by context and role separation.
- **Raw SQL/table helpers are implementation/admin tools:** maintainers, migration code,
  validators, and tests may use lower-level SQLite operations. They are not the normal
  semantic-executor contract and must not become a way to bypass protocol transitions.
- **Versioning:** the binding becomes canonical at a minor protocol version (candidate
  v4.2 or v4.5), not a major bump on its own. The current edition is 4.2 with its
  frozen oracle; 4.1.3 remains as historical reference with the 24 fixtures; conversion must round-trip and the binding needs
  its own conformance corpus before any conformance claim.

## Semantic executor capability surface

The exact tool names and payload schemas should be generated from the protocol/EIR and
may evolve. The important boundary is semantic intent rather than physical storage.

Illustrative read capabilities:

```text
get_record(id)
search_records(text?, kinds?, scope?, status?, limit?)
get_related(id, relation?, direction?, kinds?)
get_evidence(id)
get_source_excerpt(source_id, coordinate)
get_open_work(scope?)
get_context(id, depth?, kinds?)
```

Illustrative semantic-write capabilities:

```text
propose_claim(...)
classify_behavior(...)
relate_records(...)
set_disposition(...)
attach_evidence(...)
resolve_frontier(...)
raise_question(...)
```

Illustrative control capabilities:

```text
submit_work(...)
request_review(...)
abandon_work(...)
```

These operations express what the executor learned or intends. They do **not** expose
mechanical setters such as:

```text
set_fingerprint(...)
set_coverage_percentage(...)
mark_gate_passed(...)
write_certificate(...)
```

Those are runtime/MACHINE effects.

The semantic API should also avoid binding model behavior to the physical schema.
`search_records` or `get_related` may compile to joins, FTS, recursive CTEs, or later
indexes without changing the executor contract. Raw `SELECT` is unnecessary for normal
semantic work and would leak storage layout, make schema migration model-facing, and
allow accidentally unbounded queries.

## Invocation read tracking and stale safety

A semantic invocation begins with a small runtime-rendered projection. It may expand
its view through read capabilities as reasoning requires. Every successful read extends
the invocation's authoritative dependency set automatically:

```text
worker reads:
  get_record(SRC-0187)
  get_related(SRC-0187)
  get_record(SRC-0041)
  get_evidence(CLM-0022)

runtime binds:
  SRC-0187@fingerprint
  REL-0321@fingerprint
  REL-0322@fingerprint
  SRC-0041@fingerprint
  CLM-0022@fingerprint
  EVD-* @fingerprint
```

The model never has to reproduce that ledger. At submission, the runtime rechecks the
bound rows. A changed dependency rejects the semantic mutation as stale and schedules a
fresh invocation.

This preserves strong protocol read-set/fingerprint semantics while allowing model
context to scale with the semantic problem rather than with total autopsy age.

## Transaction ownership

Semantic writes are proposals. The runtime owns the transaction that turns an accepted
proposal into canonical state. One semantic change may deterministically trigger
multiple storage effects, for example:

- validate IDs, enums, scope, ownership, and allowed transition;
- insert/update the semantic record rows;
- maintain relationship rows and reciprocal integrity;
- invalidate stale dependent confirmations/reviews;
- recompute affected coverage/frontier aggregates;
- recompute canonical payloads and fingerprints;
- update gate eligibility without letting the model set gate state;
- append the invocation/`0G` audit result;
- mark derived projections stale;
- commit atomically or not at all.

The executor should not need to know which of these mechanical consequences occurred in
order to make the semantic decision.

## Lower-level maintenance surface

A small `la` helper remains useful for operators, migrations, conformance tests, and
debugging. It is not the semantic model API.

```text
la init                     create/pin the schema and snapshot
la import <workspace>       deterministic 4.x Markdown -> SQLite conversion
la inspect ...              bounded administrative inspection
la validate                 constraints plus deterministic recomputation
la render <view>            Markdown/report projections for humans
la dump                     canonical SQL dump for hashing and packaging
```

Migration/test code may have additional table-level helpers behind internal interfaces,
but normal model execution goes through protocol-aware capabilities.

## Human and machine projections

The canonical SQLite store is not intended as a human interface. The runtime may derive
multiple views from the same committed state:

- bounded semantic context for a worker;
- independently constructed review context;
- plain-language operator questions and decision batches;
- handbook chapters and navigation;
- Atlas/search/index projections;
- forensic/audit Markdown or deterministic SQL dumps;
- final certification/package material.

A projection may be discarded and regenerated. It never becomes an independent source
of truth merely because a human or model consumed it.

## Owner sign-off needed

1. Minor-version number and timing for the canonical-storage flip.
2. Whether Markdown remains the official human-review export or another rendered view
   becomes primary. Either way, Markdown is no longer the working database.
