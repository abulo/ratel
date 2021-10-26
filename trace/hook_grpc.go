package trace

import (
	"context"
	"strings"
	"time"

	"github.com/abulo/ratel/v1/ecode"
	"github.com/abulo/ratel/v1/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func grpcExtractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

//GRPCUnaryServerInterceptor 监控程序跟踪
func GRPCUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, grpcExtractAID(ctx))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, grpcExtractAID(ctx), code.GetMessage())
	return resp, err
}

//GRPCStreamServerInterceptor 监控程序跟踪
func GRPCStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, grpcExtractAID(ss.Context()))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, grpcExtractAID(ss.Context()), code.GetMessage())
	return err
}
