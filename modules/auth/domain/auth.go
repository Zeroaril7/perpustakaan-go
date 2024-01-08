package domain

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/auth/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

type AuthUsecase interface {
	AuthWithPassword(ctx context.Context, authReq models.LoginAuth) <-chan utils.Result
}
