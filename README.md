# webserver
[![Build Status](https://travis-ci.org/joaosoft/webserver.svg?branch=master)](https://travis-ci.org/joaosoft/webserver) | [![codecov](https://codecov.io/gh/joaosoft/webserver/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/webserver) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/webserver)](https://goreportcard.com/report/github.com/joaosoft/webserver) | [![GoDoc](https://godoc.org/github.com/joaosoft/webserver?status.svg)](https://godoc.org/github.com/joaosoft/webserver)

A simple and fast web server.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for methods
* HEAD
* GET
* POST
* PUT
* DELETE
* PATCH
* COPY
* CONNECT
* OPTIONS
* TRACE
* LINK
* UNLINK
* PURGE
* LOCK
* UNLOCK
* PROPFIND
* VIEW

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
		webserver.NewRoute(webserver.MethodHead, "/hello/:name", HandlerHelloForHead),
		webserver.NewRoute(webserver.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareThree()),
		webserver.NewRoute(webserver.MethodPost, "/hello/:name", HandlerHelloForPost),
		webserver.NewRoute(webserver.MethodPut, "/hello/:name", HandlerHelloForPut),
		webserver.NewRoute(webserver.MethodDelete, "/hello/:name", HandlerHelloForDelete),
		webserver.NewRoute(webserver.MethodPatch, "/hello/:name", HandlerHelloForPatch),
		webserver.NewRoute(webserver.MethodCopy, "/hello/:name", HandlerHelloForCopy),
		webserver.NewRoute(webserver.MethodConnect, "/hello/:name", HandlerHelloForConnect),
		webserver.NewRoute(webserver.MethodOptions, "/hello/:name", HandlerHelloForOptions),
		webserver.NewRoute(webserver.MethodTrace, "/hello/:name", HandlerHelloForTrace),
		webserver.NewRoute(webserver.MethodLink, "/hello/:name", HandlerHelloForLink),
		webserver.NewRoute(webserver.MethodUnlink, "/hello/:name", HandlerHelloForUnlink),
		webserver.NewRoute(webserver.MethodPurge, "/hello/:name", HandlerHelloForPurge),
		webserver.NewRoute(webserver.MethodLock, "/hello/:name", HandlerHelloForLock),
		webserver.NewRoute(webserver.MethodUnlock, "/hello/:name", HandlerHelloForUnlock),
		webserver.NewRoute(webserver.MethodPropFind, "/hello/:name", HandlerHelloForPropFind),
		webserver.NewRoute(webserver.MethodView, "/hello/:name", HandlerHelloForView),
	)

	w.AddNamespace("/p").AddRoutes(
		webserver.NewRoute(webserver.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareFour()),
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

func MyMiddlewareThree() webserver.MiddlewareFunc {
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(ctx *webserver.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE THREE")
			return next(ctx)
		}
	}
}
func MyMiddlewareFour() webserver.MiddlewareFunc {
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(ctx *webserver.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE FOUR")
			return next(ctx)
		}
	}
}
func HandlerHelloForHead(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR HEAD")

	return ctx.Response.NoContent(webserver.StatusOK)
}

func HandlerHelloForGet(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPost(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPut(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PUT")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForDelete(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DELETE")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPatch(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PATCH")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForCopy(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR COPY")

	return ctx.Response.NoContent(webserver.StatusOK)
}

func HandlerHelloForConnect(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR CONNECT")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForOptions(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR OPTIONS")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForTrace(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR TRACE")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForLink(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LINK")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlink(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLINK")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPurge(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PURGE")

	return ctx.Response.NoContent(webserver.StatusOK)
}

func HandlerHelloForLock(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LOCK")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlock(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLOCK")

	return ctx.Response.NoContent(webserver.StatusOK)
}
func HandlerHelloForPropFind(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PROPFIND")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForView(ctx *webserver.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR VIEW")

	return ctx.Response.Bytes(
		webserver.StatusOK,
		webserver.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
