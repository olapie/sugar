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

type requestContextBuilder struct {
	ctx  context.Context
	info requestContextInfo
}

func Request(ctx context.Context) *requestContextBuilder {
	return &requestContextBuilder{ctx: ctx}
}

func (b *requestContextBuilder) Build() context.Context {
	return context.WithValue(b.ctx, keyRequestInfo, &b.info)
}

func (b *requestContextBuilder) WithAppID(v string) *requestContextBuilder {
	b.info.AppID = v
	return b
}

func (b *requestContextBuilder) WithClientID(v string) *requestContextBuilder {
	b.info.ClientID = v
	return b
}

func (b *requestContextBuilder) WithHTTPHeader(v http.Header) *requestContextBuilder {
	b.info.HttpHeader = v
	return b
}

func (b *requestContextBuilder) WithServiceID(v string) *requestContextBuilder {
	b.info.ServiceID = v
	return b
}

func (b *requestContextBuilder) WithTraceID(v string) *requestContextBuilder {
	b.info.TraceID = v
	return b
}

func (b *requestContextBuilder) WithTestFlag(v bool) *requestContextBuilder {
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
