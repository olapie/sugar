package mobilex

import (
	"code.olapie.com/sugar/errorx"
	"code.olapie.com/sugar/timing"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type IntE struct {
	Val int
	Err *Error
}

type BoolE struct {
	Val bool
	Err *Error
}

type StringE struct {
	Val string
	Err *Error
}

type Int64E struct {
	Val int64
	Err *Error
}

type Float64E struct {
	Val float64
	Err *Error
}

type ByteArrayE struct {
	Val []byte
	Err *Error
}

type SecretManager interface {
	Get(key string) string
	Set(key, data string) bool
	Del(key string) bool
}

type Uptimer interface {
	Uptime() int64
}

const (
	KeyDeviceID = "mobile_device_id"
)

func GetDeviceID(m SecretManager) string {
	id := m.Get(KeyDeviceID)
	if id == "" {
		id = uuid.NewString()
		m.Set(KeyDeviceID, id)
	}
	return id
}

func SetTimeZone(name string, offset int) {
	time.Local = timing.ToLocation(name, offset)
}

func GetTimeZoneOffset() int {
	_, o := time.Now().Zone()
	return o
}

func GetTimeZoneName() string {
	n, _ := time.Now().Zone()
	return n
}

func NewUUID() string {
	return uuid.New().String()
}

type Now interface {
	Now() int64
}

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
