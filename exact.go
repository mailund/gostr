package gostr

// Help for the methods that cannot handle
// an empty pattern. It just returns every position.
func reportEmptyMatches(x string, fn func(int)) {
	for i := range x {
		fn(i)
	}

	// For empty matches, we also have one past the last
	// letter.
	fn(len(x))
}

// Naive runs the naive (duh) O(nm) times search algorithm.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func Naive(x, p string, callback func(int)) {
	var i, j int
	for i = 0; i < len(x)-len(p)+1; i++ {
		for j = 0; j < len(p); j++ {
			if x[i+j] != p[j] {
				break
			}
		}

		if j == len(p) {
			callback(i)
		}
	}
}

// BorderSearch runs the O(n+m) time algorithm based on building
// a border array and reporting when its value matches m.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func BorderSearch(x, p string, callback func(int)) {
	if p == "" {
		reportEmptyMatches(x, callback)
		return
	}

	ba := StrictBorderarray(p)
	b := 0

	for i := range x {
		for {
			if p[b] == x[i] {
				b++
				break
			}

			if b == 0 {
				break
			}

			b = ba[b-1]
		}

		if b == len(p) {
			callback(i - len(p) + 1)

			b = ba[b-1]
		}
	}
}

// Kmp runs the O(n+m) time Knuth-Morris-Prat algorithm.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func Kmp(x, p string, callback func(int)) {
	if p == "" {
		reportEmptyMatches(x, callback)
		return
	}

	ba := StrictBorderarray(p)

	var i, j int

	for i < len(x) {
		// Match...
		for i < len(x) && j < len(p) && x[i] == p[j] {
			i++
			j++
		}
		// Report...
		if j == len(p) {
			callback(i - len(p))
		}
		// Shift pattern...
		if j == 0 {
			i++
		} else {
			j = ba[j-1]
		}
	}
}

// Bmh runs the O(nm) worst-case but expected sub-linear time
// Boyer-Moore-Horspool algorithm. This version uses a table of size
// 256 to map bytes to jumps, exploiting that we get bytes out of
// indexing strings.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func Bmh(x, p string, callback func(int)) {
	if p == "" {
		reportEmptyMatches(x, callback)
		return
	}

	jump := make([]int, byteSize)
	for b := 0; b < len(jump); b++ {
		jump[b] = len(p)
	}

	for j := 0; j < len(p)-1; j++ {
		jump[p[j]] = len(p) - j - 1
	}

	for i := 0; i < len(x)-len(p)+1; i += jump[x[i+len(p)-1]] {
		for j := len(p) - 1; x[i+j] == p[j]; j-- {
			if j == 0 {
				callback(i)
				break
			}
		}
	}
}

// BmhWithMap runs the O(nm) worst-case but expected sub-linear time
// Boyer-Moore-Horspool algorithm. It uses a map for the jump table, and
// can theoretically handle more characters than bytes (but doesn't, since
// indexing into strings gives us bytes). It demonstrates the performance
// hit you get from using a map rather than an array as in Bmh.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func BmhWithMap(x, p string, callback func(int)) {
	if p == "" {
		reportEmptyMatches(x, callback)
		return
	}

	jumpTbl := map[byte]int{}
	for j := 0; j < len(p)-1; j++ {
		jumpTbl[p[j]] = len(p) - j - 1
	}

	for i := 0; i < len(x)-len(p)+1; {
		for j := len(p) - 1; x[i+j] == p[j]; j-- {
			if j == 0 {
				callback(i)
				break
			}
		}

		if jmp, ok := jumpTbl[x[i+len(p)-1]]; ok {
			i += jmp
		} else {
			i += len(p)
		}
	}
}

// BmhWithAlphabet runs the O(nm) worst-case but expected sub-linear time
// Boyer-Moore-Horspool algorithm. This version maps the input
// strings before search, so we know their alphabet size, and can
// create a jump table of the apprporiate size.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func BmhWithAlphabet(x, p string, callback func(int)) {
	if p == "" {
		reportEmptyMatches(x, callback)
		return
	}

	xb, alpha := MapString(x)

	pb, err := alpha.MapToBytes(p)
	if err != nil {
		// We can't map, so we can't match
		return
	}

	jumpTbl := make([]int, alpha.Size())

	for j := range jumpTbl {
		jumpTbl[j] = len(pb)
	}

	for j := 0; j < len(pb)-1; j++ {
		jumpTbl[pb[j]] = len(pb) - j - 1
	}

	for i := 0; i < len(xb)-len(pb)+1; {
		for j := len(pb) - 1; xb[i+j] == pb[j]; j-- {
			if j == 0 {
				callback(i)
				break
			}
		}

		i += jumpTbl[xb[i+len(pb)-1]]
	}
}
