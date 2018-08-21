package main

import (
	"builder"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	w, err := builder.NewBuilder(builder.WithReloadTime(1))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := w.Start(&wg); err != nil {
		panic(err)
	}

	<-termChan
	wg.Add(1)
	if err := w.Stop(&wg); err != nil {
		panic(err)
	}
}
