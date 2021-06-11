package gostr

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/mailund/gostr/test"
)

func checkPathLabels(
	n STNode, algo string,
	st SuffixTree, t *testing.T) {
	switch n.NodeType {
	case Leaf:
		v := n.Leaf()
		if n.PathLabel(st.String, st.Alpha) != string(st.Alpha.RevmapBytes(st.String[v.Index:])) {
			t.Errorf(`%s(%s): the path label of leaf %d should be "%s" but is "%s"`,
				algo,
				st.Alpha.RevmapBytes(st.String), v.Index,
				st.Alpha.RevmapBytes(st.String[v.Index:]),
				n.PathLabel(st.String, st.Alpha))
		}

	case Inner:
		for _, child := range n.Inner().Children {
			if child.NodeType != UnInitialised {
				checkPathLabels(child, algo, st, t)
			}
		}
	}
}

func testSuffixTree(
	algo string,
	construction func(string) SuffixTree,
	x string,
	t *testing.T) *SuffixTree {
	st := construction(x)

	leaves := []int{}
	st.Root.LeafIndices(
		func(idx int) {
			leaves = append(leaves, idx)
		})

	noLeaves := 0
	if len(leaves) > 0 {
		prev := leaves[0]
		noLeaves++
		for i := 1; i < len(leaves); i++ {
			if string(st.String[prev:]) >= string(st.String[leaves[i]:]) {
				t.Errorf(`We got the leaf "%s" before leaf "%s" in %s("%s").`,
					st.Alpha.RevmapBytes(st.String[prev:]),
					st.Alpha.RevmapBytes(st.String[leaves[i]:]),
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
	st.Search(p, func(i int) {
		hit := st.Alpha.RevmapBytes(st.String[i : i+len(p)])
		if hit != p {
			t.Errorf(`%s("%s"): While searching for "%s" I found "%s".`,
				algo, st.Alpha.RevmapBytes(st.String),
				p, st.Alpha.RevmapBytes(st.String[i:]))
		}
	})
}

func testSearchMismatch(
	algo string,
	st *SuffixTree,
	p string,
	t *testing.T) {

	st.Search(p, func(i int) {
		t.Errorf(`We shouldn't find "%s" in %s("%s").`,
			p, algo, st.Alpha.RevmapBytes(st.String))
	})
}

func testSearchMississippi(
	algo string,
	st *SuffixTree,
	t *testing.T) {
	testSearchMatch(algo, st, "ssi", t)
	testSearchMismatch(algo, st, "x", t)
	testSearchMismatch(algo, st, "spi", t)
}

func Test_NaiveConstruction(t *testing.T) {

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

func Test_McCreightConstruction(t *testing.T) {

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

func Test_STRandomStrings(t *testing.T) {
	algos := []string{"NaiveST", "McCreight"}
	constructors := []func(string) SuffixTree{NaiveST, McCreight}

	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 10       // testing 10 random strings
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		x := test.RandomStringRange(0, maxlen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		for i := range algos {
			testSuffixTree(algos[i], constructors[i], x, t)
		}
	}
}
