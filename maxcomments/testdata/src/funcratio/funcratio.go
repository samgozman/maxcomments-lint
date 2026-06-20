package funcratio

func CommentHeavy() { // want `function "CommentHeavy" has 3 body comment lines for 4 code lines, max allowed is 1 \(func.ratio\)`
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
