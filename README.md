# maxcomments-lint

[![CI](https://github.com/samgozman/maxcomments-lint/actions/workflows/ci.yml/badge.svg)](https://github.com/samgozman/maxcomments-lint/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/samgozman/maxcomments-lint.svg)](https://pkg.go.dev/github.com/samgozman/maxcomments-lint)
[![codecov](https://codecov.io/gh/samgozman/maxcomments-lint/graph/badge.svg?token=u4ogknIYUs)](https://codecov.io/gh/samgozman/maxcomments-lint)

A [golangci-lint](https://golangci-lint.run) plugin that limits how many
comment lines a function or file may contain.

It flags functions that narrate every line instead of being split up or
simplified:

```go
// BAD: a comment on nearly every line
func process(items []Item) error {
	// loop over all items
	for _, item := range items {
		// skip invalid items
		if !item.Valid() {
			continue
		}
		// save the item
		if err := save(item); err != nil {
			// return the error
			return err
		}
	}
	return nil
}
```

It checks comment *quantity*, not *style*. For style, use
[godot](https://github.com/tetafro/godot), `gocritic`, or `revive`.

The cap can be a flat number of comment lines or a *ratio* (at most one comment
line per N lines of code), enforced per function and/or per file. So a small
helper and a 200-line function can be held to proportional budgets instead of
the same fixed limit.

**Especially useful against AI-generated noise.** Coding assistants tend to
comment nearly every line. Capping comment density nudges that output back
toward self-explanatory code. And it's a compiled check, so it's
*deterministic*: it enforces the limit on every run, unlike a prompt or rules
file the model may ignore.

## How it works

Set a limit, get a warning when code exceeds it. Limits come in two flavours,
each usable per **function** and/or per **file**:

- **Hard cap:** a fixed maximum number of comment lines
  (`func.body-lines`, `func.doc-lines`, `file.lines`).
- **Ratio:** at most one comment line per *N* code lines
  (`func.ratio`, `file.ratio`); the budget is `floor(codeLines / N)`. A *code
  line* is a non-blank line that isn't itself a comment. Use `ratio-min-lines`
  to skip small scopes.

### What counts as a comment

- **Doc vs. body are counted separately.** A function's doc comment (the block
  above `func`) is governed by `func.doc-lines`; everything inside the body by
  `func.body-lines`/`func.ratio`. So a long, legitimate doc comment never trips
  the budget meant for line-by-line body narration.
- **Closures count on their own.** Each comment belongs to the *innermost*
  function containing it, so a closure's comments are never folded into the
  enclosing function.
- **Lines, not tokens.** A `/* */` block or a stack of `//` lines counts by its
  actual line span.
- **Directives are never counted.** `//nolint:...`, `//go:generate`,
  `//go:embed`, `//line ...` and the like are tooling instructions, not docs.
- **File scope** sums every comment in the file, including the package doc.

## Settings

Per-scope budgets are grouped under `func:` and `file:`; the remaining keys are
scope-independent and stay at the top level.

| Key               | Type     | Default        | Description                                                                                |
|-------------------|----------|----------------|--------------------------------------------------------------------------------------------|
| `func.body-lines` | int      | `0` (disabled) | Hard cap: max **body** comment lines allowed per function (doc comment excluded).          |
| `func.doc-lines`  | int      | `0` (disabled) | Hard cap: max **doc** comment lines allowed per function (the block above `func`).         |
| `func.ratio`      | int      | `0` (disabled) | Ratio: allow 1 **body** comment line per this many code lines, per function.               |
| `file.lines`      | int      | `0` (disabled) | Hard cap: max comment lines allowed per file.                                              |
| `file.ratio`      | int      | `0` (disabled) | Ratio: allow 1 comment line per this many code lines, per file.                            |
| `ratio-min-lines` | int      | `0` (no floor) | Skip the ratio checks for any scope with fewer than this many code lines.                  |
| `ignore`          | []string | `[]` (none)    | Regular expressions matched against each file's path; matching files are skipped entirely. |
| `check-generated` | bool     | `false`        | Check machine-generated files too. By default generated files are skipped.                 |

```yaml
settings:
  func:
    body-lines: 3
    doc-lines: 5
    ratio: 8
  file:
    lines: 120
    ratio: 5
  ratio-min-lines: 10
  ignore:
    - 'testdata/'
  check-generated: false
```

Every diagnostic ends with the name of the setting that triggered it (e.g.
`... max allowed is 3 (func.body-lines)`), so you can tell which knob to tune
when more than one check is enabled.

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

### Generated files

Machine-generated files are **skipped by default**: there's no point nudging
a code generator toward fewer comments. A file counts as generated when it
carries the [standard Go marker](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source)
(`// Code generated ... DO NOT EDIT.`) before its `package` clause.

To lint generated files like any other code, opt in:

```yaml
settings:
  check-generated: true
```

This is independent of `ignore`: explicit `ignore` patterns always apply, and
generated files are skipped on top of them unless `check-generated` is set.

## Using it in a project

Module plugins must be compiled into golangci-lint itself. See the
[Module Plugin System docs](https://golangci-lint.run/docs/plugins/module-plugins/).
You don't clone this repo; you reference it by version from your own project.

### 1. Add a `.custom-gcl.yml` to your project root

```yaml
version: v2.12.2   # the golangci-lint version to build
plugins:
  - module: github.com/samgozman/maxcomments-lint
    import: github.com/samgozman/maxcomments-lint/maxcomments
    version: v0.1.0   # pin a released tag of this plugin
```

### 2. Build a custom golangci-lint binary

```bash
golangci-lint custom   # reads .custom-gcl.yml, builds ./custom-gcl
```

### 3. Configure your project's `.golangci.yml`

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
        original-url: github.com/samgozman/maxcomments-lint
        settings:
          func:
            body-lines: 5
            doc-lines: 15
            # optional ratio mode (1 body comment line per 10 code lines):
            # ratio: 10
          file:
            lines: 150
            # ratio: 10
          # ratio-min-lines: 10
          ignore:
            - 'vendor/'
            - 'testdata/'
```

### 4. Run it

```bash
./custom-gcl run
```

## Editor integration

A custom golangci-lint binary works with your IDE just like the stock one. You
only have to point the IDE at *your* binary (`./bin/custom-gcl`) instead of the
one on your `PATH`. Build it first (see above), then configure the IDE.

<details>
<summary>JetBrains GoLand IDE</summary>

1. Open **Settings → Go → Linters**.
2. Tick **Execute 'golangci-lint run'** (and **'golangci-lint fmt'** if you
   want formatting too).
3. Set **Executable** to the absolute path of your custom binary, e.g.
   `/path/to/your/project/bin/custom-gcl`.
4. Tick **Use config** and point it at your `.golangci.yml`.
5. Click **OK**. `maxcomments` now shows up in the linters list and its
   warnings appear inline in the editor and in the **Problems** view.

![GoLand linter settings](docs/goland-linter-settings.png)

</details>

## Running in CI

CI is the same three steps as local use: install upstream golangci-lint, build
the custom binary from `.custom-gcl.yml`, then lint with it. See the
[`custom` command docs](https://golangci-lint.run/docs/plugins/module-plugins/)
for details, and this repo's own [`.github/workflows/ci.yml`](.github/workflows/ci.yml)
for a working example:

```yaml
- name: Install golangci-lint
  run: |
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
      | sh -s -- -b "$(go env GOPATH)/bin" v2.12.2

- name: Build the custom binary
  run: golangci-lint custom   # reads .custom-gcl.yml

- name: Lint
  run: ./bin/custom-gcl run
```

Keep the golangci-lint version pinned the same in `.custom-gcl.yml` and your CI.

## Development

```bash
go mod tidy
go test ./...
```

Each behaviour has its own `analysistest` fixture under
`maxcomments/testdata/src/` (one package per scenario: `funclines`,
`funcdoclines`, `directives`, `closures`, `funcratio`, `fileratio`,
`ratiomin`, `nolintfile`, `nolintfunc`, `ignore`, `generated`,
`generatedcheck`), alongside white-box unit tests for the pure helpers
(`isDirective`, `ratioViolation`, `nolintForMaxcomments`, `matchesAny`).

## Known gaps

- No autofix, by design. Fix flagged comments yourself, or ask an AI to
  summarise them down.
- A function-level `//nolint` on a `func` *signature line* is matched by line
  number, the same way golangci-lint applies `//nolint` to the diagnostics it
  receives. In the rare case where two `func` tokens share one physical source
  line (e.g. a one-line closure nested in another function), a trailing
  `//nolint` on that line suppresses both. Keep each function on its own line
  (which `gofmt` already does) to scope the directive precisely.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for the
workflow, project layout, and testing conventions. Notable changes are recorded
in [CHANGELOG.md](CHANGELOG.md).

## License

MIT. See [LICENSE](LICENSE).
