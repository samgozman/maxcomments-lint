# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0]

Initial release.

- Per-function and per-file comment-line budgets: hard caps (`func.body-lines`,
  `func.doc-lines`, `file.lines`) and ratios (`func.ratio`, `file.ratio`).
- Doc and body comments counted separately; machine directives never counted.
- `//nolint` suppression (per function and per file), `ignore` path patterns,
  and opt-in checking of generated files.

[0.1.0]: https://github.com/samgozman/maxcomments-lint/releases/tag/v0.1.0
