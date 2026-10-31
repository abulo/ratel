package ecode

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Status 是错误状态的通用占位类型，实际载体为 *spbStatus。
type Status any

type spbStatus struct {
	*spb.Status
}

// GetCodeAsInt 返回状态错误码的 int 表示。
func (s *spbStatus) GetCodeAsInt() int {
	return int(s.Code)
}

// GetCodeAsUint32 返回状态错误码的 uint32 表示。
func (s *spbStatus) GetCodeAsUint32() uint32 {
	return uint32(s.Code)
}

// GetCodeAsBool 返回当前状态是否为成功状态：剥离应用编号前缀后的错误码为 0 时视为成功。
func (s *spbStatus) GetCodeAsBool() bool {
	return s.CauseCode() == 0
}

// GetMessage 返回状态消息文本；exts 须按键、值成对传入，追加为 ",键:值" 片段，
// 参数个数为奇数时直接 panic。
func (s *spbStatus) GetMessage(exts ...any) string {
	if len(exts)%2 != 0 {
		panic("parameter must be odd")
	}

	var buf bytes.Buffer
	buf.WriteString(s.Message)

	if len(exts) > 0 {
		buf.WriteByte(',')
	}
	for i := 0; i < len(exts); i++ {
		buf.WriteString(messageAnyToString(exts[i]))
		buf.WriteByte(':')
		buf.WriteString(messageAnyToString(exts[i+1]))
		i++
	}
	return buf.String()
}

// messageAnyToString 对常见标量类型走无反射快路径，其余类型回退 fmt.Sprintf("%v")，
// 输出格式与原实现保持一致（PERF-4-15 [P2]）。
func messageAnyToString(v any) string {
	switch t := v.(type) {
	case nil:
		return "<nil>"
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float64:
		// 与 fmt.Sprintf("%v") 的 %g 语义保持一致，避免大指数场景被固定小数表示改写。
		return strconv.FormatFloat(t, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 32)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetDetailMessage 在 GetMessage 的结果之后逐行追加状态携带的每条详情消息。
func (s *spbStatus) GetDetailMessage(exts ...any) string {
	var buf bytes.Buffer
	buf.WriteString(s.GetMessage(exts...))
	for _, detail := range s.Details {
		buf.WriteByte('\n')
		buf.WriteString(detail.String())
	}
	return buf.String()
}

// String 返回状态的文本表示；PERF-4-14 [P2]：采用 code/message 常见字段手写拼接，
// 替代整结构 json.Marshal 反射序列化（Details 等其余字段不在保守子集内；预期收益需实测验证）。
func (s *spbStatus) String() string {
	if s == nil {
		return "null"
	}
	code, message := int32(0), ""
	if s.Status != nil {
		code = s.Code
		message = s.Message
	}
	return fmt.Sprintf(`{"code":%d,"message":%q}`, code, message)
}

// CauseCode 返回剥离应用编号前缀后的原始错误码（错误码对 10000 取模）。
func (s *spbStatus) CauseCode() int {
	return int(s.Code) % 10000
}

// Proto 返回内部 gRPC Status 的深拷贝；接收者为 nil 时返回 nil。
func (s *spbStatus) Proto() *spb.Status {
	if s == nil {
		return nil
	}
	return proto.Clone(s.Status).(*spb.Status)
}

// MustWithDetails 在当前状态上追加详情，失败时直接 panic；成功返回新的状态实例。
func (s *spbStatus) MustWithDetails(details ...any) *spbStatus {
	status, err := s.WithDetails(details...)
	if err != nil {
		panic(err)
	}
	return status
}

// WithDetails 克隆当前状态并追加详情消息，返回新的状态实例；原状态不被修改。
//
// 成功状态（错误码为 0）不允许附加详情，此时返回错误；details 为 proto.Message 时
// 按 Protobuf 编码打包，其余值按反射类型名与格式化文本打包，打包失败时返回首个错误。
func (s *spbStatus) WithDetails(details ...any) (*spbStatus, error) {
	if s.CauseCode() == 0 {
		return nil, errors.New("no error details for status with code OK")
	}
	p := s.Proto()
	for _, detail := range details {
		if pmsg, ok := detail.(proto.Message); ok {
			anypb, err := marshalAnyProtoMessage(pmsg)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, anypb)
		} else {
			anypb, err := marshalAny(detail)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, anypb)
		}
	}
	return &spbStatus{Status: p}, nil
}

func marshalAny(obj any) (*anypb.Any, error) {
	typ := reflect.TypeOf(obj)
	val := fmt.Sprintf("%+v", obj)

	return &anypb.Any{TypeUrl: typ.Name(), Value: []byte(val)}, nil
}

func marshalAnyProtoMessage(pb proto.Message) (*anypb.Any, error) {
	value, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return &anypb.Any{TypeUrl: string(pb.ProtoReflect().Descriptor().FullName()), Value: value}, nil
}
