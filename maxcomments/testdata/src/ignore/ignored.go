package ignore

// ignoredFunc lives in a file matched by the ignore patterns, so even though
// it has four doc comment lines over the budget of two,
// nothing
// is reported.
func ignoredFunc() {}
