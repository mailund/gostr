/*
 Straightforward implementation of the skew/DC3 algorithm

- https://www.cs.helsinki.fi/u/tpkarkka/publications/jacm05-revised.pdf
*/

package gostr

func safeIdx(x []int32, i int32) int32 {
	if i >= int32(len(x)) {
		return 0
	} else {
		return x[i]
	}
}

func sa3len(n int32) int32 {
	if n == 0 {
		return 0
	} else {
		return (n-1)/3 + 1
	}
}

func sa12len(n int32) int32 {
	return n - sa3len(n)
}

func symbcount(x []int32, idx []int32, offset int32, asize int32) []int32 {
	counts := make([]int32, asize)
	for _, i := range idx {
		counts[safeIdx(x, i+offset)]++
	}
	return counts
}

func cumsum(counts []int32) []int32 {
	res := make([]int32, len(counts))
	var acc int32
	for i, k := range counts {
		res[i] = acc
		acc += k
	}
	return res
}

func bucketSort(x []int32, asize int32, idx []int32, offset int32) []int32 {
	counts := symbcount(x, idx, offset, asize)
	buckets := cumsum(counts)
	out := make([]int32, len(idx))
	for _, i := range idx {
		bucket := safeIdx(x, i+offset)
		out[buckets[bucket]] = i
		buckets[bucket]++
	}
	return out
}

func radix3(x []int32, asize int32, idx []int32) []int32 {
	idx = bucketSort(x, asize, idx, 2)
	idx = bucketSort(x, asize, idx, 1)
	return bucketSort(x, asize, idx, 0)
}

func getSA12(x []int32) []int32 {
	SA12 := make([]int32, sa12len(int32(len(x))))
	for i, j := 0, 0; i < len(x); i++ {
		if i%3 != 0 {
			SA12[j] = int32(i)
			j++
		}
	}
	return SA12
}

func getSA3(x, SA12 []int32) []int32 {
	SA3 := make([]int32, sa3len(int32(len(x))))
	k := 0
	if len(x)%3 == 1 {
		SA3[k] = int32(len(x) - 1)
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

type triplet = [3]int32
type tripletMap = map[triplet]int32

func trip(x []int32, i int32) triplet {
	return triplet{safeIdx(x, i), safeIdx(x, i+1), safeIdx(x, i+2)}
}

func collectAlphabet(x, idx []int32) tripletMap {
	alpha := tripletMap{}
	for _, i := range idx {
		t := trip(x, i)
		if _, ok := alpha[t]; !ok {
			alpha[t] = int32(len(alpha) + 1) // + 1 for sentinel
		}
	}
	return alpha
}

func less(x []int32, i, j int32, ISA map[int32]int32) bool {
	a, b := safeIdx(x, i), safeIdx(x, j)
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

func merge(x, SA12, SA3 []int32) []int32 {
	ISA := map[int32]int32{}
	for i := 0; i < len(SA12); i++ {
		ISA[SA12[i]] = int32(i)
	}

	SA := make([]int32, len(SA12)+len(SA3))
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

func buildU(x []int32, alpha tripletMap) []int32 {
	u := make([]int32, sa12len(int32(len(x))))
	var (
		k int32
		i int32
	)
	for i = 1; i < int32(len(x)); i += 3 {
		u[k] = alpha[trip(x, i)]
		k++
	}
	for i = 2; i < int32(len(x)); i += 3 {
		u[k] = alpha[trip(x, i)]
		k++
	}
	return u
}

func uidx(i, m int32) int32 {
	if i < m {
		return 1 + 3*i
	} else {
		return 2 + 3*(i-m)
	}
}

func skew(x []int32, asize int32) []int32 {
	SA12 := radix3(x, asize, getSA12(x))
	alpha := collectAlphabet(x, SA12)
	if len(alpha) < len(SA12) {
		// Build u and its SA.
		u := buildU(x, alpha)
		usa := skew(u, int32(len(alpha)+1)) // +1 for sentinel
		// Then map back to SA12 indices
		var m int32 = int32((len(u) + 1) / 2)
		for k, i := range usa {
			SA12[k] = uidx(i, m)
		}
	}
	SA3 := bucketSort(x, asize, getSA3(x, SA12), 0)
	return merge(x, SA12, SA3)
}

// Skew builds the suffix array of a String using the skew algorithm.
func SkewWithAlphabet(x string, alpha *Alphabet) ([]int32, error) {
	x_, err := alpha.MapToIntsWithSentinel(x)
	if err != nil {
		return []int32{}, err
	}
	return skew(x_, int32(alpha.Size())), nil
}

func Skew(x string) []int32 {
	sa, _ := SkewWithAlphabet(x, NewAlphabet(x))
	return sa
}
