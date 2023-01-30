package xerror

import (
	"encoding/json"
	"testing"
)

func TestErrorString(t *testing.T) {
	err := Conflict("duplicate nickname")
	t.Log(err.Error())
	if err.Error() != "duplicate nickname" {
		t.Fail()
	}
}

func TestEmbedError(t *testing.T) {
	err := NotFound("token")
	err = Unauthorized(err.Error())
	t.Log(err)
	if err.Error() != "token" {
		t.Fail()
	}
}

func TestJSON(t *testing.T) {
	text, err := json.Marshal(BadRequest("test"))
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
}
