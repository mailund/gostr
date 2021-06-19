package main

import (
	"fmt"
	"io"
	"os"
)

type Parser = func(arg string)

type Param struct {
	name   string
	desc   string
	parser Parser
}

// FIXME: find a way to specify whether you can have more of an
// argument or whether the argset must match completely

type ArgSet struct {
	name   string
	params []*Param
	args   []string // will be set to remaining args after parsing
	out    io.Writer

	// You can assign to usage to change the help info.
	Usage UsageFunc
}

func (args *ArgSet) Output() io.Writer { return args.out }

func (args *ArgSet) PrintDefaults() {
	fmt.Fprintf(args.out, "Arguments:\n")
	for _, param := range args.params {
		fmt.Fprintf(args.out, "  %s\n\t%s\n", param.name, param.desc)
	}
}

func NewArgSet(name string) *ArgSet {
	argset := &ArgSet{name: name,
		params: []*Param{}, args: []string{}}
	argset.Usage = argset.PrintDefaults
	argset.out = os.Stderr
	return argset
}

func (argset *ArgSet) ParamNames() []string {
	names := make([]string, len(argset.params))
	for i, param := range argset.params {
		names[i] = param.name
	}
	return names
}

func (argset *ArgSet) Parse(args []string) {
	if len(args) < len(argset.params) {
		fmt.Fprintf(argset.out,
			"Insufficient arguments for command '%s'\n\n",
			argset.name)
		argset.Usage()
		os.Exit(1) // FIXME: control by flag?
	}
	for i, param := range argset.params {
		param.parser(args[i])
	}
	argset.args = args[len(argset.params):]
}

func (argset *ArgSet) Args() []string { return argset.args }

func stringParser(target *string) Parser {
	return func(arg string) { *target = arg }
}

func (args *ArgSet) StringVar(target *string, name, desc string) {
	param := &Param{name, desc, stringParser(target)}
	args.params = append(args.params, param)
}
