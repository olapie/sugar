package timing_test

import (
	"testing"

	"code.olapie.com/sugar/v2/timing"
)

func TestMonth_NumOfWeeks(t *testing.T) {
	tests := []struct {
		Month      *timing.Month
		NumOfWeeks int
	}{
		{timing.NewMonth(2020, 6),
			5,
		}, {timing.NewMonth(2020, 7),
			5,
		}, {timing.NewMonth(2020, 8),
			6,
		}, {timing.NewMonth(2020, 9),
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
		Month *timing.Month
		Week  int
		Day   int
		Date  *timing.Date
	}{
		{timing.NewMonth(2020, 9),
			1,
			3,
			timing.NewDate(2020, 9, 1),
		}, {timing.NewMonth(2020, 9),
			4,
			1,
			timing.NewDate(2020, 9, 20),
		}, {timing.NewMonth(2020, 9),
			5,
			4,
			timing.NewDate(2020, 9, 30),
		}, {timing.NewMonth(2020, 9),
			1,
			1,
			nil,
		}, {timing.NewMonth(2020, 9),
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
