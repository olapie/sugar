package jsonutil_test

import (
	"testing"

	"code.olapie.com/sugar/v2/jsonutil"
)

func TestJSONExample(t *testing.T) {
	type Embed struct {
		Field1 *string
		Field2 int
		// Field3 *time.Time
		Field4 []int
		Field5 []*string
	}

	type Foo struct {
		Field1 *Embed
		Field2 Embed
		Field3 bool
		List   []*Embed
	}

	var foo *Foo

	t.Log(jsonutil.ExampleStringOf(foo))
}
