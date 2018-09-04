package main

import "webserver"

func main() {
	w, err := webserver.NewWebServer()
	if err != nil {
		panic(err)
	}

	if err := w.Start(); err != nil {
		panic(err)
	}
}
