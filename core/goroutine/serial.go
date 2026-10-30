package goroutine

import (
	"go.uber.org/multierr"
)

// SerialWithError 将多个返回 error 的函数包装为按声明顺序串行执行的函数；
// 全部执行后聚合收集到的错误返回，无错误时返回 nil。
func SerialWithError(fns ...func() error) func() error {
	return func() error {
		var errs error
		for _, fn := range fns {
			errs = multierr.Append(errs, try(fn, nil))
		}
		return errs
	}
}

// SerialUntilError 将多个函数包装为串行执行函数；遇到第一个非 nil 错误立即停止后续执行，
// 并将该错误返回。
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
			}
		}
		return nil
	}
}

// WhenError 指定串行执行序列遇到错误时采取的处理策略。
type WhenError int

// 串行错误处理策略常量。
var (
	// ReturnWhenError 遇到错误立即返回，中止后续函数执行。
	ReturnWhenError WhenError = 1

	// ContinueWhenError 记录错误后继续执行后续函数，最终聚合全部错误。
	ContinueWhenError WhenError = 2

	// PanicWhenError 遇到错误直接 panic。
	PanicWhenError WhenError = 3

	// LastErrorWhenError 遇到错误继续执行，最终只返回最后一个错误。
	LastErrorWhenError WhenError = 4
)

// SerialWhenError 按错误处理策略 we 构造串行执行器：先接收 fns，返回一个执行函数；
// 执行时按声明顺序运行，遇到错误后的行为由 we 指定的策略决定。
func SerialWhenError(we WhenError) func(fn ...func() error) func() error {
	return func(fns ...func() error) func() error {
		return func() error {
			var errs error
			for _, fn := range fns {
				if err := try(fn, nil); err != nil {
					switch we {
					case ReturnWhenError:
						return err
					case ContinueWhenError:
						errs = multierr.Append(errs, err)
					case PanicWhenError:
						panic(err)
					case LastErrorWhenError:
						errs = err
					}
				}
			}
			return errs
		}
	}
}
