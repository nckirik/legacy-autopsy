# Development

> Non-authoritative developer guide. The CDL sources under [`protocol/`](../protocol/) define behavior; this file records the current build and command workflow.

## Prerequisites

Go 1.24 or newer is required. The built CLI has no external runtime dependencies. Tests and fixtures are local and require no network, production credentials, or raw exports.

Development tooling is pinned under `scripts/`. Install dprint with the official
installer at the CI version; `scripts/format.sh` finds it on `PATH` or at
`~/.dprint/bin/dprint`.

```sh
curl -fsSL https://dprint.dev/install.sh | sh -s 0.57.4   # pinned; ~/.dprint/bin

scripts/format.sh           # gofmt -w + dprint fmt
scripts/format.sh --check   # verify formatting only
scripts/check.sh            # full local gate; mirrors CI
scripts/goldens.sh          # regenerate CDL/analysis golden assets
```

The official installer pins by version but does not verify checksums; CI's
`dprint/check` action verifies the release build attestation.

## CLI

The commands below are the implemented M0 bootstrap interface. The target CLI becomes a headless observation/administration client to the planned local service, but it does not manually start or claim semantic work during the skill-first stage. Proposed service, adapter, question, and graph commands are intentionally absent until implemented. See the [runtime design](runtime.md) and [roadmap](roadmap.md).

```sh
go run ./cmd/legacy-autopsy doctor
go run ./cmd/legacy-autopsy protocol check
go run ./cmd/legacy-autopsy init --workspace /tmp/autopsy/.extracted
go run ./cmd/legacy-autopsy workspace check --workspace /tmp/autopsy/.extracted
go run ./cmd/legacy-autopsy id --type CMP --namespace example --owner src/example.go#Run --discriminator atomic-component
go run ./cmd/legacy-autopsy validate routing
go run ./cmd/legacy-autopsy validate fixtures
```

The M0 fixture command always runs every registered bootstrap case and lists all complete Part 17.2 families as unsupported. The accepted `--implemented` flag is currently a redundant compatibility selector; public examples omit it.

