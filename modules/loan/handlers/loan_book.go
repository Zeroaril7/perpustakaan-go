package handlers

import (
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
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

	return handler
}

// Add implements LoanBookHandler.
func (h *loanBookHandler) Add(c echo.Context) error {
	panic("unimplemented")
}

// Delete implements LoanBookHandler.
func (h *loanBookHandler) Delete(c echo.Context) error {
	panic("unimplemented")
}

// Get implements LoanBookHandler.
func (h *loanBookHandler) Get(c echo.Context) error {
	panic("unimplemented")
}

// GetByLoanID implements LoanBookHandler.
func (h *loanBookHandler) GetByLoanID(c echo.Context) error {
	panic("unimplemented")
}

// Update implements LoanBookHandler.
func (h *loanBookHandler) Update(c echo.Context) error {
	panic("unimplemented")
}
