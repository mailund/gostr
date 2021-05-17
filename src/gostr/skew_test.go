package gostr

import (
	"math/rand"
	"testing"
	"time"
)

func TestLengthCalculations(t *testing.T) {
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

func testSASorted(x string, sa []int, t *testing.T) {
	for i := 1; i < len(sa); i++ {
		if x[sa[i-1]:] >= x[sa[i]:] {
			t.Errorf("Suffix array is not sorted! %q >= %q",
				x[sa[i-1]:], x[sa[i]:])
		}
	}
}

func TestMississippiSkew(t *testing.T) {
	x := "mississippi"
	testSASorted(x, Skew(x), t)

}

func TestRandomStringsSkew(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 30       // testing 30 random strings, enough to hit all mod 3 lengths
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		slen := rng.Intn(maxlen)
		x := randomString(slen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		testSASorted(x, Skew(x), t)
	}
}
