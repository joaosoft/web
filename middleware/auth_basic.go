package middleware

import (
	"web"

	"github.com/joaosoft/auth-types/basic"
)

func CheckAuthBasic(user, password string) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) error {

			authHeader := ctx.Request.GetHeader(web.HeaderAuthorization)

			ok, err := basic.New().Check(authHeader, func(username string) (*basic.Credentials, error) {
				return &basic.Credentials{
					UserName: user,
					Password: password,
				}, nil
			})

			if !ok || err != nil {
				return ErrorInvalidAuthorization
			}

			return next(ctx)
		}
	}
}
