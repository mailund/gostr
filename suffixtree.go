package gostr

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// -- Substrings represented as intervals -----------------
type interval struct {
	i, j int
}

func (r interval) length() int {
	return r.j - r.i
}

// Chump off the first k letters of an interval
func (r interval) chump(k int) interval {
	return interval{r.i + k, r.j}
}

func (r interval) prefix(k int) interval {
	return interval{r.i, r.i + k}
}

func (r interval) substr(x string) string {
	return x[r.i:r.j]
}

// STNode represents the nodes in a suffix tree.
type STNode interface {
	// This is an opague type, so the interface
	// is interly private. There are functions below
	// for the public interface

	// Methods where it makes sense to have dynamic dispatch
	// (And where we don't lose too much performance on it)
	toDot(x string, w io.Writer)
}

// Dynamic dispatch is slower than the switch statements, and I really
// only want the interface so I can store two types of nodes with a common
// interface, so I do it this way. It isn't as pretty (especially because
// they share functionality where I have repeated code), but it is faster

// -- Private interfce to STNodes ------
func castToShared(n STNode) *sharedNode {
	switch v := n.(type) {
	case *leafNode:
		return &v.sharedNode
	case *innerNode:
		return &v.sharedNode
	default:
		panic("Unknown STNode type")
	}
}

func getInterval(n STNode) interval {
	return castToShared(n).interval
}

func setParent(n, parent STNode) {
	castToShared(n).parent = parent
}

func chumpInterval(n STNode, depth int) {
	castToShared(n).i += depth
}

// -- Public interface to STNodes ------

func IsLeaf(n STNode) bool {
	switch n.(type) {
	case *leafNode:
		return true
	case *innerNode:
		return false
	default:
		panic("Unknown STNode type")
	}
}

func Parent(n STNode) STNode {
	return castToShared(n).parent
}

func Children(n STNode) []STNode {
	switch v := n.(type) {
	case *innerNode:
		return v.getChildren()
	default:
		panic("You can only get children from an inner node")
	}
}

func LeafIndex(n STNode) int {
	switch v := n.(type) {
	case *leafNode:
		return v.leafIdx
	default:
		panic("You can only get leaf indices from a leaf")
	}
}

func ToDot(n STNode, x string, w io.Writer) {
	n.toDot(x, w)
}

func EdgeLabel(n STNode, x string) string {
	return castToShared(n).substr(x)
}

func LeafIndices(n STNode, visitor func(int)) {
	switch v := n.(type) {
	case *leafNode:
		visitor(v.leafIdx)
	case *innerNode:
		for _, child := range v.getChildren() {
			LeafIndices(child, visitor)
		}
	default:
		panic("Unknown STNode type")
	}
}

// Data both in inner STNodes and in leaf-STNodes
type sharedNode struct {
	interval
	parent STNode
}

func (n *sharedNode) edgeLabel(x string) string {
	return n.interval.substr(x)
}

type innerNode struct {
	sharedNode
	suffixLink     STNode
	children       map[byte]STNode
	sortedChildren *[]STNode // Cached sorted edges for lexicographic output
}

func newInner(inter interval) *innerNode {
	return &innerNode{
		sharedNode: sharedNode{interval: inter},
		children:   map[byte]STNode{}}
}

func (n *innerNode) sortChildren() {
	edges := []byte{}
	for k := range n.children {
		edges = append(edges, k)
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i] < edges[j]
	})
	children := make([]STNode, len(edges))
	for i, e := range edges {
		children[i] = n.children[e]
	}
	n.sortedChildren = &children
}

func (n *innerNode) getChildren() []STNode {
	if n.sortedChildren == nil {
		n.sortChildren()
	}
	return *n.sortedChildren
}

func (n *innerNode) addChild(child STNode, x string) {
	if n.sortedChildren != nil {
		panic("The edges should never be sorted while we construct the tree.")
	}
	n.children[x[getInterval(child).i]] = child
	setParent(child, n)
}

func ReplaceSentinel(x string) string {
	// Need this one for readable output
	return strings.ReplaceAll(x, "\x00", "†")
}

