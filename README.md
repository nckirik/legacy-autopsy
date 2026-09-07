# Legacy Autopsy

> **Project status — M0 foundation**
>
> The repository can validate protocol/skill routing, create and inspect an empty `.extracted/` skeleton, generate foundational IDs, assemble bounded context packets, and run bootstrap fixtures. It does **not** yet execute semantic deconstruction modes, prove Exit A or Exit E, or claim Protocol v4 conformance.

> **Normative authority**
>
> [`protocol.md`](protocol.md) is normative. All skills, schemas, tools, docs, examples, and validators are non-authoritative implementations or projections. If they conflict, `protocol.md` wins.

## What Legacy Autopsy is

Legacy Autopsy is an evidence-grounded agent/harness system for deconstructing a legacy system into a reconstruction-ready package. An LLM performs bounded semantic investigation; provider-neutral code owns deterministic mechanics and validation.

The complete protocol pipeline is:

```text
Legacy System
   ↓
Forensic Plane
   ↓
Assurance Plane
   ↓
Exit A
   ↓
Synthesis Plane
   ↓
Human Reconstruction Handbook
   ↓
Human Confirmation
   ↓
Exit E
   ↓
Certified Reconstruction Bundle
```

This diagram describes Protocol v4, not current M0 implementation coverage. Exit A and Exit E are protocol gates that only their required deterministic artifacts and validations can establish.

## Design principles

- **Normative law stays singular:** `protocol.md` defines behavior; secondary material routes to or implements it.
- **Reasoning and mechanics stay separate:** LLMs interpret evidence and semantics; code handles identity, parsing, scope, ordering, hashing, validation, and state transitions.
- **Filesystem state survives sessions:** `.extracted/` is persistent execution state; conversation memory is disposable.
- **Evidence is not inference:** model reasoning cannot promote itself into source evidence.
- **Projections are not authority:** generated sidecars and human-facing views remain traceable derivatives.
- **Trust boundaries fail closed:** raw acquisition material, credentials, probes, and operator-private state stay outside agent-visible repository content.

Legacy Autopsy is not a generic code summarizer, automatic rewrite tool, modernization-architecture generator, or system that treats LLM inference as evidence. It does not authorize production access or mutation.

## Getting started

Go 1.24 or newer is required. From the repository root:

```sh
go run ./cmd/legacy-autopsy doctor
go run ./cmd/legacy-autopsy protocol check
go run ./cmd/legacy-autopsy init --workspace /tmp/legacy-autopsy-example/.extracted
go run ./cmd/legacy-autopsy workspace check --workspace /tmp/legacy-autopsy-example/.extracted
```

This creates only an initialization scaffold; it invents no source facts, evidence, coverage, confirmation, or gate result. Continue with the [minimal example](examples/minimal/README.md), or see [development.md](docs/development.md) for IDs, context packets, validators, fixtures, and contributor checks.

## M0 implementation maturity

Implemented bootstrap capabilities:

- a thin skill/router with exactly 13 mode projections and machine-validated routing metadata;
- a provider-neutral Go CLI: `doctor`, `init`, `protocol check`, `workspace check`, `id`, `context`, and implemented `validate` targets;
- a small structural Markdown model for headings, fields, code blocks, comments, tables, and bounded sections;
- foundational typed IDs, normalized SRC-FILE paths, relative-path validation, and an explicitly limited basic Markdown hash primitive;
- atomic creation and structural checking of a non-fabricated four-plane workspace skeleton;
- workspace-confined context packets containing exact routed protocol sections and loaded workspace inputs;
- unit tests, positive/negative bootstrap fixtures with stable diagnostics, and CI.

Explicitly unsupported or deferred:

- semantic execution of all 13 invocation modes, mutation authorization, automatic transitive read closure, stale-check enforcement, and `0G` transaction handling;
- complete canonicalization, schema/enum coverage, collision extension, specialized PRF/HBK/COV/CND identity, record/envelope/package hashing, and every complete Part 19.2 family;
- claims, source/frontier semantics, tickets, probes, reconciliation, coverage arithmetic, acquisition execution, synthesis, confirmations, handbook generation, packaging, Exit A, and Exit E;
- full Protocol v4 conformance and a real-system dry run.

Routing support is not mode-execution support. Bootstrap fixtures test only registered M0 primitives and do not satisfy a complete Part 19.2 family.

## Repository guide

**First-time users**

- [Minimal example](examples/minimal/README.md): shortest implemented flow.
- [Workspace guide](docs/workspace.md): the four planes and initialization boundary.
- [Development guide](docs/development.md): full CLI and local checks.

**Agents and contributors**

- [AGENTS.md](AGENTS.md): protocol-integrity and contribution rules.
- [Skill entry point](skill/SKILL.md) and [skill guide](skill/README.md): non-authoritative invocation routing.
- [Architecture](docs/architecture.md): current M0 packages and target boundaries.
- [Conformance](docs/conformance.md) and [fixture corpus](fixtures/README.md): implemented cases versus normative completeness.
- [Roadmap](docs/roadmap.md): staged M0–M15 delivery.
- [Schema projections](schemas/README.md): current schema status and future boundary.
- [Export-helper boundary](tools/export-helper/README.md): operator-only acquisition trust model.

## License

MIT. See [LICENSE](LICENSE).
