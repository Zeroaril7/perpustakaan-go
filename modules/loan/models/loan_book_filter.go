package models

type LoanBookFilter struct {
	User         string `json:"user"`
	Status       string `json:"status"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	LoanTypeDate string `json:"loan_type_date"`
}
