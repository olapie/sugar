package ctxutil

import (
	"context"
	"net/http"
)

type RequestMetadata struct {
	AppID      string
	ClientID   string
	HttpHeader http.Header
	ServiceID  string
	TraceID    string
	TestFlag   bool
}

func WithRequestMetadata(ctx context.Context, v RequestMetadata) context.Context {
	return context.WithValue(ctx, keyRequestMetadata, &v)
}

func GetAppID(ctx context.Context) string {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).AppID
}

func GetHTTPHeader(ctx context.Context) http.Header {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).HttpHeader
}

func GetTraceID(ctx context.Context) string {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).TraceID
}

func GetServiceID(ctx context.Context) string {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).ServiceID
}

func GetClientID(ctx context.Context) string {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).ClientID
}

func IsTest(ctx context.Context) bool {
	return ctx.Value(keyRequestMetadata).(*RequestMetadata).TestFlag
}

func GetRequestMetadata(ctx context.Context) RequestMetadata {
	v, _ := ctx.Value(keyRequestMetadata).(*RequestMetadata)
	return *v
}
