package recovery

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc"
)

// RecoveryHandlerFunc 为无上下文的 panic 恢复回调，参数 p 为 recover 捕获到的值，
// 返回的 error 将作为 RPC 的错误结果。
type RecoveryHandlerFunc func(p any) (err error)

// RecoveryHandlerFuncContext 为携带上下文的 panic 恢复回调，可从 ctx 读取请求级元数据；
// 参数 p 为 recover 捕获到的值，返回的 error 将作为 RPC 的错误结果。
type RecoveryHandlerFuncContext func(ctx context.Context, p any) (err error)

// UnaryServerInterceptor 返回一元服务端 panic 恢复拦截器：处理器 panic 被捕获后交由恢复回调
// 转换为错误，未配置回调时返回附带调用栈的 PanicError；无 panic 时行为与直接调用处理器一致。
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, r, o.recoveryHandlerFunc)
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor 返回流式服务端 panic 恢复拦截器：处理器 panic 被捕获后交由恢复回调
// 转换为错误，未配置回调时返回附带调用栈的 PanicError；无 panic 时行为与直接调用处理器一致。
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(stream.Context(), r, o.recoveryHandlerFunc)
			}
		}()

		return handler(srv, stream)
	}
}

// recoverFrom 优先调用自定义恢复回调；未注册回调时收集当前 goroutine 调用栈并包装为 PanicError。
func recoverFrom(ctx context.Context, p any, r RecoveryHandlerFuncContext) error {
	if r != nil {
		return r(ctx, p)
	}
	stack := make([]byte, 64<<10)
	stack = stack[:runtime.Stack(stack, false)]
	return &PanicError{Panic: p, Stack: stack}
}

// PanicError 表示被拦截器捕获的 panic，携带原始 panic 值与对应调用栈，实现 error 接口。
type PanicError struct {
	Panic any
	Stack []byte
}

// Error 返回 panic 值与调用栈组成的错误文本，满足 error 接口约定。
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught: %v\n\n%s", e.Panic, e.Stack)
}
