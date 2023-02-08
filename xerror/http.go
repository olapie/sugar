package xerror

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func ParseHTTPResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !isText(contentType) {
		return New(resp.StatusCode, resp.Status)
	}

	body, ioErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if ioErr != nil {
		log.Printf("Failed reading response body: %v\n", ioErr)
		return nil
	}

	var err Error

	if strings.HasPrefix(contentType, "application/json") {
		var errObj errorJSONObject
		if json.Unmarshal(body, &errObj) == nil {
			err.code = errObj.Code
			err.subCode = errObj.SubCode
			err.message = errObj.Message
		}
	} else {
		err.message = string(body)
	}

	if err.code <= 0 {
		err.code = resp.StatusCode
	}

	if err.message == "" {
		err.message = resp.Status
	}

	return &err
}

var textTypes = []string{
	"text/plain", "text/html", "text/xml", "text/css", "application/xml", "application/xhtml+xml",
}

func isText(mimeType string) bool {
	for _, t := range textTypes {
		if strings.HasPrefix(mimeType, t) {
			return true
		}
	}
	return false
}

func BadRequest(format string, a ...any) *Error {
	return New(http.StatusBadRequest, format, a...)
}

func Unauthorized(format string, a ...any) *Error {
	return New(http.StatusUnauthorized, format, a...)
}

func PaymentRequired(format string, a ...any) *Error {
	return New(http.StatusPaymentRequired, format, a...)
}

func Forbidden(format string, a ...any) *Error {
	return New(http.StatusForbidden, format, a...)
}

func NotFound(format string, a ...any) *Error {
	return New(http.StatusNotFound, format, a...)
}

func MethodNotAllowed(format string, a ...any) *Error {
	return New(http.StatusMethodNotAllowed, format, a...)
}

func NotAcceptable(format string, a ...any) *Error {
	return New(http.StatusNotAcceptable, format, a...)
}

func ProxyAuthRequired(format string, a ...any) *Error {
	return New(http.StatusProxyAuthRequired, format, a...)
}

func RequestTimeout(format string, a ...any) *Error {
	return New(http.StatusRequestTimeout, format, a...)
}

func Conflict(format string, a ...any) *Error {
	return New(http.StatusConflict, format, a...)
}

func LengthRequired(format string, a ...any) *Error {
	return New(http.StatusLengthRequired, format, a...)
}

func PreconditionFailed(format string, a ...any) *Error {
	return New(http.StatusPreconditionFailed, format, a...)
}

func RequestEntityTooLarge(format string, a ...any) *Error {
	return New(http.StatusRequestEntityTooLarge, format, a...)
}

func RequestURITooLong(format string, a ...any) *Error {
	return New(http.StatusRequestURITooLong, format, a...)
}

func ExpectationFailed(format string, a ...any) *Error {
	return New(http.StatusExpectationFailed, format, a...)
}

func Teapot(format string, a ...any) *Error {
	return New(http.StatusTeapot, format, a...)
}

func MisdirectedRequest(format string, a ...any) *Error {
	return New(http.StatusMisdirectedRequest, format, a...)
}

func UnprocessableEntity(format string, a ...any) *Error {
	return New(http.StatusUnprocessableEntity, format, a...)
}

func Locked(format string, a ...any) *Error {
	return New(http.StatusLocked, format, a...)
}

func TooEarly(format string, a ...any) *Error {
	return New(http.StatusTooEarly, format, a...)
}

func UpgradeRequired(format string, a ...any) *Error {
	return New(http.StatusUpgradeRequired, format, a...)
}

func PreconditionRequired(format string, a ...any) *Error {
	return New(http.StatusPreconditionRequired, format, a...)
}

func TooManyRequests(format string, a ...any) *Error {
	return New(http.StatusTooManyRequests, format, a...)
}

func RequestHeaderFieldsTooLarge(format string, a ...any) *Error {
	return New(http.StatusRequestHeaderFieldsTooLarge, format, a...)
}

func UnavailableForLegalReasons(format string, a ...any) *Error {
	return New(http.StatusUnavailableForLegalReasons, format, a...)
}

func InternalServerError(format string, a ...any) *Error {
	return New(http.StatusInternalServerError, format, a...)
}

func NotImplemented(format string, a ...any) *Error {
	return New(http.StatusNotImplemented, format, a...)
}

func BadGateway(format string, a ...any) *Error {
	return New(http.StatusBadGateway, format, a...)
}

func ServiceUnavailable(format string, a ...any) *Error {
	return New(http.StatusServiceUnavailable, format, a...)
}

func GatewayTimeout(format string, a ...any) *Error {
	return New(http.StatusGatewayTimeout, format, a...)
}

func HTTPVersionNotSupported(format string, a ...any) *Error {
	return New(http.StatusHTTPVersionNotSupported, format, a...)
}

func VariantAlsoNegotiates(format string, a ...any) *Error {
	return New(http.StatusVariantAlsoNegotiates, format, a...)
}

func InsufficientStorage(format string, a ...any) *Error {
	return New(http.StatusInsufficientStorage, format, a...)
}

func LoopDetected(format string, a ...any) *Error {
	return New(http.StatusLoopDetected, format, a...)
}

func NotExtended(format string, a ...any) *Error {
	return New(http.StatusNotExtended, format, a...)
}

func NetworkAuthenticationRequired(format string, a ...any) *Error {
	return New(http.StatusNetworkAuthenticationRequired, format, a...)
}
