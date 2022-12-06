package timex

import (
	"fmt"
	"log"
	"strings"
	"time"

	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/rtx"
)

var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-1-2",
	"20060102",
	"2006/1/2",
	"2/1/2006",
}

// ToDuration converts an interface to a time.Duration type.
func ToDuration(i any) (time.Duration, error) {
	i = rtx.Indirect(i)

	if s, err := conv.ToString(i); err == nil {
		s = strings.ToLower(s)
		if strings.ContainsAny(s, "nsuµmh") {
			return time.ParseDuration(s)
		} else {
			return time.ParseDuration(s + "ns")
		}
	}

	if n, err := conv.ToInt64(i); err == nil {
		return time.Duration(n), nil
	}

	return 0, fmt.Errorf("cannot convert %#v of type %T to duration", i, i)
}

func ToTime(i any) (time.Time, error) {
	return ToTimeInLocation(i, time.Local)
}

func ToTimeInLocation(i any, loc *time.Location) (time.Time, error) {
	i = rtx.Indirect(i)
	s, err := conv.ToString(i)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot convert %#v of type %T to date: %w", i, i, err)
	}
	for _, df := range timeFormats {
		d, err := time.ParseInLocation(df, s, loc)
		if err == nil {
			return d, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot convert %#v of type %T to date", i, i)
}

func ToLocation(name string, offset int) *time.Location {
	// LoadLocation get failed on iOS
	loc, err := time.LoadLocation(name)
	if err == nil {
		return loc
	}
	log.Printf("Cannot load location %s: %v. Converted to a fixed zone", name, err)
	loc = time.FixedZone(name, offset)
	return loc
}

func IsDate(s string) bool {
	if len(s) > 10 {
		return false
	}
	_, err := ToTime(s)
	return err == nil
}

func ToDateString(s string) string {
	s = strings.Replace(s, "年", "-", 1)
	s = strings.Replace(s, "月", "-", 1)
	s = strings.Replace(s, "日", "", 1)
	s = strings.Replace(s, "o", "0", -1)
	s = strings.Replace(s, "O", "0", -1)
	s = strings.Replace(s, "l", "1", -1)
	s = strings.Replace(s, "I", "1", -1)
	if strings.Contains(s, "-") && len(s) >= 10 {
		s = s[:10]
	} else if len(s) >= 8 {
		s = s[:8]
	} else {
		return ""
	}

	t, err := ToTime(s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}
