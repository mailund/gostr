package gostr

import (
	"fmt"
	"io"
)

// Trie is both a trie and a node in a trie.
type Trie struct {
	// The string label of this node. We assume that
	// only one string can be mapped to the same node.
	Label int
	// The children of this node
	Children [256]*Trie // If you map the alphabet, you could save some space...
}

func insertInTrie(n *Trie, label int, x string) {
	if x == "" {
		n.Label = label
		return
	}

	if n.Children[x[0]] == nil {
		n.Children[x[0]] = &Trie{Label: -1}
	}

	insertInTrie(n.Children[x[0]], label, x[1:])
}

// BuildTrie builds a new trie from a sequence of strings.
func BuildTrie(strings []string) *Trie {
	root := &Trie{Label: -1}
	for i, x := range strings {
		insertInTrie(root, i, x)
	}

	return root
}

// Contains check if the trie t contains the string p
func (t *Trie) Contains(p string) bool {
	if p == "" {
		return t.Label >= 0
	}

	child := t.Children[p[0]]
	if child == nil {
		return false
	}

	return child.Contains(p[1:])
}

// toDot writes the trie rooted at t to the writer, but does not
// include the digrahp { ... } bit. ToDot does.
func (t *Trie) toDot(w io.Writer) {
	if t.Label >= 0 {
		fmt.Fprintf(w, "\"%p\"[label=\"%d\", shape=circle]\n", t, t.Label)
	} else {
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=point]\n", t)
	}

	for i, child := range &t.Children {
		if child != nil {
			fmt.Fprintf(w, `"%p" -> "%p"[label="%c"];`, t, child, byte(i))
			child.toDot(w)
		}
	}
}

// ToDot writes a trie structure to the writer in Dot format.
func (t *Trie) ToDot(w io.Writer) {
	fmt.Fprintln(w, `digraph { rankdir="LR" `)
	t.toDot(w)
	fmt.Fprintln(w, "}")
}
