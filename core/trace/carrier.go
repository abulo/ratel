package trace

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// 编译期断言 MetadataReaderWriter 实现 propagation.TextMapCarrier。
var _ propagation.TextMapCarrier = (*MetadataReaderWriter)(nil)

// MetadataReaderWriter 将 gRPC metadata.MD 适配为 OpenTelemetry 文本传播载体，用于从 gRPC 元数据提取或注入追踪上下文。
type MetadataReaderWriter metadata.MD

// Get 返回指定键的第一个取值；键不存在时返回空字符串。
func (m MetadataReaderWriter) Get(key string) string {
	values := metadata.MD(m).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Set 写入指定键的值，覆盖该键已有的全部取值。
func (m MetadataReaderWriter) Set(key, value string) {
	metadata.MD(m).Set(key, value)
}

// Keys 返回载体中全部键名。
func (m MetadataReaderWriter) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range metadata.MD(m) {
		keys = append(keys, k)
	}
	return keys
}
