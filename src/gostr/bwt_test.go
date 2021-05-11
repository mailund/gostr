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
	sa := Skew(x)
	fmt.Println(sa)
	for _, i := range sa {
		fmt.Println(x[i:] + x[:i])
	}
	fmt.Println()

	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	p := "is"
	L, R := BwtSearch(x, p, ctab, otab)
	fmt.Println("Hits:", L, R)
	for i := L; i < R; i++ {
		fmt.Println(sa[i], x[sa[i]:])
		if !isPrefix(p, x[sa[i]:]) {
			t.Error("We have an incorrect match.")
		}
	}
}
