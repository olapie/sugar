package xerror

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

var rawErrType = reflect.TypeOf(errors.New(""))

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
	if message == "" {
		message = http.StatusText(code)
	}
	return &Error{
		code:    code,
		subCode: subCode,
		message: message,
	}
}

func GetCode(err error) int {
	if code := getCode(err); code != 0 {
		return code
	}

	if u, ok := err.(unwrapError); ok {
		return getCode(u.Unwrap())
	}

	if u, ok := err.(unwrapErrors); ok {
		for _, v := range u.Unwrap() {
			if code := GetCode(v); code != 0 {
				return code
			}
		}
	}
	return 0
}

func getCode(err error) int {
	if err == nil {
		return 0
	}

	if reflect.TypeOf(err) == rawErrType {
		return 0
	}

	if err == NotExist || err == sql.ErrNoRows {
		return http.StatusNotFound
	}

	if coder, ok := err.(interface{ Status() int }); ok {
		return coder.Status()
	}

	if coder, ok := err.(interface{ StatusCode() int }); ok {
		return coder.StatusCode()
	}

	if coder, ok := err.(interface{ Code() int }); ok {
		return coder.Code()
	}

	if v := reflect.ValueOf(err); v.Kind() == reflect.Int {
		n := int(v.Int())
		if n > 0 {
			return n
		}
		return 0
	}

	keys := []string{"status", "Status", "status_code", "StatusCode", "statusCode", "code", "code"}
	i := indirect(err)
	k := reflect.ValueOf(i).Kind()
	if k != reflect.Struct && k != reflect.Map {
		return 0
	}

	b, jErr := json.Marshal(i)
	if jErr != nil {
		log.Printf("Cannot marshal json: %v", err)
		return 0
	}
	var m map[string]any
	jErr = json.Unmarshal(b, &m)
	if jErr != nil {
		log.Printf("Cannot unmarshal json: %v", err)
		return 0
	}

	for _, k := range keys {
		if v, ok := m[k]; ok {
			if i, err := strconv.ParseInt(fmt.Sprint(v), 0, 0); err == nil && i > 0 {
				return int(i)
			}
		}
	}
	return 0
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after resolving as many times
// as necessary to reach the base type (or nil).
func indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
