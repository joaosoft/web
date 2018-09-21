# web
[![Build Status](https://travis-ci.org/joaosoft/web.svg?branch=master)](https://travis-ci.org/joaosoft/web) | [![codecov](https://codecov.io/gh/joaosoft/web/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/web) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/web)](https://goreportcard.com/report/github.com/joaosoft/web) | [![GoDoc](https://godoc.org/github.com/joaosoft/web?status.svg)](https://godoc.org/github.com/joaosoft/web)

A simple and fast web server and client.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for 
* Common http methods
* Single/Multiple File Upload
* Single/Multiple File Download

## With attachment modes
* [default] zip files when returns more then one file 
  - on client WithClientAttachmentMode(web.MultiAttachmentModeZip)
  - on server WithServerAttachmentMode(web.MultiAttachmentModeZip)
* [experimental] returns attachmentes splited by a boundary defined on header Content-Type 
  - on client WithClientAttachmentMode(web.MultiAttachmentModeBoundary)
  - on server WithServerAttachmentMode(web.MultiAttachmentModeBoundary)

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
go get github.com/joaosoft/web
```

## Usage 
This examples are available in the project at [web/examples](https://github.com/joaosoft/web/tree/master/examples)

### Server
```
func main() {
	// create a new server
	w, err := web.NewServer(web.WithServerMultiAttachmentMode(web.MultiAttachmentModeBoundary))
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
		web.NewRoute(web.MethodGet, "/hello/:name/download", HandlerHelloForDownloadFiles),
		web.NewRoute(web.MethodGet, "/hello/:name/download/one", HandlerHelloForDownloadOneFile),
		web.NewRoute(web.MethodPost, "/hello/:name/upload", HandlerHelloForUploadFiles),
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
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPost(ctx *web.Context) error {
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

func HandlerHelloForPut(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PUT")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDelete(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DELETE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForPatch(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR PATCH")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
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
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForOptions(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR OPTIONS")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForTrace(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR TRACE")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForLink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR LINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUnlink(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UNLINK")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
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
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
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
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForView(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR VIEW")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadOneFile(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD ONE FILE")

	dir, _ := os.Getwd()
	ctx.Response.Attachment(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), "text_a.text")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadFiles(ctx *web.Context) error {
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

func HandlerHelloForUploadFiles(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UPLOAD FILES")

	fmt.Printf("\nAttachments: %+v\n", ctx.Request.Attachments)
	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}
```

### Client
```
func main() {
	// create a new client
	c, err := web.NewClient(web.WithClientMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)

	requestGetBoundary(c)

}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestGetBoundary(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao/download")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
