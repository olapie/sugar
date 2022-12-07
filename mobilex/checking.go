package mobilex

import (
	"code.olapie.com/sugar/checking"
	"code.olapie.com/sugar/timing"
)

func IsEmailAddress(s string) bool {
	return checking.IsEmailAddress(s)
}

func IsURL(s string) bool {
	return checking.IsURL(s)
}

func IsDate(s string) bool {
	return timing.IsDate(s)
}
