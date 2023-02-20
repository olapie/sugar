package httpkit

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/hashutil"
)

func Sign(req *http.Request) {
	t := time.Now().Unix()
	var b [40]byte
	binary.BigEndian.PutUint64(b[:], uint64(t))
	hash := hashutil.Hash32(fmt.Sprint(t))
	copy(b[8:], hash[:])
	sign := base62.EncodeToString(b[:])
	req.Header.Set(keySignature, sign)
}

func Verify(req *http.Request) bool {
	sign := GetHeader(req.Header, keySignature)
	if sign == "" {
		fmt.Println("missing", keySignature)
		return false
	}

	var b [40]byte
	decoded, err := base62.DecodeString(sign)
	if err != nil {
		fmt.Println("invalid", keySignature, err)
		return false
	}
	copy(b[40-len(decoded):], decoded)
	t := int64(binary.BigEndian.Uint64(b[:]))
	elapsed := time.Now().Unix() - t
	if elapsed < -3 || elapsed > 10 {
		fmt.Println("invalid timestamp", t, elapsed)
		return false
	}
	hash := hashutil.Hash32(fmt.Sprint(t))
	equal := bytes.Equal(b[8:], hash[:])
	return equal
}
