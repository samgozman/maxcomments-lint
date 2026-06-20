//nolint:maxcomments
package nolintfile

// The file-level nolint above suppresses the file budget even though this
// package has more comment lines than the configured maximum. Function-level
// checks still apply, so the function below is still reported.

// FuncDocOver has three doc comment lines
// which exceed the function budget of two,
// and no nolint of its own.
func FuncDocOver() {} // want `function "FuncDocOver" has 3 doc comment lines, max allowed is 2 \(func.doc-lines\)`
