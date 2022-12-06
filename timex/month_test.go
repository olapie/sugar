package timex_test

import (
	"testing"

	"code.olapie.com/sugar/timex"
)

func TestMonth_NumOfWeeks(t *testing.T) {
	tests := []struct {
		Month      *timex.Month
		NumOfWeeks int
	}{
		{timex.NewMonth(2020, 6),
			5,
		}, {timex.NewMonth(2020, 7),
			5,
		}, {timex.NewMonth(2020, 8),
			6,
		}, {timex.NewMonth(2020, 9),
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
		Month *timex.Month
		Week  int
		Day   int
		Date  *timex.Date
	}{
		{timex.NewMonth(2020, 9),
			1,
			3,
			timex.NewDate(2020, 9, 1),
		}, {timex.NewMonth(2020, 9),
			4,
			1,
			timex.NewDate(2020, 9, 20),
		}, {timex.NewMonth(2020, 9),
			5,
			4,
			timex.NewDate(2020, 9, 30),
		}, {timex.NewMonth(2020, 9),
			1,
			1,
			nil,
		}, {timex.NewMonth(2020, 9),
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
