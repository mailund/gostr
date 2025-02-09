package gostr

import (
	"fmt"
	"io"
	"strings"
	"unsafe"
)

// EdgeLabel is the representation of a string along an edge.
// It is a byte slice, so it takes up constant space, holding
// only a pointer, length and capacity, while the underlying
// string is shared between the edges in the tree.
type EdgeLabel []byte

// Revmap maps an EdgeLabel back to the underlying string. EdgeLabels
// are represented as slices of bytes, mapped from a string through an
// alphabet, and Revmap uses the alphabet to get the string back.
func (el EdgeLabel) Revmap(alpha *Alphabet) string {
	return alpha.RevmapBytes(el)
}

// STNodeType is a tag for identifying when we have leaves and when we
// have inner nodes.
type STNodeType int

const (
	// Leaf means that the node is a leaf
	Leaf STNodeType = iota
	// Inner means that the node is an inner node
	Inner STNodeType = iota
)

// SharedNode contains the attributes that both leaves and inner nodes have.
type SharedNode struct {
	EdgeLabel
	Parent *InnerNode
}

// LeafNode contains the additional properties that only leaves have.
type LeafNode struct {
	SharedNode
	Index int
}

// InnerNode contains the additional properties that only inner nodes have.
type InnerNode struct {
	SharedNode
	SuffixLink *InnerNode
	Children   []STNode
}

// STNode wraps either a leaf or an inner node. Use the node type determine which,
// before you access it as a node.
type STNode struct {
	NodeType STNodeType
	ptr      unsafe.Pointer
}

// IsNil returns true if the node represents a nil pointer. If it does, you cannot
// cast it to any other node type.
func (n STNode) IsNil() bool {
	return n.ptr == nil
}

// Shared returns a pointer to the shared part of leaves and inner nodes. It is an
// error to access a node if it IsNil(), but otherwise, you can get the shared part
// of both leaves and inner nodes.
func (n STNode) Shared() *SharedNode {
	return (*SharedNode)(n.ptr)
}

// Leaf casts a node to a leaf. You should only do this if the NodeType is Leaf.
func (n STNode) Leaf() *LeafNode {
	return (*LeafNode)(n.ptr)
}

// Inner casts a node to a leaf. You should only do this if the NodeType is Inner.
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

// PathLabel returns the string from the root down to a node.
// Do not call it with a NIL node, that is an error the function
// will crash.
func (n STNode) PathLabel(alpha *Alphabet) string {
	v := n.Shared()
	labels := []string{v.EdgeLabel.Revmap(alpha)}

	for p := v.Parent; p != nil; p = p.Parent {
		labels = append(labels, p.EdgeLabel.Revmap(alpha))
	}

	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}

	return strings.Join(labels, "")
}

// LeafIndices maps fn over all the leaf indices in the subtree
// rooted at n.
func (n STNode) LeafIndices(fn func(int)) {
	switch n.NodeType {
	case Leaf:
		fn(n.Leaf().Index)

	case Inner:
		for _, child := range n.Inner().Children {
			if !child.IsNil() {
				child.LeafIndices(fn)
			}
		}
	}
}

// ToDot writes the subtree starting at n to w.
//
// Parameters:
//   - alpha: The alphabet that was used to map the original string into
//     the byte representation stored in the tree. You can get it from the
//     suffix tree.
//   - w: the output stream to write the dot representation to.
func (n STNode) ToDot(alpha *Alphabet, w io.Writer) {
	switch n.NodeType {
	case Leaf:
		v := n.Leaf()
		fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
			v.Parent, v, v.Revmap(alpha))
		fmt.Fprintf(w, "\"%p\"[label=%d]\n", v, v.Index)

	case Inner:
		v := n.Inner()

		if v.Parent == nil {
			// Root
			fmt.Fprintf(w, "\"%p\"[label=\"\", shape=circle, style=filled, fillcolor=grey]\n", v)
		} else {
			fmt.Fprintf(w, "\"%p\" -> \"%p\"[label=\"%s\"]\n",
				v.Parent, v, v.Revmap(alpha))
			fmt.Fprintf(w, "\"%p\"[shape=point]\n", v)
		}

		if v.SuffixLink != nil {
			fmt.Fprintf(w, `"%p" -> "%p"[style=dotted, color=red];`, v, v.SuffixLink)
		}

		for _, child := range v.Children {
			if !child.IsNil() {
				child.ToDot(alpha, w)
			}
		}
	}
}

func (n *InnerNode) addChild(child STNode) {
	n.Children[child.Shared().EdgeLabel[0]] = child
	child.Shared().Parent = n
}

// -- Suffix tree --------------------------

// SuffixTree is the representation of a suffix tree.
type SuffixTree struct {
	Alpha  *Alphabet
	String []byte
	Root   STNode
}

func (st *SuffixTree) newLeaf(idx int, el EdgeLabel) STNode {
	return wrapLeaf(&LeafNode{
		SharedNode: SharedNode{EdgeLabel: el},
		Index:      idx})
}

func (st *SuffixTree) newInner(el EdgeLabel) STNode {
	return wrapInner(&InnerNode{
		SharedNode: SharedNode{EdgeLabel: el},
		Children:   make([]STNode, st.Alpha.Size())})
}

func (st *SuffixTree) breakEdge(n STNode, depth, leafidx int, y []byte) STNode {
	newNode := st.newInner(n.Shared().EdgeLabel[:depth])
	n.Shared().Parent.addChild(newNode)

	newLeaf := st.newLeaf(leafidx, y)
	n.Shared().EdgeLabel = n.Shared().EdgeLabel[depth:]
	newNode.Inner().addChild(newLeaf)
	newNode.Inner().addChild(n)

	return newLeaf
}

