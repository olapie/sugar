package grpcutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/hashutil"
	"google.golang.org/grpc/metadata"
)

const (
	keyClientID      = "x-client-id"
	keyAppID         = "x-app-id"
	keyTraceID       = "x-trace-id"
	keySignature     = "x-sign"
	keyAuthorization = "authorization"
)

func MatchMetadata(key string) (string, bool) {
	key = strings.ToLower(key)
	switch key {
	case keyClientID, keyAppID, keyTraceID, keySignature:
		return key, true
	default:
		return "", false
	}
}

func GetTraceID(md metadata.MD) string {
	return GetMetadata(md, keyTraceID)
}

func SetTraceID(md metadata.MD, id string) {
	md.Set(keyTraceID, id)
}

func GetClientID(md metadata.MD) string {
	return GetMetadata(md, keyClientID)
}

func SetClientID(md metadata.MD, id string) {
	md.Set(keyClientID, id)
}

func GetAppID(md metadata.MD) string {
	return GetMetadata(md, keyAppID)
}

func SetAppID(md metadata.MD, id string) {
	md.Set(keyAppID, id)
}

func GetAuthorization(md metadata.MD) string {
	return GetMetadata(md, keyAuthorization)
}

func SetAuthorization(md metadata.MD, a string) {
	md.Set(keyAuthorization, a)
}

func GetMetadata(m metadata.MD, key string) string {
	v, ok := m[key]
	if !ok {
		v = m[strings.ToLower(key)]
	}
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func Sign(md metadata.MD) {
	t := time.Now().Unix()
	var b [41]byte
	b[0] = 1
	binary.BigEndian.PutUint64(b[1:], uint64(t))
	clientID := GetClientID(md)
	traceID := GetTraceID(md)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	copy(b[9:], hash[:])
	sign := base62.EncodeToString(b[:])
	md.Set(keySignature, sign)
}

func Verify(md metadata.MD, delaySeconds int) bool {
	sign := GetMetadata(md, keySignature)
	if sign == "" {
		fmt.Println("missing", keySignature)
		return false
	}

	b, err := base62.DecodeString(sign)
	if err != nil {
		fmt.Println("invalid", keySignature, err)
		return false
	}
	t := int64(binary.BigEndian.Uint64(b[1:]))
	elapsed := time.Now().Unix() - t
	if elapsed < -3 || elapsed > int64(delaySeconds) {
		fmt.Println("invalid timestamp", t, elapsed)
		return false
	}
	clientID := GetClientID(md)
	traceID := GetTraceID(md)
	hash := hashutil.Hash32(fmt.Sprint(t) + traceID + clientID)
	return bytes.Equal(b[9:], hash[:])
}
