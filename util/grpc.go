package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func GrpcTime(x *timestamppb.Timestamp) time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).Local()
}
