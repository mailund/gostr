package gostr_test // black box testing...

import (
	"sort"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

type exactFunc = func(x, p string) []int
type exactAlgo = func(x, p string, cb func(int))

func exactWrapper(algo exactAlgo) exactFunc {
	return func(x, p string) []int {
		res := []int{}
		algo(x, p, func(i int) {
			res = append(res, i)
		})
		sort.Ints(res)
		return res
	}
}

// Give BWT search the same interface as the other exact search
// algorithms
func bwtWrapper(x, p string, cb func(int)) {
	gostr.BwtPreprocess(x)(p, cb)
}

// Same for suffix trees...
func stWrapper(algo func(string) *gostr.SuffixTree) func(x, p string, cb func(int)) {
	return func(x, p string, cb func(int)) {
		algo(x).Search(p, cb)
	}
}

var exact_algorithms = map[string]exactAlgo{
	"Naive":        gostr.Naive,
	"BorderSearch": gostr.BorderSearch,
	"KMP":          gostr.Kmp,
	"BMH":          gostr.Bmh,
	"BMH-map":      gostr.BmhWithMap,
	"BMH-String":   gostr.BmhWithAlphabet,
	"BWT":          bwtWrapper,
	"ST-Naive":     stWrapper(gostr.NaiveST),
	"ST-McCreight": stWrapper(gostr.McCreight),
}

func runBasicExactTests(algo exactAlgo) func(*testing.T) {
	return func(t *testing.T) {
		type args struct {
			x string
			p string
		}
		tests := []struct {
			name     string
			args     args
			expected []int
		}{
			{"aaa/",
				args{"aaa", ""},
				[]int{0, 1, 2, 3},
			},
			{"aaa/a",
				args{"aaa", "a"},
				[]int{0, 1, 2},
			},
			{"aaa/b",
				args{"aaa", "b"},
				[]int{},
			},
			{"aaa/aa",
				args{"aaa", "aa"},
				[]int{0, 1},
			},
			{"aa/aaa",
				args{"aa", "aaa"},
				[]int{},
			},
			{"mississippi/ssi",
				args{"mississippi", "ssi"},
				[]int{2, 5},
			},
			{"mississippi/ppi",
				args{"mississippi", "ppi"},
				[]int{8},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hits := exactWrapper(algo)(tt.args.x, tt.args.p)
				if !test.IntArraysEqual(tt.expected, hits) {
					t.Errorf("Searching for %s in %s and found %v (expected %v)\n",
						tt.args.p, tt.args.x, hits, tt.expected)
				}
			})
		}
	}
}

func TestBasicExact(t *testing.T) {
	for name, algo := range exact_algorithms {
		t.Run(name, runBasicExactTests(algo))
	}
}

func runRandomExactOccurencesTests(algo exactAlgo) func(*testing.T) {
	return func(t *testing.T) {
		rng := test.NewRandomSeed(t)
		test.GenerateTestStringsAndPatterns(100, 200, rng,
			func(x, p string) {
				hits := exactWrapper(algo)(x, p)
				if !test.CheckAllOccurrences(t, x, p, hits) {
					t.Fatalf("Incorrect results for x = %q and p = %q", x, p)
				}
			})
	}
}

func TestRandomExactOccurences(t *testing.T) {
	for name, algo := range exact_algorithms {
		t.Run(name, runRandomExactOccurencesTests(algo))
	}
}

func runCheckExactOccurencesEqual(expected exactAlgo, algo exactAlgo) func(*testing.T) {
	return func(t *testing.T) {
		rng := test.NewRandomSeed(t)
		test.GenerateTestStringsAndPatterns(100, 200, rng,
			func(x, p string) {
				expected_hits := exactWrapper(expected)(x, p)
				hits := exactWrapper(algo)(x, p)
				if !test.IntArraysEqual(expected_hits, hits) {
					t.Fatalf("Expected and actual hits disagree %v vs %v",
						expected_hits, hits)
				}
			})
	}
}

func TestExactEqual(t *testing.T) {
	naive := exact_algorithms["Naive"]
	for name, algo := range exact_algorithms {
		t.Run(name, runCheckExactOccurencesEqual(naive, algo))
	}
}
