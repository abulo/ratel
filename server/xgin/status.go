package xgin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	rstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	any "google.golang.org/protobuf/types/known/anypb"
)

// EmptyMessage 是不携带任何字段的 Protobuf 消息，用于无响应体的 RPC 场景占位，实现 proto.Message 接口。
type EmptyMessage struct{}

// Reset 将消息清空回零值，满足 proto.Message 接口约定。
func (m *EmptyMessage) Reset() { *m = EmptyMessage{} }

// String 返回消息的字符串表示，固定为 "{}"。
func (m *EmptyMessage) String() string { return "{}" }

// ProtoMessage 标记该类型为 Protobuf 消息，满足 proto.Message 接口约定。
func (*EmptyMessage) ProtoMessage() {}

// GRPCProxyMessage 为 HTTP 网关转发 gRPC 响应时使用的包装消息，携带错误码、提示消息与原始 Protobuf 负载。
type GRPCProxyMessage struct {
	Error   int           `protobuf:"varint,1,opt,name=error" json:"error"`
	Message string        `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Data    proto.Message `protobuf:"bytes,3,opt,name=data" json:"data"`
}

// Reset 将消息清空回零值，满足 proto.Message 接口约定。
func (m *GRPCProxyMessage) Reset() { *m = GRPCProxyMessage{} }

// String 返回 Data 负载的 Protobuf JSON 格式字符串。
func (m *GRPCProxyMessage) String() string { return jsonpb.Format(m.Data) }

// ProtoMessage 标记该类型为 Protobuf 消息，满足 proto.Message 接口约定。
func (*GRPCProxyMessage) ProtoMessage() {}

// MarshalJSONPB 将消息序列化为 JSON：Data 先按 Protobuf JSON 编码，再与错误码、消息一并组装输出。
//
// Data 编码失败时返回错误，返回的字节不包含尾部换行。
func (m *GRPCProxyMessage) MarshalJSONPB(jsb *jsonpb.MarshalOptions) ([]byte, error) {
	ss, err := jsonpbMarshaler.Marshal(m.Data)
	if err != nil {
		return []byte{}, err
	}

	msg := struct {
		Error   int             `json:"error"`
		Message string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
	}{
		Error:   m.Error,
		Message: m.Message,
		Data:    json.RawMessage([]byte(ss)),
	}

	return json.Marshal(msg)
}

var (
	jsonpbMarshaler = jsonpb.MarshalOptions{
		EmitUnpopulated: true,
	}
)

// statusErr 在内部包装 rstatus.Status，供微服务错误状态的格式转换使用。
type statusErr struct {
	s *rstatus.Status
}

// Proto 返回内部状态的深拷贝；内部状态为空时返回 nil。
func (e *statusErr) Proto() *rstatus.Status {
	if e.s == nil {
		return nil
	}
	return proto.Clone(e.s).(*rstatus.Status)
}
func init() {
	s, _ := status.FromError(errMicroDefault)
	de, _ := statusFromString(s.Message())
	_ = de.Proto()
}
func statusFromString(s string) (*statusErr, bool) {
	i := strings.Index(s, ":")
	if i == -1 {
		return nil, false
	}
	u64, err := strconv.ParseInt(s[:i], 10, 32)
	if err != nil {
		return nil, false
	}

	return &statusErr{
		&rstatus.Status{
			Code:    int32(u64),
			Message: s[i:],
			Details: []*any.Any{},
		},
	}, true
}
func createStatusErr(code uint32, msg string) string {
	return fmt.Sprintf("%d:%s", code, msg)
}
