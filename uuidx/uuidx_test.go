package uuidx_test

import (
	"testing"

	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/uuidx"
	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := uuid.New()
		short := uuidx.ShortString(id)
		t.Log(short)
		t.Log(id.String())
		id2 := uuidx.FromShortString(short)
		testx.Equal(t, id, id2.UUID)
		testx.Equal(t, id.String(), id2.UUID.String())
	}
}
