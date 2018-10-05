package web

import (
	"strings"
)

func MiddlewareOptions() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {

			if ctx.Request.Method == MethodOptions {
				if val, ok := ctx.Request.Headers[HeaderOrigin]; ok {
					ctx.Response.Headers[HeaderAccessControlAllowOrigin] = val
				}
				ctx.Response.Headers[HeaderAccessControlAllowMethods] = []string{strings.Join(MethodsStr, " ")}
			}

			return next(ctx)
		}
	}
}
