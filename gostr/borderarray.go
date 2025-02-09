package gostr

// Borderarray computes the border array over the string x. The border
// array ba will at index i have the length of the longest proper border
// of the string x[:i+1], i.e. the longest non-empty string that is both
// a prefix and a suffix of x[:i+1].
func Borderarray(x string) []int {
	ba := make([]int, len(x))
	for i := 1; i < len(x); i++ {
		b := ba[i-1]

		for {
			if x[b] == x[i] {
				ba[i] = b + 1
				break
			}

			if b == 0 {
				ba[i] = 0
				break
			}

			b = ba[b-1]
		}
	}

	return ba
}

// StrictBorderarray computes the strict border array over the string x.
// This is almost the same as the border array, but ba[i] will be the
// longest proper border of the string x[:i+1] such that x[ba[i]] != x[i].
func StrictBorderarray(x string) []int {
	ba := Borderarray(x)
	for i := 1; i < len(x)-1; i++ {
		if ba[i] > 0 && x[ba[i]] == x[i+1] {
			ba[i] = ba[ba[i]-1]
		}
	}

	return ba
}
