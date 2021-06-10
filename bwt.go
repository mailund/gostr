package gostr

const sentinel = '\x00'

/*
   BwtIdx gives you the ith letter in the Burrows-Wheeler transform
   of a string. Use it if you want to know the BWT without making a
   new string.
*/
func BwtIdx(x []byte, sa []int, i int) byte {
	if j := sa[i]; j != 0 {
		return x[j-1]
	} else {
		return sentinel
	}
}

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string. The
// string should have a sentinel
func BwtString(x []byte, sa []int) string {
	bwt := make([]byte, len(x))
	for i := 0; i < len(sa); i++ {
		bwt[i] = BwtIdx(x, sa, i)
	}
	return string(bwt)
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
func NewCTab(x []byte, alpha *Alphabet) *CTab {
	// First, count how often we see each character
	counts := make([]int, alpha.Size())
	for _, b := range x {
		counts[b]++
	}
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

func (otab *OTab) offset(a, i int) int {
	// -1 to a because we don't store the sentinel
	// and -1 to i because we don't store the first
	// row (which is always zero)
	return otab.ncol*(a-1) + (i - 1)
}

func (otab *OTab) get(a, i int) int {
	return otab.table[otab.offset(a, i)]
}

func (otab *OTab) set(a, i int, val int) {
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
		return otab.get(int(a), i)
	}
}

// Otab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func NewOTab(x []byte, sa []int, ctab *CTab, alpha *Alphabet) *OTab {
	// We index for all characters except $, so
	// nrow is the alphabet size minus one.
	// We index all indices [0,len(sa)], but we emulate
	// row 0, since it is always zero, so we only need
	// len(sa) columns.
	nrow, ncol := alpha.Size()-1, len(sa)
	table := make([]int, nrow*ncol)
	otab := OTab{nrow, ncol, table}

	// The character at the beginning of bwt gets a count
	// of one at row one.
	a := int(BwtIdx(x, sa, 0))
	otab.set(a, 1, 1)

	// The remaining entries either copies or increment from
	// the previous column. We count a from 1 to alpha size
	// to skip the sentinel, then -1 for the index
	for a := 1; a < alpha.Size(); a++ {
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
	alpha := NewAlphabet(x_)
	x, _ := alpha.MapToBytesWithSentinel(x_)
	sa := SaisWithAlphabet(x_, alpha)
	ctab := NewCTab(x, alpha)
	otab := NewOTab(x, sa, ctab, alpha)
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
