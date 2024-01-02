package validator

import (
	"fmt"
	"strings"

	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/go-playground/validator"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	var errorMsg string
	if cv.validator.Struct(i) != nil {
		errs := cv.validator.Struct(i).(validator.ValidationErrors)
		for _, err := range errs {
			errorMsg += fmt.Sprintf("\"%s\": %s \n ", strings.ToLower(err.Field()), err.Tag())
		}

		return httperror.Conflict(errorMsg)
	}

	return nil
}

func NewCustomValidator() *CustomValidator {
	cv := &CustomValidator{validator: validator.New()}
	return cv
}
