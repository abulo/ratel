// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecode

import (
	"sync"

	"github.com/abulo/ratel/v3/core/logger"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	any "google.golang.org/protobuf/types/known/anypb"
)

// EcodeNum 为系统错误码与业务错误码的分界值；编号小于等于该值视为系统错误码，
// 业务错误码应使用 10000 以上的编号。
const EcodeNum int32 = 9999

var (
	aid              int
	maxCustomizeCode = 9999
	_codes           sync.Map
	// OK 为共享的成功状态，对应 gRPC codes.OK。
	OK = add(int(codes.OK), "OK")
)

// Add 创建并缓存一条错误状态，最终错误码由应用编号前缀与传入 code 拼接而成。
//
// 传入 code 大于 9999 时直接 panic；返回状态可在后续错误判断中复用。
func Add(code int, message string) *spbStatus {
	if code > maxCustomizeCode {
		logger.Logger.Panic("customize code must less than 9999,", code)
	}

	return add(aid*10000+code, message)
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

// ExtractCodes 将 error 转换为内部错误状态；传入 nil 时返回共享的 OK 状态，
// 其余错误取其 gRPC 状态码与消息。
func ExtractCodes(e error) *spbStatus {
	if e == nil {
		return OK
	}
	gst, _ := status.FromError(e)
	// PERF-4-13 [P2]：Details 延迟初始化，错误路径不再预先分配空切片；
	// WithDetails 的 append 与 GetDetailMessage 的 range 对 nil 切片行为等价
	//（错误路径分配形成的 GC 压力预期收益需实测验证）。
	return &spbStatus{
		&spb.Status{
			Code:    int32(gst.Code()),
			Message: gst.Message(),
		},
	}
}
