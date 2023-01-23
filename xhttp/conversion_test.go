package xhttp_test

import (
	"net/http"
	"testing"

	"code.olapie.com/sugar/v2/xhttp"
	"code.olapie.com/sugar/v2/xtest"
	"code.olapie.com/sugar/v2/xtype"
)

func TestToMap(t *testing.T) {
	t.Run("HeaderToMap", func(t *testing.T) {
		h := http.Header{}
		h.Set("k1", "v1")
		h.Set("k2", "v2")
		h.Add("k2", "v22")
		m := xhttp.ToM(h)
		xtest.Equal(t, xtype.M{"K1": "v1", "K2": []string{"v2", "v22"}}, m)
	})
}
