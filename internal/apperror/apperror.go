// Package apperror lets services describe a failure in HTTP terms, the way a
// NestJS service throws an HttpException. Handlers just return the error and
// the exception filter in utils.HandleError turns it into a response.
package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func Wrap(status int, message string, err error) *AppError {
	return &AppError{Status: status, Message: message, Err: err}
}

func BadRequest(message string) *AppError   { return New(http.StatusBadRequest, message) }
func Unauthorized(message string) *AppError { return New(http.StatusUnauthorized, message) }
func Forbidden(message string) *AppError    { return New(http.StatusForbidden, message) }
func NotFound(message string) *AppError     { return New(http.StatusNotFound, message) }
func Conflict(message string) *AppError     { return New(http.StatusConflict, message) }

func Internal(message string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, message, err)
}
