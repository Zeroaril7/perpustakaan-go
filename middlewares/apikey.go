package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func VerifyAPIKey(keyAuth string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-KEY",
		Validator: func(auth string, c echo.Context) (bool, error) {
			return auth == keyAuth, nil
		},
	})
}