// ToDot writes a dot representation of the tree to the output writer w.
func (st *SuffixTree) ToDot(w io.Writer) {
	fmt.Fprintln(w, `digraph { rankdir="LR" `)
	st.Root.ToDot(st.Alpha, w)
	fmt.Fprintln(w, "}")
}

// Search maps visitor through all the leaves in the subtree found by a search.
func (st *SuffixTree) Search(p string, visitor func(int)) {
	pb, err := st.Alpha.MapToBytes(p)
	if err != nil {
		// We can't map, so no hits
		return
	}

	n, depth, y := sscan(st.Root, pb)
	if depth == len(y) {
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

func lenSharedPrefix(x, y []byte) int {
	i, n := 0, min(len(x), len(y))
	for ; i < n; i++ {
		if x[i] != y[i] {
			break
		}
	}

	return i
}

func sscan(n STNode, y []byte) (node STNode, depth int, search []byte) {
	for {
		if len(y) == 0 {
			return n, 0, y
		}

		// If we scan on a node, it is an inner node.
		v := n.Inner().Children[y[0]]
		if v.IsNil() {
			return n, 0, y
		}

		i := lenSharedPrefix(v.Shared().EdgeLabel, y)
		if i == len(y) || i < len(v.Shared().EdgeLabel) {
			return v, i, y
		}

		// Continue from v: sscan(v, y[i:])
		n, y = v, y[i:]
	}
}

// NaiveST is the naive O(n²) construction algorithm.
func NaiveST(x string) *SuffixTree {
	xb, alpha := MapStringWithSentinel(x)

	st := SuffixTree{Alpha: alpha, String: xb}
	st.Root = st.newInner(xb[0:0])

	for i := 0; i < len(xb); i++ {
		v, j, y := sscan(st.Root, xb[i:])
		if j == 0 {
			// A mismatch when we try to leave a node
			// means that it is an inner node
			v.Inner().addChild(st.newLeaf(i, y))
		} else {
			st.breakEdge(v, j, i, y[j:])
		}
	}

	return &st
}

func fscan(n STNode, y []byte) (node STNode, depth int, search []byte) {
	for {
		if len(y) == 0 {
			return n, 0, y
		}

		// If we scan on a node, it is an inner node
		v := n.Inner().Children[y[0]]

		i := min(len(v.Shared().EdgeLabel), len(y))
		if i == len(y) {
			return v, i, y
		}

		// Continue from v: fscan(v, y[i:])
		n, y = v, y[i:]
	}
}

func (v *SharedNode) suffix() []byte {
	// If v's parent is the root, chop
	// off one index
	if v.Parent.Parent == nil {
		return v.EdgeLabel[1:]
	}

	return v.EdgeLabel
}

// McCreight constructs a suffix tree using McCreight's algorithm.
func McCreight(x string) *SuffixTree {
	xb, alpha := MapStringWithSentinel(x)
	st := SuffixTree{Alpha: alpha, String: xb}
	st.Root = st.newInner(xb[0:0])
	st.Root.Inner().SuffixLink = st.Root.Inner()
	currLeaf := st.newLeaf(0, xb)
	st.Root.Inner().addChild(currLeaf)

	// The bits of the suffix we need to search for
	var y, z []byte
	// ynode is the node we get to when searching for y
	var ynode STNode
	// depth is how far down an edge we have searched
	var depth int

	for i := 1; i < len(xb); i++ {
		p := currLeaf.Shared().Parent

		if p.SuffixLink != nil {
			// We don't need y here, just z and ynode
			z = currLeaf.Shared().suffix()
			ynode = wrapInner(p.SuffixLink)
		} else {
			pp := p.Parent
			// this time we need to search in both y and z
			y = p.suffix()
			z = currLeaf.Shared().EdgeLabel

			ynode, depth, _ = fscan(wrapInner(pp.SuffixLink), y)
			if depth < len(ynode.Shared().EdgeLabel) {
				// ended on an edge
				currLeaf = st.breakEdge(ynode, depth, i, z)
				p.SuffixLink = currLeaf.Shared().Parent
				continue // Go to next suffix, we are done here
			}

			// Remember p's suffix link for later...
			p.SuffixLink = ynode.Inner()
		}

		// This is the slow scan part, from ynode and the rest
		// of the suffix, which is z.
		n, depth, w := sscan(ynode, z)
		if depth == 0 {
			// Landed on a node
			currLeaf = st.newLeaf(i, w)
			n.Inner().addChild(currLeaf)
		} else {
			// Landed on an edge
			currLeaf = st.breakEdge(n, depth, i, w[depth:])
		}
	}

	return &st
}

// SECTION Generating other arrays

// ComputeSuffixAndLcpArray constructs a suffix array and longest common prefix
// array from a suffix tree.
func (st *SuffixTree) ComputeSuffixAndLcpArray() (sa, lcp []int32) {
	sa = make([]int32, len(st.String))
	lcp = make([]int32, len(st.String))
	i := 0

	var traverse func(n STNode, left, depth int32)
	traverse = func(n STNode, left, depth int32) {
		switch n.NodeType {
		case Leaf:
			sa[i] = int32(n.Leaf().Index)
			lcp[i] = left
			i++

		case Inner:
			for _, child := range n.Inner().Children {
				if child.IsNil() {
					continue
				}

				traverse(child, left, depth+int32(len(child.Shared().EdgeLabel)))
				left = depth // The remaining children should use depth
			}
		}
	}

	traverse(st.Root, 0, 0)

	return sa, lcp
}

// StSaConstruction constructs a suffix array from a suffix tree.
func StSaConstruction(x string) []int32 {
	sa, _ := McCreight(x).ComputeSuffixAndLcpArray()
	return sa
}

// !SECTION
