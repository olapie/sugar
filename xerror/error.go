package xerror

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var errorRegexp1 = regexp.MustCompile(`^code:(\d+)$`)
var errorRegexp2 = regexp.MustCompile(`^code:(\d+), message:(.*)$`)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func FromString(s string) *Error {
	s = strings.TrimSpace(s)
	texts := errorRegexp1.FindStringSubmatch(s)
	if len(texts) == 2 {
		code, err := strconv.ParseInt(texts[1], 0, 64)
		if err == nil {
			return &Error{
				Code: int(code),
			}
		}
	}

	texts = errorRegexp2.FindStringSubmatch(s)
	if len(texts) != 3 {
		return nil
	}

	code, err := strconv.ParseInt(texts[1], 0, 64)
	if err != nil {
		return nil
	}

	return &Error{
		Code:    int(code),
		Message: texts[2],
	}
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Error() string {
	if e.Message == "" {
		e.Message = http.StatusText(e.Code)
		if e.Message == "" {
			e.Message = fmt.Sprint(e.Code)
		}
	}
	return e.Message
}

func (e *Error) Is(target error) bool {
	if e == target {
		return true
	}

	if t, ok := target.(*Error); ok {
		return t.Code == e.Code && t.Message == e.Message
	}
	return false
}

func (e *Error) Respond(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	body, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal json: %v", err)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Printf("write body: %v", err)
	}
}

func Format(code int, format string, a ...any) *Error {
	if len(a) == 0 {
		err := FromString(format)
		if err != nil {
			err.Code = code
			return err
		}
	}

	msg := fmt.Sprintf(format, a...)
	if msg == "" {
		msg = http.StatusText(code)
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}

var rawErrType = reflect.TypeOf(errors.New(""))

func GetCode(err error) int {
	err = Cause(err)
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

	keys := []string{"status", "Status", "status_code", "StatusCode", "statusCode", "code", "Code"}
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
// indirect returns the value, after dereferencing as many times
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

func IsNotExist(err error) bool {
	return errors.Is(err, NotExist) ||
		errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, os.ErrNotExist) ||
		GetCode(err) == http.StatusNotFound
}

func IsNilOrNotExist(err error) bool {
	return err == nil || err.Error() == "nil" || IsNotExist(err)
}

func Or(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func And(errs ...error) error {
	for _, err := range errs {
		if err == nil {
			return nil
		}
	}
	return errorSlice(errs)
}

func Combine(errs ...error) error {
	for i := len(errs) - 1; i >= 0; i-- {
		if errs[i] == nil {
			errs = append(errs[:i], errs[i+1:]...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errorSlice(errs)
}

func OrFn(err error, errFns ...func() error) error {
	if err != nil {
		return err
	}

	for _, fn := range errFns {
		if e := fn(); e != nil {
			return e
		}
	}

	return nil
}
