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
		webserver.Route{Method: webserver.MethodGet, Path: "/hello/:name", Handler: HandlerHelloForGet},
		webserver.Route{Method: webserver.MethodPost, Path: "/hello/:name", Handler: HandlerHelloForPost},
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

func HandlerHelloForGet(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	ctx.Response.Status = webserver.StatusOK
	ctx.Response.Body = []byte("{ \"test\": \"ok\" }")

	return nil
}

func HandlerHelloForPost(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	ctx.Response.Status = webserver.StatusOK
	ctx.Response.Body = []byte("{ \"test\": \"nok\" }")

	return nil
}
