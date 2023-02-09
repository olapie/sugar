package timing_test

import (
	"testing"

	"code.olapie.com/sugar/v2/timing"
)

func TestLenOfMonth(t *testing.T) {
	tests := []struct {
		Year  int
		Month int
		Len   int
	}{
		{
			Year:  2020,
			Month: 1,
			Len:   31,
		},
		{
			Year:  2020,
			Month: 2,
			Len:   29,
		},
		{
			Year:  2021,
			Month: 2,
			Len:   28,
		},
		{
			Year:  2021,
			Month: 3,
			Len:   31,
		},
		{
			Year:  2021,
			Month: 4,
			Len:   30,
		},
		{
			Year:  2021,
			Month: 5,
			Len:   31,
		},
		{
			Year:  2021,
			Month: 6,
			Len:   30,
		},
		{
			Year:  2021,
			Month: 7,
			Len:   31,
		},
		{
			Year:  2021,
			Month: 8,
			Len:   31,
		},
		{
			Year:  2021,
			Month: 9,
			Len:   30,
		},
		{
			Year:  2021,
			Month: 10,
			Len:   31,
		},
		{
			Year:  2021,
			Month: 11,
			Len:   30,
		},
		{
			Year:  2021,
			Month: 12,
			Len:   31,
		},
	}

	for _, test := range tests {
		got := timing.NewMonth(test.Year, test.Month).NumOfDays()
		if test.Len != got {
			t.Errorf("%v expect: %v, got %v", test, test.Len, got)
		}
	}
}
