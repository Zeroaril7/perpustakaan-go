package domain

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type BookRepository interface {
	Add(ctx context.Context, data models.Book) (models.Book, error)
	Get(ctx context.Context, filter models.BookFilter) ([]models.Book, int64, error)
	GetByRegisterID(ctx context.Context, register_id string) (models.Book, error)
	GetLast(ctx context.Context, genre string) (models.Book, error)
	Update(ctx context.Context, data models.Book) (models.Book, error)
	Delete(ctx context.Context, register_id string) error
}

type BookUsecase interface {
	Get(ctx context.Context, filter models.BookFilter) <-chan utils.Result
	GetLast(ctx context.Context, genre string) <-chan utils.Result
	GetByRegisterID(ctx context.Context, register_id string) <-chan utils.Result
	Add(ctx context.Context, data models.Book) <-chan utils.Result
	Update(ctx context.Context, data models.Book) <-chan utils.Result
	Delete(ctx context.Context, register_id string) <-chan utils.Result
}
