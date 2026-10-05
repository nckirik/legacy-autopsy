# Crown Parity Report

- Schema version: 1
- Generator: parityreport parityreport/0.1
- Generated at: 2026-09-25T00:00:00Z
- Protocol: 4.2.1 (sha256:fb0d55f28bfefc6bd12e6eb36ebae1ae95d81f3cf119ca136816be808b244189)
- Language: cdl/0.4 (EIR format 4)

## Section coverage

| Section | Heading | CDL source | Status |
| :-- | :-- | :-- | :-- |
| 0.1 | Normative language | protocol/foundations.cdl | implemented |
| 0.2 | Guarantees | protocol/foundations.cdl | implemented |
| 0.3 | Non-guarantees | protocol/foundations.cdl | implemented |
| 0.4 | Conceptual flow | protocol/foundations.cdl | implemented |
| 1.1 | Plane 1 - Forensic Plane | protocol/foundations.cdl | implemented |
| 1.2 | Plane 2 - Assurance Plane | protocol/foundations.cdl | implemented |
| 1.3 | Plane 3 - Synthesis Plane | protocol/foundations.cdl | implemented |
| 1.4 | Plane 4 - Human Reconstruction Handbook | protocol/foundations.cdl | implemented |
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
| 4.3 | Source inventory denominator - `10-SOURCE-INVENTORY.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.4 | Traversal frontier - `11-TRAVERSAL-FRONTIER.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.5 | Scope-specific source coverage - `16-SOURCE-COVERAGE.md` | protocol/inventory-and-frontier.cdl | implemented |
| 4.6 | Stable supersession and tombstones | protocol/inventory-and-frontier.cdl | implemented |
| 5.1 | Dedicated prefix groups | protocol/status-taxonomy.cdl | implemented |
| 5.2 | Claim-level evidence - `12-CLAIM-EVIDENCE.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.3 | Sanitized projections and helper limitations | protocol/evidence-and-decisions.cdl | implemented |
| 5.4 | Contradictions - `14-CONTRADICTIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.5 | Authoritative decisions - `15-DECISIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.6 | Acquisition-candidate reconciliation - `17-ACQUISITION-CANDIDATES.md` | protocol/evidence-and-decisions.cdl | implemented |
| 5.7 | Human confirmations - `18-CONFIRMATIONS.md` | protocol/evidence-and-decisions.cdl | implemented |
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
| 7.7 | Export reconciliation - `13-EXPORT-RECONCILIATION.md` | protocol/export-reconciliation.cdl | implemented |
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
| 10.3 | Composite Exit A - Deconstruction Closure | protocol/coverage-and-exits.cdl | implemented |
| 10.4 | Exit B - Stagnation | protocol/coverage-and-exits.cdl | implemented |
| 10.5 | Exit C - Oscillation | protocol/coverage-and-exits.cdl | implemented |
| 10.6 | Exit D - Limit exhaustion and deterministic iteration accounting | protocol/coverage-and-exits.cdl | implemented |
| 11.1 | Partial versus Final Synthesis | protocol/synthesis-catalogs.cdl | implemented |
| 11.2 | Common synthesis block header | protocol/synthesis-catalogs.cdl | implemented |
| 11.3 | Architecture blueprint - `90-ARCH-BLUEPRINT.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.4 | Entity - `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.5 | Relationship - `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.6 | State machine - `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.7 | Database routine - `91-DATA-MODEL.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.8 | Business rule - `92-BUSINESS-RULES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.9 | Use case - `93-USE-CASES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.10 | Interface - `94-INTERFACES.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.11 | Deployment, configuration, and scheduling - `95-DEPLOYMENT.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.12 | NFR and security - `96-NON-FUNCTIONAL-SECURITY.md` | protocol/synthesis-catalogs.cdl | implemented |
| 11.13 | Persona profile | protocol/synthesis-catalogs.cdl | implemented |
| 12.1 | Traceability - `20-TRACEABILITY.md` | protocol/traceability.cdl | implemented |
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
| 16.1 | Security, Privacy, and Operational Guardrails | protocol/conformance.cdl | implemented |
| 17.1 | Conforming implementation capabilities | protocol/conformance.cdl | implemented |
| 17.2 | Normative conformance cases | protocol/conformance.cdl | implemented |
| 17.3 | Corpus acceptance | protocol/conformance.cdl | implemented |
| 18.1 | Compact Artifact Ownership Matrix | protocol/conformance.cdl | implemented |

## Checks

| Check | Status | Detail |
| :-- | :-- | :-- |
| assembly-compiles | pass | 27 sources |
| drift-tests | pass | 24 tests in 24 files |
| generated-protocol | pass | sha256:fb0d55f28bfefc6bd12e6eb36ebae1ae95d81f3cf119ca136816be808b244189 |
| golden-eir | pass | sha256:d0638c1e72676d31e1f8e7ee616013dc53ed85697b657a45204be342b0e56868 |
| golden-prompt | pass | sha256:9d142db280e569c78737b13a2d354ac9c599a1f2aea8a08f00135dd71ec0b5ce |
| fixture-oracle | pass | 24/24 passed |
| identity-ledger | pass | 536 identities across 27 sources |

## Fixtures

24/24 passed, 0 failed.

## Sources