The additive Spec-track pilot verbs (`spec compile|render|run`) are documented under [Spec track](#spec-track). They cover only §7.7 and do not change the M0 boundary above.

### Context packets

`context` is an advanced M0 command. It requires common §8.1 identity inputs, explicit allowed/forbidden write targets, and the exact mode. Repeat `--read-target` for mode-specific or transitive inputs. Strict single-scope modes additionally require `--persona`, `--cluster`, `--track`, `--pov`, `--pov-file`, and positive `--max-depth`; Human Hatch requires `--human-hatch-action`.

```sh
go run ./cmd/legacy-autopsy context \
  --mode Preflight \
  --system-namespace example \
  --iteration ALFA \
  --invocation-id INV-000000000001 \
  --scope preflight \
  --environment-snapshot dev/snapshot-1 \
  --write-target 0A-PREFLIGHT.md \
  --write-target 0G-DECONSTRUCTION-STATE.md \
  --forbidden-write-target "all other workspace paths" \
  --workspace /tmp/autopsy/.extracted
```

Reads are normalized and confined to the workspace, including symlink resolution; missing mandatory inputs fail closed. Strict persona-bound modes require `--persona`, `--persona-prefix`, and the canonical `--persona-directory <persona-prefix>-<persona-slug>` binding; the directory prefix must agree with `--persona-prefix`. The emitted packet embeds exact workspace bytes and fingerprints plus routed protocol text. It is **not sanitized**: redirect it only to an operator-approved location outside the repository, review it before transmission, and never commit it or send it to an unapproved provider. Automatic transitive closure, stale comparison, mutation authorization, and semantic execution remain deferred.

## Spec track

The Spec track completed the section-by-section migration and crown: the CDL sources
under [`protocol/`](../protocol/) are normative, `protocol.md` is their generated
edition, and [`protocol/legacy/protocol-4.1.3.md`](../protocol/legacy/protocol-4.1.3.md)
is the frozen parity oracle. The frozen source language is [`cdl.md`](cdl.md);
execution contracts are [`eir.md`](eir.md),
[`execution-semantics.md`](execution-semantics.md), and
[`runtime-architecture.md`](runtime-architecture.md); sequencing is in the
[roadmap](roadmap.md); migration and freeze rules are in the
[migration audit](migration-audit.md). The runtime still implements only the pilot
executable slice; that is not a conformance or Exit A/E claim.

```sh
go run ./cmd/legacy-autopsy spec compile
go run ./cmd/legacy-autopsy spec compile --golden protocol/golden/export-reconciliation.eir.json
go run ./cmd/legacy-autopsy spec render
go run ./cmd/legacy-autopsy spec render --out /tmp/export-reconciliation.light.prompt.md \
  --generated-at 2026-09-25T00:00:00Z
go run ./cmd/legacy-autopsy spec run --fixture protocol/fixtures/incomplete.json
go run ./cmd/legacy-autopsy spec run --fixture protocol/fixtures/reconciled.json --json
go run ./cmd/legacy-autopsy spec schema --out /tmp/schema.sql        # print/write DDL
go run ./cmd/legacy-autopsy spec schema --apply /tmp/autopsy.db      # apply via sqlite3
go run ./cmd/legacy-autopsy spec schema --audit                      # enum-registry completeness
go run ./cmd/legacy-autopsy spec store --db /tmp/autopsy.db          # canonical dump + fingerprint
go run ./cmd/legacy-autopsy spec store --db /tmp/autopsy.db --export /tmp/autopsy.sql
go run ./cmd/legacy-autopsy spec store --db /tmp/autopsy-2.db --import /tmp/autopsy.sql
```

`spec store --export` writes the canonical dump plus a `.sha256` content-fingerprint
sidecar; `--import` restores into a fresh database, refuses non-empty targets and
content after the canonical COMMIT block, verifies the sidecar (or an explicit
`--expect-fingerprint`) before applying anything, and re-verifies the live schema
against the EIR-derived schema after applying. A dump that exports, imports, and
re-exports is byte-identical.
go run ./cmd/legacy-autopsy spec store --db /tmp/autopsy.db # canonical dump + fingerprint

````
`spec compile` validates the identity ledger and prints the source fingerprint and EIR
hash; `spec render` renders one projection/channel with a provenance header; `spec run`
executes the native VM over a synthetic fixture and prints the canonical trace and hash.
All three default to the §7.7 source, ledger, and repository root, and accept
`--source`/`--ledger`/`--repo`.

`spec schema --audit` checks the store-facing enum registry: every declared enum
carries unique values, and every enum-typed column in the generated DDL has a
CHECK constraint listing exactly the declared values for that table context;
same-named columns on different tables never cross-check.

The generated schema (version 3) derives tables from declared storage classes:
`record`, `registry`, `singleton`, `log`, and `child` become tables; `value` and
`artifact` shapes stay embedded. Child rows carry a `la_parent_id` foreign key
and a `UNIQUE (parent, key...)` natural key over their declared `KEY` columns,
declared `REF` clauses become foreign keys (the target column must be the target
FIELD's declared `KEY`), and every write path enforces foreign keys
(`PRAGMA foreign_keys=ON`). `store.IngestChild` inserts owned rows with
`<parent>|<key>|...` identity; dumps and imports carry the schema version and
refuse incompatible versions.

The same checks are available as Go tests:

```sh
go test ./cdl/               # compiler, renderer, ledger, negative fixtures
go test ./internal/parity/   # test A (backend) and test B (protocol parity)
go test ./internal/analysis/ # classification coverage and metrics golden
````

S2 analysis lives in [`coverage-profile.md`](coverage-profile.md) and
[`capability-contracts.md`](capability-contracts.md). Regenerate the metrics golden with
`go test ./internal/analysis/ -update`; never hand-edit
[`analysis/metrics.golden.json`](../analysis/metrics.golden.json).

Generated pilot assets are golden-checked and must be regenerated deterministically,
never hand-edited:

```sh
go test ./cdl/ -run TestUpdateAssets -update
```

The M0 code and the 24-case fixture oracle remain the frozen parity baseline: bug fixes
that preserve all fixture outcomes, additive wiring, and no new protocol behavior in the
runtime. Protocol sources change through the CDL toolchain and drift tests, not through
Go edits.

## Reconstruction chain commands

The bounded semantic pipelines and packaging chain have their own verbs (all
require `--db` pointing at a workspace store):

```sh
go run ./cmd/legacy-autopsy spec handoff status                        # chain overview + next command
go run ./cmd/legacy-autopsy spec handoff report --handbook-dir preview # emit the candidate report
go run ./cmd/legacy-autopsy spec handoff manifest                      # emit the candidate payload manifest
go run ./cmd/legacy-autopsy spec handoff readiness --equivalence preview/package/equivalence/acceptance-suite.md
go run ./cmd/legacy-autopsy spec handoff equivalence                   # regenerate the acceptance suite
go run ./cmd/legacy-autopsy spec handoff certificate                   # stage-5 scope certificate
go run ./cmd/legacy-autopsy spec handoff outer --package-dir preview/bundle
go run ./cmd/legacy-autopsy spec handoff sign --package-dir preview/bundle
```

Candidate condition 15 fails when discovered schema objects or user-visible
outputs are not represented by their record family or an exact approved
exclusion (the PI-25/PI-26 coverage obligations).

Human acts are interactive by design and never machine-authored:
`spec handoff review readability|equivalence|decisions|confirm` prompt the
operator (per-chapter Pass/Fail, per-mapping accept/exclude, per-decision
approve/reject/defer, per-group confirm) and record the result with reviewer
identity and timestamp. `spec handbook mint` writes HBK
semantic payloads to a separate directory, keeping the human chapters prose-first;
`spec handbook mint --restore` rebuilds the clean views.

Decision content follows the same human-authority split:

```sh
go run ./cmd/legacy-autopsy spec handoff decisions draft --db <workspace db>
go run ./cmd/legacy-autopsy spec handoff review decisions
go run ./cmd/legacy-autopsy spec handoff decisions commit --db <workspace db>
```

`decisions draft` drafts candidate content for every unbound decision subject
(synthesis records with an empty target decision, pending capability decisions,
and contradictions whose resolution requires a decision) from exact observed
inputs. `review decisions` records the human verdicts. `decisions commit`
creates one canonical `DECISION` record per approved draft: the
decision-content fingerprint follows §5.6 (fingerprint carrier and the exact
decision-approval envelope excluded), the approval fields come from the
recorded review, and the candidate payload manifest binds the content
version/fingerprint; readiness condition 03 rejects references whose approval
binds stale content or lacks an approval envelope.

Closure validators for dead-code and lifecycle eligibility are deterministic:

```sh
go run ./cmd/legacy-autopsy spec assurance closure --db <workspace db>
```

`assurance closure` rejects `[S-DEAD-CODE]` components with missing denominator
dimensions, uncommitted source units, unresolved frontiers, active tickets, or
any discovered caller in the committed component graph, and reports
lifecycle-eligibility rows with duplicate scopes, missing provenance,
inconsistent snapshots, or `Unknown` classifications without a ticket. It also
reports schema/output coverage diagnostics: declared `SRC-TABLE`/`SRC-VIEW`
units must be claim-backed and represented by an entity record, `SRC-DR` units
by a database-routine record, and discovered output `SRC-ACTION` units must
carry a committed claim. Components declaring `[S-STAGING-VERIFIED]` or
`[S-PROMOTED]` are checked against the §6.3 promotion locks (exactly one
persona, no active ticket or open contradiction, no dependent component,
C-COVERED source units, and no material `[E-UNKNOWN]` claim). Ticket writes are validated against the declared
`TICKET-STATE` machine at commit (initial state, declared transitions, trace
history, escalation reasons), and §9.2 actor rights are enforced for
POV-bound units: only the target POV moves a ticket to pending or probe
states, only the origin POV accepts or rejects, unbound units cannot claim
either right, and transitions recorded as a human `owner-review` are exempt.

Discovery enumerates files and, separately, the non-file units they declare:

```sh
go run ./cmd/legacy-autopsy spec discovery hits --target <repo> --db <workspace db>
```

`spec discovery hits` scans the enumerated file window lexically with declared
extensions and bounds, and extracts source units of these kinds: route
registrations (`ROUTE`), SQL statements (`QUERY`), schema object declarations
(`TABLE`, `VIEW`, `DR` for triggers/procedures/functions), scheduling markers
(`JOB`), and user-visible output/export constructs such as downloads
(`ACTION`). It then runs bounded inventory units over the hits. The model
authors one `SOURCE-INVENTORY-RECORD` per hit; the runtime computes the
matching `SRC-*` identity from the copied coordinate. Hit extraction is
deterministic and re-runnable; absence of a hit kind still needs its own
evidence-backed enumerator record.

Hit and component units carry the §8.2 strict single-scope bindings: one entry
cluster (resolved from the cluster's concrete entry coordinates), its traversal
track, one POV (routes -> `POV-5`, schema objects and queries -> `POV-4`,
jobs -> `POV-3`, output/export actions -> `POV-1`), and a persona prefix with
its canonical directory when the binding resolves. A cluster whose canonical
owner is not a registered persona resolves its participants from committed
`SHARED-ENTRY` records (coordination owner first, then registered participants
in sorted order) and binds the first registered one; when no committed
shared-entry record resolves a registered participant, the context records
`Persona=None` with the active POV file under `personas/_shared/`, so ambiguity
stays fail-inclusive instead of inventing an owner.

```sh
go run ./cmd/legacy-autopsy spec discovery components --target <repo> --db <workspace db>
```

`spec discovery components` authors one `ATOMIC-COMPONENT` per hit whose
identity is anchored in the hit's committed `SRC` unit (the model fills the
§6.1 fields from a bounded file excerpt; staging status starts at
`[S-DISCOVERED]` with `[C-GAP]` coverage). It then materializes the §6.2
dual-entry `SHARED-REFERENCE` pointer rows mechanically per persona and POV,
with the discovering invocation recovered from `0G`; components without a
resolvable persona produce no pointer row and are reported instead.

```sh
go run ./cmd/legacy-autopsy spec discovery pov-manifest --db <workspace db> [--out <personas root>]
```

`spec discovery pov-manifest` deterministically projects the committed shared
references into persona-local POV manifest files
(`personas/{persona-directory}/{NN}-{NAME}-POV.md`, §6.2 pointer layout). The
projection is non-authoritative navigation; references whose component or
persona is not committed are skipped instead of guessed.

Handbook chapters carry the v4.7 content contract: the designated flow-heavy
chapters (02, 03, 05, 08) must include a valid Mermaid diagram, chapter 03 an
entity-relationship diagram of the entity/table inventory, and the
procedure-bearing chapters (03, 08) a language-agnostic pseudocode block, each
with a reason-bearing Not-Applicable escape where nothing applies. Every
Mermaid block is parsed deterministically before candidacy (declared head,
structure, balanced delimiters); unknown heads or broken blocks fail closed.
Discovered schema objects and user-visible outputs must be represented by their
record family or an exact approved exclusion.

Per-unit telemetry is appended to `<workspace db>.metrics.jsonl` (non-authoritative;
duration, retries, context bytes, result) so rehearsal comparisons stay mechanical.

## CDL formatting

`cdl fmt` is the deterministic canonical printer for `.cdl` sources. It normalizes
declaration order, indentation, alignment, and spacing, preserves opaque normative text
(`GOAL`, `TEXT`, `REASON`, and `REQUIRE """..."""`) exactly, and refuses sources with
inline or nested comments it cannot place. `protocol/*.cdl` are kept canonical;
`TestProtocolSourcesCanonical` enforces it in CI.

```sh
go run ./cmd/cdl fmt protocol/*.cdl              # print canonical form
go run ./cmd/cdl fmt -w protocol/*.cdl           # rewrite in place
go run ./cmd/cdl fmt --check protocol/*.cdl      # fail if not canonical
```

The printer is non-authoritative: it never changes semantics, and generated artifacts
are regenerated, not reformatted. dprint remains the Markdown runner; hosting `cdl fmt`
through dprint's Exec plugin is possible but not wired into CI.

VS Code support is installed locally with:

```sh
scripts/install-vscode-extension.sh [--copy] [--editor cursor] [--uninstall] [--dry-run]
```

## Markdown formatting

Non-Go Markdown is formatted with [dprint](https://dprint.dev) using the pinned plugin in
[`dprint.json`](../dprint.json). CI runs `dprint check`; formatting is deterministic and
version-pinned. dprint is a development prerequisite only; it is not part of the Go build
or runtime, and it never formats Go code.

```sh
scripts/format.sh           # format Go and Markdown
scripts/format.sh --check   # verify only
```

`protocol.md`, `protocol/golden/`, and generated analysis artifacts are excluded because
their bytes are fingerprint-bound. `docs/cdl.md` uses intentional ASCII section rules,
and `fixtures/` and `analysis/` are source or generated data - they are regenerated,
never reformatted. `protocol/*.cdl` are formatted with `cdl fmt`, not dprint. Go
formatting stays with the Go toolchain.

VS Code support is committed in `.vscode/`: the dprint extension is recommended and
Markdown format-on-save is configured.

## Local checks

`scripts/check.sh` runs the full gate; the raw commands are:

```sh
scripts/format.sh --check
go vet ./...
go test ./...
go build -o /tmp/legacy-autopsy ./cmd/legacy-autopsy
go run ./cmd/legacy-autopsy protocol check
go run ./cmd/legacy-autopsy validate routing
go run ./cmd/legacy-autopsy validate fixtures
```

For workspace smoke testing, use a disposable directory and run both `init` and `workspace check`. Unit tests verify that initialization fabricates no covered, confirmed, or passed state. When `--db` points at a workspace store, `workspace check` also validates persona ownership: every registered 0A persona has exactly one `<prefix>-<slug>` directory, every non-reserved directory is registered, and the reserved `_shared` directory is never treated as a persona.

## Change discipline

1. Read and cite the exact affected protocol sections.
2. Keep the change inside its requested milestone and capability boundary.
3. Update deterministic code without duplicating the protocol.
4. Add positive/negative fixtures for each newly implemented normative rule.
5. Run the relevant checks above.
6. Report support and deferred behavior honestly; never create green placeholders.

See [AGENTS.md](../AGENTS.md) for artifact, trust-boundary, and generated-file rules.
