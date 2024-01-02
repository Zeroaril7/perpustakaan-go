package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zeroaril7/perpustakaan-go/pkg/httperror"
	"github.com/labstack/echo/v4"
)

type Result struct {
	Data  interface{}
	Error interface{}
	Total int64
}

// Pagination data structure
type PaginationRequest struct {
	Page              int64 `json:"page" query:"page"`
	PerPage           int64 `json:"per_page" query:"per_page"`
	DisablePagination bool  `json:"disable_pagination" query:"disable_pagination"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

type BaseWrapperModel struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int64     `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

func (q *PaginationRequest) GetOffset() int64 {
	if q.Page <= 1 {
		return 0
	}
	return (q.Page - 1) * q.PerPage
}

func (q *PaginationRequest) GetLimit() int64 {
	return q.PerPage
}

func (q *PaginationRequest) SetDefault() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PerPage == 0 {
		q.PerPage = 10
	}
}

func (q *PaginationRequest) SetPaginationResponse(page int64, perPage int64) PaginationRequest {
	return PaginationRequest{
		Page:    page,
		PerPage: perPage,
	}
}

func (q *PaginationRequest) GetPaginationRequest() PaginationRequest {
	return PaginationRequest{
		Page:              q.Page,
		PerPage:           q.PerPage,
		DisablePagination: q.DisablePagination,
	}
}

func Response(data interface{}, message string, code int, c echo.Context) error {
	success := false
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Request().Method,
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: c.Request().ContentLength,
		Ip:            c.RealIP(),
	}
	byteMeta, _ := json.Marshal(meta)
	LogDefault(string(byteMeta))

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	return c.JSON(code, result)
}

func ResponseWithPagination(data interface{}, message string, code int, total int64, pagination PaginationRequest, c echo.Context) error {
	success := false
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Request().Method,
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: c.Request().ContentLength,
		Ip:            c.RealIP(),
	}
	byteMeta, _ := json.Marshal(meta)
	LogDefault(string(byteMeta))

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	if !pagination.DisablePagination {
		result.Meta = PaginationResponse{
			Total:   total,
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		}
	}

	return c.JSON(code, result)
}

func ResponseError(err interface{}, c echo.Context) error {

	errObj := getErrStatusCode(err)

	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Request().Method,
		Code:          fmt.Sprintf("%v", errObj),
		Ip:            c.RealIP(),
		ContentLength: c.Request().ContentLength,
	}

	result := BaseWrapperModel{
		Success: false,
		Message: errObj.Message,
		Code:    errObj.Code,
	}

	byteMeta, _ := json.Marshal(meta)

	LogError(string(byteMeta))

	return c.JSON(errObj.ResponseCode, result)
}

func LogDefault(meta string) {
	log.Default().Println("service-info", "Logging service...", "audit-log", meta)
}

func LogError(meta string) {
	log.Default().Println("service-error", "Logging service...", "audit-log", meta)
}

func getErrStatusCode(err interface{}) httperror.CommonErrorData {
	errData := httperror.CommonErrorData{}

	switch obj := err.(type) {
	case httperror.BadRequestData:
		errData.ResponseCode = http.StatusBadRequest
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	case httperror.UnauthorizedData:
		errData.ResponseCode = http.StatusUnauthorized
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	case httperror.ForbiddenErrorData:
		errData.ResponseCode = http.StatusForbidden
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	case httperror.NotFoundData:
		errData.ResponseCode = http.StatusNotFound
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	case httperror.ConflictData:
		errData.ResponseCode = http.StatusConflict
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	case httperror.InternalServerErrorData:
		errData.ResponseCode = http.StatusInternalServerError
		errData.Code = obj.Code()
		errData.Message = obj.Message()
		return errData
	default:
		errData.Code = http.StatusConflict
		return errData
	}
}
