package ximage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"code.olapie.com/sugar/v2/xerror"
	"github.com/disintegration/imaging"
)

type ThumbnailOptions struct {
	Width   int
	Height  int
	Quality int
}

func (t *ThumbnailOptions) Validate() error {
	if t.Width <= 0 {
		return fmt.Errorf("negative width %d", t.Width)
	}
	if t.Height <= 0 {
		return fmt.Errorf("negative height %d", t.Height)
	}
	if t.Quality <= 0 || t.Quality > 100 {
		return fmt.Errorf("invalid quality %d, expected (0, 100]", t.Quality)
	}
	return nil
}

func GenerateThumbnail(origin []byte, w, h int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(origin))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, xerror.BadRequest("decode: %v", err)
		}
		return nil, fmt.Errorf("decode: %w", err)
	}
	return getImageThumbnail(img, &ThumbnailOptions{Width: w, Height: h, Quality: 100})
}

func getImageThumbnail(img image.Image, t *ThumbnailOptions) ([]byte, error) {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()

	var tImg image.Image
	if dx < t.Width || dy < t.Height {
		tImg = img
	} else {
		if dx*t.Height > dy*t.Width {
			tImg = imaging.Resize(img, 0, t.Height, imaging.Lanczos)
		} else {
			tImg = imaging.Resize(img, t.Width, 0, imaging.Lanczos)
		}
	}

	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, tImg, &jpeg.Options{Quality: t.Quality})
	if err != nil {
		return nil, fmt.Errorf("encode image to jpeg: %w", err)
	}
	return buf.Bytes(), nil
}
