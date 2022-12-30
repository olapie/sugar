package mob

import (
	"fmt"
	"reflect"

	"code.olapie.com/sugar/v2/xerror"
)

type Error xerror.Error

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

	cause := xerror.Cause(err)
	if e, ok := cause.(*Error); ok && e != nil {
		return NewError((*xerror.Error)(e).Code, err.Error())
	}

	if e, ok := cause.(*xerror.Error); ok && e != nil {
		return NewError(e.Code, e.Message)
	}

	return NewError(0, err.Error())
}
