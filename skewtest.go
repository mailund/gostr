package main

import (
	"fmt"
	"skew"
)

func main() {
	x := "mississippi"
	sa := skew.Skew(x)
	for _, i := range sa {
		fmt.Printf("%d %s\n", i, x[i:])
	}
}
