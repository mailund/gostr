package gostr

func Naive(x, p string, callback func(i int)) {
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
