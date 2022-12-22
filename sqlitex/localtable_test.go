package sqlitex_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"code.olapie.com/sugar/sqlitex"
	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
	"github.com/google/uuid"
)

type localTableItem struct {
	ID     int64
	Text   string
	Number float64
	List   []int
}

func setupLocalTable(t testing.TB) *sqlitex.LocalTable[*localTableItem] {
	filename := "testdata/localtable" + types.NextID().Pretty() + ".db"
	db, err := sqlitex.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() {
			db.Close()
			os.Remove(filename)
		})
	return sqlitex.NewLocalTable[*localTableItem](db, "", nil)
}

func newLocalTableItem() *localTableItem {
	return &localTableItem{
		ID:     rand.Int63(),
		Text:   time.Now().String(),
		Number: rand.Float64(),
		List:   []int{rand.Int(), rand.Int()},
	}
}

func TestLocalTable_SaveRemote(t *testing.T) {
	ctx := context.TODO()
	table := setupLocalTable(t)
	item := newLocalTableItem()
	localID := uuid.NewString()
	err := table.SaveRemote(ctx, localID, item, time.Now().Unix())
	testx.NoError(t, err)
	record, err := table.Get(ctx, localID)
	testx.NoError(t, err)
	testx.Equal(t, item, record)

	item.Text = time.Now().String() + "new"
	err = table.SaveRemote(ctx, localID, item, time.Now().Unix())
	testx.NoError(t, err)
	record, err = table.Get(ctx, localID)
	testx.NoError(t, err)
	testx.Equal(t, item, record)

	records, err := table.ListLocals(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 0, len(records))

	records, err = table.ListDeletions(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 0, len(records))

	records, err = table.ListUpdates(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 0, len(records))

	records, err = table.ListRemotes(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 1, len(records))

	err = table.Delete(ctx, localID)
	testx.NoError(t, err)
	records, err = table.ListRemotes(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 0, len(records))

	records, err = table.ListDeletions(ctx)
	testx.NoError(t, err)
	testx.Equal(t, 1, len(records))
	testx.Equal(t, item, records[0])

}

func TestLocalTable_SaveLocal(t *testing.T) {
	ctx := context.TODO()
	table := setupLocalTable(t)
	t.Run("SyncedRemote", func(t *testing.T) {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, item)
		testx.NoError(t, err)
		record, err := table.Get(ctx, localID)
		testx.NoError(t, err)
		testx.Equal(t, item, record)

		item.Text = time.Now().String() + "new"
		err = table.SaveLocal(ctx, localID, item)
		testx.NoError(t, err)
		record, err = table.Get(ctx, localID)
		testx.NoError(t, err)
		testx.Equal(t, item, record)

		locals, err := table.ListLocals(ctx)
		testx.NoError(t, err)
		testx.Equal(t, 1, len(locals))

		remoteItem := newLocalTableItem()
		err = table.SaveRemote(ctx, localID, remoteItem, time.Now().Unix())
		testx.NoError(t, err)

		locals, err = table.ListLocals(ctx)
		testx.NoError(t, err)
		testx.Equal(t, 0, len(locals))
	})

	t.Run("DeleteLocal", func(t *testing.T) {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, item)
		testx.NoError(t, err)
		record, err := table.Get(ctx, localID)
		testx.NoError(t, err)
		testx.Equal(t, item, record)

		item.Text = time.Now().String() + "new"
		err = table.SaveLocal(ctx, localID, item)
		testx.NoError(t, err)
		record, err = table.Get(ctx, localID)
		testx.NoError(t, err)
		testx.Equal(t, item, record)

		locals, err := table.ListLocals(ctx)
		testx.NoError(t, err)
		testx.Equal(t, 1, len(locals))

		err = table.Delete(ctx, localID)
		testx.NoError(t, err)

		locals, err = table.ListLocals(ctx)
		testx.NoError(t, err)
		testx.Equal(t, 0, len(locals))

		deletes, err := table.ListDeletions(ctx)
		testx.NoError(t, err)
		testx.Equal(t, 0, len(deletes))
	})
}

func BenchmarkLocalTable_SaveLocal(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, item)
		testx.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, item, time.Now().Unix())
		testx.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		testx.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		table.SaveLocal(ctx, localID, item)
	}

	for i := 0; i < b.N; i++ {
		table.Get(ctx, ids[i%len(ids)])
	}
}

func BenchmarkLocalTable_Get(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, item)
		testx.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, item, time.Now().Unix())
		testx.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		testx.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		table.Get(ctx, ids[i%len(ids)])
	}
}
