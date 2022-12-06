package timex

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	_ json.Marshaler   = (*Range)(nil)
	_ json.Unmarshaler = (*Range)(nil)
)

type Range struct {
	begin time.Time // inclusive
	end   time.Time // inclusive
}

// NewRange : return a range [begin, end]
func NewRange(begin, end time.Time) *Range {
	if begin.After(end) {
		panic("timex: expect begin <= end")
	}

	begin = begin.Local()
	end = end.Local()

	return &Range{
		begin: begin.Local(),
		end:   end.Local(),
	}
}

func (r *Range) Set(begin, end time.Time) {
	if begin.After(end) {
		panic("timex: expect begin <= end")
	}

	r.begin = begin.Local()
	r.end = end.Local()
}

func (r *Range) SetBegin(t time.Time) {
	if t.After(r.end) {
		panic("timex: expect begin <= end")
	}
	r.begin = t
}

func (r *Range) SetEnd(t time.Time) {
	if t.Before(r.begin) {
		panic("timex: expect end >= begin")
	}
	r.end = t
}

func (r *Range) Begin() time.Time {
	return r.begin
}

func (r *Range) End() time.Time {
	return r.end
}

func (r *Range) BeginUnix() int64 {
	return r.begin.Unix()
}

func (r *Range) EndUnix() int64 {
	return r.end.Unix()
}

func (r *Range) Duration() time.Duration {
	return r.end.Sub(r.begin) + time.Nanosecond
}

func (r *Range) AddDate(years, months, days int) *Range {
	return NewRange(r.begin.AddDate(years, months, days), r.end.AddDate(years, months, days))
}

func (r *Range) Before(ra *Range) bool {
	return r.begin.Before(ra.begin)
}

func (r *Range) After(ra *Range) bool {
	return r.begin.After(ra.begin)
}

func (r *Range) Equals(ra *Range) bool {
	return r.begin.Equal(ra.begin) && r.end.Equal(ra.end)
}

func (r *Range) Contains(ra *Range) bool {
	// r.begin <= ra.begin && r.end >= ra.end
	return !r.begin.After(ra.begin) && !r.end.Before(ra.end)
}

func (r *Range) Intersects(ra *Range) *Range {
	begin, end := r.begin, r.end
	if ra.begin.After(begin) {
		begin = ra.begin
	}
	if ra.end.Before(end) {
		end = ra.end
	}
	if begin.After(end) {
		return nil
	}
	return NewRange(begin, end)
}

func (r *Range) Overlap(ra *Range) bool {
	subBegin := r.begin.Sub(ra.begin)
	switch {
	case subBegin < 0: // r.begin < ra.begin
		return r.end.After(ra.begin) // r.end > ra.begin
	case subBegin > 0: // r.begin > ra.begin
		return !r.end.After(ra.end) // r.end <= ra.end
	default: // r.begin == ra.begin
		return true
	}
}

func (r *Range) ContainsTime(t time.Time) bool {
	return (t.After(r.begin) && t.Before(r.end)) || t == r.begin || t == r.end
}

func (r *Range) ContainsDate(d *Date) bool {
	if r.FirstDay().After(d) {
		return false
	}
	if r.LastDay().Before(d) {
		return false
	}
	return true
}

func (r *Range) IsAllDay() bool {
	return r.InDay() && r.Duration() == time.Hour*24
}

