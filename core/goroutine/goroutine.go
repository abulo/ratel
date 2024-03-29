package goroutine

import (
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/codegangsta/inject"
	"github.com/sirupsen/logrus"
)

// Serial 串行
func Serial(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

// Parallel 并发执行
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

// RestrictParallel 并发,最大并发量restrict
func RestrictParallel(restrict int, fns ...func()) func() {
	var channel = make(chan struct{}, restrict)
	return func() {
		var wg sync.WaitGroup
		for _, fn := range fns {
			wg.Add(1)
			go func(fn func()) {
				defer wg.Done()
				channel <- struct{}{}
				try2(fn, nil)
				<-channel
			}(fn)
		}
		wg.Wait()
		close(channel)
	}
}

// GoDirect ...
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
		// 忽略返回值, goroutine执行的返回值通常都会忽略掉
		_, err := inj.Invoke(fn)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("inject")
			return
		}
	}()
}

// Go goroutine
func Go(fn func()) {
	go try2(fn, nil)
}

// DelayGo goroutine
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

// SafeGo safe go
func SafeGo(fn func(), rec func(error)) {
	go func() {
		err := try2(fn, nil)
		if err != nil {
			rec(err)
		}
	}()
}
