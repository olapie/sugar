package xtime_test

import (
	"testing"
	"time"

	"code.olapie.com/sugar/xtime"
)

func TestRange_SplitInDay(t *testing.T) {
	tz := time.FixedZone("PST", -7*3600)
	begin := time.Date(2002, 5, 3, 17, 0, 0, 0, tz)
	end := time.Date(2002, 5, 3, 18, 0, 0, 0, tz)
	r := xtime.NewRange(begin, end)
	for _, dr := range r.SplitInDay() {
		t.Log(dr.Begin(), dr.End())
	}
}
