package httpkit

import (
	"context"

	"code.olapie.com/sugar/v2/urlutil"
)

const dateEndpoint = "ola-debug/date"

// GetServerTime only works if server is powered by ola
func GetServerTime(ctx context.Context, serverURL string) (int64, error) {
	type Response struct {
		Timestamp int64 `json:"timestamp"`
	}
	c := NewGet[struct{}, *Response](urlutil.Join(serverURL, dateEndpoint))
	resp, err := c.Call(ctx, struct{}{})
	if err != nil {
		return 0, err
	}
	return resp.Timestamp, nil
}
