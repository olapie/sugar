package conv

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/testutil"
)

func TestToBool(t *testing.T) {
	type Bool bool
	type Int int
	goodCases := []struct {
		Value  any
		Result bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
		{
			"TRUE",
			true,
		},
		{
			"FALSE",
			false,
		},
		{
			true,
			true,
		},
		{
			false,
			false,
		},
		{
			false,
			false,
		},
		{
			1,
			true,
		},
		{
			0,
			false,
		},
		{
			-1,
			true,
		},
		{
			10.01,
			true,
		},
		{
			-10.90,
			true,
		},
		{
			0.0,
			false,
		},
		{
			Bool(true),
			true,
		},
		{
			Bool(false),
			false,
		},
		{
			Int(100),
			true,
		},
		{
			Int(0),
			false,
		},
		{
			"t",
			true,
		},
		{
			"f",
			false,
		},
		{
			"T",
			true,
		},
		{
			"F",
			false,
		},
		{
			[]byte("true"),
			true,
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToBool(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		"hello", time.Now(), Foo{},
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToBool(c)
			if err == nil {
				t.Error("should fail", c)
			}
			testutil.Equal(t, false, res)
		}
	})
}

func TestToFloat64(t *testing.T) {
	type Bool bool
	type Int int
	type Float64 float64
	goodCases := []struct {
		Value  any
		Result float64
	}{
		{
			true,
			1,
		}, {
			false,
			0,
		}, {
			123456,
			123456,
		}, {
			Int(123456),
			123456,
		}, {
			int8(8),
			8,
		}, {
			int16(16),
			16,
		}, {
			int32(32),
			32,
		}, {
			int64(64),
			64,
		}, {
			uint8(80),
			80,
		}, {
			uint16(160),
			160,
		}, {
			uint32(320),
			320,
		}, {
			uint64(123),
			123,
		}, {
			math.MaxFloat64,
			math.MaxFloat64,
		}, {
			-math.MaxFloat64,
			-math.MaxFloat64,
		}, {
			0.0,
			0,
		},
		{
			-123.1,
			-123.1,
		},
		{
			Bool(false),
			0,
		},
		{
			[]byte("123.309"),
			123.309,
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToFloat64(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		"hello", time.Now(), Foo{},
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToFloat64(c)
			if err == nil {
				t.Error("should fail", c)
			}
			if res != 0 {
				t.Error("expect empty")
			}
		}
	})
}

func TestToInt64(t *testing.T) {
	type Bool bool
	type Int int
	goodCases := []struct {
		Value  any
		Result int64
	}{
		{
			true,
			1,
		}, {
			false,
			0,
		}, {
			123456,
			123456,
		}, {
			Int(123456),
			123456,
		}, {
			int8(123),
			123,
		}, {
			int16(123),
			123,
		}, {
			int32(123),
			123,
		}, {
			int64(123),
			123,
		}, {
			uint8(123),
			123,
		}, {
			uint16(123),
			123,
		}, {
			uint32(123),
			123,
		}, {
			uint64(123),
			123,
		}, {
			math.MaxInt64,
			math.MaxInt64,
		}, {
			math.MinInt64,
			math.MinInt64,
		}, {
			0.0,
			0,
		},
		{
			-123.1,
			-123,
		},
		{
			Bool(false),
			0,
		},
		{
			[]byte("12"),
			12,
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToInt64(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		"hello", time.Now(), Foo{}, uint64(math.MaxInt64 + 1), uint64(math.MaxUint64),
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToInt64(c)
			if err == nil {
				t.Error("should fail", c)
			}

			if res != 0 {
				t.Error("expect empty")
			}
		}
	})
}

func TestToInt(t *testing.T) {
	type Bool bool
	type Int int
	goodCases := []struct {
		Value  any
		Result int
	}{
		{
			true,
			1,
		}, {
			false,
			0,
		}, {
			123456,
			123456,
		}, {
			Int(123456),
			123456,
		}, {
			int8(123),
			123,
		}, {
			int16(123),
			123,
		}, {
			int32(123),
			123,
		}, {
			int64(123),
			123,
		}, {
			uint8(123),
			123,
		}, {
			uint16(123),
			123,
		}, {
			uint32(123),
			123,
		}, {
			uint64(123),
			123,
		}, {
			0.0,
			0,
		},
		{
			-123.1,
			-123,
		},
		{
			Bool(false),
			0,
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToInt(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		"hello", time.Now(), Foo{}, uint64(math.MaxInt64 + 1), uint64(math.MaxUint64),
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToInt(c)
			if err == nil {
				t.Error("should fail", c)
			}
			if res != 0 {
				t.Error("expect empty")
			}
		}
	})
}

