package httpkit

import (
	"net/http"
)

type joinHandler struct {
	handlers []http.Handler
}

func (j *joinHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	w, ok := writer.(*Writer)
	if !ok {
		w = NewWriter(writer)
	}

	for _, h := range j.handlers {
		h.ServeHTTP(w, request)
		if w.Status() != 0 {
			return
		}
	}
}

var _ http.Handler = (*joinHandler)(nil)

func JoinHandlers(handlers ...http.Handler) http.Handler {
	return &joinHandler{
		handlers: handlers,
	}
}

func JoinHandlerFuncs(funcs ...http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		w, ok := writer.(*Writer)
		if !ok {
			w = NewWriter(writer)
		}

		for _, f := range funcs {
			f.ServeHTTP(w, request)
			if w.Status() != 0 {
				return
			}
		}
	}
}
