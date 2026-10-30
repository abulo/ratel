// Package trace 封装 OpenTelemetry 链路追踪的初始化、上下文传播与 gRPC 元数据载体。
//
// 包内支持 HTTP/gRPC 两种 OTLP 导出协议，同时兼容 Jaeger 与 W3C TraceContext
// 传播格式，并提供面向服务端、客户端、生产者、消费者角色的 Tracer 包装。
package trace
