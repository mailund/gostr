package gostr

import (
	"sort"
)

func remap(s string) (x []int, asize int) {
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

// Simple bit-array implementation...
var (
	mask = []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}
)

type bitArray struct {
	length int
	bytes  []byte
}

func newBitArray(size int, bits ...bool) *bitArray {
	ba := bitArray{length: size, bytes: make([]byte, (size+1)/8+1)}
	for i, b := range bits {
		ba.set(i, b)
	}
	return &ba
}

func (a *bitArray) get(i int) bool {
	return (a.bytes[i/8] & mask[i%8]) != 0
}

func (a *bitArray) set(i int, b bool) {
	if b {
		a.bytes[i/8] = a.bytes[i/8] | mask[i%8]
	} else {
		a.bytes[i/8] = a.bytes[(i)/8] & ^mask[(i)%8]
	}
}

func classifyS(isS *bitArray, x []int) {
	// Last element always exists, it is the sentinel and is S
	isS.set(len(x)-1, true)

	// Otherwise, an index is S if the first letter is smaller
	// or the first letters are the same and the next is S.
	for i := len(x) - 1; i > 0; i-- {
		isS.set(i-1, x[i-1] < x[i] || (x[i-1] == x[i] && isS.get(i)))
	}
}

func isLMS(isS *bitArray, i int) bool {
	if i == 0 {
		return false
	} else {
		return isS.get(i) && !isS.get(i-1)
	}
}

func equalLMS(x []int, isS *bitArray, i, j int) bool {
	if i == len(x) || j == len(x) {
		return false
	}
	for k := 0; ; k++ {
		iLMS := isLMS(isS, i+k)
		jLMS := isLMS(isS, j+k)
		if k > 0 && iLMS && jLMS {
			return true // reached end of the strings without diff.
		}
		if iLMS != jLMS || x[i+k] != x[j+k] {
			return false // mismatch
		}
	}
}

func countBuckets(x []int, asize int) []int {
	buckets := make([]int, asize)
	for _, a := range x {
		buckets[a]++
	}
	return buckets
}

func bucketsFronts(fronts, buckets []int) {
	sum := 0
	for i := range buckets {
		fronts[i] = sum
		sum += buckets[i]
	}
}

func bucketsEnd(ends, buckets []int) {
	sum := 0
	for i := range buckets {
		sum += buckets[i]
		ends[i] = sum
	}
}

func insertFrontBucket(out []int, fronts []int, bucket, val int) {
	out[fronts[bucket]] = val
	fronts[bucket]++
}

func insertEndBucket(out []int, ends []int, bucket, val int) {
	ends[bucket]--
	out[ends[bucket]] = val
}

const (
	undefined = -1
)

func clearToUndefined(SA []int) {
	for i := range SA {
		SA[i] = undefined
	}
}

func bucketLMS(
	x, SA []int,
	buckets, bucketEnds []int,
	isS *bitArray) {
	bucketsEnd(bucketEnds, buckets)
	for i := len(x) - 1; i >= 0; i-- {
		if isLMS(isS, i) {
			insertEndBucket(SA, bucketEnds, x[i], i)
		}
	}
}

func induceLS(x, SA, buckets, bucketEnds []int, isS *bitArray) {
	// Induce L sorting
	bucketsFronts(bucketEnds, buckets)
	for i := 0; i < len(x); i++ {
		if SA[i] == 0 || SA[i] == undefined {
			continue
		}
		j := SA[i] - 1
		if !isS.get(j) {
			insertFrontBucket(SA, bucketEnds, x[j], j)
		}
	}

	// Induce S sorting
	bucketsEnd(bucketEnds, buckets)
	for i := len(x) - 1; i >= 0; i-- {
		if SA[i] == 0 {
			continue
		}
		j := SA[i] - 1
		if isS.get(j) {
			insertEndBucket(SA, bucketEnds, x[j], j)
		}
	}
}

func compactLMS(SA []int, isS *bitArray) ([]int, []int) {
	k := 0
	for _, j := range SA {
		if isLMS(isS, j) {
			SA[k] = j
			k++
		}
	}
	// slice out the part with the LMS strings and the rest
	return SA[:k], SA[k:]
}

