package gostr

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func testSuffixTree(
	construction func(string) SuffixTree,
	x string,
	t *testing.T) *SuffixTree {
	st := construction(x)

	for i := range LeafLabels(st.Root) {
		fmt.Println(i)
	}

	leaves := LeafLabels(st.Root)
	noLeaves := 0
	prev, ok := <-leaves
	if ok {
		noLeaves++
		for i := range leaves {
			if st.String[prev:] >= st.String[i:] {
				t.Errorf("We got the leaf %s before leaf %s.",
					ReplaceSentinel(st.String[prev:]), ReplaceSentinel(st.String[i:]))
			}
			noLeaves++
			prev = i
		}
	}
	if noLeaves != len(st.String) {
		t.Errorf("We didn't get all the leaves, we got %d but expected %d.\n",
			noLeaves, len(st.String))
	}

	return &st
}

func testSearchMatch(
	st *SuffixTree,
	p string,
	t *testing.T) {
	for i := range st.Search(p) {
		if st.String[i:i+len(p)] != p {
			t.Errorf("While searching for %s I found %s.",
				p, ReplaceSentinel(st.String[i:]))
		}
	}
}

func testSearchMismatch(
	st *SuffixTree,
	p string,
	t *testing.T) {
	res := st.Search(p)
	if _, ok := <-res; ok {
		t.Errorf("We shouldn't find '%s' in '%s.'", p, ReplaceSentinel(st.String))
	}
}

func testSearchMississippi(
	st *SuffixTree,
	t *testing.T) {
	testSearchMatch(st, "ssi", t)
	testSearchMismatch(st, "x", t)
	testSearchMismatch(st, "spi", t)
}

func TestNaiveConstruction(t *testing.T) {
	x := "mississippi"
	st := testSuffixTree(NaiveST, x, t)

	testSearchMississippi(st, t)

	f, err := os.Create("naive-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()

}

func TestMcCreightConstruction(t *testing.T) {
	x := "mississippi"
	st := testSuffixTree(McCreight, x, t)

	testSearchMississippi(st, t)

	f, err := os.Create("mccreight-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()
}
