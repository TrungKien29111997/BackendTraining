package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrBadRequest          ErrorCode = "bad_request"
	ErrNotFound            ErrorCode = "not_found"
	ErrConflict            ErrorCode = "conflict"
	ErrInternalServerError ErrorCode = "internal_server_error"
)

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (err *AppError) Error() string {
	return ""
}

func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func WrapError(err error, message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err}
}

func ResponeError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"message": appErr.Message,
			"code":    appErr.Code,
		}
		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrInternalServerError,
	})
}

func ReponseSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ReponseStatusCode(c *gin.Context, status int) {
	c.Status(status)
}
