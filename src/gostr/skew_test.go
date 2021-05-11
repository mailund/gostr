package gostr

import (
	"fmt"
	"testing"
)

func TestMississippiSkew(t *testing.T) {
	x := "mississippi"
	sa := Skew(x)

	fmt.Println("Sorted suffixes:")
	for _, i := range sa {
		fmt.Printf("%d %s\n", i, x[i:])
	}

	fmt.Println("Testing that they are sorted...")
	for i := 1; i < len(sa); i++ {
		if x[sa[i-1]:] >= x[sa[i]:] {
			t.Errorf("Suffix array is not sorted! %q >= %q",
				x[sa[i-1]:], x[sa[i]:])
		}
	}
}
