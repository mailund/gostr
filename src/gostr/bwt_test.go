package gostr

import (
	"fmt"
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
	sa := Skew(x, true) // We must include the sentinel handle sentinel here
	for i, j := range sa {
		fmt.Println(i, x[j:]+"$"+x[:j])
	}
	ctab := Ctab(x)
	fmt.Println("got the c table")
	otab := Otab(x, sa, ctab)
	fmt.Println("got the o table")
	p := "is"
	L, R := BwtSearch(x, p, ctab, otab)
	fmt.Println("Searched and found", L, R)
	for i := L; i < R; i++ {
		fmt.Println(p, x[sa[i]:])
		if !isPrefix(p, x[sa[i]:]) {
			t.Errorf(`We have an incorrect match: "%s" doesn't match "%s"`,
				p, x[sa[i]:])
		}
	}

}
