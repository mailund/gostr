package test

import (
	"math/rand"
	"testing"
	"time"
)

func NewRandomSeed(tb testing.TB) *rand.Rand {
	seed := time.Now().UTC().UnixNano()
	tb.Logf("Random seed: %d", seed)
	return rand.New(rand.NewSource(seed))
}

/*
	RandomString constructs a random string of
	length in n, over the alphabet alpha.
*/
func RandomStringN(
	n int,
	alpha string,
	rng *rand.Rand) string {

	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = alpha[rng.Intn(len(alpha))]
	}
	return string(bytes)
}

/*
	RandomString constructs a random string of
	length in [min, max), over the alphabet alpha.
*/
func RandomStringRange(
	min, max int,
	alpha string,
	rng *rand.Rand) string {

	n := min + rng.Intn(max-min)
	return RandomStringN(n, alpha, rng)
}

/*
	SingletonString generates a string of length n
	consisting only of the letter a
*/
func SingletonString(n int, a byte) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = a
	}
	return string(bytes)
}

/*
	PickRandomPrefix returns a random prefix of the string x.
*/
func PickRandomPrefix(x string, rng *rand.Rand) string {
	return x[:rng.Intn(len(x))]
}

/*
	PickRandomSufix returns a random sufix of the string x.
*/
func PickRandomSuffix(x string, rng *rand.Rand) string {
	return x[rng.Intn(len(x)):]
}

/*
	PickRandomSubstring returns a random substring of the string x.
*/
func PickRandomSubstring(x string, rng *rand.Rand) string {
	i := rng.Intn(len(x) - 1)
	j := rng.Intn(len(x) - i)
	return x[i : i+j]
}

/*
	GenerateTestStrings generates strings of length between
	min and max and calls callback with them.
*/
func GenerateRandomTestStrings(
	min, max int,
	rng *rand.Rand,
	callback func(x string)) {
	n := 50 // number of random strings (maybe parameterise)
	for i := 0; i < n; i++ {
		callback(RandomStringRange(min, max, "abcdefg", rng))
	}
}

/*
	GenerateSingletonTestStrings generate singeton strings with length
	between min and max
*/
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

/*
	GenerateTestStrings generates strings of length between
	min and max and calls callback with them.
*/
func GenerateTestStrings(
	min, max int,
	rng *rand.Rand,
	callback func(x string)) {
	GenerateRandomTestStrings(min, max, rng, callback)
	GenerateSingletonTestStrings(min, max, rng, callback)
}

func GenerateTestStringsAndPatterns(
	min, max int,
	rng *rand.Rand,
	callback func(x, p string)) {

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
}
