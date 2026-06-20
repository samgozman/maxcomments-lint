package directives

// WithDirective's doc comment is not counted toward the body budget, and the
// directives in its body are never counted, so its two real body comment lines
// stay within the budget of two.
func WithDirective() {
	//go:noinline
	//nolint:something
	// real one
	// real two
	x := 1
	_ = x
}

// TooMany has three real body comment lines, which is over the budget of two,
// even though the directives below are never counted.
func TooMany() { // want `function "TooMany" has 3 body comment lines, max allowed is 2 \(func.body-lines\)`
	//go:noinline
	// real one (the trailing comment on the func line above also counts)
	// real two
	x := 1
	_ = x
}

// ColonProse has comments that look directive-ish but are prose: a space or a
// slash follows the colon, so they are counted, not skipped. With the trailing
// comment on the func line that is three body comment lines, over the budget.
func ColonProse() { // want `function "ColonProse" has 3 body comment lines, max allowed is 2 \(func.body-lines\)`
	//note: this is prose, not a directive
	//https://example.com is a url, not a directive
	x := 1
	_ = x
}
