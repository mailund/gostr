package gostr

import (
	"testing"
)

func TestMississippiBWT(t *testing.T) {
	x := "mississippi"
	sa := Skew(x, true) // We must include the sentinel handle sentinel here
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	p := "is"
	L, R := BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		testOccurrence(x, p, sa[i], t)
	}
}
