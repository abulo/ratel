package trace

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// assert that MetadataReaderWriter implements the TextMapCarrier interface
var _ propagation.TextMapCarrier = (*MetadataReaderWriter)(nil)

// MetadataReaderWriter ...
type MetadataReaderWriter metadata.MD

func (m MetadataReaderWriter) Get(key string) string {
	values := metadata.MD(m).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (m MetadataReaderWriter) Set(key, value string) {
	metadata.MD(m).Set(key, value)
}

func (m MetadataReaderWriter) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range metadata.MD(m) {
		keys = append(keys, k)
	}
	return keys
}
