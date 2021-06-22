package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func benchmarkSAconstruction(
	constr func(string) []int32,
	n int,
	b *testing.B) {
	rng := test.NewRandomSeed(b)
	for i := 0; i < b.N; i++ {
		x := test.RandomStringN(n, "abcdefg", rng)
		constr(x)
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
