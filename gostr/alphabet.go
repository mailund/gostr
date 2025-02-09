package gostr

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const (
	// Sentinel is a unique byte (zero) used in several algorithms
	// to be different and smaller than all other letters.
	Sentinel = '\x00'
	// SentinelSymbol is a rune you can use to display the sentinel
	// when you translate byte slices into strings.
	SentinelSymbol = rune('𝕊')
)

// Alphabet handles mapping from strings to smaller alphabets of bytes.
type Alphabet struct {
	_map    [256]byte
	_revmap [256]byte
	size    int
}

// NewAlphabet creates an alphabet consisting of the bytes in ref only.
func NewAlphabet(ref string) *Alphabet {
	alpha := Alphabet{
		_map:    [256]byte{},
		_revmap: [256]byte{},
		size:    0,
	}

	alpha._map[0] = 1 // sentinel is always here
	for i := 0; i < len(ref); i++ {
		alpha._map[ref[i]] = 1
	}

	var alphaSize byte // tracks alphabet size...

	for a, tag := range alpha._map {
		if tag == 1 {
			alpha._map[a] = alphaSize
			alpha._revmap[alphaSize] = byte(a)
			alphaSize++
		}
	}

	alpha.size = int(alphaSize)

	return &alpha
}

// Size gives the number of letters in the alphabet. It will always be at least
// one, since all alphabets contain the sentinel character (zero).
func (alpha *Alphabet) Size() int {
	return alpha.size
}

// Contains checks if a is contained in the alphabet
func (alpha *Alphabet) Contains(a byte) bool {
	return a == Sentinel || alpha._map[a] != 0
}

// Generics would be really nice here... unfortunately, not in the
// language yet
func (alpha *Alphabet) mapBytes(x string, out []byte) ([]byte, error) {
	for i := 0; i < len(x); i++ {
		b := alpha._map[x[i]]
		if b == 0 && x[i] != 0 {
			return []byte{}, &AlphabetLookupError{x[i]}
		}

		out[i] = b
	}

	return out, nil
}

func (alpha *Alphabet) mapInts(x string, out []int32) ([]int32, error) {
	for i := 0; i < len(x); i++ {
		b := alpha._map[x[i]]
		if b == 0 && x[i] != 0 {
			return []int32{}, &AlphabetLookupError{b}
		}

		out[i] = int32(b)
	}

	return out, nil
}

// MapToBytes translates a string into a byte slice, mapping characters according
// to the alphabet
func (alpha *Alphabet) MapToBytes(x string) ([]byte, error) {
	return alpha.mapBytes(x, make([]byte, len(x)))
}

// MapToBytesWithSentinel translates a string into a byte slice, mapping characters according
// to the alphabet. The resulting byte slice has a terminal zero, acting as a sentinel.
func (alpha *Alphabet) MapToBytesWithSentinel(x string) ([]byte, error) {
	return alpha.mapBytes(x, make([]byte, len(x)+1))
}

// MapToInts translates a string into an integer slice, mapping characters according
// to the alphabet. It only differs from MapToBytes in the type of the output, but is
// used in algorithms where we need to operate on integers rather than the smaller
// bytes.
func (alpha *Alphabet) MapToInts(x string) ([]int32, error) {
	return alpha.mapInts(x, make([]int32, len(x)))
}

// MapToIntsWithSentinel translates a string into an integer slice, mapping characters according
// to the alphabet. The resulting int slice has a terminal zero, acting as a sentinel.
// It only differs from MapToBytes in the type of the output, but is
// used in algorithms where we need to operate on integers rather than the smaller
// bytes.
func (alpha *Alphabet) MapToIntsWithSentinel(x string) ([]int32, error) {
	return alpha.mapInts(x, make([]int32, len(x)+1))
}

func (alpha *Alphabet) revmapBytes(x []byte, stripSentinel bool) string {
	strip := 0 // If we have a sentinel, we remove it again
	if stripSentinel && len(x) > 0 && x[len(x)-1] == 0 {
		strip++
	}

	out := make([]rune, len(x)-strip)

	for i := 0; i < len(out); i++ {
		if x[i] == Sentinel {
			out[i] = SentinelSymbol
		} else {
			out[i] = rune(alpha._revmap[x[i]])
		}
	}

	return string(out)
}

// RevmapBytes maps a byte slice back into a string according to the alphabet
// that was used to translate the string into bytes in the first place.
func (alpha *Alphabet) RevmapBytes(x []byte) string {
	return alpha.revmapBytes(x, false)
}

// RevmapBytesStripSentinel maps a byte slice back into a string according to the alphabet
// that was used to translate the string into bytes in the first place. RevmapBytesStripSentinel
// will strip the last character in the input from the output, getting rid of a sentinel that
// (hopefully) was added when the byte slice was created.
func (alpha *Alphabet) RevmapBytesStripSentinel(x []byte) string {
	return alpha.revmapBytes(x, true)
}

// MapString creates an alphabet from the input and then maps the string through it,
// returning both resulting byte slice and alphabet.
func MapString(x string) ([]byte, *Alphabet) {
	alpha := NewAlphabet(x)
	xb, _ := alpha.MapToBytes(x)

	return xb, alpha
}

// MapStringWithSentinel creates an alphabet from the input and then maps the string through it,
// returning both resulting byte slice and alphabet. Unlike MapString, MapStringWithSentinel
// will add a terminal zero byte to the byte slice it returns.
func MapStringWithSentinel(x string) ([]byte, *Alphabet) {
	alpha := NewAlphabet(x)
	xb, _ := alpha.MapToBytesWithSentinel(x)

	return xb, alpha
}

// GobEncode implements the encoder interface for serialising an alphabet to a stream of bytes
func (alpha Alphabet) GobEncode() (res []byte, err error) { 
	defer catchError(&err)

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
	)

	checkError(enc.Encode(alpha._map))
	checkError(enc.Encode(alpha._revmap))
	checkError(enc.Encode(alpha.size))

	res = buf.Bytes()

	return res, nil
}

// GobDecode implements the encoder interface for serialising an alphabet to a stream of bytes
func (alpha *Alphabet) GobDecode(b []byte) (error) {
	
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)

	checkError(dec.Decode(&alpha._map))
	checkError(dec.Decode(&alpha._revmap))

	if err := dec.Decode(&alpha.size); err != nil {
		return fmt.Errorf("failed to decode alphabet size: %w", err)
	} 

	return nil
}
