package gostr

import (
	"math/rand"
	"testing"
	"time"
)

func benchmarkExactSearchRandom(
	n, m int,
	algo func(x, p string, cb func(int)),
	b *testing.B) {
	b.StopTimer()
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 10; i++ {
		x := randomString(n, "abcdefg", rng)
		for j := 0; j < 100; j++ {
			p := pickRandomSubstring(x, rng)
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
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 10; i++ {
		x := equalString(n)
		for j := 0; j < 100; j++ {
			p := pickRandomSubstring(x, rng)
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
