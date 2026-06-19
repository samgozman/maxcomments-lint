# maxcomments-lint

A [golangci-lint](https://golangci-lint.run) module plugin that caps the
number of **comment lines** allowed per function and/or per file.

It doesn't judge comment *style* (use [godot](https://github.com/tetafro/godot)
or `gocritic`/`revive` for that) — only comment *quantity*. The idea is to
catch functions that have drifted into narrating every line instead of being
simplified or split up.

## How it counts

For each function declaration, it sums:
- the function's doc comment (the block directly above `func ...`)
- every comment group whose source range falls inside the function body

Multi-line `/* */` blocks and stacks of consecutive `//` lines are each
counted by their actual line span, not by "number of `//` tokens."

File-level counting (optional) sums every comment group in the file,
including the package doc comment.

## Settings

| Key | Type | Default | Description |
|---|---|---|---|
| `max-func-lines` | int | `0` (disabled) | Max comment lines allowed per function. |
| `max-file-lines` | int | `0` (disabled) | Max comment lines allowed per file. |

## Using it in a project

### 1. Build a custom golangci-lint binary

Module plugins must be compiled into golangci-lint itself — see the
[Module Plugin System docs](https://golangci-lint.run/docs/plugins/module-plugins/).

```bash
git clone https://github.com/samgozman/maxcomments-lint.git
cd maxcomments-lint
golangci-lint custom   # reads .custom-gcl.yml, builds ./bin/custom-gcl
```

### 2. Configure your project's `.golangci.yml`

```yaml
version: "2"

linters:
  default: none
  enable:
    - maxcomments
  settings:
    custom:
      maxcomments:
        type: module
        description: Limits the number of comment lines per function and per file.
        settings:
          max-func-lines: 15
          max-file-lines: 150
```

### 3. Run it

```bash
./bin/custom-gcl run
```

## Development

```bash
go mod tidy
go test ./...
```

`maxcomments/testdata/src/a/a.go` has a minimal `analysistest` fixture —
add more cases there as you harden the rule (e.g. methods on receivers,
nested funcs/closures, `//nolint` interaction, block comments).

## Known gaps (sketch-stage)

- Doesn't yet special-case `//nolint:` directive lines (they'll count
  toward the budget like any other comment).
- Closures and nested function literals are only checked if they're
  themselves top-level `FuncDecl`s — anonymous funcs inside a body are
  currently folded into their enclosing function's total, not checked
  individually.
- No autofix.

## License

MIT — see [LICENSE](LICENSE).
