package gostr

import (
	"math/rand"
	"testing"
	"time"
)

func benchmarkSAconstruction(
	constr func(string, bool) []int,
	n int,
	b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < b.N; i++ {
		constr(randomString(n, "abcdefg", rng), false)
	}
}

func BenchmarkSkew10000(b *testing.B) {
	benchmarkSAconstruction(Skew, 10000, b)
}

func BenchmarkSkew100000(b *testing.B) {
	benchmarkSAconstruction(Skew, 100000, b)
}

func BenchmarkSkew1000000(b *testing.B) {
	benchmarkSAconstruction(Skew, 1000000, b)
}

func BenchmarkSais10000(b *testing.B) {
	benchmarkSAconstruction(Sais, 10000, b)
}
func BenchmarkSais100000(b *testing.B) {
	benchmarkSAconstruction(Sais, 100000, b)
}

func BenchmarkSais1000000(b *testing.B) {
	benchmarkSAconstruction(Sais, 1000000, b)
}

func Test_adccacacbbccdccdbccb(t *testing.T) {
	x := "adccacacbbccdccdbccb"
	testSASorted(x, Sais(x, false), t)
}
