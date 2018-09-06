# webserver
[![Build Status](https://travis-ci.org/joaosoft/webserver.svg?branch=master)](https://travis-ci.org/joaosoft/webserver) | [![codecov](https://codecov.io/gh/joaosoft/webserver/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/webserver) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/webserver)](https://goreportcard.com/report/github.com/joaosoft/webserver) | [![GoDoc](https://godoc.org/github.com/joaosoft/webserver?status.svg)](https://godoc.org/github.com/joaosoft/webserver)

A simple web server. [UNDER DEVELOPMENT]

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Methods (HEAD, GET, POST, PUT, CONNECT, PATCH, DELETE, OPTIONS, TRACE)

>### Go
```
go get github.com/joaosoft/webserver
```

## Usage 
```
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
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
