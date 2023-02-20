package httpkit

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"time"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/hashutil"
)

func Sign(req *http.Request) {
	t := time.Now().Unix()
	var b [41]byte
	b[0] = 1
	binary.BigEndian.PutUint64(b[:], uint64(t))
	clientID := GetClientID(req.Header)
	traceID := GetTraceID(req.Header)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	copy(b[9:], hash[:])
	sign := base62.EncodeToString(b[:])
	req.Header.Set(keySignature, sign)
}

func Verify(req *http.Request, delaySeconds int) bool {
	sign := GetHeader(req.Header, keySignature)
	if sign == "" {
		log.Println("missing", keySignature)
		return false
	}

	b, err := base62.DecodeString(sign)
	if err != nil {
		log.Println("invalid", keySignature, err)
		return false
	}

	t := int64(binary.BigEndian.Uint64(b[:]))
	elapsed := time.Now().Unix() - t
	if elapsed < -3 || elapsed > int64(delaySeconds) {
		log.Println("invalid timestamp", t, elapsed)
		return false
	}
	clientID := GetClientID(req.Header)
	traceID := GetTraceID(req.Header)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	return bytes.Equal(b[9:], hash[:])
}
