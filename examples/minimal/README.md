# Minimal example

> Non-authoritative and not a conformance fixture.

Run the implemented minimal flow from the repository root:

```sh
go run ./cmd/legacy-autopsy doctor
go run ./cmd/legacy-autopsy init --workspace /tmp/legacy-autopsy-minimal/.extracted
go run ./cmd/legacy-autopsy workspace check --workspace /tmp/legacy-autopsy-minimal/.extracted
go run ./cmd/legacy-autopsy context --mode Preflight --system-namespace example --iteration ALFA --invocation-id INV-000000000001 --scope preflight --environment-snapshot dev/snapshot-1 --write-target 0A-PREFLIGHT.md --write-target 0G-DECONSTRUCTION-STATE.md --forbidden-write-target "all other workspace paths" --workspace /tmp/legacy-autopsy-minimal/.extracted
```

The initializer emits only headings and explicit empty scaffolds: it invents no source fact, claim, persona, ownership, `[C-COVERED]`, `[B-CONFIRMED]`, confirmation, passed gate, or certified bundle. This example remains separate from `fixtures/`; examples teach usage, while fixtures prove deterministic acceptance/rejection rules.
