package main

import (
	"dependency/service"
	"fmt"
	"os"
)

func main() {
	cmd := service.CmdDependencyGet

	args := os.Args
	if len(args) > 1 {
		cmd = service.CmdDependency(args[1])
	}

	d, err := service.NewDependency()
	if err != nil {
		panic(err)
	}

	switch cmd {
	case service.CmdDependencyGet:
		if err := d.Get(); err != nil {
			panic(err)
		}
	case service.CmdDependencyUpdate:
		if err := d.Update(); err != nil {
			panic(err)
		}
	case service.CmdDependencyReset:
		if err := d.Reset(); err != nil {
			panic(err)
		}
	default:
		fmt.Printf("invalid command! available commands are [%s, %s]", service.CmdDependencyGet, service.CmdDependencyReset)
	}

}
