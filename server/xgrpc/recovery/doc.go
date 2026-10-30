// Package recovery 提供 gRPC 服务端的 panic 恢复拦截器。
//
// 该包提供一元与流式两种拦截器，捕获处理器抛出的 panic，可通过选项注册自定义恢复回调；
// 未配置回调时将 panic 值与调用栈包装为 PanicError 返回，避免单个请求 panic 导致进程退出。
package recovery
