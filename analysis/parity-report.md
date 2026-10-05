# Crown Parity Report

- Schema version: 1
- Generator: parityreport parityreport/0.1
- Generated at: 2026-09-25T00:00:00Z
- Protocol: 4.2.1 (sha256:f2005f08760334f874f8cc26382c4933ffea181ff34eef4560224cba340aef7d)
- Language: cdl/0.3 (EIR format 3)

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
| generated-protocol | pass | sha256:f2005f08760334f874f8cc26382c4933ffea181ff34eef4560224cba340aef7d |
| golden-eir | pass | sha256:e766ba55b66a9fc54adf27dbecff23ca4c4b51fe65bfb42eff75a14869d8b4c3 |
| golden-prompt | pass | sha256:7b3a7bce1b945ab0d5c171fd1eb18b7966bf3782033b55a35daab4966dbf7c8a |
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
| protocol/personas-povs.cdl | 6 | 18 | sha256:0589cce227972317bd7f63f271001125d4743fe7997c550ce596e7878f64c360 |
| protocol/preflight-registry.cdl | 1 | 9 | sha256:3616d5a269131c8990557c9ae6cc07d2aef3a49e7964adf2216f291ddf4ec466 |
| protocol/workspace-registries.cdl | 1 | 9 | sha256:a63a48f23ddb869889298620bb1ae892311f3bece06dea9062fe8d9583557a97 |
| protocol/typed-id.cdl | 1 | 9 | sha256:cc2efdcefb008626001bba599c2202b6606fa64db5e1c85ecb4b8706cd93fb65 |
| protocol/semantic-payload-identity.cdl | 1 | 7 | sha256:76fabcee6b8afcf3c37907f03c65aad9929badd62c8a3a27420c41a64f6f99d9 |
| protocol/canonical-hash-profile.cdl | 1 | 13 | sha256:b9f59a77a8c0c0e5e5c649c2e8792193746fdb079cdaadc599ff0ab119e8686d |
| protocol/inventory-and-frontier.cdl | 5 | 23 | sha256:174b2b15ca2df2604c311c78d85a70ed69055f2e8a77c2e8b2bb512fb4813754 |
| protocol/status-taxonomy.cdl | 1 | 112 | sha256:f57baa323d813e24e2886f322967736e28e42f8e2817294b0cc50f7b7fd1cce3 |
| protocol/evidence-and-decisions.cdl | 6 | 13 | sha256:09a2ff0b2c8435d42ed78e7bf18a78f91b554f5266f4a26801fc75834012c6e5 |
| protocol/extraction-evolution.cdl | 5 | 15 | sha256:355235731afd4f07809f5133d2cbd3c1c7c32aaeaf82917f341957d04d0a21f2 |
| protocol/acquisition-trust.cdl | 4 | 10 | sha256:0f68499f4ba1a9ddee5884735bd5718c2b3fb9d7c878c514899e66943522c610 |
| protocol/export-acquisition-loop.cdl | 1 | 10 | sha256:3b4657ae694ade4a3325651b6256443f5c3eb233a50af73b31308dc2fd51baf7 |
| protocol/normalized-maps.cdl | 1 | 4 | sha256:2ee0fd332ff549c8df63ecdf788942ac9ab1941f1b8b2721d46876a6c5532528 |
| protocol/export-reconciliation.cdl | 1 | 22 | sha256:0bcf9505c4284408be0a48561179b6f5002ca47395d415551dd34f3874e3f692 |
| protocol/invocation-context.cdl | 4 | 17 | sha256:f5f821ffe7446d76fe677541cce40381541e0d1246c83e08822ddf2d34edac81 |
| protocol/cold-resume.cdl | 1 | 4 | sha256:cbf8a049d18abe287b4163f2a47e856912463f132dfa9b1e899609fa666faa0c |
| protocol/invocation-lifecycle.cdl | 3 | 8 | sha256:25dd23cf1d8aabc960e67c34a13b518af2ceb73bfc0bc8b156cf7739de64e63e |
| protocol/ticket-fsm.cdl | 6 | 20 | sha256:ee9df2e138cb46840e8680243d58f64620acdafa96d7ef69ed7e4adfa5be025e |
| protocol/coverage-and-exits.cdl | 6 | 19 | sha256:72cb2b377c1fc6f743c0371db785876b0ad750bbe764a43c79bfcbe470a3b2aa |
| protocol/synthesis-catalogs.cdl | 13 | 66 | sha256:060d24b070b059d395b40bd93154bef296dfdc097f07bfd34a13a5073a0058c4 |
| protocol/traceability.cdl | 5 | 15 | sha256:e6732880f8e333282c6b427f1e7fcb45f2c7769838c7d839a00bf0904ec968e6 |
| protocol/handbook-and-decisions.cdl | 10 | 23 | sha256:56c6ad208fbc054697e2cd51b3b8c15e3b4872407ac456328c5c2b175273c066 |
| protocol/packaging.cdl | 9 | 50 | sha256:465c0a368f3f16e4930d860a5dc23dfef3b161edd3c405d09750ccf0f1e937e1 |
| protocol/conformance.cdl | 5 | 10 | sha256:f23f75ad0bc6e191d58280e5dc3bec54bc5e51b1b35af0fda981f106cefd5697 |

## Declared deferrals

- standalone protocol-quality test on the generated protocol.md edition is pending (content completeness is enforced by TestGeneratedProtocolCoversOracle)
- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

The CDL sources under `protocol/` are the normative authority; `protocol.md` is a generated render (fingerprint above). `protocol/legacy/protocol-4.1.3.md` is the frozen bootstrap parity oracle.
