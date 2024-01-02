package usecases

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type bookUsecase struct {
	bookRepository domain.BookRepository
}

// Add implements domain.BookUsecase.
func (u *bookUsecase) Add(ctx context.Context, data models.Book) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.bookRepository.Add(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// Delete implements domain.BookUsecase.
func (u *bookUsecase) Delete(ctx context.Context, register_id string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		err := u.bookRepository.Delete(ctx, register_id)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{}
	}()

	return output
}

// Get implements domain.BookUsecase.
func (u *bookUsecase) Get(ctx context.Context, filter models.BookFilter) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, total, err := u.bookRepository.Get(ctx, filter)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result, Total: total}
	}()

	return output
}

// GetLast implements domain.BookUsecase.
func (u *bookUsecase) GetLast(ctx context.Context, genre string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.bookRepository.GetLast(ctx, genre)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// GetByRegisterID implements domain.BookUsecase.
func (u *bookUsecase) GetByRegisterID(ctx context.Context, register_id string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.bookRepository.GetByRegisterID(ctx, register_id)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// Update implements domain.BookUsecase.
func (u *bookUsecase) Update(ctx context.Context, data models.Book) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.bookRepository.Update(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

func NewBookUsecase(bookRepository domain.BookRepository) domain.BookUsecase {
	return &bookUsecase{bookRepository: bookRepository}
}
