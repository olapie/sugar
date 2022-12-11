package sqlitex_test

import (
	"code.olapie.com/sugar/sqlitex"
	"code.olapie.com/sugar/testx"
	"context"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"testing"
	"time"
)

type localTableItem struct {
	ID     int64
	Text   string
	Number float64
	List   []int
}

func setupLocalTable(t *testing.T) *sqlitex.LocalTable[*localTableItem] {
	filename := "testdata/localtable.db"
	t.Cleanup(
		func() {
			os.Remove(filename)
		})
	db, err := sqlitex.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
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
}
