package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func benchmarkExactSearchRandom(
	n, m int,
	algo func(x, p string, cb func(int)),
	b *testing.B) {
	b.StopTimer()
	rng := test.NewRandomSeed(b)
	for i := 0; i < 10; i++ {
		x := test.RandomStringN(n, "abcdefg", rng)
		for j := 0; j < 100; j++ {
			p := test.PickRandomSubstring(x, rng)
			b.StartTimer()
			algo(x, p, func(int) {})
			b.StopTimer()
		}
	}
}

func benchmarkExactSearchEqual(
	n, m int,
	algo func(x, p string, cb func(int)),
	b *testing.B) {
	b.StopTimer()
	rng := test.NewRandomSeed(b)
	for i := 0; i < 10; i++ {
		x := test.SingletonString(n, 'a')
		for j := 0; j < 100; j++ {
			p := test.PickRandomSubstring(x, rng)
			b.StartTimer()
			algo(x, p, func(int) {})
			b.StopTimer()
		}
	}
}
func Benchmark_ExactSearch_Random_Naive(b *testing.B) {
	benchmarkExactSearchRandom(10000, 100, Naive, b)
}
func Benchmark_ExactSearch_Equal_Naive(b *testing.B) {
	benchmarkExactSearchEqual(10000, 100, Naive, b)
}

func Benchmark_ExactSearch_Random_BorderSearch(b *testing.B) {
	benchmarkExactSearchRandom(10000, 100, BorderSearch, b)
}
func Benchmark_ExactSearch_Equal_BorderSearch(b *testing.B) {
	benchmarkExactSearchEqual(10000, 100, BorderSearch, b)
}
