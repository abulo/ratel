package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/pkg/errors"
)

// Bool 是可空的布尔值，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Bool struct {
	// Bool 为内部保存的布尔取值。
	Bool bool
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewBool 使用给定取值与有效性标记创建 Bool；valid 为 false 时得到空值。
func NewBool(b, valid bool) Bool {
	return Bool{
		Bool:  b,
		Valid: valid,
		Set:   true,
	}
}

// BoolFrom 创建一个始终有效的 Bool。
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

// BoolFromPtr 根据指针创建 Bool；指针为 nil 时得到空值。
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (b Bool) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (b Bool) IsSet() bool {
	return b.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按布尔值解析，
// 解析失败时返回对应错误。
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.Set = true

	if bytes.Equal(data, NullBytes) {
		b.Bool = false
		b.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return err
	}

	b.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，仅接受 "true" 与 "false"，
// 其余文本返回错误。
func (b *Bool) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Valid = false
		return nil
	}

	str := string(text)
	switch str {
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = true
	return nil
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return json.Marshal(nil)
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid 将内部取值置为 v，并同时标记为有效、已赋值。
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
	b.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (b Bool) IsZero() bool {
	return !b.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (b *Bool) Scan(value any) error {
	if value == nil {
		b.Bool, b.Valid, b.Set = false, false, false
		return nil
	}
	b.Valid, b.Set = true, true
	return convert.ConvertAssign(&b.Bool, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (b Bool) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (t Bool) ValueOrDefault() bool {
	if !t.Valid {
		return false
	}
	return t.Bool
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Bool) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return strconv.FormatBool(t.Bool)
}
