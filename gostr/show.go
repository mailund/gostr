package main

import (
	"os"

	"github.com/mailund/cli"
	"github.com/mailund/gostr"
)

type argShowST struct {
	X string `pos:"x" descr:"the string to build the suffix tree from"`
}

var showST = cli.NewCommand(cli.CommandSpec{
	Name:  "st",
	Short: "shows the structure of a suffix tree for string x",
	Long:  "Shows the structure of a suffix tree for string x.",
	Init:  func() interface{} { return new(argShowST) },
	Action: func(args interface{}) {
		gostr.McCreight(args.(*argShowST).X).ToDot(os.Stdout)
	},
})

var show = cli.NewMenu(
	"show",
	"shows algorithms and data structures",
	"Shows algorithms and data structures. Pick a subcommand",
	showST)
