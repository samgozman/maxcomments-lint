package funcratio

func CommentHeavy() { // want `function "CommentHeavy" has 3 body comment lines for 5 code lines, max allowed is 1 \(func.ratio\)`
	// extra one
	// extra two
	a := 1
	b := 2
	_ = a + b
}

func CodeHeavy() {
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	f := 6
	// just one note for plenty of code
	_ = a + b + c + d + e + f
}

// TrailingCountsAsCode keeps its comments on the same lines as code. Those are
// code lines, not comment-only lines, so the 1-per-3 ratio holds and nothing is
// reported. (If trailing-comment lines were dropped from the code count, the
// shrunken denominator would falsely trip the ratio.)
func TrailingCountsAsCode() {
	a := 1 // note a
	b := 2 // note b
	c := 3
	_ = a + b + c
}
