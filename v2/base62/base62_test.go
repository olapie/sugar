package base62_test

import (
	"code.olapie.com/sugar/base62"
	"code.olapie.com/sugar/xtest"
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestEncodeToString(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		id := uuid.New()
		idStr := base62.EncodeToString(id[:])
		t.Log(idStr)
		t.Log(base64.StdEncoding.EncodeToString(id[:]))
		t.Log(strings.ReplaceAll(id.String(), "-", ""))
		t.Log(id.String())
		parsed, err := base62.DecodeString(idStr)
		xtest.NoError(t, err)
		xtest.Equal(t, id[:], parsed)
	})
}
