# Legacy Autopsy Invocation Router

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION**

`protocol.md` is the sole normative authority. This skill routes one Protocol v4.0 invocation for an agent/harness; it does not define protocol behavior. If this file, `modes.json`, or a mode projection conflicts with `protocol.md`, follow `protocol.md`.

**M0 support boundary:** the current CLI validates routing and assembles bounded context packets. It does not execute semantic mode effects, authorize mutations, append `0G`, enforce complete cold resume/staleness, or run protocol gates. The sequence below is the target invocation contract an executing agent/harness must obey as those capabilities are implemented.

Never execute protocol rules from remembered summaries when `protocol.md` is available. Load the authoritative sections needed for the current invocation.

## State model

```text
conversation memory = disposable
.extracted/ = persistent execution state
protocol.md = normative execution law
```

Reconstruct invocation state from the filesystem. Routing metadata and conversation summaries may locate authoritative material but never replace it.

## Target invocation sequence

1. Determine the requested Invocation Mode.
2. Resolve the relevant normative protocol sections.
3. Load the required workspace/checkpoint context.
4. Construct or validate the Resume Identity Header.
5. Verify allowed read/write scope.
6. Execute exactly one protocol invocation.
7. Invoke deterministic helpers/validators where applicable.
8. Produce required forensic/assurance side effects.
9. Append `0G` invocation history.
10. Stop when that invocation's responsibility is complete.

Do not silently chain modes. Before mutation, apply §8.1 `Resume identity header`, §8.4 `Mandatory read sets`, §8.5 `Cold resume summary`, and §8.7 `Stale checkpoint guard`. Enumerate exact targets for every multi-file mode. After an implemented invocation, apply §8.8 `` `0G` invocation log `` and stop.

## Mode routing

Resolve the exact §8.3 mode through `modes.json`, then load its projection and every cited normative section:

- `Preflight`
- `Export Acquisition`
- `Discovery`
- `Ticket Resolution`
- `Promotion Review`
- `Sequential Reconciliation`
- `Profile Synchronization`
- `Cross-Reference-Reconciliation`
- `Partial-Synthesis`
- `Final-Synthesis`
- `Reconstruction-Handoff`
- `Human Hatch`
- `Validation/Gate`

Mode projections are navigation aids, not permission or implementation-support claims. Generated context must carry the exact protocol sections routed by the manifest. A mode name never grants undeclared read or write access.

## Deterministic boundary

Use deterministic helpers for identity, canonicalization, hashing, schemas, enums, ordering, allocation, references, coverage/frontier arithmetic, fingerprints, stale checks, gates, and packaging only as they are implemented. Unsupported checks remain explicitly unsupported; never infer a pass. Semantic interpretation stays bounded by loaded protocol text and evidence.