| Source | Sections | Identities | Fingerprint |
| :-- | :-- | --: | :-- |
| protocol/globals.cdl | 0 | 0 | sha256:29a4522f11e1da6eff75430b1b507b4486e4764ff2969f28df3c7a03e70be8e2 |
| protocol/invocation-modes.cdl | 0 | 0 | sha256:4a8970adcf98246d33bc61684d6ca954b3a9e1f2d1c5e4df62788095f5cd7f5b |
| protocol/foundations.cdl | 9 | 30 | sha256:e872677e91ed6a42be4037037497cc097f9875529e6e2e19492f7b83774c1e0c |
| protocol/personas-povs.cdl | 6 | 18 | sha256:3c2ddd4f7664e2b9e032dfa4c06b3e8043ff94b28a6e6ebc0e5fedfe7982dd5a |
| protocol/preflight-registry.cdl | 1 | 9 | sha256:4bd97105f0101ec55a50a2521d1794a35bb5968c3a89feb6fe8a83865b3387e9 |
| protocol/workspace-registries.cdl | 1 | 9 | sha256:c01cfabc67f97647d593f37897c9f5e39df3eab585207c2b2440feb67ea70436 |
| protocol/typed-id.cdl | 1 | 9 | sha256:cc2efdcefb008626001bba599c2202b6606fa64db5e1c85ecb4b8706cd93fb65 |
| protocol/semantic-payload-identity.cdl | 1 | 7 | sha256:76fabcee6b8afcf3c37907f03c65aad9929badd62c8a3a27420c41a64f6f99d9 |
| protocol/canonical-hash-profile.cdl | 1 | 13 | sha256:b9f59a77a8c0c0e5e5c649c2e8792193746fdb079cdaadc599ff0ab119e8686d |
| protocol/inventory-and-frontier.cdl | 5 | 23 | sha256:7aed82c1b17e2c06d0d4b7c849fecf73facfa36d3063607cafa31e17c6a32824 |
| protocol/status-taxonomy.cdl | 1 | 112 | sha256:f57baa323d813e24e2886f322967736e28e42f8e2817294b0cc50f7b7fd1cce3 |
| protocol/evidence-and-decisions.cdl | 6 | 13 | sha256:efea46abc017ca9302cc5c2c1c450a3006f3b252b20c412a2b44d0073e5b4068 |
| protocol/extraction-evolution.cdl | 5 | 15 | sha256:a71016d8659dabb0b6320076bbe032065635338b305ec3531647c61ad1c7cc08 |
| protocol/acquisition-trust.cdl | 4 | 10 | sha256:049bd75c300689cee7f69fba81e0889d534e7a3e45a6076cc03a68e3e88f0f0a |
| protocol/export-acquisition-loop.cdl | 1 | 10 | sha256:3b4657ae694ade4a3325651b6256443f5c3eb233a50af73b31308dc2fd51baf7 |
| protocol/normalized-maps.cdl | 1 | 4 | sha256:967e4904801286ad48784e03ccfdf33ec9b671c3abfd54141546f8992172c440 |
| protocol/export-reconciliation.cdl | 1 | 22 | sha256:f030ff80bc682f12ddfcdd970741c021e16f17e68070c2a234e901c729f493f3 |
| protocol/invocation-context.cdl | 4 | 17 | sha256:f71efc685b1b5f9b0a683ee06e76ea401894129198a477a2d95b1920fdbb4222 |
| protocol/cold-resume.cdl | 1 | 4 | sha256:cbf8a049d18abe287b4163f2a47e856912463f132dfa9b1e899609fa666faa0c |
| protocol/invocation-lifecycle.cdl | 3 | 8 | sha256:3e04a6cb7e01e352e10c58411b427893abf41dc810684228acfb18a6062227ed |
| protocol/ticket-fsm.cdl | 6 | 20 | sha256:33e51665b071fedafb359fd074b95b64ce76296f26502b01fae3cd281f035125 |
| protocol/coverage-and-exits.cdl | 6 | 19 | sha256:ab4c7c59b85afa3f085b495dd516733a3724b45771e3a7e3c983e48b361cd402 |
| protocol/synthesis-catalogs.cdl | 13 | 66 | sha256:57f7d469e4dd1eabe8cd0122a55cfe042cba64231f0acb6c82a2796953fb666d |
| protocol/traceability.cdl | 5 | 15 | sha256:2b888630ef79d80f6a6dacdcddae453f56572997162c6c2970643794ed713e80 |
| protocol/handbook-and-decisions.cdl | 10 | 23 | sha256:56c6ad208fbc054697e2cd51b3b8c15e3b4872407ac456328c5c2b175273c066 |
| protocol/packaging.cdl | 9 | 50 | sha256:36493cd62b0d502d2b2c5455ae648b6e6e9d461d31c4d73a6e7558b9c1e15fd0 |
| protocol/conformance.cdl | 5 | 10 | sha256:9a8fb1d1d5212eb46e646fc31b84b0b2869ff1581c85f0c1ec5f8c598656407a |

## Declared deferrals

- standalone protocol-quality test on the generated protocol.md edition is pending (content completeness is enforced by TestGeneratedProtocolCoversOracle)
- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

The CDL sources under `protocol/` are the normative authority; `protocol.md` is a generated render (fingerprint above). `protocol/legacy/protocol-4.1.3.md` is the frozen bootstrap parity oracle.
