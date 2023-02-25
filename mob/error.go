package mob

import (
	"errors"
	"reflect"

	"code.olapie.com/sugar/v2/xerror"
)

type Error xerror.Error

func (e *Error) Code() int {
	return (*xerror.Error)(e).Code()
}

func (e *Error) Message() string {
	return (*xerror.Error)(e).Message()
}

func (e *Error) Error() string {
	return (*xerror.Error)(e).Error()
}

func (e *Error) String() string {
	return (*xerror.Error)(e).String()
}

func NewError(code int, message string) *Error {
	return (*Error)(xerror.New(code, message))
}

func ToError(err error) *Error {
	if err == nil {
		return nil
	}

	if v := reflect.ValueOf(err); !v.IsValid() || v.IsZero() {
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		return e
	}

	var xe *xerror.Error
	if errors.As(err, &xe) {
		return NewError(xe.Code(), xe.Message())
	}

	return NewError(0, err.Error())
}
