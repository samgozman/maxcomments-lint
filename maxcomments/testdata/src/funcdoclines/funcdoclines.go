package funcdoclines

// ShortDoc is within budget.
func ShortDoc() {
	// body comments do not count toward the doc budget
	// so this pair is ignored here
	x := 1
	_ = x
}

// LongDoc has three doc comment lines
// which is over the doc budget of two,
// so it is reported (body comments are not counted here).
func LongDoc() { // want `function "LongDoc" has 3 doc comment lines, max allowed is 2`
	x := 1
	_ = x
}
