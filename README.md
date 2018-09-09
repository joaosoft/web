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
	w, err := web.NewWebServer()
	if err != nil {
		panic(err)
	}

	// add middleware's
	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())
	w.AddRoutes(
		web.NewRoute(web.MethodHead, "/hello/:name", HandlerHelloForHead),
		web.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareThree()),
		web.NewRoute(web.MethodPost, "/hello/:name", HandlerHelloForPost),
		web.NewRoute(web.MethodPut, "/hello/:name", HandlerHelloForPut),
		web.NewRoute(web.MethodDelete, "/hello/:name", HandlerHelloForDelete),
		web.NewRoute(web.MethodPatch, "/hello/:name", HandlerHelloForPatch),
		web.NewRoute(web.MethodCopy, "/hello/:name", HandlerHelloForCopy),
		web.NewRoute(web.MethodConnect, "/hello/:name", HandlerHelloForConnect),
		web.NewRoute(web.MethodOptions, "/hello/:name", HandlerHelloForOptions),
		web.NewRoute(web.MethodTrace, "/hello/:name", HandlerHelloForTrace),
		web.NewRoute(web.MethodLink, "/hello/:name", HandlerHelloForLink),
		web.NewRoute(web.MethodUnlink, "/hello/:name", HandlerHelloForUnlink),
		web.NewRoute(web.MethodPurge, "/hello/:name", HandlerHelloForPurge),
		web.NewRoute(web.MethodLock, "/hello/:name", HandlerHelloForLock),
		web.NewRoute(web.MethodUnlock, "/hello/:name", HandlerHelloForUnlock),
		web.NewRoute(web.MethodPropFind, "/hello/:name", HandlerHelloForPropFind),
		web.NewRoute(web.MethodView, "/hello/:name", HandlerHelloForView),
	)

	w.AddNamespace("/p").AddRoutes(
		web.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareFour()),
	)

	// start the server
	if err := w.Start(); err != nil {
		panic(err)
	}
}

func MyMiddlewareOne() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE ONE")
			return next(ctx)
		}
	}
}

func MyMiddlewareTwo() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE TWO")
			return next(ctx)
		}
	}
}

func MyMiddlewareThree() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE THREE")
			return next(ctx)
		}
	}
}
func MyMiddlewareFour() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE FOUR")
			return next(ctx)
		}
	}
}
func HandlerHelloForHead(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR HEAD")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForGet(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPost(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPut(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PUT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForDelete(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DELETE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPatch(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PATCH")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForCopy(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR COPY")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForConnect(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR CONNECT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForOptions(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR OPTIONS")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForTrace(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR TRACE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForLink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForPurge(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PURGE")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForLock(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LOCK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlock(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLOCK")

	return ctx.Response.NoContent(web.StatusOK)
}
func HandlerHelloForPropFind(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PROPFIND")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}

func HandlerHelloForView(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR VIEW")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParms["name"][0]+"\" }"),
	)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
