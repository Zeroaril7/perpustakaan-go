package models

import (
	"fmt"
	"strings"

	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

func (m *LoanBookAdd) ToLoanBook(e LoanBook) LoanBook {
	e.LoanID = generateLoanID(e)
	e.BookID = m.BookID
	e.LoanStartDate = m.LoanStartDate
	e.LoanEndDate = m.LoanEndDate
	e.Username = m.Username
	e.Status = m.Status

	return e
}

func generateLoanID(data LoanBook) string {

	var loanID string
	var number string

	if data.LoanID == "" {
		number = utils.GetFourDigitsNumber("1")
		loanID = fmt.Sprintf("%s-%s-%v", constant.Loan, data.Username, number)
	} else {
		lastBookData := strings.Split(data.LoanID, "-")

		lastNumber := utils.ConvertInt(lastBookData[2]) + 1
		number = utils.GetFourDigitsNumber(utils.ConvertString(lastNumber))
		loanID = fmt.Sprintf("%s-%s-%v", constant.Loan, data.Username, number)
	}

	return loanID
}
