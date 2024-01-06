package usecases

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/models"
	userDomain "github.com/Zeroaril7/perpustakaan-go/modules/user/domain"
	userModel "github.com/Zeroaril7/perpustakaan-go/modules/user/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/Zeroaril7/perpustakaan-go/pkg/sdk/jwtrsa"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"gorm.io/gorm"
)

type authUsecase struct {
	userRepository userDomain.UserRepository
}

// AuthWithPassword implements domain.AuthUsecase.
func (u *authUsecase) AuthWithPassword(ctx context.Context, authReq models.LoginAuth) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		user, err := u.userRepository.GetByUsername(ctx, authReq.Username)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				output <- utils.Result{Error: httperror.NewUnauthorized(httperror.InvalidLoginMsg)}
				return
			} else {
				output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
				return
			}
		} else if user.Username != authReq.Username {
			output <- utils.Result{Error: httperror.NewUnauthorized(httperror.InvalidLoginMsg)}
			return
		}

		if !utils.CheckPasswordHash(authReq.Password, user.Password) {
			output <- utils.Result{Error: httperror.NewUnauthorized(httperror.InvalidLoginMsg)}
			return
		}

		authResponse, err := u.createAuthResponse(ctx, user)
		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: authResponse}

	}()

	return output
}

// verifyPassword implements domain.AuthUsecase.
func (u *authUsecase) verifyPassword(password string, hash string) bool {
	err := utils.CheckPasswordHash(password, hash)
	return err
}

func (u *authUsecase) createAccessToken(ctx context.Context, user userModel.User, accessTokenTTL time.Duration) (accessToken string, err error) {
	userIDStr := strconv.Itoa(int(user.ID))

	accessTokenClaims := models.AccessTokenClaims{
		Aud:      userIDStr,
		Username: user.Username,
		Role:     user.Role,
	}

	claims := generateMapClaims(accessTokenClaims)

	inputJWT := jwtrsa.GenerateInputJWT{
		PrivateKey: config.Config().PrivateKey,
		Claims:     claims,
		TimeExpire: accessTokenTTL,
	}

	accessToken, _, err = jwtrsa.GenerateJWT(inputJWT)
	return
}

func (u *authUsecase) createAuthResponse(ctx context.Context, user userModel.User) (token models.AuthResponse, err error) {
	accessTokenTTL := 2 * time.Hour

	accessToken, err := u.createAccessToken(ctx, user, accessTokenTTL)

	token.TokenType = "Bearer"
	token.AccessToken = accessToken
	token.ExpiresIn = int(accessTokenTTL.Minutes())

	return
}

func generateMapClaims(claims interface{}) map[string]interface{} {
	types := reflect.TypeOf(claims)
	values := reflect.ValueOf(claims)

	result := make(map[string]interface{})

	for i := 0; i < types.NumField(); i++ {
		key := types.Field(i).Tag.Get("claim")
		value := values.Field(i).Interface()
		result[key] = value
	}

	return result
}

func NewAuthUsecase(userRepository userDomain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
	}
}
