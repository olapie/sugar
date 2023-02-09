package stringutil

import (
	"testing"
	"time"

	"code.olapie.com/sugar/v2/testutil"
)

func TestSquish(t *testing.T) {
	s := "\n\n\t\t1 \t2\n3  4\n\t \n5   "
	testutil.Equal(t, "1 2 3 4 5", Squish(s))
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
	SquishFields(u)
	testutil.Equal(t, "Tom", u.Name.FirstName)
	testutil.Equal(t, "Jim", u.Name.LastName)
	testutil.Equal(t, "Toronto Canada", u.Address)
}

func TestRemoveAllSpaces(t *testing.T) {
	s := "\n\ra b c  d "
	res := RemoveAllSpaces(s)
	testutil.Equal(t, "abcd", res)
}

func TestRemoveBullet(t *testing.T) {
	s := "1. hello"
	testutil.Equal(t, "hello", RemoveBullet(s))

	s = "12. hello"
	testutil.Equal(t, "hello", RemoveBullet(s))

	s = ". hello"
	testutil.Equal(t, "hello", RemoveBullet(s))

	s = "* hello"
	testutil.Equal(t, "hello", RemoveBullet(s))

	s = "*hello"
	testutil.Equal(t, "hello", RemoveBullet(s))

	s = "12.hello"
	testutil.Equal(t, "hello", RemoveBullet(s))
}
