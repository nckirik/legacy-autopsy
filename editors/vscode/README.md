# CDL for VS Code

Non-authoritative editor support for `.cdl` sources. The grammar and language
configuration describe the language surface; they never define protocol behavior
and are validated against the parser keyword surface by `internal/editing`.

## Install (development)

Copy or symlink this directory into `~/.vscode/extensions/cdl-language`, or use
the "Developer: Install Extension from Location..." command.

## Scope

- syntax highlighting for declarations, opaque `GOAL`/`TEXT`/`REQUIRE` blocks
  (embedded Markdown), comments, strings, numbers, and operators;
- comment and bracket configuration.
