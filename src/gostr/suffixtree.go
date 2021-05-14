package gostr

import (
	"fmt"
	"io"
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

// -- Suffix tree nodes --------------------
type Node interface {
	// Private interface
	getInterval() interval
	chumpInterval(int)
	setParent(parent Node)
	leafLabels(res *[]int)

	// Public interface
	Parent() Node
	EdgeLabel(x string) string
	ToDot(x string, w io.Writer)
	LeafLabels() []int
}

// helper function to implement the Node
// interface
func nodeLeafLabels(n Node) []int {
	res := []int{}
	n.leafLabels(&res)
	return res
}

// Data both in inner nodes and in leaf-nodes
type sharedNode struct {
	interval
	parent Node
}

func (n *sharedNode) Parent() Node {
	return n.parent
}

func (n *sharedNode) setParent(parent Node) {
	n.parent = parent
}

func (n *sharedNode) getInterval() interval {
	return n.interval
}

func (n *sharedNode) chumpInterval(i int) {
	n.i += i
}

func (n *sharedNode) EdgeLabel(x string) string {
	return n.interval.substr(x)
}

type innerNode struct {
	sharedNode
	suffixLink Node
	children   map[byte]Node
}

func newInner(inter interval) *innerNode {
	return &innerNode{
		sharedNode: sharedNode{interval: inter},
		children:   map[byte]Node{}}
}

func (n *innerNode) addChild(child Node, x string) {
	n.children[x[child.getInterval().i]] = child
	child.setParent(n)
}

func replaceSentinel(x string) string {
	// Need this one for dot output
	return strings.ReplaceAll(x, "\x00", "†")
}

func (n *innerNode) ToDot(x string, w io.Writer) {
	if n.parent == nil {
		// Root
		fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", n)
	} else {
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			n.Parent(), n, replaceSentinel(n.EdgeLabel(x)))
		fmt.Fprintf(w, "\"%p\"[shape=point]\n", n)
	}
	for _, child := range n.children {
		child.ToDot(x, w)
	}
}

func (n *innerNode) leafLabels(res *[]int) {
	for _, child := range n.children {
		child.leafLabels(res)
	}
}

func (n *innerNode) LeafLabels() []int {
	return nodeLeafLabels(n)
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

func (n *leafNode) ToDot(x string, w io.Writer) {
	fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
		n.parent, n, replaceSentinel(n.EdgeLabel(x)))
	fmt.Fprintf(w, "\"%p\"[label=%d]\n", n, n.leafIdx)
}

func (n *leafNode) leafLabels(res *[]int) {
	*res = append(*res, n.leafIdx)
}

func (n *leafNode) LeafLabels() []int {
	return nodeLeafLabels(n)
}

// -- Suffix tree --------------------------

type SuffixTree struct {
	String string
	Root   Node
}

func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, "digraph {")
	st.Root.ToDot(st.String, w)
	fmt.Fprintln(w, "}")
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

func lenSharedPrefix(i1, i2 interval, x string) int {
	i, n := 0, min(i1.length(), i2.length())
	for ; i < n; i++ {
		if x[i1.i+i] != x[i2.i+i] {
			break
		}
	}
	return i
}

func sscan(n Node, inter interval, x string) (Node, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	// If we scan on a node, it is an inner node.
	v, ok := n.(*innerNode).children[x[inter.i]]
	if !ok {
		return n, 0, inter
	}
	i := lenSharedPrefix(v.getInterval(), inter, x)
	if i == inter.length() || i < v.getInterval().length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, inter.chump(i), x)
}

func breakEdge(n Node, depth, leafidx int, y interval, x string) *leafNode {
	if n.Parent() == nil {
		panic("A node must have a parent when we break its edge.")
	}
	new_node := newInner(n.getInterval().prefix(depth))
	n.Parent().(*innerNode).addChild(new_node, x)
	new_leaf := newLeaf(leafidx, y)
	n.chumpInterval(depth)
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
			// A mismatch when we try to leave a node
			// means that it is an inner node
			v.(*innerNode).addChild(newLeaf(i, y), x)
		} else {
			breakEdge(v, j, i, y.chump(j), x)
		}
	}
	return SuffixTree{x, root}
}

func fscan(n Node, inter interval, x string) (Node, int, interval) {
	if inter.length() == 0 {
		return n, 0, inter
	}
	// If we scan on a node, it is an inner node
	v, ok := n.(*innerNode).children[x[inter.i]]
	if !ok {
		panic("With fscan there should always be an out-edge")
	}
	i := min(v.getInterval().length(), inter.length())
	if i == inter.length() {
		return v, i, inter
	}
	// Continue from v (exploiting tail call optimisation)
	return sscan(v, inter.chump(i), x)
}

func (v *sharedNode) suffix() interval {
	// If v's parent is the root, chop
	// off one index
	if v.parent.Parent() == nil {
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
	var ynode Node

	for i := 1; i < len(x); i++ {
		p := currLeaf.Parent().(*innerNode)

		if p.suffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.suffix()
			ynode = p.suffixLink
		} else {
			pp := p.parent.(*innerNode)
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.interval
			ynode, depth, _ := fscan(pp.suffixLink, y, x)
			if depth < ynode.getInterval().length() {
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
		n, depth, w := sscan(ynode, z, x)
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
