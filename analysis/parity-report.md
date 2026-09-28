# Crown Parity Report

- Schema version: 1
- Generator: parityreport parityreport/0.1
- Generated at: 2026-09-25T00:00:00Z
- Protocol: 4.1.2 (sha256:518c9e4922200ae6e8c01a1568aa448c9c28dcde3bf50559d0a82b4cd5fa267c)
- Language: cdl/0.2 (EIR format 2)

## Section coverage

| Section | Heading | CDL source | Status |
| :-- | :-- | :-- | :-- |
| 0 | Purpose, Guarantees, and Conceptual Flow |  | wrapper |
| 1.1 | Plane 1 — Forensic Plane | examples/spec/foundations.cdl | implemented |
| 1.2 | Plane 2 — Assurance Plane | examples/spec/foundations.cdl | implemented |
| 1.3 | Plane 3 — Synthesis Plane | examples/spec/foundations.cdl | implemented |
| 1.4 | Plane 4 — Human Reconstruction Handbook | examples/spec/foundations.cdl | implemented |
| 1.5 | Non-authoritative sidecars | examples/spec/foundations.cdl | implemented |
| 2.1 | Persona definition and purity | examples/spec/personas-povs.cdl | implemented |
| 2.2 | Canonical entry ownership | examples/spec/personas-povs.cdl | implemented |
| 2.3 | Five isolated POVs | examples/spec/personas-povs.cdl | implemented |
| 2.4 | Traversal tracks | examples/spec/personas-povs.cdl | implemented |
| 2.5 | Capability and runtime state | examples/spec/personas-povs.cdl | implemented |
| 2.6 | Capability grouping and target decision | examples/spec/personas-povs.cdl | implemented |
| 3.1 | `0A-PREFLIGHT.md` | examples/spec/preflight-registry.cdl | implemented |
| 3.2 | Foundational registries | examples/spec/workspace-registries.cdl | implemented |
| 4.1 | ID generation | examples/spec/typed-id.cdl | implemented |
| 4.1.1 | Semantic payload identity and certification envelopes | examples/spec/semantic-payload-identity.cdl | implemented |
| 4.1.2 | Normative canonical hash profile and packaging artifacts | examples/spec/canonical-hash-profile.cdl | implemented |
| 4.2 | Atomic unit rule | examples/spec/inventory-and-frontier.cdl | implemented |
| 4.3 | Source inventory denominator — `10-SOURCE-INVENTORY.md` | examples/spec/inventory-and-frontier.cdl | implemented |
| 4.4 | Traversal frontier — `11-TRAVERSAL-FRONTIER.md` | examples/spec/inventory-and-frontier.cdl | implemented |
| 4.5 | Scope-specific source coverage — `16-SOURCE-COVERAGE.md` | examples/spec/inventory-and-frontier.cdl | implemented |
| 4.6 | Stable supersession and tombstones | examples/spec/inventory-and-frontier.cdl | implemented |
| 5.1 | Dedicated prefix groups | examples/spec/status-taxonomy.cdl | implemented |
| 5.2 | Claim-level evidence — `12-CLAIM-EVIDENCE.md` | examples/spec/evidence-and-decisions.cdl | implemented |
| 5.3 | Sanitized projections and helper limitations | examples/spec/evidence-and-decisions.cdl | implemented |
| 5.4 | Contradictions — `14-CONTRADICTIONS.md` | examples/spec/evidence-and-decisions.cdl | implemented |
| 5.5 | Authoritative decisions — `15-DECISIONS.md` | examples/spec/evidence-and-decisions.cdl | implemented |
| 5.6 | Acquisition-candidate reconciliation — `17-ACQUISITION-CANDIDATES.md` | examples/spec/evidence-and-decisions.cdl | implemented |
| 5.7 | Human confirmations — `18-CONFIRMATIONS.md` | examples/spec/evidence-and-decisions.cdl | implemented |
| 6.1 | Atomic component schema | examples/spec/extraction-evolution.cdl | implemented |
| 6.2 | Dual-entry discovery | examples/spec/extraction-evolution.cdl | implemented |
| 6.3 | Staging and promotion | examples/spec/extraction-evolution.cdl | implemented |
| 6.4 | Depromotion | examples/spec/extraction-evolution.cdl | implemented |
| 6.5 | Surgical evolution | examples/spec/extraction-evolution.cdl | implemented |
| 7.1 | Trust boundary | examples/spec/acquisition-trust.cdl | implemented |
| 7.2 | Acquisition register | examples/spec/acquisition-trust.cdl | implemented |
| 7.3 | Logical explosion and virtual coordinates | examples/spec/acquisition-trust.cdl | implemented |
| 7.4 | Required explosion coverage | examples/spec/acquisition-trust.cdl | implemented |
| 7.5 | n8n human acquisition loop | examples/spec/export-acquisition-loop.cdl | implemented |
| 7.6 | Normalized maps | examples/spec/normalized-maps.cdl | implemented |
| 7.7 | Export reconciliation — `13-EXPORT-RECONCILIATION.md` | examples/spec/export-reconciliation.cdl | implemented |
| 8.1 | Resume identity header | examples/spec/invocation-context.cdl | implemented |
| 8.2 | Strict single-scope modes | examples/spec/invocation-context.cdl | implemented |
| 8.3 | Invocation-mode enum and explicit multi-file modes | examples/spec/invocation-context.cdl | implemented |
| 8.4 | Mandatory read sets | examples/spec/invocation-context.cdl | implemented |
| 8.5 | Cold resume check | examples/spec/cold-resume.cdl | implemented |
| 8.6 | Concurrent persona traversal | examples/spec/invocation-context.cdl | implemented |
| 8.7 | Stale checkpoint guard | examples/spec/invocation-context.cdl | implemented |
| 8.8 | `0G` invocation log | examples/spec/invocation-context.cdl | implemented |
| 9.1 | Ticket schema and canonical IDs | examples/spec/ticket-fsm.cdl | implemented |
| 9.2 | FSM and write rights | examples/spec/ticket-fsm.cdl | implemented |
| 9.3 | Human Hatch action prerequisites | examples/spec/ticket-fsm.cdl | implemented |
| 9.4 | Probe specification and asynchronous handoff | examples/spec/ticket-fsm.cdl | implemented |
| 9.5 | No-mock integrity | examples/spec/ticket-fsm.cdl | implemented |
| 9.6 | No-mock fallback payload | examples/spec/ticket-fsm.cdl | implemented |
| 10.1 | Evidence-backed `[R-SWEPT]` | examples/spec/coverage-and-exits.cdl | implemented |
| 10.2 | Per-kind coverage arithmetic | examples/spec/coverage-and-exits.cdl | implemented |
| 10.3 | Composite Exit A — Deconstruction Closure | examples/spec/coverage-and-exits.cdl | implemented |
| 10.4 | Exit B — Stagnation | examples/spec/coverage-and-exits.cdl | implemented |
| 10.5 | Exit C — Oscillation | examples/spec/coverage-and-exits.cdl | implemented |
| 10.6 | Exit D — Limit exhaustion and deterministic iteration accounting | examples/spec/coverage-and-exits.cdl | implemented |
| 11.1 | Partial versus Final Synthesis | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.2 | Common synthesis block header | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.3 | Architecture blueprint — `90-ARCH-BLUEPRINT.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.4 | Entity — `91-DATA-MODEL.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.5 | Relationship — `91-DATA-MODEL.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.6 | State machine — `91-DATA-MODEL.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.7 | Database routine — `91-DATA-MODEL.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.8 | Business rule — `92-BUSINESS-RULES.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.9 | Use case — `93-USE-CASES.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.10 | Interface — `94-INTERFACES.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.11 | Deployment, configuration, and scheduling — `95-DEPLOYMENT.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.12 | NFR and security — `96-NON-FUNCTIONAL-SECURITY.md` | examples/spec/synthesis-catalogs.cdl | implemented |
| 11.13 | Persona profile | examples/spec/synthesis-catalogs.cdl | implemented |
| 12.1 | Traceability — `20-TRACEABILITY.md` | examples/spec/traceability.cdl | implemented |
| 12.2 | Required reciprocal references | examples/spec/traceability.cdl | implemented |
| 12.3 | Synthesis-ID timing | examples/spec/traceability.cdl | implemented |
| 12.4 | Dependency and impact provenance | examples/spec/traceability.cdl | implemented |
| 12.5 | Deterministic validation summary | examples/spec/traceability.cdl | implemented |
| 13.1 | Audience paths and progressive disclosure | examples/spec/handbook-and-decisions.cdl | implemented |
| 13.2 | Required chapter behavior | examples/spec/handbook-and-decisions.cdl | implemented |
| 13.3 | Disabled and Dormant Features chapter | examples/spec/handbook-and-decisions.cdl | implemented |
| 13.4 | Handbook quality gate | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.1 | Decision log | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.2 | Modernization mapping templates | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.3 | Business-rule ownership | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.4 | Interface compatibility | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.5 | Equivalence and acceptance suite | examples/spec/handbook-and-decisions.cdl | implemented |
| 14.6 | Sole-input invariant | examples/spec/handbook-and-decisions.cdl | implemented |
| 15.1 | Strict conditions | examples/spec/packaging.cdl | implemented |
| 15.1.1 | Universal deterministic packaging rules | examples/spec/packaging.cdl | implemented |
| 15.1.2 | Exit E Candidate Report ordered payload schema | examples/spec/packaging.cdl | implemented |
| 15.1.3 | Candidate Payload Manifest ordered payload schema | examples/spec/packaging.cdl | implemented |
| 15.1.4 | Exit E Content-Readiness Report ordered payload schema | examples/spec/packaging.cdl | implemented |
| 15.2 | Scope certificate | examples/spec/packaging.cdl | implemented |
| 15.3 | Package contents | examples/spec/packaging.cdl | implemented |
| 15.4 | Acyclic certification sequence and hash domains | examples/spec/packaging.cdl | implemented |
| 15.5 | Reproducible final-bundle verification | examples/spec/packaging.cdl | implemented |
| 17.1 | Required order | examples/spec/migration-and-conformance.cdl | implemented |
| 17.2 | Migration guarantees | examples/spec/migration-and-conformance.cdl | implemented |
| 17.3 | v3 final-artifact compatibility | examples/spec/migration-and-conformance.cdl | implemented |
| 19.1 | Conforming implementation capabilities | examples/spec/migration-and-conformance.cdl | implemented |
| 19.2 | Normative conformance cases | examples/spec/migration-and-conformance.cdl | implemented |
| 19.3 | Corpus acceptance | examples/spec/migration-and-conformance.cdl | implemented |

## Checks

| Check | Status | Detail |
| :-- | :-- | :-- |
| assembly-compiles | pass | 26 sources |
| drift-tests | pass | 23 tests in 23 files |
| golden-eir | pass | sha256:1398d922ff88e2a725c86efe4728d70fc1c90115ae8d6562ba7e75e896010d37 |
| golden-prompt | pass | sha256:4e669e5bc2b60de6cf0afbd8bb296d47cafdca61f8b17e5bb74b9d161badd076 |
| fixture-oracle | pass | 24/24 passed |
| identity-ledger | pass | 542 identities across 26 sources |

## Fixtures

24/24 passed, 0 failed.

## Sources

| Source | Sections | Identities | Fingerprint |
| :-- | :-- | --: | :-- |
| examples/spec/globals.cdl | 0 | 0 | sha256:f7f351b217c8c2ae284ebb06ea7a5f49b5451c0643c0aaf453680cc24fb7d53b |
| examples/spec/foundations.cdl | 9 | 30 | sha256:5bfc57211c460060f03d736e65bd267c01e621f2b7b8d602cf0c0f3336b505a0 |
| examples/spec/personas-povs.cdl | 6 | 18 | sha256:9446d63a4b9d9d326d1b2e2cd0e583d17cecb600153ede61c2c9f886efe16259 |
| examples/spec/preflight-registry.cdl | 1 | 9 | sha256:4bc76cea0713269561a43f884b89432b5cb5b2998c040847d0e4d11e9b0ef9e4 |
| examples/spec/typed-id.cdl | 1 | 9 | sha256:8191941ec70330771c9d6e0b88809e67b2e44258c112c29d905381e83eb57912 |
| examples/spec/semantic-payload-identity.cdl | 1 | 7 | sha256:93fd4a060373fbdbd88a5bc4973f3d73d661bd902bd10a00e87c8cd85191f362 |
| examples/spec/canonical-hash-profile.cdl | 1 | 13 | sha256:551851233d40b93595606db5c01925ef5e808a8dcb33749cb85e221cce6d4213 |
| examples/spec/inventory-and-frontier.cdl | 5 | 23 | sha256:70ea5219b5db253f332f03c20301af1118252c6bfff53f2cbc68dece16ffc898 |
| examples/spec/status-taxonomy.cdl | 1 | 112 | sha256:726a44b5d7ded376a2b76e7f81840db74e013e68bedef20e30d74863769f02b8 |
| examples/spec/evidence-and-decisions.cdl | 6 | 13 | sha256:a1005c84673292408a24b95d99664f68563c3538645779f36e714551fd25f93c |
| examples/spec/extraction-evolution.cdl | 5 | 15 | sha256:ca15d2c508ca414d19f1fceaeddcf3a7954f14a8ebd03e9b68d27d4859564750 |
| examples/spec/workspace-registries.cdl | 1 | 9 | sha256:72aaa4f50e79ca9df15ff27692c13fe1f6ab83b711e544b91ab9104396024ec6 |
| examples/spec/acquisition-trust.cdl | 4 | 10 | sha256:abde8c00516c4468b670d6780ed2425f2aef430f7860d6f614ddce78fc152fa5 |
| examples/spec/export-acquisition-loop.cdl | 1 | 10 | sha256:973802a4a41edc9068823623e3b22d6b7eb8bad75082f61ea6b5c2a62cc1485f |
| examples/spec/normalized-maps.cdl | 1 | 4 | sha256:beb7585061fe6e045e377ab0ab7ee961b935b437c87a7e9c128cfdeee257267d |
| examples/spec/export-reconciliation.cdl | 1 | 22 | sha256:2c424497e45d2e3dd204ff9992e5f36a3f4391154080588e527b87c5ed8000c2 |
| examples/spec/invocation-modes.cdl | 0 | 0 | sha256:c599bd58c22d06daf961ddd6952a1d5e78598178f762efc7c2242f792ed75b57 |
| examples/spec/cold-resume.cdl | 1 | 4 | sha256:9ed5048be3887d41d5ef49a7a93e49dfcc39bf30870e4499ca97190a4758cf7b |
| examples/spec/invocation-context.cdl | 7 | 25 | sha256:670a49c66ac1baee7398c369cb0262a03f84a4655f783f5020cf86b109b2df62 |
| examples/spec/ticket-fsm.cdl | 6 | 20 | sha256:6c6b7b4f631a8a2b6e5b30f8f955530941f8b36608285d063cdc19a5860fdd47 |
| examples/spec/coverage-and-exits.cdl | 6 | 19 | sha256:10d78954cb51a04517f949477857375b8811996b0fa77cabda3cc3e1265efeb1 |
| examples/spec/synthesis-catalogs.cdl | 13 | 66 | sha256:909ba2f8e46d460acea5d509f51f349593b2b0ce1d2e06227fdc74d0776a1435 |
| examples/spec/traceability.cdl | 5 | 15 | sha256:482415d3dd89057fe0fc0bd3966ee8ae2f73186899951d08b2915153bf126b0a |
| examples/spec/handbook-and-decisions.cdl | 10 | 23 | sha256:eaae9e34a17395e90e28b7ea8b455d3b147989b437e78bd3a9206d35940b8458 |
| examples/spec/packaging.cdl | 9 | 50 | sha256:01b891065ab12be512a0866943f14f43da6b915f9dbc20db7bee9ad801d2dd60 |
| examples/spec/migration-and-conformance.cdl | 7 | 16 | sha256:0209969885a4c48d7cd1933a0666cb5c19bc8a97e0839689bc09b250b6468e24 |

## Declared deferrals

- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

Human acceptance of this report is the S4 crown gate; until acceptance, `protocol.md` remains the sole normative authority.
