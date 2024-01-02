package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func VerifyBasicAuth(username, password string) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(user, pwd string, ctx echo.Context) (bool, error) {
			return user == username && pwd == password, nil
		},
	})
}
