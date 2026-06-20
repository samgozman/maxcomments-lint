package ignore

// checkedFunc has three doc comment lines
// over the budget of two, and this file is
// not ignored, so it is reported.
func checkedFunc() {} // want `function "checkedFunc" has 3 comment lines, max allowed is 2`
