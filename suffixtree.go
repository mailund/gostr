package gostr

import (
	"fmt"
	"io"
	"strings"
	"unsafe"
)

// -- Substrings represented as intervals -----------------
type Range struct {
	From, To int
}

func (r Range) length() int {
	return r.To - r.From
}

// Chump off the first k letters of an interval
func (r Range) chump(k int) Range {
	return Range{r.From + k, r.To}
}

func (r Range) prefix(k int) Range {
	return Range{r.From, r.From + k}
}

func (r Range) substr(x []byte, alpha *Alphabet) string {
	return alpha.RevmapBytes(x[r.From:r.To])
}

type STNodeType int

const (
	UnInitialised STNodeType = iota
	Leaf          STNodeType = iota
	Inner         STNodeType = iota
)

type SharedNode struct {
	Range
	Parent *InnerNode
}

type LeafNode struct {
	SharedNode
	Index int
}

type InnerNode struct {
	SharedNode
	SuffixLink *InnerNode
	Children   []STNode
}

type STNode struct {
	NodeType STNodeType
	ptr      unsafe.Pointer
}

func (n STNode) isNil() bool {
	return n.NodeType == UnInitialised
}

func (n STNode) Shared() *SharedNode {
	return (*SharedNode)(n.ptr)
}

func (n STNode) Leaf() *LeafNode {
	return (*LeafNode)(n.ptr)
}

func (n STNode) Inner() *InnerNode {
	return (*InnerNode)(n.ptr)
}

func wrapLeaf(n *LeafNode) STNode {
	return STNode{
		NodeType: Leaf,
		ptr:      unsafe.Pointer(n),
	}
}

func wrapInner(n *InnerNode) STNode {
	return STNode{
		NodeType: Inner,
		ptr:      unsafe.Pointer(n),
	}
}

func (n STNode) EdgeLabel(x []byte, alpha *Alphabet) string {
	return n.Shared().substr(x, alpha)
}

func (n STNode) PathLabel(x []byte, alpha *Alphabet) string {
	labels := []string{n.Shared().substr(x, alpha)}
	for p := n.Shared().Parent; p != nil; p = p.Parent {
		labels = append(labels, p.substr(x, alpha))
	}
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}
	return strings.Join(labels, "")
}

func (n STNode) LeafIndices(visitor func(int)) {
	switch n.NodeType {
	case Leaf:
		visitor(n.Leaf().Index)
	case Inner:
		for _, child := range n.Inner().Children {
			child.LeafIndices(visitor)
		}
	}
}

func (n STNode) ToDot(x []byte, alpha *Alphabet, w io.Writer) {
	switch n.NodeType {
	case Leaf:
		v := n.Leaf()
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			v.Parent, v, v.substr(x, alpha))
		fmt.Fprintf(w, "\"%p\"[label=%d]\n", v, v.Index)

	case Inner:
		v := n.Inner()
		if v.Parent == nil {
			// Root
			fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", v)
		} else {
			fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
				v.Parent, v, v.substr(x, alpha))
			fmt.Fprintf(w, "\"%p\"[shape=point]\n", v)
		}
		if v.SuffixLink != nil {
			fmt.Fprintf(w, `"%p" -> "%p"[style=dotted, color=red];`, v, v.SuffixLink)
		}
		for _, child := range v.Children {
			child.ToDot(x, alpha, w)
		}
	}
}

func (n *InnerNode) addChild(child STNode, x []byte) {
	n.Children[x[child.Shared().From]] = child
	child.Shared().Parent = n
}

// -- Suffix tree --------------------------

type SuffixTree struct {
	Alpha  *Alphabet
	String []byte
	Root   STNode
}

func (st *SuffixTree) newLeaf(idx int, r Range) STNode {
	leaf := LeafNode{
		SharedNode: SharedNode{Range: r},
		Index:      idx}
	return wrapLeaf(&leaf)
}

func (st *SuffixTree) newInner(r Range) STNode {
	node := InnerNode{
		SharedNode: SharedNode{Range: r},
		Children:   make([]STNode, st.Alpha.Size())}
	return wrapInner(&node)
}

func (st *SuffixTree) breakEdge(n STNode, depth, leafidx int, y Range, x []byte) STNode {
	new_node := st.newInner(n.Shared().Range.prefix(depth))
	n.Shared().Parent.addChild(new_node, x)
	new_leaf := st.newLeaf(leafidx, y)
	n.Shared().From += depth
	new_node.Inner().addChild(new_leaf, x)
	new_node.Inner().addChild(n, x)
	return new_leaf
}

func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, `digraph { rankdir="LR" `)
	st.Root.ToDot(st.String, st.Alpha, w)
	fmt.Fprintln(w, "}")
}

