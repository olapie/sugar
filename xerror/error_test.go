package xerror

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestErrorString(t *testing.T) {
	err := New(http.StatusConflict, "duplicate nickname")
	t.Log(err.Error())
	if err.Error() != "duplicate nickname" {
		t.Fail()
	}
}

func TestEmbedError(t *testing.T) {
	err := New(http.StatusNotFound, "token")
	err = New(http.StatusUnauthorized, err.Error())
	t.Log(err)
	if err.Error() != "token" {
		t.Fail()
	}
}

func TestJSON(t *testing.T) {
	text, err := json.Marshal(New(http.StatusBadRequest, "test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(text))

	var e *Error
	err = json.Unmarshal(text, &e)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e.Code(), e.message)

	obj := &errorJSONObject{
		Code:    e.code,
		Message: e.message,
	}
	text, _ = json.Marshal(obj)
	t.Log(string(text))
}
