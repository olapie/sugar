package mob

import (
	"fmt"

	"code.olapie.com/sugar/v2/xtime"
)

func GetDateTimeString(t int64) string {
	tm := xtime.TimeWithUnix(t)
	return tm.Date().PrettyText() + " " + tm.TimeTextWithZero()
}

func GetRelativeDateTimeString(t int64) string {
	tm := xtime.TimeWithUnix(t)
	return tm.RelativeDateTimeText()
}

type Time = xtime.Time

func NowTime() *Time {
	return (*Time)(xtime.NewTime())
}

func TimeWithUnix(seconds int64) *Time {
	return xtime.TimeWithUnix(seconds)
}

func TimerText(elapse int64) string {
	h := elapse / 3600
	elapse %= 3600
	m := elapse / 60
	s := elapse % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
