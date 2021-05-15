package gostr

import (
	"testing"
)

func TestMississippiSkew(t *testing.T) {
	x := "mississippi"
	sa := Skew(x)

	for i := 1; i < len(sa); i++ {
		if x[sa[i-1]:] >= x[sa[i]:] {
			t.Errorf("Suffix array is not sorted! %q >= %q",
				x[sa[i-1]:], x[sa[i]:])
		}
	}
}
