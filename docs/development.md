# Development

> Non-authoritative developer guide. `protocol.md` defines behavior; this file records the current build and command workflow.

## Prerequisites

Go 1.24 or newer is required. The built CLI has no external runtime dependencies. Tests and fixtures are local and require no network, production credentials, or raw exports.

## CLI

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

Reads are normalized and confined to the workspace, including symlink resolution; missing mandatory inputs fail closed. The emitted packet embeds exact workspace bytes and fingerprints plus routed protocol text. It is **not sanitized**: redirect it only to an operator-approved location outside the repository, review it before transmission, and never commit it or send it to an unapproved provider. Automatic transitive closure, stale comparison, mutation authorization, and semantic execution remain deferred.

## Local checks

```sh
test -z "$(gofmt -l cmd internal)"
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
