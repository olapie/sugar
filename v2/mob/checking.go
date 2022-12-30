package mob

import (
	"strings"
	"unicode"

	"code.olapie.com/sugar/xcheck"
	"code.olapie.com/sugar/xtime"
)

func IsEmailAddress(s string) bool {
	return xcheck.IsEmailAddress(s)
}

func IsURL(s string) bool {
	return xcheck.IsURL(s)
}

func IsDate(s string) bool {
	return xtime.IsDate(s)
}

var (
	MinPasswordLen int = 6
	MinUsernameLen int = 4
	MaxUsernameLen int = 20
)

func IsValidPassword(password string) bool {
	if len(password) < MinPasswordLen {
		return false
	}

	hasDigit := false
	hasAlpha := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasAlpha = true
		}
	}
	return hasDigit && hasAlpha
}

func IsValidUsername(username string) bool {
	username = strings.ToLower(username)
	s := []rune(username)
	if len(s) > MaxUsernameLen {
		return false
	}

	if len(s) < MinUsernameLen {
		return false
	}

	for _, c := range s {
		if unicode.IsDigit(c) || c == '_' || (c >= 'a' && c <= 'z') {
			continue
		}
		return false
	}
	return true
}
