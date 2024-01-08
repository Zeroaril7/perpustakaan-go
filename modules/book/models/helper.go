package models

import (
	"fmt"
	"strings"

	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

func (m *BookAdd) ToBook(e Book) Book {
	e.BookID = generateBookID(e)
	e.Title = m.Title
	e.Author = m.Author
	e.Genre = m.Genre
	e.Publisher = m.Publisher
	e.PublicationYear = m.PublicationYear
	e.Status = constant.AvailableStatus
	e.Timestamp = utils.ConvertString(utils.GetLocalTime())

	return e
}

func generateBookID(data Book) string {

	var bookID string
	var number string

	if data.BookID == "" {
		number = utils.GetFourDigitsNumber("1")
		bookID = fmt.Sprintf("%s-%s-%v", constant.Institute, data.Genre, number)
	} else {
		lastBookData := strings.Split(data.BookID, "-")

		lastNumber := utils.ConvertInt(lastBookData[2]) + 1
		number = utils.GetFourDigitsNumber(utils.ConvertString(lastNumber))
		bookID = fmt.Sprintf("%s-%s-%v", constant.Institute, lastBookData[1], number)
	}

	return bookID
}
