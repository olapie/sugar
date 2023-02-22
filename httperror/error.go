package httperror

import (
	"code.olapie.com/sugar/v2/xerror"
	"io"
	"log"
	"net/http"
	"strings"
)

func ParseResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !isText(contentType) {
		return xerror.New(resp.StatusCode, resp.Status)
	}

	body, ioErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if ioErr != nil {
		log.Printf("failed reading response body: %v\n", ioErr)
	}

	if strings.HasPrefix(contentType, "application/json") {
		var respError xerror.Error
		if err := respError.UnmarshalJSON(body); err != nil {
			log.Printf("unmarshal json body: %v\n", err)
		} else {
			return &respError
		}
	}

	code := resp.StatusCode
	message := string(body)
	if message == "" {
		message = resp.Status
	}

	return xerror.New(code, message)
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

func BadRequest(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusBadRequest, format, a...)
}

func Unauthorized(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusUnauthorized, format, a...)
}

func PaymentRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusPaymentRequired, format, a...)
}

func Forbidden(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusForbidden, format, a...)
}

func NotFound(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusNotFound, format, a...)
}

func MethodNotAllowed(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusMethodNotAllowed, format, a...)
}

func NotAcceptable(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusNotAcceptable, format, a...)
}

func ProxyAuthRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusProxyAuthRequired, format, a...)
}

func RequestTimeout(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusRequestTimeout, format, a...)
}

func Conflict(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusConflict, format, a...)
}

func LengthRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusLengthRequired, format, a...)
}

func PreconditionFailed(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusPreconditionFailed, format, a...)
}

func RequestEntityTooLarge(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusRequestEntityTooLarge, format, a...)
}

func RequestURITooLong(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusRequestURITooLong, format, a...)
}

func ExpectationFailed(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusExpectationFailed, format, a...)
}

func Teapot(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusTeapot, format, a...)
}

func MisdirectedRequest(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusMisdirectedRequest, format, a...)
}

func UnprocessableEntity(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusUnprocessableEntity, format, a...)
}

func Locked(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusLocked, format, a...)
}

func TooEarly(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusTooEarly, format, a...)
}

func UpgradeRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusUpgradeRequired, format, a...)
}

func PreconditionRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusPreconditionRequired, format, a...)
}

func TooManyRequests(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusTooManyRequests, format, a...)
}

func RequestHeaderFieldsTooLarge(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusRequestHeaderFieldsTooLarge, format, a...)
}

func UnavailableForLegalReasons(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusUnavailableForLegalReasons, format, a...)
}

func InternalServerError(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusInternalServerError, format, a...)
}

func NotImplemented(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusNotImplemented, format, a...)
}

func BadGateway(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusBadGateway, format, a...)
}

func ServiceUnavailable(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusServiceUnavailable, format, a...)
}

func GatewayTimeout(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusGatewayTimeout, format, a...)
}

func HTTPVersionNotSupported(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusHTTPVersionNotSupported, format, a...)
}

func VariantAlsoNegotiates(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusVariantAlsoNegotiates, format, a...)
}

func InsufficientStorage(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusInsufficientStorage, format, a...)
}

func LoopDetected(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusLoopDetected, format, a...)
}

func NotExtended(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusNotExtended, format, a...)
}

func NetworkAuthenticationRequired(format string, a ...any) *xerror.Error {
	return xerror.New(http.StatusNetworkAuthenticationRequired, format, a...)
}
