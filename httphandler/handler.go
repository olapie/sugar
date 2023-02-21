package httphandler

import (
	"net/http"

	"code.olapie.com/sugar/v2/httpwriter"
)

type joinHandler struct {
	handlers []http.Handler
}

func (j *joinHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	w, ok := writer.(*httpwriter.Wrapper)
	if !ok {
		w = httpwriter.NewWrapper(writer)
	}

	for _, h := range j.handlers {
		h.ServeHTTP(w, request)
		if w.Status() != 0 {
			return
		}
	}
}

var _ http.Handler = (*joinHandler)(nil)

func Join(handlers ...http.Handler) http.Handler {
	return &joinHandler{
		handlers: handlers,
	}
}

func JoinFuncs(funcs ...http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		w, ok := writer.(*httpwriter.Wrapper)
		if !ok {
			w = httpwriter.NewWrapper(writer)
		}

		for _, f := range funcs {
			f.ServeHTTP(w, request)
			if w.Status() != 0 {
				return
			}
		}
	}
}
