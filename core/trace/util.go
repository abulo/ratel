package trace

import (
	"net"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// ParseFullMethod 解析 gRPC 全量方法名（/package.service/method），返回简短方法名与服务、方法属性。
//
// 格式不符合两段式时仅返回原名，属性切片为 nil。
func ParseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		// 不符合 /service/method 两段格式，无法拆分出服务与方法属性。
		return name, []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}

	return name, attrs
}

// PeerAttr 由 host:port 地址构造对端网络属性；地址格式非法时返回 nil。
//
// 主机名为空时回填 127.0.0.1。
func PeerAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	if len(host) == 0 {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetHostNameKey.String(host),
		semconv.NetPeerPortKey.String(port),
	}
}
