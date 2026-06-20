# maxcomments-lint

A [golangci-lint](https://golangci-lint.run) module plugin that caps the
number of **comment lines** allowed per function and/or per file.

It doesn't judge comment *style* (use [godot](https://github.com/tetafro/godot)
or `gocritic`/`revive` for that) — only comment *quantity*. The idea is to
catch functions that have drifted into narrating every line instead of being
simplified or split up.

## How it counts

Every function is checked independently. A "function" here is any named
function/method (`FuncDecl`) **and** any anonymous function literal/closure.
Each comment is attributed to the *innermost* function that contains it, so a
closure's comments are counted against the closure and never folded into the
function that encloses it.

For each function it sums:

- the function's doc comment (the block directly above `func ...`)
- every comment group whose source range falls inside the function body
  (minus anything that belongs to a nested closure)

Multi-line `/* */` blocks and stacks of consecutive `//` lines are each
counted by their actual line span, not by "number of `//` tokens."

File-level counting (optional) sums every comment group in the file,
including the package doc comment.

**Directive lines are never counted.** Machine directives such as
`//nolint:...`, `//go:generate`, `//go:embed`, and `//line ...` are tooling
instructions, not documentation, so they are excluded from every total.

## Two modes

You can use either or both, at function and/or file scope:

1. **Hard cap** — a fixed maximum number of comment lines
   (`max-func-lines`, `max-file-lines`).
2. **Ratio** — at most one comment line per *N* code lines
   (`max-func-ratio`, `max-file-ratio`). The allowed budget is
   `floor(codeLines / N)`. A **code line** is a physical, non-blank line that
   is not itself a comment line. Use `ratio-min-lines` to exempt small scopes.

## Settings

| Key               | Type     | Default        | Description                                                                                |
|-------------------|----------|----------------|--------------------------------------------------------------------------------------------|
| `max-func-lines`  | int      | `0` (disabled) | Hard cap: max comment lines allowed per function.                                          |
| `max-file-lines`  | int      | `0` (disabled) | Hard cap: max comment lines allowed per file.                                              |
| `max-func-ratio`  | int      | `0` (disabled) | Ratio: allow 1 comment line per this many code lines, per function.                        |
| `max-file-ratio`  | int      | `0` (disabled) | Ratio: allow 1 comment line per this many code lines, per file.                            |
| `ratio-min-lines` | int      | `0` (no floor) | Skip the ratio checks for any scope with fewer than this many code lines.                  |
| `ignore`          | []string | `[]` (none)    | Regular expressions matched against each file's path; matching files are skipped entirely. |

### Suppressing with `//nolint`

This plugin honours golangci-lint's `//nolint` directives directly:

- **Per function:** a `//nolint:maxcomments` (or bare `//nolint` / `//nolint:all`)
  in a function's doc comment or trailing on its `func` line suppresses both
  the cap and ratio reports for that function.
- **Per file:** the same directive placed **before the `package` clause** at
  the top of a file suppresses the file-level checks. Function-level checks
  still apply.

### Ignoring files and folders

```yaml
settings:
  ignore:
    - 'vendor/'
    - 'testdata/'
    - '_test\.go$'
    - '\.pb\.go$'   # generated code
```

An invalid regex is reported as an error rather than silently ignored.

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
          # optional ratio mode (1 comment line per 10 code lines):
          # max-func-ratio: 10
          # max-file-ratio: 10
          # ratio-min-lines: 10
          ignore:
            - 'vendor/'
            - 'testdata/'
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

Each behaviour has its own `analysistest` fixture under
`maxcomments/testdata/src/` (one package per scenario: `directives`,
`closures`, `funcratio`, `fileratio`, `ratiomin`, `nolintfile`, `nolintfunc`,
`ignore`), alongside white-box unit tests for the pure helpers
(`isDirective`, `ratioViolation`, `nolintForMaxcomments`, `matchesAny`).

## Known gaps

- No autofix.

## License

MIT — see [LICENSE](LICENSE).
