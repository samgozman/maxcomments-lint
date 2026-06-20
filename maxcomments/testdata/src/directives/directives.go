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
func TooMany() { // want `function "TooMany" has 3 body comment lines, max allowed is 2`
	//go:noinline
	// real one (the trailing comment on the func line above also counts)
	// real two
	x := 1
	_ = x
}
