package sqlite

import (
	"database/sql"
	"errors"
	"testing"

	"code.olapie.com/sugar/v2/xtest"
	"code.olapie.com/sugar/v2/xtype"
	_ "github.com/mattn/go-sqlite3"
)

func createTable[K SimpleKey, R SimpleTableRecord[K]](t *testing.T, name string) *SimpleTable[K, R] {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Error(err)
	}

	tbl, err := NewSimpleTable[K, R](db, name)
	xtest.NoError(t, err)
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
		ID:    xtype.RandomID().Int(),
		Name:  xtype.RandomID().Pretty(),
		Score: float64(xtype.RandomID()) / float64(3),
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
		ID:    xtype.RandomID().Pretty(),
		Name:  xtype.RandomID().Pretty(),
		Score: float64(xtype.RandomID()) / float64(3),
	}
}

func TestIntTable(t *testing.T) {
	tbl := createTable[int64, *IntItem](t, "tbl"+xtype.RandomID().Pretty())
	var items []*IntItem
	item := newIntItem()
	items = append(items, item)
	err := tbl.Insert(item)
	xtest.NoError(t, err)
	v, err := tbl.Get(item.ID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, v)

	item = newIntItem()
	item.ID = items[0].ID + 1
	err = tbl.Insert(item)
	xtest.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	xtest.NoError(t, err)
	xtest.True(t, len(l) != 0)
	xtest.Equal(t, items, l)

	l, err = tbl.ListGreaterThan(item.ID, 10)
	xtest.NoError(t, err)
	xtest.True(t, len(l) == 0)

	l, err = tbl.ListLessThan(item.ID+1, 10)
	xtest.NoError(t, err)
	xtest.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	xtest.True(t, l[0].ID < l[1].ID)

	err = tbl.Delete(item.ID)
	xtest.NoError(t, err)

	v, err = tbl.Get(item.ID)
	xtest.Error(t, err)
	xtest.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}

func TestStringTable(t *testing.T) {
	tbl := createTable[string, *StringItem](t, "tbl"+xtype.RandomID().Pretty())
	var items []*StringItem
	item := newStringItem()
	items = append(items, item)
	err := tbl.Insert(item)
	xtest.NoError(t, err)
	v, err := tbl.Get(item.ID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, v)

	item = newStringItem()
	err = tbl.Insert(item)
	xtest.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	xtest.NoError(t, err)
	xtest.True(t, len(l) != 0)
	xtest.Equal(t, items, l)

	l, err = tbl.ListGreaterThan("a", 10)
	xtest.NoError(t, err)
	xtest.True(t, len(l) == 0)

	l, err = tbl.ListLessThan("Z", 10)
	xtest.NoError(t, err)
	xtest.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	xtest.True(t, l[0].ID < l[1].ID)

	err = tbl.Delete(item.ID)
	xtest.NoError(t, err)

	v, err = tbl.Get(item.ID)
	xtest.Error(t, err)
	xtest.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}
