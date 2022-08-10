package xgin

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// errBadRequest     = status.Errorf(codes.InvalidArgument, createStatusErr(codeMSInvalidParam, "bad request"))
	errMicroDefault = status.Errorf(codes.Internal, createStatusErr(codeMS, "micro default"))
	// errMicroInvoke  = status.Errorf(codes.Internal, createStatusErr(codeMSInvoke, "invoke failed"))
	// errMicroInvokeLen = status.Errorf(codes.Internal, createStatusErr(codeMSInvokeLen, "invoke result not 2 item"))
	// errMicroInvokeInvalid = status.Errorf(codes.Internal, createStatusErr(codeMSSecondItemNotError, "second invoke res not a error"))
	// errMicroResInvalid = status.Errorf(codes.Internal, createStatusErr(codeMSResErr, "response is not valid"))
)

// HTTPError wraps handler error.
type HTTPError struct {
	Code    int
	Message string
}

// NewHTTPError constructs a new HTTPError instance.
func NewHTTPError(code int, msg ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: StatusText(code)}
	if len(msg) > 0 {
		he.Message = msg[0]
	}

	return he
}

// Errord return error message.
func (e HTTPError) Error() string {
	return e.Message
}

// ErrNotFound defines StatusNotFound error.
var ErrNotFound = HTTPError{
	Code:    StatusNotFound,
	Message: "not found",
}

// ErrGRPCResponseValid ...
var (
	// ErrGRPCResponseValid ...
	ErrGRPCResponseValid = status.Errorf(codes.Internal, "response valid")
	// ErrGRPCInvokeLen ...
	ErrGRPCInvokeLen = status.Errorf(codes.Internal, "invoke request without len 2 res")
)
