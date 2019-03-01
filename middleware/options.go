package middleware

import (
	"strings"
	"web"
)

func Options() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			if ctx.Request.Method == web.MethodOptions {
				var method web.Method
				if val, ok := ctx.Request.Headers[web.HeaderAccessControlRequestMethod]; ok {
					method = web.Method(val[0])
				} else {
					return ctx.Response.NoContent(web.StatusBadRequest)
				}

				route, err := ctx.Request.Server.GetRoute(method, ctx.Request.Address.Url)
				if err == nil && route != nil {
					ctx.Response.Headers[web.HeaderAccessControlAllowMethods] = []string{string(ctx.Request.Method)}
					ctx.Response.Headers[web.HeaderAccessControlAllowHeaders] = []string{strings.Join([]string{
						string(web.HeaderContentType),
						string(web.HeaderAccessControlAllowHeaders),
						string(web.HeaderAuthorization),
						string(web.HeaderXRequestedWith),
					}, ", ")}
				} else if err != web.ErrorNotFound {
					return ctx.Response.NoContent(web.StatusNotFound)
				} else {
					return ctx.Response.NoContent(web.StatusBadRequest)
				}
			}

			return next(ctx)
		}
	}
}