func compactDefined(x []int) []int {
	k := 0
	for _, i := range x {
		if i != undefined {
			x[k] = i
			k++
		}
	}
	// Slice out the piece we used
	return x[:k]
}

func reduceLMSString(x, SA []int, isS *bitArray) ([]int, []int, int) {
	// We split the input SA into two bits, one that is large
	// enough to hold the LMS indices and one that can hold the
	// indices if we divide them by two. The LMS strings are in the
	// first slice after the compaction, in sorted order. Using
	// compact and buffer, we can compute the reduced string.
	compact, buffer := compactLMS(SA, isS)

	clearToUndefined(buffer)
	prevLMS := compact[0]
	letter := 0 // the first LMS is the sentinel
	buffer[prevLMS/2] = 0
	for i := 1; i < len(compact); i++ {
		j := compact[i]
		if !equalLMS(x, isS, prevLMS, j) {
			letter++
		}
		buffer[j/2] = letter
		prevLMS = j
	}
	reduced := compactDefined(buffer)

	// The compact slice is big enough to store the SA for the
	// reduced string, so that is what we return for it.
	// The new alphabet size is the largest letter we have assigned
	// plus one (size == largest value + 1)
	return reduced, compact, letter + 1
}

func reverseLMSMap(
	x, SA []int,
	reducedSA, offsets []int,
	buckets, bucketEnds []int,
	isS *bitArray) {

	// Remap the reduced suffixes to the indices in the longer
	// string. Reset the rest of SA so its ready for imputing.

	// This figures out the original indices for the LMS strings.
	// They originally came in the same order, so index by index
	// they match, we just have to skip over the non-LMS indices
	k := 0
	for i := 1; i < len(x); i++ {
		if isLMS(isS, i) {
			offsets[k] = i
			k++
		}
	}

	// The replace the indices into the reduced string with
	// the indices into the original string
	for i, j := range reducedSA {
		SA[i] = offsets[j]
	}

	// Move the LMS suffixes to the correct buckets and leave
	// the rest of SA undefined. Going right-to-left we are
	// ensured that we cannot overwrite a LMS suffix we need to
	// move later.
	clearToUndefined(SA[len(reducedSA):])
	bucketsEnd(bucketEnds, buckets)
	var j int
	for i := len(reducedSA) - 1; i >= 0; i-- {
		j, reducedSA[i] = reducedSA[i], undefined
		insertEndBucket(SA, bucketEnds, x[j], j)
	}
}

func recSAIS(x, SA []int, asize int, isS *bitArray) {
	// Base case of recursion: unique characters
	if len(x) == asize {
		for i, a := range x {
			SA[a] = i
		}
		return
	}

	// Recursive case...
	classifyS(isS, x)
	buckets := countBuckets(x, asize)
	bucketEnds := make([]int, len(buckets))

	// Induce first sorting
	clearToUndefined(SA)
	bucketLMS(x, SA, buckets, bucketEnds, isS)
	induceLS(x, SA, buckets, bucketEnds, isS)

	// Recursion
	redX, redSA, redSize := reduceLMSString(x, SA, isS)
	if redSize == len(redX) {
		// Save some memory if we are going to recurse further
		buckets = nil
		bucketEnds = nil
	}
	recSAIS(redX, redSA, redSize, isS)
	classifyS(isS, x) // Recompute S/L types for this function
	if redSize == len(redX) {
		// Restore the tables we need again now
		buckets = countBuckets(x, asize)
		bucketEnds = make([]int, len(buckets))
	}

	// Second impute
	reverseLMSMap(x, SA, redSA, redX, buckets, bucketEnds, isS)
	induceLS(x, SA, buckets, bucketEnds, isS)
}

func SAIS(s string) (SA []int) {
	x, asize := remap(s)
	SA = make([]int, len(x))
	isS := newBitArray(len(x))
	recSAIS(x, SA, asize, isS)

	// slicing away the sentinel that we no longer need
	return SA[1:]
}
