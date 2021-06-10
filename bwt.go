package gostr

const sentinel = '\x00'

/*
   BwtIdx gives you the ith letter in the Burrows-Wheeler transform
   of a string. Use it if you want to know the BWT without making a
   new string.
*/
func BwtIdx(x *String, sa []int, i int) byte {
	if j := sa[i]; j != 0 {
		return x.At(j - 1)
	} else {
		return sentinel
	}
}

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string.
func BwtString(x *String, sa []int) string {
	bwt := make([]byte, x.Length())
	for i := 0; i < len(sa); i++ {
		bwt[i] = BwtIdx(x, sa, i)
	}
	return string(bwt)
}

// CTAB Structor holding the C-table for BWT search.
// This is a map from letters in the alphabet to the
// cumulative sum of how often we see letters in the
// BWT
type CTAB struct {
	Alpha  *Alphabet
	CumSum []int
}

// Contains tells you whether the table contains
// the byte a. If a is the sentinel, this is always
// true
func (ctab *CTAB) Contains(a byte) bool {
	return ctab.Alpha.Contains(a)
}

// Rank How many times does the BWT hold a letter smaller
// than a? Undefined behaviour if a isn't in the table.
func (ctab *CTAB) Rank(a byte) int {
	return ctab.CumSum[a]
}

// Ctab builds the c-table from a string.
func Ctab(x *String) *CTAB {
	// First, count how often we see each character
	counts := make([]int, x.Alpha.Size())
	for _, b := range x.bytes {
		counts[b]++
	}
	var n int = 0
	for i, count := range counts {
		counts[i] = n
		n += count
	}
	return &CTAB{x.Alpha, counts}
}

// OTAB Holds the o-table (rank table) from a BWT string
type OTAB struct {
	nrow, ncol int
	table      []int
}

func (otab *OTAB) offset(a, i int) int {
	// -1 to a because we don't store the sentinel
	// and -1 to i because we don't store the first
	// row (which is always zero)
	return otab.ncol*(a-1) + (i - 1)
}

func (otab *OTAB) get(a, i int) int {
	return otab.table[otab.offset(a, i)]
}

func (otab *OTAB) set(a, i int, val int) {
	otab.table[otab.offset(a, i)] = val
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTAB) Rank(a byte, i int) int {
	// We don't explicitly store the first column,
	// since it is always empty anyway.
	if i == 0 {
		return 0
	} else {
		return otab.get(int(a), i)
	}
}

// Otab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func Otab(x *String, sa []int, ctab *CTAB) *OTAB {
	// We index for all characters except $, so
	// nrow is the alphabet size minus one.
	// We index all indices [0,len(sa)], but we emulate
	// row 0, since it is always zero, so we only need
	// len(sa) columns.
	nrow, ncol := x.Alpha.Size()-1, len(sa)
	table := make([]int, nrow*ncol)
	otab := OTAB{nrow, ncol, table}

	// The character at the beginning of bwt gets a count
	// of one at row one.
	a := int(BwtIdx(x, sa, 0))
	otab.set(a, 1, 1)

	// The remaining entries either copies or increment from
	// the previous column. We count a from 1 to alpha size
	// to skip the sentinel, then -1 for the index
	for a := 1; a < ctab.Alpha.Size(); a++ {
		for i := 2; i <= len(sa); i++ {
			val := otab.get(a, i-1)
			if BwtIdx(x, sa, i-1) == byte(a) {
				val++
			}
			otab.set(a, i, val)
		}
	}
	return &otab
}

// BwtSearch finds all occurrences of p in x via a c-table
// and an o-table.
func BwtSearch(x, p *String, ctab *CTAB, otab *OTAB) (int, int) {
	xlen := x.Length()     // include sentinel
	plen := p.Length() - 1 // exclude sentinel

	L, R := 0, xlen
	for i := plen - 1; i >= 0; i-- {
		a := p.At(i)
		L = ctab.Rank(a) + otab.Rank(a, L)
		R = ctab.Rank(a) + otab.Rank(a, R)
		if L >= R {
			return 0, 0
		}
	}
	return L, R
}

func BwtPreprocess(x *String) func(p string, cb func(i int)) {
	sa := Sais(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	return func(p_ string, cb func(i int)) {
		p, err := NewString(p_, x.Alpha)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}
		L, R := BwtSearch(x, p, ctab, otab)
		for i := L; i < R; i++ {
			cb(sa[i])
		}
	}
}

func BwtPreprocessGoString(x_ string) func(p string, cb func(i int)) {
	x, _ := NewString(x_, nil)
	return BwtPreprocess(x)
}
