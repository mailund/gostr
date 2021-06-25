package gostr

import "testing"

func checkAlphabetMaps(t *testing.T, x string, letters []byte) {
	t.Helper()

	alpha := NewAlphabet(x)
	if alpha.Size() != len(letters)+1 { // +1 for sentinel
		t.Fatalf("Expected size %d, got %d", len(letters)+1, alpha.Size())
	}

	for i := 0; i < 256; i++ {
		a := byte(i)
		b := alpha._revmap[alpha._map[a]]

		if b != 0 && a != b {
			t.Fatalf("Mapping %d and then back gave us %d", a, b)
		}
	}

	for i := 0; i < 256; i++ {
		a := byte(i)
		b := alpha._map[alpha._revmap[a]]

		if b != 0 && a != b {
			t.Fatalf("Reverse mapping %d and then back gave us %d", a, b)
		}
	}
}

func Test_Alphabet_Maps(t *testing.T) {
	checkAlphabetMaps(t, "foo", []byte{'f', 'o'})
	checkAlphabetMaps(t, "foobar", []byte{'a', 'b', 'f', 'o', 'r'})

	// including the sentinel. Doesn't affect the alphabet, but
	// it will be included in the mapped string
	checkAlphabetMaps(t, "foobar\x00", []byte{'a', 'b', 'f', 'o', 'r'})
	checkAlphabetMaps(t, "foo\x00bar", []byte{'a', 'b', 'f', 'o', 'r'})
}
