package types_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/testutil"
	"code.olapie.com/sugar/v2/types"
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
		types.RegisterAnyType(ID(0))

		v := types.NewAny(ID(10))
		b, err := json.Marshal(v)
		testutil.NoError(t, err)

		var vv *types.Any
		err = json.Unmarshal(b, &vv)
		testutil.NoError(t, err)

		testutil.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("String", func(t *testing.T) {
		v := types.NewAny("hello")
		b, err := json.Marshal(v)
		testutil.NoError(t, err)

		var vv *types.Any
		err = json.Unmarshal(b, &vv)
		testutil.NoError(t, err)

		testutil.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("Array", func(t *testing.T) {
		var l []*types.Any
		l = append(l, types.NewAny("hello"))
		l = append(l, types.NewAny(nextImage()))
		b, err := json.Marshal(l)
		testutil.NoError(t, err)

		var ll []*types.Any
		err = json.Unmarshal(b, &ll)
		testutil.NoError(t, err)

		testutil.Equal(t, jsonString(l), jsonString(ll))
	})
}
