package xcontact

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func (x *PhoneNumber) ToString() string {
	return fmt.Sprintf("+%d%d", x.Code, x.Number)
}

func (x *PhoneNumber) InternationalFormat() string {
	pn, err := phonenumbers.Parse(x.ToString(), "")
	if err != nil {
		return ""
	}
	return phonenumbers.Format(pn, phonenumbers.INTERNATIONAL)
}

func (x *PhoneNumber) MaskString() string {
	nnBytes := []byte(fmt.Sprint(x.Number))
	maskLen := (len(nnBytes) + 2) / 3
	start := len(nnBytes) - 2*maskLen
	for i := 0; i < maskLen; i++ {
		nnBytes[start+i] = '*'
	}

	nn := string(nnBytes)

	if len(x.Extension) == 0 {
		return fmt.Sprintf("+%d%s", x.Code, nn)
	}

	return fmt.Sprintf("+%d%s-%s", x.Code, nn, x.Extension)
}

func (x *PhoneNumber) AccountType() string {
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
		Code:      int32(pn.GetCountryCode()),
		Number:    int64(pn.GetNationalNumber()),
		Extension: pn.GetExtension(),
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
