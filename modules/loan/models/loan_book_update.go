package models

type LoanBookUpdate struct {
	BookID string `json:"book_id" validate:"required"`
	Status string `json:"status"`
}
