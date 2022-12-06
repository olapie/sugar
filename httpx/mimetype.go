package httpkit

import (
	"net/http"
)

const (
	Plain    = "text/plain"
	HTML     = "text/html"
	XML2     = "text/xml"
	CSS      = "text/css"
	XML      = "application/xml"
	XHTML    = "application/xhtml+xml"
	Protobuf = "application/x-protobuf"

	FormData = "multipart/form-data"
	GIF      = "image/gif"
	JPEG     = "image/jpeg"
	PNG      = "image/png"
	WEBP     = "image/webp"
	ICON     = "image/x-icon"

	MPEG = "video/mpeg"

	FormURLEncoded = "application/x-www-form-urlencoded"
	OctetStream    = "application/octet-stream"
	JSON           = "application/json"
	PDF            = "application/pdf"
	MSWord         = "application/msword"
	GZIP           = "application/x-gzip"
	WASM           = "application/wasm"
)

const (
	CharsetUTF8 = "charset=utf-8"

	charsetSuffix = "; " + CharsetUTF8

	PlainUTF8 = Plain + charsetSuffix

	// HtmlUTF8 is better than HTMLUTF8, etc.
	HtmlUTF8 = HTML + charsetSuffix
	JsonUTF8 = JSON + charsetSuffix
	XmlUTF8  = XML + charsetSuffix
)

func IsText[T string | http.Header](typeOrHeader T) bool {
	switch v := any(typeOrHeader).(type) {
	case string:
		switch v {
		case Plain, HTML, CSS, XML, XML2, XHTML, JSON, PlainUTF8, HtmlUTF8, JsonUTF8, XmlUTF8:
			return true
		default:
			return false
		}
	case http.Header:
		return IsText(GetContentType(v))
	default:
		return false
	}
}

func IsXML[T string | http.Header](typeOrHeader T) bool {
	switch v := any(typeOrHeader).(type) {
	case string:
		switch v {
		case XML, XML2, XmlUTF8:
			return true
		default:
			return false
		}
	case http.Header:
		return IsXML(GetContentType(v))
	default:
		return false
	}
}

func IsJSON[T string | http.Header](typeOrHeader T) bool {
	switch v := any(typeOrHeader).(type) {
	case string:
		switch v {
		case JSON, JsonUTF8:
			return true
		default:
			return false
		}
	case http.Header:
		return IsXML(GetContentType(v))
	default:
		return false
	}
}
