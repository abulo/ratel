// Package singleton 提供按模块与键分组的进程内单例存储，供客户端、注册中心等组件
// 在初始化后复用同一实例；存储基于 sync.Map，读写操作并发安全。
package singleton
