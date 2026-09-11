# Conformance

> Non-authoritative Legacy Autopsy implementation policy. Protocol §19.2 defines normative conformance cases and wins on conflict; this repository's JSON files, Go runner, CLI output, and CI workflow are one concrete test realization, not protocol prerequisites.

Protocol v4.1.1 conformance requires satisfying every applicable positive and negative case in §19.2. A passing Legacy Autopsy unit test, file-presence check, bootstrap fixture, or unsupported/skipped family is not conformance and cannot establish a gate result.

## Standalone protocol quality test

Give a fresh capable coding harness only `protocol.md`, the target repository and approved evidence roots, and a short instruction to execute the protocol. Without Legacy Autopsy context, service APIs, CLI commands, UI, Atlas, hidden validators, repository-specific fixtures, or prior conversation, it should be able to identify required inputs, create and resume `.extracted/`, perform semantic work, implement or reproduce deterministic operations with ordinary local tools, ask required human questions, and report unsupported behavior rather than invent success.

This is a protocol-quality regression test, not by itself proof that an implementation is conforming or that a workspace has passed Exit A or Exit E.

## Legacy Autopsy fixture contract

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
- prefix-qualified persona-directory basename parsing and reserved `_shared` rejection;
- the explicitly limited basic Markdown hash primitive;
- generated workspace-skeleton structure.

A case's protocol heading is provenance, not a claim that the case implements that complete section. No complete §19.2 family is implemented.

## Unsupported complete families

1. persona workspace layout and ownership;
2. context-qualified finite enums and bidirectional registry completeness;
3. invocation ownership and cold resume;
4. ticket escalation and honest unresolved coverage;
5. dead-code and semantic-predicate closure;
6. deterministic iteration accounting;
7. PRF identity and semantic hashing;
8. HBK identity and path/anchor normalization;
9. Profile Synchronization dependency closure;
10. record-versus-artifact hashing and file-transport separation;
11. canonicalization and exact carrier/envelope exclusions;
12. final envelope and package-member integrity;
13. packaging schemas, gate completeness, evidence bindings, and deterministic ordering;
14. the acyclic, snapshot-consistent Exit E sequence and reproducible final verification.

The repository supports neither Exit A nor Exit E and makes no Protocol v4.1.1 conformance claim. A complete family becomes implemented only when its production behavior and all required positive/negative fixtures run successfully in CI.
