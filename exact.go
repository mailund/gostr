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
		for ; b > 0 && p[b] != x[i]; b = ba[b-1] {
		}

		if p[b] == x[i] {
			b++
		} else {
			b = 0
		}

		if b == len(p) {
			callback(i - len(p) + 1)
			b = ba[b-1]
		}
	}
}
