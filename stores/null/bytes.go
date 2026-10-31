package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// NullBytes 为 JSON null 字面量的字节表示，供各可空类型判断输入是否为 null。
var NullBytes = []byte("null")

// Bytes 是可空的 []byte，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Bytes struct {
	// Bytes 为内部保存的字节切片取值。
	Bytes []byte
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewBytes 使用给定取值与有效性标记创建 Bytes；valid 为 false 时得到空值。
func NewBytes(b []byte, valid bool) Bytes {
	return Bytes{
		Bytes: b,
		Valid: valid,
		Set:   true,
	}
}

// BytesFrom 根据字节切片创建 Bytes；传入 nil 时得到空值，空切片仍视为有效。
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, b != nil)
}

// BytesFromPtr 根据指针创建 Bytes；指针为 nil 时得到空值。
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil {
		return NewBytes(nil, false)
	}
	n := NewBytes(*b, true)
	return n
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (b Bytes) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (b Bytes) IsSet() bool {
	return b.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按字节切片解析，
// 解析失败时返回对应错误。
func (b *Bytes) UnmarshalJSON(data []byte) error {
	b.Set = true

	if bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Bytes = nil
		return nil
	}

	var bv []byte
	if err := json.Unmarshal(data, &bv); err != nil {
		return err
	}

	b.Bytes = bv
	b.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，非空文本覆盖写入内部切片。
func (b *Bytes) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON 实现 json.Marshaler；空值或空切片序列化为 JSON null。
func (b Bytes) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	if len(b.Bytes) == 0 {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return json.Marshal(b.Bytes)
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 nil。
func (b Bytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (b *Bytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
	b.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (b Bytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (b Bytes) IsZero() bool {
	return !b.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (b *Bytes) Scan(value any) error {
	if value == nil {
		b.Bytes, b.Valid, b.Set = nil, false, false
		return nil
	}
	b.Valid, b.Set = true, true
	return convert.ConvertAssign(&b.Bytes, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (b Bytes) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// ValueOrDefault 返回内部取值；空值时返回空字节切片。
func (t Bytes) ValueOrDefault() []byte {
	if !t.Valid {
		return []byte{}
	}
	return t.Bytes
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Bytes) Result() string {
	if !t.Valid {
		return "NULL"
	}
	// PERF-2-19 [P2] 已知类型为 []byte，直接转换字符串，去掉 cast 反射分发。
	return string(t.Bytes)
}
