package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GrpcTime 将 protobuf Timestamp 转换为本地时区的 time.Time。
//
// x 为 nil 时返回零值时刻对应的本地时间。
func GrpcTime(x *timestamppb.Timestamp) time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).Local()
}
