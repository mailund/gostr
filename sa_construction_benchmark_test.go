package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func benchmarkSAconstruction(b *testing.B, constr func(string) []int32, n int) {
	b.Helper()

	rng := test.NewRandomSeed(b)
	for i := 0; i < b.N; i++ {
		x := test.RandomStringN(n, "abcdefg", rng)
		constr(x)
	}
}

func BenchmarkSkew10000(b *testing.B) {
	benchmarkSAconstruction(b, Skew, 10000)
}

func BenchmarkSkew100000(b *testing.B) {
	benchmarkSAconstruction(b, Skew, 100000)
}

func BenchmarkSkew1000000(b *testing.B) {
	benchmarkSAconstruction(b, Skew, 1000000)
}

func BenchmarkSais10000(b *testing.B) {
	benchmarkSAconstruction(b, Sais, 10000)
}
func BenchmarkSais100000(b *testing.B) {
	benchmarkSAconstruction(b, Sais, 100000)
}

func BenchmarkSais1000000(b *testing.B) {
	benchmarkSAconstruction(b, Sais, 1000000)
}
