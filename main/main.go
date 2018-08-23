package main

import (
	"dependency"
	"fmt"
	"os"
)

func main() {
	cmd := dependency.CmdDependencyGet

	args := os.Args
	if len(args) > 1 {
		cmd = dependency.CmdDependency(args[1])
	}

	d := dependency.NewDependency()

	switch cmd {
	case dependency.CmdDependencyGet:
		if err := d.Get(); err != nil {
			panic(err)
		}
	case dependency.CmdDependencyReset:
		if err := d.Reset(); err != nil {
			panic(err)
		}
	default:
		fmt.Printf("invalid command! available commands are [%s, %s]", dependency.CmdDependencyGet, dependency.CmdDependencyReset)
	}

}
