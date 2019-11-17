package errutil

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type Error interface {
	Error() string
	GetCode() int
}

type GeneralError struct {
	Err    error
	Status int
	Msg    string
	Field  string
}

func NewFromCode(s int, msg string) Error {
	return &GeneralError{Msg: msg, Status: s}
}

func (e *GeneralError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Msg
}

func (e *GeneralError) GetCode() int {
	return e.Status
}

type NotFoundError struct {
	GeneralError
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{GeneralError{Err: err, Status: http.StatusNotFound}}
}

func (e *NotFoundError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "not found"
}

func (e *NotFoundError) GetCode() int {
	return e.Status
}

type InternalServerError struct {
	GeneralError
}

func (e *InternalServerError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "internal"
}

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{GeneralError{Err: err, Status: http.StatusInternalServerError}}
}

type BadRequestError struct {
	GeneralError
}

func (e *BadRequestError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "bad request"
}

func NewBadRequestError(err error) *BadRequestError {
	return &BadRequestError{GeneralError{Err: err, Status: http.StatusBadRequest}}
}

func NewBadRequestFieldError(field, msg string) *BadRequestError {
	return &BadRequestError{GeneralError{Field: field, Msg: msg, Status: http.StatusBadRequest}}
}
