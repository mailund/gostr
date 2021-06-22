package main

import (
	"os"

	"github.com/mailund/cli"
	"github.com/mailund/gostr"
)

func ShowSuffixTreeCommand() *cli.Command {
	init := func(cmd *cli.Command) func() {
		var x string

		cmd.Params.StringVar(&x, "x", "string to build the suffix tree from")

		return func() { gostr.McCreight(x).ToDot(os.Stdout) }
	}

	return cli.NewCommand("st", "display a suffix tree",
		"Display a suffix tree", init)
}

func ShowMenu() *cli.Command {
	return cli.NewMenu("show", "display data structures", `
	Display data structures for various string algorithms.
	`,
		ShowSuffixTreeCommand())
}
