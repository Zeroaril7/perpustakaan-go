package middlewares

import (
	"log"

	"github.com/Zeroaril7/perpustakaan-go/pkg/sdk/jwtrsa"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func VerifyJWTRSA(publicKey string) echo.MiddlewareFunc {
	verifyPublicKey, err := jwtrsa.GetPublicKey(publicKey)

	if err != nil {
		log.Default().Printf("%s", err.Error())
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    verifyPublicKey,
		SigningMethod: "RS256",
	})
}
