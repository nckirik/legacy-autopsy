# Crown Parity Report

- Schema version: 1
- Generator: parityreport parityreport/0.1
- Generated at: 2026-09-25T00:00:00Z
- Protocol: 4.1.2 (sha256:c5c310f6e2047398d177ee5eb691ada726887ad0aa6bac676b83d76d3129c825)
- Language: cdl/0.2 (EIR format 2)

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
| 8.6 | Concurrent persona traversal | protocol/invocation-context.cdl | implemented |
| 8.7 | Stale checkpoint guard | protocol/invocation-context.cdl | implemented |
| 8.8 | `0G` invocation log | protocol/invocation-context.cdl | implemented |
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
| 19.1 | Conforming implementation capabilities | protocol/conformance.cdl | implemented |
| 19.2 | Normative conformance cases | protocol/conformance.cdl | implemented |
| 19.3 | Corpus acceptance | protocol/conformance.cdl | implemented |

## Checks

| Check | Status | Detail |
| :-- | :-- | :-- |
| assembly-compiles | pass | 26 sources |
| drift-tests | pass | 23 tests in 23 files |
| generated-protocol | pass | sha256:2478eecfab1b884bf69c1519482841d6edf2793c6f0c63ab4d26669d49351086 |
| golden-eir | pass | sha256:e89f43e346edadf332e4959c6ffaa361f5c88237294ada6a4a4c4cbba2e761c6 |
| golden-prompt | pass | sha256:b9577733729a5695a131d68ebfc9e75e567ade58eb0836b5afebfdfdc25bc42f |
| fixture-oracle | pass | 24/24 passed |
| identity-ledger | pass | 536 identities across 26 sources |

## Fixtures

24/24 passed, 0 failed.

## Sources

| Source | Sections | Identities | Fingerprint |
| :-- | :-- | --: | :-- |
| protocol/globals.cdl | 0 | 0 | sha256:40b869843eb972f22ce7b37502efa2ffbcf0c9ba284fc9a81d04221d8734fcf0 |
| protocol/foundations.cdl | 9 | 30 | sha256:f50326f1874e7dcb46c7b22ec18ec4713a9081f11b65e8c426979aa7011d561d |
| protocol/personas-povs.cdl | 6 | 18 | sha256:8472228d5485d5f8de4b24b7ece20fa97a47f43fecb1dc570265ca2e732ce4d5 |
| protocol/preflight-registry.cdl | 1 | 9 | sha256:6418c22b018c485945e55804664176f02cae1fc0d607610dc6ed57bc0e2994eb |
| protocol/typed-id.cdl | 1 | 9 | sha256:49efe885b452536afa13e32562f8bcd64abd109ae2bf35f9ff09847637c32018 |
| protocol/semantic-payload-identity.cdl | 1 | 7 | sha256:047c65826b9b1aa98afd9413ec51402ea592bfdcf9b655787faab06b38e04631 |
| protocol/canonical-hash-profile.cdl | 1 | 13 | sha256:af4fd0181fcefe46b6689e2ca48c897779d03263a682840b5b3c73ec8bfdbfb7 |
| protocol/inventory-and-frontier.cdl | 5 | 23 | sha256:eacfbd4cad5cde149c270d3b67809a9b39f917c208443d197394ede1e9d375b9 |
| protocol/status-taxonomy.cdl | 1 | 112 | sha256:c4a695e8d8d76ce9e782a909d49c8e9ad62a72d31f57b28b71e9c19078b47d2e |
| protocol/evidence-and-decisions.cdl | 6 | 13 | sha256:822b003067814ded8302e989ca78e17cf15959e6e887a1e9230fd9dbc992c27c |
| protocol/extraction-evolution.cdl | 5 | 15 | sha256:32938e7b7f862fbc81ef144b99efe3edfaa622386e5e1e01451517aa1312b309 |
| protocol/workspace-registries.cdl | 1 | 9 | sha256:248cc6fce1b2df54fc8cf52823259d15ecdf6b3b3b217633c9e544efde8331f7 |
| protocol/acquisition-trust.cdl | 4 | 10 | sha256:a43b2fd0489321e5108c58b8a2e78a5713e7146c93808b821fede517e8caa809 |
| protocol/export-acquisition-loop.cdl | 1 | 10 | sha256:165ca179fe3a9a7bfd1a97d67491cafad44bdde76a2f7112a12376f25801d2a8 |
| protocol/normalized-maps.cdl | 1 | 4 | sha256:4afabe97b491546b5d1dc78ffeb9e35f43ff4b2378873a0202ddc8debc8bcff9 |
| protocol/export-reconciliation.cdl | 1 | 22 | sha256:12bed9f7c30f871446665180b1a35e647b3ec88e2c98500a8a28c45f97caea4e |
| protocol/invocation-modes.cdl | 0 | 0 | sha256:c599bd58c22d06daf961ddd6952a1d5e78598178f762efc7c2242f792ed75b57 |
| protocol/cold-resume.cdl | 1 | 4 | sha256:63288693e96ac576e139a553576cd50c8bba68490b52afb656bc9f4a8b0cb9f1 |
| protocol/invocation-context.cdl | 7 | 25 | sha256:931368ef1ed1706499823d809ec8ec6e493ec3bf76141c71e0ee137ffafd1381 |
| protocol/ticket-fsm.cdl | 6 | 20 | sha256:c6505ccc1220da4d0c8e05aff07a16af8aa2f970358a1bba3d5e3c8c2e8bea05 |
| protocol/coverage-and-exits.cdl | 6 | 19 | sha256:04d57842984f9af28b40ce8d77da11e9114be9589b1e4c1d100bc5c4d73c7661 |
| protocol/synthesis-catalogs.cdl | 13 | 66 | sha256:2b3b17b37125578e959eeb7b84ba2566037575bffbda23591b7d566dc8db0d35 |
| protocol/traceability.cdl | 5 | 15 | sha256:3d1795b04582dfd07dfbf02e6d0c9af30dda3387ebbcdfa6721a843df88ea38b |
| protocol/handbook-and-decisions.cdl | 10 | 23 | sha256:812b1ccf5fd199e62fe48b46465c1563f7d29a06e349b192dee0aa2b9bb73c39 |
| protocol/packaging.cdl | 9 | 50 | sha256:628e9034db97140c78ccf6b569b15689efd012762852e65eb1b769bc8a32062c |
| protocol/conformance.cdl | 5 | 10 | sha256:6aaee306adc0d3c525c8e7b68b9484ae3611902f5a41b70d5bd8f6f97b389f98 |

## Declared deferrals

- standalone protocol-quality test on the generated protocol.md edition is pending (content completeness is enforced by TestGeneratedProtocolCoversOracle)
- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

The CDL sources under `protocol/` are the normative authority; `protocol.md` is a generated render (fingerprint above). `protocol/legacy/protocol-4.1.2.md` is the frozen bootstrap parity oracle.
