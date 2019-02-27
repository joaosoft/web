package middleware

import (
	"encoding/base64"
	"strings"
	"web"
)

func AuthBasic(username, password string) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			authHeader := ctx.Request.GetHeader(web.HeaderAuthorization)

			authDecoded, err := base64.StdEncoding.DecodeString(authHeader)
			if err != nil {
				return web.ErrorInvalidAuth
			}

			split := strings.SplitN(string(authDecoded), ":", 2)

			if len(split) == 2 && split[0] == username && split[1] == password {
				return next(ctx)
			}

			return web.ErrorInvalidAuth
		}
	}
}
