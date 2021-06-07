package gostr

// Help for the methods that cannot handle
// an empty pattern. It just returns every position.
func reportEmptyMatches(x string, cb func(int)) {
	for i := range x {
		cb(i)
	}
	cb(len(x))
}

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

func Kmp(x, p string, callback func(int)) {
	if len(p) == 0 {
		reportEmptyMatches(x, callback)
		return
	}
	ba := StrictBorderarray(p)
	var i int = 0
	var j int = 0
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
