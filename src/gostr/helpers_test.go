package gostr

import (
	"math/rand"
	"testing"
	"time"
)

func equal_arrays(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isPrefix(x, y string) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func testOccurrence(x, p string, i int, t *testing.T) {
	if !isPrefix(p, x[i:]) {
		t.Errorf(`We have an incorrect match: "%s" doesn't match "%s"`,
			p, x[i:])
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

func randomString(n int, alpha string, rng *rand.Rand) string {
	runes := make([]byte, n)
	for i := 0; i < n; i++ {
		runes[i] = alpha[rng.Intn(len(alpha))]
	}
	return string(runes)
}

func pickRandomPrefix(x string, rng *rand.Rand) string {
	return x[:rng.Intn(len(x))]
}

func pickRandomSuffix(x string, rng *rand.Rand) string {
	return x[rng.Intn(len(x)):]
}

func pickRandomSubstring(x string, rng *rand.Rand) string {
	i := rng.Intn(len(x) - 1)
	j := rng.Intn(len(x) - i)
	return x[i : i+j]
}

func testRandomSASorted(
	constr func(x string, senti bool) (sa []int),
	t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 30       // testing 30 random strings, enough to hit all mod 3 lengths
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		slen := rng.Intn(maxlen)
		x := randomString(slen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		testSASorted(x, constr(x, false), t)
	}
}
