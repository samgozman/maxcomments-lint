package closures

// Outer stays within budget: the closure's comments are not folded into it.
func Outer() {
	// only this line counts toward Outer
	bad := func() { // want `function literal has 3 body comment lines, max allowed is 2 \(func.body-lines\)`
		// the want-directive line above is itself a comment line, and these
		// two lines make three counted for the closure alone
		_ = 0
	}
	_ = bad
}

// Nested checks that a closure inside a closure is counted on its own.
func Nested() {
	_ = func() {
		_ = func() { // want `function literal has 3 body comment lines, max allowed is 2 \(func.body-lines\)`
			// the want line above counts, plus these two lines, makes three
			// attributed to the innermost closure only
			_ = 0
		}
	}
}
