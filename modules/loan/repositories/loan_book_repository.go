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
	panic("unimplemented")
}

// GetByLoanID implements domain.LoanBookRepository.
func (r *loanBookRepository) GetByLoanID(ctx context.Context, loan_id string) (result models.LoanBook, err error) {
	err = r.db.WithContext(ctx).Find(&result).Where("loan_id = ?", loan_id).Error
	return
}

// GetLast implements domain.LoanBookRepository.
func (r *loanBookRepository) GetLast(ctx context.Context, user string) (result models.LoanBook, err error) {
	db := r.db.WithContext(ctx)

	if err = db.Select("loan_id").Last(&result, "user = ?", user).Error; err != nil {
		return result, nil
	}

	return
}

// Update implements domain.LoanBookRepository.
func (r *loanBookRepository) Update(ctx context.Context, data models.LoanBook) (models.LoanBook, error) {
	panic("unimplemented")
}

func NewLoanBookRepository(db *gorm.DB) domain.LoanBookRepository {
	return &loanBookRepository{db: db}
}