func (n *innerNode) toDot(x string, w io.Writer) {
	if n.parent == nil {
		// Root
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", n)
	} else {
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			n.parent, n, ReplaceSentinel(n.edgeLabel(x)))
		fmt.Fprintf(w, "\"%p\"[shape=point]\n", n)
	}
	if n.suffixLink != nil {
		fmt.Fprintf(w, `"%p" -> "%p"[style=dotted, color=red];`, n, n.suffixLink)
	}
	for _, child := range n.children {
		child.toDot(x, w)
	}
}

type leafNode struct {
	sharedNode
	leafIdx int
}

func newLeaf(idx int, inter interval) *leafNode {
	return &leafNode{
		sharedNode: sharedNode{interval: inter},
		leafIdx:    idx}
}

func (n *leafNode) toDot(x string, w io.Writer) {
	fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
		n.parent, n, ReplaceSentinel(n.edgeLabel(x)))
	fmt.Fprintf(w, "\"%p\"[label=%d]\n", n, n.leafIdx)
}

// -- Suffix tree --------------------------

type SuffixTree struct {
	String string
	Root   STNode
}

func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, `digraph { rankdir="LR" `)
	ToDot(st.Root, st.String, w)
	fmt.Fprintln(w, "}")
}

// Search maps visitor through all the leaves in the subtree found by a search.
func (st *SuffixTree) Search(p string, visitor func(int)) {
	n, depth, y := sscan(st.Root, interval{0, len(p)}, st.String, p)
	if depth == y.length() {
		LeafIndices(n, visitor)
	}
}

func PathLabel(n STNode, x string) string {
	labels := []string{EdgeLabel(n, x)}
	for p := Parent(n); p != nil; p = Parent(p) {
		labels = append(labels, EdgeLabel(p, x))
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

func lenSharedPrefix(i1, i2 interval, x, y string) int {
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
func sscan(n STNode, inter interval, x, y string) (STNode, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	// If we scan on a node, it is an inner node.
	v, ok := n.(*innerNode).children[y[inter.i]]
	if !ok {
		return n, 0, inter
	}
	i := lenSharedPrefix(getInterval(v), inter, x, y)
	if i == inter.length() || i < getInterval(v).length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, inter.chump(i), x, y)
}

func breakEdge(n STNode, depth, leafidx int, y interval, x string) *leafNode {
	new_node := newInner(getInterval(n).prefix(depth))
	Parent(n).(*innerNode).addChild(new_node, x)
	new_leaf := newLeaf(leafidx, y)
	chumpInterval(n, depth)
	new_node.addChild(new_leaf, x)
	new_node.addChild(n, x)
	return new_leaf
}

func NaiveST(x string) SuffixTree {
	// Add sentinel
	x += "\x00"
	root := newInner(interval{0, 0})
	for i := 0; i < len(x); i++ {
		v, j, y := sscan(root, interval{i, len(x)}, x, x)
		if j == 0 {
			// A mismatch when we try to leave a node
			// means that it is an inner node
			v.(*innerNode).addChild(newLeaf(i, y), x)
		} else {
			breakEdge(v, j, i, y.chump(j), x)
		}
	}
	return SuffixTree{x, root}
}

func fscan(n STNode, inter interval, x string) (STNode, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	// If we scan on a node, it is an inner node
	v, ok := n.(*innerNode).children[x[inter.i]]
	if !ok {
		panic("With fscan there should always be an out-edge")
	}
	i := min(getInterval(v).length(), inter.length())
	if i == inter.length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return fscan(v, inter.chump(i), x)
}

func (v *sharedNode) suffix() interval {
	// If v's parent is the root, chop
	// off one index
	if Parent(v.parent) == nil {
		return v.chump(1)
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
	var ynode STNode
	// depth is how far down an edge we have searched
	var depth int

	for i := 1; i < len(x); i++ {
		p := currLeaf.parent.(*innerNode)

		if p.suffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.suffix()
			ynode = p.suffixLink

		} else {
			pp := p.parent.(*innerNode)
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.interval

			ynode, depth, _ = fscan(pp.suffixLink, y, x)
			if depth < getInterval(ynode).length() {
				// ended on an edge
				currLeaf = breakEdge(ynode, depth, i, z, x)
				p.suffixLink = currLeaf.parent
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
			currLeaf = newLeaf(i, w)
			n.(*innerNode).addChild(currLeaf, x)

		} else {
			// Landed on an edge
			currLeaf = breakEdge(n, depth, i, w.chump(depth), x)
		}
	}
	return SuffixTree{x, root}
}
