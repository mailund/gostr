package gostr_test

import (
	"fmt"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

func runExactBenchmarkRandom(algo exactAlgo, n int) func(*testing.B) {
	return func(b *testing.B) {
		b.StopTimer()
		rng := test.NewRandomSeed(b)
		for i := 0; i < 5; i++ {
			x := test.RandomStringN(n, "abcde", rng)
			p := test.PickRandomSubstring(x, rng)
			b.StartTimer()
			algo(x, p, func(int) {})
			b.StopTimer()
		}
	}
}

func Benchmark_ExactSearchRandomStrings(b *testing.B) {
	b.StopTimer()
	ns := []int{5000, 10000}
	for name, algo := range exact_algorithms {
		for _, n := range ns {
			b.Run(fmt.Sprintf("%s:n=%d", name, n),
				runExactBenchmarkRandom(algo, n))
		}
	}
}

func Benchmark_BMH_100000(b *testing.B) {
	runExactBenchmarkRandom(gostr.Bmh, 100000)(b)
}

func Benchmark_BMH_map_100000(b *testing.B) {
	runExactBenchmarkRandom(gostr.BmhWithMap, 100000)(b)
}

func Benchmark_BMH_String_100000(b *testing.B) {
	runExactBenchmarkRandom(gostr.BmhWithAlphabet, 100000)(b)
}

func Benchmark_BWT_100000(b *testing.B) {
	runExactBenchmarkRandom(bwtWrapper, 100000)(b)
}
