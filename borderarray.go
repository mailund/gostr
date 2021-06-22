package gostr

func Borderarray(x string) []int {
	ba := make([]int, len(x))
	for i := 1; i < len(x); i++ {
		b := ba[i-1]

		for {
			if x[b] == x[i] {
				ba[i] = b + 1
				break
			}

			if b == 0 {
				ba[i] = 0
				break
			}

			b = ba[b-1]
		}
	}

	return ba
}

func StrictBorderarray(x string) []int {
	ba := Borderarray(x)
	for i := 1; i < len(x)-1; i++ {
		if ba[i] > 0 && x[ba[i]] == x[i+1] {
			ba[i] = ba[ba[i]-1]
		}
	}

	return ba
}
