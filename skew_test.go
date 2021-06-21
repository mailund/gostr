package gostr

import (
	"testing"
)

func Test_LengthCalculations(t *testing.T) {
	if sa3len(0) != 0 || sa12len(0) != 0 {
		t.Errorf("If the length is zero, both lengths should be zero")
	}
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
