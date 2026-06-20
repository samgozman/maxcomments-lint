package directives

// WithDirective has two real comment lines, and the directives in its body
// must not push it over the budget of two.
func WithDirective() {
	//go:noinline
	//nolint:something
	x := 1
	_ = x
}

// TooMany has three real comment lines which is over the budget of two,
// even though
// directives elsewhere are never counted.
func TooMany() {} // want `function "TooMany" has 3 comment lines, max allowed is 2`
