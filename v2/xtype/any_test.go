package xtype_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/xtest"
	"code.olapie.com/sugar/v2/xtype"
)

func jsonString(i any) string {
	b, _ := json.Marshal(i)
	return string(b)
}

type Image struct {
	Url    string `json:"url"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Format string `json:"format"`
}

func nextImage() *Image {
	return &Image{
		Url:    "https://www.image.com/" + fmt.Sprint(time.Now().Unix()),
		Width:  rand.Int31(),
		Height: rand.Int31(),
		Format: "png",
	}
}

func TestAny(t *testing.T) {
	t.Run("AliasType", func(t *testing.T) {
		type ID int
		xtype.RegisterAnyType(ID(0))

		v := xtype.NewAny(ID(10))
		b, err := json.Marshal(v)
		xtest.NoError(t, err)

		var vv *xtype.Any
		err = json.Unmarshal(b, &vv)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("String", func(t *testing.T) {
		v := xtype.NewAny("hello")
		b, err := json.Marshal(v)
		xtest.NoError(t, err)

		var vv *xtype.Any
		err = json.Unmarshal(b, &vv)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("Array", func(t *testing.T) {
		var l []*xtype.Any
		l = append(l, xtype.NewAny("hello"))
		l = append(l, xtype.NewAny(nextImage()))
		b, err := json.Marshal(l)
		xtest.NoError(t, err)

		var ll []*xtype.Any
		err = json.Unmarshal(b, &ll)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(l), jsonString(ll))
	})
}
