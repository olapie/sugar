package types_test

import (
	"encoding/json"
	"sort"
	"testing"

	"code.olapie.com/sugar/v2/testutil"
	"code.olapie.com/sugar/v2/types"
)

func TestMarshalJSON(t *testing.T) {
	s1 := types.NewSet[int](10)
	a0 := []int{1, 2, 3, 5, 9}
	for _, v := range a0 {
		s1.Add(v)
	}
	d1, err := s1.MarshalJSON()
	testutil.NoError(t, err)
	var s2 *types.Set[int]
	err = json.Unmarshal(d1, &s2)
	testutil.NoError(t, err)
	a1 := s1.Slice()
	a2 := s2.Slice()
	sort.IntSlice(a1).Sort()
	sort.IntSlice(a2).Sort()
	testutil.Equal(t, a0, a1)
	testutil.Equal(t, a1, a2)
}
