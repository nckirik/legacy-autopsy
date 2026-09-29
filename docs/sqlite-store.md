# SQLite record store

> Non-authoritative proposal for the next protocol version's storage binding. The
> protocol semantics and record schemas are unchanged; the current 4.1.3 Markdown
> workspace remains valid until the version flips.

## Decision summary

- **Canonical store:** `.extracted/legacy-autopsy.db` (SQLite), single writer, never
  hand-edited. This is where records are written; there is no parallel Markdown state.
- **Canonical serialization:** a deterministic SQL dump (schema plus rows in declared
  order). Fingerprints and package bindings are computed over dump bytes; the binary
  database file is derived and rebuildable. A text dump keeps the bundle auditable,
  diffable, and hashable; no second fingerprint domain is introduced.
- **Records:** one table per record class with the existing declared fields;
  status/enum columns are `CHECK`-constrained from the §5.1 registries; typed IDs are
  unique per type; scope tuples (`COV`, cluster cells) are unique; reciprocal
  references are keys; `0G` is append-only via triggers. Illegal tokens and dangling
  references fail at write time.
- **Fingerprints unchanged:** each record stores the canonical payload bytes the
  protocol already hashes, next to its queryable columns, so semantic fingerprints stay
  reproducible. Dump-level hashes replace per-file transport fingerprints for packaging
  and manifests; that binding is the only new definition.
- **Light harness:** `legacy-autopsy-light/`, an optional archive with schema DDL, the
  `la` helper CLI, a validator, a renderer, and the prompt. It is not a framework: the
  user's harness may replace it, and a prompt-only run creates the schema first with
  ordinary `sqlite3`.
- **Human output:** rendered views emit Markdown and reports on demand; operator
  questions and reports lead with plain language (what is needed, why, options and
  consequences, authority, what happens next), IDs secondary (PI-6).
- **Versioning:** the binding becomes canonical at a minor protocol version (candidate
  v4.2 or v4.5), not a major bump on its own. 4.1.3 remains the Markdown reference with
  its frozen oracle and 24 fixtures;
  conversion must round-trip and the binding needs its own conformance corpus before any
  conformance claim.

## Helper surface

```
la init                     create the schema; pin namespace and snapshot
la import <workspace>       convert an existing Markdown workspace (deterministic)
la write <table> --json     append a record (constraints enforced)
la read <table> [--where]   query records
la update <table> --json    update a mutable record (not append-only tables)
la log                      append a 0G invocation entry
la checkpoint               write the derived 0H summary
la validate                 constraints plus deterministic recomputation
la render <view>            Markdown/report projections for humans
la dump                     canonical SQL dump for hashing and packaging
```

## Owner sign-off needed

1. Minor-version number and timing for the flip (proposed: v4.2 or v4.5).
2. Whether a Markdown export stays the official human-review format or rendered views
   replace it.
