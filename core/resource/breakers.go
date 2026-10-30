package resource

import "sync"

var (
	lock     sync.RWMutex
	breakers = make(map[string]Breaker)
)

// Do 获取名为 name 的共享熔断器并在其保护下执行 req；熔断器拒绝时返回对应错误。
func Do(name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.Do(req)
	})
}

// DoWithAcceptable 行为同 Do，acceptable 用于自定义错误是否视为成功。
func DoWithAcceptable(name string, req func() error, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptable(req, acceptable)
	})
}

// DoWithFallback 行为同 Do，熔断器拒绝时执行 fallback 降级并返回其错误。
func DoWithFallback(name string, req func() error, fallback func(err error) error) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallback(req, fallback)
	})
}

// DoWithFallbackAcceptable 组合 fallback 降级与 acceptable 成功判定，其余行为同 Do。
func DoWithFallbackAcceptable(name string, req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackAcceptable(req, fallback, acceptable)
	})
}

// GetBreaker 返回名为 name 的共享熔断器；不存在时按默认配置创建并缓存，
// 后续以相同名称获取的均为同一实例。
func GetBreaker(name string) Breaker {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return b
	}

	lock.Lock()
	b, ok = breakers[name]
	if !ok {
		b = NewBreaker(WithName(name))
		breakers[name] = b
	}
	lock.Unlock()

	return b
}

// NoBreakerFor 将 name 对应的熔断器替换为空操作熔断器，使该名称下的请求
// 不再经过熔断保护。
func NoBreakerFor(name string) {
	lock.Lock()
	breakers[name] = newNoOpBreaker()
	lock.Unlock()
}

func do(name string, execute func(b Breaker) error) error {
	return execute(GetBreaker(name))
}
