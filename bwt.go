package gostr

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string. The
// string should have a sentinel
func Bwt(x []byte, sa []int) []byte {
	bwt := make([]byte, len(x))
	for i := 0; i < len(sa); i++ {
		j := sa[i]
		if j == 0 {
			bwt[i] = 0
		} else {
			bwt[i] = x[j-1]
		}
	}
	return bwt
}

// CTab Structur holding the C-table for BWT search.
// This is a map from letters in the alphabet to the
// cumulative sum of how often we see letters in the
// BWT
type CTab struct {
	CumSum []int
}

// Rank How many times does the BWT hold a letter smaller
// than a? Undefined behaviour if a isn't in the table.
func (ctab *CTab) Rank(a byte) int {
	return ctab.CumSum[a]
}

// Ctab builds the c-table from a string.
func NewCTab(bwt []byte, asize int) *CTab {
	// First, count how often we see each character
	counts := make([]int, asize)
	for _, b := range bwt {
		counts[b]++
	}
	// Then get the accumulative sum
	var n int = 0
	for i, count := range counts {
		counts[i] = n
		n += count
	}
	return &CTab{counts}
}

// OTab Holds the o-table (rank table) from a BWT string
type OTab struct {
	nrow, ncol int
	table      []int
}

func (otab *OTab) offset(a byte, i int) int {
	// -1 to a because we don't store the sentinel
	// and -1 to i because we don't store the first
	// row (which is always zero)
	return otab.ncol*(int(a)-1) + (i - 1)
}

func (otab *OTab) get(a byte, i int) int {
	return otab.table[otab.offset(a, i)]
}

func (otab *OTab) set(a byte, i int, val int) {
	otab.table[otab.offset(a, i)] = val
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTab) Rank(a byte, i int) int {
	// We don't explicitly store the first column,
	// since it is always empty anyway.
	if i == 0 {
		return 0
	} else {
		return otab.get(a, i)
	}
}

// Otab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func NewOTab(bwt []byte, asize int) *OTab {
	// We index for all characters except $, so
	// nrow is the alphabet size minus one.
	// We index all indices [0,len(sa)], but we emulate
	// row 0, since it is always zero, so we only need
	// len(sa) columns.
	nrow, ncol := asize-1, len(bwt)
	table := make([]int, nrow*ncol)
	otab := OTab{nrow, ncol, table}

	// The character at the beginning of bwt gets a count
	// of one at row one.
	otab.set(bwt[0], 1, 1)

	// The remaining entries either copies or increment from
	// the previous column. We count a from 1 to alpha size
	// to skip the sentinel, then -1 for the index
	for a := 1; a < asize; a++ {
		ba := byte(a) // get the right type for accessing otab
		for i := 2; i <= len(bwt); i++ {
			val := otab.get(ba, i-1)
			if bwt[i-1] == ba {
				val++
			}
			otab.set(ba, i, val)
		}
	}
	return &otab
}

func countLetters(x []byte) int {
	observed := make([]int, 256)
	for _, a := range x {
		observed[a] = 1
	}
	asize := 0
	for _, i := range observed {
		asize += i
	}
	return asize
}

func ReverseBwt(bwt []byte) []byte {
	asize := countLetters(bwt)
	ctab := NewCTab(bwt, asize)
	otab := NewOTab(bwt, asize)

	var x []byte = make([]byte, len(bwt))
	var i int = 0
	// We start at len(bwt) - 2 because we already
	// (implicitly) have the sentinel at len(bwt) - 1
	// and this way we don't need to start at the index
	// in bwt that has the sentinel (so we save a search).
	for j := len(bwt) - 2; j >= 0; j-- {
		a := bwt[i]
		x[j] = a
		i = ctab.Rank(a) + otab.Rank(a, i)
	}

	return x
}

// BwtSearch finds all occurrences of p in x via a c-table
// and an o-table.
func BwtSearch(x, p []byte, ctab *CTab, otab *OTab) (int, int) {
	L, R := 0, len(x)
	for i := len(p) - 1; i >= 0; i-- {
		a := p[i]
		L = ctab.Rank(a) + otab.Rank(a, L)
		R = ctab.Rank(a) + otab.Rank(a, R)
		if L >= R {
			return 0, 0
		}
	}
	return L, R
}

func BwtPreprocess(x_ string) func(p string, cb func(i int)) {
	x, alpha := MapStringWithSentinel(x_)
	sa, _ := SaisWithAlphabet(x_, alpha)
	bwt := Bwt(x, sa)
	ctab := NewCTab(bwt, alpha.Size())
	otab := NewOTab(bwt, alpha.Size())
	return func(p_ string, cb func(i int)) {
		p, err := alpha.MapToBytes(p_)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}
		L, R := BwtSearch(x, p, ctab, otab)
		for i := L; i < R; i++ {
			cb(sa[i])
		}
	}
}
