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
| 17.1 | Required order | protocol/migration-and-conformance.cdl | implemented |
| 17.2 | Migration guarantees | protocol/migration-and-conformance.cdl | implemented |
| 17.3 | v3 final-artifact compatibility | protocol/migration-and-conformance.cdl | implemented |
| 19.1 | Conforming implementation capabilities | protocol/migration-and-conformance.cdl | implemented |
| 19.2 | Normative conformance cases | protocol/migration-and-conformance.cdl | implemented |
| 19.3 | Corpus acceptance | protocol/migration-and-conformance.cdl | implemented |

## Checks

| Check | Status | Detail |
| :-- | :-- | :-- |
| assembly-compiles | pass | 26 sources |
| drift-tests | pass | 23 tests in 23 files |
| generated-protocol | pass | sha256:eb301ede9212a96019ed7a954f67e7b2db6c2c5797e8612a61f3eaa1fc236e72 |
| golden-eir | pass | sha256:ed9c391cef76da15af913d6a5c0c3032f2b11ae138ca6e080d041c63c5c64855 |
| golden-prompt | pass | sha256:119c21b56526c8f646615769ebfd762c7b1171f80c4910cc47195e21155dec05 |
| fixture-oracle | pass | 24/24 passed |
| identity-ledger | pass | 542 identities across 26 sources |

## Fixtures

24/24 passed, 0 failed.

## Sources

| Source | Sections | Identities | Fingerprint |
| :-- | :-- | --: | :-- |
| protocol/globals.cdl | 0 | 0 | sha256:f7f351b217c8c2ae284ebb06ea7a5f49b5451c0643c0aaf453680cc24fb7d53b |
| protocol/foundations.cdl | 9 | 30 | sha256:5add40f4812f174b2c96dd5a86c531188961fefbe734af6aaa353f1faabd7af1 |
| protocol/personas-povs.cdl | 6 | 18 | sha256:ab4662dedd2b1ff543b7dfbd941278e4d87881e7cded529a4cbe57129d518a5c |
| protocol/preflight-registry.cdl | 1 | 9 | sha256:ac4141a8b9be0c11b9c3f1d5d0136edfa6e61c2923468cb621d4942393979656 |
| protocol/typed-id.cdl | 1 | 9 | sha256:8191941ec70330771c9d6e0b88809e67b2e44258c112c29d905381e83eb57912 |
| protocol/semantic-payload-identity.cdl | 1 | 7 | sha256:93fd4a060373fbdbd88a5bc4973f3d73d661bd902bd10a00e87c8cd85191f362 |
| protocol/canonical-hash-profile.cdl | 1 | 13 | sha256:551851233d40b93595606db5c01925ef5e808a8dcb33749cb85e221cce6d4213 |
| protocol/inventory-and-frontier.cdl | 5 | 23 | sha256:fb192c748dc2aafc853f130bb68f808e761a024bf6effa2d324170797941565d |
| protocol/status-taxonomy.cdl | 1 | 112 | sha256:e73565b556dad4fc6e42701f252cfd26a0cabb5cc8cc80f6837b69fcbf4efb60 |
| protocol/evidence-and-decisions.cdl | 6 | 13 | sha256:88c67d05353ed105f72bbc0ab59dd82d76b95ebd740328d9d6f6f8fae94fbc00 |
| protocol/extraction-evolution.cdl | 5 | 15 | sha256:7b58a6e46c19555620b1c499ecec05cf9a1def94b15205e05dd35c01bfc8c11b |
| protocol/workspace-registries.cdl | 1 | 9 | sha256:72aaa4f50e79ca9df15ff27692c13fe1f6ab83b711e544b91ab9104396024ec6 |
| protocol/acquisition-trust.cdl | 4 | 10 | sha256:abde8c00516c4468b670d6780ed2425f2aef430f7860d6f614ddce78fc152fa5 |
| protocol/export-acquisition-loop.cdl | 1 | 10 | sha256:973802a4a41edc9068823623e3b22d6b7eb8bad75082f61ea6b5c2a62cc1485f |
| protocol/normalized-maps.cdl | 1 | 4 | sha256:beb7585061fe6e045e377ab0ab7ee961b935b437c87a7e9c128cfdeee257267d |
| protocol/export-reconciliation.cdl | 1 | 22 | sha256:7d9c5101fcfec9bfd18cd764ec981970e16e93d47023b69a9ec21c49813c0d6b |
| protocol/invocation-modes.cdl | 0 | 0 | sha256:c599bd58c22d06daf961ddd6952a1d5e78598178f762efc7c2242f792ed75b57 |
| protocol/cold-resume.cdl | 1 | 4 | sha256:9ed5048be3887d41d5ef49a7a93e49dfcc39bf30870e4499ca97190a4758cf7b |
| protocol/invocation-context.cdl | 7 | 25 | sha256:a55d9af5fa050c2a7150e226479a1a2b0884f12a7ba657df0aad66dce213c25c |
| protocol/ticket-fsm.cdl | 6 | 20 | sha256:2e10a42108f55fc04eaa118b8acaa34b1ec982c5a9cbc15dad45c5739206025c |
| protocol/coverage-and-exits.cdl | 6 | 19 | sha256:cdca1eeff0c25fadccc2fe233cf8a1bd55d1adfe6d7463af6026f0ebb245562b |
| protocol/synthesis-catalogs.cdl | 13 | 66 | sha256:51131344e48774509018fea153ce505683cf8dcd7686d2e7112d4fee1c409636 |
| protocol/traceability.cdl | 5 | 15 | sha256:a4cbc6d9175f181d0a84adeae8c9337ef4a7cb35356c993ad0452f7ebf9944b3 |
| protocol/handbook-and-decisions.cdl | 10 | 23 | sha256:6a498e9797958134e594923db8bd609834bb7e2993c99a68ac6c8ef053747e4f |
| protocol/packaging.cdl | 9 | 50 | sha256:8a1ae5f10cc4fbb6dbfb69677aef545a2519ad9d4d684bb472d2006f8ecd42df |
| protocol/migration-and-conformance.cdl | 7 | 16 | sha256:6a042953eb424077ca834dce18bf21f2fa89511360b6779679af4232a557820c |

## Declared deferrals

- generated protocol.md is a structural render with condensed rule text; the frozen oracle remains the usable standalone prompt until verbatim enrichment and the standalone protocol-quality test pass
- canonical.markdown engine remains a bounded v1 implementation (internal/capabilities/canonical.go)
- skill projections stay hand-maintained through S4; EIR reference-view generation is deferred
- CDL editor tooling (highlighting, canonical printer) is planned after crown as S5
- renderer emits EIR reference and prompt editions only; no HTML/backend editions yet

The CDL sources under `protocol/` are the normative authority; `protocol.md` is a generated render (fingerprint above). `protocol/legacy/protocol-4.1.2.md` is the frozen bootstrap parity oracle.
