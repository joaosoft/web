package main

import (
	"builder"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	quit := make(chan int)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	w, err := builder.NewBuilder()
	if err != nil {
		panic(err)
	}

	if err := w.Start(&sync.WaitGroup{}); err != nil {
		panic(err)
	}

	<-termChan
	quit <- 1
}
