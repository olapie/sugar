package mob

import (
	"strings"

	"code.olapie.com/sugar/v2/httpkit"
)

func IsTextFile(f FileInfo) bool {
	return IsMIMEText(f.MIMEType())
}

func IsImageFile(f FileInfo) bool {
	return IsMIMEImage(f.MIMEType())
}

func IsVideoFile(f FileInfo) bool {
	return IsMIMEVideo(f.MIMEType())
}

func IsAudioFile(f FileInfo) bool {
	return IsMIMEAudio(f.MIMEType())
}

func IsMIMEText(t string) bool {
	return httpkit.IsText(t)
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
