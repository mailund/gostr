package gostr

import (
	"fmt"
	"io"
	"strings"
)

// -- Substrings represented as intervals -----------------
type Range struct {
	i, j int
}

func (r Range) length() int {
	return r.j - r.i
}

// Chump off the first k letters of an interval
func (r Range) chump(k int) Range {
	return Range{r.i + k, r.j}
}

func (r Range) prefix(k int) Range {
	return Range{r.i, r.i + k}
}

func (r Range) substr(x []byte, alpha *Alphabet) string {
	return alpha.RevmapBytes(x[r.i:r.j])
}

// STNode represents the nodes in a suffix tree.
type STNode interface {
	// This is an opague type, so the interface
	// is interly private. There are functions below
	// for the public interface

	// Methods where it makes sense to have dynamic dispatch
	// (And where we don't lose too much performance on it)
	toDot(x []byte, alpha *Alphabet, w io.Writer)
}

// Dynamic dispatch is slower than the switch statements, and I really
// only want the interface so I can store two types of nodes with a common
// interface, so I do it this way. It isn't as pretty (especially because
// they share functionality where I have repeated code), but it is faster

// -- Private interfce to STNodes ------
func castToShared(n STNode) *SharedNode {
	switch v := n.(type) {
	case *LeafNode:
		return &v.SharedNode
	case *InnerNode:
		return &v.SharedNode
	default:
		panic("Unknown STNode type")
	}
}

func getInterval(n STNode) Range {
	return castToShared(n).Range
}

func setParent(n, parent STNode) {
	castToShared(n).Parent = parent
}

func chumpInterval(n STNode, depth int) {
	castToShared(n).Range.i += depth
}

// -- Public interface to STNodes ------

func IsLeaf(n STNode) bool {
	switch n.(type) {
	case *LeafNode:
		return true
	case *InnerNode:
		return false
	default:
		panic("Unknown STNode type")
	}
}

func Parent(n STNode) STNode {
	return castToShared(n).Parent
}

func Children(n STNode) []STNode {
	switch v := n.(type) {
	case *InnerNode:
		return v.children
	default:
		panic("You can only get children from an inner node")
	}
}

func LeafIndex(n STNode) int {
	switch v := n.(type) {
	case *LeafNode:
		return v.leafIdx
	default:
		panic("You can only get leaf indices from a leaf")
	}
}

func EdgeLabel(n STNode, x []byte, alpha *Alphabet) string {
	return castToShared(n).substr(x, alpha)
}

func LeafIndices(n STNode, visitor func(int)) {
	switch v := n.(type) {
	case *LeafNode:
		visitor(v.leafIdx)
	case *InnerNode:
		for _, child := range v.children {
			if child != nil {
				LeafIndices(child, visitor)
			}
		}
	default:
		panic("Unknown STNode type")
	}
}

// Data both in inner STNodes and in leaf-STNodes
type SharedNode struct {
	Range
	Parent STNode
}

func (n *SharedNode) edgeLabel(x []byte, alpha *Alphabet) string {
	return n.Range.substr(x, alpha)
}

type InnerNode struct {
	SharedNode
	suffixLink STNode
	children   []STNode
}

func (n *InnerNode) addChild(child STNode, x []byte) {
	n.children[x[getInterval(child).i]] = child
	setParent(child, n)
}

func ReplaceSentinelBytes(x []byte, alpha *Alphabet) string {
	// Need this one for readable output
	return ReplaceSentinelString(string(alpha.RevmapBytes(x)))
}

func ReplaceSentinelString(x string) string {
	// Need this one for readable output
	return strings.ReplaceAll(x, "\x00", "†")
}

func (n *InnerNode) toDot(x []byte, alpha *Alphabet, w io.Writer) {
	if n.Parent == nil {
		// Root
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", n)
	} else {
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			n.Parent, n, ReplaceSentinelString(n.edgeLabel(x, alpha)))
		fmt.Fprintf(w, "\"%p\"[shape=point]\n", n)
	}
	if n.suffixLink != nil {
		fmt.Fprintf(w, `"%p" -> "%p"[style=dotted, color=red];`, n, n.suffixLink)
	}
	for _, child := range n.children {
		if child != nil {
			child.toDot(x, alpha, w)
		}
	}
}

type LeafNode struct {
	SharedNode
	leafIdx int
}

func (n *LeafNode) toDot(x []byte, alpha *Alphabet, w io.Writer) {
	fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
		n.Parent, n, ReplaceSentinelString(n.edgeLabel(x, alpha)))
	fmt.Fprintf(w, "\"%p\"[label=%d]\n", n, n.leafIdx)
}

// -- Suffix tree --------------------------

type SuffixTree struct {
	Alpha  *Alphabet
	String []byte
	Root   STNode
}

func (st *SuffixTree) newLeaf(idx int, r Range) *LeafNode {
	if r.i >= len(st.String) {
		panic("penis")
	}
	if r.j > len(st.String) {
		panic("narko")
	}
	return &LeafNode{
		SharedNode: SharedNode{Range: r},
		leafIdx:    idx}
}

