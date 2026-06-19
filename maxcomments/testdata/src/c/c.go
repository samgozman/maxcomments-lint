// Package c contains test fixtures for anonymous function literal checking.
package c

// OuterWithLiteral has a two-line doc comment
// and contains a func literal with too many comments.
func OuterWithLiteral() {
	fn := func() { // want `function literal has 4 comment lines, max allowed is 2`
		// comment one
		// comment two
		// comment three
	}
	_ = fn
}

// OuterExceedsWithBodyComment has a two-line
// doc comment and a body comment outside the literal.
func OuterExceedsWithBodyComment() { // want `function "OuterExceedsWithBodyComment" has 4 comment lines, max allowed is 2`
	// outer body comment
	fn := func() {
		// inner literal comment — excluded from the outer count
	}
	_ = fn
}
