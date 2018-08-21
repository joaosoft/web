package main

import (
	"fmt"
	"watcher/service"
)

func main() {
	c := make(chan *service.Event)
	w, err := service.NewWatcher(service.WithEventChannel(c))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	if err := w.Start(); err != nil {
		panic(err)
	}
}