func (r *Range) InDay() bool {
	y1, m1, d1 := r.begin.Date()
	y2, m2, d2 := r.end.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (r *Range) Dates() []*Date {
	begin := DateWithTime(r.begin)
	end := DateWithTime(r.end)
	var l []*Date
	for d := begin; d.Equals(begin) || d.Before(end); d = d.Add(0, 0, 1) {
		l = append(l, d)
	}
	return l
}

func (r *Range) Months() []*Month {
	y, m, _ := r.begin.Date()
	first := NewMonth(y, int(m))
	y, m, _ = r.end.Date()
	last := NewMonth(y, int(m))
	var l []*Month
	for v := first; !v.After(last); v.Add(0, 1) {
		l = append(l, v)
	}
	return l
}

func (r *Range) NumOfMonths() int {
	y, m, _ := r.begin.Date()
	first := NewMonth(y, int(m))
	y, m, _ = r.end.Date()
	last := NewMonth(y, int(m))
	return last.Since(first) + 1
}

func (r *Range) IndexOfMonths(m *Month) int {
	for i, mo := range r.Months() {
		if mo.Equals(m) {
			return i
		}
	}
	return -1
}

func (r *Range) FirstMonth() *Month {
	y, m, _ := r.begin.Date()
	return NewMonth(y, int(m))
}

func (r *Range) LastMonth() *Month {
	y, m, _ := r.end.Date()
	return NewMonth(y, int(m))
}

func (r *Range) FirstDay() *Date {
	return r.BeginT().Date()
}

func (r *Range) LastDay() *Date {
	return r.EndT().Date()
}

func (r *Range) SplitInDay() []*Range {
	dates := r.Dates()
	l := make([]*Range, len(dates))
	beginTime, endTime := GetDayTime(r.begin), GetDayTime(r.end)
	for i, d := range dates {
		begin, end := d.Begin(), d.End()
		if i == 0 {
			begin = begin.Add(beginTime)
		}
		if i == len(dates)-1 {
			if endTime == 0 {
				endTime = Day
			}
			end = d.Begin().Add(endTime)
		}
		l[i] = NewRange(begin, end)
	}
	return l
}

func (r *Range) UnmarshalJSON(b []byte) error {
	var rr struct {
		Begin time.Time `json:"begin"`
		End   time.Time `json:"end"`
	}
	err := json.Unmarshal(b, &rr)
	if err != nil {
		return err
	}
	r.begin = rr.Begin
	r.end = rr.End
	return nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	var rr struct {
		Begin time.Time `json:"begin"`
		End   time.Time `json:"end"`
	}
	rr.Begin = r.begin
	rr.End = r.end
	return json.Marshal(rr)
}

var (
	_ driver.Valuer = (*Range)(nil)
	_ sql.Scanner   = (*Range)(nil)
)

const (
	sqlTimeLayout = "2006-01-02 15:04:05.999999999-07"
	timeLayout    = "2006-01-02 15:04:05-07"
)

func (r *Range) Scan(src interface{}) error {
	var err error
	var s string
	if str, ok := src.(string); ok {
		s = str
	} else if stringer, ok := src.(fmt.Stringer); ok {
		s = stringer.String()
	} else if stringer, ok := src.(fmt.GoStringer); ok {
		s = stringer.GoString()
	} else if b, ok := src.([]byte); ok {
		s = string(b)
	} else {
		return errors.New(fmt.Sprintf("not string: %T", src))
	}

	if s == "" {
		return nil
	}

	s = strings.Replace(s, `"`, "", -1)

	if s[0] != '[' {
		return fmt.Errorf("cannot parse %s", s)
	}

	if c := s[len(s)-1]; c != ']' {
		return fmt.Errorf("cannot parse %s", s)
	}

	s = s[1 : len(s)-1]

	fields := strings.Split(s, ",")
	if len(fields) != 2 {
		return fmt.Errorf("parse composite fields %s", s)
	}
	r.begin, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[0]))
	if err != nil {
		return fmt.Errorf("parse begin %s: %w", fields[0], err)
	}
	r.end, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[1]))
	if err != nil {
		return fmt.Errorf("parse begin %s: %w", fields[1], err)
	}
	if r.begin.After(r.end) {
		return fmt.Errorf("begin %v is after end %v", r.begin, r.end)
	}
	r.begin = r.begin.Local()
	r.end = r.end.Local()
	return nil
}

func (r Range) Value() (driver.Value, error) {
	return fmt.Sprintf("[%s, %s]", r.begin.UTC().Format(sqlTimeLayout), r.end.UTC().Format(sqlTimeLayout)), nil
}

func (r *Range) String() string {
	return fmt.Sprintf("[%s, %s)", r.begin.Format(timeLayout), r.end.Format(timeLayout))
}

func (r *Range) In(loc *time.Location) *Range {
	return &Range{
		begin: r.begin.In(loc),
		end:   r.end.In(loc),
	}
}

func (r *Range) NextRepeat(repeat Repeat) *Range {
	switch repeat {
	case Daily:
		return r.AddDate(0, 0, 1)
	case Weekly:
		return r.AddDate(0, 0, 7)
	case Monthly:
		return r.AddDate(0, 1, 0)
	case Yearly:
		return r.AddDate(1, 0, 0)
	default:
		return nil
	}
}

func (r *Range) PrevRepeat(repeat Repeat) *Range {
	switch repeat {
	case Daily:
		return r.AddDate(0, 0, -1)
	case Weekly:
		return r.AddDate(0, 0, -7)
	case Monthly:
		return r.AddDate(0, -1, 0)
	case Yearly:
		return r.AddDate(-1, 0, 0)
	default:
		return nil
	}
}

func (r *Range) RelativeText() string {
	hans := IsSimplifiedChinese()
	begin, end := r.BeginT(), r.EndT()
	beginText := begin.Date().ShortText()
	if !begin.IsBeginOfDay() {
		beginText += " " + strings.ToLower(begin.TimeText())
	}
	endText := end.Date().ShortText()
	if !end.IsBeginOfDay() {
		endText += " " + strings.ToLower(end.TimeText())
	}
	if r.InDay() {
		switch {
		case r.IsAllDay():
			if hans {
				return beginText + "全天"
			}
			return beginText + " all day"
		case begin.Equals(end):
			return beginText
		case begin.IsBeginOfDay():
			if hans {
				return endText + " 结束"
			}
			return endText + " ends"
		case end.IsEndOfDay():
			if hans {
				return beginText + " 开始"
			}
			return beginText + " begins"
		default:
			return beginText + " - " + strings.ToLower(end.TimeText())
		}
	}
	if begin.Year() == end.Year() || end.Year()-1 == time.Now().Year() {
		endText = strings.Replace(endText, fmt.Sprintf("%d-", begin.Year()), "", 1)
	}
	return beginText + " - " + endText
}
