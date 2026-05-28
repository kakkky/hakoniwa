package domain

import (
	"fmt"
)

var (
	NotfoundErr = newErr("NOT FOUND ERROR")
)

type err struct {
	msg     string
	wrapped error
}

func newErr(msg string) *err {
	return &err{
		msg: msg,
	}
}

func (e *err) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("%s : %s", e.msg, e.wrapped.Error())
	}
	return e.msg
}

func (e *err) With(msg string) *err {
	e.msg = msg
	return e
}

func (e *err) Wrap(err error) error {
	e.wrapped = err
	return e
}

func (e *err) Unwrap() error {
	return e.wrapped
}
