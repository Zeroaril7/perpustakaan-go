package httperror

import "net/http"

type CommonErrorData struct {
	Code         int    `json:"code"`
	ResponseCode int    `json:"responseCode,omitempty"`
	Message      string `json:"message"`
}

type ErrorString struct {
	code    int
	message string
}

type (
	BadRequestData struct {
		ErrorString
	}

	UnauthorizedData struct {
		ErrorString
	}

	ForbiddenErrorData struct {
		ErrorString
	}

	NotFoundData struct {
		ErrorString
	}

	ConflictData struct {
		ErrorString
	}

	InternalServerErrorData struct {
		ErrorString
	}
)

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Message() string {
	return e.message
}

func NewBadRequest(msg string) BadRequestData {
	err := BadRequestData{}
	if msg != "" {
		err.message = msg
	} else {
		err.message = "Bad Request"
	}

	err.code = http.StatusBadRequest

	return err
}

func NewUnauthorized(msg string) UnauthorizedData {
	err := UnauthorizedData{}

	if msg != "" {
		err.message = msg
	} else {
		err.message = "Unauthorized"
	}

	err.code = http.StatusUnauthorized

	return err
}

func NewForbiddenError(msg string) ForbiddenErrorData {
	err := ForbiddenErrorData{}

	if msg != "" {
		err.message = msg
	} else {
		err.message = "Forbidden"
	}

	err.code = http.StatusForbidden

	return err
}

func NewNotFound(msg string) NotFoundData {
	err := NotFoundData{}

	if msg != "" {
		err.message = msg
	} else {
		err.message = "Not Found"
	}

	err.code = http.StatusNotFound

	return err
}

func NewConflict(msg string) ConflictData {
	err := ConflictData{}

	if msg != "" {
		err.message = msg
	} else {
		err.message = "Conflict"
	}

	err.code = http.StatusConflict

	return err
}

func NewInternalServerError(msg string) InternalServerErrorData {
	err := InternalServerErrorData{}

	if msg != "" {
		err.message = msg
	} else {
		err.message = "Internal Server Error"
	}

	err.code = http.StatusInternalServerError

	return err
}

func BadRequest(msg string) error {
	return NewBadRequest(msg)
}

func Unauthorized(msg string) error {
	return NewUnauthorized(msg)
}

func Forbidden(msg string) error {
	return NewForbiddenError(msg)
}

func NotFound(msg string) error {
	return NewNotFound(msg)
}

func Conflict(msg string) error {
	return NewConflict(msg)
}

func InternalServerError(msg string) error {
	return NewInternalServerError(msg)
}
