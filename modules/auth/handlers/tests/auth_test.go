package tests

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/handlers"
	"github.com/Zeroaril7/perpustakaan-go/modules/auth/usecases"
	userDomain "github.com/Zeroaril7/perpustakaan-go/modules/user/domain"
	userRepo "github.com/Zeroaril7/perpustakaan-go/modules/user/repositories"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/Zeroaril7/perpustakaan-go/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	authEndpoint                    = "/auth"
	loginEndpoint                   = "/login"
	invalidUserAuthBodyFilePath     = "test_data/auth_invalid_username_login_body_req.json"
	invalidPasswordAuthBodyFilePath = "test_data/auth_invalid_password_login_body_req.json"
	authBodyFilePath                = "test_data/auth_login_body_req.json"
	authBodyInvalidFilePath         = "test_data/auth_login_body_invalid_req.json"
	authBodyEmptyFilePath           = "test_data/auth_login_body_empty_req.json"
	privateKeyPath                  = "test_data/private.pem"
	publicKeyPath                   = "test_data/public.pem"
	userRows                        = []string{"username", "password", "role"}
	userResult                      = []driver.Value{"test", utils.HashPassword("test"), "ADMIN"}
)

type Suite struct {
	suite.Suite
	e              *echo.Echo
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	authUsecase    domain.AuthUsecase
	userRepository userDomain.UserRepository
	authHandler    handlers.AuthHandler
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	s.e = echo.New()
	s.e.Validator = validator.NewCustomValidator()
	db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	s.Require().NoError(err)

	s.userRepository = userRepo.NewUserRepository(s.DB)

	s.authUsecase = usecases.NewAuthUsecase(s.userRepository)
	s.authHandler = handlers.NewAuthHandler(s.e, s.authUsecase)

	config.LoadConfig()
}

func (s *Suite) TearDownSuite() {
	db, err := s.DB.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *Suite) TestLogin() {
	var tests = []struct {
		name            string
		bindErr         bool
		validatorErr    bool
		authUsernameErr bool
		authPasswordErr bool
		jwtErr          bool
		sqlErr          error
		expectedStatus  int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: gorm.ErrRecordNotFound, expectedStatus: http.StatusUnauthorized},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "auth username error", authUsernameErr: true, expectedStatus: http.StatusUnauthorized},
		{name: "auth password error", authPasswordErr: true, expectedStatus: http.StatusUnauthorized},
		{name: "jwt error", jwtErr: true, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		if !tt.jwtErr {
			privateKeyFile, _ := os.Open(privateKeyPath)
			privateKey, _ := io.ReadAll(privateKeyFile)

			publicKeyFile, _ := os.Open(publicKeyPath)
			publicKey, _ := io.ReadAll(publicKeyFile)

			config.Config().PrivateKey = string(privateKey)
			config.Config().PublicKey = string(publicKey)
		} else {
			config.Config().PrivateKey = ""
			config.Config().PublicKey = ""
		}

		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = authBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = authBodyEmptyFilePath
		} else if tt.authUsernameErr {
			bodyFilepath = invalidUserAuthBodyFilePath
		} else if tt.authPasswordErr {
			bodyFilepath = invalidPasswordAuthBodyFilePath
		} else {
			bodyFilepath = authBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPost, authEndpoint+loginEndpoint, jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetPath(authEndpoint + loginEndpoint)

		if tt.sqlErr != nil && !tt.bindErr && !tt.validatorErr && !tt.authPasswordErr && !tt.authUsernameErr && !tt.jwtErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if !tt.bindErr && !tt.validatorErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(userRows).AddRow(userResult...))
		}

		err = s.authHandler.Login(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)

	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
