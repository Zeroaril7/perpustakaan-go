package models

import "github.com/Zeroaril7/perpustakaan-go/pkg/utils"

type BookFilter struct {
	Author          []string `json:"author" query:"author"`
	Publisher       []string `json:"publisher" query:"publisher"`
	PublicationYear string   `json:"publication_year" query:"publication_year"`
	utils.PaginationRequest
}
