package gostr_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/mailund/gostr"
)

func TestTrie_Contains(t *testing.T) {
	input := []string{"foo", "foobar", "bar", "baz", "abc", "bca", "a", "cb", "b"}
	trie := gostr.BuildTrie(input)

	for _, x := range input {
		if !trie.Contains(x) {
			t.Errorf("the trie should contain %s", x)
		}
	}

	if trie.Contains("qux") {
		t.Error("the trie should not contain qux")
	}
}

func TestTrieToDot(t *testing.T) {
	input := []string{"foo", "foobar", "bar", "baz"}
	trie := gostr.BuildTrie(input)

	f, err := os.Create("trie.dot")
	if err != nil {
		log.Fatal(err)
	}

	trie.ToDot(f)
	f.Close()
}

func buildInEdgesTable(trie *gostr.Trie) map[*gostr.Trie]byte {
	var (
		rec func(trie *gostr.Trie) // forward decl.
		tbl = map[*gostr.Trie]byte{}
	)

	rec = func(trie *gostr.Trie) {
		for b, child := range &trie.Children {
			if child != nil {
				tbl[child] = byte(b)

				rec(child)
			}
		}
	}

	rec(trie)

	return tbl
}

func getTriePath(trie *gostr.Trie, tbl map[*gostr.Trie]byte) string {
	path := make([]byte, 0)

	for !trie.IsRoot() {
		path = append(path, tbl[trie])
		trie = trie.Parent
	}

	return gostr.ReverseString(string(path))
}

func TestTrieSuffixLinks(t *testing.T) {
	var (
		input = []string{"foo", "foobar", "bar", "baz"} // FIXME: generate some random strings
		trie  = gostr.BuildTrie(input)
		inTbl = buildInEdgesTable(trie)
		rec   func(*gostr.Trie) // forward decl
	)

	rec = func(trie *gostr.Trie) {
		if trie.Suffix == nil {
			if !trie.IsRoot() {
				t.Error("only the root should not have a suffix link")
			}
		} else {
			myPath := getTriePath(trie, inTbl)
			sufPath := getTriePath(trie.Suffix, inTbl)

			if !strings.HasSuffix(myPath, sufPath) {
				t.Errorf("The suffix doesn't have a suffix path: '%s' '%s'\n", myPath, sufPath)
			}
		}

		// test recursively
		for _, child := range &trie.Children {
			if child != nil {
				rec(child)
			}
		}
	}

	rec(trie) // now do the recursive testing
}
