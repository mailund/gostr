package main

import (
	"fmt"
	"os"

	"github.com/mailund/gostr"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "%s takes one argument, the string to build the suffix tree from.\n", os.Args[0])
		os.Exit(1)
	}
	st := gostr.McCreight(args[0])
	st.ToDot(os.Stdout)
}
