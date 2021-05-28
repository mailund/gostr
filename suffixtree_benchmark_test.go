package gostr

import (
	"math/rand"
	"testing"
	"time"
)

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
	LeafIndices(n,
		func(idx int) { res += idx })
	return res
}

func Test_Traversal(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := randomString(1000, "abcdefg", rng)
	st := McCreight(x)

	public := publicTraversal(st.Root)
	private := privateTraversal(st.Root)
	visitor := 0
	LeafIndices(st.Root, func(i int) { visitor += i })

	if public != private || public != visitor {
		t.Errorf("The public/private/visitor traversal gave different resuls: %d/%d/%d",
			public, private, visitor)
	}
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
