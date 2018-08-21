package main

import (
	"builder"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	build := builder.NewBuilder(builder.WithReloadTime(1))

	if err := build.Start(nil); err != nil {
		panic(err)
	}

	<-termChan
	if err := build.Stop(nil); err != nil {
		panic(err)
	}
}
