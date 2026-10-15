# Relational mapping semantics for the canonical store

> Implemented as the CDL `cdl/0.4` language revision with EIR format 4: every `FIELD`
> declares a storage class, and `internal/store.Schema` (schema version 3) derives
> tables from the declared classes, projects `child` rows as parent-keyed tables, and
> emits declared `REF` foreign keys. Nothing here changes protocol semantics.

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

## The implemented CDL extension

Storage clauses live on `FIELD` only; no new declaration kind was introduced.

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

- `STORAGE record | child | singleton | registry | log | artifact | value`
  - `record`: identity-bearing row with `la_record_id`, version, semantic
    fingerprint, and canonical payload.
  - `child`: owned rows with exactly one parent; primary key is the parent foreign key
    plus the declared natural/ordinal key; a child has no independent identity.
  - `singleton`: at most one row per declared scope (for example per workspace).
  - `registry`: addressable natural-key rows that other records reference.
  - `log`: append-only; update/delete rejected by trigger as `0G` is today.
  - `artifact`: a validation/emit shape for one packaging artifact; never a store
    table.
  - `value`: an embedded fragment or row shape shared by records; never a store
    table.
- `PARENT <FIELD-ID>[.<column>]` on `child`/`singleton`: declares ownership and the
  foreign-key target. The parent foreign-key column defaults to `la_parent_id`; an
  explicit `<column>` in the clause overrides it.
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

## Implemented generator rules

- `record`, `registry`, `singleton`, `log`, and `child` produce tables; `value` and
  `artifact` do not (value shapes stay embedded in the owning canonical payload, and
  artifact schemas compile to validators/emitters).
- Metadata columns derive from `STORAGE`; no per-field guessing. An unknown or missing
  class fails closed.
- `CHECK`s come from enums, declared `KEY` columns become `UNIQUE` natural keys,
  append-only triggers come from `log`.
- Child rows project with a `la_parent_id` foreign key, a `UNIQUE (parent, key...)`
  natural key over the declared `KEY` columns, and `store.IngestChild`
  (`<parent>|<key>|...` record identity); `ReplaceRecord` preserves the parent link.
  A child without `PARENT` or `KEY` fails generation.
- Declared `REF` clauses emit foreign keys. The target FIELD must be
  `record`/`registry`/`singleton`, and the target column must exist and be the
  target's declared `KEY` (so it is unique); anything else fails generation. A
  registry that is addressed by a column must declare it as its `KEY` (see PI-27 for
  `ENTRY-CLUSTER`).
- Foreign keys are enforced on every write path (`PRAGMA foreign_keys=ON` per
  transaction).
- `set<string>` writes are canonicalized; `list<string>` preserves input order.
- The `23`-`27` artifact payload fields compile to validators/emitters and do not
  enter the store schema.

## Version and migration impact

- language `cdl/0.4`, EIR format 4 (`EIRField` gains storage, parent, key, and
  reference declarations);
- `internal/store` `SchemaVersion` 3 (2 added storage classes, 3 adds child tables
  and REF foreign keys); the canonical dump serialization is additive, and dumps
  declare their schema version so imports fail closed across versions;
- protocol semantics are unchanged: these are language-level declarations. Until
  every `FIELD` declares `STORAGE`, table generation may only be used as a provisional
  development artifact.

## Non-goals

- no ORM, migration framework, or store-side semantic inference;
- once the language version flips, a `FIELD` that would become a table without
  `STORAGE` is a compile error rather than a guess.
