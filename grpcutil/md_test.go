package grpcutil

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/hashutil"

	"google.golang.org/grpc/metadata"
)

func TestSign(t *testing.T) {
	md := make(metadata.MD)
	Sign(md)
	t.Log(md)
	if !Verify(md, 1) {
		t.FailNow()
	}

	ts := time.Now().Unix() - 20
	var b [36]byte
	binary.BigEndian.PutUint64(b[:], uint64(ts))
	hash := hashutil.Hash32(fmt.Sprintf("%x", ts))
	copy(b[4:], hash[:])
	sign := base62.EncodeToString(b[:])
	md.Set(keySignature, sign)
	if Verify(md, 1) {
		t.FailNow()
	}
}
