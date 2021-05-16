package gostr

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func checkPathLabels(n STNode, algo string, st SuffixTree, t *testing.T) {
	switch v := n.(type) {
	case *innerNode:
		for _, child := range v.children {
			checkPathLabels(child, algo, st, t)
		}
	case *leafNode:
		if PathLabel(v, st.String) != st.String[v.leafIdx:] {
			t.Errorf(`%s(%s): the path label of leaf %d should be "%s" but is "%s"`,
				algo,
				ReplaceSentinel(st.String), v.leafIdx,
				ReplaceSentinel(st.String[v.leafIdx:]),
				ReplaceSentinel(PathLabel(v, st.String)))
		}
	}
}

func testSuffixTree(
	algo string,
	construction func(string) SuffixTree,
	x string,
	t *testing.T) *SuffixTree {
	st := construction(x)

	t.Logf("The leaf labels iterator is out of commission for now")

	leaves := []int{}
	LeafIndicesVisitor(
		st.Root,
		func(idx int) { leaves = append(leaves, idx) })

	noLeaves := 0
	if len(leaves) > 0 {
		prev := leaves[0]
		noLeaves++
		for i := 1; i < len(leaves); i++ {
			if st.String[prev:] >= st.String[leaves[i]:] {
				t.Errorf(`We got the leaf "%s" before leaf "%s" in %s("%s").`,
					ReplaceSentinel(st.String[prev:]),
					ReplaceSentinel(st.String[leaves[i]:]),
					algo, x)
			}
			noLeaves++
			prev = leaves[i]
		}
	}
	if noLeaves != len(st.String) {
		t.Errorf(`%s("%s"): We got %d leaves but expected %d.\n`,
			algo, x, noLeaves, len(st.String))
	}

	leaves = []int{}
	for iter := LeafIndicesIterator(st.Root); iter.HasMore(); {
		leaves = append(leaves, iter.Next())
	}
	noLeaves = 0
	if len(leaves) > 0 {
		prev := leaves[0]
		noLeaves++
		for i := 1; i < len(leaves); i++ {
			if st.String[prev:] >= st.String[leaves[i]:] {
				t.Errorf(`We got the leaf "%s" before leaf "%s" in %s("%s").`,
					ReplaceSentinel(st.String[prev:]),
					ReplaceSentinel(st.String[leaves[i]:]),
					algo, x)
			}
			noLeaves++
			prev = leaves[i]
		}
	}
	if noLeaves != len(st.String) {
		t.Errorf(`%s("%s"): We got %d leaves but expected %d.\n`,
			algo, x, noLeaves, len(st.String))
	}

	checkPathLabels(st.Root, algo, st, t)

	return &st
}

func testSearchMatch(
	algo string,
	st *SuffixTree,
	p string,
	t *testing.T) {
	for iter := st.Search(p); iter.HasMore(); {
		i := iter.Next()
		if st.String[i:i+len(p)] != p {
			t.Errorf(`%s("%s"): While searching for "%s" I found "%s".`,
				algo, ReplaceSentinel(st.String), p, ReplaceSentinel(st.String[i:]))
		}
	}
}

func testSearchMismatch(
	algo string,
	st *SuffixTree,
	p string,
	t *testing.T) {

	iter := st.Search(p)
	if iter.HasMore() {
		t.Errorf(`We shouldn't find "%s" in %s("%s").`,
			p, algo, ReplaceSentinel(st.String))
	}
}

func testSearchMississippi(
	algo string,
	st *SuffixTree,
	t *testing.T) {
	testSearchMatch(algo, st, "ssi", t)
	testSearchMismatch(algo, st, "x", t)
	testSearchMismatch(algo, st, "spi", t)
}

func TestNaiveConstruction(t *testing.T) {

	x := "mississippi"
	st := testSuffixTree("NaiveST", NaiveST, x, t)
	testSearchMississippi("NaiveST", st, t)
	f, err := os.Create("naive-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()
}

func TestMcCreightConstruction(t *testing.T) {

	x := "mississippi"
	st := testSuffixTree("McCreight", McCreight, x, t)
	testSearchMississippi("McCreight", st, t)
	f, err := os.Create("mccreight-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()
}

func randomString(n int, alpha string, rng *rand.Rand) string {
	runes := make([]byte, n)
	for i := 0; i < n; i++ {
		runes[i] = alpha[rng.Intn(len(alpha))]
	}
	return string(runes)
}

func TestSTRandomStrings(t *testing.T) {
	algos := []string{"NaiveST", "McCreight"}
	constructors := []func(string) SuffixTree{NaiveST, McCreight}

	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 10       // testing 10 random strings
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		slen := rng.Intn(maxlen)
		x := randomString(slen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		for i := range algos {
			testSuffixTree(algos[i], constructors[i], x, t)
		}
	}
}

func benchmarkConstruction(
	b *testing.B,
	constructor func(string) SuffixTree,
	n int) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := randomString(n, "abcdefg", rng)
	for i := 0; i < b.N; i++ {
		constructor(x)
	}
}

func BenchmarkNaive10000(b *testing.B)   { benchmarkConstruction(b, NaiveST, 10000) }
func BenchmarkNaive100000(b *testing.B)  { benchmarkConstruction(b, NaiveST, 100000) }
func BenchmarkNaive1000000(b *testing.B) { benchmarkConstruction(b, NaiveST, 1000000) }

func BenchmarkMcCreight10000(b *testing.B)   { benchmarkConstruction(b, McCreight, 10000) }
func BenchmarkMcCreight100000(b *testing.B)  { benchmarkConstruction(b, McCreight, 100000) }
func BenchmarkMcCreight1000000(b *testing.B) { benchmarkConstruction(b, McCreight, 1000000) }

func publicTraversal(n STNode) int {
	if IsLeaf(n) {
		return LeafIndex(n)
	} else {
		val := 0
		for _, child := range Children(n) {
			val += publicTraversal(child)
		}
		return val
	}
}

func privateTraversal(n STNode) int {
	switch v := n.(type) {
	case *leafNode:
		return v.leafIdx
	case *innerNode:
		val := 0
		for _, child := range v.getChildren() {
			val += privateTraversal(child)
		}
		return val
	}
	return -1
}

func visitorTraversal(n STNode) int {
	res := 0
	LeafIndicesVisitor(n,
		func(idx int) { res += idx })
	return res
}

func iteratorTraversal(n STNode) int {
	res := 0
	for iter := LeafIndicesIterator(n); iter.HasMore(); {
		res += iter.Next()
	}
	return res
}

func benchmarkTraversal(traversal func(STNode) int, b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := randomString(1000, "abcdefg", rng)
	st := McCreight(x)
	traversal(st.Root) // first traversal sorts the children...

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		traversal(st.Root)
	}
}

func BenchmarkPublicTraversal(b *testing.B) {
	benchmarkTraversal(publicTraversal, b)
}
func BenchmarkPrivateTraversal(b *testing.B) {
	benchmarkTraversal(privateTraversal, b)
}

func BenchmarkVisitorTraversal(b *testing.B) {
	benchmarkTraversal(visitorTraversal, b)
}

func BenchmarkIteratorTraversal(b *testing.B) {
	benchmarkTraversal(iteratorTraversal, b)
}
