package gostr

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string. The
// string should have a sentinel
func Bwt(x []byte, sa []int32) []byte {
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

// NewCTab builds the c-table from a string.
func NewCTab(bwt []byte, asize int) *CTab {
	// First, count how often we see each character
	counts := make([]int, asize)
	for _, b := range bwt {
		counts[b]++
	}
	// Then get the accumulative sum
	var n int
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

func (otab *OTab) set(a byte, i, val int) {
	otab.table[otab.offset(a, i)] = val
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTab) Rank(a byte, i int) int {
	// We don't explicitly store the first column,
	// since it is always empty anyway.
	if i == 0 {
		return 0
	}

	return otab.get(a, i)
}

// NewOTab builds the o-table from a string. It uses
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
	const noBytes = 256

	observed := make([]int, noBytes)
	for _, a := range x {
		observed[a] = 1
	}

	asize := 0
	for _, i := range observed {
		asize += i
	}

	return asize
}

// ReverseBwt reconstructs the original string from the bwt string.
func ReverseBwt(bwt []byte) []byte {
	asize := countLetters(bwt)
	ctab := NewCTab(bwt, asize)
	otab := NewOTab(bwt, asize)

	x := make([]byte, len(bwt))
	i := 0
	// We start at len(bwt) - 2 because we already
	// (implicitly) have the sentinel at len(bwt) - 1
	// and this way we don't need to start at the index
	// in bwt that has the sentinel (so we save a search).
	secondToLast := len(bwt) - 2 //nolint:gomnd // 2 isn't magic now...

	for j := secondToLast; j >= 0; j-- {
		a := bwt[i]
		x[j] = a
		i = ctab.Rank(a) + otab.Rank(a, i)
	}

	return x
}

// fmIndexTables contains the preprocessed tables used for FM-index
// searching
type fmIndexTables struct {
	alpha *Alphabet
	sa    []int32
	ctab  *CTab
	otab  *OTab

	// for approx matching
	rotab *OTab
}

// buildFMIndexExactTables builds the preprocessing tables for exact FM-index
// searching.
func buildFMIndexExactTables(x string) *fmIndexTables {
	xb, alpha := MapStringWithSentinel(x)
	sa, _ := SaisWithAlphabet(x, alpha)
	bwt := Bwt(xb, sa)
	ctab := NewCTab(bwt, alpha.Size())
	otab := NewOTab(bwt, alpha.Size())

	return &fmIndexTables{
		alpha: alpha,
		sa:    sa,
		ctab:  ctab,
		otab:  otab,
	}
}

// buildFMIndexExactTables builds the preprocessing tables for exact FM-index
// searching.
func buildFMIndexApproxTables(x string) *fmIndexTables {
	tbls := buildFMIndexExactTables(x)

	// Reverse string x and build the reverse O-table.
	revx := []byte(x)
	for i, j := 0, len(revx)-1; i < j; i, j = i+1, j-1 {
		revx[i], revx[j] = revx[j], revx[i]
	}

	sa, _ := SaisWithAlphabet(string(revx), tbls.alpha)
	revb, _ := tbls.alpha.MapToBytesWithSentinel(string(revx))
	tbls.rotab = NewOTab(Bwt(revb, sa), tbls.alpha.Size())

	return tbls
}

// FMIndexExactPreprocess preprocesses the string x and returns a function
// that you can use to efficiently search in x.
func FMIndexExactPreprocess(x string) func(p string, cb func(i int)) {
	tbls := buildFMIndexExactTables(x)

	return func(p string, cb func(i int)) {
		pb, err := tbls.alpha.MapToBytes(p)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}

		left, right := 0, len(tbls.sa)

		for i := len(pb) - 1; i >= 0; i-- {
			a := pb[i]
			left = tbls.ctab.Rank(a) + tbls.otab.Rank(a, left)
			right = tbls.ctab.Rank(a) + tbls.otab.Rank(a, right)

			if left >= right {
				return // no match
			}
		}

		for i := left; i < right; i++ {
			cb(int(tbls.sa[i]))
		}
	}
}

