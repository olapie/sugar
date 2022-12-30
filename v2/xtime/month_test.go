package xtime_test

import (
	"testing"

	"code.olapie.com/sugar/xtime"
)

func TestMonth_NumOfWeeks(t *testing.T) {
	tests := []struct {
		Month      *xtime.Month
		NumOfWeeks int
	}{
		{xtime.NewMonth(2020, 6),
			5,
		}, {xtime.NewMonth(2020, 7),
			5,
		}, {xtime.NewMonth(2020, 8),
			6,
		}, {xtime.NewMonth(2020, 9),
			5,
		},
	}

	for _, te := range tests {
		if te.NumOfWeeks != te.Month.NumOfWeeks() {
			t.Fatal(te.NumOfWeeks, te.Month.NumOfWeeks(), te.Month)
		}
	}
}

func TestMonth_GetCalendarDate(t *testing.T) {
	tests := []struct {
		Month *xtime.Month
		Week  int
		Day   int
		Date  *xtime.Date
	}{
		{xtime.NewMonth(2020, 9),
			1,
			3,
			xtime.NewDate(2020, 9, 1),
		}, {xtime.NewMonth(2020, 9),
			4,
			1,
			xtime.NewDate(2020, 9, 20),
		}, {xtime.NewMonth(2020, 9),
			5,
			4,
			xtime.NewDate(2020, 9, 30),
		}, {xtime.NewMonth(2020, 9),
			1,
			1,
			nil,
		}, {xtime.NewMonth(2020, 9),
			5,
			7,
			nil,
		},
	}

	for _, te := range tests {
		date := te.Month.GetCalendarDate(te.Week, te.Day)
		if date == te.Date {
			continue
		}
		if date == nil {
			t.Fatal(te.Month.String(), te.Week, te.Day)
		}
		if te.Date == nil {
			t.FailNow()
		}

		if !date.Equals(te.Date) {
			t.Fatal(date.String(), te.Date.String())
		}
	}
}
