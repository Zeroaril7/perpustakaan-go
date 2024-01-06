package models

type LoanBook struct {
	ID            int64  `json:"id" gorm:"primaryKey"`
	LoanID        string `json:"loan_id"`
	BookID        string `json:"book_id"`
	Title         string `json:"title"`
	Username      string `json:"username"`
	LoanStartDate string `json:"loan_start_date"`
	LoanEndDate   string `json:"loan_end_date"`
	Status        string `json:"status"`
}

func (LoanBook) TableName() string {
	return "loan_book"
}
