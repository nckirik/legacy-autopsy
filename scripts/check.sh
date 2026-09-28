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

extension_tmp="$(mktemp -d)"
VSCODE_EXTENSIONS="${extension_tmp}/ext" scripts/install-vscode-extension.sh --copy >/dev/null
test -f "${extension_tmp}/ext/legacy-autopsy.cdl-language-0.1.0/package.json"
VSCODE_EXTENSIONS="${extension_tmp}/ext" scripts/install-vscode-extension.sh --uninstall >/dev/null
test ! -e "${extension_tmp}/ext/legacy-autopsy.cdl-language-0.1.0"
rm -rf "${extension_tmp}"
echo "all checks passed"
