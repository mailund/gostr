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
	for {
		// Match up...
		for i < len(x) && j < len(p) && x[i] == p[j] {
			i++
			j++
		}
		// Report...
		if j == len(p) {
			callback(i - len(p))
		}
		if i == len(x) {
			break // We are done
		}
		// Shift pattern down...
		if j > 0 {
			for j = ba[j-1]; j > 0 && x[i] != p[j]; {
				j = ba[j-1]
			}
		}
		// And increment if we can't find a hit at
		// index zero after shifting.
		if j == 0 && x[i] != p[j] {
			i++
		}
	}
}
