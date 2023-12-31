package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zeroaril7/perpustakaan-go/config"
	authDomain "github.com/Zeroaril7/perpustakaan-go/modules/auth/domain"
	authHandler "github.com/Zeroaril7/perpustakaan-go/modules/auth/handlers"
	authUsecase "github.com/Zeroaril7/perpustakaan-go/modules/auth/usecases"
	bookDomain "github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	bookHandler "github.com/Zeroaril7/perpustakaan-go/modules/book/handlers"
	bookRepository "github.com/Zeroaril7/perpustakaan-go/modules/book/repositories"
	bookUsecase "github.com/Zeroaril7/perpustakaan-go/modules/book/usecases"
	loanBookDomain "github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
	loanBookHandler "github.com/Zeroaril7/perpustakaan-go/modules/loan/handlers"
	loanBookRepository "github.com/Zeroaril7/perpustakaan-go/modules/loan/repositories"
	loanBookUsecase "github.com/Zeroaril7/perpustakaan-go/modules/loan/usecases"
	userDomain "github.com/Zeroaril7/perpustakaan-go/modules/user/domain"
	userHandler "github.com/Zeroaril7/perpustakaan-go/modules/user/handlers"
	userRepository "github.com/Zeroaril7/perpustakaan-go/modules/user/repositories"
	userUsecase "github.com/Zeroaril7/perpustakaan-go/modules/user/usecases"
	mysqlgorm "github.com/Zeroaril7/perpustakaan-go/pkg/databases"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
	"github.com/Zeroaril7/perpustakaan-go/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type repositories struct {
	bookRepository    bookDomain.BookRepository
	userRepository    userDomain.UserRepository
	loanBokRepository loanBookDomain.LoanBookRepository
}

type usecase struct {
	bookUsecase     bookDomain.BookUsecase
	userUsecase     userDomain.UserUsecase
	authUsecase     authDomain.AuthUsecase
	loanBookUsecase loanBookDomain.LoanBookUsecase
}

type packages struct {
	repositories repositories
	usecase      usecase
}

var pkg packages

func setPackages() {
	// repository
	pkg.repositories.bookRepository = bookRepository.NewBookRepository(mysqlgorm.DBConnect.Connection)
	pkg.repositories.userRepository = userRepository.NewUserRepository(mysqlgorm.DBConnect.Connection)
	pkg.repositories.loanBokRepository = loanBookRepository.NewLoanBookRepository(mysqlgorm.DBConnect.Connection)

	// usecase
	pkg.usecase.bookUsecase = bookUsecase.NewBookUsecase(pkg.repositories.bookRepository)
	pkg.usecase.userUsecase = userUsecase.NewUserUsecase(pkg.repositories.userRepository)
	pkg.usecase.authUsecase = authUsecase.NewAuthUsecase(pkg.repositories.userRepository)
	pkg.usecase.loanBookUsecase = loanBookUsecase.NewLoanBookUsecase(pkg.repositories.loanBokRepository, pkg.repositories.bookRepository)

}

func setHttp(e *echo.Echo) {
	e.GET("/v1/health-check", func(c echo.Context) error {
		log.Default().Println("main", "This service is running properly")
		return utils.Response(nil, "This service is running properly", 200, c)
	})

	// Book
	bookHandler.NewBookHandler(e, pkg.usecase.bookUsecase)

	// User
	userHandler.NewUserHandler(e, pkg.usecase.userUsecase)

	// Auth
	authHandler.NewAuthHandler(e, pkg.usecase.authUsecase)

	// Loan
	loanBookHandler.NewLoanBookHandler(e, pkg.usecase.loanBookUsecase)

}

func main() {

	path, _ := os.Getwd()
	utils.LogDefault(path)

	mysqlgorm.InitConnection(config.Config().MySQLDSN())

	e := echo.New()

	e.Validator = validator.NewCustomValidator()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           `[ROUTE] ${time_rfc3339} | ${status} | ${latency_human} ${latency} | ${method} | ${uri}` + "\n",
		CustomTimeFormat: "2000-01-01 10:10:01.00000",
	}))

	e.Use(middleware.Recover())
	setPackages()
	setHttp(e)

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	listenerPort := fmt.Sprintf(":%s", config.Config().AppPort)
	e.Logger.Fatal(e.Start(listenerPort))

	server := &http.Server{
		Addr:         listenerPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Default().Println("main", fmt.Sprintf("Could not listen on %s: %v\n", config.Config().AppPort, err))
	}

}
