package types_test

import (
	"testing"

	"code.olapie.com/sugar/v2/testutil"
	"code.olapie.com/sugar/v2/types"
)

func TestM_AddStruct(t *testing.T) {
	m := types.M{}
	m["id"] = 1
	m["name"] = "Smith"
	var foo struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	foo.Name = "Mike"
	foo.Age = 19
	err := m.AddStruct(foo)
	testutil.NoError(t, err)
	testutil.Equal(t, 1, m.Int("id"))
	testutil.Equal(t, 19, m.Int("age"))
	testutil.Equal(t, "Mike", m.String("name"))
}
