# webserver
[![Build Status](https://travis-ci.org/joaosoft/webserver.svg?branch=master)](https://travis-ci.org/joaosoft/webserver) | [![codecov](https://codecov.io/gh/joaosoft/webserver/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/webserver) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/webserver)](https://goreportcard.com/report/github.com/joaosoft/webserver) | [![GoDoc](https://godoc.org/github.com/joaosoft/webserver?status.svg)](https://godoc.org/github.com/joaosoft/webserver)

A simple and fast web server.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for 
* Common http methods
* Single/Multiple File Upload
* Single/Multiple File Download

## With attachment modes
* [default] zip files when returns more then one file (WithServerAttachmentMode(webserver.MultiAttachmentModeZip))
* [experimental] returns attachmentes splited by a boundary defined on header Content-Type (WithServerAttachmentMode(webserver.MultiAttachmentModeBoundary))

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
	w, err := server.NewWebServer()
	if err != nil {
		panic(err)
	}

	// add middleware's
	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())
	w.AddRoutes(
		server.NewRoute(web.MethodHead, "/hello/:name", HandlerHelloForHead),
		server.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareThree()),
		server.NewRoute(web.MethodPost, "/hello/:name", HandlerHelloForPost),
		server.NewRoute(web.MethodPut, "/hello/:name", HandlerHelloForPut),
		server.NewRoute(web.MethodDelete, "/hello/:name", HandlerHelloForDelete),
		server.NewRoute(web.MethodPatch, "/hello/:name", HandlerHelloForPatch),
		server.NewRoute(web.MethodCopy, "/hello/:name", HandlerHelloForCopy),
		server.NewRoute(web.MethodConnect, "/hello/:name", HandlerHelloForConnect),
		server.NewRoute(web.MethodOptions, "/hello/:name", HandlerHelloForOptions),
		server.NewRoute(web.MethodTrace, "/hello/:name", HandlerHelloForTrace),
		server.NewRoute(web.MethodLink, "/hello/:name", HandlerHelloForLink),
		server.NewRoute(web.MethodUnlink, "/hello/:name", HandlerHelloForUnlink),
		server.NewRoute(web.MethodPurge, "/hello/:name", HandlerHelloForPurge),
		server.NewRoute(web.MethodLock, "/hello/:name", HandlerHelloForLock),
		server.NewRoute(web.MethodUnlock, "/hello/:name", HandlerHelloForUnlock),
		server.NewRoute(web.MethodPropFind, "/hello/:name", HandlerHelloForPropFind),
		server.NewRoute(web.MethodView, "/hello/:name", HandlerHelloForView),
		server.NewRoute(web.MethodGet, "/hello/:name/download", HandlerHelloForDownloadFiles),
		server.NewRoute(web.MethodGet, "/hello/:name/download/one", HandlerHelloForDownloadOneFile),
		server.NewRoute(web.MethodPost, "/hello/:name/upload", HandlerHelloForUploadFiles),
	)

	w.AddNamespace("/p").AddRoutes(
		server.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareFour()),
	)

	// start the server
	if err := w.Start(); err != nil {
		panic(err)
	}
}

func MyMiddlewareOne() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE ONE")
			return next(ctx)
		}
	}
}

func MyMiddlewareTwo() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE TWO")
			return next(ctx)
		}
	}
}

func MyMiddlewareThree() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE THREE")
			return next(ctx)
		}
	}
}
func MyMiddlewareFour() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE FOUR")
			return next(ctx)
		}
	}
}
func HandlerHelloForHead(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR HEAD")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForGet(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPost(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	ctx.Request.Bind(&data)
	fmt.Printf("%+v", data)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPut(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PUT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDelete(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DELETE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPatch(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PATCH")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForCopy(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR COPY")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForConnect(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR CONNECT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForOptions(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR OPTIONS")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForTrace(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR TRACE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForLink(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlink(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPurge(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PURGE")

	return ctx.Response.NoContent(web.StatusOK)
}

func HandlerHelloForLock(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LOCK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlock(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLOCK")

	return ctx.Response.NoContent(web.StatusOK)
}
func HandlerHelloForPropFind(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PROPFIND")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForView(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR VIEW")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadOneFile(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD ONE FILE")

	dir, _ := os.Getwd()
	ctx.Response.Attachment(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), "text_a.text")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadFiles(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD FILES")

	dir, _ := os.Getwd()
	ctx.Response.Attachment(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), "text_a.text")
	ctx.Response.Attachment(fmt.Sprintf("%s%s", dir, "/examples/data/b.text"), "text_b.text")
	ctx.Response.Inline(fmt.Sprintf("%s%s", dir, "/examples/data/c.text"), "text_c.text")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUploadFiles(ctx *server.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UPLOAD FILES")

	fmt.Printf("\nAttachments: %+v\n", ctx.Request.Attachments)
	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
