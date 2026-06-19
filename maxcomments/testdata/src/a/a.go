package a

// Short is fine.
func Short() {}

// LongDoc has too many
// comment lines for the
// configured budget.
func LongDoc() {} // want `function "LongDoc" has 3 comment lines, max allowed is 2`
