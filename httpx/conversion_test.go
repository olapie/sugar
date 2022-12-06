package httpx_test

import (
	"net/http"
	"testing"

	"code.olapie.com/sugar/httpx"
	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
)

func TestToMap(t *testing.T) {
	t.Run("HeaderToMap", func(t *testing.T) {
		h := http.Header{}
		h.Set("k1", "v1")
		h.Set("k2", "v2")
		h.Add("k2", "v22")
		m := httpx.ToM(h)
		testx.Equal(t, types.M{"K1": "v1", "K2": []string{"v2", "v22"}}, m)
	})
}
