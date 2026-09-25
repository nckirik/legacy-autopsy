#!/usr/bin/env bash
# Regenerates deterministic golden assets. Never hand-edit goldens.
set -euo pipefail
cd "$(dirname "$0")/.."

go test ./cdl -run TestUpdateAssets -update
go test ./internal/analysis -update
echo "goldens regenerated; review the diff before committing"
