package response

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("not found")
	ErrExists       = errors.New("already exists")
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type CustomError struct {
	Err    error
	Msg    string
	Status int
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func NewCustomError(err error, msg string, status int) *CustomError {
	return &CustomError{
		Err:    err,
		Msg:    msg,
		Status: status,
	}
}

func AsCustomErr(err error) (*CustomError, bool) {
	var customErr *CustomError
	ok := errors.As(err, &customErr)
	return customErr, ok
}
