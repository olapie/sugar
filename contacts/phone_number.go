package contacts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	Code   int32 `json:"code,omitempty"`
	Number int64 `json:"number,omitempty"`
}

func (pn *PhoneNumber) String() string {
	return fmt.Sprintf("+%d%d", pn.Code, pn.Number)
}

func (pn *PhoneNumber) InternationalFormat() string {
	n, err := phonenumbers.Parse(pn.String(), "")
	if err != nil {
		return ""
	}
	return phonenumbers.Format(n, phonenumbers.INTERNATIONAL)
}

func (pn *PhoneNumber) MaskString() string {
	nnBytes := []byte(fmt.Sprint(pn.Number))
	maskLen := (len(nnBytes) + 2) / 3
	start := len(nnBytes) - 2*maskLen
	for i := 0; i < maskLen; i++ {
		nnBytes[start+i] = '*'
	}

	nn := string(nnBytes)
	return fmt.Sprintf("+%d%s", pn.Code, nn)
}

func (pn *PhoneNumber) AccountType() string {
	return "phone_number"
}

func NewPhoneNumber(s string) (*PhoneNumber, error) {
	s = trimPhoneNumberString(s)
	pn, err := phonenumbers.Parse(s, "")
	if err != nil {
		return nil, err
	}

	if !phonenumbers.IsValidNumber(pn) {
		return nil, errors.New("invalid phone number")
	}

	return &PhoneNumber{
		Code:   pn.GetCountryCode(),
		Number: int64(pn.GetNationalNumber()),
	}, nil
}

func NewPhoneNumberV2(s string, code int) (*PhoneNumber, error) {
	s = trimPhoneNumberString(s)
	pn, err := NewPhoneNumber(s)
	if err == nil {
		return pn, nil
	}

	if s == "" {
		return nil, errors.New("invalid phone number")
	}

	if s[0] != '+' && code != 0 {
		s = fmt.Sprintf("+%d%s", code, s)
		return NewPhoneNumber(s)

	}
	return nil, err
}

func MustPhoneNumber(s string) *PhoneNumber {
	pn, err := NewPhoneNumber(s)
	if err != nil {
		panic(err)
	}
	return pn
}

func trimPhoneNumberString(s string) string {
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return s
}

func IsPhoneNumber(s string) bool {
	if s == "" {
		return false
	}
	n, err := phonenumbers.Parse(s, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(n)
}
