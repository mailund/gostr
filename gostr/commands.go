package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type (
	UsageFunc = func()
	RunFunc   = func(args []string)
	InitFunc  = func(f *flag.FlagSet, a *ArgSet) (UsageFunc, RunFunc)
)

type Command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Args        *ArgSet
	UsageCB     UsageFunc
	RunCB       RunFunc
}

func (cmd Command) Usage() { cmd.UsageCB() }
func (cmd Command) Run(args []string) {
	cmd.Flags.Parse(args)
	cmd.Args.Parse(cmd.Flags.Args())
	cmd.RunCB(cmd.Args.Args())
}

func NewCommand(
	name, description string,
	init InitFunc) *Command {
	cmd := &Command{Name: name, Description: description}
	cmd.Flags = flag.NewFlagSet(name, flag.ExitOnError)
	cmd.Args = NewArgSet(name) // flags?
	cmd.UsageCB, cmd.RunCB = init(cmd.Flags, cmd.Args)
	cmd.Flags.Usage = cmd.UsageCB
	cmd.Args.Usage = cmd.UsageCB

	// There is always a help command when we parse, but the usage won't
	// show it unless we make it explicit
	cmd.Flags.Bool("help", false, fmt.Sprintf("show help for %s", name))

	return cmd
}

func DefaultUsage(name string, f *flag.FlagSet, a *ArgSet) UsageFunc {
	return func() {
		fmt.Fprintf(f.Output(),
			"Usage: %s [options] %s\n\n",
			name, strings.Join(a.ParamNames(), " "))
		fmt.Fprintf(f.Output(), "Options:\n")
		f.PrintDefaults()
		fmt.Fprintf(f.Output(), "\n")
		a.PrintDefaults()
	}
}

func NewMenu(name, description string, subcmds ...*Command) *Command {
	subcommands := map[string]*Command{}
	for _, cmd := range subcmds {
		subcommands[cmd.Name] = cmd
	}

	init := func(f *flag.FlagSet, a *ArgSet) (UsageFunc, RunFunc) {
		var command string
		a.StringVar(&command, "cmd", "command to run")

		usage := func() {
			DefaultUsage(name, f, a)()
			fmt.Fprintf(f.Output(), "\nCommands:\n")
			for name, cmd := range subcommands {
				fmt.Fprintf(f.Output(), "\t%s\t%s\n", name, cmd.Description)
			}
			fmt.Fprintf(f.Output(), "\n")
		}
		run := func(args []string) {
			if subcmd, ok := subcommands[command]; ok {
				subcmd.Run(args)
			} else {
				fmt.Fprintf(f.Output(),
					"'%s' is not a valid command for %s.\n\n", command, name)
				f.Usage()
				os.Exit(1)
			}
		}
		return usage, run
	}

	return NewCommand(name, description, init)
}
