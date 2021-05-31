package gostr

import (
	"testing"
)

func Test_LengthCalculations(t *testing.T) {
	n12, n3 := 0, 0
	for lastIdx := 0; lastIdx < 100; lastIdx++ {
		if lastIdx%3 == 0 {
			n3++
		} else {
			n12++
		}
		n := lastIdx + 1
		if sa12len(n) != n12 {
			t.Errorf(`sa12len(%d) = %d (expected %d)`, n, sa12len(n), n12)
		}
		if sa3len(n) != n3 {
			t.Errorf(`sa3len(%d) = %d (expected %d)`, n, sa3len(n), n3)
		}
	}
}

func Test_SkewMississippi(t *testing.T) {
	x := "mississippi"
	testSASorted(x, Skew(x, false), t)
}

func Test_SkewRandomStrings(t *testing.T) {
	testRandomSASorted(Skew, t)
}