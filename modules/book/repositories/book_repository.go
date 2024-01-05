package repositories

import (
	"context"

	"github.com/Zeroaril7/perpustakaan-go/modules/book/domain"
	"github.com/Zeroaril7/perpustakaan-go/modules/book/models"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

// Add implements domain.BookRepository.
func (r *bookRepository) Add(ctx context.Context, data models.Book) (result models.Book, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	return data, err
}

// Delete implements domain.BookRepository.
func (r *bookRepository) Delete(ctx context.Context, book_id string) error {
	return r.db.WithContext(ctx).Where("book_id = ?", book_id).Delete(&models.Book{}).Error
}

// GetLast implements domain.BookRepository.
func (r *bookRepository) GetLast(ctx context.Context, genre string) (result models.Book, err error) {
	db := r.db.WithContext(ctx)

	if err = db.Select("book_id").Last(&result, "genre = ?", genre).Error; err != nil {
		return result, nil
	}

	return
}

// GetByBookID implements domain.BookRepository.
func (r *bookRepository) GetByBookID(ctx context.Context, book_id string) (result models.Book, err error) {
	err = r.db.WithContext(ctx).Where("book_id = ?", book_id).First(&result).Error
	return
}

// Get implements domain.BookRepository.
func (r *bookRepository) Get(ctx context.Context, filter models.BookFilter) (result []models.Book, total int64, err error) {
	db := r.db.WithContext(ctx)
	db = buildFilterQuery(db, filter)

	if err = db.Model(&models.Book{}).Count(&total).Error; err != nil {
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

// Update implements domain.BookRepository.
func (r *bookRepository) Update(ctx context.Context, data models.Book) (result models.Book, err error) {
	err = r.db.WithContext(ctx).Save(&data).Error
	return data, err
}

func NewBookRepository(db *gorm.DB) domain.BookRepository {
	return &bookRepository{db: db}
}
