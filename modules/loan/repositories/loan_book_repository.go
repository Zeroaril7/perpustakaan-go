package repositories

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/loan/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/loan/models"
	"gorm.io/gorm"
)

type loanBookRepository struct {
	db *gorm.DB
}

// Add implements domain.LoanBookRepository.
func (r *loanBookRepository) Add(ctx context.Context, data models.LoanBook) (result models.LoanBook, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	return data, err
}

// Delete implements domain.LoanBookRepository.
func (r *loanBookRepository) Delete(ctx context.Context, loan_id string) error {
	return r.db.WithContext(ctx).Where("loan_id = ?", loan_id).Delete(&models.LoanBook{}).Error
}

// Get implements domain.LoanBookRepository.
func (r *loanBookRepository) Get(ctx context.Context, filter models.LoanBookFilter) (result []models.LoanBook, total int64, err error) {
	db := r.db.WithContext(ctx)
	db = buildFilterQuery(db, filter)

	if err = db.Model(&models.LoanBook{}).Count(&total).Error; err != nil {
		return
	}

	if !filter.DisablePagination {
		db = db.Offset(int(filter.GetOffset())).Limit(int(filter.GetLimit()))
	}

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

// GetByLoanID implements domain.LoanBookRepository.
func (r *loanBookRepository) GetByLoanID(ctx context.Context, loan_id string) (result models.LoanBook, err error) {
	err = r.db.WithContext(ctx).Where("loan_id = ?", loan_id).First(&result).Error
	return
}

// GetLast implements domain.LoanBookRepository.
func (r *loanBookRepository) GetLast(ctx context.Context, username string) (result models.LoanBook, err error) {
	err = r.db.WithContext(ctx).Select("loan_id", "username").Last(&result, "username = ?", username).Error
	return
}

// Update implements domain.LoanBookRepository.
func (r *loanBookRepository) Update(ctx context.Context, data models.LoanBook) (result models.LoanBook, err error) {
	err = r.db.WithContext(ctx).Save(&data).Error
	return data, err
}

func NewLoanBookRepository(db *gorm.DB) domain.LoanBookRepository {
	return &loanBookRepository{db: db}
}
