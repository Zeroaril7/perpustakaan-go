package models

type BookAdd struct {
	Title           string `json:"title" validate:"required"`
	Genre           string `json:"genre" validate:"required"`
	Author          string `json:"author" validate:"required"`
	Publisher       string `json:"publisher" validate:"required"`
	PublicationYear string `json:"publication_year"`
	Status			string `json:"status"`
}
