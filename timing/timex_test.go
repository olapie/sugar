package timing_test

import (
	"testing"
	"time"

	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/timing"
)

func TestToDuration(t *testing.T) {
	type Test struct {
		Value  string
		Result time.Duration
	}
	tests := []Test{
		{
			"1s",
			time.Second,
		},
		{
			"5m",
			5 * time.Minute,
		},
		{
			"3h",
			3 * time.Hour,
		},
		{
			"3h24m",
			3*time.Hour + 24*time.Minute,
		},
		{
			"1us",
			time.Microsecond,
		},
		{
			"3.5m",
			3*time.Minute + 30*time.Second,
		},
	}

	for _, te := range tests {
		d, err := timing.ToDuration(te.Value)
		if err != nil {
			t.Error(err)
		}
		testx.Equal(t, te.Result, d)
	}
}

func TestToTime(t *testing.T) {
	type Test struct {
		Value  string
		Result time.Time
	}
	date20180102 := time.Date(2018, 1, 2, 0, 0, 0, 0, time.Local)
	tests := []Test{
		{
			"2018-1-2",
			date20180102,
		},
		{
			"2018-01-02",
			date20180102,
		},
		{
			"2018/1/2",
			date20180102,
		},
		{
			"2018/01/02",
			date20180102,
		},
		{
			"20180102",
			date20180102,
		},
		{
			"02/01/2018",
			date20180102,
		},
		{
			"2/1/2018",
			date20180102,
		},
	}

	for _, te := range tests {
		d, err := timing.ToTime(te.Value)
		if err != nil {
			t.Error(err)
		}
		testx.Equal(t, te.Result, d)
	}
}
