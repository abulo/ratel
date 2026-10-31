package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

// Byte 是可空的 byte，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Byte struct {
	// Byte 为内部保存的字节取值。
	Byte byte
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewByte 使用给定取值与有效性标记创建 Byte；valid 为 false 时得到空值。
func NewByte(b byte, valid bool) Byte {
	return Byte{
		Byte:  b,
		Valid: valid,
		Set:   true,
	}
}

// ByteFrom 创建一个始终有效的 Byte。
func ByteFrom(b byte) Byte {
	return NewByte(b, true)
}

// ByteFromPtr 根据指针创建 Byte；指针为 nil 时得到空值。
func ByteFromPtr(b *byte) Byte {
	if b == nil {
		return NewByte(0, false)
	}
	return NewByte(*b, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (b Byte) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (b Byte) IsSet() bool {
	return b.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为空或 JSON null 时置为空值；其余内容按单字符字符串解析，
// 长度大于 1 或解析失败时返回错误。
func (b *Byte) UnmarshalJSON(data []byte) error {
	b.Set = true

	if len(data) == 0 || bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Byte = 0
		return nil
	}

	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if len(x) > 1 {
		return errors.New("json: cannot convert to byte, text len is greater than one")
	}

	b.Byte = x[0]
	b.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，文本长度大于 1 时返回错误，
// 否则取首字节写入内部取值。
func (b *Byte) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Valid = false
		return nil
	}

	if len(text) > 1 {
		return errors.New("text: cannot convert to byte, text len is greater than one")
	}

	b.Valid = true
	b.Byte = text[0]
	return nil
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null，有效值序列化为带引号的单字符字符串。
func (b Byte) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte{'"', b.Byte, '"'}, nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (b Byte) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return []byte{b.Byte}, nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (b *Byte) SetValid(n byte) {
	b.Byte = n
	b.Valid = true
	b.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (b Byte) Ptr() *byte {
	if !b.Valid {
		return nil
	}
	return &b.Byte
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (b Byte) IsZero() bool {
	return !b.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 或空字符串时清空状态，否则取字符串首字节写入内部取值。
func (b *Byte) Scan(value any) error {
	if value == nil {
		b.Byte, b.Valid, b.Set = 0, false, false
		return nil
	}

	val := value.(string)
	if len(val) == 0 {
		b.Byte, b.Valid, b.Set = 0, false, false
		return nil
	}

	b.Byte, b.Valid, b.Set = val[0], true, true
	return nil
}

// Value 实现 driver.Valuer；空值返回 nil。
func (b Byte) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return []byte{b.Byte}, nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (t Byte) ValueOrDefault() byte {
	if !t.Valid {
		return 0
	}
	return t.Byte
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Byte) Result() string {
	if !t.Valid {
		return "NULL"
	}
	// PERF-2-19 [P2] 已知类型为 byte，直接十进制格式化，去掉 cast 反射分发。
	return strconv.FormatUint(uint64(t.Byte), 10)
}
