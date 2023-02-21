package httphandler

import (
	"code.olapie.com/sugar/v2/types"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type Upload struct {
	Filename string
	Size     int64
	Content  []byte
}

func NewUploadHandler(maxMemory int64, store func(*Upload) error) http.Handler {
	if maxMemory <= 0 {
		maxMemory = 10 * types.MB.Int64()
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(maxMemory)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, fileHeaders := range req.MultipartForm.File {
			for _, fh := range fileHeaders {
				if upload, err := processMultipartFileHeader(fh); err != nil {
					log.Println(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					if err = store(upload); err != nil {
						log.Println(err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		}
	})
}

func processMultipartFileHeader(h *multipart.FileHeader) (*Upload, error) {
	f, err := h.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	upload := &Upload{
		Filename: filepath.Base(h.Filename),
		Size:     h.Size,
		Content:  content,
	}
	return upload, nil
}