// Search maps visitor through all the leaves in the subtree found by a search.
func (st *SuffixTree) Search(p_ string, visitor func(int)) {
	p, err := st.Alpha.MapToBytes(p_)
	if err != nil {
		// We can't map, so no hits
		return
	}
	n, depth, y := sscan(st.Root, Range{0, len(p)}, st.String, p)
	if depth == y.length() {
		n.LeafIndices(visitor)
	}
}

// -- Construction algorithms --------------------------

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

func lenSharedPrefix(r1, r2 Range, x, y []byte) int {
	i, n := 0, min(r1.length(), r2.length())
	for ; i < n; i++ {
		if x[r1.From+i] != y[r2.From+i] {
			break
		}
	}
	return i
}

// x is the underlying strings for nodes, y is the string
// for inter (which when we construct is also x, but when we
// search it is likely another string).
func sscan(n STNode, r Range, x, y []byte) (STNode, int, Range) {
	if r.length() == 0 {
		return n, 0, r
	}
	// If we scan on a node, it is an inner node.
	v := n.Inner().Children[y[r.From]]
	if v.isNil() {
		return n, 0, r
	}
	i := lenSharedPrefix(v.Shared().Range, r, x, y)
	if i == r.length() || i < v.Shared().Range.length() {
		return v, i, r
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, r.chump(i), x, y)
}

func NaiveST(x_ string) SuffixTree {
	x, alpha := MapStringWithSentinel(x_)

	st := SuffixTree{Alpha: alpha, String: x}
	st.Root = st.newInner(Range{0, 0})

	for i := 0; i < len(x); i++ {
		v, j, y := sscan(st.Root, Range{i, len(x)}, x, x)
		if j == 0 {
			// A mismatch when we try to leave a node
			// means that it is an inner node
			v.Inner().addChild(st.newLeaf(i, y), x)
		} else {
			st.breakEdge(v, j, i, y.chump(j), x)
		}
	}
	return st
}

func fscan(n STNode, r Range, x []byte) (STNode, int, Range) {
	if r.length() == 0 {
		return n, 0, r
	}
	// If we scan on a node, it is an inner node
	v := n.Inner().Children[x[r.From]]
	if v.isNil() {
		panic("With fscan there should always be an out-edge")
	}
	i := min(v.Shared().Range.length(), r.length())
	if i == r.length() {
		return v, i, r
	}
	// Continue from v (exploiting tail call optimisation)
	return fscan(v, r.chump(i), x)
}

func (v *SharedNode) suffix() Range {
	// If v's parent is the root, chop
	// off one index
	if v.Parent.Parent == nil {
		return v.chump(1)
	} else {
		return v.Range
	}
}

func McCreight(x_ string) SuffixTree {
	x, alpha := MapStringWithSentinel(x_)
	st := SuffixTree{Alpha: alpha, String: x}
	st.Root = st.newInner(Range{0, 0})
	st.Root.Inner().SuffixLink = st.Root.Inner()
	currLeaf := st.newLeaf(0, Range{0, len(x)})
	st.Root.Inner().addChild(currLeaf, x)

	// The bits of the suffix we need to search for
	var y, z Range
	// ynode is the node we get to when searching for y
	var ynode STNode
	// depth is how far down an edge we have searched
	var depth int

	for i := 1; i < len(x); i++ {
		p := currLeaf.Shared().Parent

		if p.SuffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.Shared().suffix()
			ynode = wrapInner(p.SuffixLink)

		} else {
			pp := p.Parent
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.Shared().Range

			ynode, depth, _ = fscan(wrapInner(pp.SuffixLink), y, x)
			if depth < ynode.Shared().Range.length() {
				// ended on an edge
				currLeaf = st.breakEdge(ynode, depth, i, z, x)
				p.SuffixLink = currLeaf.Shared().Parent
				continue // Go to next suffix, we are done here
			}

			// Remember p's suffix link for later...
			p.SuffixLink = ynode.Inner()
		}

		// This is the slow scan part, from ynode and the rest
		// of the suffix, which is z.
		n, depth, w := sscan(ynode, z, x, x)
		if depth == 0 {
			// Landed on a node
			currLeaf = st.newLeaf(i, w)
			n.Inner().addChild(currLeaf, x)

		} else {
			// Landed on an edge
			currLeaf = st.breakEdge(n, depth, i, w.chump(depth), x)
		}
	}
	return st
}

/*
func (st *SuffixTree) ComputeSuffixAndLcpArray() ([]int, []int) {
	sa := make([]int, len(st.String))
	lcp := make([]int, len(st.String))
	// FIXME
	return sa, lcp
}
*/
