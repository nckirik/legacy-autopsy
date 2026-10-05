# Relational mapping semantics for the canonical store

> Non-authoritative proposal. The current `internal/store.Schema` mapping is a
> **bootstrap heuristic** (one table per declared `FIELD` block) and is marked
> provisional in code. Language changes require the normal CDL revision process and a
> language version bump; nothing here changes protocol semantics.

## Why table-per-FIELD is not enough

`FIELD` blocks currently describe at least five different things:

- identity-bearing records (`SRC`, `CLM`, `DEC`, ticket and capability records);
- owned child rows and value shapes (`ENTITY-COLUMN`, `USE-CASE-STEP`,
  `STATE-TRANSITION`, `AUTHORIZATION-MATRIX-ROW`);
- singleton metadata blocks (`PREFLIGHT-HEADER`-style records);
- append-only logs (`INVOCATION-LOG`, `0G`);
- packaging artifact payload schemas (the `23`-`27` reports), which are validation and
  emit shapes, not store tables.

Generating one table per `FIELD` therefore produces tables for value shapes and
artifact schemas, gives every table an invented `la_record_id`, and cannot know foreign
keys, ownership, ordinals, or which parent column corresponds to an owned collection.
At the same time parents still carry those structures as opaque text (for example
`columns-fields -> string`), so the same structure is both serialized and exploded into
disconnected tables. Storage semantics belong in the language, not in generator
heuristics.

## Proposed smallest coherent CDL extension

Add storage clauses to `FIELD` only. No new declaration kind is introduced.

```cdl
FIELD ENTITY-COLUMN
  STORAGE child
  PARENT ENTITY-COLUMN-LIST
  KEY ordinal, column-name
  ...
END

FIELD INVOCATION-LOG
  STORAGE log
  KEY sequence
  ...
END

FIELD TICKET
  STORAGE record
  REF owning-decision -> DECISION.dec-id
  ...
END

FIELD EXIT-E-CANDIDATE-REPORT
  STORAGE artifact
  ...
END
```

### Clauses

- `STORAGE record | child | singleton | registry | log | artifact`
  - `record`: identity-bearing row with `la_record_id`, version, semantic
    fingerprint, and canonical payload.
  - `child`: owned rows with exactly one parent; primary key is the parent foreign key
    plus the declared natural/ordinal key; a child has no independent identity.
  - `singleton`: at most one row per declared scope (for example per workspace).
  - `registry`: addressable natural-key rows that other records reference.
  - `log`: append-only; update/delete rejected by trigger as `0G` is today.
  - `artifact`: a validation/emit shape for one packaging artifact; never a store
    table.
- `PARENT <FIELD-ID>[.<column>]` on `child`/`singleton`: declares ownership and the
  foreign-key target. The parent column name derives from the child field name unless
  an explicit `<column>` is given.
- `KEY <column>[, ...]`: primary or natural key. Defaults to `la_record_id` for
  `record`/`registry`; required for `child` and `log`.
- `REF <column> -> <FIELD-ID>.<column>`: foreign key. The target must be
  `record`, `registry`, or `singleton`; target existence and acyclicity are compile-time
  errors.
- Collection semantics stay in the existing type system: `list<string>` preserves order
  and duplicates; `set<string>` is stored in the one canonical form (ascending, unique;
  see `certify.CanonicalStringSet`). A collection of structured rows uses `STORAGE
  child`, never a string column.
- `VERSIONED` (`record`/`singleton` only, optional): keep row versions instead of
  overwriting.

### Parent text columns

The child's `PARENT` clause is authoritative for storage. A parent's declared text
field for the same collection (for example `columns-fields -> string`) remains part of
the canonical payload as its human/audit rendering and is not a queryable column. The
generator skips it when a matching owned child exists and keeps it as a text column
otherwise.

## Mechanical generator rules

- Only `record`, `registry`, `singleton`, `log`, and `child` produce tables; `artifact`
  produces a validation/emit schema only.
- Metadata columns derive from `STORAGE`; no per-field guessing.
- `CHECK`s come from enums, `UNIQUE`/primary keys from `KEY`, foreign keys from `REF`,
  append-only triggers from `log`.
- `set<string>` writes are canonicalized; `list<string>` preserves input order.
- The `23`-`27` artifact payload fields compile to validators/emitters and must not
  enter the store schema.

## Version and migration impact

- language `cdl/0.4`, EIR format 4 (`EIRField` gains storage, parent, key, and
  reference declarations);
- `internal/store` `SchemaVersion` 2; the canonical dump serialization is additive;
- protocol semantics are unchanged: these are language-level declarations. Until
  every `FIELD` declares `STORAGE`, table generation may only be used as a provisional
  development artifact.

## Non-goals

- no ORM, migration framework, or store-side semantic inference;
- once the language version flips, a `FIELD` that would become a table without
  `STORAGE` is a compile error rather than a guess.
