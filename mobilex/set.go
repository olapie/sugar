package mobilex

import (
	"code.olapie.com/sugar/jsonx"
	"code.olapie.com/sugar/types"
)

type IntSet struct {
	types.Set[int]
}

func (s *IntSet) String() string {
	return jsonx.ToString(s.Slice())
}

type Int64Set struct {
	types.Set[int64]
}

func (s *Int64Set) String() string {
	return jsonx.ToString(s.Slice())
}

type Float64Set struct {
	types.Set[float64]
}

func (s *Float64Set) String() string {
	return jsonx.ToString(s.Slice())
}

type StringSet struct {
	types.Set[string]
}

func NewStringSet() *StringSet {
	return new(StringSet)
}
