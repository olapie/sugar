package xtype

import (
	"code.olapie.com/sugar/conv"
)

type ByteUnit int64

const (
	_           = iota
	KB ByteUnit = 1 << (10 * iota) // 1 << (10*1)
	MB                             // 1 << (10*2)
	GB                             // 1 << (10*3)
	TB                             // 1 << (10*4)
	PB                             // 1 << (10*5)
)

func (b ByteUnit) HumanReadable() string {
	return conv.SizeToHumanReadable(int64(b))
}

func (b ByteUnit) Int64() int64 {
	return int64(b)
}
