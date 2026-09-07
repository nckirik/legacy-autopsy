# Stateless skill runtime

> Non-authoritative operational guide. It routes to `protocol.md`; it does not reproduce or replace the protocol.

## M0 support boundary

M0 validates the official mode registry, skill projections, identity inputs, confined workspace reads, and exact protocol-section assembly. It does **not** execute semantic mode effects, authorize writes, append `0G`, prove a checkpoint current, perform automatic transitive closure, or run Exit A/Exit E. The lifecycle below is the target contract for future executing agents and harness services.

## Target invocation lifecycle

1. **Select the invocation mode.** Resolve the request to one official §8.3 mode.
2. **Route authoritative sections.** Load the minimum sufficient sections from `protocol.md`; never execute remembered summaries while the protocol is available.
3. **Load persistent state.** Read the complete mode-specific workspace set and active source/snapshot inputs.
4. **Construct resume identity.** Build or validate the §8.1 header and §8.5 cold-resume summary.
5. **Verify scope.** Check planned reads and writes against identity, ownership, and strict single-scope constraints.
6. **Execute once.** Perform only the selected invocation; a mode name is not broad write permission.
7. **Use deterministic helpers.** Delegate implemented parsing, identity, ordering, hashing, validation, and atomic mechanics to code; unsupported helpers block rather than being guessed.
8. **Commit protocol effects.** Apply only authorized plane-specific and audit effects, atomically where required.
9. **Append `0G`.** Record identity, loaded fingerprints, mutations, blockers, and next loads.
10. **Stop.** Do not roll into another mode or scope.

```text
conversation memory = disposable
.extracted/ = persistent execution state
protocol.md = normative execution law
```

Filesystem contradictions create records or blockers; conversation memory cannot settle them. M0 context packets quote exact sections routed by `skill/modes.json` and embed loaded workspace bytes and fingerprints. Treat those packets as potentially sensitive operator material, not repository documentation or a sanitized export.
