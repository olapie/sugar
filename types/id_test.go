package types

import (
	"math"
	"testing"
)

func TestID(t *testing.T) {
	for i := 0; i < 256; i++ {
		id := NextID()
		t.Logf("%d %0X %s", id, id, id.Short())
		//time.Sleep(time.Millisecond * 1)
	}

	var id ID = 123
	if id.Short() != "1z" {
		t.Log(id.Short())
		t.FailNow()
	}

	id = 62
	if id.Short() != "10" {
		t.Log(id.Short())
		t.FailNow()
	}

	id = math.MaxInt64
	if i, _ := NewIDFromString(id.Short(), ShortIDFormat); i != id {
		t.Log(id.Short(), i)
		t.FailNow()
	}
}

func TestID_Pretty(t *testing.T) {
	for i := 0; i < 256; i++ {
		id := NextID()
		t.Logf("%d %0X %s", id, id, id.Pretty())
		//time.Sleep(time.Millisecond * 1)
	}

	var id ID = 123
	if id.Pretty() != "4M" {
		t.Log(id.Pretty())
		t.FailNow()
	}

	id = 34
	if id.Pretty() != "21" {
		t.Log(id.Pretty())
		t.FailNow()
	}

	id = math.MaxInt64
	if i, _ := NewIDFromString(id.Pretty(), PrettyIDFormat); i != id {
		t.Log(id.Pretty(), i)
		t.FailNow()
	}
}
