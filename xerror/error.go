package xerror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type String string

func (s String) Error() string {
	return string(s)
}

const (
	NotExist String = "not exist"
)

type Error struct {
	code    int
	subCode int
	message string
}

type errorJSONObject struct {
	Code    int    `json:"code,omitempty"`
	SubCode int    `json:"sub_code,omitempty"`
	Message string `json:"message,omitempty"`
}

var _ json.Marshaler = (*Error)(nil)
var _ json.Unmarshaler = (*Error)(nil)

func (e *Error) Code() int {
	return e.code
}

func (e *Error) SubCode() int {
	return e.subCode
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Error() string {
	if e.message == "" {
		e.message = http.StatusText(e.code)
		if e.message == "" {
			e.message = fmt.Sprint(e.code)
		} else if e.subCode > 0 {
			e.message = fmt.Sprintf("%s (%d)", e.message, e.subCode)
		}
	}
	return e.message
}

func (e *Error) Is(target error) bool {
	if e == target {
		return true
	}

	if t, ok := target.(*Error); ok {
		return t.code == e.code && t.subCode == e.subCode && t.message == e.message
	}
	return false
}

func (e *Error) MarshalJSON() (text []byte, err error) {
	obj := &errorJSONObject{
		Code:    e.code,
		SubCode: e.subCode,
		Message: e.message,
	}
	return json.Marshal(obj)
}

func (e *Error) UnmarshalJSON(text []byte) error {
	var obj errorJSONObject
	err := json.Unmarshal(text, &obj)
	if err != nil {
		return err
	}
	e.code = obj.Code
	e.subCode = obj.SubCode
	e.message = obj.Message
	return nil
}

func New(code int, format string, a ...any) *Error {
	if code <= 0 {
		panic("invalid code")
	}
	msg := fmt.Sprintf(format, a...)
	if msg == "" {
		msg = http.StatusText(code)
	}
	return &Error{
		code:    code,
		message: msg,
	}
}

func NewSub(code, subCode int, message string) *Error {
	if code <= 0 {
		panic("invalid code")
	}

	if subCode <= 0 {
		panic("invalid subCode")
	}

	if message == "" {
		message = http.StatusText(code)
	}
	return &Error{
		code:    code,
		subCode: subCode,
		message: message,
	}
}
