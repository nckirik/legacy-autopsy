#!/usr/bin/env bash
# Formats Go and Markdown, or verifies formatting with --check.
set -euo pipefail
cd "$(dirname "$0")/.."

check=false
if [[ "${1:-}" == "--check" ]]; then
  check=true
fi

dprint_bin=""
if command -v dprint >/dev/null 2>&1; then
  dprint_bin="dprint"
elif [[ -x "${HOME}/.dprint/bin/dprint" ]]; then
  dprint_bin="${HOME}/.dprint/bin/dprint"
else
  echo "dprint not found: install it with the official installer (see docs/development.md)" >&2
  exit 1
fi

mapfile -t go_files < <(git ls-files '*.go')
if [[ ${#go_files[@]} -eq 0 ]]; then
  echo "no tracked Go files found" >&2
  exit 1
fi

if ${check}; then
  unformatted="$(gofmt -l "${go_files[@]}")"
  if [[ -n "${unformatted}" ]]; then
    printf 'Unformatted Go files:\n%s\n' "${unformatted}" >&2
    exit 1
  fi
  "${dprint_bin}" check
  echo "format check passed (gofmt, dprint)"
else
  gofmt -w "${go_files[@]}"
  "${dprint_bin}" fmt
  echo "formatted Go and Markdown"
fi
