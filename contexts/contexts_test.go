package contexts_test

import (
	"code.olapie.com/sugar/contexts"
	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
	"context"
	"testing"
)

func TestGetLogin(t *testing.T) {
	t.Run("int64ToID", func(t *testing.T) {
		ctx := contexts.WithLogin(context.TODO(), int64(1))
		id := contexts.GetLogin[types.ID](ctx)
		testx.Equal(t, types.ID(1), id)
	})

	t.Run("int64ToString", func(t *testing.T) {
		ctx := contexts.WithLogin(context.TODO(), int64(1))
		id := contexts.GetLogin[string](ctx)
		testx.Equal(t, "1", id)
	})
}
