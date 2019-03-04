package main

import (
	"fmt"
	"os"
	"web"

	"github.com/joaosoft/auth-types/jwt"
)

func main() {
	// create a new server
	w, err := web.NewServer(web.WithServerMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	claims := jwt.Claims{"name": "joao", "age": 30}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte("bananas"), nil
	}

	checkFunc := func(c jwt.Claims) (bool, error) {
		if claims["name"] == c["name"].(string) &&
			claims["age"] == int(c["age"].(float64)) {
			return true, nil
		}
		return false, fmt.Errorf("invalid jwt session token")
	}

	// add middleware's
	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())
	w.AddRoutes(
		web.NewRoute(web.MethodOptions, "*", HandlerHelloForOptions, web.MiddlewareOptions()),
		web.NewRoute(web.MethodGet, "/auth-basic", HandlerForGet, web.MiddlewareCheckAuthBasic("joao", "ribeiro")),
		web.NewRoute(web.MethodGet, "/auth-jwt", HandlerForGet, web.MiddlewareCheckAuthJwt(keyFunc, checkFunc)),
		web.NewRoute(web.MethodHead, "/hello/:name", HandlerHelloForHead),
		web.NewRoute(web.MethodGet, "/hello/:name", HandlerHelloForGet, MyMiddlewareThree()),
		web.NewRoute(web.MethodPost, "/hello/:name", HandlerHelloForPost),
		web.NewRoute(web.MethodPut, "/hello/:name", HandlerHelloForPut),
		web.NewRoute(web.MethodDelete, "/hello/:name", HandlerHelloForDelete),
		web.NewRoute(web.MethodPatch, "/hello/:name", HandlerHelloForPatch),
		web.NewRoute(web.MethodCopy, "/hello/:name", HandlerHelloForCopy),
		web.NewRoute(web.MethodConnect, "/hello/:name", HandlerHelloForConnect),
		web.NewRoute(web.MethodOptions, "/hello/:name", HandlerHelloForOptions, web.MiddlewareOptions()),
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

func HandlerForGet(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR GET")

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \"guest\" }"),
	)
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

	return ctx.Response.NoContent(web.StatusNoContent)
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
	body, _ := web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), nil)
	ctx.Response.Attachment("text_a.text", body)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForDownloadFiles(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR DOWNLOAD FILES")

	dir, _ := os.Getwd()
	body, _ := web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/a.text"), nil)
	ctx.Response.Attachment("text_a.text", body)

	body, _ = web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/b.text"), nil)
	ctx.Response.Attachment("text_b.text", body)

	body, _ = web.ReadFile(fmt.Sprintf("%s%s", dir, "/examples/data/c.text"), nil)
	ctx.Response.Inline("text_c.text", body)

	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}

func HandlerHelloForUploadFiles(ctx *web.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR UPLOAD FILES")

	fmt.Printf("\nAttachments: %+v\n", ctx.Request.FormData)
	return ctx.Response.Bytes(
		web.StatusOK,
		web.ContentTypeApplicationJSON,
		[]byte("{ \"welcome\": \""+ctx.Request.UrlParams["name"][0]+"\" }"),
	)
}
