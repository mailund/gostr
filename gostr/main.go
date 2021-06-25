package main

import (
	"os"

	"github.com/mailund/cli"
)

func main() {
	var main = cli.NewMenu(
		"gostr",
		"shows string algorithms",
		"Examples of various string algorithms and data structures.",
		show)

	main.Run(os.Args[1:])
}
