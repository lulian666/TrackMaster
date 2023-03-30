package pkg

import (
	"fmt"
	"net/http"
)

const (
	Success = iota
	ServerError
	InvalidParams
	NotFound
	TooManyRequests
)

type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

func NewError(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code: %d, error message:: %s", e.Code, e.Msg)
}

func (e *Error) WithDetails(details ...string) *Error {
	e.Details = []string{}
	for _, d := range details {
		e.Details = append(e.Details, d)
	}

	return e
}

func (e *Error) StatusCode() int {
	switch e.Code {
	case Success:
		return http.StatusOK
	case ServerError:
		return http.StatusInternalServerError
	case InvalidParams:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case TooManyRequests:
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
