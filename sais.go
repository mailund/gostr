package gostr

const bytebits = 8

type bitArray struct {
	length int
	bytes  []byte
}

func newBitArray(size int, bits ...bool) *bitArray {
	ba := bitArray{length: size, bytes: make([]byte, (size+bytebits-1)/bytebits)}

	for i, b := range bits {
		ba.set(int32(i), b)
	}

	return &ba
}

func (a *bitArray) get(i int32) bool {
	return (a.bytes[i/bytebits] & (1 << (i % bytebits))) != 0
}

func (a *bitArray) set(i int32, b bool) {
	if b {
		a.bytes[i/bytebits] |= 1 << (i % bytebits)
	} else {
		a.bytes[i/bytebits] = a.bytes[(i)/bytebits] & ^(1 << (i % bytebits))
	}
}

func classifyS(isS *bitArray, x []int32) {
	// Last element always exists, it is the sentinel and is S
	isS.set(int32(len(x)-1), true)

	// Otherwise, an index is S if the first letter is smaller
	// or the first letters are the same and the next is S.
	var secondToLastIndex = int32(len(x) - 2) //nolint:gomnd // -2 because we want second last

	for i := secondToLastIndex; i >= 0; i-- {
		isS.set(i, x[i] < x[i+1] || (x[i] == x[i+1] && isS.get(i+1)))
	}
}

func isLMS(isS *bitArray, i int32) bool {
	return (i != 0) && isS.get(i) && !isS.get(i-1)
}

