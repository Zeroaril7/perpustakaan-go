package handlers

import (
	"net/http"

	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/middlewares"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/labstack/echo/v4"
)

type BookHandler interface {
	Add(c echo.Context) error
	Get(c echo.Context) error
	GetByBookID(c echo.Context) error
	Delete(c echo.Context) error
	Update(c echo.Context) error
}

type bookHandler struct {
	bookUsecase domain.BookUsecase
}

func NewBookHandler(e *echo.Echo, bookUsecase domain.BookUsecase) BookHandler {
	handler := &bookHandler{
		bookUsecase: bookUsecase,
	}

	group := e.Group("/book")
	group.DELETE("/:book-id", handler.Delete, middlewares.VerifyBasicAuth(config.Config().BasicAuthUsername, config.Config().BasicAuthPassword))
	group.GET("", handler.Get)
	group.GET("/:book-id", handler.GetByBookID)
	group.POST("", handler.Add, middlewares.VerifyBasicAuth(config.Config().BasicAuthUsername, config.Config().BasicAuthPassword))
	group.PUT("/:book-id", handler.Update, middlewares.VerifyBasicAuth(config.Config().BasicAuthUsername, config.Config().BasicAuthPassword))
	return handler
}

func (h *bookHandler) Add(c echo.Context) error {
	data := new(models.BookAdd)

	if err := c.Bind(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	result := <-h.bookUsecase.GetLast(c.Request().Context(), data.Genre)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	expend := result.Data.(models.Book)

	if expend == (models.Book{}) {
		expend.Genre = data.Genre
	}
	expend = data.ToBook(expend)

	result = <-h.bookUsecase.Add(c.Request().Context(), expend)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Add book success", http.StatusOK, c)
}

// Delete implements BookHandler.
func (h *bookHandler) Delete(c echo.Context) error {
	bookID := utils.ConvertString(c.Param("book-id"))

	role := c.Get("role").(string)

	if role != constant.Admin && role != constant.SuperAdmin {
		return utils.ResponseError(httperror.Unauthorized(httperror.UnauthorizedErrorMessage), c)
	}

	result := <-h.bookUsecase.Delete(c.Request().Context(), bookID)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(nil, "Delete book success", http.StatusOK, c)
}

func (h *bookHandler) Get(c echo.Context) error {
	filter := new(models.BookFilter)

	if err := c.Bind(filter); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if !filter.DisablePagination {
		filter.SetDefault()
	}

	result := <-h.bookUsecase.Get(c.Request().Context(), *filter)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.ResponseWithPagination(result.Data, "Get book success", http.StatusOK, result.Total, filter.GetPaginationRequest(), c)
}

// GetByBookID implements BookHandler.
func (h *bookHandler) GetByBookID(c echo.Context) error {
	bookID := utils.ConvertString(c.Param("book-id"))

	result := <-h.bookUsecase.GetByBookID(c.Request().Context(), bookID)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Get book success", http.StatusOK, c)
}

// Update implements BookHandler.
func (h *bookHandler) Update(c echo.Context) error {
	bookID := utils.ConvertString(c.Param("book-id"))

	result := <-h.bookUsecase.GetByBookID(c.Request().Context(), bookID)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	expend := result.Data.(models.Book)
	if expend == (models.Book{}) {
		return utils.ResponseError(httperror.NotFound(httperror.NotFoundErrorMessage), c)
	}

	data := new(models.BookAdd)

	if err := c.Bind(data); err != nil {
		utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	expend = data.ToBook(expend)

	result = <-h.bookUsecase.Update(c.Request().Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Update book success", http.StatusOK, c)
}
