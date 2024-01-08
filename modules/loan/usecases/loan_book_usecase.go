package usecases

import (
	"context"

	bookDomain "github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	bookModel "github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type loanBookUsecase struct {
	loanBookRepository domain.LoanBookRepository
	bookRepository     bookDomain.BookRepository
}

// Add implements domain.LoanBookUsecase.
func (u *loanBookUsecase) Add(ctx context.Context, data models.LoanBook) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		expend, err := u.bookRepository.GetByBookID(ctx, data.BookID)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		if expend == (bookModel.Book{}) {
			output <- utils.Result{Error: httperror.NotFound("Book ID not found")}
			return
		}

		data.Title = expend.Title

		result, err := u.loanBookRepository.Add(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		expend.Status = constant.NotAvailableStatus

		_, err = u.bookRepository.Update(ctx, expend)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// Delete implements domain.LoanBookUsecase.
func (u *loanBookUsecase) Delete(ctx context.Context, loan_id string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		err := u.loanBookRepository.Delete(ctx, loan_id)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{}
	}()

	return output
}

// Get implements domain.LoanBookUsecase.
func (u *loanBookUsecase) Get(ctx context.Context, filter models.LoanBookFilter) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, total, err := u.loanBookRepository.Get(ctx, filter)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result, Total: total}
	}()

	return output
}

// GetByLoanID implements domain.LoanBookUsecase.
func (u *loanBookUsecase) GetByLoanID(ctx context.Context, loan_id string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.loanBookRepository.GetByLoanID(ctx, loan_id)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// GetLast implements domain.LoanBookUsecase.
func (u *loanBookUsecase) GetLast(ctx context.Context, user string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.loanBookRepository.GetLast(ctx, user)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// Update implements domain.LoanBookUsecase.
func (u *loanBookUsecase) Update(ctx context.Context, data models.LoanBook) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		expend, err := u.bookRepository.GetByBookID(ctx, data.BookID)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		if expend == (bookModel.Book{}) {
			output <- utils.Result{Error: httperror.NotFound("Book ID not found")}
			return
		}

		result, err := u.loanBookRepository.Update(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		switch result.Status {
		case constant.LoanBorrowedStatus:
			expend.Status = constant.NotAvailableStatus
		case constant.LoanReturnedStatus:
			expend.Status = constant.AvailableStatus
		}

		_, err = u.bookRepository.Update(ctx, expend)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

func NewLoanBookUsecase(loanBokRepository domain.LoanBookRepository, bookRepository bookDomain.BookRepository) domain.LoanBookUsecase {
	return &loanBookUsecase{loanBookRepository: loanBokRepository, bookRepository: bookRepository}
}
