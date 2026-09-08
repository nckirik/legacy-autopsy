# Legacy Autopsy Invocation Router

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION**

`protocol.md` is the sole normative authority. This skill routes one Protocol v4.0 invocation for an agent/harness; it does not define protocol behavior. If this file, `modes.json`, or a mode projection conflicts with `protocol.md`, follow `protocol.md`.

**M0 support boundary:** the current CLI validates routing and assembles bounded context packets. It does not discover or start a service, attach to an autopsy, claim scheduled work, execute semantic mode effects, authorize mutations, append `0G`, enforce complete cold resume/staleness, or run protocol gates.

The target skill is the first and initially exclusive semantic-work ingress to the local Legacy Autopsy service. It must work from any compatible coding harness without a harness-specific adapter. The service owns scheduling, leases, bounded context, stale checks, validation, workspace commits, and state transitions; the host coding harness owns only the semantic work for one claimed invocation. The UI and CLI cannot manually start or claim semantic work in this stage.

Named harness adapters are a later managed-runner layer, and the minimal internal direct-model runner comes last. Both must reuse the skill-proven logical invocation/result contract without gaining broader authority.

Never execute protocol rules from remembered summaries when `protocol.md` is available. Load the authoritative sections needed for the current invocation.

## State model

```text
conversation memory = disposable
.extracted/ = persistent execution state
protocol.md = normative execution law
service/UI state = operational, not semantic authority
```

Reconstruct protocol state from the filesystem. Routing metadata, service indexes, and conversation summaries may locate authoritative material but never replace it.

## Target service-backed sequence

1. Detect or start the configured local service and wait for readiness.
2. Identify or attach to the autopsy bound to the current repository and pinned snapshot.
3. Ask the service for the next permitted invocation and claim exactly one task.
4. Retrieve its bounded context, exact normative sections, identity, scope, and fingerprints.
5. Execute only that semantic invocation in the coding harness.
6. Submit the structured result to the service.
7. Let the service revalidate stale inputs, ownership, scope, and implemented rules before commit.
8. Surface committed, rejected, stale, or unsupported status without inferring success.
9. Repeat only by claiming another permitted invocation.

Every adapter operation must carry explicit `autopsy_id`; every invocation still obeys the complete protocol identity and write restrictions. Do not silently chain modes or scopes. A mode name never grants undeclared access. For every started invocation, including one whose semantic result is rejected, stale, blocked, or unsupported, the service applies the mandatory §8.8 `` `0G` invocation log `` with exact blockers and next loads. Authorized semantic effects and the audit append commit under the applicable transaction boundary; derivatives and events update only afterward.

See [`docs/skill-runtime.md`](../docs/skill-runtime.md) for the detailed ownership split and [`docs/runtime.md`](../docs/runtime.md) for the target service/workbench design.

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
