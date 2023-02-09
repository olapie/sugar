package base62_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/testutil"
	"github.com/google/uuid"
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
		testutil.NoError(t, err)
		testutil.Equal(t, id[:], parsed)
	})
}
