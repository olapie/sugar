package ctxutil

import (
	"context"
	"net/http"
)

type requestContextInfo struct {
	AppID      string
	ClientID   string
	HttpHeader http.Header
	ServiceID  string
	TraceID    string
	TestFlag   bool
}

type RequestContextBuilder interface {
	Build() context.Context
	WithAppID(v string) RequestContextBuilder
	WithClientID(v string) RequestContextBuilder
	WithHTTPHeader(v http.Header) RequestContextBuilder
	WithServiceID(v string) RequestContextBuilder
	WithTraceID(v string) RequestContextBuilder
	WithTestFlag(v bool) RequestContextBuilder
}

type requestContextBuilderImpl struct {
	ctx  context.Context
	info requestContextInfo
}

func Request(ctx context.Context) RequestContextBuilder {
	return &requestContextBuilderImpl{ctx: ctx}
}

func (b *requestContextBuilderImpl) Build() context.Context {
	return context.WithValue(b.ctx, keyRequestInfo, &b.info)
}

func (b *requestContextBuilderImpl) WithAppID(v string) RequestContextBuilder {
	b.info.AppID = v
	return b
}

func (b *requestContextBuilderImpl) WithClientID(v string) RequestContextBuilder {
	b.info.ClientID = v
	return b
}

func (b *requestContextBuilderImpl) WithHTTPHeader(v http.Header) RequestContextBuilder {
	b.info.HttpHeader = v
	return b
}

func (b *requestContextBuilderImpl) WithServiceID(v string) RequestContextBuilder {
	b.info.ServiceID = v
	return b
}

func (b *requestContextBuilderImpl) WithTraceID(v string) RequestContextBuilder {
	b.info.TraceID = v
	return b
}

func (b *requestContextBuilderImpl) WithTestFlag(v bool) RequestContextBuilder {
	b.info.TestFlag = v
	return b
}

func GetAppID(ctx context.Context) string {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return ""
	}
	return info.AppID
}

func GetHTTPHeader(ctx context.Context) http.Header {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return nil
	}
	return info.HttpHeader
}

func GetTraceID(ctx context.Context) string {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return ""
	}
	return info.TraceID
}

func GetServiceID(ctx context.Context) string {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return ""
	}
	return info.ServiceID
}

func GetClientID(ctx context.Context) string {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return ""
	}
	return info.ClientID
}

func IsTest(ctx context.Context) bool {
	info, _ := ctx.Value(keyRequestInfo).(*requestContextInfo)
	if info == nil {
		return false
	}
	return info.TestFlag
}
