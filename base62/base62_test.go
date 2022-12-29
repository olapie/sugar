package base62_test

import (
	"code.olapie.com/sugar/base62"
	"code.olapie.com/sugar/testx"
	"github.com/google/uuid"
	"testing"
)

func TestEncodeToString(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		id := uuid.New()
		idStr := base62.EncodeToString(id[:])
		t.Log(idStr)
		parsed, err := base62.DecodeString(idStr)
		testx.NoError(t, err)
		testx.Equal(t, id[:], parsed)
	})
}
