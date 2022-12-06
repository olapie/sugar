package httpx

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"code.olapie.com/sugar/assigning"
	"code.olapie.com/sugar/errorx"
)

type UnmarshalFunc func([]byte, any) error

var contentTypeToUnmarshalFunc sync.Map

func init() {
	RegisterUnmarshalFunc(JSON, json.Unmarshal)
	RegisterUnmarshalFunc(JsonUTF8, json.Unmarshal)
	RegisterUnmarshalFunc(XML, xml.Unmarshal)
	RegisterUnmarshalFunc(XML2, xml.Unmarshal)
	RegisterUnmarshalFunc(XmlUTF8, xml.Unmarshal)
}

func RegisterUnmarshalFunc(contentType string, f UnmarshalFunc) {
	contentTypeToUnmarshalFunc.Store(contentType, f)
}

func GetUnmarshalFunc(contentType string) UnmarshalFunc {
	v, ok := contentTypeToUnmarshalFunc.Load(contentType)
	if ok {
		u, _ := v.(UnmarshalFunc)
		return u
	}
	return nil
}

func GetResponseResult[T any](resp *http.Response) (T, error) {
	var res T
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return res, fmt.Errorf("read resp body: %v", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return res, errorx.Format(resp.StatusCode, string(body))
	}

	if any(res) == nil {
		return res, nil
	}

	if val := reflect.ValueOf(res); val.Kind() == reflect.Struct && val.Type().NumField() == 0 {
		return res, nil
	}

	ct := GetContentType(resp.Header)
	if f := GetUnmarshalFunc(ct); f != nil {
		err = f(body, &res)
		return res, errorx.Wrapf(err, "unmarshal")
	}

	if len(body) == 0 {
		err = errors.New("no data")
	} else if _, ok := any(res).([]byte); ok {
		res = any(body).(T)
	} else {
		if err = assigning.SetBytes(&res, body); err != nil {
			err = fmt.Errorf("cannot handle %s: %w ", ct, err)
		}
	}
	return res, err
}

func RequireBasicAuthenticate(realm string, w http.ResponseWriter) {
	a := "Basic realm=" + strconv.Quote(realm)
	w.Header().Set(KeyWWWAuthenticate, a)
	w.WriteHeader(http.StatusUnauthorized)
}
