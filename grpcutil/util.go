package grpcutil

import (
	"code.olapie.com/log"
	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/ctxutil"
	"code.olapie.com/sugar/v2/xerror"
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"reflect"
	"time"
)

// Refer to https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md

var statusErrorType = reflect.TypeOf(status.Error(codes.Unknown, ""))
var statusToCodeMap = map[int]codes.Code{
	http.StatusUnauthorized: codes.Unauthenticated,
	http.StatusForbidden:    codes.PermissionDenied,
	http.StatusBadRequest:   codes.InvalidArgument,
	http.StatusNotFound:     codes.NotFound,
	http.StatusConflict:     codes.AlreadyExists,

	http.StatusNotImplemented: codes.Unimplemented,

	http.StatusInternalServerError: codes.Internal,
	http.StatusBadGateway:          codes.Unavailable,
	http.StatusServiceUnavailable:  codes.Unavailable,
}

func ServerStart(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed reading request metadata")
	}

	if !VerifyAPIKey(md, 10) {
		log.FromContext(ctx).Sugar().Infoln(md)
		return nil, status.Error(codes.InvalidArgument, "failed verifying")
	}

	traceID := GetTraceID(md)
	ctx = ctxutil.Request(ctx).WithAppID(GetAppID(md)).
		WithClientID(GetClientID(md)).
		WithTraceID(traceID).
		Build()
	logger := log.FromContext(ctx).With(log.String("trace_id", traceID))
	fields := make([]log.Field, 0, len(md)+1)
	fields = append(fields, log.String("full_method", info.FullMethod))
	for k, v := range md {
		if len(v) == 0 {
			continue
		}
		fields = append(fields, log.String(k, v[0]))
	}
	logger.Info("start", fields...)
	ctx = log.BuildContext(ctx, logger)
	return ctx, nil
}

func ServerFinish(resp any, err error, logger *log.Logger, startAt time.Time) (any, error) {
	if err == nil {
		logger.Info("finished", log.Duration("cost", time.Now().Sub(startAt)))
		return resp, nil
	}

	logger.Error("failed", log.Err(err))

	if reflect.TypeOf(err) == statusErrorType {
		return nil, err
	}

	code, ok := statusToCodeMap[xerror.GetCode(err)]
	if !ok {
		code = codes.Unknown
	}
	return nil, status.Error(code, err.Error())
}

func SignClientContext(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}

	if traceID := GetTraceID(md); traceID == "" {
		traceID = ctxutil.GetTraceID(ctx)
		if traceID == "" {
			traceID = base62.NewUUIDString()
		}
		SetTraceID(md, traceID)
	}

	if clientID := GetClientID(md); clientID == "" {
		clientID = ctxutil.GetClientID(ctx)
		if clientID != "" {
			SetClientID(md, clientID)
		}
	}

	if appID := GetAppID(md); appID == "" {
		appID = ctxutil.GetAppID(ctx)
		if appID != "" {
			SetAppID(md, appID)
		}
	}

	SetAPIKey(md)
	return metadata.NewOutgoingContext(ctx, md)
}

func WithClientCert(cert []byte) grpc.DialOption {
	config := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
			},
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(config))
}
