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
	bookDomain "github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	bookRepo "github.com/Zeroaril7/perpustakaan-go/modules/book/repositories"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/handlers"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/repositories"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/usecases"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	loanBookEndpoint            = "/loan-book"
	loanBookBodyFilePath        = "test_data/loan_book_body_req.json"
	loanBookBodyInvalidFilePath = "test_data/loan_book_body_invalid_req.json"
	loanBookBodyEmptyFilePath   = "test_data/loan_book_body_empty_req.json"
	loanBookRows                = []string{"id", "loan_id", "book_id", "title", "username", "loan_start_date", "loan_end_date", "status"}
	loanBookResult              = []driver.Value{1, "LOAN-TEST-0001", "TEST-DRAMA-0001", testStr, testStr, dateStr, dateStr, constant.LoanBorrowedStatus}
	emptyLoanBookResult         = []driver.Value{0, "", "", "", "", "", "", ""}
	bookRows                    = []string{"id", "book_id", "title", "genre", "author", "publisher", "publication_year", "status", "timestamp"}
	bookResult                  = []driver.Value{1, "TEST-DRAMA-0001", testStr, testStr, testStr, testStr, dateStr, constant.AvailableStatus, dateStr}
	emptyBookResult             = []driver.Value{0, "", "", "", "", "", "", "", ""}
	testStr                     = "test"
	dateStr                     = "2024-01-01"
)

type Suite struct {
	suite.Suite
	e                  *echo.Echo
	DB                 *gorm.DB
	mock               sqlmock.Sqlmock
	bookRepository     bookDomain.BookRepository
	loanBookRepository domain.LoanBookRepository
	loanBookUsecase    domain.LoanBookUsecase
	loanBookHandler    handlers.LoanBookHandler
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

	s.bookRepository = bookRepo.NewBookRepository(s.DB)
	s.loanBookRepository = repositories.NewLoanBookRepository(s.DB)
	s.loanBookUsecase = usecases.NewLoanBookUsecase(s.loanBookRepository, s.bookRepository)
	s.loanBookHandler = handlers.NewLoanBookHandler(s.e, s.loanBookUsecase)
}

