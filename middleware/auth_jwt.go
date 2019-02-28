package middleware

import (
	"strings"
	"web"

	"github.com/joaosoft/auth-types/jwt"
)

func CheckAuthJwt(keyFunc jwt.KeyFunc, checkFunc jwt.CheckFunc) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			authHeader := ctx.Request.GetHeader(web.HeaderAuthorization)
			token := strings.Replace(authHeader, "Bearer ", "", 1)

			ok, err := jwt.New().Check(token, keyFunc, checkFunc, jwt.Claims{}, false)

			if !ok || err != nil {
				return ErrorInvalidAuthorization
			}

			return next(ctx)
		}
	}
}
