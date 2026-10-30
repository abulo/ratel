package xgin

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errMicroDefault = status.Errorf(codes.Internal, createStatusErr(codeMS, "micro default"))
)

// HTTPError 包装处理器返回的业务错误，携带 HTTP 状态码与错误消息，并实现 error 接口。
type HTTPError struct {
	// Code 为对应的 HTTP 状态码。
	Code    int
	// Message 为对外返回的错误消息。
	Message string
}

// NewHTTPError 按状态码构造 HTTPError；未传入自定义消息时，消息取该状态码对应的 StatusText。
func NewHTTPError(code int, msg ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: StatusText(code)}
	if len(msg) > 0 {
		he.Message = msg[0]
	}

	return he
}

// Error 返回错误消息文本，满足 error 接口约定。
func (e HTTPError) Error() string {
	return e.Message
}

// ErrNotFound 为 404 Not Found 状态的预置错误实例。
var ErrNotFound = HTTPError{
	Code:    StatusNotFound,
	Message: "not found",
}

var (
	// ErrGRPCResponseValid 为 gRPC 响应内容校验不通过时返回的内部错误。
	ErrGRPCResponseValid = status.Errorf(codes.Internal, "response valid")
	// ErrGRPCInvokeLen 为 gRPC 调用返回结果数量不符合预期（应为两项）时返回的内部错误。
	ErrGRPCInvokeLen = status.Errorf(codes.Internal, "invoke request without len 2 res")
)
