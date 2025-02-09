package gostr_test

import (
	"testing"

	"github.com/mailund/gostr/gostr"
	"github.com/mailund/gostr/testutils"
)

func benchmarkSAconstruction(b *testing.B, constr func(string) []int32, n int) {
	b.Helper()

	rng := testutils.NewRandomSeed(b)
	for i := 0; i < b.N; i++ {
		x := testutils.RandomStringN(n, "abcdefg", rng)
		constr(x)
	}
}

func BenchmarkSkew10000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Skew, 10000)
}

func BenchmarkSkew100000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Skew, 100000)
}

func BenchmarkSkew1000000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Skew, 1000000)
}

func BenchmarkSais10000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Sais, 10000)
}
func BenchmarkSais100000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Sais, 100000)
}

func BenchmarkSais1000000(b *testing.B) {
	benchmarkSAconstruction(b, gostr.Sais, 1000000)
}
