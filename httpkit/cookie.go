package httpkit

import "net/http"

const (
	CookieNameClientID = "client_id"
	CookieNameUserID   = "user_id"
)

func GetCookieValue(req *http.Request, name string) string {
	cookie, err := req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
