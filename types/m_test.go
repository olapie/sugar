package types_test

import (
	"testing"

	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
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
	testx.NoError(t, err)
	testx.Equal(t, 1, m.Int("id"))
	testx.Equal(t, 19, m.Int("age"))
	testx.Equal(t, "Mike", m.String("name"))
}
