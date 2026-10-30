package xgin

const (
	// codeMS 为微服务自定义错误码的基值。
	codeMS = 1000
)

// HTTP 状态码常量，取值与标准库 net/http 保持一致，供 HTTPError 与响应状态判断使用。
const (
	StatusContinue = 100
	StatusSwitchingProtocols = 101

	StatusOK = 200
	StatusCreated = 201
	StatusAccepted = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent = 204
	StatusResetContent = 205
	StatusPartialContent = 206

	StatusMultipleChoices = 300
	StatusMovedPermanently = 301
	StatusFound = 302
	StatusSeeOther = 303
	StatusNotModified = 304
	StatusUseProxy = 305
	StatusTemporaryRedirect = 307

	StatusBadRequest = 400
	StatusUnauthorized = 401
	StatusPaymentRequired = 402
	StatusForbidden = 403
	StatusNotFound = 404
	StatusMethodNotAllowed = 405
	StatusNotAcceptable = 406
	StatusProxyAuthRequired = 407
	StatusRequestTimeout = 408
	StatusConflict = 409
	StatusGone = 410
	StatusLengthRequired = 411
	StatusPreconditionFailed = 412
	StatusRequestEntityTooLarge = 413
	StatusRequestURITooLong = 414
	StatusUnsupportedMediaType = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed = 417
	StatusTeapot = 418
	StatusPreconditionRequired = 428
	StatusTooManyRequests = 429
	StatusRequestHeaderFieldsTooLarge = 431
	StatusUnavailableForLegalReasons = 451
	StatusInternalServerError = 500
	StatusNotImplemented = 501
	StatusBadGateway = 502
	StatusServiceUnavailable = 503
	StatusGatewayTimeout = 504
	StatusHTTPVersionNotSupported = 505
	StatusNetworkAuthenticationRequired = 511

	// StatusErrorCodeReturned 为微服务自定义错误返回使用的状态码。
	StatusErrorCodeReturned = 800
)

// StatusText 返回 HTTP 状态码对应的英文描述文本；状态码未登记时返回空字符串。
func StatusText(code int) string {
	return statusText[code]
}

var statusText = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Errord",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// HTTP 头部字段名称常量。
const (
	// HeaderAcceptEncoding 为客户端可接受内容编码的请求头。
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderContentType 为请求/响应内容类型的头部字段。
	HeaderContentType = "Content-Type"
	// HeaderHRPCErr 为微服务框架用于透传错误信息的自定义响应头。
	HeaderHRPCErr = "HRPC-Errord"
)

// MIME 内容类型常量。
const (
	// MIMEApplicationJSON 为 JSON 内容类型。
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 为带 UTF-8 字符集声明的 JSON 内容类型。
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	// MIMEApplicationProtobuf 为 Protobuf 内容类型。
	MIMEApplicationProtobuf = "application/protobuf"
)
const (
	charsetUTF8 = "charset=utf-8"
)
