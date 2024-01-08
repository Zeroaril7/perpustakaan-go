package models

type Book struct {
	ID              int64  `json:"id" gorm:"primaryKey"`
	BookID          string `json:"book_id"`
	Title           string `json:"title"`
	Genre           string `json:"genre"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher"`
	PublicationYear string `json:"publication_year"`
	Status          string `json:"status"`
	Timestamp       string `json:"timestamp"`
}

func (Book) TableName() string {
	return "book"
}
