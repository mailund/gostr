package test

import (
	"math/rand"
	"testing"
	"time"
)

// NewRandomSeed creates a new random number generator
func NewRandomSeed(tb testing.TB) *rand.Rand {
	tb.Helper()

	seed := time.Now().UTC().UnixNano()
	// maybe enable this again if it is useful,
	// but right now I don't want it in the benchmarks
	// tb.Logf("Random seed: %d", seed)
	return rand.New(rand.NewSource(seed))
}

// RandomStringN constructs a random string of length in n, over the alphabet alpha.
func RandomStringN(n int, alpha string, rng *rand.Rand) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = alpha[rng.Intn(len(alpha))]
	}

	return string(bytes)
}

// RandomStringRange constructs a random string of length in [min, max), over the alphabet alpha.
func RandomStringRange(min, max int, alpha string, rng *rand.Rand) string {
	n := min + rng.Intn(max-min)
	return RandomStringN(n, alpha, rng)
}

// FibonacciString returns the n'th Fibonacci string.
func FibonacciString(n int) string {
	const (
		fib0 = "a"
		fib1 = "b"
	)

	switch n {
	case 0:
		return fib0

	case 1:
		return fib1

	default:
		a, b := fib0, fib1
		for ; n > 1; n-- {
			a, b = b, a+b
		}

		return b
	}
}

// SingletonString generates a string of length n consisting only of the letter a
func SingletonString(n int, a byte) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = a
	}

	return string(bytes)
}

// PickRandomPrefix returns a random prefix of the string x.
func PickRandomPrefix(x string, rng *rand.Rand) string {
	return x[:rng.Intn(len(x))]
}

// PickRandomSuffix returns a random sufix of the string x.
func PickRandomSuffix(x string, rng *rand.Rand) string {
	return x[rng.Intn(len(x)):]
}

// PickRandomSubstring returns a random substring of the string x.
func PickRandomSubstring(x string, rng *rand.Rand) string {
	i := rng.Intn(len(x) - 1)
	j := rng.Intn(len(x) - i)

	return x[i : i+j]
}

// GenerateRandomTestStrings generates strings of length between
// min and max and calls callback with them.
func GenerateRandomTestStrings(
	min, max int,
	rng *rand.Rand,
	callback func(x string)) {
	n := 50 // number of random strings (maybe parameterise)
	for i := 0; i < n; i++ {
		callback(RandomStringRange(min, max, "abcdefg", rng))
	}
}

// GenerateSingletonTestStrings generate singeton strings with length
// between min and max
func GenerateSingletonTestStrings(
	min, max int,
	rng *rand.Rand,
	callback func(x string)) {
	n := 50 // number of random strings (maybe parameterise)
	for i := 0; i < n; i++ {
		// maybe it is a little overkill to generate this many
		// singletons?
		callback(SingletonString(min+rng.Intn(max-min), 'a'))
	}
}

// GenerateTestStrings generates strings of length between min
// and max and calls callback with them.
func GenerateTestStrings(
	min, max int,
	rng *rand.Rand,
	callback func(x string)) {
	GenerateRandomTestStrings(min, max, rng, callback)
	GenerateSingletonTestStrings(min, max, rng, callback)

	for n := 0; n < 10; n++ {
		callback(FibonacciString(n))
	}
}

// GenerateTestStringsAndPatterns generates a set of strings (x, p) where x is a string
// to search in and p is a string to search for.
func GenerateTestStringsAndPatterns(min, max int, rng *rand.Rand, callback func(x, p string)) {
	GenerateRandomTestStrings(min, max, rng,
		func(x string) {
			for j := 0; j < 10; j++ {
				// random patterns, they might have a character that
				// doesn't exist in x, to make sure we test that
				callback(x, RandomStringRange(0, len(x), "abcdefgx", rng))
			}

			for j := 0; j < 10; j++ {
				callback(x, PickRandomPrefix(x, rng))
			}

			for j := 0; j < 10; j++ {
				callback(x, PickRandomSuffix(x, rng))
			}

			for j := 0; j < 10; j++ {
				callback(x, PickRandomSubstring(x, rng))
			}
		})
	GenerateSingletonTestStrings(min, max, rng,
		func(x string) {
			for j := 0; j < 10; j++ {
				// random patterns, they might have a character that
				// doesn't exist in x, to make sure we test that
				callback(x, RandomStringRange(0, len(x), "abc", rng))
			}

			for j := 0; j < 10; j++ {
				callback(x, PickRandomPrefix(x, rng))
			}
		})

	for n := 3; n < 10; n++ {
		x := FibonacciString(n)

		for j := 0; j < 10; j++ {
			callback(x, PickRandomPrefix(x, rng))
		}

		for j := 0; j < 10; j++ {
			callback(x, PickRandomSuffix(x, rng))
		}

		for j := 0; j < 10; j++ {
			callback(x, PickRandomSubstring(x, rng))
		}
	}
}
