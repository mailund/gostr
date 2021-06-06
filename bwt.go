package gostr

import (
	"sort"
)

const sentinel = 0

/* In this code, the suffix array is one longer than the string
   because it includes the sentinel implicitly at index zero.
   We need the sentinel with BWT to jump correctly in the search,
   but we don't need to store it explicitly in the string.
*/

/*
   BwtIdx gives you the ith letter in the Burrows-Wheeler transform
   of a string. Use it if you want to know the BWT without making a
   new string.
*/
func BwtIdx(x string, sa []int, i int) byte {
	j := sa[i]
	if j == 0 {
		return sentinel
	} else {
		return x[j-1]
	}
}

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string.
func Bwt(x string, sa []int) string {
	bwt := make([]byte, len(sa))
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
	cumsum map[byte]int
	// Alpha The sorted bytes in the string.
	Alpha []byte
}

// Contains Does the c-table contain a byte?
func (ctab *CTAB) Contains(a byte) bool {
	_, ok := ctab.cumsum[a]
	return ok
}

// Rank How many times does the BWT hold a letter smaller
// than a?
func (ctab *CTAB) Rank(a byte) int {
	return ctab.cumsum[a]
}

// Ctab builds the c-table from a string.
func Ctab(x string) CTAB {
	// First, count how often we see each character
	counts := make(map[byte]int)
	for _, b := range x {
		counts[byte(b)]++
	}

	// Computer the sorted alphabet from this.
	// We could avoid sorting if we knew the alphabet
	// up front, but I won't bother with that here.
	alpha := make([]byte, len(counts))
	i := 0
	for k := range counts {
		alpha[i] = k
		i++
	}
	sort.Slice(alpha, func(i, j int) bool {
		return alpha[i] < alpha[j]
	})

	// Then do the cumulative sum
	cumsum := make(map[byte]int)
	acc := 1 // we start at 1 bcs of sentinel
	for _, a := range alpha {
		cumsum[a] = acc
		acc += counts[a]
	}

	return CTAB{cumsum, alpha}
}

// OTAB Holds the o-table (rank table) from a BWT string
type OTAB struct {
	table map[byte][]int
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTAB) Rank(a byte, i int) int {
	if i == 0 {
		return 0
	} else {
		return otab.table[a][i-1]
	}
}

// Otab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func Otab(x string, sa []int, ctab CTAB) OTAB {
	bwt := Bwt(x, sa)

	table := map[byte][]int{}
	for _, a := range ctab.Alpha {
		table[a] = make([]int, len(sa))
	}
	table[bwt[0]][0] = 1

	for _, a := range ctab.Alpha {
		for i := 1; i < len(sa); i++ {
			table[a][i] = table[a][i-1]
			if bwt[i] == a {
				table[a][i]++
			}
		}
	}
	return OTAB{table}
}

// BwtSearch finds all occurrences of p in x via a c-table
// and an o-table.
func BwtSearch(x, p string, ctab CTAB, otab OTAB) (int, int) {
	L, R := 0, len(x)+1 // + 1 to get the range including the sentinel
	for i := len(p) - 1; i >= 0; i-- {
		a := p[i]
		if !ctab.Contains(a) {
			return 0, 0
		}
		L = ctab.Rank(a) + otab.Rank(a, L)
		R = ctab.Rank(a) + otab.Rank(a, R)
		if L >= R {
			return 0, 0
		}
	}
	return L, R
}

func BwtPreprocess(x string) func(p string, cb func(i int)) {
	sa := Sais(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	return func(p string, cb func(i int)) {
		L, R := BwtSearch(x, p, ctab, otab)
		for i := L; i < R; i++ {
			cb(sa[i])
		}
	}
}
