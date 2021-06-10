package gostr

import "fmt"

type Alphabet struct {
	_map    [256]byte
	_revmap [256]byte
	size    int
}

func NewAlphabet(ref string) *Alphabet {
	alpha := Alphabet{}

	alpha._map[0] = 1 // sentinel is always here
	for i := 0; i < len(ref); i++ {
		alpha._map[ref[i]] = 1
	}

	var n byte = 0
	for a, tag := range alpha._map {
		if tag == 1 {
			alpha._map[a] = n
			alpha._revmap[n] = byte(a)
			n++
		}
	}
	alpha.size = int(n)

	return &alpha
}

func (alpha *Alphabet) Size() int {
	return alpha.size
}

func (alpha *Alphabet) Contains(a byte) bool {
	return a == 0 || alpha._map[a] != 0
}

// Generics would be really nice here... unfortunately, not in the
// language yet
func (alpha *Alphabet) mapBytes(x string, out []byte) ([]byte, error) {
	for i := 0; i < len(x); i++ {
		b := alpha._map[x[i]]
		if b == 0 && x[i] != 0 {
			return []byte{}, fmt.Errorf("character %q is not in the alphabet", x[i])
		}
		out[i] = b
	}
	return out, nil
}

func (alpha *Alphabet) mapInts(x string, out []int) ([]int, error) {
	for i := 0; i < len(x); i++ {
		b := alpha._map[x[i]]
		if b == 0 && x[i] != 0 {
			return []int{}, fmt.Errorf("character %q is not in the alphabet", x[i])
		}
		out[i] = int(b)
	}
	return out, nil
}

func (alpha *Alphabet) MapToBytes(x string) ([]byte, error) {
	return alpha.mapBytes(x, make([]byte, len(x)))
}

func (alpha *Alphabet) MapToBytesWithSentinel(x string) ([]byte, error) {
	return alpha.mapBytes(x, make([]byte, len(x)+1))
}

func (alpha *Alphabet) MapToInts(x string) ([]int, error) {
	return alpha.mapInts(x, make([]int, len(x)))
}

func (alpha *Alphabet) MapToIntsWithSentinel(x string) ([]int, error) {
	return alpha.mapInts(x, make([]int, len(x)+1))
}

func (alpha *Alphabet) revmapBytes(x []byte, strip_sentinel bool) string {
	strip := 0 // If we have a sentinel, we remove it again
	if strip_sentinel && len(x) > 0 && x[len(x)-1] == 0 {
		strip++
	}
	out := make([]byte, len(x)-strip)
	for i := 0; i < len(out); i++ {
		b := alpha._revmap[x[i]]
		out[i] = b
	}
	return string(out)
}

func (alpha *Alphabet) RevmapBytes(x []byte) string {
	return alpha.revmapBytes(x, false)
}

func (alpha *Alphabet) RevmapBytesStripSentinel(x []byte) string {
	return alpha.revmapBytes(x, true)
}
