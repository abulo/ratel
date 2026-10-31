package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/pkg/errors"
)

// JSON 是保存原始 JSON 字节的可空 []byte，适用于数据库中可为 null 的 JSON 列。
//
// 它需要区分两层空值：数据库列本身为 null，以及 JSON 内容为 null；序列化时无法把
// JSON 值 null 单独表达，因此以 Valid 标记内容是否有效、以 Set 标记字段是否出现。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type JSON struct {
	// JSON 为内部保存的原始 JSON 字节。
	JSON []byte
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewJSON 使用给定字节与有效性标记创建 JSON；valid 为 false 时得到空值。
func NewJSON(b []byte, valid bool) JSON {
	return JSON{
		JSON:  b,
		Valid: valid,
		Set:   true,
	}
}

// JSONFrom 根据字节切片创建 JSON；传入 nil 时得到空值，空切片仍视为有效。
func JSONFrom(b []byte) JSON {
	return NewJSON(b, b != nil)
}

// JSONFromPtr 根据指针创建 JSON；指针为 nil 时得到空值。
func JSONFromPtr(b *[]byte) JSON {
	if b == nil {
		return NewJSON(nil, false)
	}
	n := NewJSON(*b, true)
	return n
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (j JSON) IsValid() bool {
	return j.Set && j.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (j JSON) IsSet() bool {
	return j.Set
}

// Unmarshal 将内部保存的 JSON 解码到 dest 指向的对象；dest 不能为 nil，
// 返回解码过程中产生的错误。
func (j JSON) Unmarshal(dest any) error {
	if dest == nil {
		return errors.New("destination is nil, not a valid pointer to an object")
	}
	// PERF-2-18 [P2] 直接解码内部原始字节，去掉原先对接收者再 marshal 一轮的冗余往返（预期收益需实测验证）；
	// 空值/空内容保持原语义：解码为 no-op，null 不解构目标对象。
	if !j.Valid || len(j.JSON) == 0 {
		return nil
	}
	return json.Unmarshal(j.JSON, dest)
}

// UnmarshalJSON 实现 json.Unmarshaler。
//
// 以结构体中名为 v 的 JSON 字段为例，三种输入状态分别为：
//
//	{}          未调用本方法，Set 与 Valid 均为 false；
//	{"v": null} 调用本方法，Set 为 true、Valid 为 false；
//	{"v": {}}   调用本方法，Set 与 Valid 均为 true，JSON 内容为 "{}"。
//
// JSON 层面传入 null 时不保留该内容，而是置为值层面的空标记，使数据库落库结果为 SQL NULL。
func (j *JSON) UnmarshalJSON(data []byte) error {
	if data == nil {
		return fmt.Errorf("null: cannot unmarshal nil into Go value of type null.JSON")
	}

	j.Set = true

	if bytes.Equal(data, NullBytes) {
		j.JSON = nil
		j.Valid = false
		return nil
	}

	j.Valid = true
	j.JSON = make([]byte, len(data))
	copy(j.JSON, data)

	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，非空文本覆盖写入内部字节。
func (j *JSON) UnmarshalText(text []byte) error {
	j.Set = true
	if len(text) == 0 {
		j.JSON = nil
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], text...)
		j.Valid = true
	}

	return nil
}

// Marshal 将 obj 编码为 JSON 并保存到接收者内部；返回编码过程中产生的错误。
func (j *JSON) Marshal(obj any) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	// PERF-2-18 [P2] 编码结果直接写入内部字节，去掉原先再 unmarshal 一轮的冗余往返（预期收益需实测验证）；
	// 保持原语义：结果为 JSON null 时 Valid=false，其余内容 Valid=true。
	j.JSON = res
	j.Valid = !bytes.Equal(res, NullBytes)
	j.Set = true
	return nil
}

// MarshalJSON 实现 json.Marshaler；空值或空字节序列化为 JSON null，有效值原样返回内部字节。
func (j JSON) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	if len(j.JSON) == 0 || j.JSON == nil {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return j.JSON, nil
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 nil。
func (j JSON) MarshalText() ([]byte, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (j *JSON) SetValid(n []byte) {
	j.JSON = n
	j.Valid = true
	j.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (j JSON) Ptr() *[]byte {
	if !j.Valid {
		return nil
	}
	return &j.JSON
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (j JSON) IsZero() bool {
	return !j.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (j *JSON) Scan(value any) error {
	if value == nil {
		j.JSON, j.Valid, j.Set = nil, false, false
		return nil
	}
	j.Valid, j.Set = true, true
	return convert.ConvertAssign(&j.JSON, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (j JSON) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// ValueOrDefault 返回内部取值；空值时返回空字节切片。
func (t JSON) ValueOrDefault() []byte {
	if !t.Valid {
		return []byte{}
	}
	return t.JSON
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t JSON) Result() string {
	if !t.Valid {
		return "NULL"
	}
	// PERF-2-19 [P2] 已知类型为 []byte，直接转换字符串，去掉 cast 反射分发。
	return string(t.JSON)
}
