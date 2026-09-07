# Agent contribution rules

This is operational guidance, not a protocol specification. [`protocol.md`](protocol.md) is normative and wins every conflict.

## Protocol integrity and scope

- Treat `protocol.md` as read-only by default. Change it only for a demonstrated protocol defect, never to simplify implementation; record proposed issues separately first.
- Before changing normative behavior, read the affected sections and identify them in the change description.
- Never weaken `MUST`/`MUST NOT`, duplicate normative rules into secondary sources, or present a projection as authority.
- Keep work inside the requested milestone and feature boundary. Record adjacent protocol work in the roadmap rather than expanding scope opportunistically.

## Deterministic implementation

- Keep semantic interpretation, claim formulation, and ambiguity recognition with the LLM; put identity, parsing, scope, ordering, canonicalization, hashing, validation, and state transitions in deterministic code.
- Parse protocol-significant Markdown structurally. Whole-document regex is not a substitute for Markdown-AST semantics.
- Treat conversation memory as disposable and `.extracted/` filesystem state as persistent execution state.
- Unsupported behavior must fail or report `unsupported`; never add placeholder success paths.

## Artifacts and trust boundaries

- Never commit real `.extracted/` workspaces, raw acquisition data, probe output, context packets, credentials, secrets, production data, private matching hints, or operator-private material. Fixtures and examples must remain synthetic and non-sensitive.
- Generated derivatives are not semantic authority. Future checked-in generated artifacts must identify their generator/version, source inputs and fingerprints, schema version, and generation time; regenerate them deterministically instead of hand-editing them.
- Hand-maintained skill projections and routing metadata are not generated sidecars. Keep their non-authoritative status explicit and validate them against `protocol.md`.

## Fixtures and claims

- Every newly implemented normative rule requires traceability to exact protocol sections plus positive and negative fixtures with stable expected results or diagnostics.
- Do not claim an implemented Part 19 family, Exit A, Exit E, or Protocol v4 conformance until the complete applicable deterministic validators and normative fixtures pass.

Before finishing, run the narrowest relevant commands from [docs/development.md](docs/development.md). Report checks run, protocol sections affected, implemented support, and deferred behavior.
