package middleware

import (
	"fmt"
	"web"
)

func AuthJwt(token string) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			authHeader := ctx.Request.GetHeader(web.HeaderAuthorization)
			fmt.Println(authHeader)
			return next(ctx)
		}
	}
}
