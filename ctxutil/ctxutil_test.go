package ctxutil_test

import (
	"context"
	"testing"

	"code.olapie.com/sugar/ctxutil"
	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
)

func TestGetLogin(t *testing.T) {
	t.Run("int64ToID", func(t *testing.T) {
		ctx := ctxutil.WithLogin(context.TODO(), int64(1))
		id := ctxutil.GetLogin[types.ID](ctx)
		testx.Equal(t, types.ID(1), id)
	})

	t.Run("int64ToString", func(t *testing.T) {
		ctx := ctxutil.WithLogin(context.TODO(), int64(1))
		id := ctxutil.GetLogin[string](ctx)
		testx.Equal(t, "1", id)
	})
}
