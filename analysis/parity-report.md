# Crown Parity Report

- Schema version: 1
- Generator: parityreport parityreport/0.1
- Generated at: 2026-09-25T00:00:00Z
- Protocol: 4.1.3 (sha256:561b5d78a8e062a89c79309dc2ae46f7ef12b6f91c573ad4b892195ec5a5e999)
- Language: cdl/0.3 (EIR format 2)

## Section coverage

| Section | Heading | CDL source | Status |
| :-- | :-- | :-- | :-- |
| 0 | Purpose, Guarantees, and Conceptual Flow |  | wrapper |
| 1.1 | Plane 1 — Forensic Plane | protocol/foundations.cdl | implemented |
| 1.2 | Plane 2 — Assurance Plane | protocol/foundations.cdl | implemented |
| 1.3 | Plane 3 — Synthesis Plane | protocol/foundations.cdl | implemented |
| 1.4 | Plane 4 — Human Reconstruction Handbook | protocol/foundations.cdl | implemented |
| 1.5 | Non-authoritative sidecars | protocol/foundations.cdl | implemented |
| 2.1 | Persona definition and purity | protocol/personas-povs.cdl | implemented |
| 2.2 | Canonical entry ownership | protocol/personas-povs.cdl | implemented |
| 2.3 | Five isolated POVs | protocol/personas-povs.cdl | implemented |
| 2.4 | Traversal tracks | protocol/personas-povs.cdl | implemented |
| 2.5 | Capability and runtime state | protocol/personas-povs.cdl | implemented |
| 2.6 | Capability grouping and target decision | protocol/personas-povs.cdl | implemented |
| 3.1 | `0A-PREFLIGHT.md` | protocol/preflight-registry.cdl | implemented |
| 3.2 | Foundational registries | protocol/workspace-registries.cdl | implemented |
| 4.1 | ID generation | protocol/typed-id.cdl | implemented |
| 4.1.1 | Semantic payload identity and certification envelopes | protocol/semantic-payload-identity.cdl | implemented |
| 4.1.2 | Normative canonical hash profile and packaging artifacts | protocol/canonical-hash-profile.cdl | implemented |
| 4.2 | Atomic unit rule | protocol/inventory-and-frontier.cdl | implemented |
| 4.3 | Source inventory denominator — `10-SOURCE-INVENTORY.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.4 | Traversal frontier — `11-TRAVERSAL-FRONTIER.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.5 | Scope-specific source coverage — `16-SOURCE-COVERAGE.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.6 | Stable supersession and tombstones | protocol/inventory-and-frontier.cdl | implemented |
| 5.1 | Dedicated prefix groups | protocol/status-taxonomy.cdl | implemented |
| 5.2 | Claim-level evidence — `12-CLAIM-EVIDENCE.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.3 | Sanitized projections and helper limitations | protocol/evidence-and-decisions.cdl | implemented |
| 5.4 | Contradictions — `14-CONTRADICTIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.5 | Authoritative decisions — `15-DECISIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.6 | Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.7 | Human confirmations — `18-CONFIRMATIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
| 6.1 | Atomic component schema | protocol/extraction-evolution.cdl | implemented |
| 6.2 | Dual-entry discovery | protocol/extraction-evolution.cdl | implemented |
| 6.3 | Staging and promotion | protocol/extraction-evolution.cdl | implemented |
| 6.4 | Depromotion | protocol/extraction-evolution.cdl | implemented |
| 6.5 | Surgical evolution | protocol/extraction-evolution.cdl | implemented |
| 7.1 | Trust boundary | protocol/acquisition-trust.cdl | implemented |
| 7.2 | Acquisition register | protocol/acquisition-trust.cdl | implemented |
| 7.3 | Logical explosion and virtual coordinates | protocol/acquisition-trust.cdl | implemented |
| 7.4 | Required explosion coverage | protocol/acquisition-trust.cdl | implemented |
| 7.5 | n8n human acquisition loop | protocol/export-acquisition-loop.cdl | implemented |
| 7.6 | Normalized maps | protocol/normalized-maps.cdl | implemented |
| 7.7 | Export reconciliation — `13-EXPORT-RECONCILIATION.md` | protocol/export-reconciliation.cdl | implemented |
| 8.1 | Resume identity header | protocol/invocation-context.cdl | implemented |
| 8.2 | Strict single-scope modes | protocol/invocation-context.cdl | implemented |
| 8.3 | Invocation-mode enum and explicit multi-file modes | protocol/invocation-context.cdl | implemented |
| 8.4 | Mandatory read sets | protocol/invocation-context.cdl | implemented |
| 8.5 | Cold resume check | protocol/cold-resume.cdl | implemented |
| 8.6 | Concurrent persona traversal | protocol/invocation-lifecycle.cdl | implemented |
| 8.7 | Stale checkpoint guard | protocol/invocation-lifecycle.cdl | implemented |
| 8.8 | `0G` invocation log | protocol/invocation-lifecycle.cdl | implemented |
| 9.1 | Ticket schema and canonical IDs | protocol/ticket-fsm.cdl | implemented |
| 9.2 | FSM and write rights | protocol/ticket-fsm.cdl | implemented |
| 9.3 | Human Hatch action prerequisites | protocol/ticket-fsm.cdl | implemented |
| 9.4 | Probe specification and asynchronous handoff | protocol/ticket-fsm.cdl | implemented |
| 9.5 | No-mock integrity | protocol/ticket-fsm.cdl | implemented |
| 9.6 | No-mock fallback payload | protocol/ticket-fsm.cdl | implemented |
| 10.1 | Evidence-backed `[R-SWEPT]` | protocol/coverage-and-exits.cdl | implemented |
| 10.2 | Per-kind coverage arithmetic | protocol/coverage-and-exits.cdl | implemented |
| 10.3 | Composite Exit A — Deconstruction Closure | protocol/coverage-and-exits.cdl | implemented |
| 10.4 | Exit B — Stagnation | protocol/coverage-and-exits.cdl | implemented |
| 10.5 | Exit C — Oscillation | protocol/coverage-and-exits.cdl | implemented |
| 10.6 | Exit D — Limit exhaustion and deterministic iteration accounting | protocol/coverage-and-exits.cdl | implemented |
| 11.1 | Partial versus Final Synthesis | protocol/synthesis-catalogs.cdl | implemented |
| 11.2 | Common synthesis block header | protocol/synthesis-catalogs.cdl | implemented |
| 11.3 | Architecture blueprint — `90-ARCH-BLUEPRINT.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.4 | Entity — `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.5 | Relationship — `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.6 | State machine — `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.7 | Database routine — `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.8 | Business rule — `92-BUSINESS-RULES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.9 | Use case — `93-USE-CASES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.10 | Interface — `94-INTERFACES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.11 | Deployment, configuration, and scheduling — `95-DEPLOYMENT.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.12 | NFR and security — `96-NON-FUNCTIONAL-SECURITY.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.13 | Persona profile | protocol/synthesis-catalogs.cdl | implemented |
| 12.1 | Traceability — `20-TRACEABILITY.md` | protocol/traceability.cdl | implemented |
| 12.2 | Required reciprocal references | protocol/traceability.cdl | implemented |
| 12.3 | Synthesis-ID timing | protocol/traceability.cdl | implemented |
| 12.4 | Dependency and impact provenance | protocol/traceability.cdl | implemented |
| 12.5 | Deterministic validation summary | protocol/traceability.cdl | implemented |
| 13.1 | Audience paths and progressive disclosure | protocol/handbook-and-decisions.cdl | implemented |
| 13.2 | Required chapter behavior | protocol/handbook-and-decisions.cdl | implemented |
| 13.3 | Disabled and Dormant Features chapter | protocol/handbook-and-decisions.cdl | implemented |
| 13.4 | Handbook quality gate | protocol/handbook-and-decisions.cdl | implemented |
| 14.1 | Decision log | protocol/handbook-and-decisions.cdl | implemented |
| 14.2 | Modernization mapping templates | protocol/handbook-and-decisions.cdl | implemented |
| 14.3 | Business-rule ownership | protocol/handbook-and-decisions.cdl | implemented |
| 14.4 | Interface compatibility | protocol/handbook-and-decisions.cdl | implemented |
| 14.5 | Equivalence and acceptance suite | protocol/handbook-and-decisions.cdl | implemented |
| 14.6 | Sole-input invariant | protocol/handbook-and-decisions.cdl | implemented |
| 15.1 | Strict conditions | protocol/packaging.cdl | implemented |
| 15.1.1 | Universal deterministic packaging rules | protocol/packaging.cdl | implemented |
| 15.1.2 | Exit E Candidate Report ordered payload schema | protocol/packaging.cdl | implemented |
| 15.1.3 | Candidate Payload Manifest ordered payload schema | protocol/packaging.cdl | implemented |
| 15.1.4 | Exit E Content-Readiness Report ordered payload schema | protocol/packaging.cdl | implemented |
| 15.2 | Scope certificate | protocol/packaging.cdl | implemented |
| 15.3 | Package contents | protocol/packaging.cdl | implemented |
| 15.4 | Acyclic certification sequence and hash domains | protocol/packaging.cdl | implemented |
| 15.5 | Reproducible final-bundle verification | protocol/packaging.cdl | implemented |
| 17.1 | Conforming implementation capabilities | protocol/conformance.cdl | implemented |
| 17.2 | Normative conformance cases | protocol/conformance.cdl | implemented |
| 17.3 | Corpus acceptance | protocol/conformance.cdl | implemented |

## Checks

| Check | Status | Detail |
| :-- | :-- | :-- |
| assembly-compiles | pass | 27 sources |
| drift-tests | pass | 23 tests in 23 files |
| generated-protocol | pass | sha256:cbe7f2a768b635bbf410498e4e2e8dce6c31e91b26e2a3c76c207cf15e42b3e8 |
| golden-eir | pass | sha256:35db4396e0dfcfbf653041bc2bb88d221af86354f5981d542973e9a092fe58ac |
| golden-prompt | pass | sha256:469688636148b488fce03232c36ab7c64439fca83db2b43b7a339b7845c63db5 |
| fixture-oracle | pass | 24/24 passed |
| identity-ledger | pass | 536 identities across 27 sources |

## Fixtures

24/24 passed, 0 failed.

## Sources

| Source | Sections | Identities | Fingerprint |
| :-- | :-- | --: | :-- |
| protocol/globals.cdl | 0 | 0 | sha256:f7f351b217c8c2ae284ebb06ea7a5f49b5451c0643c0aaf453680cc24fb7d53b |
| protocol/invocation-modes.cdl | 0 | 0 | sha256:c599bd58c22d06daf961ddd6952a1d5e78598178f762efc7c2242f792ed75b57 |
| protocol/foundations.cdl | 9 | 30 | sha256:c0f47838ac1040d0c357f3fc5789bdc29a556a30b472ba2d65bc06d417565610 |
| protocol/personas-povs.cdl | 6 | 18 | sha256:4337b5ac7044015f99b57dfc7f4d911180e2ed0c00a0b9f12a2fe81a04ad7d9b |
| protocol/preflight-registry.cdl | 1 | 9 | sha256:53bc2e4845b667aa6edcf2018765efdbca89e77a6dd49d4ee51afb0f75f25115 |
| protocol/workspace-registries.cdl | 1 | 9 | sha256:e14fcf5bef07a0eb24bce78f6dc89057fbc29ecee4ba5eadf2b9beebfb5443cc |
| protocol/typed-id.cdl | 1 | 9 | sha256:ea12122d8eee3e1bcc35d7b7f8198510215aef1ea6e3c276a6d26f9bc59c5730 |
| protocol/semantic-payload-identity.cdl | 1 | 7 | sha256:07747c57e064fc03ec144bec16d87a8be378b613990e2711caa7e0742e2579d9 |
| protocol/canonical-hash-profile.cdl | 1 | 13 | sha256:a5712ff3df7c69210cd914ddb769c75ce746d128a5aadbaa59da5becf487f200 |
| protocol/inventory-and-frontier.cdl | 5 | 23 | sha256:c78e877ce0b1e919dc1eabe424f9c3df2ebac1df19a26af02ce1fbf541b05437 |
| protocol/status-taxonomy.cdl | 1 | 112 | sha256:a4f149fcb7777b4474ca0181fc16c21c467e7d035de117aa2d9f529bea0611b9 |
| protocol/evidence-and-decisions.cdl | 6 | 13 | sha256:c5a90679f6dff72b9c5c906f13065ceb7faae8802c4f13cffbb29faa62b75b2e |
| protocol/extraction-evolution.cdl | 5 | 15 | sha256:3ed11f0af3a85ac3e31d19b339376c3ad4ad4ce6d66cd5246a8449e15aeb239b |
| protocol/acquisition-trust.cdl | 4 | 10 | sha256:3df535ce5dc7213a68c728f99750e82fd9a4ad4add1408057acc663c7fc004a4 |
| protocol/export-acquisition-loop.cdl | 1 | 10 | sha256:67ef3cacc4057f22318775d9d3ecaffc62e760dcbd92eecc0c37112bd0083962 |
| protocol/normalized-maps.cdl | 1 | 4 | sha256:909965bb69994fa0b11f24ec99f6b9081f7de558f16f82972849e09730b067ba |
| protocol/export-reconciliation.cdl | 1 | 22 | sha256:448a3605fcb7470a17c98a77f6a00b7d7c78386323edbec7632df9ccda546d73 |
| protocol/invocation-context.cdl | 4 | 17 | sha256:f10345d91db0d68df07885dac803a6b6ee9a53ce339457b8739fcca0a1db9613 |
| protocol/cold-resume.cdl | 1 | 4 | sha256:97410ec60b1e39177da5b3732372354becb46a39ba144b9558b3334aa28fd350 |
| protocol/invocation-lifecycle.cdl | 3 | 8 | sha256:111d8724fe190b872b60950ce1073ad150a6e6ae410a3d2296b6b7589ada3895 |
| protocol/ticket-fsm.cdl | 6 | 20 | sha256:caec43e35f16d1fe2f820e6dc22e7c338d514eefd6148b4e0d688405cf793c2a |
| protocol/coverage-and-exits.cdl | 6 | 19 | sha256:ae942630aea2aad82b4bce1afd7f84b3a0bc19b24d42fe7f7d13210aedeac7ef |
| protocol/synthesis-catalogs.cdl | 13 | 66 | sha256:8484ade63bf6ba3ac89ae3cbc47b13f0abde5928b3295451367b25a9accf5028 |
| protocol/traceability.cdl | 5 | 15 | sha256:a376905e75a5ab3066d0755a9ec727626a94db12cc9d06504b8169656a78caf6 |
| protocol/handbook-and-decisions.cdl | 10 | 23 | sha256:46e416664481be6b52170df20268f2b1f060ab96e6ffa88ab2ab42be4755c321 |
| protocol/packaging.cdl | 9 | 50 | sha256:874196e684060da59fc7df7a6cca3253347b520e04df014ab504c36f166fb7f6 |
| protocol/conformance.cdl | 5 | 10 | sha256:b45b409c398f209b7e72d894108b741ebd4da2272f2871f7c9e5bde5f8977a4c |

## Declared deferrals

- standalone protocol-quality test on the generated protocol.md edition is pending (content completeness is enforced by TestGeneratedProtocolCoversOracle)
- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

The CDL sources under `protocol/` are the normative authority; `protocol.md` is a generated render (fingerprint above). `protocol/legacy/protocol-4.1.3.md` is the frozen bootstrap parity oracle.
