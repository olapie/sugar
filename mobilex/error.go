package mobilex

import (
	"code.olapie.com/sugar/errorx"
	"fmt"
	"reflect"
)

type Error errorx.Error

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) String() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ToError(err error) *Error {
	if err == nil {
		return nil
	}

	if v := reflect.ValueOf(err); !v.IsValid() || v.IsZero() {
		return nil
	}

	cause := errorx.Cause(err)
	if e, ok := cause.(*Error); ok && e != nil {
		return NewError((*errorx.Error)(e).Code, err.Error())
	}

	if e, ok := cause.(*errorx.Error); ok && e != nil {
		return NewError(e.Code, e.Message)
	}

	return NewError(0, err.Error())
}
