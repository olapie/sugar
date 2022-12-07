package mobilex

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type StorageHandler interface {
	CacheImage(link string, data []byte)
}

type Storage struct {
	dir      string
	imageURL string
	handler  StorageHandler
}

func NewStorage(dir, imageURL string, handler StorageHandler) *Storage {
	MustMkdir(dir)
	return &Storage{
		dir:      dir,
		handler:  handler,
		imageURL: imageURL,
	}
}

type ProgressHandler interface {
	OnProgress(size int)
}

type progressBuffer struct {
	buf bytes.Buffer
	h   ProgressHandler
}

func (b *progressBuffer) Write(p []byte) (n int, err error) {
	n, err = b.buf.Write(p)
	if n > 0 && b.h != nil {
		b.h.OnProgress(n)
	}
	return
}

func (b *progressBuffer) Read(p []byte) (n int, err error) {
	n, err = b.buf.Read(p)
	if n > 0 && b.h != nil {
		b.h.OnProgress(n)
	}
	return
}

func (s *Storage) UploadImage(name string, data []byte, handler ProgressHandler) *StringE {
	res := new(StringE)
	buf := new(progressBuffer)
	w := multipart.NewWriter(buf)
	part, err := w.CreateFormFile("images", name)
	if err != nil {
		res.Error = ToError(err)
		return res
	}
	if _, err = part.Write(data); err != nil {
		res.Error = ToError(err)
		return res
	}
	w.Close()

	buf.h = handler
	resp, err := http.Post(s.imageURL, w.FormDataContentType(), buf)
	if err != nil {
		res.Error = ToError(err)
		return res
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		res.Error = ToError(err)
		return res
	}
	resp.Body.Close()
	var url string
	if err = json.Unmarshal(body, &url); err != nil {
		res.Error = ToError(err)
		return res
	}
	res.Value = url
	if s.handler != nil {
		s.handler.CacheImage(res.Value, data)
		s.Delete(name)
	}
	return res
}

func (s *Storage) Save(data []byte) *StringE {
	res := new(StringE)
	res.Value = uuid.NewString()
	res.Error = ToError(os.WriteFile(s.GetFilePath(res.Value), data, 0644))
	return res
}

func (s *Storage) GetFilePath(name string) string {
	return filepath.Join(s.dir, name)
}

func (s *Storage) Delete(name string) *Error {
	return ToError(os.Remove(s.GetFilePath(name)))
}
