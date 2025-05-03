package server

import "fmt"

type Error struct {
	Err     error
	Method  string
	Problem string
	Status  uint
}

func NewError(Err error,
	Method string,
	Problem string,
	Type uint) *Error {
	return &Error{
		Err:     Err,
		Method:  Method,
		Problem: Problem,
		Status:  Type,
	}
}

const (
	ErrUnknown = iota
)

func (e *Error) Error() string {
	return fmt.Sprintf("Method %s encountered error: %s %s", e.Method, e.Problem, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
