package repositories

import (
	"github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"gorm.io/gorm"
)

func buildFilterQuery(db *gorm.DB, f models.BookFilter) *gorm.DB {
	if len(f.Author) > 0 {
		db = db.Where("author IN ?", f.Author)
	}

	if len(f.Publisher) > 0 {
		db = db.Where("publisher IN ?", f.Publisher)
	}

	if f.PublicationYear != "" {
		db = db.Where("publication_year = ?", f.PublicationYear)
	}

	return db
}
