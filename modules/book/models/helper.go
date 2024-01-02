package models

import (
	"fmt"
	"strings"

	"github.com/Zeroaril7/perpustakaan-go/pkg/constant"
	"github.com/Zeroaril7/perpustakaan-go/pkg/utils"
)

func (m *BookAdd) ToBook(e Book) Book {
	e.RegisterID = generateReqisterID(e)
	e.Title = m.Title
	e.Author = m.Author
	e.Genre = m.Genre
	e.Publisher = m.Publisher
	e.PublicationYear = m.PublicationYear
	e.Status = constant.AvailableStatus
	e.Timestamp = utils.ConvertString(utils.GetLocalTime())

	return e
}

func generateReqisterID(data Book) string {

	var registerID string
	var number string

	if data.RegisterID == "" {
		number = utils.GetFourDigitsNumber("1")
		registerID = fmt.Sprintf("%s-%s-%v", constant.Institute, data.Genre, number)
	} else {
		lastRegisterData := strings.Split(data.RegisterID, "-")

		lastNumber := utils.ConvertInt(lastRegisterData[1]) + 1
		number = utils.GetFourDigitsNumber(utils.ConvertString(lastNumber))
		registerID = fmt.Sprintf("%s-%s-%v", constant.Institute, lastRegisterData[0], number)
	}

	return registerID
}