func equalLMS(x []int32, isS *bitArray, i, j int32) bool {
	if i == j {
		// The same index is obviously the same...
		return true
	}
	// they can't be equal now, so only one is the
	// sentinel LMS, thus they cannot be equal
	if i == int32(len(x)) || j == int32(len(x)) {
		return false
	}

	// From here on, we assume that neither index points past the end.
	for k := int32(0); ; k++ {
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

func countBuckets(x []int32, asize int) []int32 {
	buckets := make([]int32, asize)
	for _, a := range x {
		buckets[a]++
	}

	return buckets
}

func bucketsFronts(fronts, buckets []int32) {
	var sum int32
	for i := range buckets {
		fronts[i] = sum
		sum += buckets[i]
	}
}

func bucketsEnd(ends, buckets []int32) {
	var sum int32
	for i := range buckets {
		sum += buckets[i]
		ends[i] = sum
	}
}

func insertBucketFront(out, fronts []int32, bucket, val int32) {
	out[fronts[bucket]] = val
	fronts[bucket]++
}

func insertBucketEnd(out, ends []int32, bucket, val int32) {
	ends[bucket]--
	out[ends[bucket]] = val //nolint:wsl // There is nothing wrong here!
}

const (
	undefined = -1
)

func clearToUndefined(sa []int32) {
	for i := range sa {
		sa[i] = undefined
	}
}

func bucketLMS(x, sa, buckets, bucketEnds []int32, isS *bitArray) {
	bucketsEnd(bucketEnds, buckets)

	for i := int32(len(x) - 1); i >= 0; i-- {
		if isLMS(isS, i) {
			insertBucketEnd(sa, bucketEnds, x[i], i)
		}
	}
}

func induceLS(x, sa, buckets, bucketEnds []int32, isS *bitArray) {
	// Induce L sorting
	bucketsFronts(bucketEnds, buckets)

	for i := 0; i < len(x); i++ {
		if sa[i] == 0 || sa[i] == undefined {
			continue
		}

		j := sa[i] - 1
		if !isS.get(j) {
			insertBucketFront(sa, bucketEnds, x[j], j)
		}
	}

	// Induce S sorting
	bucketsEnd(bucketEnds, buckets)

	for i := len(x) - 1; i >= 0; i-- {
		if sa[i] == 0 {
			continue
		}

		j := sa[i] - 1
		if isS.get(j) {
			insertBucketEnd(sa, bucketEnds, x[j], j)
		}
	}
}

func compactLMS(sa []int32, isS *bitArray) (compact, rest []int32) {
	k := 0

	for _, j := range sa {
		if isLMS(isS, j) {
			sa[k] = j
			k++
		}
	}

	// slice out the part with the LMS strings and the rest
	return sa[:k], sa[k:]
}

func compactDefined(x []int32) []int32 {
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

func reduceLMSString(x, sa []int32, isS *bitArray) (reduced, compact []int32, asize int) {
	// We split the input SA into two bits, one that is large
	// enough to hold the LMS indices and one that can hold the
	// indices if we divide them by two. The LMS strings are in the
	// first slice after the compaction, in sorted order. Using
	// compact and buffer, we can compute the reduced string.
	compact, buffer := compactLMS(sa, isS)

	clearToUndefined(buffer)

	var letter int32

	prevLMS := compact[0]
	buffer[prevLMS/2] = 0

	for i := 1; i < len(compact); i++ {
		j := compact[i]
		if !equalLMS(x, isS, prevLMS, j) {
			letter++
		}

		buffer[j/2] = letter
		prevLMS = j
	}

	reduced = compactDefined(buffer)

	// The compact slice is big enough to store the SA for the
	// reduced string, so that is what we return for it.
	// The new alphabet size is the largest letter we have assigned
	// plus one (size == largest value + 1)
	return reduced, compact, int(letter + 1)
}

func reverseLMSMap(x, sa, reducedSA, offsets, buckets, bucketEnds []int32, isS *bitArray) {
	// Remap the reduced suffixes to the indices in the longer
	// string. Reset the rest of SA so its ready for imputing.
	// This figures out the original indices for the LMS strings.
	// They originally came in the same order, so index by index
	// they match, we just have to skip over the non-LMS indices
	var k, i int32
	for i = 1; i < int32(len(x)); i++ {
		if isLMS(isS, i) {
			offsets[k] = i
			k++
		}
	}

	// The replace the indices into the reduced string with
	// the indices into the original string
	for i, j := range reducedSA {
		sa[i] = offsets[j]
	}

	// Move the LMS suffixes to the correct buckets and leave
	// the rest of SA undefined. Going right-to-left we are
	// ensured that we cannot overwrite a LMS suffix we need to
	// move later.
	clearToUndefined(sa[len(reducedSA):])
	bucketsEnd(bucketEnds, buckets)

	var j int32

	for i := len(reducedSA) - 1; i >= 0; i-- {
		j, reducedSA[i] = reducedSA[i], undefined
		insertBucketEnd(sa, bucketEnds, x[j], j)
	}
}

func recSais(x, sa []int32, asize int, isS *bitArray) {
	// Base case of recursion: unique characters
	if len(x) == asize {
		for i, a := range x {
			sa[a] = int32(i)
		}

		return
	}

	// Recursive case...
	classifyS(isS, x)
	buckets := countBuckets(x, asize)
	bucketEnds := make([]int32, len(buckets))

	// Induce first sorting
	clearToUndefined(sa)
	bucketLMS(x, sa, buckets, bucketEnds, isS)
	induceLS(x, sa, buckets, bucketEnds, isS)

	// Recursion
	redX, redSA, redSize := reduceLMSString(x, sa, isS)
	if redSize != len(redX) {
		// Save some memory if we are going to recurse further
		buckets = nil
		bucketEnds = nil
	}

	recSais(redX, redSA, redSize, isS)
	classifyS(isS, x) // Recompute S/L types for this function

	if redSize != len(redX) {
		// Restore the tables we need again now
		buckets = countBuckets(x, asize)
		bucketEnds = make([]int32, len(buckets))
	}

	// Second impute
	reverseLMSMap(x, sa, redSA, redX, buckets, bucketEnds, isS)
	induceLS(x, sa, buckets, bucketEnds, isS)
}

// SaisWithAlphabet builds a suffix array from the string x, first
// mapping it to bytes using the alphabet alpha.
func SaisWithAlphabet(x string, alpha *Alphabet) ([]int32, error) {
	xb, err := alpha.MapToIntsWithSentinel(x)
	if err != nil {
		return []int32{}, err
	}

	sa := make([]int32, len(xb))
	isS := newBitArray(len(xb))

	recSais(xb, sa, alpha.Size(), isS)

	return sa, nil
}

// Sais builds a suffix array from the string x
func Sais(x string) (sa []int32) {
	sa, _ = SaisWithAlphabet(x, NewAlphabet(x))
	return sa
}
