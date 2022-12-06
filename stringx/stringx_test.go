package stringx_test

import (
	"testing"
	"time"

	"code.olapie.com/sugar/stringx"
	"code.olapie.com/sugar/testx"
)

func TestSquish(t *testing.T) {
	s := "\n\n\t\t1 \t2\n3  4\n\t \n5   "
	testx.Equal(t, "1 2 3 4 5", stringx.Squish(s))
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
	stringx.SquishFields(u)
	testx.Equal(t, "Tom", u.Name.FirstName)
	testx.Equal(t, "Jim", u.Name.LastName)
	testx.Equal(t, "Toronto Canada", u.Address)
}

func TestRemoveAllSpaces(t *testing.T) {
	s := "\n\ra b c  d "
	res := stringx.RemoveAllSpaces(s)
	testx.Equal(t, "abcd", res)
}

func TestRemoveBullet(t *testing.T) {
	s := "1. hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))

	s = "12. hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))

	s = ". hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))

	s = "* hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))

	s = "*hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))

	s = "12.hello"
	testx.Equal(t, "hello", stringx.RemoveBullet(s))
}
