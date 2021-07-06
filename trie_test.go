package gostr_test

import (
	"log"
	"os"
	"testing"

	"github.com/mailund/gostr"
)

func TestTrie_Contains(t *testing.T) {
	strings := []string{"foo", "foobar", "bar", "baz", "abc", "bca", "a", "cb", "b"}
	trie := gostr.BuildTrie(strings)

	for _, x := range strings {
		if !trie.Contains(x) {
			t.Errorf("the trie should contain %s", x)
		}
	}

	if trie.Contains("qux") {
		t.Error("the trie should not contain qux")
	}
}

func TestTrieToDot(t *testing.T) {
	strings := []string{"foo", "foobar", "bar", "baz"}
	trie := gostr.BuildTrie(strings)

	f, err := os.Create("trie.dot")
	if err != nil {
		log.Fatal(err)
	}

	trie.ToDot(f)
	f.Close()
}
