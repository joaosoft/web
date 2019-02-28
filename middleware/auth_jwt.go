package middleware

import (
	"github.com/joaosoft/auth-types/jwt"
	"strings"
	"web"
)

func CheckAuthJwt(keyFunc jwt.KeyFunc, checkFunc jwt.CheckFunc) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {
			authHeader := ctx.Request.GetHeader(web.HeaderAuthorization)
			token := strings.Replace(authHeader, "Bearer ", "", 1)

			ok, err := jwt.Check(token, keyFunc, checkFunc, jwt.Claims{}, true)

			if !ok || err != nil {
				return ErrorInvalidAuthorization
			}

			return next(ctx)
		}
	}
}
