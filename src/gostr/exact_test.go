package gostr

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func basicExactTests(
	algoname string,
	exact func(x, p string, cb func(int)),
	t *testing.T) {
	res := []int{}
	append_res := func(i int) {
		res = append(res, i)
	}

	type args struct {
		x        string
		p        string
		callback func(i int)
	}
	tests := []struct {
		name     string
		args     args
		expected []int
	}{
		{"aaa/",
			args{"aaa", "", append_res},
			[]int{0, 1, 2, 3},
		},
		{"aaa/a",
			args{"aaa", "a", append_res},
			[]int{0, 1, 2},
		},
		{"aaa/b",
			args{"aaa", "b", append_res},
			[]int{},
		},
		{"aaa/aa",
			args{"aaa", "aa", append_res},
			[]int{0, 1},
		},
		{"mississippi/ssi",
			args{"mississippi", "ssi", append_res},
			[]int{2, 5},
		},
		{"mississippi/ppi",
			args{"mississippi", "ppi", append_res},
			[]int{8},
		},
	}
	for _, tt := range tests {
		t.Run(algoname+":"+tt.name, func(t *testing.T) {
			exact(tt.args.x, tt.args.p, tt.args.callback)
			sort.Ints(res)
			if !equal_arrays(tt.expected, res) {
				t.Errorf("Searching for %s in %s and found %v (expected %v)\n",
					tt.args.p, tt.args.x, res, tt.expected)
			}
			res = []int{} // Reset
		})
	}
}

func testExactOccurrences(
	x, p string,
	algo func(x, p string, cb func(int)),
	t *testing.T) {

	t.Logf(`Searching for "%s" in "%s"\n`, p, x)
	algo(x, p, func(i int) { testOccurrence(x, p, i, t) })
}

func testExactOccurrencesRandomStrings(
	algo func(x, p string, cb func(int)),
	t *testing.T) {

	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d\n", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 100 // The length of the random strings
	for i := 0; i < 10; i++ {
		x := randomString(n, "abcdefg", rng)
		t.Logf(`Random string is x = %s\n`, x)

		t.Logf("Picking random patterns...\n")
		for j := 0; j < 10; j++ {
			// random patterns, they have a character that
			// doesn't exist in x, to make sure we test that
			m := rng.Intn(len(x))
			p := randomString(m, "abcdefgx", rng)
			t.Logf(`Random pattern: "%s"\n`, p)
			testExactOccurrences(x, p, algo, t)
		}

		t.Logf("Picking random prefixes...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomPrefix(x, rng)
			t.Logf(`Prefix: "%s"\n`, p)
			testExactOccurrences(x, p, algo, t)
		}

		t.Logf("Picking random suffixes...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomSuffix(x, rng)
			t.Logf(`Sufix: "%s"\n`, p)
			testExactOccurrences(x, p, algo, t)
		}

		t.Logf("Picking random substrings...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomSubstring(x, rng)
			t.Logf(`Substring: "%s"\n`, p)
			testExactOccurrences(x, p, algo, t)
		}
	}
}

func testEqualResults(
	x, p string,
	algo1 func(x, p string, cb func(int)),
	algo2 func(x, p string, cb func(int)),
	t *testing.T) {
	res1, res2 := []int{}, []int{}
	algo1(x, p, func(i int) { res1 = append(res1, i) })
	algo2(x, p, func(i int) { res2 = append(res2, i) })
	sort.Ints(res1)
	sort.Ints(res2)
	if !equal_arrays(res1, res2) {
		t.Errorf("Unequal results: %v %v\n", res1, res2)
	}
}

func testEqualRandomStrings(
	algo1 func(x, p string, cb func(int)),
	algo2 func(x, p string, cb func(int)),
	t *testing.T) {

	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d\n", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 100 // The length of the random strings
	for i := 0; i < 10; i++ {
		x := randomString(n, "abcdefg", rng)
		t.Logf(`Random string is x = %s\n`, x)

		t.Logf("Picking random patterns...\n")
		for j := 0; j < 10; j++ {
			// random patterns, they have a character that
			// doesn't exist in x, to make sure we test that
			m := rng.Intn(len(x))
			p := randomString(m, "abcdefgx", rng)
			t.Logf(`Random pattern: "%s"\n`, p)
			testEqualResults(x, p, algo1, algo2, t)
		}

		t.Logf("Picking random prefixes...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomPrefix(x, rng)
			t.Logf(`Prefix: "%s"\n`, p)
			testEqualResults(x, p, algo1, algo2, t)
		}

		t.Logf("Picking random suffixes...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomSuffix(x, rng)
			t.Logf(`Sufix: "%s"\n`, p)
			testEqualResults(x, p, algo1, algo2, t)
		}

		t.Logf("Picking random substrings...\n")
		for j := 0; j < 10; j++ {
			p := pickRandomSubstring(x, rng)
			t.Logf(`Substring: "%s"\n`, p)
			testEqualResults(x, p, algo1, algo2, t)
		}
	}
}

// Give BWT search the same interface as the other exact search
// algorithms
func bwtWrapper(x, p string, cb func(int)) {
	BwtPreprocess(x)(p, cb)
}

// Same for suffix trees...
func stWrapper(algo func(string) SuffixTree) func(x, p string, cb func(int)) {
	return func(x, p string, cb func(int)) {
		st := algo(x)
		st.Search(p, cb)
	}
}

var naiveStWrapper = stWrapper(NaiveST)
var mccreightWrapper = stWrapper(McCreight)

func Test_NaiveBasic(t *testing.T)        { basicExactTests("Naive", Naive, t) }
func Test_BorderSearchBasic(t *testing.T) { basicExactTests("BorderSearch", BorderSearch, t) }
func Test_BwtBasic(t *testing.T)          { basicExactTests("BWT", bwtWrapper, t) }
func Test_NaiveSTBasic(t *testing.T)      { basicExactTests("NaiveST", naiveStWrapper, t) }
func Test_McCreightBasic(t *testing.T)    { basicExactTests("McCreight", mccreightWrapper, t) }

func Test_NaiveRandom(t *testing.T)        { testExactOccurrencesRandomStrings(Naive, t) }
func Test_BorderSearchRandom(t *testing.T) { testExactOccurrencesRandomStrings(BorderSearch, t) }
func Test_BwtRandom(t *testing.T)          { testExactOccurrencesRandomStrings(bwtWrapper, t) }
func Test_NaiveSTRandom(t *testing.T)      { testExactOccurrencesRandomStrings(naiveStWrapper, t) }
func Test_McCreightRandom(t *testing.T)    { testExactOccurrencesRandomStrings(mccreightWrapper, t) }

func Test_EqualRandomStrings_Naive_BorderSearch(t *testing.T) {
	testEqualRandomStrings(Naive, BorderSearch, t)
}
func Test_EqualRandomStrings_Naive_BWT(t *testing.T) { testEqualRandomStrings(Naive, bwtWrapper, t) }
func Test_EqualRandomStrings_Naive_NaiveST(t *testing.T) {
	testEqualRandomStrings(Naive, naiveStWrapper, t)
}
func Test_EqualRandomStrings_Naive_McCreight(t *testing.T) {
	testEqualRandomStrings(Naive, mccreightWrapper, t)
}
