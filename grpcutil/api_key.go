package grpcutil

import (
	"bytes"
	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/hashutil"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func SetAPIKey(md metadata.MD) {
	t := time.Now().Unix()
	var b [41]byte
	b[0] = 1
	binary.BigEndian.PutUint64(b[1:], uint64(t))
	clientID := GetClientID(md)
	traceID := GetTraceID(md)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	copy(b[9:], hash[:])
	sign := base62.EncodeToString(b[:])
	md.Set(keyAPIKey, sign)
}

func VerifyAPIKey(md metadata.MD, delaySeconds int) bool {
	sign := GetMetadata(md, keyAPIKey)
	if sign == "" {
		log.Println("missing", keyAPIKey)
		return false
	}

	b, err := base62.DecodeString(sign)
	if err != nil {
		log.Println("invalid", keyAPIKey, err)
		return false
	}
	t := int64(binary.BigEndian.Uint64(b[1:]))
	elapsed := time.Now().Unix() - t
	if elapsed < -3 || elapsed > int64(delaySeconds) {
		log.Println("invalid timestamp", t, elapsed)
		return false
	}
	clientID := GetClientID(md)
	traceID := GetTraceID(md)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	return bytes.Equal(b[9:], hash[:])
}
