package main

import (
	"flag"
	"os"

	"github.com/mailund/gostr"
)

func ShowSuffixTreeCommand() *Command {
	init := func(f *flag.FlagSet, a *ArgSet) (UsageFunc, RunFunc) {
		var x string
		a.StringVar(&x, "x", "string to build the suffix tree from")

		run := func(args []string) {
			gostr.McCreight(x).ToDot(os.Stdout)
		}

		return DefaultUsage("st", f, a), run
	}
	return NewCommand("st", "display a suffix tree", init)
}

func ShowMenu() *Command {
	return NewMenu("show", "display data structures",
		ShowSuffixTreeCommand())
}
