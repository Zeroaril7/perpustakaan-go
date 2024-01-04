package handlers

import (
	"net/http"

	"github.com/Zeroaril7/perpustakaan-go/modules/auth/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Login(c echo.Context) error
}

type authHandler struct {
	authUsecase domain.AuthUsecase
}

// Login implements AuthHandler.
func (h *authHandler) Login(c echo.Context) error {
	authRequest := new(models.LoginAuth)

	if err := c.Bind(authRequest); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(authRequest); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	result := <-h.authUsecase.AuthWithPassword(c.Request().Context(), *authRequest)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Login success", http.StatusOK, c)
}

func NewAuthHandler(e *echo.Echo, authUsecase domain.AuthUsecase) AuthHandler {
	handler := &authHandler{
		authUsecase: authUsecase,
	}

	group := e.Group("/auth")
	group.POST("/login", handler.Login)

	return handler
}
