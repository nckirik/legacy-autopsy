#!/usr/bin/env bash
# Installs the local CDL VS Code extension (symlink by default).
#
# Usage: scripts/install-vscode-extension.sh [options]
#   --editor code|code-insiders|cursor|vscodium   target editor (default: code)
#   --extensions-dir DIR                          override the extensions directory
#                                                 (or set VSCODE_EXTENSIONS)
#   --copy                                        copy instead of symlinking
#   --force                                       replace an existing install
#   --uninstall                                   remove installed copies/symlinks
#   --dry-run                                     print actions without executing
set -euo pipefail
cd "$(dirname "$0")/.."

src_dir="editors/vscode"
editor="code"
extensions_dir=""
mode="link"
force=false
uninstall=false
dry_run=false

usage() {
  sed -n '2,12p' "$0" | sed 's/^# \{0,1\}//'
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --editor) editor="${2:?--editor requires a value}"; shift 2 ;;
    --extensions-dir) extensions_dir="${2:?--extensions-dir requires a value}"; shift 2 ;;
    --copy) mode="copy"; shift ;;
    --force) force=true; shift ;;
    --uninstall) uninstall=true; shift ;;
    --dry-run) dry_run=true; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

case "${editor}" in
  code) default_dir="${HOME}/.vscode/extensions" ;;
  code-insiders) default_dir="${HOME}/.vscode-insiders/extensions" ;;
  cursor) default_dir="${HOME}/.cursor/extensions" ;;
  vscodium) default_dir="${HOME}/.vscode-oss/extensions" ;;
  *)
    echo "unknown editor '${editor}'; supported: code, code-insiders, cursor, vscodium" >&2
    exit 2
    ;;
esac
extensions_dir="${extensions_dir:-${VSCODE_EXTENSIONS:-${default_dir}}}"

pkg="${src_dir}/package.json"
if [[ ! -f "${pkg}" ]]; then
  echo "missing ${pkg}: run this script from the repository" >&2
  exit 1
fi

field() {
  sed -n "s/^[[:space:]]*\"$1\":[[:space:]]*\"\\([^\"]*\\)\".*/\\1/p" "${pkg}" | head -n 1
}
name="$(field name)"
version="$(field version)"
publisher="$(field publisher)"
if [[ -z "${name}" || -z "${version}" || -z "${publisher}" ]]; then
  echo "${pkg}: missing name, version, or publisher" >&2
  exit 1
fi

target="${extensions_dir}/${publisher}.${name}-${version}"
src_abs="$(cd "${src_dir}" && pwd)"

run() {
  if ${dry_run}; then
    printf 'dry-run:'
    printf ' %q' "$@"
    printf '\n'
  else
    "$@"
  fi
}

if ${uninstall}; then
  if [[ ! -d "${extensions_dir}" ]]; then
    echo "extensions directory not found: ${extensions_dir}"
    exit 0
  fi
  shopt -s nullglob
  removed=0
  for dir in "${extensions_dir}/${publisher}.${name}-"*; do
    if [[ -L "${dir}" ]]; then
      run rm -f "${dir}"
      removed=$((removed + 1))
    elif [[ -f "${dir}/package.json" ]] &&
      grep -q "\"name\": \"${name}\"" "${dir}/package.json" &&
      grep -q "\"publisher\": \"${publisher}\"" "${dir}/package.json"; then
      run rm -rf "${dir}"
      removed=$((removed + 1))
    fi
  done
  if [[ ${removed} -eq 0 ]]; then
    echo "no installed ${publisher}.${name} extension found in ${extensions_dir}"
  elif ${dry_run}; then
    echo "dry-run: would remove ${removed} install(s) from ${extensions_dir}"
  else
    echo "removed ${removed} install(s) from ${extensions_dir}; reload ${editor} windows"
  fi
  exit 0
fi

if [[ -L "${target}" ]] && [[ "$(readlink "${target}")" == "${src_abs}" ]]; then
  echo "already linked: ${target} -> ${src_abs}"
  exit 0
fi
if [[ -e "${target}" || -L "${target}" ]]; then
  if ! ${force}; then
    echo "${target} exists; rerun with --force to replace it" >&2
    exit 1
  fi
  run rm -rf "${target}"
fi

run mkdir -p "${extensions_dir}"
if [[ "${mode}" == "copy" ]]; then
  run cp -R "${src_abs}" "${target}"
  message="copied ${publisher}.${name} ${version} to ${target}"
else
  run ln -s "${src_abs}" "${target}"
  message="linked ${publisher}.${name} ${version}: ${target} -> ${src_abs}"
fi
if ${dry_run}; then
  echo "dry-run: no changes written"
else
  echo "${message}"
  echo "reload ${editor} windows (Developer: Reload Window) to activate"
fi
