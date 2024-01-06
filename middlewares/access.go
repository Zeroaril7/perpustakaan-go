package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func EchoSetCredential() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, ok := c.Get("user").(*jwt.Token); !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Token is invalid")
			}

			claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])

			return next(c)
		}
	}
}
