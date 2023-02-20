package httpkit

import (
	"code.olapie.com/sugar/v2/jsonutil"
	"code.olapie.com/sugar/v2/xerror"
	"log"
	"net/http"
	"strconv"
)

func RequireBasicAuthenticate(realm string, w http.ResponseWriter) {
	a := "Basic realm=" + strconv.Quote(realm)
	w.Header().Set(KeyWWWAuthenticate, a)
	w.WriteHeader(http.StatusUnauthorized)
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	status := xerror.GetCode(err)
	if status < 100 || status > 599 {
		log.Println("invalid status:", status)
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Println(err)
	}
}

func WriteJSON(w http.ResponseWriter, v any) {
	SetContentType(w.Header(), JSON)
	_, err := w.Write(jsonutil.ToBytes(v))
	if err != nil {
		log.Println(err)
	}
}
