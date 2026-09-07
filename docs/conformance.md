# Conformance

> Non-authoritative implementation policy. Protocol §19.2 is the normative fixture requirement and wins on conflict.

Protocol v4 conformance requires every applicable positive and negative fixture in §19.2. A passing unit test, file-presence check, bootstrap case, or unsupported/skipped family is not conformance and cannot establish a gate result.

## Fixture contract

Each implemented normative rule needs:

- positive acceptance and negative fail-closed coverage;
- a stable expected result or diagnostic identity;
- provenance to the exact protocol heading and implementation capability;
- deterministic, synthetic inputs containing no raw exports, credentials, secrets, production identifiers, private hints, or production data.

The M0 runner validates registered operation/group pairs and unique protocol-heading provenance. Current operations compare stable string results or diagnostic classes. An unrelated error does not satisfy a negative case. CI runs every registered bootstrap case, and human-readable CLI/workflow output lists complete Part 19.2 families as unsupported rather than counting them as passing.

## Implemented bootstrap cases

Positive and negative cases currently cover:

- structural discovery of protocol version and invocation-mode count;
- foundational 12-character typed IDs;
- relative-path normalization and traversal rejection;
- the explicitly limited basic Markdown hash primitive;
- generated workspace-skeleton structure.

A case's protocol heading is provenance, not a claim that the case implements that complete section. No complete §19.2 family is implemented.

## Unsupported complete families

1. persona workspace layout and ownership;
2. context-qualified finite enums and bidirectional registry completeness;
3. PRF identity and semantic hashing;
4. HBK identity and path/anchor normalization;
5. Profile Synchronization dependency closure;
6. record-versus-artifact hashing and file-transport separation;
7. canonicalization and exact carrier/envelope exclusions;
8. final envelope and package-member integrity;
9. all five packaging schemas and deterministic ordering;
10. the acyclic Exit E sequence.

The repository supports neither Exit A nor Exit E and makes no Protocol v4 conformance claim. A complete family becomes implemented only when its production behavior and all required positive/negative fixtures run successfully in CI.
