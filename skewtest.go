package main

import (
	"fmt"
	"gostr"
)

func main() {
	x := "mississippi"
	sa := gostr.Skew(x)
	for _, i := range sa {
		fmt.Printf("%d %s\n", i, x[i:])
	}
}
