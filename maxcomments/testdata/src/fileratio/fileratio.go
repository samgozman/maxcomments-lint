package fileratio // want `file ".*fileratio.go" has 4 comment lines for 6 code lines, max allowed is 2 \(file.ratio\)`

// comment one
// comment two
// comment three
func f() {
	a := 1
	b := 2
	_ = a + b
}