func (s *Suite) TearDownSuite() {
	db, err := s.DB.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *Suite) TestAddLoanBook() {
	var tests = []struct {
		name           string
		bindErr        bool
		validatorErr   bool
		notFound       bool
		sqlGetDataErr  error
		sqlGetLastErr  error
		sqlErr         error
		sqlUpdateErr   error
		expectedStatus int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql get last error", sqlGetLastErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql get data error", sqlGetDataErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "not found", notFound: true, expectedStatus: http.StatusNotFound},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql update error", sqlUpdateErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = loanBookBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = loanBookBodyEmptyFilePath
		} else {
			bodyFilepath = loanBookBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPost, loanBookEndpoint, jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)

		c.SetPath(loanBookEndpoint)

		if tt.sqlErr == nil && tt.sqlGetLastErr != nil && tt.sqlGetDataErr == nil && tt.sqlUpdateErr == nil && !tt.bindErr && !tt.validatorErr && !tt.notFound {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetLastErr)
		} else if !tt.bindErr && !tt.validatorErr && tt.sqlGetDataErr != nil && tt.sqlErr == nil && tt.sqlUpdateErr == nil && tt.sqlGetLastErr == nil && !tt.notFound {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetDataErr)
		} else if tt.sqlErr == nil && tt.sqlUpdateErr == nil && !tt.bindErr && !tt.validatorErr && tt.sqlGetDataErr == nil && tt.sqlGetLastErr == nil && tt.notFound {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(emptyBookResult...))
		} else if !tt.bindErr && !tt.validatorErr && tt.sqlGetDataErr == nil && tt.sqlUpdateErr == nil && tt.sqlErr != nil && tt.sqlGetLastErr == nil && !tt.notFound {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetDataErr == nil && tt.sqlUpdateErr != nil && tt.sqlErr == nil && tt.sqlGetLastErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlUpdateErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetDataErr == nil && tt.sqlUpdateErr == nil && tt.sqlErr == nil && tt.sqlGetLastErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}

		err = s.loanBookHandler.Add(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestDeleteLoanBook() {
	tests := []struct {
		name           string
		expectedStatus int
		bindErr        bool
		validatorErr   bool
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
		req := httptest.NewRequest(http.MethodDelete, loanBookEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(loanBookEndpoint + "/:loan-id")
		c.SetParamNames("loan-id")
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

		err := s.loanBookHandler.Delete(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetLoanBook() {
	tests := []struct {
		name           string
		LoanTypeDate   string
		bindErr        bool
		totalErr       bool
		sqlErr         error
		expectedStatus int
	}{
		{name: "success", LoanTypeDate: constant.LoanStartDate, expectedStatus: http.StatusOK},
		{name: "success", LoanTypeDate: constant.LoanEndDate, expectedStatus: http.StatusOK},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "total error", totalErr: true, sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("user", testStr)
		q.Set("status", constant.LoanReturnedStatus)
		q.Set("start_date", dateStr)
		q.Set("end_date", dateStr)

		switch tt.LoanTypeDate {
		case constant.LoanStartDate:
			q.Set("loan_type_date", constant.LoanStartDate)
		case constant.LoanEndDate:
			q.Set("loan_type_date", constant.LoanEndDate)
		}

		if tt.bindErr {
			q.Set("per_page", "a")
		} else {
			q.Set("per_page", "10")
		}

		req := httptest.NewRequest(http.MethodGet, loanBookEndpoint+"?"+q.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(loanBookEndpoint)

		if tt.sqlErr != nil && !tt.bindErr && tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if tt.sqlErr != nil && !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else if !tt.bindErr && !tt.totalErr {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
		}

		err := s.loanBookHandler.Get(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestGetByLoanID() {
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
		req := httptest.NewRequest(http.MethodGet, loanBookEndpoint+"/test", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.SetPath(loanBookEndpoint + "/:loan-id")
		c.SetParamNames("loan-id")
		c.SetParamValues(testStr)

		if tt.sqlErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlErr)
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
		}

		err := s.loanBookHandler.GetByLoanID(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func (s *Suite) TestUpdateLoanBook() {
	var tests = []struct {
		name             string
		bindErr          bool
		validatorErr     bool
		notFound         bool
		sqlGetBookIDErr  error
		sqlGetLoanIDErr  error
		sqlErr           error
		sqlUpdateBookErr error
		expectedStatus   int
	}{
		{name: "success", expectedStatus: http.StatusOK},
		{name: "success", expectedStatus: http.StatusOK},
		{name: "sql get loan id error", sqlGetLoanIDErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "not found loan id", sqlGetLoanIDErr: sql.ErrNoRows, notFound: true, expectedStatus: http.StatusNotFound},
		{name: "bind error", bindErr: true, expectedStatus: http.StatusBadRequest},
		{name: "validator error", validatorErr: true, expectedStatus: http.StatusBadRequest},
		{name: "sql get book id data error", sqlGetBookIDErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "not found book id", sqlGetBookIDErr: sql.ErrNoRows, notFound: true, expectedStatus: http.StatusNotFound},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql error", sqlErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
		{name: "sql update book error", sqlUpdateBookErr: sql.ErrNoRows, expectedStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		var bodyFilepath string
		if tt.bindErr {
			bodyFilepath = loanBookBodyInvalidFilePath
		} else if tt.validatorErr {
			bodyFilepath = loanBookBodyEmptyFilePath
		} else {
			bodyFilepath = loanBookBodyFilePath
		}

		jsonFile, err := os.Open(bodyFilepath)
		s.Require().NoError(err)
		defer jsonFile.Close()

		req := httptest.NewRequest(http.MethodPut, loanBookEndpoint+"/test", jsonFile)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)

		c.SetPath(loanBookEndpoint + "/:loan-id")
		c.SetParamNames("loan-id")
		c.SetParamValues(testStr)

		if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetBookIDErr == nil && tt.sqlUpdateBookErr == nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetBookIDErr == nil && tt.sqlUpdateBookErr == nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetLoanIDErr)
		} else if !tt.bindErr && !tt.validatorErr && tt.notFound && tt.sqlGetBookIDErr == nil && tt.sqlUpdateBookErr == nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr != nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(emptyLoanBookResult...))
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetBookIDErr != nil && tt.sqlUpdateBookErr == nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnError(tt.sqlGetBookIDErr)
		} else if !tt.bindErr && !tt.validatorErr && tt.notFound && tt.sqlGetBookIDErr != nil && tt.sqlUpdateBookErr == nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(emptyBookResult...))
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetBookIDErr == nil && tt.sqlUpdateBookErr == nil && tt.sqlErr != nil && tt.sqlGetLoanIDErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlErr)
			s.mock.ExpectRollback()
		} else if !tt.bindErr && !tt.validatorErr && !tt.notFound && tt.sqlGetBookIDErr == nil && tt.sqlUpdateBookErr != nil && tt.sqlErr == nil && tt.sqlGetLoanIDErr == nil {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(bookRows).AddRow(bookResult...))
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
			s.mock.ExpectBegin()
			s.mock.ExpectExec("").WithArgs().WillReturnError(tt.sqlUpdateBookErr)
			s.mock.ExpectRollback()
		} else {
			s.mock.ExpectQuery("").WithArgs().WillReturnRows(sqlmock.NewRows(loanBookRows).AddRow(loanBookResult...))
		}

		err = s.loanBookHandler.Update(c)
		s.Require().NoError(err)
		s.Require().Equal(tt.expectedStatus, rec.Code)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
