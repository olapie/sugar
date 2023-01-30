package mob

import (
	"net/http"
	"time"

	"code.olapie.com/sugar/v2/conv"
	"code.olapie.com/sugar/v2/xstring"
	"code.olapie.com/sugar/v2/xtime"
	"github.com/google/uuid"
)

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
	time.Local = xtime.ToLocation(name, offset)
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

func SmartLen(s string) int {
	n := 0
	for _, c := range s {
		if c <= 255 {
			n++
		} else {
			n += 2
		}
	}

	return n
}

func SquishString(s string) string {
	return xstring.Squish(s)
}

type Handler interface {
	SaveSecret(name, value string) bool
	DeleteSecret(name string) bool
	GetSecret(name string) string
	NeedSignIn()
}

type AuthErrorChecker struct {
	h Handler
}

func NewAuthErrorChecker(h Handler) *AuthErrorChecker {
	return &AuthErrorChecker{
		h: h,
	}
}

func (c *AuthErrorChecker) Check(err error) {
	if err == nil {
		return
	}

	code := ToError(err).Code()
	if code == http.StatusUnauthorized {
		go c.h.NeedSignIn()
	}
}

func GetSizeString(n int64) string {
	return conv.SizeToHumanReadable(n)
}
