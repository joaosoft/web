package main

import (
	"fmt"
	"webserver"
)

func main() {
	// create a new server
	w, err := webserver.NewWebServer()
	if err != nil {
		panic(err)
	}

	// add middleware's
	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())
	w.AddRoutes(
		webserver.Route{Method: webserver.MethodGet, Path: "/hello/:name", Handler: HandlerHello},
		webserver.Route{Method: webserver.MethodPost, Path: "/hello/:name", Handler: HandlerHello},
	)

	// start the server
	if err := w.Start(); err != nil {
		panic(err)
	}

}

func MyMiddlewareOne() webserver.MiddlewareFunc {
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(ctx *webserver.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE ONE")
			return next(ctx)
		}
	}
}

func MyMiddlewareTwo() webserver.MiddlewareFunc {
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(ctx *webserver.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE TWO")
			return next(ctx)
		}
	}
}

func HandlerHello(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER")
	return nil
}
