package models

import "github.com/Zeroaril7/perpustakaan-go/pkg/utils"

type LoanBookFilter struct {
	User         string `json:"user" query:"user"`
	Status       string `json:"status" query:"status"`
	StartDate    string `json:"start_date" query:"start_date"`
	EndDate      string `json:"end_date" query:"end_date"`
	LoanTypeDate string `json:"loan_type_date" query:"loan_type_date"`
	utils.PaginationRequest
}
