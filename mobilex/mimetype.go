package mobilex

import (
	"strings"

	"code.olapie.com/sugar/httpx"
)

func IsMIMEText(t string) bool {
	return httpx.IsText(t)
}

func IsMIMEImage(t string) bool {
	return strings.HasPrefix(t, "image") || strings.HasPrefix(t, "photo")
}

func IsMIMEVideo(t string) bool {
	return strings.HasPrefix(t, "video")
}

func IsMIMEAudio(t string) bool {
	return strings.HasPrefix(t, "audio")
}
