package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func testRandomSASorted(
	constr func(x string, senti bool) (sa []int),
	t *testing.T) {

	rng := test.NewRandomSeed(t)

	n := 30       // testing 30 random strings, enough to hit all mod 3 lengths
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		x := test.RandomStringRange(0, maxlen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		test.CheckSuffixArray(t, x, constr(x, false))
	}
}
