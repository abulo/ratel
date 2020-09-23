package ecode

import (
	"sync"

	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_codes sync.Map
	// OK ...
	OK = add(int(codes.OK), "OK")
)

type spbStatus struct {
	*spb.Status
}

// ExtractCodes cause from error to ecode.
func ExtractCodes(e error) *spbStatus {
	if e == nil {
		return OK
	}
	// todo 不想做code类型转换，所以全部用grpc标准码处理
	// 如果存在标准的grpc的错误，直接返回自定义的ecode编码
	gst, _ := status.FromError(e)
	return &spbStatus{
		&spb.Status{
			Code:    int32(gst.Code()),
			Message: gst.Message(),
			Details: make([]*any.Any, 0),
		},
	}
}

func add(code int, message string) *spbStatus {
	status := &spbStatus{
		&spb.Status{
			Code:    int32(code),
			Message: message,
			Details: make([]*any.Any, 0),
		},
	}
	_codes.Store(code, status)
	return status
}