func (st *SuffixTree) newInner(r Range) *InnerNode {
	return &InnerNode{
		SharedNode: SharedNode{Range: r},
		children:   make([]STNode, st.Alpha.Size())}
}

func (st *SuffixTree) breakEdge(n STNode, depth, leafidx int, y Range, x []byte) *LeafNode {
	if y.length() == 0 {
		panic("how the hell did this happen?")
	}
	new_node := st.newInner(getInterval(n).prefix(depth))
	Parent(n).(*InnerNode).addChild(new_node, x)
	new_leaf := st.newLeaf(leafidx, y)
	chumpInterval(n, depth)
	new_node.addChild(new_leaf, x)
	new_node.addChild(n, x)
	return new_leaf
}

func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, `digraph { rankdir="LR" `)
	st.Root.toDot(st.String, st.Alpha, w)
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
		LeafIndices(n, visitor)
	}
}

func PathLabel(n STNode, x []byte, alpha *Alphabet) string {
	labels := []string{EdgeLabel(n, x, alpha)}
	for p := Parent(n); p != nil; p = Parent(p) {
		labels = append(labels, EdgeLabel(p, x, alpha))
	}
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}
	return strings.Join(labels, "")
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

func lenSharedPrefix(i1, i2 Range, x, y []byte) int {
	i, n := 0, min(i1.length(), i2.length())
	for ; i < n; i++ {
		if x[i1.i+i] != y[i2.i+i] {
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
	v := n.(*InnerNode).children[y[r.i]]
	if v == nil {
		return n, 0, r
	}
	i := lenSharedPrefix(getInterval(v), r, x, y)
	if i == r.length() || i < getInterval(v).length() {
		return v, i, r
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, r.chump(i), x, y)
}

func NaiveST(x_ string) SuffixTree {
	alpha := NewAlphabet(x_)
	x, _ := alpha.MapToBytesWithSentinel(x_)

	st := SuffixTree{alpha, x, nil}
	st.Root = st.newInner(Range{0, 0})

	for i := 0; i < len(x); i++ {
		v, j, y := sscan(st.Root, Range{i, len(x)}, x, x)
		if j == 0 {
			// A mismatch when we try to leave a node
			// means that it is an inner node
			v.(*InnerNode).addChild(st.newLeaf(i, y), x)
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
	v := n.(*InnerNode).children[x[r.i]]
	if v == nil {
		panic("With fscan there should always be an out-edge")
	}
	i := min(getInterval(v).length(), r.length())
	if i == r.length() {
		return v, i, r
	}
	// Continue from v (exploiting tail call optimisation)
	return fscan(v, r.chump(i), x)
}

func (v *SharedNode) suffix() Range {
	// If v's parent is the root, chop
	// off one index
	if Parent(v.Parent) == nil {
		return v.chump(1)
	} else {
		return v.Range
	}
}

func McCreight(x_ string) SuffixTree {
	alpha := NewAlphabet(x_)
	x, _ := alpha.MapToBytesWithSentinel(x_)
	st := SuffixTree{alpha, x, nil}
	st.Root = st.newInner(Range{0, 0})
	st.Root.(*InnerNode).suffixLink = st.Root
	currLeaf := st.newLeaf(0, Range{0, len(x)})
	st.Root.(*InnerNode).addChild(currLeaf, x)

	// The bits of the suffix we need to search for
	var y, z Range
	// ynode is the node we get to when searching for y
	var ynode STNode
	// depth is how far down an edge we have searched
	var depth int

	for i := 1; i < len(x); i++ {
		p := currLeaf.Parent.(*InnerNode)

		if p.suffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.suffix()
			ynode = p.suffixLink

		} else {
			pp := p.Parent.(*InnerNode)
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.Range

			ynode, depth, _ = fscan(pp.suffixLink, y, x)
			if depth < getInterval(ynode).length() {
				// ended on an edge
				currLeaf = st.breakEdge(ynode, depth, i, z, x)
				p.suffixLink = currLeaf.Parent
				continue // Go to next suffix, we are done here
			}

			// Remember p's suffix link for later...
			p.suffixLink = ynode
		}

		// This is the slow scan part, from ynode and the rest
		// of the suffix, which is z.
		n, depth, w := sscan(ynode, z, x, x)
		if depth == 0 {
			// Landed on a node
			currLeaf = st.newLeaf(i, w)
			n.(*InnerNode).addChild(currLeaf, x)

		} else {
			// Landed on an edge
			currLeaf = st.breakEdge(n, depth, i, w.chump(depth), x)
		}
	}
	return st
}

/*
func (st *SuffixTree) computeSuffixAndLcpArray() ([]int, []int) {
	sa := make([]int, len(st.String))
	lcp := make([]int, len(st.String))
	// FIXME
	return sa, lcp
}
*/
