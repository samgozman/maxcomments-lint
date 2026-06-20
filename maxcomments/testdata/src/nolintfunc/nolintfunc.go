package nolintfunc

// SuppressedByDoc has three doc comment lines
// which exceed the budget, but the nolint
// directive below exempts the function.
//
//nolint:maxcomments
func SuppressedByDoc() {}

// SuppressedByTrailing also has three doc
// comment lines over the budget, but a
// trailing nolint on its signature exempts it.
func SuppressedByTrailing() {} //nolint:maxcomments

// Reported has three doc comment lines over
// the budget of two and carries no nolint
// directive, so it is reported.
func Reported() {} // want `function "Reported" has 3 comment lines, max allowed is 2`
