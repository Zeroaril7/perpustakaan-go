package models

import "github.com/Zeroaril7/perpustakaan-go/pkg/utils"

type UserFilter struct {
	Role string `json:"role" query:"role"`
	utils.PaginationRequest
}
