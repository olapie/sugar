package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.olapie.com/sugar/errorx"
)

func DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, errorx.ParseHTTPResponse(resp)
}

func ParseRequest(req *http.Request, memInBytes int64) (map[string]any, []byte, error) {
	typ := GetContentType(req.Header)
	params := map[string]any{}
	switch typ {
	case HTML, Plain:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return params, nil, fmt.Errorf("read html or plain body: %w", err)
		}
		return params, body, nil
	case JSON:
		body, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return params, nil, fmt.Errorf("read json body: %w", err)
		}
		if len(body) == 0 {
			return params, nil, nil
		}
		decoder := json.NewDecoder(bytes.NewBuffer(body))
		decoder.UseNumber()
		err = decoder.Decode(&params)
		if err != nil {
			var obj any
			err = json.Unmarshal(body, &obj)
			if err != nil {
				return params, body, fmt.Errorf("unmarshal json %s: %w", string(body), err)
			}
		}
		return params, body, nil
	case FormURLEncoded:
		// TODO: will crash
		//body, err := req.GetBody()
		//if err != nil {
		//	return params, nil, fmt.Errorf("get body: %w", err)
		//}
		//bodyData, err := ioutil.Read(body)
		//body.Close()
		//if err != nil {
		//	return params, nil, fmt.Errorf("read form body: %w", err)
		//}
		if err := req.ParseForm(); err != nil {
			return params, nil, fmt.Errorf("parse form: %w", err)
		}
		return valuesToMap(req.Form), nil, nil
	case FormData:
		err := req.ParseMultipartForm(memInBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("parse multipart form: %w", err)
		}

		if req.MultipartForm != nil && req.MultipartForm.File != nil {
			return valuesToMap(req.MultipartForm.Value), nil, nil
		}
		return params, nil, nil
	default:
		body, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return params, nil, fmt.Errorf("read json body: %w", err)
		}
		return params, body, nil
	}
}

type RequestInterceptorFunc func(req *http.Request) error
