package main

import (
	"fmt"
	"os"

	"github.com/mailund/cli"
)

func main() {
	var main *cli.Command = cli.NewMenu(
		"gostr", "run gostr commands", "",
		ShowMenu())

	if len(os.Args) < 1 {
		fmt.Println("no args")
		main.Usage()
		return
	} else {
		main.Run(os.Args[1:])
	}
}
