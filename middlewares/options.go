package middlewares

import (
	"strings"
	"web"
)

func Options() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			if ctx.Request.Method == web.MethodOptions {
				if val, ok := ctx.Request.Headers[web.HeaderOrigin]; ok {
					ctx.Response.Headers[web.HeaderAccessControlAllowOrigin] = val
				}
				ctx.Response.Headers[web.HeaderAccessControlAllowMethods] = []string{strings.Join(web.MethodsStr, " ")}
			}

			return next(ctx)
		}
	}
}
