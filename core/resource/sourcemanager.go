package resource

import (
	"io"
	"sync"
)

// ResourceManager 按名称管理一组 io.Closer 资源，并通过单飞合并并发的重复创建请求。
//
// 方法为并发安全；Close 之后不可再调用 GetResource 或 Inject。
type ResourceManager struct {
	resources    map[string]io.Closer
	singleFlight SingleFlight
	lock         sync.RWMutex
}

// NewResourceManager 创建一个空的资源管理器。
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources:    make(map[string]io.Closer),
		singleFlight: NewSingleFlight(),
	}
}

// Close 依次关闭全部已注册资源，聚合并返回所有关闭过程中产生的错误；
// 任何单个资源关闭失败都不会中断其余资源的关闭。
//
// Close 返回后资源表被清空，调用方不得再次使用本管理器。
func (manager *ResourceManager) Close() error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var be BatchError
	for _, resource := range manager.resources {
		if err := resource.Close(); err != nil {
			be.Add(err)
		}
	}

	// 清空资源表，防止 Close 完成后被误用已关闭的资源。
	manager.resources = nil

	return be.Err()
}

// GetResource 返回 key 对应的资源；资源不存在时调用 create 创建并缓存。
//
// 同一 key 的并发创建请求由单飞合并，仅执行一次 create；create 返回错误时
// 不缓存该资源，错误原样返回给调用方。
func (manager *ResourceManager) GetResource(key string, create func() (io.Closer, error)) (io.Closer, error) {
	val, err := manager.singleFlight.Do(key, func() (any, error) {
		manager.lock.RLock()
		resource, ok := manager.resources[key]
		manager.lock.RUnlock()
		if ok {
			return resource, nil
		}

		resource, err := create()
		if err != nil {
			return nil, err
		}

		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.resources[key] = resource

		return resource, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(io.Closer), nil
}

// Inject 将 resource 直接注册到 key 名下，已存在的同名资源会被覆盖。
func (manager *ResourceManager) Inject(key string, resource io.Closer) {
	manager.lock.Lock()
	manager.resources[key] = resource
	manager.lock.Unlock()
}
