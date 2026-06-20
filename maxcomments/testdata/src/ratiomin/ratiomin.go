package ratiomin

// SmallExempt would break the 1-per-3 ratio but has too few code lines, so
// the ratio-min-lines floor exempts it.
func SmallExempt() {
	// note a
	// note b
	a := 1
	_ = a
}

func LargeFlagged() { // want `function "LargeFlagged" has 3 body comment lines for 8 code lines, max allowed is 2 \(func.ratio\)`
	// note one
	// note two
	b := 1
	c := 2
	d := 3
	e := 4
	f := 5
	_ = b + c + d + e + f
}
