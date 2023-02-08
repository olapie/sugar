package xerror

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func ParseHTTPResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !isText(contentType) {
		return New(resp.StatusCode, resp.Status)
	}

	body, ioErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if ioErr != nil {
		log.Printf("Failed reading response body: %v\n", ioErr)
		return nil
	}

	var err Error

	if strings.HasPrefix(contentType, "application/json") {
		var errObj errorJSONObject
		if json.Unmarshal(body, &errObj) == nil {
			err.code = errObj.Code
			err.subCode = errObj.SubCode
			err.message = errObj.Message
		}
	} else {
		err.message = string(body)
	}

	if err.code <= 0 {
		err.code = resp.StatusCode
	}

	if err.message == "" {
		err.message = resp.Status
	}

	return &err
}

var textTypes = []string{
	"text/plain", "text/html", "text/xml", "text/css", "application/xml", "application/xhtml+xml",
}

func isText(mimeType string) bool {
	for _, t := range textTypes {
		if strings.HasPrefix(mimeType, t) {
			return true
		}
	}
	return false
}
