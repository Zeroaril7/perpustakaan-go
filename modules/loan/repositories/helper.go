package repositories

import (
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/models"
	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"gorm.io/gorm"
)

func buildFilterQuery(db *gorm.DB, f models.LoanBookFilter) *gorm.DB {
	if f.User != "" {
		db = db.Where("user = ?", f.User)
	}

	if f.Status != "" {
		db = db.Where("status = ?", f.Status)
	}

	if f.StartDate != "" && f.EndDate != "" {
		switch f.LoanTypeDate {
		case constant.LoanStartDate:
			db = db.Where("loan_start_date BETWEEN ? AND ?", f.StartDate, f.EndDate)
		case constant.LoanEndDate:
			db = db.Where("loan_end_date BETWEEN ? AND ?", f.StartDate, f.EndDate)
		}
	}

	return db
}
