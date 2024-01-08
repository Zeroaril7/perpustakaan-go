package tests

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zeroaril7/perpustakaan-go/config"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/handlers"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/repositories"
	"github.com/Zeroaril7/perpustakaan-go/modules/user/usecases"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/Zeroaril7/perpustakaan-go/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	userEndpoint            = "/user"
	userRows                = []string{"id", "username", "password", "role"}
	userResult              = []driver.Value{1, "test", utils.HashPassword("test123"), "ADMIN"}
	emptyResult             = []driver.Value{0, "", "", ""}
	userBodyFilePath        = "test_data/user_body_req.json"
	userBodyInvalidFilePath = "test_data/user_body_invalid_req.json"
	userBodyEmptyFilePath   = "test_data/user_body_empty_req.json"
	testStr                 = "test"
)

type Suite struct {
	suite.Suite
	e              *echo.Echo
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	userRepository domain.UserRepository
	userUsecase    domain.UserUsecase
	userHandler    handlers.UserHandler
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

	s.userRepository = repositories.NewUserRepository(s.DB)
	s.userUsecase = usecases.NewUserUsecase(s.userRepository)
	s.userHandler = handlers.NewUserHandler(s.e, s.userUsecase)

	config.LoadConfig()
}

func (s *Suite) TearDownSuite() {
	db, err := s.DB.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *Suite) TestAddUser() {
	tests := []struct {
		name           string
		expectedStatus int
		sqlErr         error
		bindErr        bool
		validatorErr   bool
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = userBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = userBodyEmptyFilePath
		} else {
			bodyFilepath = userBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPost, userEndpoint, jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)

		c.SetPath(userEndpoint)

		if tt.sqlErr != nil && !tt.bindErr && !tt.validatorErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err = s.userHandler.Add(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)

	}
}

func (s *Suite) TestDeleteUser() {
	tests := []struct {
		name           string
		expectedStatus int
		sqlErr         error
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodDelete, userEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(userEndpoint + "/:username")
		c.SetParamNames("username")
		c.SetParamValues(testStr)

		if tt.sqlErr != nil {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err := s.userHandler.Delete(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetUser() {
	tests := []struct {
		name           string
		bindErr        bool
		totalErr       bool
		sqlErr         error
		expectedStatus int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "total error", totalErr: true, sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("role", constant.Karyawan)

		if tt.bindErr {
			q.Set("per_page", "a")
		} else {
			q.Set("per_page", "10")
		}

		req := httptest.NewRequest(http.MethodGet, userEndpoint+"?"+q.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(userEndpoint)

		if tt.sqlErr != nil && !tt.bindErr && tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if tt.sqlErr != nil && !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(userRows).AddRow(userResult...))
		}

		err := s.userHandler.Get(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetByUsername() {
	tests := []struct {
		name           string
		sqlErr         error
		expectedStatus int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, userEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(userEndpoint + "/:username")
		c.SetParamNames("username")
		c.SetParamValues(testStr)

		if tt.sqlErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(userRows).AddRow(userResult...))
		}

		err := s.userHandler.GetByUsername(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestUpdateUser() {
	tests := []struct {
		name           string
		expectedStatus int
		sqlErr         error
		sqlGetUserErr  error
		notFound       bool
		bindErr        bool
		validatorErr   bool
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql get user error", sqlGetUserErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "not found", notFound: true, expectedStatus: http.StatusNotFound},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = userBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = userBodyEmptyFilePath
		} else {
			bodyFilepath = userBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPut, userEndpoint+"/test", jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(userEndpoint + "/:username")
		c.SetParamNames("username")
		c.SetParamValues(testStr)

		if tt.notFound && !tt.bindErr && !tt.validatorErr && tt.sqlErr == nil && tt.sqlGetUserErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(userRows).AddRow(emptyResult...))
		} else if !tt.notFound && !tt.bindErr && !tt.validatorErr && tt.sqlErr == nil && tt.sqlGetUserErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetUserErr)
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(userRows).AddRow(userResult...))
		}

		if tt.sqlErr != nil && tt.sqlGetUserErr == nil && !tt.notFound && !tt.bindErr && !tt.validatorErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if tt.sqlErr == nil && tt.sqlGetUserErr == nil && !tt.notFound && !tt.bindErr && !tt.validatorErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err = s.userHandler.Update(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
