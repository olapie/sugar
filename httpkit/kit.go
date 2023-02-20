package httpkit

import (
	"net/http"
	"strconv"
)

func RequireBasicAuthenticate(realm string, w http.ResponseWriter) {
	a := "Basic realm=" + strconv.Quote(realm)
	w.Header().Set(KeyWWWAuthenticate, a)
	w.WriteHeader(http.StatusUnauthorized)
}
