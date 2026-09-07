# Ticket Resolution

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION** — `protocol.md` is normative.

## Purpose
Apply one legal ticket-FSM action, including probe ingestion, with surgical updates to its bound scope.

## Required identity bindings
Bind exactly one persona/prefix, cluster, track, POV, POV file, zero/one ledger, one ticket, snapshot/environment, and actor origin/target role. Neutral arbitration retains original scope and records the neutral POV.

## Normative protocol sections
- §5.4. `Contradictions — `14-CONTRADICTIONS.md``; §6.5. `Surgical evolution`
- §8.1. `Resume identity header`; §8.2. `Strict single-scope modes`; §8.3. `Invocation-mode enum and explicit multi-file modes`; §8.4. `Mandatory read sets`; §8.5. `Cold resume summary`; §8.7. `Stale checkpoint guard`; §8.8. `` `0G` invocation log ``
- §9.1. `Ticket schema and canonical IDs`; §9.2. `FSM and write rights`; §9.4. `Probe specification and asynchronous handoff`; §9.5. `No-mock integrity`; §9.6. `No-mock fallback payload`
- §12.5. `Deterministic validation summary`

## Mandatory read set
Load §8.4 strict-scope inputs, ticket ledger/history, affected records/claims/frontiers/contradictions, static source, and matching scrubbed fingerprinted probe log when applicable.

## Semantic write targets
The bound ledger transition and only affected assumptions/records in the bound POV.

## Allowed assurance/audit side effects
Directly produced claims, contradictions, inventory/frontier/coverage rows, local buffers, and append-only `0G`.

## Forbidden mutations
No illegal actor transition, direct `[T-HUMAN-REQUIRED]`, guessing, fabricated mocks/evidence, unsanitized probe evidence, another POV’s semantics, clean-sweep rewrite, or terminal-ticket-as-coverage inference.

## Deterministic validation
Validate ticket allocation/uniqueness, FSM and actor rights, static investigation, probe authorization/safety/fingerprint/sanitization, evidence class, surgical versioning, and stale effects.

## Completion condition
Exactly one valid ticket action and corresponding evidence/affected-record updates commit; unresolved work remains in its correct active state or is ready for authorized Human Hatch handling.
