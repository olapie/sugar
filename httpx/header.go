package httpx

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"code.olapie.com/sugar/must"
)

const (
	KeyAuthorization       = "Authorization"
	KeyAcceptEncoding      = "Accept-Encoding"
	KeyACLAllowCredentials = "Access-Control-Allow-Credentials"
	KeyACLAllowHeaders     = "Access-Control-Allow-Headers"
	KeyACLAllowMethods     = "Access-Control-Allow-Methods"
	KeyACLAllowOrigin      = "Access-Control-Allow-Origin"
	KeyACLExposeHeaders    = "Access-Control-Expose-Headers"
	KeyContentType         = "Content-Type"
	KeyContentDisposition  = "Content-Disposition"
	KeyContentEncoding     = "Content-Encoding"
	KeyCookies             = "Cookies"
	KeyLocation            = "Location"
	KeyReferrer            = "Referer"
	KeyReferrerPolicy      = "Referrer-Policy"
	KeyUserAgent           = "User-Agent"
	KeyWWWAuthenticate     = "WWW-Authenticate"
	KeyAcceptLanguage      = "Accept-Language"

	KeyClientID      = "X-Client-Id"
	KeyApplicationID = "X-Application-Id"
	KeyServiceID     = "X-Service-Id"
	KeyTraceID       = "X-Trace-Id"
	KeyUserID        = "X-User-Id"
)

const (
	Bearer = "Bearer"
	Basic  = "Basic"
)

type HeaderTypeSet interface {
	http.Header | Header | *Header | map[string]string | map[string][]string
}

func GetHeader[H HeaderTypeSet](h H, key string) string {
	switch m := any(h).(type) {
	case map[string]string:
		v := m[key]
		if v == "" {
			v = m[strings.ToLower(key)]
		}
		return v
	case map[string][]string:
		hh := http.Header(m)
		v := hh.Get(key)
		if v == "" {
			v = hh.Get(strings.ToLower(key))
		}
		return v
	case http.Header:
		v := m.Get(key)
		if v == "" {
			v = m.Get(strings.ToLower(key))
		}
		return v
	case Header:
		v := m.Get(key)
		if v == "" {
			v = m.Get(strings.ToLower(key))
		}
		return v
	case *Header:
		v := m.Get(key)
		if v == "" {
			v = m.Get(strings.ToLower(key))
		}
		return v
	}
	return ""
}

func SetHeader[H HeaderTypeSet](h H, key, value string) {
	switch m := any(h).(type) {
	case map[string]string:
		m[key] = value
	case map[string][]string:
		hh := http.Header(m)
		hh.Set(key, value)
	case http.Header:
		m.Set(key, value)
	case Header:
		m.Set(key, value)
	case *Header:
		m.Set(key, value)
	default:
		panic(fmt.Sprintf("unsupported type %T", h))
	}
}

func SetHeaderNX[H HeaderTypeSet](h H, key, value string) {
	if GetHeader(h, key) != "" {
		return
	}
	SetHeader(h, key, value)
}

