package xhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"code.olapie.com/sugar/xerror"
)

func DoWithResponse(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}

	if resp.StatusCode < 400 {
		return resp, nil
	}

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	resp.Body.Close()
	return nil, &xerror.Error{
		Code:    resp.StatusCode,
		Message: string(message),
	}
}

func Do(ctx context.Context, method, url string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}

	if resp.StatusCode < 400 {
		return nil
	}

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	resp.Body.Close()
	return &xerror.Error{
		Code:    resp.StatusCode,
		Message: string(message),
	}
}

func Post(ctx context.Context, url string, body io.Reader) error {
	return Do(ctx, http.MethodPost, url, body)
}

func Put(ctx context.Context, url string, body io.Reader) error {
	return Do(ctx, http.MethodPut, url, body)
}

func Patch(ctx context.Context, url string, body io.Reader) error {
	return Do(ctx, http.MethodPut, url, body)
}

func Get(ctx context.Context, url string) error {
	return Do(ctx, http.MethodDelete, url, nil)
}

func Delete(ctx context.Context, url string) error {
	return Do(ctx, http.MethodDelete, url, nil)
}
