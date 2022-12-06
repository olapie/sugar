package errorx_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	var err = errorx.New("List of xerror")
	var b strings.Builder
	b.WriteString("List of xerror")
	for i := 1; i < 10; i++ {
		s := fmt.Errorf("#%d error", i)
		b.WriteString("; ")
		b.WriteString(s.Error())
		err = errorx.Append(err, s)
		t.Log(b.String())
		if err.Error() != b.String() {
			t.Fatalf("Expect %s instead of %s", b.String(), err.Error())
		}
	}
}
