package gostr_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

func benchmarkConstruction(
	b *testing.B,
	constructor func(string) *gostr.SuffixTree,
	n int) {
	b.Helper()

	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := test.RandomStringN(n, "abcdefg", rng)

	for i := 0; i < b.N; i++ {
		constructor(x)
	}
}

func BenchmarkNaive10000(b *testing.B)   { benchmarkConstruction(b, gostr.NaiveST, 10000) }
func BenchmarkNaive100000(b *testing.B)  { benchmarkConstruction(b, gostr.NaiveST, 100000) }
func BenchmarkNaive1000000(b *testing.B) { benchmarkConstruction(b, gostr.NaiveST, 1000000) }

func BenchmarkMcCreight10000(b *testing.B)   { benchmarkConstruction(b, gostr.McCreight, 10000) }
func BenchmarkMcCreight100000(b *testing.B)  { benchmarkConstruction(b, gostr.McCreight, 100000) }
func BenchmarkMcCreight1000000(b *testing.B) { benchmarkConstruction(b, gostr.McCreight, 1000000) }

func publicTraversal(n gostr.STNode) int {
	switch n.NodeType {
	case gostr.Leaf:
		return n.Leaf().Index

	case gostr.Inner:
		val := 0

		for _, child := range n.Inner().Children {
			if child.NodeType != gostr.UnInitialised {
				val += publicTraversal(child)
			}
		}

		return val

	case gostr.UnInitialised:
		// do nothing
		return 0
	}

	return 0 // Unreachable, but we need to return...
}

func visitorTraversal(n gostr.STNode) int {
	res := 0

	n.LeafIndices(func(idx int) { res += idx })

	return res
}

func Test_Traversal(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := test.RandomStringRange(500, 1000, "abcdefg", rng)
	st := gostr.NaiveST(x)

	public := publicTraversal(st.Root)
	visitor := 0

	st.Root.LeafIndices(func(i int) { visitor += i })

	if public != visitor {
		t.Errorf("The public/visitor traversal gave different resuls: %d/%d",
			public, visitor)
	}
}

func benchmarkTraversal(b *testing.B, traversal func(gostr.STNode) int) {
	b.Helper()

	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := test.RandomStringN(1000, "abcdefg", rng)
	st := gostr.NaiveST(x)

	traversal(st.Root) // first traversal sorts the children...

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		traversal(st.Root)
	}
}

func BenchmarkPublicTraversal(b *testing.B) {
	benchmarkTraversal(b, publicTraversal)
}

func BenchmarkVisitorTraversal(b *testing.B) {
	benchmarkTraversal(b, visitorTraversal)
}
