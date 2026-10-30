package goroutine

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// ParallelWithError 将多个返回 error 的函数包装为并发执行函数，等待全部函数结束；
// 返回首个非 nil 错误，全部成功时返回 nil。
func ParallelWithError(fns ...func() error) func() error {
	return func() error {
		eg := errgroup.Group{}
		for _, fn := range fns {
			eg.Go(fn)
		}

		return eg.Wait()
	}
}

// ParallelWithErrorChan 并发执行 fns，通过返回的带缓冲错误通道逐个报告错误；
// 全部执行完毕后通道自动关闭。单个函数的 panic 被恢复并转换为错误送入通道。
func ParallelWithErrorChan(fns ...func() error) chan error {
	total := len(fns)
	errs := make(chan error, total)

	var wg sync.WaitGroup
	wg.Add(total)

	go func(errs chan error) {
		wg.Wait()
		close(errs)
	}(errs)

	for _, fn := range fns {
		go func(fn func() error, errs chan error) {
			defer wg.Done()
			errs <- try(fn, nil)
		}(fn, errs)
	}

	return errs
}

// RestrictParallelWithErrorChan 以 concurrency 为并发上限并发执行 fns，通过返回的错误通道
// 逐个报告错误，全部执行完毕后通道自动关闭。
//
// concurrency 小于等于 0 时按 1 处理，大于函数个数时按函数个数处理。
func RestrictParallelWithErrorChan(concurrency int, fns ...func() error) chan error {
	total := len(fns)
	if concurrency <= 0 {
		concurrency = 1
	}
	if concurrency > total {
		concurrency = total
	}
	var wg sync.WaitGroup
	errs := make(chan error, total)
	jobs := make(chan func() error, concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(jobs chan func() error, errs chan error) {
			defer wg.Done()
			for fn := range jobs {
				errs <- try(fn, nil)
			}
		}(jobs, errs)
	}
	go func(errs chan error) {
		for _, fn := range fns {
			jobs <- fn
		}
		close(jobs)
		wg.Wait()
		close(errs)
	}(errs)
	return errs
}
