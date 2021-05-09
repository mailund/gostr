package gostr

import (
	"sort"
)

// If you want to know the BWT without making a new string.
func BwtIdx(x string, sa []int, i int) byte {
	j := sa[i]
	if j == 0 {
		return x[len(x)-1]
	} else {
		return x[j-1]
	}
}

// If you want BWT as a string
func Bwt(x string, sa []int) string {
	bwt := make([]byte, len(x))
	for i := 0; i < len(sa); i++ {
		bwt[i] = BwtIdx(x, sa, i)
	}
	return string(bwt)
}

type CTAB struct {
	cumsum map[byte]int
	// Alpha The sorted bytes in the string.
	Alpha []byte
}

func (ctab *CTAB) Contains(a byte) bool {
	_, ok := ctab.cumsum[a]
	return ok
}

func (ctab *CTAB) Rank(a byte) int {
	return ctab.cumsum[a]
}

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
	acc := 0
	for _, a := range alpha {
		cumsum[a] = acc
		acc += counts[a]
	}

	return CTAB{cumsum, alpha}
}

type OTAB struct {
	table map[byte][]int
}

func (otab *OTAB) Rank(a byte, i int) int {
	if i == 0 {
		return 0
	} else {
		return otab.table[a][i-1]
	}
}

func Otab(x string, sa []int, ctab CTAB) OTAB {
	bwt := Bwt(x, sa)

	table := map[byte][]int{}
	for _, a := range ctab.Alpha {
		table[a] = make([]int, len(x))
	}
	table[bwt[0]][0] = 1

	for _, a := range ctab.Alpha {
		for i := 1; i < len(x); i++ {
			table[a][i] = table[a][i-1]
			if bwt[i] == a {
				table[a][i]++
			}
		}
	}
	return OTAB{table}
}

func BwtSearch(x, p string, ctab CTAB, otab OTAB) (int, int) {
	L, R := 0, len(x)
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
