package gostr

import (
	"fmt"
	"io"
	"strings"
)

type interval struct {
	i, j int
}

// With slicing, we get a substring in constant
// time
func (r interval) substr(x string) string {
	return x[r.i:r.j]
}

func (r interval) length() int {
	return r.j - r.i
}

func (r interval) slice(i, j int) interval {
	if j < 0 {
		j = r.length()
	}
	if i < 0 || i > r.length() || j > r.length() {
		panic("Slicing outside of range")
	}
	if j < i {
		panic("Interval must end after it begins")
	}
	return interval{r.i + i, r.i + j}
}

type STNode struct {
	interval
	LeafIdx    int
	Parent     *STNode
	suffixLink *STNode
	Children   map[byte]*STNode
}

func newInner(inter interval) *STNode {
	return &STNode{
		interval: inter,
		LeafIdx:  -1,
		Children: map[byte]*STNode{}}
}

func newLeaf(idx int, inter interval) *STNode {
	return &STNode{
		interval: inter,
		LeafIdx:  idx}
}

func (n *STNode) IsSTInner() bool {
	return n.LeafIdx == -1
}

func (n *STNode) IsSTLeaf() bool {
	return n.LeafIdx != -1
}

func (n *STNode) EdgeLabel(x string) string {
	return n.interval.substr(x)
}

func (n *STNode) addChild(child *STNode, x string) {
	n.Children[x[child.i]] = child
	child.Parent = n
}

func replaceSentinel(x string) string {
	// Need this one for dot output
	return strings.ReplaceAll(x, "\x00", "†")
}

func (n *STNode) ToDot(x string, w io.Writer) {
	if n.Parent == nil {
		// Root
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", n)
	} else {
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			n.Parent, n, replaceSentinel(n.EdgeLabel(x)))
		if n.IsSTLeaf() {
			fmt.Fprintf(w, "\"%p\"[label=%d]\n", n, n.LeafIdx)
		} else {
			fmt.Fprintf(w, "\"%p\"[shape=point]\n", n)
		}
	}
	for _, child := range n.Children {
		child.ToDot(x, w)
	}
}

type SuffixTree struct {
	String string
	Root   *STNode
}

func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, "digraph {")
	st.Root.ToDot(st.String, w)
	fmt.Fprintln(w, "}")
}

// This function doesn't really belong with suffix trees,
// but this is where I need it...
func min(vars ...int) int {
	m := vars[0]
	for _, n := range vars {
		if n < m {
			m = n
		}
	}
	return m
}

func lenSharedPrefix(i1, i2 interval, x string) int {
	i, n := 0, min(i1.length(), i2.length())
	for ; i < n; i++ {
		if x[i1.i+i] != x[i2.i+i] {
			break
		}
	}
	return i
}

func sscan(n *STNode, inter interval, x string) (*STNode, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	v, ok := n.Children[x[inter.i]]
	if !ok {
		return n, 0, inter
	}
	i := lenSharedPrefix(v.interval, inter, x)
	if i == inter.length() || i < v.interval.length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, inter.slice(i, -1), x)
}

func fscan(n *STNode, inter interval, x string) (*STNode, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	v, ok := n.Children[x[inter.i]]
	if !ok {
		panic("With fscan there should always be an out-edge")
	}
	i := min(v.interval.length(), inter.length())
	if i == inter.length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, inter.slice(i, -1), x)
}

func breakEdge(n *STNode, depth, leafidx int, y interval, x string) *STNode {
	if n.Parent == nil {
		panic("A node must have a parent when we break its edge.")
	}
	new_node := newInner(n.interval.slice(0, depth))
	n.Parent.addChild(new_node, x)
	new_leaf := newLeaf(leafidx, y)
	n.i += depth
	new_node.addChild(new_leaf, x)
	new_node.addChild(n, x)
	return new_leaf
}

func NaiveST(x string) SuffixTree {
	// Add sentinel
	x += "\x00"
	root := newInner(interval{0, 0})
	for i := 0; i < len(x); i++ {
		v, j, y := sscan(root, interval{i, len(x)}, x)
		if j == 0 {
			v.addChild(newLeaf(i, y), x)
		} else {
			breakEdge(v, j, i, y.slice(j, -1), x)
		}
	}
	return SuffixTree{x, root}
}

func (v *STNode) suffix() interval {
	// If v's parent is the root, chop
	// off one index
	if v.Parent.Parent == nil {
		return v.interval.slice(1, -1)
	} else {
		return v.interval
	}
}

func McCreight(x string) SuffixTree {
	x += "\x00"
	root := newInner(interval{0, 0})
	root.suffixLink = root
	currLeaf := newLeaf(0, interval{0, len(x)})
	root.addChild(currLeaf, x)

	// The bits of the suffix we need to search for
	var y, z interval
	// ynode is the node we get to when searching for y
	var ynode *STNode

	for i := 1; i < len(x); i++ {
		p := currLeaf.Parent

		if p.suffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.suffix()
			ynode = p.suffixLink
		} else {
			pp := p.Parent
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.interval
			ynode, depth, _ := fscan(pp.suffixLink, y, x)
			if depth < ynode.interval.length() {
				// ended on an edge
				currLeaf = breakEdge(ynode, depth, i, z, x)
				p.suffixLink = currLeaf.Parent
				continue // Go to next suffix, we are done here
			}

			// Remember p's suffix link for later...
			p.suffixLink = ynode
		}

		// This is the slow scan part, from ynode and the rest
		// of the suffix, which is z.
		n, depth, w := sscan(ynode, z, x)
		if depth == 0 {
			// Landed on a node
			currLeaf = newLeaf(i, w)
			n.addChild(currLeaf, x)
		} else {
			// Landed on an edge
			currLeaf = breakEdge(n, depth, i, w.slice(depth, -1), x)
		}
	}
	return SuffixTree{x, root}
}
