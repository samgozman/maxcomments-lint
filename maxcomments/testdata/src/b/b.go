// Package b contains test fixtures for //nolint exclusion.
package b

// NolintSkip has a two-line doc comment
// and a nolint directive in the body.
func NolintSkip() {
	//nolint:somecheck
}

// NolintInline has a two-line doc comment
// and an inline nolint on the opening brace.
func NolintInline() { //nolint:maxcomments
}

// OverBudgetThreeLineDoc has a three-line
// doc comment that exceeds the budget.
// No nolint is present.
func OverBudgetThreeLineDoc() {} // want `function "OverBudgetThreeLineDoc" has 3 comment lines, max allowed is 2`
