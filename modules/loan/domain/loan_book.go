package domain

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/loan/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type LoanBookRepository interface {
	Add(ctx context.Context, data models.LoanBook) (models.LoanBook, error)
	Get(ctx context.Context, filter models.LoanBookFilter) ([]models.LoanBook, int64, error)
	GetByLoanID(ctx context.Context, loan_id string) (models.LoanBook, error)
	GetLast(ctx context.Context, username string) (models.LoanBook, error)
	Update(ctx context.Context, data models.LoanBook) (models.LoanBook, error)
	Delete(ctx context.Context, loan_id string) error
}

type LoanBookUsecase interface {
	Get(ctx context.Context, filter models.LoanBookFilter) <-chan utils.Result
	GetLast(ctx context.Context, user string) <-chan utils.Result
	GetByLoanID(ctx context.Context, loan_id string) <-chan utils.Result
	Add(ctx context.Context, data models.LoanBook) <-chan utils.Result
	Update(ctx context.Context, data models.LoanBook) <-chan utils.Result
	Delete(ctx context.Context, loan_id string) <-chan utils.Result
}