func GetAcceptEncodings[H HeaderTypeSet](h H) []string {
	a := strings.Split(GetHeader(h, KeyAcceptEncoding), ",")
	for i, s := range a {
		a[i] = strings.TrimSpace(s)
	}

	// Remove empty strings
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func GetContentType[H HeaderTypeSet](h H) string {
	t := GetHeader(h, KeyContentType)
	for i, ch := range t {
		if ch == ' ' || ch == ';' {
			t = t[:i]
			break
		}
	}
	return t
}

func SetContentType(h http.Header, contentType string) {
	h.Set(KeyContentType, contentType)
}

func SetContentTypeIfNX(h http.Header, contentType string) {
	SetHeaderNX(h, KeyContentType, contentType)
}

func GetAuthorization[H HeaderTypeSet](h H) string {
	return GetHeader(h, KeyAuthorization)
}

func SetAuthorization[H HeaderTypeSet](h H, contentType string) {
	SetHeader(h, KeyContentType, contentType)
}

func SetAuthorizationNX[H HeaderTypeSet](h H, contentType string) {
	SetHeaderNX(h, KeyContentType, contentType)
}

func GetBasicAccount[H HeaderTypeSet](h H) (user string, password string) {
	s := GetAuthorization(h)
	l := strings.Split(s, " ")
	if len(l) != 2 {
		return
	}

	if l[0] != Basic {
		return
	}

	b, err := base64.StdEncoding.DecodeString(l[1])
	if err != nil {
		return
	}

	userAndPass := strings.Split(string(b), ":")
	if len(userAndPass) != 2 {
		return
	}
	return userAndPass[0], userAndPass[1]
}

// GetBearer returns bearer token in header
func GetBearer[H HeaderTypeSet](h H) string {
	s := GetAuthorization(h)
	l := strings.Split(s, " ")
	if len(l) != 2 {
		return ""
	}
	if l[0] == Bearer {
		return l[1]
	}
	return ""
}

func SetBearer[H HeaderTypeSet](h H, bearer string) {
	authorization := Bearer + " " + bearer
	SetHeader(h, KeyAuthorization, authorization)
}

func GetContentEncoding(h http.Header, encoding string) string {
	return GetHeader(h, KeyContentEncoding)
}

func SetContentEncoding[H HeaderTypeSet](h H, encoding string) {
	SetHeader(h, KeyContentEncoding, encoding)
}

func GetUserID[H HeaderTypeSet](h H) string {
	return GetHeader(h, KeyUserID)
}

func SetUserID[H HeaderTypeSet](h H, id string) {
	SetHeader(h, KeyUserID, id)
}

func GetTraceID[H HeaderTypeSet](h H) string {
	return GetHeader(h, KeyTraceID)
}

func SetTraceID[H HeaderTypeSet](h H, id string) {
	SetHeader(h, KeyTraceID, id)
}

func GetClientID[H HeaderTypeSet](h H) string {
	return GetHeader(h, KeyClientID)
}

func SetClientID[H HeaderTypeSet](h H, id string) {
	SetHeader(h, KeyClientID, id)
}

func GetApplicationID[H HeaderTypeSet](h H) string {
	return GetHeader(h, KeyApplicationID)
}

func SetApplicationID[H HeaderTypeSet](h H, id string) {
	SetHeader(h, KeyApplicationID, id)
}

func IsWebsocket(h http.Header) bool {
	conn := strings.ToLower(h.Get("Connection"))
	if conn != "upgrade" {
		return false
	}
	return strings.EqualFold(h.Get("Upgrade"), "websocket")
}

// ToAttachment returns value for Content-Disposition
// e.g. Content-Disposition: attachment; filename=test.txt
func ToAttachment(filename string) string {
	return fmt.Sprintf(`attachment; filename="%s"`, filename)
}

func CreateUserAuthorizations(userToPassword map[string]string) map[string]string {
	userToAuthorization := make(map[string]string)
	for user, password := range userToPassword {
		if user == "" || password == "" {
			panic("empty user or password")
		}
		account := user + ":" + password
		userToAuthorization[user] = "Basic " + base64.StdEncoding.EncodeToString([]byte(account))
	}
	return userToAuthorization
}

type Header struct {
	http.Header
}

func NewHeader() *Header {
	return &Header{
		Header: make(http.Header),
	}
}

func (h *Header) WriteTo(rw http.ResponseWriter) {
	for k, v := range h.Header {
		rw.Header()[k] = v
	}
}

func (h *Header) Clone() *Header {
	c := &Header{
		Header: make(http.Header),
	}
	for k, v := range h.Header {
		c.Header[k] = v
	}
	return c
}

func (h *Header) AllowOrigins(origins ...string) {
	h.Header[KeyACLAllowOrigin] = origins
}

func (h *Header) SetAuthorization(credential string) {
	h.Set(KeyAuthorization, credential)
}

func (h *Header) SetBasicAuthorization(account, password string) {
	credential := []byte(account + ":" + password)
	h.Set(KeyAuthorization, "Basic "+base64.StdEncoding.EncodeToString(credential))
}

func (h *Header) Authorization() string {
	return h.Get(KeyAuthorization)
}

func (h *Header) SetClientID(id string) {
	h.Set(KeyClientID, id)
}

func (h *Header) ClientID() string {
	return h.Get(KeyClientID)
}

func (h *Header) SetApplicationID(id string) {
	h.Set(KeyApplicationID, id)
}

func (h *Header) ApplicationID() string {
	return h.Get(KeyApplicationID)
}

func (h *Header) AllowMethods(methods ...string) {
	// Combine multiple values separated by comma. Multiple lines style is also fine.
	h.Header.Set(KeyACLAllowMethods, strings.Join(methods, ","))
}

func (h *Header) AllowCredentials(b bool) {
	h.Header.Set(KeyACLAllowCredentials, must.ToString(b))
}

func (h *Header) AllowHeaders(headers ...string) {
	h.Header.Set(KeyACLAllowHeaders, strings.Join(headers, ","))
}

func (h *Header) ExposeHeaders(headers ...string) {
	h.Header.Set(KeyACLExposeHeaders, strings.Join(headers, ","))
}

func (h *Header) SetContentEncoding(encoding string) {
	h.Header.Set(KeyContentEncoding, encoding)
}
