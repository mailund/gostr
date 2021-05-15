package gostr

import (
	"testing"
)

func isPrefix(x, y string) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func TestMississippiBWT(t *testing.T) {
	x := "mississippi"
	sa := Skew(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	p := "is"
	L, R := BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		if !isPrefix(p, x[sa[i]:]) {
			t.Errorf(`We have an incorrect match: "%s" doesn't match "%s"`,
				p, x[sa[i]:])
		}
	}

}
