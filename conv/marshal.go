package conv

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"code.olapie.com/sugar/v2/rt"
	"google.golang.org/protobuf/proto"
)

func Marshal(i any) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}

	var data []byte
	srcType := reflect.TypeOf(i)
	dstType := reflect.TypeOf(data)
	if srcType.AssignableTo(dstType) {
		reflect.ValueOf(&data).Elem().Set(reflect.ValueOf(i))
		return data, nil
	}

	if srcType.ConvertibleTo(dstType) {
		reflect.ValueOf(&data).Elem().Set(reflect.ValueOf(i).Convert(dstType))
		return data, nil
	}

	if m, ok := i.(encoding.BinaryMarshaler); ok {
		return m.MarshalBinary()
	}

	if m, ok := i.(encoding.TextMarshaler); ok {
		return m.MarshalText()
	}

	if m, ok := i.(json.Marshaler); ok {
		return m.MarshalJSON()
	}

	if m, ok := i.(gob.GobEncoder); ok {
		return m.GobEncode()
	}
	if m, ok := i.(proto.Message); ok {
		return proto.Marshal(m)
	}

	return nil, errors.New("cannot convert ")
}

func Unmarshal(data []byte, i any) (err error) {
	if reflect.ValueOf(i).Kind() != reflect.Pointer {
		return fmt.Errorf("cannot unmarshal to non pointer type: %T", i)
	}

	if p, ok := i.(*[]byte); ok {
		*p = data
		return nil
	}

	srcType := reflect.TypeOf(data)
	dstType := reflect.TypeOf(i).Elem()
	if srcType.AssignableTo(dstType) {
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	if srcType.ConvertibleTo(dstType) {
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(data).Convert(dstType))
		return nil
	}

	// i is a pointer
	// v is pointer of the same type
	v := rt.DeepNew(reflect.TypeOf(i).Elem())
	defer func() {
		if err == nil {
			// assign v to i
			// i is a parameter, it cannot be set, the value it points to can be set
			// assign v.Elem to i.Elem
			reflect.ValueOf(i).Elem().Set(v.Elem())
		}
	}()

	for p := v; p.Kind() == reflect.Pointer && !p.IsNil(); p = p.Elem() {
		if u, ok := p.Interface().(encoding.BinaryUnmarshaler); ok {
			return u.UnmarshalBinary(data)
		}

		if u, ok := p.Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText(data)
		}

		if u, ok := p.Interface().(json.Unmarshaler); ok {
			return u.UnmarshalJSON(data)
		}

		if d, ok := p.Interface().(gob.GobDecoder); ok {
			return d.GobDecode(data)
		}

		if m, ok := p.Interface().(proto.Message); ok {
			return proto.Unmarshal(data, m)
		}
	}

	return fmt.Errorf("cannot unmarshal into: %T", i)
}
