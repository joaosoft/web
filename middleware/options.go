package middleware

import (
	"strings"
	"web"
)

func Options() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			if ctx.Request.Method == web.MethodOptions {
				ctx.Response.Headers[web.HeaderAccessControlAllowMethods] = []string{strings.Join(web.MethodsStr, ", ")}
				ctx.Response.Headers[web.HeaderAccessControlAllowHeaders] = []string{strings.Join([]string{
					string(web.HeaderContentType),
					string(web.HeaderAccessControlAllowHeaders),
					string(web.HeaderAuthorization),
					string(web.HeaderXRequestedWith),
				}, ", ")}
			}

			return next(ctx)
		}
	}
}
