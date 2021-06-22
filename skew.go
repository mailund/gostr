/*
 Straightforward implementation of the skew/DC3 algorithm

- https://www.cs.helsinki.fi/u/tpkarkka/publications/jacm05-revised.pdf
*/

package gostr

func safeIdx(x []int32, i int32) int32 {
	if i >= int32(len(x)) {
		return 0
	}

	return x[i]
}

func sa3len(n int32) int32 {
	if n == 0 {
		return 0
	}

	return (n-1)/3 + 1
}

func sa12len(n int32) int32 {
	return n - sa3len(n)
}

func symbcount(x, idx []int32, offset, asize int32) []int32 {
	counts := make([]int32, asize)
	for _, i := range idx {
		counts[safeIdx(x, i+offset)]++
	}

	return counts
}

func cumsum(counts []int32) []int32 {
	var acc int32

	res := make([]int32, len(counts))
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
	idx = bucketSort(x, asize, idx, 2) //nolint:gomnd // 2 is an offset
	idx = bucketSort(x, asize, idx, 1)
	idx = bucketSort(x, asize, idx, 0)

	return idx
}

func getSA12(x []int32) []int32 {
	sa12 := make([]int32, sa12len(int32(len(x))))

	for i, j := 0, 0; i < len(x); i++ {
		if i%3 != 0 {
			sa12[j] = int32(i)
			j++
		}
	}

	return sa12
}

func getSA3(x, sa12 []int32) []int32 {
	k, sa3 := 0, make([]int32, sa3len(int32(len(x))))

	if len(x)%3 == 1 {
		sa3[k] = int32(len(x) - 1)
		k++
	}

	for _, i := range sa12 {
		if i%3 == 1 {
			sa3[k] = i - 1
			k++
		}
	}

	return sa3
}

type triplet = [3]int32
type tripletMap = map[triplet]int32

func trip(x []int32, i int32) triplet {
	return triplet{safeIdx(x, i), safeIdx(x, i+1), safeIdx(x, i+2)} //nolint:gomnd // 2 is an offset
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

func less(x []int32, i, j int32, isa map[int32]int32) bool {
	a, b := safeIdx(x, i), safeIdx(x, j)

	if a < b {
		return true
	}

	if a > b {
		return false
	}

	if i%3 != 0 && j%3 != 0 {
		return isa[i] < isa[j]
	}

	return less(x, i+1, j+1, isa)
}

func merge(x, sa12, sa3 []int32) []int32 {
	isa := map[int32]int32{}
	for i := 0; i < len(sa12); i++ {
		isa[sa12[i]] = int32(i)
	}

	sa := make([]int32, len(sa12)+len(sa3))

	i, j, k := 0, 0, 0
	for i < len(sa12) && j < len(sa3) {
		if less(x, sa12[i], sa3[j], isa) {
			sa[k] = sa12[i]
			i++
			k++
		} else {
			sa[k] = sa3[j]
			j++
			k++
		}
	}

	for ; i < len(sa12); i++ {
		sa[k] = sa12[i]
		k++
	}

	for ; j < len(sa3); j++ {
		sa[k] = sa3[j]
		k++
	}

	return sa
}

func buildU(x []int32, alpha tripletMap) []int32 {
	var (
		u = make([]int32, sa12len(int32(len(x))))
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
	}

	return 2 + 3*(i-m) //nolint:gomnd // 2 and 3 makes sense in the context of the algo.
}

func skew(x []int32, asize int32) []int32 {
	sa12 := radix3(x, asize, getSA12(x))
	alpha := collectAlphabet(x, sa12)

	if len(alpha) < len(sa12) {
		// Build u and its SA.
		u := buildU(x, alpha)
		usa := skew(u, int32(len(alpha)+1)) // +1 for sentinel

		// Then map back to SA12 indices
		m := int32((len(u) + 1) / 2) //nolint:gomnd // 2 is just half, not magic
		for k, i := range usa {
			sa12[k] = uidx(i, m)
		}
	}

	SA3 := bucketSort(x, asize, getSA3(x, sa12), 0)

	return merge(x, sa12, SA3)
}

// SkewWithAlphabet builds the suffix array of a string, first mapping
// it to a byte slice using the alphabet alpha.
func SkewWithAlphabet(x string, alpha *Alphabet) ([]int32, error) {
	xb, err := alpha.MapToIntsWithSentinel(x)
	if err != nil {
		return []int32{}, err
	}

	return skew(xb, int32(alpha.Size())), nil
}

// Skew builds the suffix array of a string using the skew algorithm.
func Skew(x string) []int32 {
	sa, _ := SkewWithAlphabet(x, NewAlphabet(x))
	return sa
}
