package handlers

import (
	"net/http"

	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/middlewares"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Add(c echo.Context) error
	Delete(c echo.Context) error
	Get(c echo.Context) error
	GetByUsername(c echo.Context) error
	Update(c echo.Context) error
}

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, userUsecase domain.UserUsecase) UserHandler {
	handler := &userHandler{
		userUsecase: userUsecase,
	}

	group := e.Group("/user", middlewares.VerifyBasicAuth(config.Config().BasicAuthUsername, config.Config().BasicAuthPassword))
	group.DELETE("/:username", handler.Delete)
	group.GET("", handler.Get)
	group.GET("/:username", handler.GetByUsername)
	group.POST("", handler.Add)
	group.PUT("/:username", handler.Update)

	return handler
}

// Add implements UserHandler.
func (h *userHandler) Add(c echo.Context) error {
	data := new(models.UserAdd)

	if err := c.Bind(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	expend := models.User{}
	expend = data.ToUser(expend)

	result := <-h.userUsecase.Add(c.Request().Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Add user success", http.StatusOK, c)
}

// Delete implements UserHandler.
func (h *userHandler) Delete(c echo.Context) error {
	username := utils.ConvertString(c.Param("username"))

	result := <-h.userUsecase.Delete(c.Request().Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(nil, "Delete User success", http.StatusOK, c)
}

// Get implements UserHandler.
func (h *userHandler) Get(c echo.Context) error {
	filter := new(models.UserFilter)

	if err := c.Bind(filter); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if !filter.DisablePagination {
		filter.SetDefault()
	}

	result := <-h.userUsecase.Get(c.Request().Context(), *filter)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.ResponseWithPagination(result.Data, "Get user success", http.StatusOK, result.Total, filter.GetPaginationRequest(), c)
}

// GetByUsername implements UserHandler.
func (h *userHandler) GetByUsername(c echo.Context) error {
	username := utils.ConvertString(c.Param("username"))

	result := <-h.userUsecase.GetByUsername(c.Request().Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Get user success", http.StatusOK, c)
}

// Update implements UserHandler.
func (h *userHandler) Update(c echo.Context) error {
	username := utils.ConvertString(c.Param("username"))

	result := <-h.userUsecase.GetByUsername(c.Request().Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	expend := result.Data.(models.User)

	if expend == (models.User{}) {
		return utils.ResponseError(httperror.NotFound("User not found"), c)
	}

	data := new(models.UserAdd)
	if err := c.Bind(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	expend = data.ToUser(expend)

	result = <-h.userUsecase.Update(c.Request().Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Update user success", http.StatusOK, c)
}
