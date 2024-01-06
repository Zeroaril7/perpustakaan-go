package handlers

import (
	"net/http"

	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/middlewares"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/labstack/echo/v4"
)

type LoanBookHandler interface {
	Add(c echo.Context) error
	Delete(c echo.Context) error
	Get(c echo.Context) error
	GetByLoanID(c echo.Context) error
	Update(c echo.Context) error
}

type loanBookHandler struct {
	loanBookUsecase domain.LoanBookUsecase
}

func NewLoanBookHandler(e *echo.Echo, loanBookUsecase domain.LoanBookUsecase) LoanBookHandler {
	handler := &loanBookHandler{loanBookUsecase: loanBookUsecase}

	group := e.Group("/loan-book", middlewares.VerifyJWTRSA(config.Config().PublicKey), middlewares.EchoSetCredential())
	group.POST("", handler.Add)
	group.DELETE("/:loan-id", handler.Delete)
	group.GET("", handler.Get)
	group.GET("/:loan-id", handler.GetByLoanID)
	group.PUT("/:loan-id", handler.Update)

	return handler
}

// Add implements LoanBookHandler.
func (h *loanBookHandler) Add(c echo.Context) error {
	data := new(models.LoanBookAdd)

	if err := c.Bind(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	result := <-h.loanBookUsecase.GetLast(c.Request().Context(), data.Username)

	expend := result.Data.(models.LoanBook)

	if expend == (models.LoanBook{}) {
		expend.Username = data.Username
	}

	expend = data.ToLoanBook(expend)
	expend.Status = constant.LoanBorrowedStatus

	result = <-h.loanBookUsecase.Add(c.Request().Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Add loan book success", http.StatusOK, c)
}

// Delete implements LoanBookHandler.
func (h *loanBookHandler) Delete(c echo.Context) error {
	loan_id := utils.ConvertString(c.Param("loan-id"))

	role := c.Get("role").(string)

	if role != constant.Admin && role != constant.SuperAdmin {
		return utils.ResponseError(httperror.Unauthorized(httperror.UnauthorizedErrorMessage), c)
	}

	result := <-h.loanBookUsecase.Delete(c.Request().Context(), loan_id)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(nil, "Delete loan book success", http.StatusOK, c)
}

// Get implements LoanBookHandler.
func (h *loanBookHandler) Get(c echo.Context) error {
	filter := new(models.LoanBookFilter)

	if err := c.Bind(filter); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if !filter.DisablePagination {
		filter.SetDefault()
	}

	result := <-h.loanBookUsecase.Get(c.Request().Context(), *filter)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.ResponseWithPagination(result.Data, "Get loan book success", http.StatusOK, result.Total, filter.GetPaginationRequest(), c)
}

// GetByLoanID implements LoanBookHandler.
func (h *loanBookHandler) GetByLoanID(c echo.Context) error {
	loan_id := utils.ConvertString(c.Param("loan-id"))

	result := <-h.loanBookUsecase.GetByLoanID(c.Request().Context(), loan_id)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Get loan book success", http.StatusOK, c)
}

// Update implements LoanBookHandler.
func (h *loanBookHandler) Update(c echo.Context) error {
	loan_id := utils.ConvertString(c.Param("loan-id"))

	result := <-h.loanBookUsecase.GetByLoanID(c.Request().Context(), loan_id)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	expend := result.Data.(models.LoanBook)

	if expend == (models.LoanBook{}) {
		return utils.ResponseError(httperror.NotFound(httperror.NotFoundErrorMessage), c)
	}

	data := new(models.LoanBookAdd)

	if err := c.Bind(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if err := c.Validate(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	expend = data.ToLoanBook(expend)

	result = <-h.loanBookUsecase.Update(c.Request().Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Update loan book success", http.StatusOK, c)
}
