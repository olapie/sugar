package sqlx_test

import (
	"fmt"
	"testing"

	"code.olapie.com/sugar/sqlx"
)

func TestParsePSQLCompositeFields(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		res, err := sqlx.ParsePSQLCompositeFields("(abc)")
		if err != nil {
			t.Error(err)
		}

		if len(res) != 1 || res[0] != "abc" {
			t.Errorf("Expect [abc] got %v", res)
		}

		_, err = sqlx.ParsePSQLCompositeFields("(abc\")")
		if err == nil {
			t.Error("Expect error")
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		res, err := sqlx.ParsePSQLCompositeFields("(abc,123)")
		if err != nil {
			t.Error(err)
		}

		expect := fmt.Sprint([]string{"abc", "123"})
		got := fmt.Sprint(res)
		if expect != got {
			t.Errorf("Expect %s, got %s", expect, got)
		}

		res, err = sqlx.ParsePSQLCompositeFields("(abc,)")
		if err != nil {
			t.Error(err)
		}
		expect = fmt.Sprint([]string{"abc", ""})
		got = fmt.Sprint(res)
		if expect != got {
			t.Errorf("Expect %s, got %s", expect, got)
		}
	})

	t.Run("Embedded", func(t *testing.T) {
		res, err := sqlx.ParsePSQLCompositeFields("(abc,123,\"(19,20)\")")
		if err != nil {
			t.Error(err)
		}
		expect := fmt.Sprint([]string{"abc", "123", "(19,20)"})
		got := fmt.Sprint(res)
		if expect != got {
			t.Errorf("Expect %s, got %s", expect, got)
		}

		res, err = sqlx.ParsePSQLCompositeFields("(\"(19,20)\",abc,123,)")
		if err != nil {
			t.Error(err)
		}
		expect = fmt.Sprint([]string{"(19,20)", "abc", "123", ""})
		got = fmt.Sprint(res)
		if expect != got {
			t.Errorf("Expect %s, got %s", expect, got)
		}
	})

	t.Run("Quoted", func(t *testing.T) {
		res, err := sqlx.ParsePSQLCompositeFields("(\"abc\"\", \",123)")
		if err != nil {
			t.Error(err)
		}
		expect := fmt.Sprint([]string{"abc\", ", "123"})
		got := fmt.Sprint(res)
		if expect != got {
			t.Errorf("Expect %s, got %s", expect, got)
		}
	})
}
