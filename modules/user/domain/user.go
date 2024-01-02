package domain

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/user/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type UserRepository interface {
	Add(ctx context.Context, data models.User) (models.User, error)
	Delete(ctx context.Context, username string) error
	Get(ctx context.Context, filter models.UserFilter) ([]models.User, int64, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	Update(ctx context.Context, data models.User) (models.User, error)
}

type UserUsecase interface {
	Add(ctx context.Context, data models.User) <-chan utils.Result
	Delete(ctx context.Context, username string) <-chan utils.Result
	Get(ctx context.Context, filter models.UserFilter) <-chan utils.Result
	GetByUsername(ctx context.Context, username string) <-chan utils.Result
	Update(ctx context.Context, data models.User) <-chan utils.Result
}
