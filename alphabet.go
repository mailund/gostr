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

func (alpha *Alphabet) Map(x string) ([]byte, error) {
	out := make([]byte, len(x)+1) // for convinience we always include sentinel
	for i := 0; i < len(x); i++ {
		b := alpha._map[x[i]]
		if b == 0 && x[i] != 0 {
			return []byte{}, fmt.Errorf("character %q is not in the alphabet", x[i])
		}
		out[i] = b
	}
	return out, nil
}

func (alpha *Alphabet) Revmap(x []byte) string {
	if len(x) == 0 || x[len(x)-1] != 0 {
		panic("Mapped strings must have a terminal sentinel")
	}
	out := make([]byte, len(x)-1)
	for i := 0; i < len(x)-1; i++ {
		b := alpha._revmap[x[i]]
		out[i] = b
	}
	return string(out)
}

type String struct {
	bytes []byte
	Alpha *Alphabet
}

func NewString(x string, alpha *Alphabet) (*String, error) {
	if alpha == nil {
		alpha = NewAlphabet(x)
	}
	if bytes, err := alpha.Map(x); err == nil {
		return &String{bytes, alpha}, nil
	} else {
		return nil, err
	}
}

func (x *String) At(i int) byte {
	return x.bytes[i]
}

func (x *String) Length() int {
	return len(x.bytes)
}

func (x *String) ToGoString() string {
	return string(x.Alpha.Revmap(x.bytes))
}

// ToInts returns the bytes in a string as an integer slice.
// This is mostly useful in suffix array construction algorithms where
// we need to work with slices of integers in the recursions.
func (x *String) ToInts() []int {
	res := make([]int, len(x.bytes))
	for i, b := range x.bytes {
		res[i] = int(b)
	}
	return res
}
