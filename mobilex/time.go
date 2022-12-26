package mobilex

import "code.olapie.com/sugar/timing"

func GetDateTimeString(t int64) string {
	tm := timing.TimeWithUnix(t)
	return tm.Date().PrettyText() + " " + tm.TimeTextWithZero()
}

func GetRelativeDateTimeString(t int64) string {
	tm := timing.TimeWithUnix(t)
	return tm.RelativeDateTimeText()
}

type Time = timing.Time

func NowTime() *Time {
	return (*Time)(timing.NewTime())
}

func TimeWithUnix(seconds int64) *Time {
	return timing.TimeWithUnix(seconds)
}
