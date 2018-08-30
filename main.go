package main

import (
	"dependency/service"
	"fmt"
	"os"
)

func main() {
	cmd := service.CmdDependencyGet
	protocol := service.ProtocolSSH

	args := os.Args
	if len(args) > 1 {
		cmd = service.CmdDependency(args[1])
	}

	if len(args) > 2 {
		protocol = service.Protocol(args[2])
	}

	d, err := service.NewDependency(service.WithProtocol(protocol))
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	switch cmd {
	case service.CmdDependencyGet:
		if err := d.Get(); err != nil {
			panic(err)
			os.Exit(1)
		}
	case service.CmdDependencyUpdate:
		if err := d.Update(); err != nil {
			panic(err)
			os.Exit(1)
		}
	case service.CmdDependencyReset:
		if err := d.Reset(); err != nil {
			panic(err)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid command! available commands are [%s, %s, %s]", service.CmdDependencyGet, service.CmdDependencyUpdate, service.CmdDependencyReset)
		os.Exit(1)
	}

	os.Exit(0)
}
