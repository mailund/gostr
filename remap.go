package gostr

import "sort"

func Remap(s string) (x []int, asize int) {
	// Identify the letters in the string
	alpha := map[byte]int{}
	for i := 0; i < len(s); i++ {
		alpha[s[i]] = 1
	}

	// then sort those letters and assign them numbers
	// respecting the original order
	letters := make([]byte, len(alpha))
	var i int = 0
	for a := range alpha {
		letters[i] = a
		i++
	}
	sort.Slice(letters, func(i, j int) bool {
		return letters[i] < letters[j]
	})
	for i, a := range letters {
		alpha[a] = i + 1 // +1 to reserve zero for $
	}

	// Finally, output the string with the new alphabet
	res := make([]int, len(s)+1) // +1 for $
	for i := 0; i < len(s); i++ {
		res[i] = alpha[s[i]]
	}
	// The last index is already zero (sentinel) from default
	// int values.

	return res, len(alpha) + 1
}
