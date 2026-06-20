# Contributing

Thanks for taking the time to contribute! This is a small, focused
golangci-lint module plugin, so the workflow is intentionally simple.

## Getting started

```bash
go mod tidy
make test        # go test -race -v ./...
```

The [`Makefile`](Makefile) is the canonical entry point for every common task —
run `make help` to list the targets:

| Target             | What it does                                                     |
|--------------------|------------------------------------------------------------------|
| `make test`        | Run the test suite with the race detector.                       |
| `make cover-check` | Run tests and fail if coverage drops below 90%.                  |
| `make vet`         | `go vet ./...`.                                                  |
| `make tidy`        | `go mod tidy`.                                                   |
| `make build`       | Build a custom golangci-lint binary with the plugin compiled in. |
| `make run`         | Self-lint this repo with that custom binary.                     |

## How the linter is organised

The whole plugin lives in [`maxcomments/`](maxcomments):

- `plugin.go` — registers the plugin with golangci-lint's module system.
- `analyzer.go` — settings, the `analysis.Analyzer`, and the per-scope checks.
- `scopes.go` — collects functions/closures and attributes comments to the
  innermost scope.
- `counting.go` — counts comment vs. code lines and recognises directives.
- `nolint.go` — honours `//nolint` suppression.
- `ignore.go` — compiles and applies the `ignore` path patterns.

## Tests

Behaviour is verified with [`analysistest`](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest)
fixtures, one package per scenario under `maxcomments/testdata/src/`
(`funclines`, `funcdoclines`, `directives`, `closures`, `funcratio`,
`fileratio`, `ratiomin`, `nolintfile`, `nolintfunc`, `ignore`, `generated`,
`generatedcheck`). Pure helpers have white-box unit tests alongside the code.

When you change behaviour:

1. Add or update a `testdata` fixture with the expected `// want` diagnostics
   (the fixture *is* the spec — make it read clearly).
2. Add a unit test for any new pure helper.
3. Keep coverage at or above 90% (`make cover-check` enforces this in CI).

## Pull requests

- Keep changes small and focused; one behavioural change per PR where possible.
- Make sure `make cover-check`, `make vet`, and `make tidy` are clean — CI runs
  all three (plus a build-the-plugin-and-self-lint job) on every PR.
- Add an entry to [`CHANGELOG.md`](CHANGELOG.md) and, when behaviour changes,
  update the [`README.md`](README.md).

## Releasing (maintainers)

1. Add a new version heading with the changes in `CHANGELOG.md`.
2. Tag the release (`git tag vX.Y.Z && git push --tags`); users pin this tag in
   their `.custom-gcl.yml`.
3. Keep the golangci-lint version pinned in `.custom-gcl.yml` and the CI
   workflow in sync.
