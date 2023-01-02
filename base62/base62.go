package base62

import (
	"errors"
	"math/big"

	"github.com/google/uuid"
)

func EncodeToString(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	var i big.Int
	i.SetBytes(src)
	return i.Text(62)
}

func DecodeString(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	var i big.Int
	_, ok := i.SetString(s, 62)
	if !ok {
		return nil, errors.New("illegal base62 data")
	}
	return i.Bytes(), nil
}

func NewUUIDString() string {
	id := uuid.New()
	return EncodeToString(id[:])
}

func UUIDFromString(s string) (id uuid.UUID, err error) {
	b, err := DecodeString(s)
	if err != nil {
		return id, err
	}
	return uuid.FromBytes(b)
}
