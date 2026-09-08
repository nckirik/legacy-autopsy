# Skill layer

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION**

This directory routes an agent/harness to the normative root [`protocol.md`](../protocol.md). M0 implements routing and context assembly, not semantic execution of the 13 modes. Its first target role is the generic service-backed skill that works inside any compatible coding harness and is initially the only semantic-work ingress. Named harness adapters follow later; the minimal internal direct-model runner comes last. The service—not the skill, UI, CLI, or executor—owns scheduling, validation, workspace commits, and protocol/runtime state.

- `SKILL.md` describes the future service-backed claim/context/submit flow and current support boundary.
- `modes.json` is strict JSON routing metadata, parsed by Go `encoding/json`.
- `modes/` contains exactly one concise navigation projection for each official §8.3 invocation mode.

Start with `SKILL.md`, select the exact mode identity in `modes.json`, open its projection, and load the cited sections from `protocol.md`. The manifest and projections are hand-maintained indexes: they are neither generated sidecars nor normative behavior, and they grant no permissions. The routing validator checks their structural relationship to the protocol; it cannot prove the prose semantically equivalent.

Mode filenames are lowercase slugs; each file's H1 and manifest `mode` preserve the exact §8.3 identity, including spaces, hyphens, and `/`.
