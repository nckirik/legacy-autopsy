#!/usr/bin/env bash
# Full local gate. Mirrors the CI jobs for the deterministic surface.
set -euo pipefail
cd "$(dirname "$0")/.."

scripts/format.sh --check
go vet ./...
go test ./...
go build -o /tmp/legacy-autopsy ./cmd/legacy-autopsy
go run ./cmd/legacy-autopsy protocol check
go run ./cmd/legacy-autopsy validate routing
go run ./cmd/legacy-autopsy validate fixtures
echo "all checks passed"
