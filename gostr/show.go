package main

import (
	"os"

	"github.com/mailund/cli"
	"github.com/mailund/gostr"
)

type argShowST struct {
	Out cli.OutFile `flag:"out" short:"o" descr:"file to write the dot description to"`
	X   string      `pos:"x" descr:"the string to build the suffix tree from"`
}

var showST = cli.NewCommand(cli.CommandSpec{
	Name:  "st",
	Short: "shows the structure of a suffix tree for string x",
	Long:  "Shows the structure of a suffix tree for string x.",
	Init: func() interface{} {
		return &argShowST{Out: cli.OutFile{Writer: os.Stdout}}
	},
	Action: func(i interface{}) {
		args, _ := i.(*argShowST)
		gostr.McCreight(args.X).ToDot(args.Out)
	},
})

type argShowTrie struct {
	Out cli.OutFile `flag:"out" short:"o" descr:"file to write the dot description to"`
	X   []string    `pos:"strings" descr:"the strings to build the trie from"`
}

var showTrie = cli.NewCommand(cli.CommandSpec{
	Name:  "trie",
	Short: "shows the structure of a trie",
	Long:  "Shows the structure of a trie.",
	Init: func() interface{} {
		return &argShowTrie{Out: cli.OutFile{Writer: os.Stdout}}
	},
	Action: func(i interface{}) {
		args, _ := i.(*argShowTrie)
		gostr.BuildTrie(args.X).ToDot(args.Out)
	},
})

var show = cli.NewMenu(
	"show",
	"shows algorithms and data structures",
	"Shows algorithms and data structures. Pick a subcommand",
	showST, showTrie)
