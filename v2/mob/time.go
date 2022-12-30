package mob

import "code.olapie.com/sugar/v2/xtime"

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
