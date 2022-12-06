package imagex

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"code.olapie.com/sugar/errorx"
)

type Resizer interface {
	// Resize resizes the image to the specified width and height using the specified resampling
	// filter and returns the transformed image. If one of width or height is 0, the image aspect
	// ratio is preserved.
	Resize(img image.Image, width int, height int) image.Image
}

func Resize(origin []byte, width, height, quality int, resizer Resizer) (resized []byte, err error) {
	img, typ, err := image.Decode(bytes.NewReader(origin))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, errorx.BadRequest("decode: %v", err)
		}
		return nil, fmt.Errorf("decode: %w", err)
	}

	var resizedImg image.Image
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	if dx < width || dy < height {
		resizedImg = img
	} else {
		if dx*height > dy*width {
			resizedImg = resizer.Resize(img, 0, height)
		} else {
			resizedImg = resizer.Resize(img, width, 0)
		}
	}

	buf := bytes.NewBuffer(nil)
	switch typ {
	case "jpeg", "jpg":
		err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("cannot encode jpeg: %w", err)
		}
	case "png":
		err = png.Encode(buf, resizedImg)
		if err != nil {
			return nil, fmt.Errorf("cannot encode png: %w", err)
		}
	case "gif":
		err = gif.Encode(buf, resizedImg, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot encode gif: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupport image format: %s", typ)
	}
	return buf.Bytes(), nil
}
