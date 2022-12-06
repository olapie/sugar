package httpkit_test

import (
	"code.olapie.com/sugar/testx"
	"net/http"
	"testing"

	"code.olapie.com/sugar/httpx"
	"code.olapie.com/sugar/types"
)

func TestToMap(t *testing.T) {
	t.Run("HeaderToMap", func(t *testing.T) {
		h := http.Header{}
		h.Set("k1", "v1")
		h.Set("k2", "v2")
		h.Add("k2", "v22")
		m := httpkit.ToM(h)
		testx.Equal(t, types.M{"K1": "v1", "K2": []string{"v2", "v22"}}, m)
	})
}