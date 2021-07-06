package gostr

import (
	"bytes"
	"encoding/gob"
)

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

// FMIndexTables contains the preprocessed tables used for FM-index
// searching
type FMIndexTables struct {
	Alpha *Alphabet
	Sa    []int32
	Ctab  *CTab
	Otab  *OTab

	// for approx matching
	Rotab *OTab
}

// BuildFMIndexExactTables builds the preprocessing tables for exact FM-index
// searching.
func BuildFMIndexExactTables(x string) *FMIndexTables {
	xb, alpha := MapStringWithSentinel(x)
	sa, _ := SaisWithAlphabet(x, alpha)
	bwt := Bwt(xb, sa)
	ctab := NewCTab(bwt, alpha.Size())
	otab := NewOTab(bwt, alpha.Size())

	return &FMIndexTables{
		Alpha: alpha,
		Sa:    sa,
		Ctab:  ctab,
		Otab:  otab,
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// BuildFMIndexExactTables builds the preprocessing tables for exact FM-index
// searching.
func BuildFMIndexApproxTables(x string) *FMIndexTables {
	tbls := BuildFMIndexExactTables(x)

	// Reverse string x and build the reverse O-table.
	revx := reverseString(x)
	sa, _ := SaisWithAlphabet(revx, tbls.Alpha)
	revb, _ := tbls.Alpha.MapToBytesWithSentinel(revx)
	tbls.Rotab = NewOTab(Bwt(revb, sa), tbls.Alpha.Size())

	return tbls
}

// FMIndexExactFromTables returns a search function based
// on the preprocessed tables
func FMIndexExactFromTables(tbls *FMIndexTables) func(p string, cb func(i int)) {
	return func(p string, cb func(i int)) {
		pb, err := tbls.Alpha.MapToBytes(p)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}

		left, right := 0, len(tbls.Sa)

		for i := len(pb) - 1; i >= 0; i-- {
			a := pb[i]
			left = tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, left)
			right = tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, right)

			if left >= right {
				return // no match
			}
		}

		for i := left; i < right; i++ {
			cb(int(tbls.Sa[i]))
		}
	}
}

// FMIndexExactPreprocess preprocesses the string x and returns a function
// that you can use to efficiently search in x.
func FMIndexExactPreprocess(x string) func(p string, cb func(i int)) {
	return FMIndexExactFromTables(BuildFMIndexExactTables(x))
}

func buildDtab(p []byte, tbls *FMIndexTables) []int {
	dtab := make([]int, len(p))

	minEdits := 0
	left, right := 0, len(tbls.Sa)

	for i := range p {
		a := p[i]
		left = tbls.Ctab.Rank(a) + tbls.Rotab.Rank(a, left)
		right = tbls.Ctab.Rank(a) + tbls.Rotab.Rank(a, right)

		if left >= right {
			minEdits++

			left, right = 0, len(tbls.Sa)
		}

		dtab[i] = minEdits
	}

	return dtab
}

func withOp(ops *EditOps, op ApproxEdit, fn func()) {
	*ops = append(*ops, op) // remember operation

	fn() // do computation

	*ops = (*ops)[:len(*ops)-1] // then pop operation
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

func fmApproxReport(left, right int, ops *EditOps, sa *[]int32, cb func(i int, cigar string)) {
	// Reverse the ops because we build them in reverse, then
	// convert them into a cigar for reporting
	cigar := OpsToCigar(revOps(ops))
	for j := left; j < right; j++ {
		cb(int((*sa)[j]), cigar)
	}
}

// FMIndexApproxFromTables return a search function from the preprocessed tables.
func FMIndexApproxFromTables(tbls *FMIndexTables) func(p string, edits int, cb func(i int, cigar string)) {
	return func(p string, edits int, cb func(i int, cigar string)) {
		pb, err := tbls.Alpha.MapToBytes(p)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}

		ops := make(EditOps, 0, len(p)+edits) // keep track of operations
		dtab := buildDtab(pb, tbls)           // D-table for early termination

		// closure for handling the real operations. There is a lot less
		// wrapping tables in structs or parsing them as parameters
		// with a closure, even if it might look a bit ugly.
		var rec func(i, left, right, edits int)

		rec = func(i, left, right, edits int) {
			if i < 0 {
				if edits >= 0 {
					fmApproxReport(left, right, &ops, &tbls.Sa, cb)
				}

				return
			}

			if edits < dtab[i] {
				return // not sufficient edits left
			}

			for a := byte(1); a < byte(tbls.Alpha.Size()); a++ {
				nextLeft := tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, left)
				nextRight := tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, right)

				if nextLeft == nextRight {
					continue
				}

				// Do an M operation
				withOp(&ops, M, func() {
					if a == pb[i] {
						rec(i-1, nextLeft, nextRight, edits)
					} else {
						rec(i-1, nextLeft, nextRight, edits-1)
					}
				})

				// Do a D operation, as long as it is not the first op
				if len(ops) > 0 {
					withOp(&ops, D, func() { rec(i, nextLeft, nextRight, edits-1) })
				}
			}

			// Do an I operation
			withOp(&ops, I, func() { rec(i-1, left, right, edits-1) })
		}

		// finally, fire away with the first recursive call!
		i, left, right := len(p)-1, 0, len(tbls.Sa)
		rec(i, left, right, edits)
	}
}

// FMIndexApproxPreprocess preprocesses the string x and returns a function
// that you can use to efficiently search in x.
func FMIndexApproxPreprocess(x string) func(p string, edits int, cb func(i int, cigar string)) {
	return FMIndexApproxFromTables(BuildFMIndexApproxTables(x))
}

// GobEncode implements the encoder interface for serialising to a stream of bytes
func (otab OTab) GobEncode() (res []byte, err error) {
	defer func() { err = catchError() }()

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
	)

	checkError(enc.Encode(otab.nrow))
	checkError(enc.Encode(otab.ncol))
	checkError(enc.Encode(otab.table))

	return buf.Bytes(), nil
}

// GobDecode implements the decoder interface for serialising to a stream of bytes
func (otab *OTab) GobDecode(b []byte) (err error) {
	defer func() { err = catchError() }()

	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)

	checkError(dec.Decode(&otab.nrow))
	checkError(dec.Decode(&otab.ncol))

	return dec.Decode(&otab.table)
}
