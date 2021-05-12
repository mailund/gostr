/*
 Straightforward implementation of the skew/DC3 algorithm

- https://www.cs.helsinki.fi/u/tpkarkka/publications/jacm05-revised.pdf
*/

package gostr

func safe_idx(x []int, i int) int {
	if i >= len(x) {
		return 0
	} else {
		return x[i]
	}
}

func symbcount(x []int, idx []int, offset int, asize int) []int {
	counts := make([]int, asize)
	for _, i := range idx {
		counts[safe_idx(x, i+offset)]++
	}
	return counts
}

func cumsum(counts []int) []int {
	res := make([]int, len(counts))
	acc := 0
	for i, k := range counts {
		res[i] = acc
		acc += k
	}
	return res
}

func bucketSort(x []int, asize int, idx []int, offset int) []int {
	counts := symbcount(x, idx, offset, asize)
	buckets := cumsum(counts)
	out := make([]int, len(idx))
	for _, i := range idx {
		bucket := safe_idx(x, i+offset)
		out[buckets[bucket]] = i
		buckets[bucket]++
	}
	return out
}

func radix3(x []int, asize int, idx []int) []int {
	idx = bucketSort(x, asize, idx, 2)
	idx = bucketSort(x, asize, idx, 1)
	return bucketSort(x, asize, idx, 0)
}

func getSA12(x []int) []int {
	/*
			You can append to the slice here:
		SA12 := []int{}
		for i := 0; i < len(x); i++ {
			if i%3 != 0 {
				SA12 = append(SA12, i)
			}
		}
			but that allocates extra space for the capacity
			so preallocating is better.
	*/
	SA12 := make([]int, len(x)-((len(x)-1)/3+1))
	for i, j := 0, 0; i < len(x); i++ {
		if i%3 != 0 {
			SA12[j] = i
			j++
		}
	}
	return SA12
}

func getSA3(x, SA12 []int) []int {
	/*You can append to the slice here:
	SA3 := []int{}
	if len(x)%3 == 1 {
		SA3 = append(SA3, len(x)-1)
	}
	for _, i := range SA12 {
		if i%3 == 1 {
			SA3 = append(SA3, i-1)
		}
	}
	but preallocating is better
	*/
	SA3 := make([]int, (len(x)-1)/3+1)
	k := 0
	if len(x)%3 == 1 {
		SA3[k] = len(x) - 1
		k++
	}
	for _, i := range SA12 {
		if i%3 == 1 {
			SA3[k] = i - 1
			k++
		}
	}
	return SA3
}

type triplet = [3]int
type tripletMap = map[triplet]int

func trip(x []int, i int) triplet {
	return triplet{safe_idx(x, i), safe_idx(x, i+1), safe_idx(x, i+2)}
}

func collectAlphabet(x []int, idx []int) tripletMap {
	alpha := tripletMap{}
	for _, i := range idx {
		t := trip(x, i)
		if _, ok := alpha[t]; !ok {
			alpha[t] = len(alpha) + 2
		}
	}
	return alpha
}

func less(x []int, i int, j int, ISA map[int]int) bool {
	a, b := safe_idx(x, i), safe_idx(x, j)
	if a < b {
		return true
	}
	if a > b {
		return false
	}
	if i%3 != 0 && j%3 != 0 {
		return ISA[i] < ISA[j]
	}
	return less(x, i+1, j+1, ISA)
}

func merge(x []int, SA12 []int, SA3 []int) []int {
	ISA := map[int]int{}
	for i := 0; i < len(SA12); i++ {
		ISA[SA12[i]] = i
	}
	/*
			Using append:
		SA := []int{}
		i, j := 0, 0
		for i < len(SA12) && j < len(SA3) {
			if less(x, SA12[i], SA3[j], ISA) {
				SA = append(SA, SA12[i])
				i++
			} else {
				SA = append(SA, SA3[j])
				j++
			}
		}
		SA = append(SA, SA12[i:]...)
		SA = append(SA, SA3[j:]...)
	*/

	SA := make([]int, len(SA12)+len(SA3))
	i, j, k := 0, 0, 0
	for i < len(SA12) && j < len(SA3) {
		if less(x, SA12[i], SA3[j], ISA) {
			SA[k] = SA12[i]
			i++
			k++
		} else {
			SA[k] = SA3[j]
			j++
			k++
		}
	}
	for ; i < len(SA12); i++ {
		SA[k] = SA12[i]
		k++
	}
	for ; j < len(SA3); j++ {
		SA[k] = SA3[j]
		k++
	}
	return SA
}

func buildU(x []int, alpha tripletMap) []int {
	/*
		With append:
			u := []int{}
			for i := 1; i < len(x); i += 3 {
				u = append(u, alpha[trip(x, i)])
			}
			u = append(u, 1)
			for i := 2; i < len(x); i += 3 {
				u = append(u, alpha[trip(x, i)])
			}
	*/
	// The length is len(SA12) which is len(x)-(len(x)/3 + 1)
	// but then plus the central sentinel, so don't subtract the
	// 1.
	u := make([]int, len(x)-len(x)/3)
	k := 0
	for i := 1; i < len(x); i += 3 {
		u[k] = alpha[trip(x, i)]
		k++
	}
	u[k] = 1
	k++
	for i := 2; i < len(x); i += 3 {
		u[k] = alpha[trip(x, i)]
		k++
	}

	return u
}

func uidx(i int, m int) int {
	if i < m {
		return 1 + 3*i
	} else {
		return 2 + 3*(i-m-1)
	}
}

func skew(x []int, asize int) []int {
	// Some of the length calculations assume that x isn't
	// empty, so handle that explicitly. If we append to slices
	// instead of pre-allocating, we don't need it, but it
	// will still save some time to handle this case as soon
	// as possible.
	if len(x) == 0 {
		return []int{}
	}

	SA12 := radix3(x, asize, getSA12(x))
	alpha := collectAlphabet(x, SA12)
	if len(alpha) < len(SA12) {
		// Build u and its SA.
		u := buildU(x, alpha)
		usa := skew(u, len(alpha)+2) // +2 for sentinels
		// Then map back to SA12 indices
		m := len(usa) / 2
		k := 0
		for _, i := range usa {
			if i != m {
				SA12[k] = uidx(i, m)
				k++
			}
		}
	}
	SA3 := bucketSort(x, asize, getSA3(x, SA12), 0)
	return merge(x, SA12, SA3)
}

func str2int(x string) []int {
	out := make([]int, len(x))
	for i, c := range x {
		out[i] = int(c)
	}
	return out
}

// Skew builds the suffix array of a string using the skew algorithm.
func Skew(x string) []int {
	/*
		Skew algorithm for a string."
		The skew() function wants a list of integers,
		so we convert the string in the first call.
		I am assuming that the alphabet size is 256 here, although
		of course it might not be. It is a simplification instead of
		remapping the string.
	*/
	return skew(str2int(x), 256)
}
