package jsonx

import (
	"encoding/json"
	"os"
	"reflect"

	"code.olapie.com/sugar/rtx"
)

func ToString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func ToBytes(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func ExampleStringOf(i any) string {
	v := rtx.DeepNew(reflect.TypeOf(i))
	return ToString(v.Interface())
}

func UnmarshalFile(filename string, v any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func MarshalToFile(v any, filename string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