func TestToUint64(t *testing.T) {
	type Bool bool
	type Int int
	goodCases := []struct {
		Value  any
		Result uint64
	}{
		{
			true,
			1,
		}, {
			false,
			0,
		}, {
			123456,
			123456,
		}, {
			Int(123456),
			123456,
		}, {
			int8(123),
			123,
		}, {
			int16(123),
			123,
		}, {
			int32(123),
			123,
		}, {
			int64(123),
			123,
		}, {
			uint8(123),
			123,
		}, {
			uint16(123),
			123,
		}, {
			uint32(123),
			123,
		}, {
			uint64(123),
			123,
		}, {
			math.MaxInt64,
			math.MaxInt64,
		}, {
			0.0,
			0,
		},
		{
			Bool(false),
			0,
		},
		{
			[]byte("12"),
			12,
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToUint64(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		"hello", time.Now(), Foo{}, "18446744073709551616", -1,
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToUint64(c)
			if err == nil {
				t.Error("should fail", c)
			}
			if res != 0 {
				t.Error("expect empty")
			}
		}
	})
}

func TestToIntSlice(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		l := []any{"1", 12, -13.9, json.Number("100")}
		res, err := ToIntSlice(l)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []int{1, 12, -13, 100}, res)
	})
	t.Run("nil", func(t *testing.T) {
		res, err := ToIntSlice(nil)
		if err != nil {
			t.Error(err)
		}
		if len(res) != 0 {
			t.Error("expect empty")
		}
	})
}

func TestToString(t *testing.T) {
	type String string
	goodCases := []struct {
		Value  any
		Result string
	}{
		{
			"This is a string",
			"This is a string",
		},
		{
			String("Typed string"),
			"Typed string",
		},
		{
			fmt.Errorf("error string"),
			"error string",
		},
		{
			10,
			"10",
		},
		{
			-10,
			"-10",
		},
		{
			10.2,
			"10.2",
		},
		{
			-10.2,
			"-10.2",
		},
		{
			[]byte("This is a byte slice"),
			"This is a byte slice",
		},
		{
			json.Number("10e6"),
			"10e6",
		},
	}

	t.Run("Good", func(t *testing.T) {
		for _, c := range goodCases {
			res, err := ToString(c.Value)
			if err != nil {
				t.Error(err, c.Value)
			}
			testutil.Equal(t, c.Result, res)
		}
	})

	type Foo struct{}
	badCases := []any{
		Foo{}, &Foo{},
	}
	t.Run("Bad", func(t *testing.T) {
		for _, c := range badCases {
			res, err := ToString(c)
			if err == nil {
				t.Error("should fail", c)
			}
			testutil.Equal(t, "", res)
		}
	})
}

func TestToSlice(t *testing.T) {
	t.Run("SingleString", func(t *testing.T) {
		s := "123"
		l, err := ToStringSlice(s)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []string{s}, l)
	})
	t.Run("SingleInt", func(t *testing.T) {
		s := 123
		l, err := ToStringSlice(s)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []string{fmt.Sprint(s)}, l)
	})
	t.Run("IntSlice", func(t *testing.T) {
		s := []int{123, -1, 9}
		l, err := ToStringSlice(s)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []string{"123", "-1", "9"}, l)
	})
	t.Run("MixSlice", func(t *testing.T) {
		s := []any{123, "hello", "0x123"}
		l, err := ToStringSlice(s)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []string{"123", "hello", "0x123"}, l)
	})
	t.Run("MixArray", func(t *testing.T) {
		s := [3]any{123, "hello", "0x123"}
		l, err := ToStringSlice(s)
		if err != nil {
			t.Error(err)
		}
		testutil.Equal(t, []string{"123", "hello", "0x123"}, l)
	})
}
