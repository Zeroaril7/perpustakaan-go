package repositories

import (
	"github.com/Zeroaril7/perpustakaan-go/modules/user/models"
	"gorm.io/gorm"
)

func buildFilterQuery(db *gorm.DB, f models.UserFilter) *gorm.DB {
	if f.Role != "" {
		db = db.Where("role = ?", f.Role)
	}

	return db
}
