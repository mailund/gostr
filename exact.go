package gostr

// Help for the methods that cannot handle
// an empty pattern. It just returns every position.
func reportEmptyMatches(x string, cb func(int)) {
	for i := range x {
		cb(i)
	}
	cb(len(x))
}

// Naive runs the naive (duh) O(nm) times search algorithm.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func Naive(x, p string, callback func(int)) {
	var i, j = 0, 0
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
	if len(p) == 0 {
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
	if len(p) == 0 {
		reportEmptyMatches(x, callback)
		return
	}
	ba := StrictBorderarray(p)
	var i int
	var j int
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
	if len(p) == 0 {
		reportEmptyMatches(x, callback)
		return
	}
	// There are 256 bytes, so that is what we use...
	jump := make([]int, 256)
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
	if len(p) == 0 {
		reportEmptyMatches(x, callback)
		return
	}

	jump_tbl := map[byte]int{}
	for j := 0; j < len(p)-1; j++ {
		jump_tbl[p[j]] = len(p) - j - 1
	}

	for i := 0; i < len(x)-len(p)+1; {
		for j := len(p) - 1; x[i+j] == p[j]; j-- {
			if j == 0 {
				callback(i)
				break
			}
		}

		if jmp, ok := jump_tbl[x[i+len(p)-1]]; ok {
			i += jmp
		} else {
			i += len(p)
		}
	}
}

// Bmh runs the O(nm) worst-case but expected sub-linear time
// Boyer-Moore-Horspool algorithm. This version maps the input
// strings before search, so we know their alphabet size, and can
// create a jump table of the apprporiate size.
//
// Parameters:
//   - x: the string we search in.
//   - p: the string we search for
//   - callback: a function called for each occurrence
func BmhWithAlphabet(x_, p_ string, callback func(int)) {
	if len(p_) == 0 {
		reportEmptyMatches(x_, callback)
		return
	}

	x, alpha := MapString(x_)
	p, err := alpha.MapToBytes(p_)
	if err != nil {
		// We can't map, so we can't match
		return
	}

	jump_tbl := make([]int, alpha.Size())
	for j := range jump_tbl {
		jump_tbl[j] = len(p)
	}
	for j := 0; j < len(p)-1; j++ {
		jump_tbl[p[j]] = len(p) - j - 1
	}

	for i := 0; i < len(x)-len(p)+1; {
		for j := len(p) - 1; x[i+j] == p[j]; j-- {
			if j == 0 {
				callback(i)
				break
			}
		}
		i += jump_tbl[x[i+len(p)-1]]
	}
}
