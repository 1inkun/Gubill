package models

import (
	"fmt"
)

type BusinessError struct {
	Code int
	Msg  string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Msg)
}

func NewBusinessError(code int, msg string) *BusinessError {
	return &BusinessError{
		Code: code,
		Msg:  msg,
	}
}

type InternalError struct {
	Code     int
	Msg      string
	InterErr error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Msg)
}

func NewInternalError(code int, msg string, err error) *InternalError {
	return &InternalError{
		Code:     code,
		Msg:      msg,
		InterErr: err,
	}
}

func NewDatabaseErr(err error) *InternalError {
	return &InternalError{
		Code:     500,
		Msg:      "数据库操作出错",
		InterErr: err,
	}
}
