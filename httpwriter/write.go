package httpwriter

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"code.olapie.com/sugar/v2/httpheader"
	"code.olapie.com/sugar/v2/jsonutil"
	"code.olapie.com/sugar/v2/xerror"
)

func RequireBasicAuthenticate(w http.ResponseWriter, realm string) {
	a := "Basic realm=" + strconv.Quote(realm)
	w.Header().Set(httpheader.KeyWWWAuthenticate, a)
	w.WriteHeader(http.StatusUnauthorized)
}

func Error(w http.ResponseWriter, err error) {
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

func JSON(w http.ResponseWriter, v any) {
	httpheader.SetContentType(w.Header(), httpheader.JSON)
	_, err := w.Write(jsonutil.ToBytes(v))
	if err != nil {
		log.Println(err)
	}
}

func StreamFile(w http.ResponseWriter, name string, f io.ReadCloser) {
	defer f.Close()
	httpheader.SetContentType(w.Header(), httpheader.OctetStream)
	if name != "" {
		w.Header().Set(httpheader.KeyContentDisposition, httpheader.ToAttachment(name))
	}
	_, err := io.Copy(w, f)
	if err != nil {
		if err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
