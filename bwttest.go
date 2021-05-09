package main

import (
	"fmt"
	"gostr"
)

func main() {
	x := "mis"
	sa := gostr.Skew(x)
	fmt.Println(sa)
	for _, i := range sa {
		fmt.Println(x[i:] + x[:i])
	}
	fmt.Println()

	ctab := gostr.Ctab(x)
	otab := gostr.Otab(x, sa, ctab)
	L, R := gostr.BwtSearch(x, "mis", ctab, otab)
	fmt.Println("Hits:", L, R)
	for i := L; i < R; i++ {
		fmt.Println(sa[i], x[sa[i]:])
	}
}
