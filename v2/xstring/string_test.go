package xstring_test

import (
	"testing"
	"time"

	"code.olapie.com/sugar/xstring"
	"code.olapie.com/sugar/xtest"
)

func TestSquish(t *testing.T) {
	s := "\n\n\t\t1 \t2\n3  4\n\t \n5   "
	xtest.Equal(t, "1 2 3 4 5", xstring.Squish(s))
}

func TestSquishFields(t *testing.T) {
	type FullName struct {
		FirstName string
		LastName  string
	}

	type User struct {
		ID        int64
		Name      FullName
		Address   string
		CreatedAt time.Time
	}

	u := &User{
		ID: 1,
		Name: FullName{
			FirstName: "Tom  ",
			LastName:  "  Jim ",
		},
		Address: "\n\n \t Toronto   Canada   ",
	}
	xstring.SquishFields(u)
	xtest.Equal(t, "Tom", u.Name.FirstName)
	xtest.Equal(t, "Jim", u.Name.LastName)
	xtest.Equal(t, "Toronto Canada", u.Address)
}

func TestRemoveAllSpaces(t *testing.T) {
	s := "\n\ra b c  d "
	res := xstring.RemoveAllSpaces(s)
	xtest.Equal(t, "abcd", res)
}

func TestRemoveBullet(t *testing.T) {
	s := "1. hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))

	s = "12. hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))

	s = ". hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))

	s = "* hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))

	s = "*hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))

	s = "12.hello"
	xtest.Equal(t, "hello", xstring.RemoveBullet(s))
}
