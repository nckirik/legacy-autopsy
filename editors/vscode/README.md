# CDL for VS Code

Non-authoritative editor support for `.cdl` sources. The grammar and language
configuration describe the language surface; they never define protocol behavior
and are validated against the parser keyword surface by `internal/editing`.

## Install (development)

From the repository root:

```sh
scripts/install-vscode-extension.sh                 # symlink into ~/.vscode/extensions
scripts/install-vscode-extension.sh --copy          # copy instead of symlink
scripts/install-vscode-extension.sh --editor cursor  # code (default), code-insiders, cursor, vscodium
scripts/install-vscode-extension.sh --uninstall
```

Then reload the editor window. A distributable VSIX is out of scope for the local
installer; package one with `npx @vscode/vsce package` if needed.

## Scope

- syntax highlighting for declarations, opaque `GOAL`/`TEXT`/`REQUIRE` blocks
  (embedded Markdown), comments, strings, numbers, and operators;
- comment and bracket configuration.
