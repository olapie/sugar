package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"code.olapie.com/sugar/v2/xcheck"
	"code.olapie.com/sugar/v2/xerror"
	"code.olapie.com/sugar/v2/xruntime"
	"code.olapie.com/sugar/v2/xurl"
	"github.com/google/uuid"
)

type void = struct{}

type CallResult[R any] struct {
	Value  R
	Header http.Header
	Error  error
}

type Caller[IN any, OUT any] struct {
	Client     *http.Client
	Method     string
	Endpoint   string
	BeforeCall RequestInterceptorFunc
}

func NewCaller[IN any, OUT any](method string, endpoint string) *Caller[IN, OUT] {
	e := &Caller[IN, OUT]{
		Method:   method,
		Endpoint: endpoint,
	}
	return e
}

func (c *Caller[IN, OUT]) WithQuery(query url.Values) *Caller[IN, OUT] {
	cc := *c
	var err error
	cc.Endpoint, err = xurl.AppendQuery(c.Endpoint, query)
	if err != nil {
		fmt.Println("httpkit.Caller.WithQuery", err)
	}
	return &cc
}

func (c *Caller[IN, OUT]) WithQueryArgs(keysAndValues ...any) *Caller[IN, OUT] {
	n := len(keysAndValues)
	if n%2 != 0 {
		panic("keyAndValues is not paired")
	}

	query := url.Values{}
	for i := 0; i < n; i += 2 {
		k := keysAndValues[i]
		v := keysAndValues[i+1]
		ks, ok := k.(string)
		if !ok {
			if stringer, ok := k.(fmt.Stringer); ok {
				ks = stringer.String()
			}
		}

		if ks == "" {
			panic(fmt.Sprintf("keysAndValues[%d] is not a string key", i))
		}

		vs, ok := v.(string)
		if !ok {
			if stringer, ok := v.(fmt.Stringer); ok {
				vs = stringer.String()
			} else if xcheck.IsNumber(v) {
				vs = fmt.Sprint(v)
			}
		}
		if vs == "" {
			panic(fmt.Sprintf("keysAndValues[%d] is not a string or number value", i+1))
		}
		query.Set(ks, vs)
	}
	return c.WithQuery(query)
}

func (c *Caller[IN, OUT]) Call(ctx context.Context, input IN) (OUT, error) {
	var out OUT
	resp, err := c.call(ctx, input)
	if err != nil {
		return out, err
	}
	return GetResponseResult[OUT](resp)
}

func (c *Caller[IN, OUT]) GetResult(ctx context.Context, input IN) *CallResult[OUT] {
	res := new(CallResult[OUT])
	resp, err := c.call(ctx, input)
	if err != nil {
		res.Error = err
		return res
	}
	res.Header = resp.Header
	out, err := GetResponseResult[OUT](resp)
	if err != nil {
		res.Error = err
		return res
	}
	res.Value = out
	return res
}

func (c *Caller[IN, OUT]) CallAndRewrite(ctx context.Context, input IN, writer io.Writer) error {
	resp, err := c.call(ctx, input)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, resp.Body)
	return err
}

func (c *Caller[IN, OUT]) call(ctx context.Context, input IN) (*http.Response, error) {
	var contentType string
	endpoint, err := url.PathUnescape(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("unescape path: %w", err)
	}
	body, err := c.parseInput(&contentType, &endpoint, input)
	if err != nil {
		return nil, fmt.Errorf("parse input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, c.Method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request %s %s: %w", c.Method, endpoint, err)
	}
	req.Header.Set(KeyContentType, contentType)
	req.Header.Set(KeyTraceID, uuid.NewString())

	client := http.DefaultClient
	if c.Client != nil {
		client = c.Client
	}

	if c.BeforeCall != nil {
		if err = c.BeforeCall(req); err != nil {
			return nil, fmt.Errorf("before call: %w", err)
		}
	}

	fmt.Println(req.Method, req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			err = xerror.RequestTimeout(err.Error())
		} else {
			if tr, ok := err.(interface{ Timeout() bool }); ok && tr.Timeout() {
				err = xerror.RequestTimeout(err.Error())
			} else {
				err = xerror.New(600, err.Error())
			}
		}
		return nil, fmt.Errorf("send request: %w", err)
	}
	return resp, nil
}

func (c *Caller[IN, OUT]) parseInput(contentType *string, endpoint *string, input any) (io.Reader, error) {
	if input == nil {
		return nil, nil
	}

	if b, ok := input.([]byte); ok {
		return bytes.NewReader(b), nil
	}

	body, ok := input.(io.Reader)
	if ok {
		if *contentType == "" {
			*contentType = OctetStream
		}
		return body, nil
	}

	if v, ok := input.(url.Values); ok {
		newEndpoint, err := xurl.AppendQuery(*endpoint, v)
		if err != nil {
			return nil, err
		}
		*endpoint = newEndpoint
		return nil, nil
	}

	newEndpoint, remain := xurl.SetPathParams(*endpoint, input)
	*endpoint = newEndpoint

	if remain == nil {
		return nil, nil
	}

	kindOfRemain := xruntime.IndirectKind(remain)
	switch kindOfRemain {
	case reflect.Struct, reflect.Map, reflect.Slice:
		*contentType = JsonUTF8
		data, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}
		return bytes.NewBuffer(data), nil
	default:
		if xcheck.IsNumber(remain) || xcheck.IsString(remain) {
			*contentType = PlainUTF8
			return bytes.NewReader([]byte(fmt.Sprint(remain))), nil
		}
		return nil, fmt.Errorf("unsupported value type: %T", input)
	}
}

func NewGet[IN any, OUT any](endpoint string) *Caller[IN, OUT] {
	return NewCaller[IN, OUT](http.MethodGet, endpoint)
}

func NewPost[IN any, OUT any](endpoint string) *Caller[IN, OUT] {
	return NewCaller[IN, OUT](http.MethodPost, endpoint)
}

func NewPut[IN any, OUT any](endpoint string) *Caller[IN, OUT] {
	return NewCaller[IN, OUT](http.MethodPut, endpoint)
}

func NewPatch[IN any, OUT any](endpoint string) *Caller[IN, OUT] {
	return NewCaller[IN, OUT](http.MethodPatch, endpoint)
}

func NewDelete[IN any](endpoint string) *Caller[IN, void] {
	return NewCaller[IN, void](http.MethodDelete, endpoint)
}

func NewHead(endpoint string) *Caller[void, void] {
	return NewCaller[void, void](http.MethodHead, endpoint)
}

func NewOptions(endpoint string) *Caller[void, void] {
	return NewCaller[void, void](http.MethodOptions, endpoint)
}

func NewTrace[T any](endpoint string) *Caller[T, T] {
	return NewCaller[T, T](http.MethodTrace, endpoint)
}

func NewConnect(endpoint string) *Caller[void, void] {
	return NewCaller[void, void](http.MethodConnect, endpoint)
}
