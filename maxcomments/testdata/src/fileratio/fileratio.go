package fileratio // want `file ".*fileratio.go" has 4 comment lines for 5 code lines, max allowed is 1`

// comment one
// comment two
// comment three
func f() {
	a := 1
	b := 2
	_ = a + b
}
