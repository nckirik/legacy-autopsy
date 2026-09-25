# Development

> Non-authoritative developer guide. `protocol.md` defines behavior; this file records the current build and command workflow.

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

The M0 fixture command always runs every registered bootstrap case and lists all complete Part 19.2 families as unsupported. The accepted `--implemented` flag is currently a redundant compatibility selector; public examples omit it.

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

The S1 §7.7 pilot is implemented for exactly one section. The frozen source language is
[`cdl.md`](cdl.md); execution contracts are [`eir.md`](eir.md),
[`execution-semantics.md`](execution-semantics.md), and
[`runtime-architecture.md`](runtime-architecture.md); sequencing is in the
[roadmap](roadmap.md); migration and freeze rules are in the
[migration audit](migration-audit.md). Everything outside §7.7 remains unimplemented, and
the pilot is not a conformance or Exit A/E claim.

```sh
go run ./cmd/legacy-autopsy spec compile
go run ./cmd/legacy-autopsy spec compile --golden examples/spec/golden/export-reconciliation.eir.json
go run ./cmd/legacy-autopsy spec render
go run ./cmd/legacy-autopsy spec render --out /tmp/export-reconciliation.light.prompt.md \
  --generated-at 2026-09-25T00:00:00Z
go run ./cmd/legacy-autopsy spec run --fixture examples/spec/fixtures/incomplete.json
go run ./cmd/legacy-autopsy spec run --fixture examples/spec/fixtures/reconciled.json --json
```

`spec compile` validates the identity ledger and prints the source fingerprint and EIR
hash; `spec render` renders one projection/channel with a provenance header; `spec run`
executes the native VM over a synthetic fixture and prints the canonical trace and hash.
All three default to the §7.7 source, ledger, and repository root, and accept
`--source`/`--ledger`/`--repo`.

The same checks are available as Go tests:

```sh
go test ./cdl/               # compiler, renderer, ledger, negative fixtures
go test ./internal/parity/   # test A (backend) and test B (protocol parity)
go test ./internal/analysis/ # classification coverage and metrics golden
```

S2 analysis lives in [`coverage-profile.md`](coverage-profile.md) and
[`capability-contracts.md`](capability-contracts.md). Regenerate the metrics golden with
`go test ./internal/analysis/ -update`; never hand-edit
[`analysis/metrics.golden.json`](../analysis/metrics.golden.json).

Generated pilot assets are golden-checked and must be regenerated deterministically,
never hand-edited:

```sh
go test ./cdl/ -run TestUpdateAssets -update
```

The M0 code, fixtures, and skill projections remain frozen: bug fixes that preserve all
fixture outcomes, additive pilot wiring, and no new protocol behavior or refactors.

## Markdown formatting

Non-Go Markdown is formatted with [dprint](https://dprint.dev) using the pinned plugin in
[`dprint.json`](../dprint.json). CI runs `dprint check`; formatting is deterministic and
version-pinned. dprint is a development prerequisite only; it is not part of the Go build
or runtime, and it never formats Go code.

```sh
scripts/format.sh           # format Go and Markdown
scripts/format.sh --check   # verify only
```

`protocol.md` is excluded because its bytes are fingerprint-bound. `docs/cdl.md` uses
intentional ASCII section rules, and `fixtures/`, `examples/spec/` (`.cdl` sources and
generated goldens), and `analysis/` are source or generated data — they are regenerated,
never reformatted. Go formatting stays with the Go toolchain.

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

For workspace smoke testing, use a disposable directory and run both `init` and `workspace check`. Unit tests verify that initialization fabricates no covered, confirmed, or passed state.

## Change discipline

1. Read and cite the exact affected protocol sections.
2. Keep the change inside its requested milestone and capability boundary.
3. Update deterministic code without duplicating the protocol.
4. Add positive/negative fixtures for each newly implemented normative rule.
5. Run the relevant checks above.
6. Report support and deferred behavior honestly; never create green placeholders.

See [AGENTS.md](../AGENTS.md) for artifact, trust-boundary, and generated-file rules.
