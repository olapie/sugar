package sqlitex_test

import (
	"database/sql"
	"errors"
	"testing"

	"code.olapie.com/sugar/sqlitex"
	"code.olapie.com/sugar/testx"
	"code.olapie.com/sugar/types"
	_ "github.com/mattn/go-sqlite3"
)

func createTable[K sqlitex.SimpleKey, R sqlitex.SimpleTableRecord[K]](t *testing.T, name string) *sqlitex.SimpleTable[K, R] {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Error(err)
	}

	tbl, err := sqlitex.NewSimpleTable[K, R](db, name)
	testx.NoError(t, err)
	return tbl
}

type IntItem struct {
	ID    int64
	Name  string
	Score float64
}

func (i *IntItem) PrimaryKey() int64 {
	return i.ID
}

func newIntItem() *IntItem {
	return &IntItem{
		ID:    types.RandomID().Int(),
		Name:  types.RandomID().Pretty(),
		Score: float64(types.RandomID()) / float64(3),
	}
}

type StringItem struct {
	ID    string
	Name  string
	Score float64
}

func (i *StringItem) PrimaryKey() string {
	return i.ID
}

func newStringItem() *StringItem {
	return &StringItem{
		ID:    types.RandomID().Pretty(),
		Name:  types.RandomID().Pretty(),
		Score: float64(types.RandomID()) / float64(3),
	}
}

func TestIntTable(t *testing.T) {
	tbl := createTable[int64, *IntItem](t, "tbl"+types.RandomID().Pretty())
	var items []*IntItem
	item := newIntItem()
	items = append(items, item)
	err := tbl.Insert(item)
	testx.NoError(t, err)
	v, err := tbl.Get(item.ID)
	testx.NoError(t, err)
	testx.Equal(t, item, v)

	item = newIntItem()
	item.ID = items[0].ID + 1
	err = tbl.Insert(item)
	testx.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	testx.NoError(t, err)
	testx.True(t, len(l) != 0)
	testx.Equal(t, items, l)

	l, err = tbl.ListGreaterThan(item.ID, 10)
	testx.NoError(t, err)
	testx.True(t, len(l) == 0)

	l, err = tbl.ListLessThan(item.ID+1, 10)
	testx.NoError(t, err)
	testx.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	testx.True(t, l[0].ID < l[1].ID)

	err = tbl.Delete(item.ID)
	testx.NoError(t, err)

	v, err = tbl.Get(item.ID)
	testx.Error(t, err)
	testx.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}

func TestStringTable(t *testing.T) {
	tbl := createTable[string, *StringItem](t, "tbl"+types.RandomID().Pretty())
	var items []*StringItem
	item := newStringItem()
	items = append(items, item)
	err := tbl.Insert(item)
	testx.NoError(t, err)
	v, err := tbl.Get(item.ID)
	testx.NoError(t, err)
	testx.Equal(t, item, v)

	item = newStringItem()
	err = tbl.Insert(item)
	testx.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	testx.NoError(t, err)
	testx.True(t, len(l) != 0)
	testx.Equal(t, items, l)

	l, err = tbl.ListGreaterThan("a", 10)
	testx.NoError(t, err)
	testx.True(t, len(l) == 0)

	l, err = tbl.ListLessThan("Z", 10)
	testx.NoError(t, err)
	testx.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	testx.True(t, l[0].ID < l[1].ID)

	err = tbl.Delete(item.ID)
	testx.NoError(t, err)

	v, err = tbl.Get(item.ID)
	testx.Error(t, err)
	testx.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}
