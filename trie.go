package gostr

import (
	"fmt"
	"io"
)

// Trie is both a trie and a node in a trie.
type Trie struct {
	// The parent of the node, unless this node is the root
	Parent *Trie
	// The suffix link of this node
	Suffix *Trie
	// Outlist for Aho-Corasick
	Outlist *Trie

	// The string label of this node. We assume that
	// only one string can be mapped to the same node.
	Label int
	// The children of this node
	Children [256]*Trie // If you map the alphabet, you could save some space...
}

// IsRoot returns true if and only if this trie is the root node.
func (t *Trie) IsRoot() bool {
	return t.Parent == nil
}

// setSuffixAndOutput sets the suffix link and output, provided
// that all nodes closer to the root have their links set.
func (t *Trie) setSuffixAndOutput(edge byte) {
	for slink := t.Parent; ; slink = slink.Suffix {
		if slink.Children[edge] != nil && slink.Children[edge] != t {
			// If we can extend, and it is not to ourselves (from our parent)
			// then we have the suffix link
			t.Suffix = slink.Children[edge]
			break
		}

		if slink.IsRoot() {
			// If we couldn't extend, but got to the root, then the suffix link
			// is the root
			t.Suffix = slink
			break
		}
	}

	if t.Suffix.Label >= 0 {
		t.Outlist = t.Suffix
	} else {
		t.Outlist = t.Suffix.Outlist
	}
}

// SetSuffixAndOutput sets the suffix link and output for
// Aho-Corasick
func (t *Trie) SetSuffixAndOutput() {
	queue := newTrieQueue(10) //nolint:gomnd // 10 is an arbitrary initial capacity

	queue.enqueue(t)

	for !queue.isEmpty() {
		t := queue.dequeue()
		for e, child := range &t.Children {
			if child != nil {
				child.setSuffixAndOutput(byte(e))
				queue.enqueue(child)
			}
		}
	}
}

// FindNode finds the node at the end of the string,
// or returns nil if there isn't one.
func (t *Trie) FindNode(p string) *Trie {
	switch {
	case t == nil:
		return nil

	case p == "":
		return t
	}

	return t.Children[p[0]].FindNode(p[1:])
}

// Contains check if the trie contains the string p
func (t *Trie) Contains(p string) bool {
	n := t.FindNode(p)
	return n != nil && n.Label >= 0
}

// toDot writes the trie rooted at t to the writer, but does not
// include the "digrahp { ... }" bit. ToDot does.
func (t *Trie) toDot(w io.Writer) {
	if t.Label >= 0 {
		fmt.Fprintf(w, "\"%p\"[label=\"%d\", shape=circle]\n", t, t.Label)
	} else {
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=point]\n", t)
	}

	if t.Suffix != nil && !t.Suffix.IsRoot() {
		fmt.Fprintf(w, `"%p" -> "%p"[style=dotted, color=red];`, t, t.Suffix)
	}

	if t.Outlist != nil {
		fmt.Fprintf(w, `"%p" -> "%p"[style=dashed, color=green];`, t, t.Outlist)
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

func insertInTrie(n *Trie, label int, x string) {
	if x == "" {
		n.Label = label
		return
	}

	if n.Children[x[0]] == nil {
		n.Children[x[0]] = &Trie{Parent: n, Label: -1}
	}

	insertInTrie(n.Children[x[0]], label, x[1:])
}

// BuildTrie builds a new trie from a sequence of strings.
func BuildTrie(strings []string) *Trie {
	root := &Trie{Label: -1}
	for i, x := range strings {
		insertInTrie(root, i, x)
	}

	root.SetSuffixAndOutput()

	return root
}
