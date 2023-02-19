package grpcutil

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	keyClientID = "x-client-id"
	keyAppID    = "x-app-id"
	keyTraceID  = "x-trace-id"
	keyUserID   = "x-user-id"
)

func GetUserID(md metadata.MD) string {
	return GetMetadata(md, keyUserID)
}

func SetUserID(md metadata.MD, id string) {
	md.Set(keyUserID, id)
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
