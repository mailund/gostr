package gostr

func Borderarray(x string) []int {
	ba := make([]int, len(x))
	for i := 1; i < len(x); i++ {
		b := ba[i-1]
		for ; b > 0 && x[b] != x[i]; b = ba[b-1] {
		}
		if x[b] == x[i] {
			ba[i] = b + 1
		} else {
			ba[i] = 0
		}
	}
	return ba
}

func StrictBorderarray(x string) []int {
	ba := Borderarray(x)
	strict := make([]int, len(x))
	for i := 0; i < len(x); i++ {
		// I'm handling the last index inside the loop
		// so I don't have to deal with empty strings
		// outside of the loop.
		if i == len(x)-1 || x[ba[i]] != x[i+1] {
			strict[i] = ba[i]
		} else {
			strict[i] = strict[ba[i]]
		}
	}
	return strict
}
