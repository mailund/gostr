package main

import (
	"flag"
	"os"

	"github.com/mailund/cli"
	"github.com/mailund/cli/params"
	"github.com/mailund/gostr"
)

func ShowSuffixTreeCommand() *cli.Command {
	init := func(f *flag.FlagSet, p *params.ParamSet) (cli.UsageFunc, cli.RunFunc) {
		var x string
		p.StringVar(&x, "x", "string to build the suffix tree from")

		run := func(args []string) {
			gostr.McCreight(x).ToDot(os.Stdout)
		}

		return cli.DefaultUsage("st", f, p), run
	}
	return cli.NewCommand("st", "display a suffix tree", init)
}

func ShowMenu() *cli.Command {
	return cli.NewMenu("show", "display data structures",
		ShowSuffixTreeCommand())
}
