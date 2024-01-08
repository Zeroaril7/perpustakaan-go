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
	"github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/handlers"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/repositories"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/usecases"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	bookEndpoint            = "/book"
	bookBodyFilePath        = "test_data/book_body_req.json"
	bookBodyInvalidFilePath = "test_data/book_body_invalid_req.json"
	bookBodyEmptyFilePath   = "test_data/book_body_empty_req.json"
	bookRows                = []string{"id", "book_id", "title", "genre", "author", "publisher", "publication_year", "status", "timestamp"}
	bookResult              = []driver.Value{1, "TEST-DRAMA-0001", testStr, testStr, testStr, testStr, dateStr, constant.AvailableStatus, dateStr}
	emptyResult             = []driver.Value{0, "", "", "", "", "", "", "", ""}
	testStr                 = "test"
	dateStr                 = "2024-01-01"
)

type Suite struct {
	suite.Suite
	e              *echo.Echo
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	bookRepository domain.BookRepository
	bookUsecase    domain.BookUsecase
	bookHandler    handlers.BookHandler
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

	s.bookRepository = repositories.NewBookRepository(s.DB)
	s.bookUsecase = usecases.NewBookUsecase(s.bookRepository)
	s.bookHandler = handlers.NewBookHandler(s.e, s.bookUsecase)
}

func (s *Suite) TearDownSuite() {
	db, err := s.DB.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *Suite) TestAddBook() {
	tests := []struct {
		name           string
		expectedStatus int
		bindErr        bool
		validatorErr   bool
		sqlErr         error
		sqlGetLastErr  error
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "sql get last error", sqlGetLastErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = bookBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = bookBodyEmptyFilePath
		} else {
			bodyFilepath = bookBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPost, bookEndpoint, jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)

		c.SetPath(bookEndpoint)

		if tt.sqlErr == nil && !tt.bindErr && !tt.validatorErr && tt.sqlGetLastErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetLastErr)
		} else if tt.sqlErr != nil && !tt.bindErr && !tt.validatorErr && tt.sqlGetLastErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(emptyResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr && tt.sqlGetLastErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(emptyResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err = s.bookHandler.Add(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestDeleteBook() {
	tests := []struct {
		name           string
		expectedStatus int
		roleErr        bool
		sqlErr         error
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "role error", roleErr: true, expectedStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodDelete, bookEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(bookEndpoint + "/:book-id")
		c.SetParamNames("book-id")
		c.SetParamValues(testStr)

		var role string
		if tt.roleErr {
			role = constant.Karyawan
		} else {
			role = constant.Admin
		}

		c.Set("role", role)

		if tt.sqlErr != nil && !tt.roleErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.roleErr {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err := s.bookHandler.Delete(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetBook() {
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
		q.Set("author", testStr)
		q.Set("publisher", testStr)
		q.Set("publication_year", dateStr)

		if tt.bindErr {
			q.Set("per_page", "a")
		} else {
			q.Set("per_page", "10")
		}

		req := httptest.NewRequest(http.MethodGet, bookEndpoint+"?"+q.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(bookEndpoint)

		if tt.sqlErr != nil && !tt.bindErr && tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if tt.sqlErr != nil && !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
		}

		err := s.bookHandler.Get(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetByBookID() {
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
		req := httptest.NewRequest(http.MethodGet, bookEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(bookEndpoint + "/:book-id")
		c.SetParamNames("book-id")
		c.SetParamValues(testStr)

		if tt.sqlErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
		}

		err := s.bookHandler.GetByBookID(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestUpdateBook() {
	var tests = []struct {
		name           string
		bindErr        bool
		validatorErr   bool
		notFound       bool
		sqlGetDataErr  error
		sqlErr         error
		expectedStatus int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql get data error", sqlGetDataErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql get data error", sqlGetDataErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "not found", notFound: true, expectedStatus: http.StatusNotFound},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = bookBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = bookBodyEmptyFilePath
		} else {
			bodyFilepath = bookBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPut, bookEndpoint+"/test", jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(bookEndpoint + "/:book-id")
		c.SetParamNames("book-id")
		c.SetParamValues(testStr)

		if tt.sqlGetDataErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetDataErr)
		} else if tt.notFound {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(emptyResult...))
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
		}

		if tt.sqlErr != nil && !tt.bindErr && !tt.validatorErr && tt.sqlGetDataErr == nil && !tt.notFound {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr && tt.sqlGetDataErr == nil && !tt.notFound {
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err = s.bookHandler.Update(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
