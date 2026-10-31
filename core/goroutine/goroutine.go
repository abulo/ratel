package goroutine

import (
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/codegangsta/inject"
	"github.com/sirupsen/logrus"
)

// Serial 将多个无参函数包装为一个按声明顺序串行执行的函数；函数内部 panic 不做恢复，
// 会中止后续执行并向上抛出。
func Serial(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

// Parallel 将多个无参函数包装为一个并发执行的函数：每个函数在独立 goroutine 中运行，
// 等待全部结束后返回；单个函数的 panic 由内部恢复并记录日志，不影响其余函数。
func Parallel(fns ...func()) func() {
	var wg sync.WaitGroup
	return func() {
		wg.Add(len(fns))
		for _, fn := range fns {
			go try2(fn, wg.Done)
		}
		wg.Wait()
	}
}

// RestrictParallel 将多个函数包装为受并发数限制的并发执行函数，restrict 为同时运行的
// goroutine 上限；函数等待全部执行结束后返回，单个函数的 panic 由内部恢复并记录日志。
func RestrictParallel(restrict int, fns ...func()) func() {
	return func() {
		// PERF-4-9 [P2]：固定 restrict 个 worker 从任务通道取活，替代原先
		// 每个任务先创建一个 goroutine、再阻塞排队抢信号量的做法，限制并发的同时
		// 把在途 goroutine 数量也约束在 worker 规模内。
		if restrict <= 0 {
			restrict = 1
		}
		jobs := make(chan func(), len(fns))
		var wg sync.WaitGroup
		workers := min(restrict, len(fns))
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				for fn := range jobs {
					try2(fn, nil)
				}
			}()
		}
		for _, fn := range fns {
			jobs <- fn
		}
		close(jobs)
		wg.Wait()
	}
}

// GoDirect 在新 goroutine 中通过依赖注入调用 fn，args 作为注入参数；fn 的返回值被忽略。
// 执行过程中的 panic 与注入调用错误都会记录到日志，不会向外抛出。
func GoDirect(fn any, args ...any) {
	var inj = inject.New()
	for _, arg := range args {
		inj.Map(arg)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("recover")
			}
		}()
		// goroutine 内函数的返回值通常由调用方另行处理，这里直接忽略注入调用的返回值。
		_, err := inj.Invoke(fn)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("inject")
			return
		}
	}()
}

// Go 在新 goroutine 中执行 fn；fn 内部 panic 被恢复并记录日志，不会导致进程崩溃。
func Go(fn func()) {
	go try2(fn, nil)
}

// DelayGo 等待 delay 时长后在新 goroutine 中执行 fn；fn 内部 panic 被恢复并记录日志。
func DelayGo(delay time.Duration, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("inject")
			}
		}()
		time.Sleep(delay)
		fn()
	}()
}

// SafeGo 在新 goroutine 中执行 fn；fn 内部 panic 被恢复为错误后交给 rec 回调处理，
// 仅在出现错误时调用 rec，调用方须保证 rec 非 nil。
func SafeGo(fn func(), rec func(error)) {
	go func() {
		err := try2(fn, nil)
		if err != nil {
			rec(err)
		}
	}()
}
