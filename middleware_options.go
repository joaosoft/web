package web

import (
	"strings"
)

func MiddlewareOptions() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {

			if ctx.Request.Method == MethodOptions {
				ctx.Response.Headers[HeaderAccessControlAllowMethods] = []string{strings.Join(MethodsStr, ", ")}
				ctx.Response.Headers[HeaderAccessControlAllowHeaders] = []string{strings.Join([]string{
					string(HeaderContentType),
					string(HeaderAccessControlAllowHeaders),
					string(HeaderAuthorization),
					string(HeaderXRequestedWith),
				}, ", ")}
			}

			return next(ctx)
		}
	}
}
