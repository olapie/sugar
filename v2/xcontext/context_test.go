package xcontext_test

import (
	"context"
	"testing"

	"code.olapie.com/sugar/v2/xcontext"
	"code.olapie.com/sugar/v2/xtest"
	"code.olapie.com/sugar/v2/xtype"
)

func TestGetLogin(t *testing.T) {
	t.Run("int64ToID", func(t *testing.T) {
		ctx := xcontext.WithLogin(context.TODO(), int64(1))
		id := xcontext.GetLogin[xtype.ID](ctx)
		xtest.Equal(t, xtype.ID(1), id)
	})

	t.Run("int64ToString", func(t *testing.T) {
		ctx := xcontext.WithLogin(context.TODO(), int64(1))
		id := xcontext.GetLogin[string](ctx)
		xtest.Equal(t, "1", id)
	})
}
