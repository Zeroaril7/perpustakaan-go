package models

import "github.com/Zeroaril7/perpustakaan-go/pkg/utils"

type LoanBookFilter struct {
	User         string `json:"user"`
	Status       string `json:"status"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	LoanTypeDate string `json:"loan_type_date"`
	utils.PaginationRequest
}
