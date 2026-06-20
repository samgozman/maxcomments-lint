package funclines

// Short is fine: its doc comment is not counted toward the body budget, and
// the single body comment is within budget.
func Short() {
	// only one body comment
	x := 1
	_ = x
}

// LongBody has a long doc comment that does NOT count toward the body budget,
// so only the three body comment lines below are measured against it.
func LongBody() { // want `function "LongBody" has 3 body comment lines, max allowed is 2 \(func.body-lines\)`
	// body one (the trailing comment on the func line above also counts)
	// body two
	x := 1
	_ = x
}
