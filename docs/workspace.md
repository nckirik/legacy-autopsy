# Workspace guide

> Non-authoritative orientation only. `protocol.md` defines the actual files, schemas, ownership, write rights, and gate conditions.

`.extracted/` is the persistent execution state for one pinned system namespace and its declared environments/snapshots. It is designed to survive stateless invocations and must be validated from disk before mutation.

## Four planes

1. **Forensic** — records what was found, where, how it is invoked, its evidence, and uncertainty. Original snapshots and approved sanitized projections are evidence sources; structured forensic records are the canonical extracted representation.
2. **Assurance** — records denominator membership, source/frontier closure, claim evidence, export reconciliation, contradictions, decisions, coverage, traceability, confirmations, and gate reports. This plane is authoritative for those assurance concerns.
3. **Synthesis** — derives reconstruction semantics from closed forensic and assurance inputs. It may not invent a legacy fact. Semantic records carry derivation bindings and become stale when dependencies change.
4. **Human Reconstruction Handbook** — projects populated or confirmed synthesis and assurance information into navigable human guidance. It is synchronized output, not an independent semantic authority.

The exact tree belongs to `protocol.md`. Implementations should discover and validate it structurally rather than treating this guide as a schema.

## History and checkpoints

`0G-DECONSTRUCTION-STATE.md` is append-only invocation history. Every invocation appends its identity, loaded fingerprints, effects, blockers, and next required loads.

`0H-CHECKPOINT-SUMMARY.md` is a derived acceleration structure. It must identify the source `0G` range and fingerprint and never replace `0G`. A stale, absent, or insufficient `0H` requires deterministic reads of the authoritative records specified by the protocol.

## Initialization

Workspace initialization means **ready for deconstruction**, not **partially deconstructed**. An initializer may create required directories, empty registries, headings, and schema scaffolds, but it must not fabricate source facts or semantic progress.

In particular, templates must never invent:

- `[C-COVERED]` status;
- `[B-CONFIRMED]` status;
- passed or completed gates;
- evidence claims or source coordinates;
- personas, ownership, exclusions, decisions, confirmations, or snapshots not supplied and validated by the operator.

Missing information should remain explicitly absent, pending, unknown, or unsupported only where the normative schema permits that representation.

## Authority and sidecars

Plane authority is scoped, not global: assurance does not replace forensic evidence, synthesis governs derived semantics only, and the handbook remains a projection. Machine sidecars are non-authoritative derivatives governed by Protocol §1.5; future checked-in generated derivatives must identify their generator, source fingerprints, schema version, and generation time and must be regenerated rather than hand-edited.

Never commit a real `.extracted/` workspace, probe output, raw acquisition material, context packet, or operator-private state. Synthetic fixtures/examples belong only in their designated directories and carry no evidence authority.

`legacy-autopsy init` creates the four-plane Markdown skeleton atomically, including canonical shared POV/question scaffolds and handbook chapters. `workspace check` verifies required M0 paths and parses every required Markdown file structurally. It does not validate persona ownership, semantic schemas, record content, or protocol gates; those remain roadmap work.
