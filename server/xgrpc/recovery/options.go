package recovery

import "context"

var (
	defaultOptions = &options{
		recoveryHandlerFunc: nil,
	}
)

// options 保存 panic 恢复拦截器的内部配置。
type options struct {
	recoveryHandlerFunc RecoveryHandlerFuncContext
}

// evaluateOptions 以默认配置为基础依次套用传入选项，返回最终生效的配置副本。
func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option 是恢复拦截器的配置选项函数，在构造拦截器时修改内部 options。
type Option func(*options)

// WithRecoveryHandler 注册无上下文的 panic 恢复回调；回调返回的错误将作为拦截器的错误结果。
func WithRecoveryHandler(f RecoveryHandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = RecoveryHandlerFuncContext(func(ctx context.Context, p any) error {
			return f(p)
		})
	}
}

// WithRecoveryHandlerContext 注册携带上下文的 panic 恢复回调，可在回调中读取请求级元数据；
// 回调返回的错误将作为拦截器的错误结果。
func WithRecoveryHandlerContext(f RecoveryHandlerFuncContext) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
