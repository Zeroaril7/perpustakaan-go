package models

type LoanBookAdd struct {
	BookID        string `json:"book_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
	User          string `json:"user" validate:"required"`
	LoanStartDate string `json:"loan_start_date" validate:"required"`
	LoanEndDate   string `json:"loan_end_date" validate:"required"`
}