func buildDtab(p []byte, tbls *fmIndexTables) []int {
	dtab := make([]int, len(p))

	minEdits := 0
	left, right := 0, len(tbls.sa)

	for i := len(p) - 1; i >= 0; i-- {
		a := p[i]
		left = tbls.ctab.Rank(a) + tbls.otab.Rank(a, left)
		right = tbls.ctab.Rank(a) + tbls.otab.Rank(a, right)

		if left >= right {
			minEdits++

			left, right = 0, len(tbls.sa)
		}

		dtab[i] = minEdits
	}

	return dtab
}

// FMIndexApproxPreprocess preprocesses the string x and returns a function
// that you can use to efficiently search in x.
func FMIndexApproxPreprocess(x string) func(p string, edits int, cb func(i int, cigar string)) {
	tbls := buildFMIndexApproxTables(x)

	return func(p string, edits int, cb func(i int, cigar string)) {
		if p == "" {
			return // we can't handle empty strings...
		}

		pb, err := tbls.alpha.MapToBytes(p)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}

		// first step should not include D, so we explicitly doM and doI
		ops := make(EditOps, 0, len(p)+edits)
		dtab := buildDtab(pb, tbls)
		i, left, right := len(pb)-1, 0, len(tbls.sa)
		doM(tbls, pb, i, left, right, edits, dtab, &ops, cb)
		doI(tbls, pb, i, left, right, edits, dtab, &ops, cb)
	}
}

// we need to reverse to ops to make a cigar when we perform
// the search in reverse
func revOps(ops *EditOps) EditOps {
	rev := make(EditOps, len(*ops))

	copy(rev, *ops)

	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}

	return rev
}

func recApproxFMIndex(tbls *fmIndexTables,
	p []byte,
	i, left, right, edits int,
	dtab []int,
	ops *EditOps,
	fn func(int, string)) {
	if i < 0 {
		if edits >= 0 {
			cigar := OpsToCigar(revOps(ops))
			for j := left; j < right; j++ {
				fn(int(tbls.sa[j]), cigar)
			}
		}

		return
	}

	if edits < dtab[i] {
		return // not sufficient edits left
	}

	doM(tbls, p, i, left, right, edits, dtab, ops, fn)
	doI(tbls, p, i, left, right, edits, dtab, ops, fn)
	doD(tbls, p, i, left, right, edits, dtab, ops, fn)
}

func doM(tbls *fmIndexTables,
	p []byte,
	i, left, right, edits int,
	dtab []int,
	ops *EditOps,
	fn func(int, string)) {
	// record the M operation...
	(*ops) = append(*ops, M)

	for a := byte(1); a < byte(tbls.alpha.Size()); a++ {
		nextLeft := tbls.ctab.Rank(a) + tbls.otab.Rank(a, left)
		nextRight := tbls.ctab.Rank(a) + tbls.otab.Rank(a, right)

		if nextLeft >= nextRight {
			continue
		}

		nextEdits := edits
		if a != p[i] {
			nextEdits--
		}

		recApproxFMIndex(tbls, p, i-1, nextLeft, nextRight, nextEdits, dtab, ops, fn)
	}

	(*ops) = (*ops)[:len(*ops)-1]
}

func doI(tbls *fmIndexTables,
	p []byte,
	i, left, right, edits int,
	dtab []int,
	ops *EditOps,
	fn func(int, string)) {
	// record the I operation...
	(*ops) = append(*ops, I)
	recApproxFMIndex(tbls, p, i-1, left, right, edits-1, dtab, ops, fn)
	(*ops) = (*ops)[:len(*ops)-1]
}

func doD(tbls *fmIndexTables,
	p []byte,
	i, left, right, edits int,
	dtab []int,
	ops *EditOps,
	fn func(int, string)) {
	// record the D operation...
	(*ops) = append(*ops, D)

	for a := byte(1); a < byte(tbls.alpha.Size()); a++ {
		nextLeft := tbls.ctab.Rank(a) + tbls.otab.Rank(a, left)
		nextRight := tbls.ctab.Rank(a) + tbls.otab.Rank(a, right)

		if nextLeft >= nextRight {
			continue
		}

		recApproxFMIndex(tbls, p, i, nextLeft, nextRight, edits-1, dtab, ops, fn)
	}

	(*ops) = (*ops)[:len(*ops)-1]
}
