# Skill layer

> **NON-AUTHORITATIVE OPERATIONAL PROJECTION**

This directory routes an agent/harness to the normative root [`protocol.md`](../protocol.md). M0 implements routing and context assembly, not semantic execution of the 13 modes.

- `SKILL.md` describes the target one-invocation procedure and current support boundary.
- `modes.json` is strict JSON routing metadata, parsed by Go `encoding/json`.
- `modes/` contains exactly one concise navigation projection for each official §8.3 invocation mode.

Start with `SKILL.md`, select the exact mode identity in `modes.json`, open its projection, and load the cited sections from `protocol.md`. The manifest and projections are hand-maintained indexes: they are neither generated sidecars nor normative behavior, and they grant no permissions. The routing validator checks their structural relationship to the protocol; it cannot prove the prose semantically equivalent.

Mode filenames are lowercase slugs; each file's H1 and manifest `mode` preserve the exact §8.3 identity, including spaces, hyphens, and `/`.
