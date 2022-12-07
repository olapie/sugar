package mobilex

import (
	"time"

	"code.olapie.com/sugar/stringx"
	"code.olapie.com/sugar/timing"
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
	return stringx.Squish(s)
}
