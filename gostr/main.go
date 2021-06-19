package main

import (
	"fmt"
	"os"
)

func main() {
	var main *Command = NewMenu(
		"gostr", "run gostr commands",
		ShowMenu())

	if len(os.Args) < 1 {
		fmt.Println("no args")
		main.Usage()
		return
	} else {
		main.Run(os.Args[1:])
	}

}
