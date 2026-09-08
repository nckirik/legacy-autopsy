# Skill adapter and executor invocation lifecycle

> Non-authoritative operational guide. It routes to [`protocol.md`](../protocol.md); it does not reproduce or replace the protocol.

## M0 support boundary

M0 validates the official mode registry, skill projections, identity inputs, confined workspace reads, and exact protocol-section assembly. It does **not** start or attach to a service, claim scheduled work, execute semantic mode effects, authorize writes, append `0G`, prove a checkpoint current, perform automatic transitive closure, or run Exit A/Exit E.

The current skill is a router/context guide. The service-backed adapter and executing lifecycle below are target behavior, not implemented capability.

## Target ownership

> **Legacy Autopsy owns orchestration and protocol state. Executors own semantic work.**

The local service selects permitted work, owns leases and scheduling, assembles bounded context, validates submitted results, commits authorized workspace transactions, updates projections, and emits events. In the first user-facing stage, a generic skill invoked inside any compatible coding harness is the only semantic-work ingress. The host harness performs one invocation and returns a result; the UI and CLI cannot manually start or claim semantic work.

Thin service-managed adapters for named coding harnesses come second. A minimal internal direct-model runner comes last. Every stage reuses the same logical invocation/result contract, and executor choice cannot weaken protocol identity, scope, evidence, ownership, stale-input, or write restrictions.

## Future bring-your-own-harness flow

When invoked through a user's coding harness, the thin skill adapter should:

1. detect whether the local Legacy Autopsy service is ready;
2. start it in the background when configured and needed, then wait for readiness;
3. identify or attach to the autopsy bound to the current repository and snapshot;
4. ask the service for the next permitted invocation;
5. claim exactly one invocation and retrieve its bounded, fingerprinted context;
6. perform only that semantic task in the coding harness;
7. submit the structured result to the service for validation and commit;
8. continue only by claiming another permitted invocation.

Every operation is scoped by explicit `autopsy_id` plus rooted repository/workspace bindings. The service may reject a result as stale, unauthorized, invalid, or unsupported. The adapter must surface that result rather than infer success or edit workspace state around it.

Named adapter settings introduced in the second stage belong under the project's reserved `.legacy-autopsy/` operational root. They may select a harness, executable, limits, and non-secret preferences; they are not evidence or protocol state. Configuration format remains undecided, mutable machine state must not be committed, and credentials remain outside the repository and `.extracted/`.

## One-invocation lifecycle

For each claimed task, the combined service/executor flow must:

1. **Select the invocation mode.** Resolve the task to one official §8.3 mode.
2. **Route authoritative sections.** Load the minimum sufficient sections from `protocol.md`; never execute remembered summaries while the protocol is available.
3. **Load persistent state.** Read the complete mode-specific workspace set and active source/snapshot inputs.
4. **Construct resume identity.** Build or validate the §8.1 header and §8.5 cold-resume summary.
5. **Verify scope.** Check planned reads and writes against autopsy identity, protocol identity, ownership, and strict single-scope constraints.
6. **Execute once.** The executor performs only the selected semantic invocation; a mode name is not broad write permission.
7. **Validate deterministically.** The service applies implemented parsing, identity, ordering, hashing, stale, schema, and authorization checks; unsupported checks block.
8. **Commit the invocation record.** For every started invocation, atomically append the mandatory `0G` entry with mutations, blockers, and next loads. Apply authorized semantic effects in the applicable transaction; a rejected, stale, blocked, or unsupported semantic result commits no semantic mutation but still records the invocation.
9. **Project and stop.** Update derivative indexes after the final workspace transaction, emit events, release the claim, and do not silently roll into another scope.

```text
conversation memory = disposable
.extracted/ = persistent execution state
protocol.md = normative execution law
service scheduling/UI state = operational, not semantic authority
```

Filesystem contradictions create records or blockers; conversation memory cannot settle them. M0 context packets quote exact sections routed by `skill/modes.json` and embed loaded workspace bytes and fingerprints. Treat those packets as potentially sensitive operator material, not repository documentation or a sanitized export.

See [runtime.md](runtime.md) for the multi-autopsy service, executor, scheduler, workbench, and Atlas target design.
