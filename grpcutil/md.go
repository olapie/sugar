package grpcutil

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	keyClientID      = "x-client-id"
	keyAppID         = "x-app-id"
	keyTraceID       = "x-trace-id"
	keyAPIKey        = "x-api-key"
	keyAuthorization = "authorization"
)

func MatchMetadata(key string) (string, bool) {
	key = strings.ToLower(key)
	switch key {
	case keyClientID, keyAppID, keyTraceID, keyAPIKey:
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
