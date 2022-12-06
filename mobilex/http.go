package mobilex

import (
	"net/http"
)

const (
	StatusContinue           int = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols int = 101 // RFC 7231, 6.2.2
	StatusProcessing         int = 102 // RFC 2518, 10.1
	StatusEarlyHints         int = 103 // RFC 8297

	StatusOK                   int = 200 // RFC 7231, 6.3.1
	StatusCreated              int = 201 // RFC 7231, 6.3.2
	StatusAccepted             int = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo int = 203 // RFC 7231, 6.3.4
	StatusNoContent            int = 204 // RFC 7231, 6.3.5
	StatusResetContent         int = 205 // RFC 7231, 6.3.6
	StatusPartialContent       int = 206 // RFC 7233, 4.1
	StatusMultiStatus          int = 207 // RFC 4918, 11.1
	StatusAlreadyReported      int = 208 // RFC 5842, 7.1
	StatusIMUsed               int = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   int = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  int = 301 // RFC 7231, 6.4.2
	StatusFound             int = 302 // RFC 7231, 6.4.3
	StatusSeeOther          int = 303 // RFC 7231, 6.4.4
	StatusNotModified       int = 304 // RFC 7232, 4.1
	StatusUseProxy          int = 305 // RFC 7231, 6.4.5
	_                       int = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect int = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect int = 308 // RFC 7538, 3

	StatusBadRequest                   int = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 int = 401 // RFC 7235, 3.1
	StatusPaymentRequired              int = 402 // RFC 7231, 6.5.2
	StatusForbidden                    int = 403 // RFC 7231, 6.5.3
	StatusNotFound                     int = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             int = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                int = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            int = 407 // RFC 7235, 3.2
	StatusRequestTimeout               int = 408 // RFC 7231, 6.5.7
	StatusConflict                     int = 409 // RFC 7231, 6.5.8
	StatusGone                         int = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               int = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           int = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        int = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            int = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         int = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable int = 416 // RFC 7233, 4.4
	StatusExpectationFailed            int = 417 // RFC 7231, 6.5.14
	StatusTeapot                       int = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           int = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          int = 422 // RFC 4918, 11.2
	StatusLocked                       int = 423 // RFC 4918, 11.3
	StatusFailedDependency             int = 424 // RFC 4918, 11.4
	StatusTooEarly                     int = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              int = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         int = 428 // RFC 6585, 3
	StatusTooManyRequests              int = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  int = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   int = 451 // RFC 7725, 3

	StatusInternalServerError           int = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                int = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    int = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            int = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                int = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       int = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         int = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           int = 507 // RFC 4918, 11.5
	StatusLoopDetected                  int = 508 // RFC 5842, 7.2
	StatusNotExtended                   int = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired int = 511 // RFC 6585, 6

	// ---------------------------------------------------------
	// Custom status codes

	StatusTransportFailed int = 600
	StatusInvalidResponse int = 601

	// ----------------------------------------------------------
)

func StatusText(code int) string {
	switch code {
	case StatusTransportFailed:
		return "transport failure"
	case StatusInvalidResponse:
		return "invalid response"
	default:
		return http.StatusText(code)
	}
}
