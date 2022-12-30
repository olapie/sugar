package xtype

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"code.olapie.com/sugar/conv"
)

const prettyTableSize = 34

var prettyTable = [prettyTableSize]byte{
	'1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G',
	'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'P', 'Q',
	'R', 'S', 'T',
	'U', 'V', 'W',
	'X', 'Y', 'Z'}

type IDFormat int

const (
	ShortIDFormat IDFormat = iota
	PrettyIDFormat
)

type ID int64

// Int converts ID into int64. Just make it easier to edit code
func (i ID) Int() int64 {
	return int64(i)
}

// Short returns a short representation of id
func (i ID) Short() string {
	if i < 0 {
		panic("invalid id")
	}
	var bytes [16]byte
	k := int64(i)
	n := 15
	for {
		j := k % 62
		switch {
		case j <= 9:
			bytes[n] = byte('0' + j)
		case j <= 35:
			bytes[n] = byte('A' + j - 10)
		default:
			bytes[n] = byte('a' + j - 36)
		}
		k /= 62
		if k == 0 {
			return string(bytes[n:])
		}
		n--
	}
}

// Pretty returns a incasesensitive pretty representation of id
func (i ID) Pretty() string {
	if i < 0 {
		panic("invalid id")
	}
	var bytes [16]byte
	k := int64(i)
	n := 15

	for {
		bytes[n] = prettyTable[k%prettyTableSize]
		k /= prettyTableSize
		if k == 0 {
			return string(bytes[n:])
		}
		n--
	}
}

func NewIDFromString(s string, f IDFormat) (ID, error) {
	switch f {
	case ShortIDFormat:
		return parseShortID(s)
	case PrettyIDFormat:
		return parsePrettyID(s)
	default:
		return 0, errors.New("invalid format")
	}
}

// IDFromString parse id from string
// if failed, return 0
func IDFromString(s string) ID {
	if id, err := conv.ToInt64(s); err == nil {
		return ID(id)
	}

	if id, err := parseShortID(s); err == nil {
		return id
	}

	if id, err := parsePrettyID(s); err == nil {
		return id
	}

	return 0
}

func (i ID) Salt(v string) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%s%d", v, i)))
	return hex.EncodeToString(sum[:])
}

func (i ID) IsValid() bool {
	return i > 0
}

func parseShortID(s string) (ID, error) {
	if len(s) == 0 {
		return 0, errors.New("parse error")
	}

	var bytes = []byte(s)
	var k int64
	var v int64
	for _, b := range bytes {
		switch {
		case b >= '0' && b <= '9':
			v = int64(b - '0')
		case b >= 'A' && b <= 'Z':
			v = int64(10 + b - 'A')
		case b >= 'a' && b <= 'z':
			v = int64(36 + b - 'a')
		default:
			return 0, errors.New("parse error")
		}
		k = k*62 + v
	}
	return ID(k), nil
}

func parsePrettyID(s string) (ID, error) {
	if len(s) == 0 {
		return 0, errors.New("parse error")
	}

	s = strings.ToUpper(s)
	var bytes = []byte(s)
	var k int64
	for _, b := range bytes {
		i := searchPrettyTable(b)
		if i <= 0 {
			return 0, errors.New("parse error")
		}
		k = k*prettyTableSize + int64(i)
	}
	return ID(k), nil
}

func searchPrettyTable(v byte) int {
	left := 0
	right := prettyTableSize - 1
	for right >= left {
		mid := (left + right) / 2
		if prettyTable[mid] == v {
			return mid
		} else if prettyTable[mid] > v {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return -1
}

// ------------------------------
// IDGenerator
type IDGenerator interface {
	NextID() ID
}

type NextNumber interface {
	NextNumber() int64
}

type SnakeIDGenerator struct {
	seqSize   uint
	shardSize uint

	clock    NextNumber
	sharding NextNumber

	counter int64
}

func NewSnakeIDGenerator(shardBitsSize, seqBitsSize uint, clock, sharding NextNumber) *SnakeIDGenerator {
	if seqBitsSize < 1 || seqBitsSize > 16 {
		panic("seqBitsSize should be [1,16]")
	}

	if clock == nil {
		panic("clock is nil")
	}

	if shardBitsSize > 8 {
		panic("shardBitsSize should be [0,8]")
	}

	if shardBitsSize > 0 && sharding == nil {
		panic("sharding is nil")
	}

	if shardBitsSize+seqBitsSize >= 20 {
		panic("shardBitsSize + seqBitsSize should be less than 20")
	}

	return &SnakeIDGenerator{
		seqSize:   seqBitsSize,
		shardSize: shardBitsSize,
		clock:     clock,
		sharding:  sharding,
	}
}

func (g *SnakeIDGenerator) NextID() ID {
	id := g.clock.NextNumber() << (g.seqSize + g.shardSize)
	if g.shardSize > 0 {
		id |= (g.sharding.NextNumber() % (1 << g.shardSize)) << g.seqSize
	}
	n := atomic.AddInt64(&g.counter, 1)
	id |= n % (1 << g.seqSize)
	return ID(id)
}

// Most language's JSON decoders decode number into double if type isn't explicitly specified.
// The maximum integer part of double is 2^53，so it'd better to control id bits size less than 53
// id is made of time, shard and seq
// Putting the time at the beginning can ensure the id unique and increasing in case increase shard or seq bits size in the future
var (
	epoch                   = time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)
	idGenerator IDGenerator = NewSnakeIDGenerator(0, 6, nextMilliseconds, nil)
)

type NextNumberFunc func() int64

func (f NextNumberFunc) NextNumber() int64 {
	return f()
}

var nextMilliseconds NextNumberFunc = func() int64 {
	return time.Since(epoch).Nanoseconds() / 1e6
}

func NextID() ID {
	return idGenerator.NextID()
}

func SetIDGenerator(g IDGenerator) {
	idGenerator = g
}

func RandomID() ID {
	return ID(rand.Int63())
}
