package bytex_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/bytex"
	"code.olapie.com/sugar/testx"
)

type customByteSlice []byte

func TestMarshalCustomBytesType(t *testing.T) {
	var input customByteSlice = []byte(time.Now().String())
	output, err := bytex.Marshal(input)
	testx.NoError(t, err)
	testx.Equal(t, []byte(input), output)
}

type jsonObject struct {
	ID   int64
	Text string
}

func (o *jsonObject) MarshalJSON() ([]byte, error) {
	type alias jsonObject
	obj := (*alias)(o)
	return json.Marshal(obj)
}

func (o *jsonObject) UnmarshalJSON(data []byte) error {
	type alias jsonObject
	var obj alias
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	*o = jsonObject(obj)
	return nil
}

func TestMarshalJSON(t *testing.T) {
	o := jsonObject{ID: rand.Int63(), Text: time.Now().String()}
	data, err := bytex.Marshal(&o)
	testx.NoError(t, err)
	t.Log(string(data))

	var o2 jsonObject
	err = bytex.Unmarshal(data, &o2)
	testx.NoError(t, err)
	t.Log(o2.ID, o2.Text)
}
