package xtype

type Pair[T any] struct {
	First  T
	Second T
}

type StringPair Pair[string]
type Int64Pair Pair[int64]
