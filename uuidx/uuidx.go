package uuidx

import (
	"encoding/binary"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func NewShortString() string {
	return ShortString(uuid.New())
}

func ShortString(id uuid.UUID) string {
	_ = id[15]
	i1 := binary.BigEndian.Uint64(id[0:8])
	i2 := binary.BigEndian.Uint64(id[8:])
	return shortStringOfUint64(i1) + "-" + shortStringOfUint64(i2)
}

func FromShortString(s string) uuid.NullUUID {
	strs := strings.Split(s, "-")
	i1, err := uint64FromShortString(strs[0])
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}
	i2, err := uint64FromShortString(strs[1])
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}
	var id uuid.NullUUID
	binary.BigEndian.PutUint64(id.UUID[0:8], i1)
	binary.BigEndian.PutUint64(id.UUID[8:], i2)
	id.Valid = true
	return id
}

func shortStringOfUint64(i uint64) string {
	var bytes [16]byte
	n := 15
	for {
		j := i % 62
		switch {
		case j <= 9:
			bytes[n] = byte('0' + j)
		case j <= 35:
			bytes[n] = byte('A' + j - 10)
		default:
			bytes[n] = byte('a' + j - 36)
		}
		i /= 62
		if i == 0 {
			return string(bytes[n:])
		}
		n--
	}
}

func uint64FromShortString(s string) (uint64, error) {
	if len(s) == 0 {
		return 0, errors.New("parse error")
	}

	var bytes = []byte(s)
	var k uint64
	var v uint64
	for _, b := range bytes {
		switch {
		case b >= '0' && b <= '9':
			v = uint64(b - '0')
		case b >= 'A' && b <= 'Z':
			v = uint64(10 + b - 'A')
		case b >= 'a' && b <= 'z':
			v = uint64(36 + b - 'a')
		default:
			return 0, errors.New("parse error")
		}
		k = k*62 + v
	}
	return k, nil
}
