package models

type LoanBookAdd struct {
	BookID        string `json:"book_id" validate:"required"`
	Title         string `json:"title"`
	Username      string `json:"username" validate:"required"`
	LoanStartDate string `json:"loan_start_date" validate:"required"`
	LoanEndDate   string `json:"loan_end_date" validate:"required"`
	Status        string `json:"status"`
}
